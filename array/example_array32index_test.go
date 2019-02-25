package array_test

import (
	"fmt"

	"github.com/openacid/slim/array"
)

type ArrayU16 struct {
	array.Array32Index
	data []uint16
}

func New(index []uint32, elts []uint16) (a *ArrayU16, err error) {
	a = &ArrayU16{
		data: elts,
	}

	err = a.InitIndexBitmap(index)
	if err != nil {
		a = nil
	}

	return
}

func ExampleArray32Index() {

	index := []uint32{1, 5, 9, 203}
	elts := []uint16{12, 15, 19, 120}

	a, err := New(index, elts)
	if err != nil {
		return
	}

	// Check if an index present
	fmt.Println("a[0] present:", a.Has(0))
	fmt.Println("a[1] present:", a.Has(1))

	// Get element
	i, has := a.GetEltIndex(1)
	if has {
		fmt.Println("value of a[1]:", a.data[i])
	} else {
		fmt.Println("value of a[1] does not exist")
	}

	// Output:
	// a[0] present: false
	// a[1] present: true
	// value of a[1]: 12

}
