package array

import (
	"fmt"
	"testing"

	proto "github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"
)

func TestNewBits(t *testing.T) {

	ta := require.New(t)

	cases := []struct {
		input     []int32
		wantn     int32
		wantwords []uint64
	}{
		{
			[]int32{},
			0,
			[]uint64{},
		},
		{
			[]int32{0},
			1,
			[]uint64{1},
		},
		{
			[]int32{0, 1, 2},
			3,
			[]uint64{7},
		},
		{
			[]int32{0, 1, 2, 63},
			64,
			[]uint64{(1 << 63) + 7},
		},
		{
			[]int32{64},
			65,
			[]uint64{0, 1},
		},
		{
			[]int32{1, 2, 3, 64, 129},
			130,
			[]uint64{0x0e, 1, 2},
		},
	}

	for i, c := range cases {
		for _, got := range []Bitmap{
			NewBits(c.input),
			NewDenseBits(c.input),
		} {
			ta.Equal(c.wantn, got.Len(),
				"%d-th: input: %#v; wantn: %#v; got: %#v",
				i+1, c.input, c.wantn, got.Len())

			ta.Equal(c.wantwords, got.Bits(),
				"%d-th: input: %#v; wantwords: %#v; got: %#v",
				i+1, c.input, c.wantwords, got.Bits())

			fmt.Println("input:", c.input)

			for wantrank, j := range c.input {
				gotrank := got.Rank(j)
				ta.Equal(int32(wantrank), gotrank, "%d-th: rank: j=%d", i+1, j)
			}
		}
	}
}

func TestNewBitsJoin(t *testing.T) {

	ta := require.New(t)

	subs := []uint64{1, 2, 3, 4, 5, 6, 7, 8, 9}

	for _, dense := range []bool{false, true} {

		b := NewBitsJoin(subs, 4, dense)
		ta.Equal([]uint64{1 + (2 << 4) + (3 << 8) + (4 << 12) + (5 << 16) + (6 << 20) + (7 << 24) + (8 << 28) + (9 << 32)}, b.Bits())

		b = NewBitsJoin(subs, 8, dense)
		ta.Equal([]uint64{
			1 + (2 << 8) + (3 << 16) + (4 << 24) + (5 << 32) + (6 << 40) + (7 << 48) + (8 << 56),
			9}, b.Bits())

		b = NewBitsJoin(subs, 16, dense)
		ta.Equal([]uint64{
			1 + (2 << 16) + (3 << 32) + (4 << 48),
			5 + (6 << 16) + (7 << 32) + (8 << 48),
			9}, b.Bits())
	}
}

func TestBits_Stat(t *testing.T) {

	ta := require.New(t)
	n := 1024 * 1024
	nums := make([]int32, n)
	for i := 0; i < n; i++ {
		nums[i] = int32(i * 7)
	}

	got := NewBits(nums)
	st := got.Stat()

	ta.Equal(int32(8), st["bits/one"])
	ta.True(st["mem_total"] > int32(3*n)/8)

	got = NewDenseBits(nums)
	st = got.Stat()

	ta.Equal(int32(7), st["bits/one"])
	ta.True(st["mem_total"] > int32(3*n)/8)
}

func TestBits_Has(t *testing.T) {

	ta := require.New(t)
	nums := []int32{1, 3, 64, 129}

	for _, got := range []Bitmap{
		NewBits(nums),
		NewDenseBits(nums),
	} {

		for _, j := range nums {
			ta.False(got.Has(j - 1))
			ta.True(got.Has(j))
		}

		ta.Panics(func() {
			got.Has(-1)
		})
		ta.Panics(func() {
			got.Has(130)
		})
	}
}

func TestBits_Rank_panic(t *testing.T) {

	ta := require.New(t)

	nums := []int32{1, 3, 64, 129}
	for _, got := range []Bitmap{
		NewBits(nums),
		NewDenseBits(nums),
	} {

		ta.Panics(func() {
			got.Rank(-1)
		})

		// no panic
		_ = got.Rank(130)
		_ = got.Rank(191)

		ta.Panics(func() {
			got.Rank(192)
		})
	}
}

