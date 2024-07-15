package authz

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// The ExprParser is capable of parsing expressions from a string.
type ExprParser struct{}

func (ep ExprParser) Parse(expr string) (Expr, error) {
	if len(expr) == 0 {
		return nil, errors.New("unexpected end of input")
	}

	parsed, consumed, err := ep.parseExpr(expr)
	if err != nil {
		return nil, err
	}

	// If we fail to consume the entire input, return an error.
	if consumed != len(expr) {
		return nil, errors.New("unexpected token")
	}

	return parsed, nil
}

// Parse an expression.
func (ep ExprParser) parseExpr(expr string) (Expr, int, error) {
	token, err := ep.nextToken(expr)
	if err != nil {
		return nil, 0, err
	}

	switch token {
	case "$eq":
		return ep.parseEqExpr(expr)
	case "$in":
		return ep.parseInExpr(expr)
	case "$and":
		return ep.parseAndExpr(expr)
	case "$or":
		return ep.parseOrExpr(expr)
	default:
		return ep.parseNonOperator(expr)
	}
}

// Get the next operator token from the expression.
func (ep ExprParser) nextToken(expr string) (string, error) {
	if len(expr) == 0 {
		return "", errors.New("unexpected end of input")
	}

	for i, c := range expr {
		if isTokenTerminator(c) {
			return expr[:i], nil
		}
		if i == len(expr)-1 {
			return expr, nil
		}
	}

	return "", errors.New("expected ' ', ',', '(', or end of input")
}

// ----------------------------------------------------------------------------
// Operator Expressions
// ----------------------------------------------------------------------------

// Parse an equality expression.
func (ep ExprParser) parseEqExpr(expr string) (EqExpr, int, error) {
	precondition(len(expr) > len("$eq"))

	ok, consumed := expectPrefix(expr, "$eq(")
	if !ok {
		return EqExpr{}, 0, errors.New("expected '$eq('")
	}

	left, right, n, err := ep.parseExpressionPair(expr[consumed:])
	if err != nil {
		return EqExpr{}, 0, err
	}
	consumed += n

	if len(expr[consumed:]) == 0 {
		return EqExpr{}, 0, errors.New("unexpected end of input")
	}

	// Consume the closing parenthesis
	if expr[consumed] != ')' {
		return EqExpr{}, 0, errors.New("expected ')'")
	}
	consumed++

	return EqExpr{Left: left, Right: right}, consumed, nil
}

// Parse an $in expression.
func (ep ExprParser) parseInExpr(expr string) (InExpr, int, error) {
	precondition(len(expr) > len("$in"))

	ok, consumed := expectPrefix(expr, "$in(")
	if !ok {
		return InExpr{}, 0, errors.New("expected '$in('")
	}

	query, collection, n, err := ep.parseExpressionPair(expr[consumed:])
	if err != nil {
		return InExpr{}, 0, err
	}
	consumed += n

	if len(expr[consumed:]) == 0 {
		return InExpr{}, 0, errors.New("unexpected end of input")
	}

	// Consume the closing parenthesis
	if expr[consumed] != ')' {
		return InExpr{}, 0, errors.New("expected ')'")
	}
	consumed++

	return InExpr{Element: query, Collection: collection}, consumed, nil
}

// Parse an AND expression.
func (ep ExprParser) parseAndExpr(expr string) (AndExpr, int, error) {
	precondition(len(expr) > len("$and"))

	ok, consumed := expectPrefix(expr, "$and(")
	if !ok {
		return AndExpr{}, 0, errors.New("expected '$and('")
	}

	exprs, n, err := ep.parseExpressionSequence(expr[consumed:], ')', func(token string) (Expr, int, error) {
		return ep.parseExpr(token)
	})
	if err != nil {
		return AndExpr{}, 0, err
	}
	consumed += n

	return AndExpr{Exprs: exprs}, consumed, nil
}

// Parse an OR expression.
func (ep ExprParser) parseOrExpr(expr string) (OrExpr, int, error) {
	precondition(len(expr) > len("$or"))

	ok, consumed := expectPrefix(expr, "$or(")
	if !ok {
		return OrExpr{}, 0, errors.New("expected '$or('")
	}

	exprs, n, err := ep.parseExpressionSequence(expr[consumed:], ')', func(token string) (Expr, int, error) {
		return ep.parseExpr(token)
	})
	if err != nil {
		return OrExpr{}, 0, err
	}
	consumed += n

	return OrExpr{Exprs: exprs}, consumed, nil
}

// Expect the specified prefix.
func expectPrefix(expr string, prefix string) (bool, int) {
	if len(expr) < len(prefix) {
		return false, 0
	}
	if expr[:len(prefix)] != prefix {
		return false, 0
	}
	return true, len(prefix)
}

// ----------------------------------------------------------------------------
// Non-operator Expressions
// ----------------------------------------------------------------------------

