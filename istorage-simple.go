package myfdbstorage

import "sync"

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

// IStorageSimple - simple kv storage key int64
type IStorageSimple struct {
	Data map[int64]IItem
	mx   sync.RWMutex
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
	s.Data[itm.ID()] = itm
	s.mx.Unlock()
	return nil
}

// Del - delete item
func (s *IStorageSimple) Del(key int64) (exists bool, err error) {
	s.mx.Lock()
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
	return make(map[string]interface{})
}

// ClearAndDeleteStorage table clear and then delete from storage
func (s *IStorageSimple) ClearAndDeleteStorage() error {
	return nil
}
