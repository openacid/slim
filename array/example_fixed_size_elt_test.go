package array_test

import (
	"encoding/binary"
	"fmt"

	"github.com/openacid/slim/array"
)

// Define a array

type MyElt uint32

type MyArray struct {
	array.Array
}

func (ma *MyArray) Get(i int32) (MyElt, bool) {
	bs, found := ma.GetBytes(i, 4)
	if !found {
		return MyElt(0), false
	}

	return MyElt(binary.LittleEndian.Uint32(bs)), true
}

func Example_defineArray() {

	a := &MyArray{}
	err := a.Init(
		[]int32{1, 5, 9, 203},
		[]MyElt{1, 2, 3, 4})

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
	// value of a[1]: 1
	// value of a[2] does not exist
	// value of a[5]: 2
	// value of a[6] does not exist
}
