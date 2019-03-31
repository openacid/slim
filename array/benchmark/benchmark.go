// Package benchmark provides benchmark utilities
package benchmark

import (
	"github.com/openacid/slim/array"
	"github.com/openacid/slim/benchhelper"
	"github.com/openacid/slim/encode"
)

// MemoryUsage describe settings for an Array and memory cost.
type MemoryUsage struct {

	// MaxIndex is the last and max index in an array.
	MaxIndex int32

	// EltSize is the size in byte of an array element.
	EltSize int

	// EltCnt is the number of element in an array.
	EltCnt int32

	// ArraySize is the memory cost in total for an array.
	ArraySize int64

	// UserDataSize is the user data size.
	UserDataSize int64

	// LoadFactor describe how full an array is and is roughly the ratio of
	// EltCnt to MaxIndex.
	// Its value is in (0, 1]
	LoadFactor float64

	// Overhead is how much else memory in byte is used by an array, except
	// user-data.
	// Its value is in (0, +oo)
	Overhead float64
}

// NewBytesArray creates a Array from "indexes" and "elts" of []byte.
func NewBytesArray(indexes []int32, elts [][]byte) (*array.Array, error) {
	a := &array.Array{}
	a.EltEncoder = encode.Bytes{Size: len(elts[0])}
	err := a.Init(indexes, elts)
	return a, err
}

// MemoryCostOf detect memory increment after creating "nArray" Array-s.
func MemoryCostOf(nArray int, indexes []int32, elts [][]byte) int64 {

	rss1 := benchhelper.Allocated()

	arr := []*array.Array{}

	for i := 0; i < nArray; i++ {
		a32, err := NewBytesArray(indexes, elts)
		if err != nil {
			panic(err)
		}
		arr = append(arr, a32)
	}

	rss2 := benchhelper.Allocated()

	for i := 0; i < nArray; i++ {
		_ = arr[i].Bitmaps
		_ = arr[i].Offsets
		_ = arr[i].Elts
	}

	return rss2 - rss1
}

// CollectMemoryUsage create several Array and check how much memory it costs.
func CollectMemoryUsage(factor float64, maxIdx int32, eltSize int) *MemoryUsage {
	nArray := 512

	indexes := benchhelper.RandI32SliceBetween(0, maxIdx, factor)
	eltCnt := len(indexes)

	elts := benchhelper.NewBytesSlices(eltSize, eltCnt)
	totalMem := MemoryCostOf(nArray, indexes, elts)

	userDataSize := int64(eltSize) * int64(eltCnt)
	arraySize := totalMem / int64(nArray)

	overhead := float64(arraySize)/float64(userDataSize) - 1

	return &MemoryUsage{
		MaxIndex:     indexes[len(indexes)-1],
		EltSize:      eltSize,
		EltCnt:       int32(eltCnt),
		ArraySize:    arraySize,
		UserDataSize: userDataSize,
		LoadFactor:   factor,
		Overhead:     overhead,
	}
}
