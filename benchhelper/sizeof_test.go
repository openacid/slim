package benchhelper_test

import (
	"testing"

	. "github.com/openacid/slim/benchhelper"
)

const UintSize = 32 << (^uint(0) >> 32 & 1) // 32 or 64

func TestSizeOf(t *testing.T) {

	cases := []struct {
		input interface{}
		want  int
	}{
		{"abc", 3},

		{uint8(0), 1}, {int8(0), 1},
		{uint16(0), 2}, {int16(0), 2},
		{uint32(0), 4}, {int32(0), 4},
		{uint64(0), 8}, {int64(0), 8},

		{float32(0), 4},
		{float64(0), 8},

		{complex64(complex(1, 2)), 8},
		{complex128(complex(1, 2)), 16},
		{int(0), UintSize / 8},

		{[]int32{1, 2}, 8},
		{[3]int32{1, 2, 3}, 12},

		{map[int32]string{1: "a", 2: "b"}, 10},

		{struct{ a, b int64 }{1, 2}, 16},
		{&struct{ a, b int64 }{1, 2}, 16},
	}

	for i, c := range cases {
		rst := SizeOf(c.input)
		if rst != c.want {
			t.Fatalf("%d-th: input: %#v; want: %#v; actual: %#v",
				i+1, c.input, c.want, rst)
		}
	}
}
