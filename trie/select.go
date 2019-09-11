package trie

import (
	"math/bits"

	"github.com/openacid/low/bitmap"
)

// select8Lookup is a lookup table for "select" on 8-bit bitmap:
//	   select(aByte, ith)
var select8Lookup [256 * 8]uint8

func init() {
	initSelectLookup()
}

func initSelectLookup() {

	for i := 0; i < 256; i++ {
		v := uint32(0)
		w := uint8(i)
		for j := 0; j < 8; j++ {
			// x-th 1 in w
			// if x-th 1 is not found, it is 8
			x := bits.TrailingZeros8(w)
			v |= uint32(x) << uint(j*4)
			w &= w - 1

			select8Lookup[i*8+j] = uint8(x)
		}
	}
}

// indexSelect32 creates an index to speed up "select".
// The index values is word position but not bit position.
// Thus to use it a rank index must also be provided.
//
func indexSelect32(words []uint64) []int32 {
	l := len(words) << 6
	sidx := make([]int32, 0, len(words))

	ith := -1
	for i := 0; i < l; i++ {
		if words[i>>6]&(1<<uint(i&63)) != 0 {
			ith++
			if ith&31 == 0 {
				sidx = append(sidx, int32(i>>6))
			}
		}
	}

	// clone to reduce cap to len
	sidx = append(sidx[:0:0], sidx...)
	return sidx
}

// it requires a Rank64 index and a Select32 index
func (bm *Bitmap) select32(i int32) (int32, int32) {

	a := int32(0)
	l := int32(len(bm.Words))

	wordI := bm.SelectIndex[i>>5]
	for ; bm.RankIndex[wordI+1] <= i; wordI++ {
	}

	w := bm.Words[wordI]
	ww := w
	base := wordI << 6
	findIth := int(i - bm.RankIndex[wordI])

	offset := int32(0)

	ones := bits.OnesCount32(uint32(ww))
	if ones <= findIth {
		findIth -= ones
		offset |= 32
		ww >>= 32
	}

	ones = bits.OnesCount16(uint16(ww))
	if ones <= findIth {
		findIth -= ones
		offset |= 16
		ww >>= 16
	}

	ones = bits.OnesCount8(uint8(ww))
	if ones <= findIth {
		a = int32(select8Lookup[(ww>>5)&(0x7f8)|uint64(findIth-ones)]) + offset + 8
	} else {
		a = int32(select8Lookup[(ww&0xff)<<3|uint64(findIth)]) + offset
	}

	a += base

	// "& 63" elides boundary check
	w &= bitmap.RMaskUpto[a&63]

	if w != 0 {
		return a, base + int32(bits.TrailingZeros64(w))
	}

	wordI++
	for ; wordI < l; wordI++ {
		w = bm.Words[wordI]
		if w != 0 {
			return a, wordI<<6 + int32(bits.TrailingZeros64(w))
		}
	}
	return a, l << 6
}
