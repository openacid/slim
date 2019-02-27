// Package typehelper provides with some type conversion utility functions.
package typehelper

import (
	"reflect"
)

// ToSlice converts a `interface{}` to a `[]interface{}`.
// It returns the result slice an `ok` indicating if it is a slice.
func ToSlice(arg interface{}) (rst []interface{}, ok bool) {

	s := reflect.ValueOf(arg)
	if s.Kind() != reflect.Slice {
		return nil, false
	}

	l := s.Len()
	rst = make([]interface{}, l)
	for i := 0; i < l; i++ {
		rst[i] = s.Index(i).Interface()
	}
	return rst, true
}
