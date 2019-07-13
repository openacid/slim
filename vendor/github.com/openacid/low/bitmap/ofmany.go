package bitmap

// OfMany creates a bitmap from a list of sub-bitmap bit positions.
// "sizes" specifies the total bits in every sub-bitmap.
//
// Since 0.1.9
func OfMany(subs [][]int32, sizes []int32) []uint64 {
	r := make([]int32, 0)
	base := int32(0)
	for i, e := range subs {
		for _, idx := range e {
			r = append(r, base+idx)
		}
		base += sizes[i]
	}
	return Of(r, base)
}
