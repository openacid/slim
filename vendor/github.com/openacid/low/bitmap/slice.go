package bitmap

// Slice returns a new bitmap which is a "slice" of the input bitmap.
//
// Since 0.1.9
func Slice(words []uint64, from, to int32) []uint64 {

	l := ((to - from) + 63) & (^63)
	r := make([]uint64, l)

	for i := from; i < to; i++ {
		if words[i>>6]&(1<<uint(i&63)) != 0 {
			j := i - from
			r[j>>6] |= 1 << uint(j&63)
		}
	}

	return r
}
