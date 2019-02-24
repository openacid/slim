package array

import "fmt"

func ExampleArray32Index() {

	index := []uint32{1, 5, 9, 203}
	elts := []uint16{12, 15, 19, 120}

	type arrayU16 struct {
		Array32Index
		data []uint16
	}

	// Create array
	a := &arrayU16{
		data: elts,
	}

	err := a.InitIndexBitmap(index)
	if err == nil {
		// Check if an index present
		fmt.Println("a[0] present:", a.Has(0))
		fmt.Println("a[1] present:", a.Has(1))

		// Get element
		i, has := a.GetEltIndex(1)
		if has {
			fmt.Println("value of a[1]:", a.data[i])
		}
	}

	// Output:
	// a[0] present: false
	// a[1] present: true
	// value of a[1]: 12

}
