package array

import (
	"errors"

	"github.com/openacid/slim/bits"
	"github.com/openacid/slim/prototype"
)

// Array32Index implements sparsely distributed index with bitmap.
//
// Performance note
//
// Has():          9~10 ns / call; 1 memory accesses
// GetEltIndex(): 10~20 ns / call; 2 memory accesses
//
// Most time is spent on Bitmaps and Offsets access:
// L1 or L2 cache assess costs 0.5 ns and 7 ns.
type Array32Index struct {
	prototype.Array32Storage
}

func (a *Array32Index) GetStorage() *prototype.Array32Storage {
	return &a.Array32Storage
}

// ErrIndexNotAscending means indexes to initialize a Array must be in
// ascending order.
var ErrIndexNotAscending = errors.New("index must be an ascending ordered slice")

const (
	// bmWidth defines how many bits for a bitmap word
	bmWidth = uint32(64)
	bmMask  = uint32(63)
)

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
	cnt1 := bits.OnesCount64Before(bmWord, uint(iBit))
	return base + uint32(cnt1), true
}

// Has returns true if idx is in array, else return false.
func (a *Array32Index) Has(idx uint32) bool {
	iBm := idx / bmWidth
	return iBm < uint32(len(a.Bitmaps)) && ((a.Bitmaps[iBm]>>(idx&bmMask))&1) != 0
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
