// Package array implements several space effiecient array.
//
// Unlike a normal array, it does not allocate space for an absent element.
// In some way it is like a map[int32]interface{} .
//
// Memory overhead
//
// Internally it use a bitmap to indicate at which index there is an element.
//
// This implementation allocate 1 bit for every abscent or present element to
// inidicate if there is an element at this position.
// Thus the memory overhead is about 1 bit / load-factor.
//
//                    count(elements)
//      load-factor = ----------------
//                     max(indexes)
//
// An array with load factor = 0.5 requires about 2 extra bit pre present element.
//
// Implementation note
//
// - The bottom level Array32 is a protobuf message that defines in memory
// and on-disk structure.
//
// - The second level Base provides several basic methods such as mapping an index to its position in memory.
//
// - At the top level there are several ready to use implements. "Array"
// accepts any fixed-type value as element. Thus it is easy to use but not very
// efficient. "U32" accepts only uint32 as its element thus its performance is
// much better.
//
//      Array   U32    U64         // ready-to-use types
//        `----. | .----'
//             v v v
//             Base                // access supporting methods.
//              |
//              v
//      protobuf:Array32           // in-memory and on-disk structure.
//
// Performance note
//
// A Get involves at least 2 memory access to a.Bitmaps and a.Elts.
//
// An "Array" of general type requires one additional alloc for a Get:
//   // when Unmarshal convert a concrete type to interface{}
//   a.EltMarshaler.Unmarshal(bs)
//
// An array of specific type such as "U32" does not requires additional alloc.
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
//
// Performance note
//
// A general Array.Get is implemented with reflect.
// Benchmark shows a Get() costs ~ 168 ns/op and involves 4 alloc.
//
// To achieve the best performance, use a type specific array such as array.U32.
// The performance is much better: a Get() costs ~ 12 ns/op and involves 0
// alloc.
//
// Since 0.2.0
type Array struct {
	Base
}

// NewEmpty creates an empty Array with element of type of "v".
// If v is a pointer, the value type it points to is used.
//
// Since 0.2.0
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
//
// Since 0.2.0
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
// "elts" must be a slice of fixed-size values.
//
// By default Array encodes an element with binary.Write(), in binary.LittenEndian.
//
// Since 0.2.0
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
