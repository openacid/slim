// Package array implements several space effiecient array.
//
// Unlike a normal array, it does not allocate space for an absent element.
// In some way it is like a map[int32]interface{} .
//
// Internally it use a bitmap to indicate at which index there is an element.
package array

import (
	"reflect"

	"github.com/openacid/slim/marshal"
)

// Array is a space efficient array of fixed-size element.
//
// A fixed size type could be:
//		int32
//		struct { X int32; Y uint16 }
//
// A non-fixed size type could be:
//		int
//		[]uint32
//		map
//		etc.
type Array struct {
	Base
}

// NewEmpty creates an empty Array with element of type of "v".
// If v is a pointer, the value type it points to is used.
func NewEmpty(v interface{}) (*Array, error) {
	m, err := marshal.NewTypeMarshaler(v)
	if err != nil {
		return nil, err
	}

	a := &Array{}
	a.EltMarshaler = m
	return a, nil
}

// New creates an array from specified indexes and elts.
// The length of indexes and the length of elts must be the same.
// "elts" must be a slice of fixed-size values.
func New(indexes []int32, elts interface{}) (*Array, error) {
	a := &Array{}
	err := a.Init(indexes, elts)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// Init initializes an Array.
// Length of "indexes" and length of "elts" must be the same.
// "elts" must be a slice of fixed-size value.
func (a *Array) Init(indexes []int32, elts interface{}) error {
	err := a.Base.Init(indexes, elts)
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