func TestBits_MarshalUnmarshal(t *testing.T) {

	ta := require.New(t)

	nums := []int32{1, 3, 64, 129}
	for _, a := range []Bitmap{
		NewBits(nums),
		NewDenseBits(nums),
	} {

		bytes, err := proto.Marshal(a)
		ta.Nil(err, "want no error but: %+v", err)

		b := &Bits{}

		err = proto.Unmarshal(bytes, b)
		ta.Nil(err, "want no error but: %+v", err)

		for _, j := range nums {
			ta.False(b.Has(j - 1))
			ta.True(b.Has(j))
		}
	}
}

func TestConcatBitss_panic(t *testing.T) {

	ta := require.New(t)

	ta.Panics(func() { concatBits(nil, 0) })
	ta.Panics(func() { concatBits(nil, 3) })
	ta.Panics(func() { concatBits(nil, 5) })
	ta.Panics(func() { concatBits(nil, 6) })
	ta.Panics(func() { concatBits(nil, 7) })
	ta.Panics(func() { concatBits(nil, 9) })
	ta.Panics(func() { concatBits(nil, 15) })
	ta.Panics(func() { concatBits(nil, 17) })
	ta.Panics(func() { concatBits(nil, 31) })
	ta.Panics(func() { concatBits(nil, 33) })
	ta.Panics(func() { concatBits(nil, 63) })
	ta.Panics(func() { concatBits(nil, 65) })
}

func TestConcatBitss(t *testing.T) {

	ta := require.New(t)

	cases := []struct {
		elts   []uint64
		width  int32
		wantn  int32
		wantbm []uint64
	}{
		{
			[]uint64{},
			1,
			0,
			[]uint64{},
		},
		{
			[]uint64{1},
			1,
			1,
			[]uint64{1},
		},
		{
			[]uint64{1, 1},
			1,
			2,
			[]uint64{3},
		},
		{
			[]uint64{1, 1, 1},
			1,
			3,
			[]uint64{7},
		},
		{
			[]uint64{1, 1, 1},
			4,
			9,
			[]uint64{1 + (1 << 4) + (1 << 8)},
		},
		{
			[]uint64{1, 2, 5},
			4,
			11,
			[]uint64{1 + (2 << 4) + (5 << 8)},
		},
		{
			[]uint64{1, 2, 5, 1},
			16,
			49,
			[]uint64{1 + (2 << 16) + (5 << 32) + (1 << 48)},
		},
		{
			[]uint64{1, 2, 5, 5},
			16,
			51,
			[]uint64{1 + (2 << 16) + (5 << 32) + (5 << 48)},
		},
		{
			[]uint64{1, 2, 5, 5, 1},
			16,
			65,
			[]uint64{1 + (2 << 16) + (5 << 32) + (5 << 48), 1},
		},
		{
			[]uint64{1, 2, 5, 5, 1},
			16,
			65,
			[]uint64{1 + (2 << 16) + (5 << 32) + (5 << 48), 1},
		},
	}

	for i, c := range cases {

		gotn, gotbm := concatBits(c.elts, c.width)
		ta.Equal(c.wantn, gotn,
			"%d-th: input: %#v, %#v; want: %#v; got: %#v",
			i+1, c.elts, c.width, c.wantn, gotn)

		ta.Equal(c.wantbm, gotbm,
			"%d-th: input: %#v, %#v; want: %#v; got: %#v",
			i+1, c.elts, c.width, c.wantbm, gotbm)
	}
}

var DsOutput int

func BenchmarkBits_Rank_Dense(b *testing.B) {

	got := NewDenseBits([]int32{1, 2, 3, 64, 129})

	b.ResetTimer()

	var gotrank int32
	for i := 0; i < b.N; i++ {
		gotrank += got.Rank(int32(i & 127))
	}

	DsOutput = int(gotrank)
}

func BenchmarkBits_Rank(b *testing.B) {

	got := NewBits([]int32{1, 2, 3, 64, 129})

	b.ResetTimer()

	var gotrank int32
	for i := 0; i < b.N; i++ {
		gotrank += got.Rank(int32(i & 127))
	}

	DsOutput = int(gotrank)
}

func BenchmarkBits_Rank_50k(b *testing.B) {

	n := 64 * 1024
	mask := n - 1
	indexes := make([]int32, n)
	for i := 0; i < n; i++ {
		indexes[i] = int32(i * 2)
	}
	got := NewBits(indexes)

	b.ResetTimer()

	var gotrank int32
	for i := 0; i < b.N; i++ {
		gotrank += got.Rank(int32(i & mask))
	}

	DsOutput = int(gotrank)
}
