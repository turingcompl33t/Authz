package authz

import (
	"fmt"
	"reflect"
)

// Attempt to coerce a value to a string.
func coerceStr(v interface{}) (string, error) {
	if s, ok := v.(string); ok {
		return s, nil
	}
	return "", fmt.Errorf("expected string, got %v", reflect.TypeOf(v))
}

func coerceUint(v interface{}) (uint, error) {
	intHelper := func(v int) (uint, error) {
		if v < 0 {
			return 0, fmt.Errorf("expected uint, got negative int")
		}
		return uint(v), nil
	}

	switch v := v.(type) {
	case uint:
		return v, nil
	case uint8:
		return uint(v), nil
	case uint16:
		return uint(v), nil
	case uint32:
		return uint(v), nil
	case uint64:
		return uint(v), nil
	case int:
		return intHelper(v)
	case int8:
		return intHelper(int(v))
	case int16:
		return intHelper(int(v))
	case int32:
		return intHelper(int(v))
	case int64:
		return intHelper(int(v))
	default:
		return 0, fmt.Errorf("expected uint, got %v", reflect.TypeOf(v))
	}
}

// Attempt to coerce a value to a bool.
func coerceBool(v interface{}) (bool, error) {
	if b, ok := v.(bool); ok {
		return b, nil
	}
	return false, fmt.Errorf("expected bool, got %v", reflect.TypeOf(v))
}
