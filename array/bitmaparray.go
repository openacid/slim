package array

import "math/bits"

type Bitmap16 struct {
	Base
}

const (
	ArrayFlagHasEltWidth = uint32(0x00000001)
	ArrayFlagIsBitmap    = uint32(0x00000002)
)

func NewBitmap16(index []int32, elts []uint64, eltWidth int32) (*Bitmap16, error) {

	b := &Bitmap16{}
	b.Flags = ArrayFlagHasEltWidth | ArrayFlagIsBitmap
	b.EltWidth = eltWidth

	err := b.InitIndex(index)
	if err != nil {
		return nil, err
	}

	bm := NewBitsJoin(elts, eltWidth).(*Bits)
	b.BMElts = bm

	return b, nil
}

func (b *Bitmap16) GetWithRank(idx int32) (uint64, int32, bool) {

	iBm := idx >> bmShift
	iBit := idx & bmMask

	var bm = b.Bitmaps[iBm]

	if ((bm >> uint(iBit)) & 1) == 0 {
		return 0, 0, false
	}

	cnt1 := bits.OnesCount64(bm & ((uint64(1) << uint(iBit)) - 1))

	eltBitIdx := (b.Offsets[iBm] + int32(cnt1)) << 4

	iWord := eltBitIdx >> 6
	j := eltBitIdx & 63

	w := b.BMElts.Words[iWord]

	v := (w >> uint(j)) & 0xffff

	// Bitmap16 does not use dense mode rank index
	// Bitmap16 does not use EltWidth

	rank := b.BMElts.RankIndex[(eltBitIdx+64)>>7]

	if iWord&1 == 0 {
		word := w << (64 - uint(j))
		rank += int32(bits.OnesCount64(word))
	} else {
		word := w >> uint(j)
		rank -= int32(bits.OnesCount64(word))
	}

	return v, rank, true
}
