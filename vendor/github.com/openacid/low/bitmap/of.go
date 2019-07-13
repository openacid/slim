package bitmap

// Of creates a bitmap of bits set to 1 at specified positions.
//
// An optional argument n can be provided to make the result bitmap
// have at least n bits.
//
// Since 0.1.9
func Of(bitPositions []int32, opts ...int32) []uint64 {

	n := int32(0)

	// The first opts is specified number of result bits.
	if len(opts) > 0 {
		n = opts[0]
	}

	if len(bitPositions) > 0 {
		max := bitPositions[len(bitPositions)-1] + 1
		if n < max {
			n = max
		}
	}

	nWords := (n + 63) >> 6
	words := make([]uint64, nWords)

	for _, i := range bitPositions {
		wordI := i >> 6
		i = i & 63
		words[wordI] |= 1 << uint(i)
	}
	return words

}
