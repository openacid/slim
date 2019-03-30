package array_test

import (
	"fmt"

	"github.com/openacid/slim/array"
)

type Combined struct {
	I uint8
	J uint8
}

type MyEncoder struct{}

func (c MyEncoder) Encode(d interface{}) []byte {
	v := d.(Combined)
	return []byte{v.I, v.J}
}

func (c MyEncoder) Decode(b []byte) (int, interface{}) {
	return 2, Combined{
		I: b[0],
		J: b[1],
	}
}

func (c MyEncoder) GetSize(d interface{}) int {
	return 2
}

func (c MyEncoder) GetEncodedSize(b []byte) int {
	return 2
}

func Example_withEncoder() {

	// This example shows how to define a new array type.

	indexes := []int32{1, 5, 9, 203}
	elts := []Combined{
		{1, 2},
		{3, 4},
		{6, 6},
		{7, 8},
	}

	a := &array.Array{}
	a.EltEncoder = MyEncoder{}
	err := a.Init(indexes, elts)
	if err != nil {
		panic(err)
	}

	for _, i := range []int32{1, 2, 5, 6} {
		d, found := a.Get(i)
		if found {
			fmt.Printf("value of a[%d]: %d\n", i, d)
		} else {
			fmt.Printf("value of a[%d] does not exist\n", i)
		}
	}

	// Output:
	// value of a[1]: {1 2}
	// value of a[2] does not exist
	// value of a[5]: {3 4}
	// value of a[6] does not exist
}
