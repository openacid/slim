// Package typehelper provides with some type conversion utility functions.
package typehelper

import "github.com/openacid/low/typehelper"

// ToSlice converts a `interface{}` to a `[]interface{}`.
// It returns the result slice.
// If arg is not a slice it panic.
//
// Deprecated: use github.com/openacid/low/typehelper
// Deprecated: will be removed since 1.0.0
func ToSlice(arg interface{}) []interface{} {
	return typehelper.ToSlice(arg)
}
