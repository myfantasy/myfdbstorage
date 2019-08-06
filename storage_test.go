package myfdbstorage

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestStorageLoad(t *testing.T) {

	s, err := StorageLoad([]byte("{}"))

	if err != nil {
		t.Fatalf("StorageLoad Fail")
	}

	err = s.AddTableI("tst", IStorageSimpleName, make(map[string]interface{}))
	if err != nil {
		t.Fatalf("AddTableI Fail")
	}

	st, ok, err := s.IGetTable("tst")
	if err != nil {
		t.Fatalf("IGetTable Fail")
	}
	if !ok {
		t.Fatalf("IGetTable table not exists Fail")
	}

	err = st.Set(IItemSimple{Key: 9, BData: []byte{12, 34, 45}})
	if err != nil {
		t.Fatalf("Set Fail")
	}

	err = s.StorageStructFlush("st.tmp")
	if err != nil {
		t.Fatalf("StorageStructFlush Fail")
	}

	b, err := ioutil.ReadFile("st.tmp")
	if err != nil {
		t.Fatalf("ioutil.ReadFile Fail")
	}

	s2, err := StorageLoad(b)
	if err != nil {
		t.Fatalf("StorageLoad 2 Fail")
	}

	err = os.Remove("st.tmp")
	if err != nil {
		t.Fatalf("os.Remove Fail")
	}

	st2, ok, err := s2.IGetTable("tst")
	if err != nil {
		t.Fatalf("IGetTable 2 Fail")
	}
	if !ok {
		t.Fatalf("IGetTable 2 table not exists Fail")
	}

	err = st2.Set(IItemSimple{Key: 9, BData: []byte{12, 34, 45}})
	if err != nil {
		t.Fatalf("Set 2 Fail")
	}

	ok, err = s.DropTableI("tst")
	if err != nil {
		t.Fatalf("DropTableI Fail")
	}
	if !ok {
		t.Fatalf("DropTableI table not exists Fail")
	}

	ok, err = s.DropTableI("tst")
	if err != nil {
		t.Fatalf("DropTableI 2 Fail")
	}
	if ok {
		t.Fatalf("DropTableI 2 table exists Fail")
	}

}
