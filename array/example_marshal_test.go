package array_test

import (
	"fmt"

	"github.com/openacid/slim/array"
)

func Example_marshal() {

	// This example shows how to marshal / unmarshal an array

	indexes := []uint32{1, 5, 9, 203}
	elts := []uint16{12, 15, 19, 120}

	a := &array.ArrayU16{
		Data: elts,
	}
	err := a.InitIndexBitmap(indexes)
	if err != nil {
		return
	}

	// Marshal
	md, err := array.Marshal(a)
	fmt.Println("marshal: err:", err)
	fmt.Println("marshaled bytes:", md)

	// Unmarshal
	b := &array.ArrayU16{}
	n, err := array.Unmarshal(b, md)

	fmt.Println("unmarshal result:", n, err)
	fmt.Println("a:", a.Cnt, a.Bitmaps, a.Offsets, a.Data)
	fmt.Println("b:", b.Cnt, b.Bitmaps, b.Offsets, b.Data)

	// Output:
	// marshal: err: <nil>
	// marshaled bytes: [8 4 18 6 162 4 0 0 128 16 26 4 0 0 0 3 34 8 12 0 15 0 19 0 120 0]
	// unmarshal result: 26 <nil>
	// a: 4 [546 0 0 2048] [0 0 0 3] [12 15 19 120]
	// b: 4 [546 0 0 2048] [0 0 0 3] [12 15 19 120]
}
