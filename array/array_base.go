package array

import (
	"bytes"
	"encoding/binary"
	"reflect"

	"github.com/openacid/errors"
	"github.com/openacid/slim/bits"
	"github.com/openacid/slim/marshal"
	"github.com/openacid/slim/prototype"
)

// endian is default endian for array
var endian = binary.LittleEndian

// ArrayBase is a space efficient array implementation.
//
// Unlike a normal array, it does not allocate space for a element that there is
// not data in it.
//
// Performance note
//
// Has():          9~10 ns / call; 1 memory accesses
// GetEltIndex(): 10~20 ns / call; 2 memory accesses
//
// Most time is spent on Bitmaps and Offsets access:
// L1 or L2 cache assess costs 0.5 ns and 7 ns.
// TODO test Cnt
type ArrayBase struct {
	prototype.Array32
	EltMarshaler marshal.Marshaler
}

const (
	// bmWidth defines how many bits for a bitmap word
	bmWidth = int32(64)
	bmMask  = int32(63)
)

// bmBit calculates bitamp word index and the bit index in the word.
func bmBit(idx int32) (int32, int32) {
	c := idx >> uint32(6) // == idx / bmWidth
	r := idx & int32(63)  // == idx % bmWidth
	return c, r
}

// InitIndex initializes index bitmap for a compacted array.
// Index must be a ascending array of type unit32, otherwise, return
// the ErrIndexNotAscending error
func (a *ArrayBase) InitIndex(index []int32) error {

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

// GetEltIndex returns the data position in a.Elts indexed by `idx` and a bool
// indicating existence.
// If `idx` does not present it returns `0, false`.
func (a *ArrayBase) GetEltIndex(idx int32) (int32, bool) {
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
func (a *ArrayBase) Has(idx int32) bool {
	iBm := idx / bmWidth
	return iBm < int32(len(a.Bitmaps)) && ((a.Bitmaps[iBm]>>uint32(idx&bmMask))&1) != 0
}

// Init initializes a compacted array from the slice type elts
// the indexes parameter must be a ascending array of type unit32,
// otherwise, return the ErrIndexNotAscending error
func (a *ArrayBase) Init(indexes []int32, elts interface{}) error {

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

	var marshaler marshal.Marshaler

	if a.EltMarshaler == nil {
		var err error
		marshaler, err = marshal.NewTypeMarshalerEndian(rElts.Index(0).Interface(), endian)
		if err != nil {
			// TODO wrap
			return err
		}
	} else {
		marshaler = a.EltMarshaler
	}

	_, err = a.InitElts(elts, marshaler)
	if err != nil {
		return errors.Wrapf(err, "failure Init Array")
	}

	return nil
}

func (a *ArrayBase) InitElts(elts interface{}, marshaler marshal.Marshaler) (int, error) {

	rElts := reflect.ValueOf(elts)
	n := rElts.Len()
	eltsize := marshaler.GetMarshaledSize(nil)
	sz := eltsize * n

	b := bytes.NewBuffer(make([]byte, 0, sz))
	for i := 0; i < n; i++ {
		ee := rElts.Index(i).Interface()
		bs := marshaler.Marshal(ee)
		b.Write(bs)
	}
	a.Elts = b.Bytes()

	return n, nil
}

// GetTo retrieves the value at idx.
// "v" is a pointer to a fixed size value as receiver.
// If found it returns true, otherwise false.
// When not found, "v" in intact.
//
// Performance note
//
// Involves 2 alloc:
//	 bytes.NewBuffer()
//	 binary.Read()
func (a *ArrayBase) GetTo(idx int32, v interface{}) bool {

	if a.Cnt == 0 {
		return false
	}

	sz := binary.Size(v)
	if sz < 0 {
		panic(marshal.ErrNotFixedSize)
	}

	bs, ok := a.GetBytes(idx, sz)
	if ok {
		b := bytes.NewBuffer(bs)
		err := binary.Read(b, endian, v)
		if err != nil {
			panic(err)
		}
		return true
	}

	return false
}

func (a *ArrayBase) Get(idx int32) (interface{}, bool) {

	if a.Cnt == 0 {
		return nil, false
	}

	bs, ok := a.GetBytes(idx, a.EltMarshaler.GetMarshaledSize(nil))
	if ok {
		_, v := a.EltMarshaler.Unmarshal(bs)
		return v, true
	}

	return nil, false
}

// GetBytes is similar to Get but does not return the byte slice instead of
// unmarshaled data.
func (a *ArrayBase) GetBytes(idx int32, eltsize int) ([]byte, bool) {
	dataIndex, ok := a.GetEltIndex(idx)
	if !ok {
		return nil, false
	}

	stIdx := int32(eltsize) * dataIndex
	return a.Elts[stIdx : stIdx+int32(eltsize)], true
}

// appendIndex add a index into index bitmap.
// The `index` must be greater than any existent indexes.
func (a *ArrayBase) appendIndex(index int32) {

	iBm, iBit := bmBit(index)

	var bmWord = &a.Bitmaps[iBm]
	if *bmWord == 0 {
		a.Offsets[iBm] = a.Cnt
	}

	*bmWord |= uint64(1) << uint(iBit)

	a.Cnt++
}
