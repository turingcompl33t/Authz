package authz

import (
	"errors"
	"fmt"
	"testing"
)

const (
	NA     = "na"
	String = "string"
	Int    = "int"
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

			b, err := AsBool(got)
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
		{StringExpr{"foo"}, "foo", nil},
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

			s, err := AsString(got)
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
		{IntExpr{42}, 42, nil},
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

			i, err := AsInt(got)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if i != d.want {
				t.Fatalf("got %v, want %v", got, d.want)
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
		{EqExpr{StringExpr{"foo"}, StringExpr{"foo"}}, true, nil},
		{EqExpr{StringExpr{"foo"}, StringExpr{"bar"}}, false, nil},
		{EqExpr{IntExpr{42}, IntExpr{42}}, true, nil},
		{EqExpr{IntExpr{42}, IntExpr{43}}, false, nil},
		{EqExpr{TrueExpr{}, TrueExpr{}}, true, nil},
		{EqExpr{TrueExpr{}, FalseExpr{}}, false, nil},
		{EqExpr{StringExpr{"foo"}, IntExpr{42}}, false, errors.New("")},
		{EqExpr{IntExpr{42}, StringExpr{"foo"}}, false, errors.New("")},
		{EqExpr{TrueExpr{}, StringExpr{"foo"}}, false, errors.New("")},
		{EqExpr{StringExpr{"foo"}, TrueExpr{}}, false, errors.New("")},
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

			b, err := AsBool(got)
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
		{VariableRefExpr{"foo"}, map[string]interface{}{"foo": 42}, 42, Int, nil},
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
		w, err := AsString(want)
		if err != nil {
			return errors.New("specified the wrong wantType")
		}

		g, err := AsString(got)
		if err != nil {
			return fmt.Errorf("got the wrong type; expected string")
		}

		if w != g {
			return fmt.Errorf("got %v, want %v", got, want)
		}

		return nil
	case Int:
		w, err := AsInt(want)
		if err != nil {
			return errors.New("specified the wrong wantType")
		}

		g, err := AsInt(got)
		if err != nil {
			fmt.Println(err)
			return fmt.Errorf("got the wrong type; expected int")
		}

		if w != g {
			return fmt.Errorf("got %v, want %v", got, want)
		}

		return nil
	case Bool:
		w, err := AsBool(want)
		if err != nil {
			return errors.New("specified the wrong wantType")
		}

		g, err := AsBool(got)
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
