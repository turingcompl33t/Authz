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

func (ep ExprParser) parseExpr(expr string) (Expr, int, error) {
	token, err := ep.nextToken(expr)
	if err != nil {
		return nil, 0, err
	}

	switch token {
	case "$eq":
		return ep.parseEqExpr(expr)
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
		if isTokenTerminator(c) || i == len(expr)-1 {
			return expr[:i], nil
		}
	}

	return "", errors.New("expected ' ', ',', '(', or end of input")
}

// ----------------------------------------------------------------------------
// Operator Expressions
// ----------------------------------------------------------------------------

func (ep ExprParser) parseEqExpr(expr string) (EqExpr, int, error) {
	precondition(len(expr) > len("$eq"))

	if len(expr) < len("$eq(") || expr[:len("$eq(")] != "$eq(" {
		return EqExpr{}, 0, errors.New("expected '$eq('")
	}

	consumed := len("$eq(")

	left, n, err := ep.parseExpr(expr[consumed:])
	if err != nil {
		return EqExpr{}, 0, err
	}

	consumed += n
	if len(expr[consumed:]) == 0 {
		return EqExpr{}, 0, errors.New("unexpected end of input")
	}

	// Consume the comma
	if expr[consumed] != ',' {
		return EqExpr{}, 0, errors.New("expected ','")
	}
	consumed++

	// Consume whitespace
	for consumed < len(expr) && expr[consumed] == ' ' {
		consumed++
	}

	right, n, err := ep.parseExpr(expr[consumed:])
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
		return ep.parseStringExpr(expr)
	} else if unicode.IsDigit(rune(expr[0])) {
		fmt.Println("parsing int")
		// Integer literal
		return ep.parseIntExpr(expr)
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
func (ep ExprParser) parseStringExpr(expr string) (StringExpr, int, error) {
	precondition(len(expr) > 0)
	precondition(expr[0] == '\'')

	for i, c := range expr[1:] {
		if c == '\'' {
			return StringExpr{Value: expr[1 : i+1]}, i + 2, nil
		}
	}

	return StringExpr{}, 0, errors.New("expected closing '")
}

func (ep ExprParser) parseIntExpr(expr string) (IntExpr, int, error) {
	for i, c := range expr {
		if isTokenTerminator(c) {
			v, err := strconv.ParseUint(expr[:i], 10, 32)
			if err != nil {
				return IntExpr{}, 0, errors.New("invalid integer literal")
			}
			return IntExpr{Value: uint(v)}, i, nil
		}

		// The literal consumed the entire raw expression
		if i == len(expr)-1 {
			if !unicode.IsDigit(c) {
				return IntExpr{}, 0, errors.New("expected digit")
			} else {
				v, err := strconv.ParseUint(expr, 10, 32)
				if err != nil {
					return IntExpr{}, 0, errors.New("invalid integer literal")
				}
				return IntExpr{Value: uint(v)}, len(expr), nil
			}
		}

		if !unicode.IsDigit(c) {
			return IntExpr{}, 0, errors.New("expected digit")
		}
	}

	return IntExpr{}, len(expr), nil
}

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

// Determine if a character is a token terminator.
func isTokenTerminator(c rune) bool {
	return c == ' ' || c == ',' || c == ')' || c == '('
}

// Assert a precondition.
func precondition(pred bool) {
	if !pred {
		panic("precondition failed")
	}
}
