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
	Get(i int32) int64
	// Len returns number of elements in this array.
	Len() int
	// Stat returns a stat map that describes memory usage or else.
	Stat() map[string]int64

	proto.Message
}
