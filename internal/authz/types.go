package authz

import (
	"errors"
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

// Attempt to coerce a value to a slice of string.
func coerceStrSlice(v interface{}) ([]string, error) {
	if s, ok := v.([]string); ok {
		return s, nil
	}
	return nil, errors.New("failed to coerce slice")

}

// Attempt to coerce a value to a uint.
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

// Attempt to coerce a slice to a slice of uint.
func coerceUintSlice(v interface{}) ([]uint, error) {
	switch s := v.(type) {
	case []uint:
		return s, nil
	case []uint8:
		return transform(s, func(x uint8) uint { return uint(x) }), nil
	case []uint16:
		return transform(s, func(x uint16) uint { return uint(x) }), nil
	case []uint32:
		return transform(s, func(x uint32) uint { return uint(x) }), nil
	case []uint64:
		return transform(s, func(x uint64) uint { return uint(x) }), nil
	case []int:
		return transformWithCheck(s, func(x int) uint { return uint(x) }, func(x int) bool { return x >= 0 })
	case []int8:
		return transformWithCheck(s, func(x int8) uint { return uint(x) }, func(x int8) bool { return x >= 0 })
	case []int16:
		return transformWithCheck(s, func(x int16) uint { return uint(x) }, func(x int16) bool { return x >= 0 })
	case []int32:
		return transformWithCheck(s, func(x int32) uint { return uint(x) }, func(x int32) bool { return x >= 0 })
	case []int64:
		return transformWithCheck(s, func(x int64) uint { return uint(x) }, func(x int64) bool { return x >= 0 })
	default:
		return nil, errors.New("failed to coerce slice")
	}
}

// Attempt to coerce a value to a bool.
func coerceBool(v interface{}) (bool, error) {
	if b, ok := v.(bool); ok {
		return b, nil
	}
	return false, fmt.Errorf("expected bool, got %v", reflect.TypeOf(v))
}

// Attempt to coerce a value to a slice of bool.
func coerceBoolSlice(v interface{}) ([]bool, error) {
	if s, ok := v.([]bool); ok {
		return s, nil
	}
	return nil, errors.New("failed to coerce slice")
}

// Transform a slice.
func transform[T any, U any](slice []T, fn func(T) U) []U {
	result := make([]U, 0, len(slice))
	for _, element := range slice {
		result = append(result, fn(element))
	}
	return result
}

// Transform a slice with a check.
func transformWithCheck[T any, U any](slice []T, fn func(T) U, check func(T) bool) ([]U, error) {
	result := make([]U, 0, len(slice))
	for _, element := range slice {
		if !check(element) {
			return nil, errors.New("check failed")
		}
		result = append(result, fn(element))
	}
	return result, nil
}
