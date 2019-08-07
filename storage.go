package myfdbstorage

import (
	"encoding/json"
	"io/ioutil"
	"sync"
	"time"
)

// Storage - data storage
type Storage struct {
	I  map[string]IStorage
	S  map[string]SStorage
	Mx sync.RWMutex
}

// IGetTable - get IStorage by name
func (s *Storage) IGetTable(name string) (st IStorage, ok bool, err error) {
	s.Mx.RLock()
	st, ok = s.I[name]
	s.Mx.RUnlock()
	return st, ok, nil
}

// SGetTable - get SStorage by name
func (s *Storage) SGetTable(name string) (st SStorage, ok bool, err error) {
	s.Mx.RLock()
	st, ok = s.S[name]
	s.Mx.RUnlock()
	return st, ok, nil
}

// StorageLoad - load storage from json params
func StorageLoad(js []byte) (*Storage, error) {
	s := &Storage{
		I: make(map[string]IStorage),
		S: make(map[string]SStorage),
	}

	ss := StorageSettings{}

	err := json.Unmarshal(js, &ss)

	if err != nil {
		return nil, ErrorNew("StorageLoad.json.Unmarshal", err)
	}

	for _, v := range ss.I {
		st, err := IStorageCreate(v.Type, v.Params)
		if err != nil {
			return nil, ErrorNew("StorageLoad load table "+v.Name, err)
		}
		s.I[v.Name] = st
	}

	for _, v := range ss.S {
		st, err := SStorageCreate(v.Type, v.Params)
		if err != nil {
			return nil, ErrorNew("StorageLoad load table "+v.Name, err)
		}
		s.S[v.Name] = st
	}

	return s, nil
}

// StorageStructCollect - load storage from json params
func (s *Storage) StorageStructCollect() (StorageSettings, error) {

	ss := StorageSettings{
		I: make([]StorageTableSettings, 0),
		S: make([]StorageTableSettings, 0),
	}

	for n, v := range s.I {
		ss.I = append(ss.I, StorageTableSettings{
			Name:   n,
			Type:   v.Type(),
			Params: v.Params(),
		})
	}

	for n, v := range s.S {
		ss.S = append(ss.S, StorageTableSettings{
			Name:   n,
			Type:   v.Type(),
			Params: v.Params(),
		})
	}

	return ss, nil
}

// StorageStructFlush flush storage struct to file
func (s *Storage) StorageStructFlush(file string) error {
	ss, err := s.StorageStructCollect()
	if err != nil {
		return ErrorNew("StorageStructFlush get StorageSettings", err)
	}

	js, err := json.Marshal(ss)
	if err != nil {
		return ErrorNew("StorageStructFlush get json from struct", err)
	}

	err = ioutil.WriteFile(file, js, 0660)
	if err != nil {
		return ErrorNew("StorageStructFlush write file", err)
	}

	return nil
}

// StorageSettings params for storage save
type StorageSettings struct {
	I []StorageTableSettings
	S []StorageTableSettings
}

// StorageTableSettings params for storage table save
type StorageTableSettings struct {
	Name   string
	Type   string
	Params map[string]interface{}
}

// AddTableI - add table
func (s *Storage) AddTableI(name string, typeTable string, params map[string]interface{}) (err error) {
	s.Mx.Lock()
	_, ok := s.I[name]
	if ok {
		s.Mx.Unlock()
		return ErrorNew("AddTableI table with name "+name+" already exists", nil)
	}

	st, err := IStorageCreate(typeTable, params)
	if err != nil {
		s.Mx.Unlock()
		return ErrorNew("AddTableI load table "+name, err)
	}

	s.I[name] = st

	s.Mx.Unlock()
	return nil
}

// AddTableS - add table
func (s *Storage) AddTableS(name string, typeTable string, params map[string]interface{}) (err error) {
	s.Mx.Lock()
	_, ok := s.S[name]
	if ok {
		s.Mx.Unlock()
		return ErrorNew("AddTableS table with name "+name+" already exists", nil)
	}

	st, err := SStorageCreate(typeTable, params)
	if err != nil {
		s.Mx.Unlock()
		return ErrorNew("AddTableS load table "+name, err)
	}

	s.S[name] = st

	s.Mx.Unlock()
	return nil
}

// DropTableI - drop table
func (s *Storage) DropTableI(name string) (ok bool, err error) {
	s.Mx.Lock()
	st, ok := s.I[name]
	if !ok {
		s.Mx.Unlock()
		return false, nil
	}

	err = st.ClearAndDeleteStorage()
	if err != nil {
		s.Mx.Unlock()
		return false, ErrorNew("DropTableI ClearAndDeleteStorage table "+name, err)
	}

	delete(s.I, name)

	s.Mx.Unlock()
	return true, nil
}

// DropTableS - drop table
func (s *Storage) DropTableS(name string) (ok bool, err error) {
	s.Mx.Lock()
	st, ok := s.S[name]
	if !ok {
		s.Mx.Unlock()
		return false, nil
	}

	err = st.ClearAndDeleteStorage()
	if err != nil {
		s.Mx.Unlock()
		return false, ErrorNew("DropTableS ClearAndDeleteStorage table "+name, err)
	}

	delete(s.S, name)

	s.Mx.Unlock()
	return true, nil
}

// DetachTableI - detach table without remove data
func (s *Storage) DetachTableI(name string) (ok bool, err error) {
	s.Mx.Lock()
	_, ok = s.I[name]
	if !ok {
		s.Mx.Unlock()
		return false, nil
	}

	delete(s.I, name)

	s.Mx.Unlock()
	return true, nil
}

// DetachTableS - detach table without remove data
func (s *Storage) DetachTableS(name string) (ok bool, err error) {
	s.Mx.Lock()
	_, ok = s.S[name]
	if !ok {
		s.Mx.Unlock()
		return false, nil
	}

	delete(s.S, name)

	s.Mx.Unlock()
	return true, nil
}

// TryGetString try get string value from params
func TryGetString(params map[string]interface{}, name string, val *string) (ok bool) {
	if params == nil {
		return false
	}
	v, ok := params[name]

	if !ok {
		return false
	}

	r, ok := v.(string)

	if !ok {
		return false
	}

	*val = r

	return true
}

// TryGetInt try get int value from params
func TryGetInt(params map[string]interface{}, name string, val *int) (ok bool) {
	if params == nil {
		return false
	}
	v, ok := params[name]

	if !ok {
		return false
	}

	r, ok := v.(int)

	if !ok {

		f, ok := v.(float64)

		if !ok {
			return false
		}

		r = int(f)
	}

	*val = r

	return true
}

// TryGetDuration try get time.Duration value from params
func TryGetDuration(params map[string]interface{}, name string, val *time.Duration) (ok bool) {
	if params == nil {
		return false
	}
	v, ok := params[name]

	if !ok {
		return false
	}

	r, ok := v.(time.Duration)

	if !ok {

		f, ok := v.(float64)

		if !ok {
			i64, ok := v.(int64)
			if !ok {
				i, ok := v.(int)
				if !ok {
					return false
				}
				r = time.Duration(int64(i))
			} else {
				r = time.Duration(i64)
			}
		} else {
			r = time.Duration(int64(f))
		}

	}

	*val = r

	return true
}
