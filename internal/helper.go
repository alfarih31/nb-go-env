package internal

import "reflect"

// HasZeroValue Check a variable has Zero Value
func HasZeroValue(v interface{}) bool {
	if v == nil {
		return true
	}

	t := reflect.TypeOf(v)
	if t == nil {
		return true
	}

	switch t.Kind() {
	case reflect.Map, reflect.Slice, reflect.Array:
		return false
	}

	return v == reflect.Zero(t).Interface()
}
