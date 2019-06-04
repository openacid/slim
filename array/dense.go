package array

import proto "github.com/golang/protobuf/proto"

// Denser describes behaviors of a STATIC dense array.
//
// A dense array has no empty element(unlike array.U16 could have some empty
// element in it).
// Thus normally a dense array compress data in another way.
//
// Since 0.5.2
type Denser interface {
	// Get retrieve an element at i.
	Get(i int32) int32
	// Len returns number of elements in this array.
	Len() int
	// Stat returns a stat map that describes memory usage or else.
	Stat() map[string]int32

	proto.Message
}

// NewDense creates a "dense" array from a slice of int32.
//
// It is very efficient to store a serias integers with a overall trend, such as
// a sorted array.
//
// Since 0.5.2
func NewDense(nums []int32) Denser {
	return nil
}
