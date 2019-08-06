package myfdbstorage

import "testing"

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
