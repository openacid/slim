package sigbits

// countPrefixes counts the number of n-bit prefixes:
//
// It returns the minimal bit index where there is a different bit
// and a []int32 of length maxitem + 1.
// The ith element represents the number of distinct i-bits words.
// E.g.
// There is only 1 0-bit word.
// There are at most 2 1-bit words.
//
// The first argument is made by FirstDiffBits(keys).
//
// Since 0.1.9
func countPrefixes(firstdiffs []int32, maxitem int32) (int32, []int32) {

	min := int32(0x7fffffff)
	for _, d := range firstdiffs {
		if min > d {
			min = d
		}
	}

	// counts[i] is number of i-th bits that is the first diff bit.
	counts := make([]int32, maxitem-1)
	for _, d := range firstdiffs {
		d -= min
		if d < maxitem-1 {
			counts[d]++
		}
	}

	// rst[i] means how many distinct i-bit words there are there are
	rst := make([]int32, maxitem)
	rst[0] = 1
	for i := int32(0); i < maxitem-1; i++ {
		rst[i+1] = rst[i] + counts[i]
	}

	return min, rst
}
