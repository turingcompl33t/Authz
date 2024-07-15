package authz

import (
	"errors"
	"fmt"
	"testing"
)

const (
	NA     = "na"
	String = "string"
	Uint   = "uint"
	Bool   = "bool"
)

// Evaluator can evaluate boolean literal expressions.
func TestEvalBools(t *testing.T) {
	data := []struct {
		input       Expr
		want        bool
		expectError error
	}{
		{TrueExpr{}, true, nil},
		{FalseExpr{}, false, nil},
	}
	for _, d := range data {
		t.Run(fmt.Sprintf("%v", d.input), func(t *testing.T) {
			ev := Evaluator{}

			got, err := ev.Eval(d.input, make(map[string]interface{}))
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

			b, err := coerceBool(got)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if b != d.want {
				t.Fatalf("got %v, want %v", got, d.want)
			}
		})
	}
}

// Evaluator can evaluate string literal expressions.
func TestEvalString(t *testing.T) {
	data := []struct {
		input       Expr
		want        string
		expectError error
	}{
		{StrExpr{"foo"}, "foo", nil},
	}
	for _, d := range data {
		t.Run(fmt.Sprintf("%v", d.input), func(t *testing.T) {
			ev := Evaluator{}

			got, err := ev.Eval(d.input, make(map[string]interface{}))
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

			s, err := coerceStr(got)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if s != d.want {
				t.Fatalf("got %v, want %v", got, d.want)
			}
		})
	}
}

// Evaluator can evaluate integer literal expressions.
func TestEvalInt(t *testing.T) {
	data := []struct {
		input       Expr
		want        uint
		expectError error
	}{
		{UintExpr{42}, 42, nil},
	}
	for _, d := range data {
		t.Run(fmt.Sprintf("%v", d.input), func(t *testing.T) {
			ev := Evaluator{}

			got, err := ev.Eval(d.input, make(map[string]interface{}))
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

			i, err := coerceUint(got)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if i != d.want {
				t.Fatalf("got %v, want %v", got, d.want)
			}
		})
	}
}

// Evaluator can evaluate string slice expressions.
func TestEvalStrSlice(t *testing.T) {
	data := []struct {
		input       Expr
		want        []string
		expectError error
	}{
		{StrSliceExpr{[]Expr{}}, []string{}, nil},
		{StrSliceExpr{[]Expr{StrExpr{"foo"}, StrExpr{"bar"}}}, []string{"foo", "bar"}, nil},
		{StrSliceExpr{[]Expr{StrExpr{"foo"}, UintExpr{42}}}, nil, errors.New("")},
	}
	for _, d := range data {
		t.Run(fmt.Sprintf("%v", d.input), func(t *testing.T) {
			ev := Evaluator{}

			got, err := ev.Eval(d.input, make(map[string]interface{}))
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

			s, err := coerceStrSlice(got)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(s) != len(d.want) {
				t.Fatalf("got %v, want %v", got, d.want)
			}

			for i, e := range s {
				if e != d.want[i] {
					t.Fatalf("got %v, want %v", got, d.want)
				}
			}
		})
	}
}

// Evaluator can evaluate integer slice expressions.
func TestEvalUintSlice(t *testing.T) {
	data := []struct {
		input       Expr
		want        []uint
		expectError error
	}{
		{UintSliceExpr{[]Expr{}}, []uint{}, nil},
		{UintSliceExpr{[]Expr{UintExpr{42}, UintExpr{43}}}, []uint{42, 43}, nil},
		{UintSliceExpr{[]Expr{UintExpr{42}, StrExpr{"foo"}}}, nil, errors.New("")},
	}
	for _, d := range data {
		t.Run(fmt.Sprintf("%v", d.input), func(t *testing.T) {
			ev := Evaluator{}

			got, err := ev.Eval(d.input, make(map[string]interface{}))
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

			s, err := coerceUintSlice(got)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(s) != len(d.want) {
				t.Fatalf("got %v, want %v", got, d.want)
			}

			for i, e := range s {
				if e != d.want[i] {
					t.Fatalf("got %v, want %v", got, d.want)
				}
			}
		})
	}
}

