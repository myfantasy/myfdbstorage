package myfdbstorage

import "testing"

func TestSet(t *testing.T) {
	s, err := IStorageCreate("simple", nil)
	if err != nil {
		t.Fatalf("IStorageCreate Fail")
	}

	err = s.Set(IItemSimple{Key: 9, BData: []byte{12, 34, 45}})
	if err != nil {
		t.Fatalf("Set Fail")
	}

	err = s.Set(IItemSimple{Key: 9, BData: []byte{12, 34, 45}})
	if err != nil {
		t.Fatalf("Set2 Fail")
	}

	ok, err := s.Exists(9)
	if err != nil {
		t.Fatalf("Exists Fail")
	}
	if !ok {
		t.Fatalf("Exists Fail ok")
	}

	i, ok, err := s.Get(9)
	if err != nil {
		t.Fatalf("Get Fail")
	}
	if !ok {
		t.Fatalf("Get Fail ok")
	}
	if i.Data()[1] != 34 {
		t.Fatalf("Get Fail Data")
	}

	exists, err := s.Del(9)
	if err != nil {
		t.Fatalf("Del Fail")
	}
	if !exists {
		t.Fatalf("Del Fail Exists")
	}

	exists, err = s.Del(9)
	if err != nil {
		t.Fatalf("Del2 Fail")
	}
	if exists {
		t.Fatalf("Del2 Fail Exists")
	}

	ok, err = s.Exists(9)
	if err != nil {
		t.Fatalf("Exists2 Fail")
	}
	if ok {
		t.Fatalf("Exists2 Fail ok")
	}

	_, ok, err = s.Get(9)
	if err != nil {
		t.Fatalf("Get2 Fail")
	}
	if ok {
		t.Fatalf("Get2 Fail ok")
	}
}
