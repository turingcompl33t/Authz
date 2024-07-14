package authz

import (
	"fmt"
)

type Expr interface {
	Eval(map[string]interface{}) (interface{}, error)
	Equal(Expr) bool
}

// ----------------------------------------------------------------------------
// TrueExpr
// ----------------------------------------------------------------------------

// TrueExpr represents the boolean literal `true`.
type TrueExpr struct {
}

func (t TrueExpr) Eval(env map[string]interface{}) (interface{}, error) {
	return true, nil
}

func (t TrueExpr) Equal(other Expr) bool {
	_, ok := other.(TrueExpr)
	return ok
}

// ----------------------------------------------------------------------------
// FalseExpr
// ----------------------------------------------------------------------------

// FalseExpr represents the boolean literal `false`.
type FalseExpr struct {
}

func (f FalseExpr) Eval(env map[string]interface{}) (interface{}, error) {
	return false, nil
}

func (f FalseExpr) Equal(other Expr) bool {
	_, ok := other.(FalseExpr)
	return ok
}

// ----------------------------------------------------------------------------
// StringExpr
// ----------------------------------------------------------------------------

// StringExpr represents a string literal.
type StringExpr struct {
	Value string
}

func (s StringExpr) Eval(env map[string]interface{}) (interface{}, error) {
	return s.Value, nil
}

func (s StringExpr) Equal(other Expr) bool {
	otherString, ok := other.(StringExpr)
	if !ok {
		return false
	}

	return s.Value == otherString.Value
}

// ----------------------------------------------------------------------------
// IntExpr
// ----------------------------------------------------------------------------

// IntExpr represents an integer literal.
type IntExpr struct {
	Value uint
}

func (i IntExpr) Eval(env map[string]interface{}) (interface{}, error) {
	return i.Value, nil
}

func (i IntExpr) Equal(other Expr) bool {
	otherInt, ok := other.(IntExpr)
	if !ok {
		return false
	}

	return i.Value == otherInt.Value
}

// ----------------------------------------------------------------------------
// EqExpr
// ----------------------------------------------------------------------------

// EqExpr represents an equality comparison.
type EqExpr struct {
	Left  Expr
	Right Expr
}

func (e EqExpr) Eval(env map[string]interface{}) (interface{}, error) {
	left, err := e.Left.Eval(env)
	if err != nil {
		return false, err
	}
	right, err := e.Right.Eval(env)
	if err != nil {
		return false, err
	}

	asStr, err := coerceStr(left)
	if err == nil {
		rAsStr, err := coerceStr(right)
		if err == nil {
			return asStr == rAsStr, nil
		} else {
			return nil, fmt.Errorf("mismatched types in equality comparison: %T != %T", left, right)
		}
	}

	asUint, err := coerceUint(left)
	if err == nil {
		rAsUint, err := coerceUint(right)
		if err == nil {
			return asUint == rAsUint, nil
		} else {
			return nil, fmt.Errorf("mismatched types in equality comparison: %T != %T", left, right)
		}
	}

	asBool, err := coerceBool(left)
	if err == nil {
		rAsBool, err := coerceBool(right)
		if err == nil {
			return asBool == rAsBool, nil
		} else {
			return nil, fmt.Errorf("mismatched types in equality comparison: %T != %T", left, right)
		}
	}

	return nil, fmt.Errorf("unsupported type in equality comparison: %T", left)
}

func (e EqExpr) Equal(other Expr) bool {
	otherEq, ok := other.(EqExpr)
	if !ok {
		return false
	}

	return e.Left.Equal(otherEq.Left) && e.Right.Equal(otherEq.Right)
}

// ----------------------------------------------------------------------------
// AndExpr
// ----------------------------------------------------------------------------

type AndExpr struct {
	Exprs []Expr
}

func (a AndExpr) Eval(env map[string]interface{}) (interface{}, error) {
	for _, expr := range a.Exprs {
		r, err := expr.Eval(env)
		if err != nil {
			return false, err
		}
		if r == nil {
			return false, nil
		}
		val, ok := r.(bool)
		if !ok {
			return false, nil
		}
		if !val {
			return false, nil
		}
	}

	return true, nil
}

func (a AndExpr) Equal(other Expr) bool {
	otherAnd, ok := other.(AndExpr)
	if !ok {
		return false
	}

	if len(a.Exprs) != len(otherAnd.Exprs) {
		return false
	}

	for i, expr := range a.Exprs {
		if !expr.Equal(otherAnd.Exprs[i]) {
			return false
		}
	}

	return true
}

// ----------------------------------------------------------------------------
// OrExpr
// ----------------------------------------------------------------------------

type OrExpr struct {
	Exprs []Expr
}

func (o OrExpr) Eval(env map[string]interface{}) (interface{}, error) {
	for _, expr := range o.Exprs {
		r, err := expr.Eval(env)
		if err != nil {
			return false, err
		}
		if r == nil {
			return false, nil
		}
		val, ok := r.(bool)
		if !ok {
			return false, nil
		}
		if val {
			return true, nil
		}
	}

	return false, nil
}

func (o OrExpr) Equal(other Expr) bool {
	otherOr, ok := other.(OrExpr)
	if !ok {
		return false
	}

	if len(o.Exprs) != len(otherOr.Exprs) {
		return false
	}

	for i, expr := range o.Exprs {
		if !expr.Equal(otherOr.Exprs[i]) {
			return false
		}
	}

	return true
}

// ----------------------------------------------------------------------------
// VariableRefExpr
// ----------------------------------------------------------------------------

// An expression that represents a variable reference.
type VariableRefExpr struct {
	// The name of the referenced variable
	Name string
}

func (v VariableRefExpr) Eval(params map[string]interface{}) (interface{}, error) {
	if _, ok := params[v.Name]; !ok {
		return nil, fmt.Errorf("variable %s not found", v.Name)
	}
	return params[v.Name], nil
}

func (v VariableRefExpr) Equal(other Expr) bool {
	otherValue, ok := other.(VariableRefExpr)
	if !ok {
		return false
	}

	return v.Name == otherValue.Name
}

// ----------------------------------------------------------------------------
// StructFieldRefExpr
// ----------------------------------------------------------------------------

// An expression that represents a struct field reference.
type StructFieldRefExpr struct {
	// The name of the referenced variable
	VarName string
	// The name of the field
	FieldName string
}

func (s StructFieldRefExpr) Eval(params map[string]interface{}) (interface{}, error) {
	if _, ok := params[s.VarName]; !ok {
		return nil, fmt.Errorf("variable %s not found", s.VarName)
	}

	return GetField(params[s.VarName], s.FieldName)
}

func (s StructFieldRefExpr) Equal(other Expr) bool {
	otherValue, ok := other.(StructFieldRefExpr)
	if !ok {
		return false
	}

	return s.VarName == otherValue.VarName && s.FieldName == otherValue.FieldName
}
