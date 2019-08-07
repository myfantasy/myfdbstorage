package myfdbstorage

import (
	"os"
	"testing"
	"time"
)

func STTestSet(t *testing.T) {
	s, err := SStorageCreate("simple", nil)
	if err != nil {
		t.Fatalf("SStorageCreate Fail")
	}

	err = s.Set(SItemSimple{Key: "9", BData: []byte{12, 34, 45}})
	if err != nil {
		t.Fatalf("Set Fail")
	}

	err = s.Set(SItemSimple{Key: "9", BData: []byte{12, 34, 45}})
	if err != nil {
		t.Fatalf("Set2 Fail")
	}

	ok, err := s.Exists("9")
	if err != nil {
		t.Fatalf("Exists Fail")
	}
	if !ok {
		t.Fatalf("Exists Fail ok")
	}

	i, ok, err := s.Get("9")
	if err != nil {
		t.Fatalf("Get Fail")
	}
	if !ok {
		t.Fatalf("Get Fail ok")
	}
	if i.Data()[1] != 34 {
		t.Fatalf("Get Fail Data")
	}

	exists, err := s.Del("9")
	if err != nil {
		t.Fatalf("Del Fail")
	}
	if !exists {
		t.Fatalf("Del Fail Exists")
	}

	exists, err = s.Del("9")
	if err != nil {
		t.Fatalf("Del2 Fail")
	}
	if exists {
		t.Fatalf("Del2 Fail Exists")
	}

	ok, err = s.Exists("9")
	if err != nil {
		t.Fatalf("Exists2 Fail")
	}
	if ok {
		t.Fatalf("Exists2 Fail ok")
	}

	_, ok, err = s.Get("9")
	if err != nil {
		t.Fatalf("Get2 Fail")
	}
	if ok {
		t.Fatalf("Get2 Fail ok")
	}
}

func STTestSet2(t *testing.T) {

	params := make(map[string]interface{})

	params["dump_path"] = "tmp/"
	params["dump_timeout"] = time.Second

	s, err := SStorageCreate("simple", params)
	if err != nil {
		t.Fatalf("SStorageCreate Fail, %s", err.Error())
	}

	err = s.Set(SItemSimple{Key: "9", BData: []byte{12, 34, 45}})
	if err != nil {
		t.Fatalf("Set Fail, %s", err.Error())
	}

	err = s.Set(SItemSimple{Key: "12", BData: []byte{13, 25, 46}})
	if err != nil {
		t.Fatalf("Set2 Fail, %s", err.Error())
	}

	err = s.Flush()
	if err != nil {
		t.Fatalf("Flush Fail, %s", err.Error())
	}

	s2, err := SStorageCreate("simple", params)
	if err != nil {
		t.Fatalf("SStorageCreate 2 Fail, %s", err.Error())
	}

	i, ok, err := s2.Get("9")
	if err != nil {
		t.Fatalf("Get Fail, %s", err.Error())
	}
	if !ok {
		t.Fatalf("Get Fail ok")
	}
	if i.Data()[1] != 34 {
		t.Fatalf("Get Fail Data")
	}

	err = s.ClearAndDeleteStorage()
	if err != nil {
		t.Fatalf("ClearAndDeleteStorage Fail, %s", err.Error())
	}

	err = s2.ClearAndDeleteStorage()
	if err != nil {
		t.Fatalf("ClearAndDeleteStorage 2 Fail, %s", err.Error())
	}

	err = os.Remove("tmp/")
	if err != nil {
		t.Fatalf("os.Remove(tmp/) Fail, %s", err.Error())
	}
}
