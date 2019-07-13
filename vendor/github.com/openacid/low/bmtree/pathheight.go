package bmtree

import "math/bits"

// PathHeight returns the tree height of a searching path.
//
// Since 0.1.9
func PathHeight(path uint64) int32 {
	return int32(32 - bits.LeadingZeros32(uint32(path)))
}