// Try to parse a non-operator expression.
func (ep ExprParser) parseNonOperator(expr string) (Expr, int, error) {
	precondition(len(expr) > 0)

	if len(expr) >= len("true") && expr[:len("true")] == "true" {
		// 'true' literal
		return ep.parseTrueExpr(expr)
	} else if len(expr) >= len("false") && expr[:len("false")] == "false" {
		// 'false' literal
		return ep.parseFalseExpr(expr)
	} else if expr[0] == '\'' {
		// String literal
		return ep.parseStrExpr(expr)
	} else if unicode.IsDigit(rune(expr[0])) {
		// Integer literal
		return ep.parseUintExpr(expr)
	} else if expr[0] == '[' {
		// Slice literal
		if len(expr) < len("[]") || expr[:len("[]")] != "[]" {
			return nil, 0, errors.New("expected '[]'")
		}

		if len(expr) > len("[]bool") && expr[:len("[]bool")] == "[]bool" {
			return ep.parseBoolSliceExpr(expr)
		} else if len(expr) > len("[]str") && expr[:len("[]str")] == "[]str" {
			return ep.parseStrSliceExpr(expr)
		} else if len(expr) > len("[]uint") && expr[:len("[]uint")] == "[]uint" {
			return ep.parseUintSliceExpr(expr)
		} else {
			return nil, 0, errors.New("expected 'bool', 'string', or 'uint' for slice literal")
		}
	} else {
		return ep.parseNonLiteral(expr)
	}
}

// Parse a non-literal, non-operator expression.
func (ep ExprParser) parseNonLiteral(expr string) (Expr, int, error) {
	precondition(len(expr) > 0)

	token, err := ep.nextToken(expr)
	if err != nil {
		return nil, 0, err
	}

	if strings.Contains(token, ".") {
		// Struct field reference
		return ep.parseStructFieldRefExpr(expr)
	} else {
		// Vanilla variable reference
		return ep.parseVariableRefExpr(expr)
	}
}

// Parse a boolean 'true' literal expression.
func (ep ExprParser) parseTrueExpr(expr string) (TrueExpr, int, error) {
	if len(expr) < len("true") {
		return TrueExpr{}, 0, errors.New("unexpected end of input")
	}
	if expr[:len("true")] != "true" {
		return TrueExpr{}, 0, errors.New("expected 'true'")
	}

	return TrueExpr{}, len("true"), nil
}

// Parse a boolean 'false' literal expression.
func (ep ExprParser) parseFalseExpr(expr string) (FalseExpr, int, error) {
	if len(expr) < len("false") {
		return FalseExpr{}, 0, errors.New("unexpected end of input")
	}
	if expr[:len("false")] != "false" {
		return FalseExpr{}, 0, errors.New("expected 'false'")
	}

	return FalseExpr{}, len("false"), nil
}

// Parse a string literal expression
func (ep ExprParser) parseStrExpr(expr string) (StrExpr, int, error) {
	precondition(len(expr) > 0)

	if expr[0] != '\'' {
		return StrExpr{}, 0, errors.New("expected '")
	}

	for i, c := range expr[1:] {
		if c == '\'' {
			return StrExpr{Value: expr[1 : i+1]}, i + 2, nil
		}
	}

	return StrExpr{}, 0, errors.New("expected closing '")
}

// Parse a uint literal expression.
func (ep ExprParser) parseUintExpr(expr string) (UintExpr, int, error) {
	for i, c := range expr {
		if isTokenTerminator(c) {
			v, err := strconv.ParseUint(expr[:i], 10, 32)
			if err != nil {
				return UintExpr{}, 0, errors.New("invalid integer literal")
			}
			return UintExpr{Value: uint(v)}, i, nil
		}

		// The literal consumed the entire raw expression
		if i == len(expr)-1 {
			if !unicode.IsDigit(c) {
				return UintExpr{}, 0, errors.New("expected digit")
			} else {
				v, err := strconv.ParseUint(expr, 10, 32)
				if err != nil {
					return UintExpr{}, 0, errors.New("invalid integer literal")
				}
				return UintExpr{Value: uint(v)}, len(expr), nil
			}
		}

		if !unicode.IsDigit(c) {
			return UintExpr{}, 0, errors.New("expected digit")
		}
	}

	return UintExpr{}, len(expr), nil
}

// Parse a boolean slice literal expression.
func (ep ExprParser) parseBoolSliceExpr(expr string) (BoolSliceExpr, int, error) {
	precondition(len(expr) > len("[]bool"))

	consumed := len("[]bool")

	// Consume the opening brace
	if expr[consumed] != '{' {
		return BoolSliceExpr{}, 0, errors.New("expected '{'")
	}
	consumed++

	cb := func(token string) (Expr, int, error) {
		if token == "true" {
			return ep.parseTrueExpr(token)
		} else if token == "false" {
			return ep.parseFalseExpr(token)
		} else {
			return BoolSliceExpr{}, 0, errors.New("expected 'true' or 'false' for boolean literal")
		}
	}

	exprs, n, err := ep.parseExpressionSequence(expr[consumed:], '}', cb)
	if err != nil {
		return BoolSliceExpr{}, 0, err
	}
	consumed += n

	return BoolSliceExpr{Values: exprs}, consumed, nil
}

