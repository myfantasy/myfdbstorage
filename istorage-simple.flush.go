package myfdbstorage

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"os"
)

// FlushOnDisk changes into disk
func (s *IStorageSimple) FlushOnDisk() error {

	s.mx.RLock()
	defer s.mx.RUnlock()
	s.mxDump.Lock()
	defer s.mxDump.Unlock()

	rmOld := false

	ok, err := FileExists(s.DumpPath)
	if err != nil {
		return ErrorNew("FlushOnDisk Check directory "+s.DumpPath, err)
	}
	if !ok {
		err = os.MkdirAll(s.DumpPath, 0760)
		if err != nil {
			return ErrorNew("FlushOnDisk Mkdir file "+s.DumpPath, err)
		}
	}

	fileName := s.DumpPath + "data.iss"
	fileNameOld := s.DumpPath + "data_old.iss"

	ok, err = FileExists(fileNameOld)
	if err != nil {
		return ErrorNew("FlushOnDisk Check old file "+fileNameOld, err)
	}

	if ok {
		err = os.Remove(fileNameOld)
		if err != nil {
			return ErrorNew("FlushOnDisk Remove old file "+fileNameOld, err)
		}
	}

	ok, err = FileExists(fileName)
	if err != nil {
		return ErrorNew("FlushOnDisk Check current file "+fileName, err)
	}
	if ok {
		err = os.Rename(fileName, fileNameOld)
		if err != nil {
			return ErrorNew("FlushOnDisk Rename current file "+fileName+" to old file "+fileNameOld, err)
		}
		rmOld = true
	}

	ss := s.IStorageSimpleStorGet()

	var gdata bytes.Buffer
	enc := gob.NewEncoder(&gdata)

	err = enc.Encode(ss)
	if err != nil {
		return ErrorNew("Encode data for file "+fileName, err)
	}

	err = ioutil.WriteFile(fileName, gdata.Bytes(), 0660)
	if err != nil {
		return ErrorNew("FlushOnDisk write file "+fileName, err)
	}

	if rmOld {
		ok, err = FileExists(fileNameOld)
		if err != nil {
			return ErrorNew("FlushOnDisk Check old file "+fileNameOld, err)
		}

		if ok {
			err = os.Remove(fileNameOld)
			if err != nil {
				return ErrorNew("FlushOnDisk Remove old file "+fileNameOld, err)
			}
		}
	}

	return nil

}

// IStorageSimpleStorGet Get simple storage (IStorageSimpleStor)
func (s *IStorageSimple) IStorageSimpleStorGet() IStorageSimpleStor {
	d := IStorageSimpleStor{Data: make([]IItemSimple, 0, len(s.Data))}

	for _, v := range s.Data {
		d.Data = append(d.Data, IItemSimpleConvertTo(v))
	}
	return d
}

// IStorageSimpleStor Simple structuire for save items
type IStorageSimpleStor struct {
	Data []IItemSimple
}

// LoadFromDisk Get data from disk into disk
func (s *IStorageSimple) LoadFromDisk() error {

	s.mx.RLock()
	defer s.mx.RUnlock()
	s.mxDump.Lock()
	defer s.mxDump.Unlock()

	fileName := s.DumpPath + "data.iss"

	ok, err := FileExists(fileName)
	if err != nil {
		return ErrorNew("LoadFromDisk Check current file "+fileName, err)
	}
	if !ok {
		return nil
	}

	gdata, err := ioutil.ReadFile(fileName)
	if err != nil {
		return ErrorNew("LoadFromDisk get current file data "+fileName, err)
	}

	bf := bytes.NewBuffer(gdata)

	dec := gob.NewDecoder(bf)
	var ss IStorageSimpleStor
	err = dec.Decode(&ss)
	if err != nil {
		return ErrorNew("LoadFromDisk gob decode fail "+fileName, err)
	}

	err = s.IStorageSimpleStorRestore(ss)

	if err != nil {
		return ErrorNew("LoadFromDisk IStorageSimpleStorRestore fail "+fileName, err)
	}

	return nil
}

// IStorageSimpleStorRestore Restor from simple storage(IStorageSimpleStor)
func (s *IStorageSimple) IStorageSimpleStorRestore(d IStorageSimpleStor) error {

	for _, v := range d.Data {
		s.Data[v.ID()] = v
	}
	return nil
}

// ClearFromDisk clear disk store data
func (s *IStorageSimple) ClearFromDisk() error {

	s.mx.RLock()
	defer s.mx.RUnlock()
	s.mxDump.Lock()
	defer s.mxDump.Unlock()

	fileName := s.DumpPath + "data.iss"
	fileNameOld := s.DumpPath + "data_old.iss"

	ok, err := FileExists(fileNameOld)
	if err != nil {
		return ErrorNew("ClearFromDisk Check old file "+fileNameOld, err)
	}

	if ok {
		err = os.Remove(fileNameOld)
		if err != nil {
			return ErrorNew("ClearFromDisk Remove old file "+fileNameOld, err)
		}
	}

	ok, err = FileExists(fileName)
	if err != nil {
		return ErrorNew("ClearFromDisk Check file "+fileName, err)
	}

	if ok {
		err = os.Remove(fileName)
		if err != nil {
			return ErrorNew("ClearFromDisk Remove file "+fileName, err)
		}
	}

	return nil

}
