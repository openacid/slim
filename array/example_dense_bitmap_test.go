package array_test

import (
	"fmt"

	"github.com/openacid/slim/array"
)

func ExampleNewDenseBitmap() {

	// For large bit set, use DenseBitmap:
	// There is a "one" every 7 bits.
	// DenseBitmap cost is almost still 7 bits

	n := 1024 * 1024
	bitset := make([]int32, n)
	for i := 0; i < n; i++ {
		bitset[i] = int32(i * 7)
	}

	a := array.NewDenseBitmap(bitset)
	st := a.Stat()

	for _, k := range []string{"bits/one"} {
		fmt.Printf("%10s : %d\n", k, st[k])
	}

	fmt.Printf("has set at %d: %t\n", 2, a.Has(2))
	fmt.Printf("has set at %d: %t\n", 7, a.Has(7))
	fmt.Printf("ones count: %d\n", a.Rank(a.Len()))

	// Output:
	// bits/one : 7
	// has set at 2: false
	// has set at 7: true
	// ones count: 1048576
}
