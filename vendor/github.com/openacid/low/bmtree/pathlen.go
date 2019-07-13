package bmtree

import (
	"math/bits"
)

// PathLen returns the number of bits for this varbit.
//
// Since 0.1.9
func PathLen(p uint64) int32 {
	return int32(bits.OnesCount32(uint32(p)))
}