// Parse a string slice literal expression.
func (ep ExprParser) parseStrSliceExpr(expr string) (StrSliceExpr, int, error) {
	precondition(len(expr) > len("[]str"))

	consumed := len("[]str")

	// Consume the opening brace
	if expr[consumed] != '{' {
		return StrSliceExpr{}, 0, errors.New("expected '{'")
	}
	consumed++

	cb := func(token string) (Expr, int, error) {
		return ep.parseStrExpr(token)
	}

	exprs, n, err := ep.parseExpressionSequence(expr[consumed:], '}', cb)
	if err != nil {
		return StrSliceExpr{}, 0, err
	}
	consumed += n

	return StrSliceExpr{Values: exprs}, consumed, nil
}

// Parse a uint slice literal expression.
func (ep ExprParser) parseUintSliceExpr(expr string) (UintSliceExpr, int, error) {
	precondition(len(expr) > len("[]uint"))
	consumed := len("[]uint")

	// Consume the opening brace
	if expr[consumed] != '{' {
		return UintSliceExpr{}, 0, errors.New("expected '{'")
	}
	consumed++

	cb := func(token string) (Expr, int, error) {
		return ep.parseUintExpr(token)
	}

	exprs, n, err := ep.parseExpressionSequence(expr[consumed:], '}', cb)
	if err != nil {
		return UintSliceExpr{}, 0, err
	}
	consumed += n

	return UintSliceExpr{Values: exprs}, consumed, nil
}

// Parse a variable ref expression.
func (ep ExprParser) parseVariableRefExpr(expr string) (VariableRefExpr, int, error) {
	for i, c := range expr {
		if isTokenTerminator(c) {
			return VariableRefExpr{Name: expr[:i]}, i, nil
		} else if i == len(expr)-1 {
			return VariableRefExpr{Name: expr[:i+1]}, i + 1, nil
		}
	}

	return VariableRefExpr{}, 0, errors.New("unexpected end of input")
}

// Parse a struct field reference expression.
func (ep ExprParser) parseStructFieldRefExpr(expr string) (StructFieldRefExpr, int, error) {
	var variable string
	for i, c := range expr {
		if isTokenTerminator(c) {
			variable = expr[:i]
			break
		} else if i == len(expr)-1 {
			variable = expr
		}
	}

	parts := strings.Split(variable, ".")
	if len(parts) != 2 {
		return StructFieldRefExpr{}, 0, fmt.Errorf("invalid struct field reference: %s", variable)
	}

	return StructFieldRefExpr{VarName: parts[0], FieldName: parts[1]}, len(variable), nil
}

// ----------------------------------------------------------------------------
// Parsing Utilities
// ----------------------------------------------------------------------------

// Parse a sequence of expressions separated by commas.
func (ep ExprParser) parseExpressionSequence(expr string, terminator byte, tokenCb func(string) (Expr, int, error)) ([]Expr, int, error) {
	consumed := 0

	exprs := make([]Expr, 0)
	vTerm := true
	for {
		if consumed >= len(expr) {
			return nil, 0, errors.New("unexpected end of input")
		}

		// Consume the closing brace
		if expr[consumed] == terminator {
			if vTerm {
				return exprs, consumed + 1, nil
			} else {
				return nil, 0, errors.New("invalid terminator for expression sequence")
			}
		}

		// Consume commas
		if expr[consumed] == ',' {
			vTerm = false
			consumed++
			continue
		}

		// Consume whitespace
		if expr[consumed] == ' ' {
			vTerm = true
			consumed++
			continue
		}

		token, err := ep.nextToken(expr[consumed:])
		if err != nil {
			return nil, 0, err
		}

		newExpr, n, err := tokenCb(token)
		if err != nil {
			return nil, 0, err
		}
		consumed += n

		exprs = append(exprs, newExpr)
		vTerm = true
	}
}

// Parse a pair of expressions, separated by comma and optional whitespace
func (ep ExprParser) parseExpressionPair(expr string) (Expr, Expr, int, error) {
	left, n, err := ep.parseExpr(expr)
	if err != nil {
		return nil, nil, 0, err
	}

	consumed := n
	if len(expr[consumed:]) == 0 {
		return nil, nil, 0, errors.New("unexpected end of input")
	}

	// Consume the comma
	if expr[consumed] != ',' {
		return nil, nil, 0, errors.New("expected ','")
	}
	consumed++

	// Consume whitespace
	for consumed < len(expr) && expr[consumed] == ' ' {
		consumed++
	}

	right, n, err := ep.parseExpr(expr[consumed:])
	if err != nil {
		return nil, nil, 0, err
	}
	consumed += n

	return left, right, consumed, nil
}

// Determine if a character is a token terminator.
func isTokenTerminator(c rune) bool {
	return c == ' ' || c == ',' || c == ')' || c == '(' || c == '}'
}

// Assert a precondition.
func precondition(pred bool) {
	if !pred {
		panic("precondition failed")
	}
}
