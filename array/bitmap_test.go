package array

import (
	"fmt"
	"testing"

	proto "github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"
)

func TestNewBitmap(t *testing.T) {

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
		for _, got := range []Bitmapper{
			NewBitmap(c.input),
			NewDenseBitmap(c.input),
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

func TestBitmap_Stat(t *testing.T) {

	ta := require.New(t)
	n := 1024 * 1024
	nums := make([]int32, n)
	for i := 0; i < n; i++ {
		nums[i] = int32(i * 7)
	}

	got := NewBitmap(nums)
	st := got.Stat()

	ta.Equal(int32(8), st["bits/one"])
	ta.True(st["mem_total"] > int32(3*n)/8)

	got = NewDenseBitmap(nums)
	st = got.Stat()

	ta.Equal(int32(7), st["bits/one"])
	ta.True(st["mem_total"] > int32(3*n)/8)
}

func TestBitmap_Has(t *testing.T) {

	ta := require.New(t)
	nums := []int32{1, 3, 64, 129}

	for _, got := range []Bitmapper{
		NewBitmap(nums),
		NewDenseBitmap(nums),
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

func TestBitmap_Rank_panic(t *testing.T) {

	ta := require.New(t)

	nums := []int32{1, 3, 64, 129}
	for _, got := range []Bitmapper{
		NewBitmap(nums),
		NewDenseBitmap(nums),
	} {

		ta.Panics(func() {
			got.Rank(-1)
		})

		// no panic
		_ = got.Rank(130)

		ta.Panics(func() {
			got.Rank(131)
		})
	}
}

func TestBitmap_MarshalUnmarshal(t *testing.T) {

	ta := require.New(t)

	nums := []int32{1, 3, 64, 129}
	for _, a := range []Bitmapper{
		NewBitmap(nums),
		NewDenseBitmap(nums),
	} {

		bytes, err := proto.Marshal(a)
		ta.Nil(err, "want no error but: %+v", err)

		b := &Bitmap{}

		err = proto.Unmarshal(bytes, b)
		ta.Nil(err, "want no error but: %+v", err)

		for _, j := range nums {
			ta.False(b.Has(j - 1))
			ta.True(b.Has(j))
		}
	}
}

func TestNewRandIndex(t *testing.T) {

	ta := require.New(t)

	cases := []struct {
		input []uint64
		want1 []int32
		want2 []int32
	}{
		{
			[]uint64{},
			[]int32{},
			[]int32{0},
		},
		{
			[]uint64{0},
			[]int32{0},
			[]int32{0},
		},
		{
			[]uint64{1},
			[]int32{0},
			[]int32{0},
		},
		{
			[]uint64{0xffffffffffffffff},
			[]int32{0},
			[]int32{0},
		},
		{
			[]uint64{0xffffffffffffffff, 1},
			[]int32{0, 64},
			[]int32{0, 65},
		},
		{
			[]uint64{0xffffffffffffffff, 1, 1},
			[]int32{0, 64, 65},
			[]int32{0, 65},
		},
		{
			[]uint64{0xffffffffffffffff, 1, 1, 3},
			[]int32{0, 64, 65, 66},
			[]int32{0, 65, 68},
		},
		{
			[]uint64{0xffffffffffffffff, 1, 1, 3, 4},
			[]int32{0, 64, 65, 66, 68},
			[]int32{0, 65, 68},
		},
	}

	for i, c := range cases {

		idx1 := newRankIndex1(c.input)
		ta.Equal(c.want1, idx1,
			"%d-th: input: %#v; want: %#v; got: %#v",
			i+1, c.input, c.want1, idx1)

		got := newRankIndex2(c.input)
		ta.Equal(c.want2, got,
			"%d-th: input: %#v; want: %#v; got: %#v",
			i+1, c.input, c.want2, got)
	}
}

var DsOutput int

func BenchmarkBitmap_Rank_Dense(b *testing.B) {

	got := NewDenseBitmap([]int32{1, 2, 3, 64, 129})

	b.ResetTimer()

	var gotrank int32
	for i := 0; i < b.N; i++ {
		gotrank += got.Rank(int32(i & 127))
	}

	DsOutput = int(gotrank)
}

func BenchmarkBitmap_Rank(b *testing.B) {

	got := NewBitmap([]int32{1, 2, 3, 64, 129})

	b.ResetTimer()

	var gotrank int32
	for i := 0; i < b.N; i++ {
		gotrank += got.Rank(int32(i & 127))
	}

	DsOutput = int(gotrank)
}