// Evaluattor can evaluate boolean slice expressions.
func TestEvalBoolSlice(t *testing.T) {
	data := []struct {
		input       Expr
		want        []bool
		expectError error
	}{
		{BoolSliceExpr{[]Expr{}}, []bool{}, nil},
		{BoolSliceExpr{[]Expr{TrueExpr{}, FalseExpr{}}}, []bool{true, false}, nil},
		{BoolSliceExpr{[]Expr{TrueExpr{}, StrExpr{"foo"}}}, nil, errors.New("")},
	}
	for _, d := range data {
		t.Run(fmt.Sprintf("%v", d.input), func(t *testing.T) {
			ev := Evaluator{}

			got, err := ev.Eval(d.input, make(map[string]interface{}))
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

			s, err := coerceBoolSlice(got)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(s) != len(d.want) {
				t.Fatalf("got %v, want %v", got, d.want)
			}

			for i, e := range s {
				if e != d.want[i] {
					t.Fatalf("got %v, want %v", got, d.want)
				}
			}
		})
	}
}

// Evaluator can evaluate equality expressions.
func TestEvalEq(t *testing.T) {
	data := []struct {
		input       Expr
		want        bool
		expectError error
	}{
		{EqExpr{StrExpr{"foo"}, StrExpr{"foo"}}, true, nil},
		{EqExpr{StrExpr{"foo"}, StrExpr{"bar"}}, false, nil},
		{EqExpr{UintExpr{42}, UintExpr{42}}, true, nil},
		{EqExpr{UintExpr{42}, UintExpr{43}}, false, nil},
		{EqExpr{TrueExpr{}, TrueExpr{}}, true, nil},
		{EqExpr{TrueExpr{}, FalseExpr{}}, false, nil},
		{EqExpr{StrExpr{"foo"}, UintExpr{42}}, false, errors.New("")},
		{EqExpr{UintExpr{42}, StrExpr{"foo"}}, false, errors.New("")},
		{EqExpr{TrueExpr{}, StrExpr{"foo"}}, false, errors.New("")},
		{EqExpr{StrExpr{"foo"}, TrueExpr{}}, false, errors.New("")},
	}
	for _, d := range data {
		t.Run(fmt.Sprintf("%v", d.input), func(t *testing.T) {
			ev := Evaluator{}

			got, err := ev.Eval(d.input, make(map[string]interface{}))
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

			b, err := coerceBool(got)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if b != d.want {
				t.Fatalf("got %v, want %v", got, d.want)
			}
		})
	}
}

