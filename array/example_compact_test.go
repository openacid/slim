package array

import (
	"fmt"
)

func Example() {
	indexes := []uint32{0, 5, 9, 203, 400}
	eltsData := []uint32{12, 15, 19, 120, 300}

	ca := CompactedArray{Converter: U32Conv{}}
	err := ca.Init(indexes, eltsData)
	if err != nil {
		fmt.Printf("Init compacted array error:%s\n", err)
	} else {
		if ca.Has(indexes[1]) {
			val := ca.Get(indexes[1])
			fmt.Printf("get indexed 1 value:%v\n", val)
		}
	}

	// Output: get indexed 1 value:15
}
