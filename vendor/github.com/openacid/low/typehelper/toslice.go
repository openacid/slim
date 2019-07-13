package typehelper

import (
	"reflect"
)

// ToSlice converts an `interface{}` to a `[]interface{}`.
// It returns the result slice.
// If arg is not a slice it panics.
//
// Since 0.1.10
func ToSlice(arg interface{}) []interface{} {

	s := reflect.ValueOf(arg)
	if s.Kind() != reflect.Slice {
		panic("not a slice")
	}

	l := s.Len()
	rst := make([]interface{}, l)
	for i := 0; i < l; i++ {
		rst[i] = s.Index(i).Interface()
	}
	return rst
}