// Evaluator can evaliate $in expressions.
func TestEvalIn(t *testing.T) {
	data := []struct {
		input       Expr
		want        bool
		expectError error
	}{
		{InExpr{StrExpr{"foo"}, StrSliceExpr{[]Expr{StrExpr{"foo"}, StrExpr{"bar"}}}}, true, nil},
		{InExpr{StrExpr{"foo"}, StrSliceExpr{[]Expr{StrExpr{"bar"}, StrExpr{"baz"}}}}, false, nil},
		{InExpr{UintExpr{42}, UintSliceExpr{[]Expr{UintExpr{42}, UintExpr{43}}}}, true, nil},
		{InExpr{UintExpr{42}, UintSliceExpr{[]Expr{UintExpr{43}, UintExpr{44}}}}, false, nil},
		{InExpr{StrExpr{"foo"}, UintSliceExpr{[]Expr{UintExpr{42}, UintExpr{43}}}}, false, errors.New("")},
		{InExpr{UintExpr{42}, StrSliceExpr{[]Expr{StrExpr{"foo"}, StrExpr{"bar"}}}}, false, errors.New("")},
		{InExpr{TrueExpr{}, BoolSliceExpr{[]Expr{TrueExpr{}, FalseExpr{}}}}, true, nil},
		{InExpr{TrueExpr{}, BoolSliceExpr{[]Expr{FalseExpr{}, FalseExpr{}}}}, false, nil},
		{InExpr{TrueExpr{}, StrSliceExpr{[]Expr{StrExpr{"foo"}, StrExpr{"bar"}}}}, false, errors.New("")},
		{InExpr{StrExpr{"foo"}, BoolSliceExpr{[]Expr{TrueExpr{}, FalseExpr{}}}}, false, errors.New("")},
		{InExpr{StrExpr{"foo"}, StrSliceExpr{[]Expr{StrExpr{"foo"}, UintExpr{42}}}}, false, errors.New("")},
		{InExpr{UintExpr{42}, UintSliceExpr{[]Expr{UintExpr{42}, StrExpr{"foo"}}}}, false, errors.New("")},
		{InExpr{TrueExpr{}, BoolSliceExpr{[]Expr{TrueExpr{}, StrExpr{"foo"}}}}, false, errors.New("")},
		{InExpr{StrExpr{"foo"}, BoolSliceExpr{[]Expr{TrueExpr{}, StrExpr{"foo"}}}}, false, errors.New("")},
	}
	for _, d := range data {
		t.Run(fmt.Sprintf("%v", d.input), func(t *testing.T) {
			ev := Evaluator{}

			got, err := ev.Eval(d.input, make(map[string]interface{}))
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

			b, err := coerceBool(got)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if b != d.want {
				t.Fatalf("got %v, want %v", got, d.want)
			}
		})
	}
}

// Evaluator can evaluate AND expressions.
func TestEvalAnd(t *testing.T) {
	data := []struct {
		input       Expr
		want        bool
		expectError error
	}{
		{AndExpr{[]Expr{TrueExpr{}, TrueExpr{}}}, true, nil},
		{AndExpr{[]Expr{TrueExpr{}, FalseExpr{}}}, false, nil},
		{AndExpr{[]Expr{FalseExpr{}, TrueExpr{}}}, false, nil},
		{AndExpr{[]Expr{FalseExpr{}, FalseExpr{}}}, false, nil},
		{AndExpr{[]Expr{StrExpr{"foo"}, StrExpr{"bar"}}}, true, nil},
		{AndExpr{[]Expr{StrExpr{"foo"}, StrExpr{""}}}, false, nil},
		{AndExpr{[]Expr{UintExpr{42}, UintExpr{43}}}, true, nil},
		{AndExpr{[]Expr{UintExpr{42}, UintExpr{0}}}, false, nil},
		{AndExpr{[]Expr{TrueExpr{}, StrExpr{"foo"}, UintExpr{1}}}, true, nil},
	}
	for _, d := range data {
		t.Run(fmt.Sprintf("%v", d.input), func(t *testing.T) {
			ev := Evaluator{}

			got, err := ev.Eval(d.input, make(map[string]interface{}))
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

			b, err := coerceBool(got)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if b != d.want {
				t.Fatalf("got %v, want %v", got, d.want)
			}
		})
	}
}

// Evaluator can evaluate OR expressions.
func TestEvalOr(t *testing.T) {
	data := []struct {
		input       Expr
		want        bool
		expectError error
	}{
		{OrExpr{[]Expr{TrueExpr{}, TrueExpr{}}}, true, nil},
		{OrExpr{[]Expr{TrueExpr{}, FalseExpr{}}}, true, nil},
		{OrExpr{[]Expr{FalseExpr{}, TrueExpr{}}}, true, nil},
		{OrExpr{[]Expr{FalseExpr{}, FalseExpr{}}}, false, nil},
		{OrExpr{[]Expr{StrExpr{"foo"}, StrExpr{"bar"}}}, true, nil},
		{OrExpr{[]Expr{StrExpr{"foo"}, StrExpr{""}}}, true, nil},
		{OrExpr{[]Expr{UintExpr{42}, UintExpr{43}}}, true, nil},
		{OrExpr{[]Expr{UintExpr{42}, UintExpr{0}}}, true, nil},
		{OrExpr{[]Expr{FalseExpr{}, StrExpr{""}, UintExpr{0}}}, false, nil},
	}
	for _, d := range data {
		t.Run(fmt.Sprintf("%v", d.input), func(t *testing.T) {
			ev := Evaluator{}

			got, err := ev.Eval(d.input, make(map[string]interface{}))
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

			b, err := coerceBool(got)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if b != d.want {
				t.Fatalf("got %v, want %v", got, d.want)
			}
		})
	}
}

