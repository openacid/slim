package sigbits

import (
	"math/bits"
)

// sFirstDiffBit returns the first different bit position.
// If a or b is a prefix of the other, it returns the smaller length times 8.
//
// Since 0.1.9
func sFirstDiffBit(a, b string) int32 {

	la := len(a)
	lb := len(b)

	l1 := len(a) * 8
	l2 := len(b) * 8

	minl := l1
	if minl > l2 {
		minl = l2
	}

	for i := 0; i < la && i < lb; i += 8 {

		au := get64Bits(a[i:])
		bu := get64Bits(b[i:])

		first := bits.LeadingZeros64(au ^ bu)

		if first < 64 {
			first = i<<3 + first
			if first < minl {
				return int32(first)
			} else {
				return int32(minl)
			}
		}
	}

	return int32(minl)
}

// FirstDiffBits find the first different bit position of every two adjacent
// keys.
//
// Since 0.1.9
func FirstDiffBits(keys []string) []int32 {
	l := len(keys)

	ds := make([]int32, l-1)
	for i := 0; i < l-1; i++ {
		ds[i] = sFirstDiffBit(keys[i], keys[i+1])
	}

	return ds
}

// get64Bits converts a string of length upto 8 to a uint64,
// in big endian.
// Less than 8 byte string will be filled with trailing 0.
// More than 8 bytes will be ignored.
//
// Since 0.1.9
func get64Bits(s string) uint64 {

	if len(s) >= 8 {

		return ((uint64(s[0]) << 56) +
			(uint64(s[1]) << 48) +
			(uint64(s[2]) << 40) +
			(uint64(s[3]) << 32) +
			(uint64(s[4]) << 24) +
			(uint64(s[5]) << 16) +
			(uint64(s[6]) << 8) +
			(uint64(s[7])))

	} else {

		bs := make([]byte, 8)
		copy(bs, s)
		return ((uint64(bs[0]) << 56) +
			(uint64(bs[1]) << 48) +
			(uint64(bs[2]) << 40) +
			(uint64(bs[3]) << 32) +
			(uint64(bs[4]) << 24) +
			(uint64(bs[5]) << 16) +
			(uint64(bs[6]) << 8) +
			(uint64(bs[7])))
	}
}
