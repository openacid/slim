package trie

import (
	"testing"

	"github.com/openacid/low/bitmap"
	"github.com/stretchr/testify/require"
)

var InputI32 int32 = 35
var OutputI32 int32 = 0

func TestIndexSelect32(t *testing.T) {

	ta := require.New(t)

	cases := []struct {
		input []uint64
		want  []int32
	}{
		{nil, []int32{}},
		{[]uint64{}, []int32{}},
		{[]uint64{0}, []int32{}},
		{[]uint64{1}, []int32{0}},
		{[]uint64{2}, []int32{0}},
		{[]uint64{3}, []int32{0}},
		{[]uint64{4, 0}, []int32{0}},
		{[]uint64{0xffffffff, 0xffffffff}, []int32{0, 1}},
		{[]uint64{0xffffffff, 0xffffffff, 1}, []int32{0, 1, 2}},
		{[]uint64{0xffffffffffffffff, 0xffffffffffffffff, 1}, []int32{0, 0, 1, 1, 2}},
		{[]uint64{0, 0xffffffffffffffff, 1}, []int32{1, 1, 2}},
	}

	for i, c := range cases {
		got := indexSelect32(c.input)
		ta.Equal(c.want, got, "%d-th: case: %+v", i+1, c)
	}
}

func TestSelect32(t *testing.T) {

	ta := require.New(t)

	cases := []struct {
		input []uint64
	}{
		{nil},
		{[]uint64{}},
		{[]uint64{0}},
		{[]uint64{1}},
		{[]uint64{2}},
		{[]uint64{3}},
		{[]uint64{4, 0}},
		{[]uint64{0xf, 0xf}},
		{[]uint64{0xf, 0, 0xf}},
		{[]uint64{0xfffffffffffffff0}},
		{[]uint64{0xffffffffffffffff}},
		{[]uint64{0xffffffff, 0xffffffff}},
		{[]uint64{0xffffffff, 0xffffffff, 1}},
		{[]uint64{0x6668}}, // 000101100110011
	}

	for _, c := range cases {

		bm := &Bitmap{Words: c.input}
		bm.indexit("s32")

		all := bitmap.ToArray(c.input)

		for j := 0; j < len(all)-1; j++ {
			a, b := bm.select32(int32(j))
			ta.Equal(all[j], a, "select2-1: %d, case: %+v", j, c)
			ta.Equal(all[j+1], b, "select2-2: %d, case: %+v", j, c)
		}

	}
}

func BenchmarkSelect32(b *testing.B) {
	words := []uint64{0xffffffff, 0xffffffff, 1}
	var s int32

	bm := &Bitmap{Words: words}
	bm.indexit("s32")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		a, _ := bm.select32(InputI32)
		s += a
	}
	OutputI32 = s
}
