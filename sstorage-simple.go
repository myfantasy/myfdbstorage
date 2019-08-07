package myfdbstorage

import (
	"sync"
	"time"
)

// SStorageSimpleName type
const SStorageSimpleName = "simple"

// SItemSimple - simple kv item key string
type SItemSimple struct {
	Key   string
	BData []byte
}

// ID - uniqe object identity
func (i SItemSimple) ID() string {
	return i.Key
}

// Data - object data
func (i SItemSimple) Data() []byte {
	return i.BData
}

// FullData - all data without ID neaded for restore object
func (i SItemSimple) FullData() []byte {
	return i.BData
}

// SItemSimpleGenerate generate SItemSimple from any SItem contains all data for restore
func SItemSimpleGenerate(item SItem) SItemSimple {

	t, ok := item.(SItemSimple)
	if ok {
		return t
	}

	res := SItemSimple{Key: item.ID(), BData: item.FullData()}
	return res
}

// SItemSimpleConvertTo generate SItemSimple from any SItem contains only Data()
func SItemSimpleConvertTo(item SItem) SItemSimple {

	t, ok := item.(SItemSimple)
	if ok {
		return t
	}

	res := SItemSimple{Key: item.ID(), BData: item.Data()}
	return res
}

// SStorageSimple - simple kv storage key string
type SStorageSimple struct {
	Data map[string]SItem
	mx   sync.RWMutex

	DumpPath      string
	DumpTimeout   time.Duration
	DumpStop      bool
	ExistsChanges bool
	mxDump        sync.Mutex
}

// SStorageSimpleLoad load or create simple storage
func SStorageSimpleLoad(params map[string]interface{}) (*SStorageSimple, error) {
	s := &SStorageSimple{Data: make(map[string]SItem)}

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
func (s *SStorageSimple) Exists(key string) (ok bool, err error) {
	s.mx.RLock()
	_, ok = s.Data[key]
	s.mx.RUnlock()
	return ok, nil
}

// Get item
func (s *SStorageSimple) Get(key string) (itm SItem, ok bool, err error) {
	s.mx.RLock()
	itm, ok = s.Data[key]
	s.mx.RUnlock()
	return itm, ok, nil
}

// Set - Insert or update item
func (s *SStorageSimple) Set(itm SItem) (err error) {
	s.mx.Lock()
	s.ExistsChanges = true
	s.Data[itm.ID()] = itm
	s.mx.Unlock()
	return nil
}

// Del - delete item
func (s *SStorageSimple) Del(key string) (exists bool, err error) {
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
func (s *SStorageSimple) Type() string {
	return SStorageSimpleName
}

// Params Storage for restore
func (s *SStorageSimple) Params() map[string]interface{} {
	params := make(map[string]interface{})

	params["dump_path"] = s.DumpPath
	params["dump_timeout"] = int64(s.DumpTimeout)

	return params
}

// Flush changes into storages
func (s *SStorageSimple) Flush() error {
	if s.DumpPath != "" {
		err := s.FlushOnDisk()
		return err
	}

	return nil
}

// ClearAndDeleteStorage table clear and then delete from storage
func (s *SStorageSimple) ClearAndDeleteStorage() error {
	if s.DumpPath != "" {
		s.DumpStop = true
		err := s.ClearFromDisk()
		return err
		// load data
	}

	return nil
}
