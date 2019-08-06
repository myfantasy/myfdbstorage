package myfdbstorage

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

// Storage - data storage
type Storage struct {
	I  map[string]IStorage
	Mx sync.RWMutex
}

// IGetTable - get IStorage by name
func (s *Storage) IGetTable(name string) (st IStorage, ok bool, err error) {
	s.Mx.RLock()
	st, ok = s.I[name]
	s.Mx.RUnlock()
	return st, ok, nil
}

// StorageLoad - load storage from json params
func StorageLoad(js []byte) (*Storage, error) {
	s := &Storage{
		I: make(map[string]IStorage),
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

	return s, nil
}

// StorageStructCollect - load storage from json params
func (s *Storage) StorageStructCollect() (StorageSettings, error) {

	ss := StorageSettings{I: make([]StorageTableSettings, 0)}

	for n, v := range s.I {
		ss.I = append(ss.I, StorageTableSettings{
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
