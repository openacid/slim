package array

import (
	"errors"

	"github.com/openacid/slim/bit"
	"github.com/openacid/slim/prototype"
)

// Array32Index implements sparsely distributed index with bitmap.
type Array32Index struct {
	prototype.CompactedArray
}

// ErrIndexNotAscending means indexes to initialize a Array must be in
// ascending order.
var ErrIndexNotAscending = errors.New("index must be an ascending slice")

// bmWidth defines how many bits for a bitmap word
var bmWidth = uint32(64)

// bmBit calculates bitamp word index and the bit index in the word.
func bmBit(idx uint32) (uint32, uint32) {
	c := idx >> uint32(6) // == idx / bmWidth
	r := idx & uint32(63) // == idx % bmWidth
	return c, r
}

// InitIndexBitmap initializes index bitmap for a compacted array.
// Index must be a ascending array of type unit32, otherwise, return
// the ErrIndexNotAscending error
func (a *Array32Index) InitIndexBitmap(index []uint32) error {

	capacity := uint32(0)
	if len(index) > 0 {
		capacity = index[len(index)-1] + 1
	}

	bmCnt := (capacity + bmWidth - 1) / bmWidth

	a.Bitmaps = make([]uint64, bmCnt)
	a.Offsets = make([]uint32, bmCnt)

	nxt := uint32(0)
	for i := 0; i < len(index); i++ {
		if index[i] < nxt {
			return ErrIndexNotAscending
		}
		a.appendIndex(index[i])
		nxt = index[i] + 1
	}
	return nil
}

// GetEltIndex returns the data position in a.Elts indexed by `idx` and a bool
// indicating existence.
// If `idx` does not present it returns `0, false`.
func (a *Array32Index) GetEltIndex(idx uint32) (uint32, bool) {
	iBm, iBit := bmBit(idx)

	if iBm >= uint32(len(a.Bitmaps)) {
		return 0, false
	}

	var bmWord = a.Bitmaps[iBm]

	if ((bmWord >> iBit) & 1) == 0 {
		return 0, false
	}

	base := a.Offsets[iBm]
	cnt1 := bit.PopCnt64Before(bmWord, iBit)
	return base + cnt1, true
}

// Has returns true if idx is in array, else return false
func (a *Array32Index) Has(idx uint32) bool {
	iBm, iBit := bmBit(idx)

	if iBm >= uint32(len(a.Bitmaps)) {
		return false
	}

	var bmWord = a.Bitmaps[iBm]

	return (bmWord>>iBit)&1 > 0
}

// appendIndex add a index into index bitmap.
// The `index` must be greater than any existent indexes.
func (a *Array32Index) appendIndex(index uint32) {

	iBm, iBit := bmBit(index)

	var bmWord = &a.Bitmaps[iBm]
	if *bmWord == 0 {
		a.Offsets[iBm] = a.Cnt
	}

	*bmWord |= uint64(1) << iBit

	a.Cnt++
}
