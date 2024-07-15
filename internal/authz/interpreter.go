package authz

import "fmt"

// The Interpreter is responsible for evaluating expressions.
type Interpreter struct{}

// Evaluate an expression with the given parameters.
func (i Interpreter) Eval(expr string, params map[string]interface{}) (interface{}, error) {
	p := ExprParser{}
	parsed, err := p.Parse(expr)
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	ev := Evaluator{}
	result, err := ev.Eval(parsed, params)
	if err != nil {
		return nil, fmt.Errorf("evaluation error: %w", err)
	}

	return result, nil
}

// Evaluate a boolean-valued expression with the given parameters.
func (i Interpreter) Bool(expr string, params map[string]interface{}) (bool, error) {
	result, err := i.Eval(expr, params)
	if err != nil {
		return false, err
	}

	r, err := coerceBool(result)
	if err != nil {
		return false, fmt.Errorf("failed to coerce expression result to bool: %w", err)
	}

	return r, nil
}

// Evaluate a string-valued expression with the given parameters.
func (i Interpreter) Str(expr string, params map[string]interface{}) (string, error) {
	result, err := i.Eval(expr, params)
	if err != nil {
		return "", err
	}

	r, err := coerceStr(result)
	if err != nil {
		return "", fmt.Errorf("failed to coerce expression result to string: %w", err)
	}

	return r, nil
}

// Evaluate an integer-valued expression with the given parameters.
func (i Interpreter) Uint(expr string, params map[string]interface{}) (uint, error) {
	result, err := i.Eval(expr, params)
	if err != nil {
		return 0, err
	}

	r, err := coerceUint(result)
	if err != nil {
		return 0, fmt.Errorf("failed to coerce expression result to uint: %w", err)
	}

	return r, nil
}
