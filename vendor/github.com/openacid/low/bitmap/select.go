package bitmap

import (
	"fmt"
	"math/bits"
)

// selectLookup8 is a lookup table for "select" on 8-bit bitmap:
//	select(aByte, ith)
// An element contains 8 4-bit integers.
var selectLookup8 [256]uint32
var tbl2 [256 * 8]uint8

func initSelectLookup() {

	// selectLookup8 = make([]uint32, 256)

	for i := 0; i < 256; i++ {
		v := uint32(0)
		w := uint8(i)
		for j := 0; j < 8; j++ {
			// x-th 1 in w
			// if x-th 1 is not found, it is 8
			x := bits.TrailingZeros8(w)
			v |= uint32(x) << uint(j*4)
			w &= w - 1

			tbl2[i*8+j] = uint8(x)
		}
		selectLookup8[i] = v
	}
}

// IndexSelect32 creates a index for operation "select" on a bitmap.
// select(i) returns the position of the i-th "1".
// E.g.:
//     bitmap = 100100..
//     select(bitmap, 0) = 1
//     select(bitmap, 1) = 3
//
// It returns an index of []int32.
// An element in it is the value of select(i*32)
//
// Since 0.1.9
func IndexSelect32(words []uint64) []int32 {
	l := len(words) << 6
	sidx := make([]int32, 0, len(words))

	ith := -1
	for i := 0; i < l; i++ {
		if words[i>>6]&(1<<uint(i&63)) != 0 {
			ith++
			if ith&31 == 0 {
				sidx = append(sidx, int32(i))
			}
		}
	}

	// clone to reduce cap to len
	sidx = append(sidx[:0:0], sidx...)
	return sidx
}

// select32single returns the index of the i-th "1".
// E.g.
//
//    bitmap: 10110..
//    index:  01234...
//
//    select(0):  0
//    select(1):  2
//    select(2):  3
//
// Since 0.1.9
func select32single(words []uint64, selectIndex []int32, i int32) int32 {

	if i < 0 {
		return -1
	}

	if i>>5 >= int32(len(selectIndex)) {
		return int32(len(words) * 64)
	}

	base := selectIndex[i>>5]
	findIth := int(i & 31)

	if findIth == 0 {
		return base
	}

	l := int32(len(words))
	wordI := base >> 6
	w := words[wordI]
	// remove "1" upto i64th excluding the "1" at i64th
	w = w & ^Mask[base&63]

	base = wordI << 6

	// continue search for i-th 1
	for {

		ones := bits.OnesCount64(w)
		if ones > findIth {

			ones = bits.OnesCount64(w & 0xffffffff)

			if ones <= findIth {
				findIth -= ones
				base += 32
				w >>= 32
			}

			ones = bits.OnesCount64(w & 0xffff)

			if ones <= findIth {
				findIth -= ones
				base += 16
				w >>= 16
			}

			ones = bits.OnesCount64(w & 0xff)

			if ones <= findIth {
				findIth -= ones
				base += 8
				w >>= 8
			}

			tbl := selectLookup8[w&0xff]
			return base + int32((tbl>>uint(findIth<<2))&0xf)

		} else {
			findIth -= ones
		}

		base += 64
		wordI++
		if wordI >= l {
			return l * 64
		}
		w = words[wordI]
	}
}

// Select32 returns the indexes of the i-th "1" and the (i+1)-th "1".
//
// Since 0.1.9
func Select32(words []uint64, selectIndex []int32, i int32) (int32, int32) {

	a := int32(0)

	if i < 0 || i>>5 >= int32(len(selectIndex)) {
		panic(fmt.Sprintf("i outof range: %d", i))
	}

	base := selectIndex[i>>5]
	findIth := int(i & 31)

	l := int32(len(words))
	wordI := base >> 6
	w := words[wordI]
	// remove "1" upto i64th excluding the "1" at i64th
	w = w & ^Mask[base&63]

	// continue search for i-th 1
	for {

		ones := bits.OnesCount64(w)
		if ones <= findIth {
			findIth -= ones
			wordI++
			w = words[wordI]
			continue
		}

		base := int32(0)
		ww := w

		ones = bits.OnesCount32(uint32(ww))

		if ones <= findIth {
			findIth -= ones
			base |= 32
			ww >>= 32
		}

		ones = bits.OnesCount16(uint16(ww))

		if ones <= findIth {
			findIth -= ones
			base |= 16
			ww >>= 16
		}

		ones = bits.OnesCount8(uint8(ww))

		if ones <= findIth {
			a = int32(tbl2[(ww>>5)&(0x7f8)|uint64(findIth-ones)]) + base + 8
		} else {
			a = int32(tbl2[(ww&0xff)<<3|uint64(findIth)]) + base
		}

		a += (wordI << 6)
		break
	}

	w = w & ^MaskUpto[a&63]
	if w != 0 {
		return a, wordI<<6 + int32(bits.TrailingZeros64(w))
	}

	for wordI := a>>6 + 1; wordI < l; wordI++ {
		w = words[wordI]
		if w != 0 {
			return a, wordI<<6 + int32(bits.TrailingZeros64(w))
		}
	}
	return a, l << 6
}

// indexSelectU64 create a 8 uint8 array in a uint64.
// Element a[i] is the count of "1" in the least (i+1)*8 bits.
//
// Since 0.1.9
func indexSelectU64(w uint64) uint64 {

	all1 := ^uint64(0)
	mask01 := (all1 / 3)
	mask0011 := (all1 / 5)
	mask00001111 := (all1 / 0x11)

	// a = (w & 0x5555...) + ((w >> 1) & 0x5555...)
	a := w - ((w >> 1) & mask01)

	// b = (a & 0x3333...) + ((a >> 2) & 0x3333...)
	b := (a & mask0011) + ((a >> 2) & mask0011)

	// c = (b & 0x0f0f...) + ((b >> 4) & 0x0f0f...)
	c := (b + (b >> 4)) & mask00001111

	c *= 0x0101010101010101
	return c | 0x8080808080808080
}

func selectU64Indexed(w uint64, index uint64, findIth uint64) (int32, int) {

	// if findIth > 63 {
	//     panic("findIth must be <= 63")
	// }

	v := (findIth + 1) * 0x0101010101010101
	// diff := index - v
	biggerBits := (index - v) & 0x8080808080808080

	ithU8 := bits.TrailingZeros64(biggerBits) & (^7)

	// if ithU8 == 64 {
	//     return 64, int(index >> 56 & 0x7f)
	// }

	findIth = findIth - (index>>uint(ithU8-8))&0x7f
	vv := tbl2[(w>>uint(ithU8)&0xff)<<3+findIth]

	// f2 := diff >> uint(ithU8)
	return int32(vv) + int32(ithU8), 0

}
