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
