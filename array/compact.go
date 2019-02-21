package array

import (
	"errors"
	"reflect"

	"github.com/openacid/slim/bit"
	"github.com/openacid/slim/prototype"
)

type CompactedArray struct {
	prototype.CompactedArray
	Converter
}

var bmWidth = uint32(64) // how many bits of an uint64 == 2 ^ 6

func bmBit(idx uint32) (uint32, uint32) {
	c := idx >> uint32(6) // == idx / bmWidth
	r := idx & uint32(63) // == idx % bmWidth
	return c, r
}

func (sa *CompactedArray) appendElt(index uint32, elt []byte) {
	iBm, iBit := bmBit(index)

	var bmWord = &sa.Bitmaps[iBm]
	if *bmWord == 0 {
		sa.Offsets[iBm] = sa.Cnt
	}

	*bmWord |= uint64(1) << iBit
	sa.Elts = append(sa.Elts, elt...)

	sa.Cnt++
}

var ErrIndexLen = errors.New("the length of index and elts must be equal")
var ErrIndexNotAscending = errors.New("index must be an ascending slice")

func (sa *CompactedArray) Init(index []uint32, _elts interface{}) error {

	rElts := reflect.ValueOf(_elts)
	if rElts.Kind() != reflect.Slice {
		panic("input is not a slice")
	}

	nElts := rElts.Len()

	if len(index) != nElts {
		return ErrIndexLen
	}

	capacity := uint32(0)
	if len(index) > 0 {
		capacity = index[len(index)-1] + 1
	}

	bmCnt := (capacity + bmWidth - 1) / bmWidth

	sa.Bitmaps = make([]uint64, bmCnt)
	sa.Offsets = make([]uint32, bmCnt)
	sa.Elts = make([]byte, 0, nElts*sa.GetMarshaledSize(nil))

	var prevIndex uint32
	for i := 0; i < len(index); i++ {
		if i > 0 && index[i] <= prevIndex {
			return ErrIndexNotAscending
		}

		ee := rElts.Index(i).Interface()
		sa.appendElt(index[i], sa.Marshal(ee))

		prevIndex = index[i]
	}

	return nil
}

func (sa *CompactedArray) Get(idx uint32) interface{} {
	iBm, iBit := bmBit(idx)

	if iBm >= uint32(len(sa.Bitmaps)) {
		return nil
	}

	var bmWord = sa.Bitmaps[iBm]

	if ((bmWord >> iBit) & 1) == 0 {
		return nil
	}

	base := sa.Offsets[iBm]
	cnt1 := bit.PopCnt64Before(bmWord, iBit)

	stIdx := uint32(sa.GetMarshaledSize(nil)) * (base + cnt1)

	_, val := sa.Unmarshal(sa.Elts[stIdx:])
	return val
}

func (sa *CompactedArray) Has(idx uint32) bool {
	iBm, iBit := bmBit(idx)

	if iBm >= uint32(len(sa.Bitmaps)) {
		return false
	}

	var bmWord = sa.Bitmaps[iBm]

	return (bmWord>>iBit)&1 > 0
}
