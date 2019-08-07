package myfdbstorage

import (
	"testing"
	"time"
)

func TestTryGetString(t *testing.T) {
	params := make(map[string]interface{})

	params["tst"] = "tst"

	var s string

	if !TryGetString(params, "tst", &s) {
		t.Fatalf("TryGetString Fail")
	}

	if s != "tst" {
		t.Fatalf("TryGetString value Fail")
	}

	if TryGetString(params, "tst2", &s) {
		t.Fatalf("TryGetString 2 Fail")
	}
}

func TestTryGetInt(t *testing.T) {
	params := make(map[string]interface{})

	params["tst"] = 5

	var s int

	if !TryGetInt(params, "tst", &s) {
		t.Fatalf("TryGetInt Fail")
	}

	if s != 5 {
		t.Fatalf("TryGetInt value Fail")
	}

	if TryGetInt(params, "tst2", &s) {
		t.Fatalf("TryGetInt 2 Fail")
	}
}

func TestTryGetDuration(t *testing.T) {
	params := make(map[string]interface{})

	params["tst"] = 5

	var s time.Duration

	if !TryGetDuration(params, "tst", &s) {
		t.Fatalf("TryGetDuration Fail")
	}

	if s != 5 {
		t.Fatalf("TryGetDuration value Fail")
	}

	if TryGetDuration(params, "tst2", &s) {
		t.Fatalf("TryGetDuration 2 Fail")
	}
}
