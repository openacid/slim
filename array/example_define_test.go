package array_test

import (
	"fmt"

	"github.com/openacid/slim/array"
)

type MyElt struct {
	i uint32
	f uint64
}

type MyArray struct {
	array.Array32Index
	Data []MyElt
}

// Get2 implements user defined data retrieving routine.
func (a *MyArray) Get2(i uint32) (MyElt, bool) {
	j, found := a.GetEltIndex(i)
	if found {
		return a.Data[j], true
	}
	return MyElt{}, false
}

func Example_defineArray() {

	// This exmaple shows how to define a new array type.

	indexes := []uint32{1, 5, 9, 203}
	elts := []MyElt{
		{1, 2},
		{3, 4},
		{6, 6},
		{7, 8},
	}

	a := &MyArray{Data: elts}

	err := a.InitIndexBitmap(indexes)
	if err != nil {
		return
	}

	for _, i := range []uint32{1, 2, 5, 6} {
		d, found := a.Get2(i)
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
