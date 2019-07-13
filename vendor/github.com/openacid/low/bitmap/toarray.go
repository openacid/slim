package bitmap

// ToArray creates a new slice of int32 containing all of the integers stored in
// the bitmap in sorted order.
//
// Since 0.1.9
func ToArray(words []uint64) []int32 {

	r := make([]int32, 0)
	l := int32(len(words) * 64)

	for i := int32(0); i < l; i++ {
		if words[i>>6]&(1<<uint(i&63)) != 0 {
			r = append(r, i)
		}
	}

	return r
}
