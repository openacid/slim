package bitmap

// Get returns a uint64 with the (i%64)-th bit set.
//
// Since 0.1.9
func Get(bm []uint64, i int32) uint64 {
	return (bm[i>>6] & Bit[i&63])
}

// Get1 returns a uint64 with the least significant bit set if (i%64)-th bit is
// set.
//
// Since 0.1.9
func Get1(bm []uint64, i int32) uint64 {
	return (bm[i>>6] >> uint(i&63)) & 1
}

// Getw returns the i-th w-bit word in the least significant w bits of a
// uint64.
// "w" must be one of 1, 2, 4, 8, 16, 32 and 64.
//
// Since 0.1.9
func Getw(bm []uint64, i int32, w int32) uint64 {
	i *= w
	return (bm[i>>6] >> uint(i&63)) & Mask[w]
}

// SafeGet is same as Get() except it return 0 instead of a panic when index out
// of boundary.
//
// Since 0.1.9
func SafeGet(bm []uint64, i int32) uint64 {
	wordI := i >> 6
	bitI := i & 63
	if wordI < 0 || wordI >= int32(len(bm)) {
		return 0
	}
	return (bm[wordI] & Bit[bitI])
}

// SafeGet1 is same as Get1() except it return 0 instead of a panic when index
// out of boundary.
//
// Since 0.1.9
func SafeGet1(bm []uint64, i int32) uint64 {
	wordI := i >> 6
	bitI := i & 63
	if wordI < 0 || wordI >= int32(len(bm)) {
		return 0
	}
	return (bm[wordI] >> uint(bitI)) & 1
}
