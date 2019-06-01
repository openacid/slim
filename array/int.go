// Code generated 'by go generate ./...'; DO NOT EDIT.

package array

import (
	"math/bits"
)

// U16 is an implementation of Base with uint16 element
//
// Since 0.2.0
type U16 struct {
	Base
}

// NewU16 creates a U16
//
// Since 0.2.0
func NewU16(index []int32, elts []uint16) (a *U16, err error) {
	a = &U16{}
	err = a.Init(index, elts)
	if err != nil {
		a = nil
	}
	return a, err
}

// Get returns value at "idx" and a bool indicating if the value is
// found.
//
// Since 0.2.0
func (a *U16) Get(idx int32) (uint16, bool) {

	iBm, iBit := bmBit(idx)

	var n = a.Bitmaps[iBm]

	if ((n >> uint(iBit)) & 1) == 0 {
		return 0, false
	}

	cnt1 := bits.OnesCount64(n & ((uint64(1) << uint(iBit)) - 1))

	stIdx := a.Offsets[iBm]*2 + int32(cnt1)*2

	return endian.Uint16(a.Elts[stIdx:]), true
}

// U32 is an implementation of Base with uint32 element
//
// Since 0.2.0
type U32 struct {
	Base
}

// NewU32 creates a U32
//
// Since 0.2.0
func NewU32(index []int32, elts []uint32) (a *U32, err error) {
	a = &U32{}
	err = a.Init(index, elts)
	if err != nil {
		a = nil
	}
	return a, err
}

// Get returns value at "idx" and a bool indicating if the value is
// found.
//
// Since 0.2.0
func (a *U32) Get(idx int32) (uint32, bool) {

	iBm, iBit := bmBit(idx)

	var n = a.Bitmaps[iBm]

	if ((n >> uint(iBit)) & 1) == 0 {
		return 0, false
	}

	cnt1 := bits.OnesCount64(n & ((uint64(1) << uint(iBit)) - 1))

	stIdx := a.Offsets[iBm]*4 + int32(cnt1)*4

	return endian.Uint32(a.Elts[stIdx:]), true
}

// U64 is an implementation of Base with uint64 element
//
// Since 0.2.0
type U64 struct {
	Base
}

// NewU64 creates a U64
//
// Since 0.2.0
func NewU64(index []int32, elts []uint64) (a *U64, err error) {
	a = &U64{}
	err = a.Init(index, elts)
	if err != nil {
		a = nil
	}
	return a, err
}

// Get returns value at "idx" and a bool indicating if the value is
// found.
//
// Since 0.2.0
func (a *U64) Get(idx int32) (uint64, bool) {

	iBm, iBit := bmBit(idx)

	var n = a.Bitmaps[iBm]

	if ((n >> uint(iBit)) & 1) == 0 {
		return 0, false
	}

	cnt1 := bits.OnesCount64(n & ((uint64(1) << uint(iBit)) - 1))

	stIdx := a.Offsets[iBm]*8 + int32(cnt1)*8

	return endian.Uint64(a.Elts[stIdx:]), true
}

// I16 is an implementation of Base with int16 element
//
// Since 0.2.0
type I16 struct {
	Base
}

// NewI16 creates a I16
//
// Since 0.2.0
func NewI16(index []int32, elts []int16) (a *I16, err error) {
	a = &I16{}
	err = a.Init(index, elts)
	if err != nil {
		a = nil
	}
	return a, err
}

// Get returns value at "idx" and a bool indicating if the value is
// found.
//
// Since 0.2.0
func (a *I16) Get(idx int32) (int16, bool) {

	iBm, iBit := bmBit(idx)

	var n = a.Bitmaps[iBm]

	if ((n >> uint(iBit)) & 1) == 0 {
		return 0, false
	}

	cnt1 := bits.OnesCount64(n & ((uint64(1) << uint(iBit)) - 1))

	stIdx := a.Offsets[iBm]*2 + int32(cnt1)*2

	return int16(endian.Uint16(a.Elts[stIdx:])), true
}

// I32 is an implementation of Base with int32 element
//
// Since 0.2.0
type I32 struct {
	Base
}

// NewI32 creates a I32
//
// Since 0.2.0
func NewI32(index []int32, elts []int32) (a *I32, err error) {
	a = &I32{}
	err = a.Init(index, elts)
	if err != nil {
		a = nil
	}
	return a, err
}

// Get returns value at "idx" and a bool indicating if the value is
// found.
//
// Since 0.2.0
func (a *I32) Get(idx int32) (int32, bool) {

	iBm, iBit := bmBit(idx)

	var n = a.Bitmaps[iBm]

	if ((n >> uint(iBit)) & 1) == 0 {
		return 0, false
	}

	cnt1 := bits.OnesCount64(n & ((uint64(1) << uint(iBit)) - 1))

	stIdx := a.Offsets[iBm]*4 + int32(cnt1)*4

	return int32(endian.Uint32(a.Elts[stIdx:])), true
}

// I64 is an implementation of Base with int64 element
//
// Since 0.2.0
type I64 struct {
	Base
}

// NewI64 creates a I64
//
// Since 0.2.0
func NewI64(index []int32, elts []int64) (a *I64, err error) {
	a = &I64{}
	err = a.Init(index, elts)
	if err != nil {
		a = nil
	}
	return a, err
}

// Get returns value at "idx" and a bool indicating if the value is
// found.
//
// Since 0.2.0
func (a *I64) Get(idx int32) (int64, bool) {

	iBm, iBit := bmBit(idx)

	var n = a.Bitmaps[iBm]

	if ((n >> uint(iBit)) & 1) == 0 {
		return 0, false
	}

	cnt1 := bits.OnesCount64(n & ((uint64(1) << uint(iBit)) - 1))

	stIdx := a.Offsets[iBm]*8 + int32(cnt1)*8

	return int64(endian.Uint64(a.Elts[stIdx:])), true
}
