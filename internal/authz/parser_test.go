package authz

import (
	"errors"
	"testing"
)

// ExprParser can parse $true() and $false() expressions.
func TestParseTrueFalse(t *testing.T) {
	data := []struct {
		input       string
		want        Expr
		expectError error
	}{
		{"", nil, errors.New("")},
		{"true", TrueExpr{}, nil},
		{"true,", nil, errors.New("")},
		{"true ", nil, errors.New("")},
		{"false", FalseExpr{}, nil},
		{"false,", nil, errors.New("")},
		{"false ", nil, errors.New("")},
	}
	for _, d := range data {
		t.Run(d.input, func(t *testing.T) {
			b := ExprParser{}
			got, err := b.Parse(d.input)

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

			if !got.Equal(d.want) {
				t.Fatalf("got %v, want %v", got, d.want)
			}
		})
	}
}

// ExprParser can parse string literal expressions.
func TestParseString(t *testing.T) {
	data := []struct {
		input       string
		want        Expr
		expectError error
	}{
		{"", nil, errors.New("")},
		{"'foo'", StringExpr{"foo"}, nil},
		{"'foo", nil, errors.New("")},
		{"'foo',", nil, errors.New("")},
		{"'foo' ", nil, errors.New("")},
	}
	for _, d := range data {
		t.Run(d.input, func(t *testing.T) {
			b := ExprParser{}
			got, err := b.Parse(d.input)

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

			if !got.Equal(d.want) {
				t.Fatalf("got %v, want %v", got, d.want)
			}
		})
	}
}

// ExprParser can parse integer literal expressions.
func TestParseInt(t *testing.T) {
	data := []struct {
		input       string
		want        Expr
		expectError error
	}{
		{"", nil, errors.New("")},
		{"123", IntExpr{123}, nil},
		{"123,", nil, errors.New("")},
		{"123 ", nil, errors.New("")},
	}
	for _, d := range data {
		t.Run(d.input, func(t *testing.T) {
			b := ExprParser{}
			got, err := b.Parse(d.input)

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

			if !got.Equal(d.want) {
				t.Fatalf("got %v, want %v", got, d.want)
			}
		})
	}
}

// ExprParser can parse equality expressions.
func TestParseEq(t *testing.T) {
	data := []struct {
		input       string
		want        Expr
		expectError error
	}{
		{"", nil, errors.New("")},
		{"$eq(true, true)", EqExpr{TrueExpr{}, TrueExpr{}}, nil},
		{"$eq(true, false)", EqExpr{TrueExpr{}, FalseExpr{}}, nil},
		{"$eq(true, 'foo')", EqExpr{TrueExpr{}, StringExpr{"foo"}}, nil},
		{"$eq(true, 123)", EqExpr{TrueExpr{}, IntExpr{123}}, nil},
		{"$eq(", nil, errors.New("")},
		{"$eq(true", nil, errors.New("")},
		{"$eq(true),", nil, errors.New("")},
		{"$eq() ", nil, errors.New("")},
	}
	for _, d := range data {
		t.Run(d.input, func(t *testing.T) {
			b := ExprParser{}
			got, err := b.Parse(d.input)

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

			if !got.Equal(d.want) {
				t.Fatalf("got %v, want %v", got, d.want)
			}
		})
	}
}

// ExprParser can parse variable reference expressions.
func TestParseVariableRef(t *testing.T) {
	data := []struct {
		input       string
		want        Expr
		expectError error
	}{
		{"", nil, errors.New("")},
		{"foo", VariableRefExpr{"foo"}, nil},
		{"foobar", VariableRefExpr{"foobar"}, nil},
		{"foo,", nil, errors.New("")},
		{"foo ", nil, errors.New("")},
		{"foo(", nil, errors.New("")},
	}
	for _, d := range data {
		t.Run(d.input, func(t *testing.T) {
			b := ExprParser{}
			got, err := b.Parse(d.input)

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

			if !got.Equal(d.want) {
				t.Fatalf("got %v, want %v", got, d.want)
			}
		})
	}
}

// ExprParser can parse struct field reference expressions.
func TestParseStructFieldRef(t *testing.T) {
	data := []struct {
		input       string
		want        Expr
		expectError error
	}{
		{"", nil, errors.New("")},
		{"foo.bar", StructFieldRefExpr{"foo", "bar"}, nil},
		{"foo.bar,", nil, errors.New("")},
		{"foo.bar.baz ", nil, errors.New("")},
		{"foo.bar(", nil, errors.New("")},
	}
	for _, d := range data {
		t.Run(d.input, func(t *testing.T) {
			b := ExprParser{}
			got, err := b.Parse(d.input)

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

			if !got.Equal(d.want) {
				t.Fatalf("got %v, want %v", got, d.want)
			}
		})
	}
}
