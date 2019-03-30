package array

import (
	"encoding/binary"
	"reflect"

	"github.com/openacid/errors"
	"github.com/openacid/slim/bits"
	"github.com/openacid/slim/encode"
	"github.com/openacid/slim/prototype"
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
	prototype.Array32
	EltEncoder encode.Encoder
}

const (
	// bmWidth defines how many bits for a bitmap word
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

	capacity := int32(0)
	if len(index) > 0 {
		capacity = index[len(index)-1] + 1
	}

	bmCnt := (capacity + bmWidth - 1) / bmWidth

	a.Bitmaps = make([]uint64, bmCnt)
	a.Offsets = make([]int32, bmCnt)

	nxt := int32(0)
	for i := 0; i < len(index); i++ {
		if index[i] < nxt {
			return ErrIndexNotAscending
		}
		a.appendIndex(index[i])
		nxt = index[i] + 1
	}
	return nil
}

// GetEltIndex returns the position in a.Elts of element[idx] and a bool
// indicating if found or not.
// If "idx" absents it returns "0, false".
//
// Since 0.2.0
func (a *Base) GetEltIndex(idx int32) (int32, bool) {
	iBm, iBit := bmBit(idx)

	if iBm >= int32(len(a.Bitmaps)) {
		return 0, false
	}

	var bmWord = a.Bitmaps[iBm]

	if ((bmWord >> uint(iBit)) & 1) == 0 {
		return 0, false
	}

	base := a.Offsets[iBm]
	cnt1 := bits.OnesCount64Before(bmWord, uint(iBit))
	return base + int32(cnt1), true
}

// Has returns true if idx is in array, else return false.
//
// Since 0.2.0
func (a *Base) Has(idx int32) bool {
	iBm := idx / bmWidth
	return iBm < int32(len(a.Bitmaps)) && ((a.Bitmaps[iBm]>>uint32(idx&bmMask))&1) != 0
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

// Get retrieves the value at "idx" and return it.
// If this array has a value at "idx" it returns the value and "true",
// otherwise it returns "nil" and "false".
//
// Since 0.2.0
func (a *Base) Get(idx int32) (interface{}, bool) {

	if a.Cnt == 0 {
		return nil, false
	}

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

// appendIndex add an index into index bitmap.
// The `index` must be greater than any existent indexes.
func (a *Base) appendIndex(index int32) {

	iBm, iBit := bmBit(index)

	var bmWord = &a.Bitmaps[iBm]
	if *bmWord == 0 {
		a.Offsets[iBm] = a.Cnt
	}

	*bmWord |= uint64(1) << uint(iBit)

	a.Cnt++
}
