package authz

import (
	"errors"
	"fmt"
	"reflect"
)

// Get a field from an object by name.
func GetField(obj interface{}, name string) (interface{}, error) {
	if !hasKind(obj, []reflect.Kind{reflect.Struct, reflect.Pointer}) {
		return nil, errors.New("unsupported type")
	}

	objValue := reflectValue(obj)

	field := objValue.FieldByName(name)
	if !field.IsValid() {
		return nil, fmt.Errorf("no such field: %s in obj", name)
	}

	return field.Interface(), nil
}

// Reflect a value from an object.
func reflectValue(obj interface{}) reflect.Value {
	if reflect.TypeOf(obj).Kind() == reflect.Ptr {
		return reflect.ValueOf(obj).Elem()
	} else {
		return reflect.ValueOf(obj)
	}
}

// Determine if a value has one of a set of Kinds.
func hasKind(obj interface{}, types []reflect.Kind) bool {
	for _, t := range types {
		if reflect.TypeOf(obj).Kind() == t {
			return true
		}
	}

	return false
}
