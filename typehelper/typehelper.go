// Package typehelper provides with some type conversion utility functions.
package typehelper

import (
	"reflect"
)

// ToSlice converts a `interface{}` to a `[]interface{}`.
// It returns the result slice.
// If arg is not a slice it panic.
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
