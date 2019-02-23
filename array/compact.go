// Package array implements functions for the manipulation of compacted array.
package array

import (
	"errors"
	"reflect"

	"github.com/openacid/slim/bit"
	"github.com/openacid/slim/prototype"
)

// Array32 is a space efficient array implementation.
//
// Unlike a normal array, it does not allocate space for a element that there is
// not data in it.
type Array32 struct {
	prototype.CompactedArray
	Converter
}

var bmWidth = uint32(64) // how many bits of an uint64 == 2 ^ 6

func bmBit(idx uint32) (uint32, uint32) {
	c := idx >> uint32(6) // == idx / bmWidth
	r := idx & uint32(63) // == idx % bmWidth
	return c, r
}

// New32 creates a CompactedArray and initializes it with a slice of index and a
// slice of data.
//
// The index parameter must be a ascending array of type unit32,
// otherwise, return the ErrIndexNotAscending error
func New32(conv Converter, index []uint32, _elts interface{}) (ca *Array32, err error) {

	ca = &Array32{
		Converter: conv,
	}

	err = ca.Init(index, _elts)
	if err != nil {
		return nil, err
	}

	return ca, nil
}

func (a *Array32) appendElt(index uint32, elt []byte) {
	iBm, iBit := bmBit(index)

	var bmWord = &a.Bitmaps[iBm]
	if *bmWord == 0 {
		a.Offsets[iBm] = a.Cnt
	}

	*bmWord |= uint64(1) << iBit
	a.Elts = append(a.Elts, elt...)

	a.Cnt++
}

// ErrIndexLen is returned if number of indexes does not equal the number of
// datas, when initializing a CompactedArray.
var ErrIndexLen = errors.New("the length of index and elts must be equal")

// ErrIndexNotAscending means both indexes and datas for initialize a
// CompactedArray must be in ascending order.
var ErrIndexNotAscending = errors.New("index must be an ascending slice")

// Init initializes a compacted array from the slice type elts
// the index parameter must be a ascending array of type unit32,
// otherwise, return the ErrIndexNotAscending error
func (a *Array32) Init(index []uint32, _elts interface{}) error {

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

	a.Bitmaps = make([]uint64, bmCnt)
	a.Offsets = make([]uint32, bmCnt)
	a.Elts = make([]byte, 0, nElts*a.GetMarshaledSize(nil))

	var prevIndex uint32
	for i := 0; i < len(index); i++ {
		if i > 0 && index[i] <= prevIndex {
			return ErrIndexNotAscending
		}

		ee := rElts.Index(i).Interface()
		a.appendElt(index[i], a.Marshal(ee))

		prevIndex = index[i]
	}

	return nil
}

// Get returns the value indexed by idx if it is in array, else return nil
func (a *Array32) Get(idx uint32) interface{} {
	v, _ := a.Get2(idx)
	return v
}

// Get2 returns the value indexed by `idx` and a bool indicating existence.
// If `idx` does not present it returns `nil, false`.
func (a *Array32) Get2(idx uint32) (interface{}, bool) {
	iBm, iBit := bmBit(idx)

	if iBm >= uint32(len(a.Bitmaps)) {
		return nil, false
	}

	var bmWord = a.Bitmaps[iBm]

	if ((bmWord >> iBit) & 1) == 0 {
		return nil, false
	}

	base := a.Offsets[iBm]
	cnt1 := bit.PopCnt64Before(bmWord, iBit)

	stIdx := uint32(a.GetMarshaledSize(nil)) * (base + cnt1)

	_, val := a.Unmarshal(a.Elts[stIdx:])
	return val, true
}

// Has returns true if idx is in array, else return false
func (a *Array32) Has(idx uint32) bool {
	iBm, iBit := bmBit(idx)

	if iBm >= uint32(len(a.Bitmaps)) {
		return false
	}

	var bmWord = a.Bitmaps[iBm]

	return (bmWord>>iBit)&1 > 0
}
