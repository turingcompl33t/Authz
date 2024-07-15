package authz

import (
	"errors"
	"fmt"
	"testing"
)

// Truth evalutation works as expected.
func TestTruthy(t *testing.T) {
	data := []struct {
		input       interface{}
		expect      bool
		expectError error
	}{
		{"hello", true, nil},
		{"", false, nil},
		{1, true, nil},
		{0, false, nil},
		{true, true, nil},
		{false, false, nil},
		{nil, false, errors.New("")},
	}
	for _, d := range data {
		t.Run(fmt.Sprintf("%v", d.input), func(t *testing.T) {
			got, err := truthy(d.input)
			if err != nil {
				if d.expectError == nil {
					t.Fatalf("unexpected error: %v", err)
				} else {
					return
				}
			}

			if d.expectError != nil {
				if err == nil {
					t.Fatalf("expected error: %v", d.expectError)
				} else {
					return
				}
			}

			if got != d.expect {
				t.Fatalf("got %v, want %v", got, d.expect)
			}
		})
	}
}

// Coercion to slice of string works as expected.
func TestCoerceStringSlice(t *testing.T) {
	data := []struct {
		input       interface{}
		want        []string
		expectError error
	}{
		{[]string{"a", "b"}, []string{"a", "b"}, nil},
		{[]string{"a", "b"}, []string{"a", "b"}, nil},
		{[]interface{}{"a", 1}, nil, errors.New("")},
		{[]int{1, 2}, nil, errors.New("")},
	}
	for _, d := range data {
		t.Run(fmt.Sprintf("%v", d.input), func(t *testing.T) {
			got, err := coerceStrSlice(d.input)
			if err != nil {
				if d.expectError == nil {
					t.Fatalf("unexpected error: %v", err)
				} else {
					return
				}
			}

			if d.expectError != nil {
				if err == nil {
					t.Fatalf("expected error: %v", d.expectError)
				} else {
					return
				}
			}

			if len(got) != len(d.want) {
				t.Fatalf("got %v, want %v", got, d.want)
			}
		})
	}
}

// Coercion of slice of uint works as expected.
func TestCoerceUintSlice(t *testing.T) {
	data := []struct {
		input       interface{}
		want        []uint
		expectError error
	}{
		{[]uint{1, 2}, []uint{1, 2}, nil},
		{[]int{1, 2}, []uint{1, 2}, nil},
		{[]string{"a", "b"}, nil, errors.New("")},
		{[]interface{}{"a", 1}, nil, errors.New("")},
		{"hello", nil, errors.New("")},
	}
	for _, d := range data {
		t.Run(fmt.Sprintf("%v", d.input), func(t *testing.T) {
			got, err := coerceUintSlice(d.input)
			if err != nil {
				if d.expectError == nil {
					t.Fatalf("unexpected error: %v", err)
				} else {
					return
				}
			}

			if d.expectError != nil {
				if err == nil {
					t.Fatalf("expected error: %v", d.expectError)
				} else {
					return
				}
			}

			if len(got) != len(d.want) {
				t.Fatalf("got %v, want %v", got, d.want)
			}
		})
	}
}

// Coercion of slice of bool works as expected.
func TestCoerceBoolSlice(t *testing.T) {
	data := []struct {
		input       interface{}
		want        []bool
		expectError error
	}{
		{[]bool{true, false}, []bool{true, false}, nil},
		{[]int{1, 2}, nil, errors.New("")},
		{[]string{"a", "b"}, nil, errors.New("")},
	}
	for _, d := range data {
		t.Run(fmt.Sprintf("%v", d.input), func(t *testing.T) {
			got, err := coerceBoolSlice(d.input)
			if err != nil {
				if d.expectError == nil {
					t.Fatalf("unexpected error: %v", err)
				} else {
					return
				}
			}

			if d.expectError != nil {
				if err == nil {
					t.Fatalf("expected error: %v", d.expectError)
				} else {
					return
				}
			}

			if len(got) != len(d.want) {
				t.Fatalf("got %v, want %v", got, d.want)
			}
		})
	}
}
