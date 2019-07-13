package bmtree

import (
	"math/bits"
)

// shiftMulti maps a tree node to bitmap position.
// The bitmap is with node at some level not stored, and of size 110011...
// A "0" in bitmap size indicates an absent level.
//
// It is defined as:
//	   sum = 0
//	   for i = 0; i<h; i++ {
//		 sum += (t * p[i]) << i >> h
//	   }
//
//			 t[h] t[h-1].. t[1] t[0]
//		   X      p[h-1]...p[1] p[0]
//	   -------------------------------
//			 t[h] t[h-1].. t[1] t[0]    x   p[0]
//		t[h] t[h-1].. t[1] t[0]         x   p[1]
//		...
//	   -------------------------------
//				  |
//
// Since ???
func shiftMulti(a, b, shift uint64) uint64 {

	rst := uint64(0)

	n := bits.TrailingZeros64(b)
	b >>= uint(n)
	shift -= uint64(n)

	for b != 0 {
		rst += (a >> shift)
		n := bits.TrailingZeros64(b - 1)
		b >>= uint(n)
		shift -= uint64(n)
	}

	return rst

}
