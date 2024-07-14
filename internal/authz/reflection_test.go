package authz

import (
	"testing"
)

// GetField retrieves a field from a structure by name.
func TestGetField(t *testing.T) {
	obj := struct{ Name string }{
		Name: "foo",
	}

	v, err := GetField(obj, "Name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	switch vt := v.(type) {
	case string:
		if v != "foo" {
			t.Fatalf("got %v, want %v", t, "foo")
		}
	default:
		t.Fatalf("unexpected type: %T", vt)
	}
}

// GetField fails with error when an invalid field is requested.
func TestGetFieldFail(t *testing.T) {
	obj := struct{ Name string }{
		Name: "foo",
	}
	_, err := GetField(obj, "Bar")
	if err == nil {
		t.Fatalf("expected error")
	}
}

// GetField fails with error when invoked on non-struct object.
func TestGetFieldNonStruct(t *testing.T) {
	obj := "hello"
	_, err := GetField(obj, "Name")
	if err == nil {
		t.Fatalf("expected error")
	}
}
