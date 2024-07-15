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
		{"'foo'", StrExpr{"foo"}, nil},
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
		{"123", UintExpr{123}, nil},
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

// ExprParser can parse boolean slice literal expressions.
func TestParseBoolSlice(t *testing.T) {
	data := []struct {
		input       string
		want        Expr
		expectError error
	}{
		{"", nil, errors.New("")},
		{"[]bool{}", BoolSliceExpr{[]Expr{}}, nil},
		{"[]bool{true}", BoolSliceExpr{[]Expr{TrueExpr{}}}, nil},
		{"[]bool{false}", BoolSliceExpr{[]Expr{FalseExpr{}}}, nil},
		{"[]bool{true, false}", BoolSliceExpr{[]Expr{TrueExpr{}, FalseExpr{}}}, nil},
		{"[]bool{true,}", nil, errors.New("")},
		{"[]bool{true}", BoolSliceExpr{[]Expr{TrueExpr{}}}, nil},
		{"[]bool{true", nil, errors.New("")},
		{"[]bool{1}", nil, errors.New("")},
		{"[]bool{'hello'}", nil, errors.New("")},
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

// ExprParser can parse string slice literal expressions.
func TestParseStringSlice(t *testing.T) {
	data := []struct {
		input       string
		want        Expr
		expectError error
	}{
		{"", nil, errors.New("")},
		{"[]string{}", StrSliceExpr{[]Expr{}}, nil},
		{"[]string{'foo'}", StrSliceExpr{[]Expr{StrExpr{"foo"}}}, nil},
		{"[]string{'foo', 'bar'}", StrSliceExpr{[]Expr{StrExpr{"foo"}, StrExpr{"bar"}}}, nil},
		{"[]string{'foo',}", nil, errors.New("")},
		{"[]string{'foo", nil, errors.New("")},
		{"[]string{1}", nil, errors.New("")},
		{"[]string{true}", nil, errors.New("")},
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

// ExprParser can parse uint slice literal expressions.
func TestParseUintSlice(t *testing.T) {
	data := []struct {
		input       string
		want        Expr
		expectError error
	}{
		{"", nil, errors.New("")},
		{"[]uint{}", UintSliceExpr{[]Expr{}}, nil},
		{"[]uint{123}", UintSliceExpr{[]Expr{UintExpr{123}}}, nil},
		{"[]uint{123, 456}", UintSliceExpr{[]Expr{UintExpr{123}, UintExpr{456}}}, nil},
		{"[]uint{123,}", nil, errors.New("")},
		{"[]uint{123", nil, errors.New("")},
		{"[]uint{'foo'}", nil, errors.New("")},
		{"[]uint{true}", nil, errors.New("")},
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
		{"$eq(true, 'foo')", EqExpr{TrueExpr{}, StrExpr{"foo"}}, nil},
		{"$eq(true, 123)", EqExpr{TrueExpr{}, UintExpr{123}}, nil},
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

// ExprParser can parse $in expressions.
func TestParseIn(t *testing.T) {
	data := []struct {
		input       string
		want        Expr
		expectError error
	}{
		{"", nil, errors.New("")},
		{"$in(1, []uint{1, 2})", InExpr{UintExpr{1}, UintSliceExpr{[]Expr{UintExpr{1}, UintExpr{2}}}}, nil},
		{"$in(true, []bool{true, false})", InExpr{TrueExpr{}, BoolSliceExpr{[]Expr{TrueExpr{}, FalseExpr{}}}}, nil},
		{"$in('foo', []string{'foo', 'bar'})", InExpr{StrExpr{"foo"}, StrSliceExpr{[]Expr{StrExpr{"foo"}, StrExpr{"bar"}}}}, nil},
		{"$in(", nil, errors.New("")},
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

// ExprParser can parse AND expressions.
func TestParseAnd(t *testing.T) {
	data := []struct {
		input       string
		want        Expr
		expectError error
	}{
		{"", nil, errors.New("")},
		{"$and(true, true)", AndExpr{[]Expr{TrueExpr{}, TrueExpr{}}}, nil},
		{"$and(true, false)", AndExpr{[]Expr{TrueExpr{}, FalseExpr{}}}, nil},
		{"$and(1, 2)", AndExpr{[]Expr{UintExpr{1}, UintExpr{2}}}, nil},
		{"$and('foo', 'bar')", AndExpr{[]Expr{StrExpr{"foo"}, StrExpr{"bar"}}}, nil},
		{"$and(", nil, errors.New("")},
		{"$and(true", nil, errors.New("")},
		{"$and(true),", nil, errors.New("")},
		{"$and()", AndExpr{[]Expr{}}, nil},
		{"$and() ", nil, errors.New("")},
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

// ExprParser can parse OR expressions.
func TestParseOr(t *testing.T) {
	data := []struct {
		input       string
		want        Expr
		expectError error
	}{
		{"", nil, errors.New("")},
		{"$or(true, true)", OrExpr{[]Expr{TrueExpr{}, TrueExpr{}}}, nil},
		{"$or(true, false)", OrExpr{[]Expr{TrueExpr{}, FalseExpr{}}}, nil},
		{"$or(1, 2)", OrExpr{[]Expr{UintExpr{1}, UintExpr{2}}}, nil},
		{"$or('foo', 'bar')", OrExpr{[]Expr{StrExpr{"foo"}, StrExpr{"bar"}}}, nil},
		{"$or(", nil, errors.New("")},
		{"$or(true", nil, errors.New("")},
		{"$or(true),", nil, errors.New("")},
		{"$or()", OrExpr{[]Expr{}}, nil},
		{"$or() ", nil, errors.New("")},
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
