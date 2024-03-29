package myfdbstorage

import (
	"sync"
	"time"
)

// IStorageSimpleName type
const IStorageSimpleName = "simple"

// IItemSimple - simple kv item key int64
type IItemSimple struct {
	Key   int64
	BData []byte
}

// ID - uniqe object identity
func (i IItemSimple) ID() int64 {
	return i.Key
}

// Data - object data
func (i IItemSimple) Data() []byte {
	return i.BData
}

// FullData - all data without ID neaded for restore object
func (i IItemSimple) FullData() []byte {
	return i.BData
}

// IItemSimpleGenerate generate IItemSimple from any IItem contains all data for restore
func IItemSimpleGenerate(item IItem) IItemSimple {

	t, ok := item.(IItemSimple)
	if ok {
		return t
	}

	res := IItemSimple{Key: item.ID(), BData: item.FullData()}
	return res
}

// IItemSimpleConvertTo generate IItemSimple from any IItem contains only Data()
func IItemSimpleConvertTo(item IItem) IItemSimple {

	t, ok := item.(IItemSimple)
	if ok {
		return t
	}

	res := IItemSimple{Key: item.ID(), BData: item.Data()}
	return res
}

// IStorageSimple - simple kv storage key int64
type IStorageSimple struct {
	Data map[int64]IItem
	mx   sync.RWMutex

	DumpPath      string
	DumpTimeout   time.Duration
	DumpStop      bool
	ExistsChanges bool
	mxDump        sync.Mutex
}

// IStorageSimpleLoad load or create simple storage
func IStorageSimpleLoad(params map[string]interface{}) (*IStorageSimple, error) {
	s := &IStorageSimple{Data: make(map[int64]IItem)}

	TryGetString(params, "dump_path", &(s.DumpPath))
	TryGetDuration(params, "dump_timeout", &(s.DumpTimeout))

	if s.DumpPath != "" {
		// load from disk
		err := s.LoadFromDisk()
		if err != nil {
			return s, err
		}
	}

	if s.DumpTimeout > 0 {
		// run auto flush
		go func() {
			for !s.DumpStop {
				time.Sleep(s.DumpTimeout)
				if s.ExistsChanges {
					s.ExistsChanges = false
					err := s.Flush()
					if err != nil {
						Exception(err)
					}
				}
			}
		}()

	}

	return s, nil
}

// Exists item
func (s *IStorageSimple) Exists(key int64) (ok bool, err error) {
	s.mx.RLock()
	_, ok = s.Data[key]
	s.mx.RUnlock()
	return ok, nil
}

// Get item
func (s *IStorageSimple) Get(key int64) (itm IItem, ok bool, err error) {
	s.mx.RLock()
	itm, ok = s.Data[key]
	s.mx.RUnlock()
	return itm, ok, nil
}

// Set - Insert or update item
func (s *IStorageSimple) Set(itm IItem) (err error) {
	s.mx.Lock()
	s.ExistsChanges = true
	s.Data[itm.ID()] = itm
	s.mx.Unlock()
	return nil
}

// Del - delete item
func (s *IStorageSimple) Del(key int64) (exists bool, err error) {
	s.mx.Lock()
	s.ExistsChanges = true
	_, exists = s.Data[key]
	if exists {
		delete(s.Data, key)
	}
	s.mx.Unlock()
	return exists, nil
}

// Type Storage for restore
func (s *IStorageSimple) Type() string {
	return IStorageSimpleName
}

// Params Storage for restore
func (s *IStorageSimple) Params() map[string]interface{} {
	params := make(map[string]interface{})

	params["dump_path"] = s.DumpPath
	params["dump_timeout"] = int64(s.DumpTimeout)

	return params
}

// Flush changes into storages
func (s *IStorageSimple) Flush() error {
	if s.DumpPath != "" {
		err := s.FlushOnDisk()
		return err
	}

	return nil
}

// ClearAndDeleteStorage table clear and then delete from storage
func (s *IStorageSimple) ClearAndDeleteStorage() error {
	if s.DumpPath != "" {
		s.DumpStop = true
		err := s.ClearFromDisk()
		return err
		// load data
	}

	return nil
}
