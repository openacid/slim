package bitmap

import (
	"math/bits"
)

// IndexRank64 creates a rank index for a bitmap.
// rank(i) is defined as number of "1" upto position i, excluding i.
//
// It returns an index of []int32.
// Every element in it is rank(i*64)
//
// Since 0.1.11
// An optional bool specifies whether to add a last index entry of count of all
// "1".
//
// Since 0.1.8
func IndexRank64(words []uint64, opts ...bool) []int32 {

	trailing := false
	if len(opts) > 0 {
		trailing = opts[0]
	}

	l := len(words)
	if trailing {
		l++
	}

	idx := make([]int32, l)
	n := int32(0)
	for i := 0; i < len(words); i++ {
		idx[i] = n
		n += int32(bits.OnesCount64(words[i]))
	}

	if trailing {
		idx[len(words)] = n
	}

	return idx
}

// IndexRank128 creates a rank index for a bitmap.
// rank(i) is defined as number of "1" upto position i, excluding i.
//
// It returns an index of []int32.
// Every element in it is rank(i*128).
//
// It also adds a last index item if len(words) % 2 == 0, in order to make the
// distance from any bit to the closest index be less than 64.
//
//
// Since 0.1.8
func IndexRank128(words []uint64) []int32 {

	idx := make([]int32, 0)
	n := int32(0)
	for i := 0; i < len(words); i += 2 {
		idx = append(idx, n)
		n += int32(bits.OnesCount64(words[i]))
		if i < len(words)-1 {
			n += int32(bits.OnesCount64(words[i+1]))
		}
	}

	// Need a last index to let distance from every bit to its closest index
	// <=64
	if len(words)&1 == 0 {
		idx = append(idx, n)
	}

	// clone to reduce cap to len
	idx = append(idx[:0:0], idx...)

	return idx
}

// Rank128 returns the count of 1 up to i excluding i and bit value at i,
// with the help of a 128-bit index.
//
// It takes about 2~3 ns/op.
//
// Since 0.1.9
func Rank128(words []uint64, rindex []int32, i int32) (int32, int32) {

	// A precaculated count serves for two 64-bit word on the left and right.
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

	wordI := i >> 6
	j := uint32(i & 63)
	atRight := wordI & 1

	n := rindex[(i+64)>>7]
	w := words[wordI]

	cnt1 := int32(bits.OnesCount64(w))
	c1 := n - atRight*cnt1 + int32(bits.OnesCount64(w&Mask[j]))
	return c1, int32(w>>uint(j)) & 1
}

// Rank64 returns the count of 1 up to i excluding i and the bit value at i,
// with the help of a 64-bit index.
//
// It takes about 2~3 ns/op.
//
// Since 0.1.9
func Rank64(words []uint64, rindex []int32, i int32) (int32, int32) {

	wordI := i >> 6
	j := uint32(i & 63)

	n := rindex[wordI]
	w := words[wordI]

	c1 := n + int32(bits.OnesCount64(w&Mask[j]))
	return c1, int32(w>>uint(j)) & 1
}

// Tip: use static mask to speed up rank()
//
// calculate mask ((1<<i)-1):
//
//   BenchmarkRank64_5_bits-8        1000000000               2.38 ns/op
//   BenchmarkRank128_5_bits-8       1000000000               2.98 ns/op
//   BenchmarkRank64_64k_bits-8      1000000000               2.10 ns/op
//   BenchmarkRank128_64k_bits-8     1000000000               2.72 ns/op
//
// lookup table:
//
//   BenchmarkRank64_5_bits-8        2000000000               1.63 ns/op
//   BenchmarkRank128_5_bits-8       1000000000               2.45 ns/op
//   BenchmarkRank64_64k_bits-8      2000000000               1.45 ns/op
//   BenchmarkRank128_64k_bits-8     1000000000               2.34 ns/op
