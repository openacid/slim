package bmtree

import "math/bits"

// Height returns the tree height of this bitmap.
//
// Since 0.1.9
func Height(bitmapSize int32) int32 {
	return int32(31 - bits.LeadingZeros32(uint32(bitmapSize)))
}
