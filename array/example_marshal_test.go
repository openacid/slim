package array_test

import (
	"fmt"

	proto "github.com/golang/protobuf/proto"
	"github.com/openacid/slim/array"
)

func Example_marshal() {

	// This example shows how to marshal / unmarshal an array

	indexes := []int32{1, 5, 9, 203}
	elts := []uint16{12, 15, 19, 120}

	a, err := array.NewU16(indexes, elts)
	if err != nil {
		return
	}

	// Marshal
	md, err := proto.Marshal(a)
	fmt.Println("marshal: err:", err)
	fmt.Println("marshaled bytes:", md)

	// Unmarshal
	b := &array.U16{}
	err = proto.Unmarshal(md, b)

	fmt.Println("unmarshal result:", err)
	fmt.Println("a:", a.Cnt, a.Bitmaps, a.Offsets, a.Elts)
	fmt.Println("b:", b.Cnt, b.Bitmaps, b.Offsets, b.Elts)

	// Output:
	// marshal: err: <nil>
	// marshaled bytes: [8 4 18 6 162 4 0 0 128 16 26 4 0 0 0 3 34 8 12 0 15 0 19 0 120 0]
	// unmarshal result: <nil>
	// a: 4 [546 0 0 2048] [0 0 0 3] [12 0 15 0 19 0 120 0]
	// b: 4 [546 0 0 2048] [0 0 0 3] [12 0 15 0 19 0 120 0]
}
