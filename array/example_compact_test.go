package array

import (
	"fmt"
)

func Example() {

	// arr[0]   = 12
	// arr[5]   = 15
	// arr[9]   = 19

	indexes := []uint32{0, 5, 9}
	elts := []uint32{12, 15, 19}

	arr, err := New32(U32Conv{}, indexes, elts)
	if err != nil {
		fmt.Printf("Init compacted array error:%s\n", err)
		return
	}

	if arr.Has(indexes[1]) {
		val := arr.Get(indexes[1])
		fmt.Printf("get indexed 1 value:%v\n", val)
	}

	// Output: get indexed 1 value:15
}
