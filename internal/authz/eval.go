package authz

import "errors"

// The Evaluator is responsible for evaluating expressions.
type Evaluator struct{}

func (i Evaluator) Eval(expr Expr, params map[string]interface{}) (interface{}, error) {
	return expr.Eval(params)
}

func AsBool(v interface{}) (bool, error) {
	if b, ok := v.(bool); ok {
		return b, nil
	}

	return false, errors.New("expected boolean")
}

func AsInt(v interface{}) (uint, error) {
	switch v := v.(type) {
	case int:
		return uint(v), nil
	case uint:
		return v, nil
	default:
		return 0, errors.New("expected int")
	}
}

func AsString(v interface{}) (string, error) {
	if s, ok := v.(string); ok {
		return s, nil
	}

	return "", errors.New("expected string")
}
