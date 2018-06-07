package bit

import (
	"unsafe"
)

func Cnt1Before(n uint64, iBit uint32) uint32 {
	var cnt uint32 = 0
	MaxBit := uint32(unsafe.Sizeof(n)) * 8
	if iBit == 0 {
		n = 0
	} else {
		if iBit < MaxBit {
			n <<= MaxBit - iBit
		}
	}

	for {
		if n == 0 {
			break
		}

		n &= n - 1
		cnt++
	}
	return cnt
}
