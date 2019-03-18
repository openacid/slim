// Package array implements functions for the manipulation of compacted array.
package array

import (
	"reflect"

	"github.com/openacid/slim/marshal"
)

// Array is a space efficient array implementation.
//
// Unlike a normal array, it does not allocate space for a element that there is
// not data in it.
type Array struct {
	ArrayBase
}

// NewEmpty creates an empty array with element of type of "v".
// If v is a pointer, the value type it points to is used.
// "v" must be a fixed size type, e.g.:
//		int32
//		struct { X int32; Y uint16 }
// "v" can not be:
//		int
//		[]uint32
//		map
//		etc.
func NewEmpty(v interface{}) (*Array, error) {
	m, err := marshal.NewTypeMarshaler(v)
	if err != nil {
		return nil, err
	}

	a := &Array{}
	a.EltMarshaler = m
	return a, nil
}

// New creates an array from indexes and elts.
// Length of indexes and length of elts must be the same.
// elts must be a slice of value of fixed size.
func New(indexes []int32, elts interface{}) (*Array, error) {
	a := &Array{}
	err := a.Init(indexes, elts)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// Init initializes an empty Array.
// Length of indexes and length of elts must be the same.
// elts must be a slice of value of fixed size.
func (a *Array) Init(indexes []int32, elts interface{}) error {
	err := a.ArrayBase.Init(indexes, elts)
	if err != nil {
		return err
	}

	// Only when inited with some elements, we init the Marshaler
	if a.Cnt > 0 && a.EltMarshaler == nil {

		v := reflect.ValueOf(elts).Index(0)
		marshaler, err := marshal.NewTypeMarshalerEndian(v.Interface(), endian)
		if err != nil {
			return err
		}

		a.EltMarshaler = marshaler
	}

	return nil
}
