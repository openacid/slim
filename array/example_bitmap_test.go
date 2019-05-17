package array_test

import (
	"fmt"

	"github.com/openacid/slim/array"
)

func ExampleNewBits() {

	// There is a "one" every 7 bits.
	// Bits cost is about 9 bits for a "one"
	n := 100
	bitset := make([]int32, n)
	for i := 0; i < n; i++ {
		bitset[i] = int32(i * 7)
	}

	a := array.NewBits(bitset)

	st := a.Stat()
	for _, k := range []string{"bits/one"} {
		fmt.Printf("%10s : %d\n", k, st[k])
	}

	fmt.Printf("has set at %d: %t\n", 2, a.Has(2))
	fmt.Printf("has set at %d: %t\n", 7, a.Has(7))
	fmt.Printf("ones count: %d\n", a.Rank(a.Len()))

	// Output:
	// bits/one : 9
	// has set at 2: false
	// has set at 7: true
	// ones count: 100
}
