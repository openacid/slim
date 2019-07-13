package bmtree

import (
	"math/bits"

	"github.com/openacid/low/bitmap"
)

// AllPaths returns all possible searching path a bitmap can have.
//
// Since 0.1.9
func AllPaths(bitmapSize int32, from, to uint64) []uint64 {

	height := Height(bitmapSize)

	paths := make([]uint64, 0)

	fullPathCnt := bitmap.Bit[height]
	fullPathMask := bitmap.Mask[height]

	var t uint64
	t = to>>32 + 1

	if t > fullPathCnt {
		t = fullPathCnt
	}

	for i := from >> 32; i < t; i++ {
		tz := int32(bits.TrailingZeros64(i))
		if tz > height {
			tz = height
		}

		for ; tz >= 0; tz-- {
			if bitmapSize&int32(bitmap.Bit[height-tz]) == 0 {
				continue
			}
			m := bitmap.Mask[tz]
			p := (i << 32) | (fullPathMask ^ m)
			if p < from {
				// i is full searching bits without length, that it may generate
				// smaller path.
				continue
			}
			if p >= to {
				return paths
			}
			paths = append(paths, p)
		}
	}

	return paths
}
