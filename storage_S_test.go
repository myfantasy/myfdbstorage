package myfdbstorage

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestStorageLoadS(t *testing.T) {

	s, err := StorageLoad([]byte("{}"))

	if err != nil {
		t.Fatalf("StorageLoad Fail")
	}

	err = s.AddTableS("tst", SStorageSimpleName, make(map[string]interface{}))
	if err != nil {
		t.Fatalf("AddTableS Fail")
	}

	st, ok, err := s.SGetTable("tst")
	if err != nil {
		t.Fatalf("SGetTable Fail")
	}
	if !ok {
		t.Fatalf("SGetTable table not exists Fail")
	}

	err = st.Set(SItemSimple{Key: "9", BData: []byte{12, 34, 45}})
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

	st2, ok, err := s2.SGetTable("tst")
	if err != nil {
		t.Fatalf("SGetTable 2 Fail")
	}
	if !ok {
		t.Fatalf("SGetTable 2 table not exists Fail")
	}

	err = st2.Set(SItemSimple{Key: "9", BData: []byte{12, 34, 45}})
	if err != nil {
		t.Fatalf("Set 2 Fail")
	}

	ok, err = s.DropTableS("tst")
	if err != nil {
		t.Fatalf("DropTableS Fail")
	}
	if !ok {
		t.Fatalf("DropTableS table not exists Fail")
	}

	ok, err = s.DropTableS("tst")
	if err != nil {
		t.Fatalf("DropTableS 2 Fail")
	}
	if ok {
		t.Fatalf("DropTableS 2 table exists Fail")
	}

}
