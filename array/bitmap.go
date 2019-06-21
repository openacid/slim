package array

import (
	"fmt"
	"math/bits"

	proto "github.com/golang/protobuf/proto"
	"github.com/openacid/low/bitmap"
	"github.com/openacid/low/size"
)

const (
	BitsFlagDenseRank = 0x00000001
)

// Bitmap defines behavior of a bitmap.
//
// Since 0.5.4
type Bitmap interface {

	// Stat returns a map describing memory usage.
	//
	//    bits/one  :9
	//    mem_total :1195245
	//
	// Since 0.5.4
	Stat() map[string]int32

	// Has returns a bool indicating whether a bit is set.
	//
	// Since 0.5.4
	Has(int32) bool

	// Len returns number of bits in it.
	//
	// Since 0.5.4
	Len() int32

	// Bits returns all uint64 words in it.
	//
	// Since 0.5.4
	Bits() []uint64

	// Rank returns the count of 1 up to i(exclude i)
	//
	// Since 0.5.4
	Rank(int32) int32

	proto.Message
}

// NewBits creates a new Bitmap instance from a serias of int32.
// The input specifies what bit to set to 1.
//
// Since 0.5.4
func NewBits(nums []int32) Bitmap {

	n, words, index := newBitsData(nums)

	bm := &Bits{
		Flags:     0,
		N:         n,
		Words:     words,
		RankIndex: index,
	}
	return bm
}

// NewBitsJoin creates a new Bitmap instance from a serias of sub bitmap.
//
// Since 0.5.4
func NewBitsJoin(elts []uint64, eltWidth int32, dense bool) Bitmap {

	n, words := concatBits(elts, eltWidth)
	index := bitmap.IndexRank128(words)

	bm := &Bits{
		Flags: 0,
		N:     n,
		Words: words,
	}

	if dense {
		bm.Flags |= BitsFlagDenseRank
		bm.RankIndexDense = NewPolyArray(index)
	} else {
		bm.RankIndex = index
	}

	return bm
}

// NewDenseBits creates a new Bitmap instance from a serias of int32.
// The input specifies what bit to set to 1.
//
// It compress rand index to reduce memory cost.
// But increase query time.
//
// Since 0.5.4
func NewDenseBits(nums []int32) Bitmap {

	n, words, index := newBitsData(nums)

	d := NewPolyArray(index)

	bm := &Bits{
		Flags:          BitsFlagDenseRank,
		N:              n,
		Words:          words,
		RankIndexDense: d,
	}
	return bm
}

func newBitsData(nums []int32) (int32, []uint64, []int32) {

	n, words := newBitsWords(nums)
	index := bitmap.IndexRank128(words)

	return n, words, index
}

func newBitsWords(nums []int32) (int32, []uint64) {

	n := int32(0)
	if len(nums) > 0 {
		n = nums[len(nums)-1] + 1
	}

	nWords := (n + 63) >> 6
	words := make([]uint64, nWords)

	for _, i := range nums {
		iWord := i >> 6
		i = i & 63
		words[iWord] |= 1 << uint(i)
	}
	return n, words
}

func concatBits(elts []uint64, width int32) (int32, []uint64) {

	switch width {
	case 1, 2, 4, 8, 16, 32, 64:
	default:
		panic(fmt.Sprintf("width must be 1, 2, 4, 8, 16, 32, 64 but: %d", width))
	}

	wcap := int(64 / width)
	l := len(elts)

	nWords := (l + wcap - 1) / wcap
	words := make([]uint64, nWords)

	for i, bm := range elts {
		iWord := i / wcap
		i = i % wcap
		words[iWord] |= bm << (uint(i) * uint(width))
	}

	if len(words) == 0 {
		return 0, words
	}

	last := words[len(words)-1]
	n := (nWords * 64) - bits.LeadingZeros64(last)
	return int32(n), words
}

// Stat returns a map describing memory usage.
//
//    bits/one  :9
//    mem_total :1195245
//
// Since 0.5.4
func (b *Bits) Stat() map[string]int32 {

	totalmem := size.Of(b)

	st := map[string]int32{
		"mem_total": int32(totalmem),
		"bits/one":  int32(totalmem) * 8 / b.Rank(b.N),
	}

	return st
}

// Has return a bool indicating whether a bit is set.
//
// Since 0.5.4
func (b *Bits) Has(i int32) bool {
	if i < 0 || i >= b.N {
		panic(fmt.Sprintf("i=%d out of range, n=%d", i, b.N))
	}
	iWord := i >> 6
	j := uint(i & 63)
	return ((b.Words[iWord] >> j) & 1) != 0
}

// Len returns number of bits in it.
//
// Since 0.5.4
func (b *Bits) Len() int32 {
	return b.N
}

// Bits returns all uint64 words in it.
//
// Since 0.5.4
func (b *Bits) Bits() []uint64 {
	return b.Words
}

// Rand returns the count of 1 up to i(exclude i)
//
// It takes about 4 ns/op with uncompressed rank,
// and 15 ns/op with dense rank
//
// Since 0.5.4
func (b *Bits) Rank(i int32) int32 {

	// An precaculated count serves for two 64-bit word on the left and right.
	//
	// Idx[1] = OnesCount(bits[0:128])
	// Idx[1] serves for rank query from 64 to 192
	//
	//	   0       64      128     192     256
	//	   |-------+-------+-------+-------+
	//	   Idx[0]          Idx[1]          Idx[2]
	//
	// let j = i % 128
	//
	//	   If 0 <= j < 64:   use index at i/128
	//	   If 64 <= j < 128: use index at i/128 + 1

	iWord := uint64(i >> 6)

	var n int32
	if b.Flags&BitsFlagDenseRank == 0 {
		n = b.RankIndex[(i+64)>>7]
	} else {
		n = b.RankIndexDense.Get((i + 64) >> 7)
	}

	if iWord&1 == 0 {
		word := b.Words[iWord] << (64 - uint(i&63))
		d := int32(bits.OnesCount64(word))
		return n + d
	} else {
		word := b.Words[iWord] >> uint(i&63)
		d := int32(bits.OnesCount64(word))
		return n - d
	}
}
