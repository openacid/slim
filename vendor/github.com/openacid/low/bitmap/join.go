package bitmap

// Join creates a bitmap from a list of sub-bitmaps.
// "size" specifies the size of every sub-bitmap.
// Sub bitmaps must be of equal size.
//
// Since 0.1.9
func Join(subs []uint64, size int32) []uint64 {
	l := int(size) * len(subs)
	r := make([]uint64, (l+63)&(^63)>>6)
	for i, e := range subs {
		j := i * int(size)
		r[j>>6] |= (e & Mask[size]) << uint(j&63)

	}
	return r
}
