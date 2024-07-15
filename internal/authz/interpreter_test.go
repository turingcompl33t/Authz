package authz

import "testing"

// Interpreter can interpret boolean-valued expressions.
func TestBool(t *testing.T) {
	data := []struct {
		input     string
		params    map[string]interface{}
		want      bool
		expectErr error
	}{
		{"true", nil, true, nil},
		{"false", nil, false, nil},
		{"$and(true, false)", nil, false, nil},
		{"$and(true, true)", nil, true, nil},
		{"$or(true, false)", nil, true, nil},
		{"$or(false, false)", nil, false, nil},
		{"value", map[string]interface{}{"value": true}, true, nil},
		{"value", map[string]interface{}{"value": false}, false, nil},
		{"obj.Value", map[string]interface{}{"obj": struct{ Value bool }{true}}, true, nil},
		{"obj.Value", map[string]interface{}{"obj": struct{ Value bool }{false}}, false, nil},
		{"$eq(value, true)", map[string]interface{}{"value": true}, true, nil},
		{"$eq(value, false)", map[string]interface{}{"value": false}, true, nil},
		{"$in('foo', []str{'foo' , 'bar'})", nil, true, nil},
		{"$in('baz', []str{'foo' , 'bar'})", nil, false, nil},
	}

	for _, d := range data {
		i := Interpreter{}

		got, err := i.Bool(d.input, d.params)
		if err != d.expectErr {
			t.Errorf("unexpected error: got=%v, want=%v", err, d.expectErr)
		}
		if got != d.want {
			t.Errorf("unexpected result: got=%v, want=%v", got, d.want)
		}
	}
}
