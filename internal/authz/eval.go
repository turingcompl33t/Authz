package authz

// The Evaluator is responsible for evaluating expressions.
type Evaluator struct{}

func (i Evaluator) Eval(expr Expr, params map[string]interface{}) (interface{}, error) {
	return expr.Eval(params)
}
