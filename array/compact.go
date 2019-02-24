// Package array implements functions for the manipulation of compacted array.
package array

import (
	"errors"
	"reflect"
)

// ErrIndexLen is returned if number of indexes does not equal the number of
// datas, when initializing a Array.
var ErrIndexLen = errors.New("the length of index and elts must be equal")

// Array32 is a space efficient array implementation.
//
// Unlike a normal array, it does not allocate space for a element that there is
// not data in it.
type Array32 struct {
	Array32Index
	Converter
}

// New32 creates a Array and initializes it with a slice of index and a
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

	err := a.InitIndexBitmap(index)
	if err != nil {
		return err
	}

	a.Elts = make([]byte, 0, nElts)
	for i := 0; i < nElts; i++ {
		ee := rElts.Index(i).Interface()
		raw := a.Marshal(ee)
		a.Elts = append(a.Elts, raw...)
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
	raw, ok := a.GetBytes(idx, a.GetMarshaledSize(nil))
	if ok {
		_, val := a.Unmarshal(raw)
		return val, true
	}

	return nil, false
}

// GetBytes is similar to Get2 but does not return the byte slice instead of
// unmarshaled data.
func (a *Array32) GetBytes(idx uint32, eltsize int) ([]byte, bool) {
	dataIndex, ok := a.GetEltIndex(idx)
	if !ok {
		return nil, false
	}

	stIdx := uint32(eltsize) * dataIndex
	return a.Elts[stIdx : stIdx+uint32(eltsize)], true
}
