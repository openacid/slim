package array_test

import (
	"fmt"

	"github.com/openacid/slim/array"
)

type MyElt struct {
	I uint32
	J uint64
}

func Example_useArray() {

	// This exmaple shows how to define a new array type.

	indexes := []int32{1, 5, 9, 203}
	elts := []MyElt{
		{1, 2},
		{3, 4},
		{6, 6},
		{7, 8},
	}

	a, err := array.New(indexes, elts)
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
