package benchhelper_test

import (
	"testing"

	. "github.com/openacid/slim/benchhelper"
	"github.com/stretchr/testify/require"
)

const UintSize = 32 << (^uint(0) >> 32 & 1) // 32 or 64

type myInt interface {
	Read() int
}
type myTyp int32

func (m *myTyp) Read() int {
	return 1
}

var mytyp myTyp
var myint myInt = &mytyp
var myintNil myInt = nil

func TestSizeOf(t *testing.T) {

	ta := require.New(t)

	cases := []struct {
		input interface{}
		want  int
	}{
		{"abc", 16 + 3},

		{uint8(0), 1}, {int8(0), 1},
		{uint16(0), 2}, {int16(0), 2},
		{uint32(0), 4}, {int32(0), 4},
		{uint64(0), 8}, {int64(0), 8},

		{float32(0), 4},
		{float64(0), 8},

		{complex64(complex(1, 2)), 8},
		{complex128(complex(1, 2)), 16},
		{int(0), UintSize / 8},

		{[]int32{1, 2}, 24 + 8},
		{[3]int32{1, 2, 3}, 12},

		{map[int32]string{1: "a", 2: "b"}, 8 + (4 + (16 + 1)) + (4 + (16 + 1))},

		{struct{ a, b int64 }{1, 2}, 16},
		{&struct{ a, b int64 }{1, 2}, 8 + 16},

		{myintNil, 0},
		{myint, 8 + 4},
		{struct{ a myInt }{}, 16},
		{struct{ a myInt }{&mytyp}, 16 + 8 + 4},
	}

	for i, c := range cases {
		rst := SizeOf(c.input)
		ta.Equal(c.want, rst, "%d-th: input: %+v", i+1, c.input)
	}
}

func TestSizeStat(t *testing.T) {

	ta := require.New(t)

	type my struct {
		a []int32
		b [3]int32
		c map[string]int8
		d *my
		e []*my
		f []string
		g myInt
		h myInt
	}

	v := my{
		a: []int32{1, 2, 3},
		b: [3]int32{4, 5, 6},
		c: map[string]int8{
			"abc": 3,
		},
		d: &my{
			a: []int32{1, 2},
		},
		e: []*my{
			{
				a: []int32{1, 2, 3},
			},
			{
				a: []int32{2, 3, 4},
			},
		},
		f: []string{
			"abc",
			"def",
		},
		g: nil,
		h: &mytyp,
	}

	want10 := `
benchhelper_test.my: 658
    a: []int32: 36
        0: int32: 4
        1: int32: 4
        2: int32: 4
    b: [3]int32: 12
        0: int32: 4
        1: int32: 4
        2: int32: 4
    c: map[string]int8: 28
        abc: int8: 1
    d: *benchhelper_test.my: 148
        benchhelper_test.my: 140
            a: []int32: 32
                0: int32: 4
                1: int32: 4
            b: [3]int32: 12
                0: int32: 4
                1: int32: 4
                2: int32: 4
            c: map[string]int8: 8
            d: *benchhelper_test.my: 8
            e: []*benchhelper_test.my: 24
            f: []string: 24
            g: benchhelper_test.myInt: 16
                <nil>
            h: benchhelper_test.myInt: 16
                <nil>
    e: []*benchhelper_test.my: 328
        0: *benchhelper_test.my: 152
            benchhelper_test.my: 144
                a: []int32: 36
                    0: int32: 4
                    1: int32: 4
                    2: int32: 4
                b: [3]int32: 12
                    0: int32: 4
                    1: int32: 4
                    2: int32: 4
                c: map[string]int8: 8
                d: *benchhelper_test.my: 8
                e: []*benchhelper_test.my: 24
                f: []string: 24
                g: benchhelper_test.myInt: 16
                    <nil>
                h: benchhelper_test.myInt: 16
                    <nil>
        1: *benchhelper_test.my: 152
            benchhelper_test.my: 144
                a: []int32: 36
                    0: int32: 4
                    1: int32: 4
                    2: int32: 4
                b: [3]int32: 12
                    0: int32: 4
                    1: int32: 4
                    2: int32: 4
                c: map[string]int8: 8
                d: *benchhelper_test.my: 8
                e: []*benchhelper_test.my: 24
                f: []string: 24
                g: benchhelper_test.myInt: 16
                    <nil>
                h: benchhelper_test.myInt: 16
                    <nil>
    f: []string: 62
        0: string: 19
        1: string: 19
    g: benchhelper_test.myInt: 16
        <nil>
    h: benchhelper_test.myInt: 28
        *benchhelper_test.myTyp: 12
            benchhelper_test.myTyp: 4`[1:]

	got10 := SizeStat(v, 10, 100)
	ta.Equal(want10, got10)

	want3 := `
benchhelper_test.my: 658
    a: []int32: 36
        0: int32: 4
        1: int32: 4
        2: int32: 4
    b: [3]int32: 12
        0: int32: 4
        1: int32: 4
        2: int32: 4
    c: map[string]int8: 28
        abc: int8: 1
    d: *benchhelper_test.my: 148
        benchhelper_test.my: 140
            a: []int32: 32
            b: [3]int32: 12
            c: map[string]int8: 8
            d: *benchhelper_test.my: 8
            e: []*benchhelper_test.my: 24
            f: []string: 24
            g: benchhelper_test.myInt: 16
            h: benchhelper_test.myInt: 16
    e: []*benchhelper_test.my: 328
        0: *benchhelper_test.my: 152
            benchhelper_test.my: 144
        1: *benchhelper_test.my: 152
            benchhelper_test.my: 144
    f: []string: 62
        0: string: 19
        1: string: 19
    g: benchhelper_test.myInt: 16
        <nil>
    h: benchhelper_test.myInt: 28
        *benchhelper_test.myTyp: 12
            benchhelper_test.myTyp: 4`[1:]
	got3 := SizeStat(v, 3, 100)
	ta.Equal(want3, got3)

	want32 := `
benchhelper_test.my: 658
    a: []int32: 36
        0: int32: 4
        1: int32: 4
    b: [3]int32: 12
        0: int32: 4
        1: int32: 4
    c: map[string]int8: 28
        abc: int8: 1
    d: *benchhelper_test.my: 148
        benchhelper_test.my: 140
            a: []int32: 32
            b: [3]int32: 12
            c: map[string]int8: 8
            d: *benchhelper_test.my: 8
            e: []*benchhelper_test.my: 24
            f: []string: 24
            g: benchhelper_test.myInt: 16
            h: benchhelper_test.myInt: 16
    e: []*benchhelper_test.my: 328
        0: *benchhelper_test.my: 152
            benchhelper_test.my: 144
        1: *benchhelper_test.my: 152
            benchhelper_test.my: 144
    f: []string: 62
        0: string: 19
        1: string: 19
    g: benchhelper_test.myInt: 16
        <nil>
    h: benchhelper_test.myInt: 28
        *benchhelper_test.myTyp: 12
            benchhelper_test.myTyp: 4`[1:]
	got32 := SizeStat(v, 3, 2)
	ta.Equal(want32, got32)
}
