package array

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

	bm := NewBitsJoin(elts, eltWidth, false).(*Bits)
	b.BMElts = bm

	return b, nil
}

func (b *Bitmap16) GetWithRank(idx int32) (uint64, int32, bool) {

	i, found := b.GetEltIndex(idx)
	if !found {
		return 0, 0, false
	}

	i *= b.EltWidth
	iWord := i >> 6
	j := i & 63

	w := b.BMElts.Words[iWord]

	v := (w >> uint(j)) & ((1 << uint(b.EltWidth)) - 1)
	rank := b.BMElts.Rank(i)

	return v, rank, true
}
