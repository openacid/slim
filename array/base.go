package array

import (
	"encoding/binary"
	"reflect"

	"math/bits"

	"github.com/openacid/errors"
	"github.com/openacid/low/bitmap"
	"github.com/openacid/slim/encode"
)

// endian is the default endian for array
var endian = binary.LittleEndian

// Base is the base of: Array and U16 etc.
//
// Performance note.
//
//	   Has():          9~10 ns / call; 1 memory accesses
//	   GetEltIndex(): 10~20 ns / call; 2 memory accesses
//
// Since 0.2.0
type Base struct {
	Array32
	EltEncoder encode.Encoder
}

const (
	// bmWidth defines how many bits for a bitmap word
	// Never change this
	bmWidth = int32(64)
	bmShift = uint(6) // log₂64
	bmMask  = int32(63)
)

// bmBit calculates bitamp word index and the bit index in the word.
func bmBit(idx int32) (int32, int32) {
	c := idx >> bmShift
	r := idx & bmMask
	return c, r
}

// InitIndex initializes index bitmap for an array.
// Index must be an ascending int32 slice, otherwise, it return
// the ErrIndexNotAscending error
//
// Since 0.2.0
func (a *Base) InitIndex(index []int32) error {

	for i := 0; i < len(index)-1; i++ {
		if index[i] >= index[i+1] {
			return ErrIndexNotAscending
		}
	}

	if bmWidth != 64 {
		panic("newBitmapWords only accept uint64 as bitmap word")
	}

	_, a.Bitmaps = newBitsWords(index)
	a.Offsets = bitmap.IndexRank64(a.Bitmaps)
	a.Cnt = int32(len(index))

	// Be compatible to previous issue:
	// Since v0.2.0, Offsets is not exactly the same as bitmap ranks.
	// It is 0 for empty bitmap word.
	// But bitmap ranks set rank[i*64] to rank[(i-1)*64] for empty word.
	for i, word := range a.Bitmaps {
		if word == 0 {
			a.Offsets[i] = 0
		}
	}

	return nil
}

// ExtendIndex allocaed additional 0-bits after Bitmap and Offset.
//
// Since 0.5.9
func (a *Base) ExtendIndex(n int32) {
	nword := (n + 63) >> 6

	if nword <= int32(len(a.Bitmaps)) {
		return
	}

	bitmaps := make([]uint64, nword)
	copy(bitmaps, a.Bitmaps)

	a.Bitmaps = bitmaps
	a.Offsets = bitmap.IndexRank64(a.Bitmaps)
	for i, word := range a.Bitmaps {
		if word == 0 {
			a.Offsets[i] = 0
		}
	}
}

// GetEltIndex returns the position in a.Elts of element[idx] and a bool
// indicating if found or not.
// If "idx" absents it returns "0, false".
//
// Since 0.2.0
func (a *Base) GetEltIndex(idx int32) (int32, bool) {
	iBm, iBit := bmBit(idx)

	var bmWord = a.Bitmaps[iBm]

	if ((bmWord >> uint(iBit)) & 1) == 0 {
		return 0, false
	}

	base := a.Offsets[iBm]
	cnt1 := bits.OnesCount64(bmWord & ((uint64(1) << uint(iBit)) - 1))
	return base + int32(cnt1), true
}

// Has returns true if idx is in array, else return false.
//
// Since 0.2.0
func (a *Base) Has(idx int32) bool {
	iBm := idx / bmWidth
	return ((a.Bitmaps[iBm] >> uint32(idx&bmMask)) & 1) != 0
}

// Init initializes an array from the "indexes" and "elts".
// The indexes must be an ascending int32 slice,
// otherwise, return the ErrIndexNotAscending error.
// The "elts" is a slice.
//
// Since 0.2.0
func (a *Base) Init(indexes []int32, elts interface{}) error {

	rElts := reflect.ValueOf(elts)
	if rElts.Kind() != reflect.Slice {
		panic("elts is not a slice")
	}

	n := rElts.Len()
	if len(indexes) != n {
		return ErrIndexLen
	}

	err := a.InitIndex(indexes)
	if err != nil {
		return err
	}

	if len(indexes) == 0 {
		return nil
	}

	var encoder encode.Encoder

	if a.EltEncoder == nil {
		var err error
		encoder, err = encode.NewTypeEncoderEndian(rElts.Index(0).Interface(), endian)
		if err != nil {
			// TODO wrap
			return err
		}
	} else {
		encoder = a.EltEncoder
	}

	_, err = a.InitElts(elts, encoder)
	if err != nil {
		return errors.Wrapf(err, "failure Init Array")
	}

	return nil
}

// InitElts initialized a.Elts, by encoding elements in to bytes.
//
// Since 0.2.0
func (a *Base) InitElts(elts interface{}, encoder encode.Encoder) (int, error) {

	rElts := reflect.ValueOf(elts)
	n := rElts.Len()
	eltsize := encoder.GetEncodedSize(nil)
	sz := eltsize * n

	b := make([]byte, 0, sz)
	for i := 0; i < n; i++ {
		ee := rElts.Index(i).Interface()
		bs := encoder.Encode(ee)
		b = append(b, bs...)
	}
	a.Elts = b

	return n, nil
}

// MemSize returns the memory this array occupies.
// It includes .Cnt, .Bitmaps, .Offsets and .Elts, not includes the meomory of
// the structure itself.
//
// Since 0.3.1
func (a *Base) MemSize() int {
	return 4 + 8*len(a.Bitmaps) + 4*len(a.Offsets) + len(a.Elts)
}

// Get retrieves the value at "idx" and return it.
// If this array has a value at "idx" it returns the value and "true",
// otherwise it returns "nil" and "false".
//
// Since 0.2.0
func (a *Base) Get(idx int32) (interface{}, bool) {

	bs, ok := a.GetBytes(idx, a.EltEncoder.GetEncodedSize(nil))
	if ok {
		_, v := a.EltEncoder.Decode(bs)
		return v, true
	}

	return nil, false
}

// GetBytes retrieves the raw data of value in []byte at "idx" and return it.
//
// Performance note
//
// Involves 2 memory access:
//	 a.Bitmaps
//	 a.Elts
//
// Involves 0 alloc
//
// Since 0.2.0
func (a *Base) GetBytes(idx int32, eltsize int) ([]byte, bool) {
	dataIndex, ok := a.GetEltIndex(idx)
	if !ok {
		return nil, false
	}

	stIdx := int32(eltsize) * dataIndex
	return a.Elts[stIdx : stIdx+int32(eltsize)], true
}

// Indexes returns indexes of all present elements.
//
// Since 0.5.4
func (a *Base) Indexes() []int32 {

	rst := make([]int32, a.Cnt)
	j := int32(0)

	for i := int32(0); i < a.Cnt; {
		if a.Has(j) {
			rst[i] = j
			i++
		}
		j++
	}
	return rst
}