// Evaluator can evaluate variable reference expressions.
func TestEvalVariableRef(t *testing.T) {
	data := []struct {
		input       Expr
		env         map[string]interface{}
		want        interface{}
		wantType    string
		expectError error
	}{
		{VariableRefExpr{"foo"}, map[string]interface{}{"foo": "bar"}, "bar", String, nil},
		{VariableRefExpr{"foo"}, map[string]interface{}{"foo": 42}, 42, Uint, nil},
		{VariableRefExpr{"foo"}, map[string]interface{}{"foo": true}, true, Bool, nil},
		{VariableRefExpr{"foo"}, map[string]interface{}{}, nil, NA, errors.New("")},
		{VariableRefExpr{"foo"}, map[string]interface{}{"bar": "baz"}, nil, NA, errors.New("")},
	}
	for _, d := range data {
		t.Run(fmt.Sprintf("%v", d.input), func(t *testing.T) {
			ev := Evaluator{}

			got, err := ev.Eval(d.input, d.env)
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

			if err := coerceAndCmp(d.want, got, d.wantType); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// Evaluator can evaluate struct field reference expressions.
func TestEvalStructFieldRef(t *testing.T) {
	data := []struct {
		input       Expr
		env         map[string]interface{}
		want        interface{}
		wantType    string
		expectError error
	}{
		{StructFieldRefExpr{"foo", "Bar"}, map[string]interface{}{"foo": struct{ Bar string }{"baz"}}, "baz", String, nil},
		{StructFieldRefExpr{"foo", "Hello"}, map[string]interface{}{"foo": struct{ Bar string }{"baz"}}, "", NA, errors.New("")},
	}
	for _, d := range data {
		t.Run(fmt.Sprintf("%v", d.input), func(t *testing.T) {
			ev := Evaluator{}

			got, err := ev.Eval(d.input, d.env)
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

			if err := coerceAndCmp(d.want, got, d.wantType); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// Coerce the got value to the wantType and compare it with the want value.
func coerceAndCmp(want, got interface{}, wantType string) error {
	switch wantType {
	case String:
		w, err := coerceStr(want)
		if err != nil {
			return errors.New("specified the wrong wantType")
		}

		g, err := coerceStr(got)
		if err != nil {
			return fmt.Errorf("got the wrong type; expected string")
		}

		if w != g {
			return fmt.Errorf("got %v, want %v", got, want)
		}

		return nil
	case Uint:
		w, err := coerceUint(want)
		if err != nil {
			fmt.Println(err)
			return errors.New("specified the wrong wantType")
		}

		g, err := coerceUint(got)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("got the wrong type; expected uint")
		}

		if w != g {
			return fmt.Errorf("got %v, want %v", got, want)
		}

		return nil
	case Bool:
		w, err := coerceBool(want)
		if err != nil {
			return errors.New("specified the wrong wantType")
		}

		g, err := coerceBool(got)
		if err != nil {
			return fmt.Errorf("got the wrong type; expected bool")
		}

		if w != g {
			return fmt.Errorf("got %v, want %v", got, want)
		}

		return nil
	default:
		return fmt.Errorf("unexpected wantType %v", wantType)
	}
}
