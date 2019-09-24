package strhelper

import (
	"reflect"
	"testing"

	"github.com/openacid/low/bitword"
	"github.com/stretchr/testify/require"
)

func testPanic(t *testing.T, f func(), msg string) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic: %s", msg)
		}
	}()

	f()
}

func TestSliceToAndFromBitWords(t *testing.T) {

	cases := []struct {
		input []string
		n     int
		want  [][]byte
	}{
		{[]string{"a", "bc", "d"}, 4,
			[][]byte{
				{6, 1},
				{6, 2, 6, 3},
				{6, 4},
			},
		},
		{[]string{"a", "bc", "d"}, 2,
			[][]byte{
				{1, 2, 0, 1},
				{1, 2, 0, 2, 1, 2, 0, 3},
				{1, 2, 1, 0},
			},
		},
	}

	for i, c := range cases {
		rst := bitword.BitWord[c.n].FromStrs(c.input)
		if !reflect.DeepEqual(c.want, rst) {
			t.Fatalf("%d-th: input: %v; want: %v; actual: %v",
				i+1, c.input, c.want, rst)
		}
	}
}

func TestIntToBin(t *testing.T) {

	ta := require.New(t)

	cases := []struct {
		input interface{}
		want  string
	}{
		{int8(7), "11100000"},
		{uint8(7), "11100000"},
		{int16(0x0507), "11100000 10100000"},
		{uint16(0x0507), "11100000 10100000"},
		{int32(0x01030507), "11100000 10100000 11000000 10000000"},
		{uint32(0x01030507), "11100000 10100000 11000000 10000000"},
		{int64(0x0f01030507), "11100000 10100000 11000000 10000000 11110000 00000000 00000000 00000000"},
		{uint64(0x0f01030507), "11100000 10100000 11000000 10000000 11110000 00000000 00000000 00000000"},

		{[]int8{7}, "11100000"},
		{[]uint16{0x0507, 0x0102}, "11100000 10100000,01000000 10000000"},
	}

	for i, c := range cases {
		got := ToBin(c.input)
		ta.Equal(c.want, got,
			"%d-th: input: %#v; want: %#v; got: %#v",
			i+1, c.input, c.want, got)

	}
}
