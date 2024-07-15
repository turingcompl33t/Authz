package authz

import (
	"errors"
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
// StrExpr
// ----------------------------------------------------------------------------

// StrExpr represents a string literal.
type StrExpr struct {
	Value string
}

func (s StrExpr) Eval(env map[string]interface{}) (interface{}, error) {
	return s.Value, nil
}

func (s StrExpr) Equal(other Expr) bool {
	otherString, ok := other.(StrExpr)
	if !ok {
		return false
	}

	return s.Value == otherString.Value
}

// ----------------------------------------------------------------------------
// UintExpr
// ----------------------------------------------------------------------------

// UintExpr represents an integer literal.
type UintExpr struct {
	Value uint
}

func (i UintExpr) Eval(env map[string]interface{}) (interface{}, error) {
	return i.Value, nil
}

func (i UintExpr) Equal(other Expr) bool {
	otherInt, ok := other.(UintExpr)
	if !ok {
		return false
	}

	return i.Value == otherInt.Value
}

// ----------------------------------------------------------------------------
// BoolSliceExpr
// ----------------------------------------------------------------------------

// BoolSliceExpr represents a boolean slice literal.
type BoolSliceExpr struct {
	Values []Expr
}

func (b BoolSliceExpr) Eval(env map[string]interface{}) (interface{}, error) {
	var result []bool
	for _, expr := range b.Values {
		val, err := expr.Eval(env)
		if err != nil {
			return nil, err
		}
		boolVal, err := coerceBool(val)
		if err != nil {
			return nil, fmt.Errorf("unexpected type in boolean slice: %T", val)
		}
		result = append(result, boolVal)
	}
	return result, nil
}

func (b BoolSliceExpr) Equal(other Expr) bool {
	otherBoolSlice, ok := other.(BoolSliceExpr)
	if !ok {
		return false
	}

	if len(b.Values) != len(otherBoolSlice.Values) {
		return false
	}

	for i, expr := range b.Values {
		if !expr.Equal(otherBoolSlice.Values[i]) {
			return false
		}
	}

	return true
}

// ----------------------------------------------------------------------------
// StrSliceExpr
// ----------------------------------------------------------------------------

// Represents a string slice literal.
type StrSliceExpr struct {
	Values []Expr
}

func (s StrSliceExpr) Eval(env map[string]interface{}) (interface{}, error) {
	var result []string
	for _, expr := range s.Values {
		val, err := expr.Eval(env)
		if err != nil {
			return nil, err
		}
		strVal, err := coerceStr(val)
		if err != nil {
			return nil, fmt.Errorf("unexpected type in string slice: %T", val)
		}
		result = append(result, strVal)
	}
	return result, nil
}

func (s StrSliceExpr) Equal(other Expr) bool {
	otherStrSlice, ok := other.(StrSliceExpr)
	if !ok {
		return false
	}

	if len(s.Values) != len(otherStrSlice.Values) {
		return false
	}

	for i, expr := range s.Values {
		if !expr.Equal(otherStrSlice.Values[i]) {
			return false
		}
	}

	return true
}

// ----------------------------------------------------------------------------
// UintSliceExpr
// ----------------------------------------------------------------------------

// Represents a uint slice literal.
type UintSliceExpr struct {
	Values []Expr
}

func (u UintSliceExpr) Eval(env map[string]interface{}) (interface{}, error) {
	var result []uint
	for _, expr := range u.Values {
		val, err := expr.Eval(env)
		if err != nil {
			return nil, err
		}
		uintVal, err := coerceUint(val)
		if err != nil {
			return nil, fmt.Errorf("unexpected type in uint slice: %T", val)
		}
		result = append(result, uintVal)
	}
	return result, nil
}

func (u UintSliceExpr) Equal(other Expr) bool {
	otherUintSlice, ok := other.(UintSliceExpr)
	if !ok {
		return false
	}

	if len(u.Values) != len(otherUintSlice.Values) {
		return false
	}

	for i, expr := range u.Values {
		if !expr.Equal(otherUintSlice.Values[i]) {
			return false
		}
	}

	return true
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

		ok, err := truthy(r)
		if err != nil {
			return false, err
		}

		if !ok {
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

		ok, err := truthy(r)
		if err != nil {
			return false, err
		}

		if ok {
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

// ----------------------------------------------------------------------------
// InExpr
// ----------------------------------------------------------------------------

// An expression that represents the logic to determine if an element is a member of a slice.
type InExpr struct {
	Element    Expr
	Collection Expr
}

func (i InExpr) Eval(params map[string]interface{}) (interface{}, error) {
	// Resolve the collection
	sliceVal, err := i.Collection.Eval(params)
	if err != nil {
		return nil, err
	}

	// Resolve the element
	queryVal, err := i.Element.Eval(params)
	if err != nil {
		return nil, err
	}

	queryStr, err := coerceStr(queryVal)
	if err == nil {
		sliceStr, err := coerceStrSlice(sliceVal)
		if err == nil {
			for _, s := range sliceStr {
				if s == queryStr {
					return true, nil
				}
			}
			return false, nil
		} else {
			return nil, errors.New("mismatched types for $in()")
		}
	}

	queryUint, err := coerceUint(queryVal)
	if err == nil {
		sliceUint, err := coerceUintSlice(sliceVal)
		if err == nil {
			for _, u := range sliceUint {
				if u == queryUint {
					return true, nil
				}
			}
			return false, nil
		} else {
			return nil, errors.New("mismatched types for $in()")
		}
	}

	queryBool, err := coerceBool(queryVal)
	if err == nil {
		sliceBool, err := coerceBoolSlice(sliceVal)
		if err == nil {
			for _, u := range sliceBool {
				if u == queryBool {
					return true, nil
				}
			}
			return false, nil
		} else {
			return nil, errors.New("mismatched types for $in()")
		}
	}

	return nil, fmt.Errorf("unexpected type for $in() query")
}

func (i InExpr) Equal(other Expr) bool {
	otherIn, ok := other.(InExpr)
	if !ok {
		return false
	}
	return i.Element.Equal(otherIn.Element) && i.Collection.Equal(otherIn.Collection)
}

// // Determine if any element of the slice satisfies the predicate.
// func in[T comparable](e T, slice []T) bool {
// 	for _, element := range slice {
// 		if element == e {
// 			return true
// 		}
// 	}
// 	return false
// }
