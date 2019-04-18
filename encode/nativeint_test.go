package encode_test

import (
	"math/bits"
	"testing"

	"github.com/openacid/slim/encode"
)

func TestInt(t *testing.T) {

	sz := bits.UintSize / 8

	v0 := [8]byte{}
	v1 := [8]byte{1}
	v1234 := [8]byte{0x34, 0x12}
	vneg := [8]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

	cases := []struct {
		input    int
		want     string
		wantsize int
	}{
		{0, string(v0[:sz]), sz},
		{1, string(v1[:sz]), sz},
		{0x1234, string(v1234[:sz]), sz},
		{^int(0), string(vneg[:sz]), sz},
	}

	m := encode.Int{}

	for i, c := range cases {
		rst := m.Encode(c.input)
		if string(rst) != c.want {
			t.Fatalf("%d-th: input: %v; want: %v; actual: %v",
				i+1, c.input, []byte(c.want), rst)
		}

		n := m.GetSize(c.input)
		if c.wantsize != n {
			t.Fatalf("%d-th: input: %v; wantsize: %v; actual: %v",
				i+1, c.input, c.wantsize, n)
		}

		n = m.GetEncodedSize(rst)
		if c.wantsize != n {
			t.Fatalf("%d-th: input: %v; wantsize: %v; actual: %v",
				i+1, c.input, c.wantsize, n)
		}

		n, u64 := m.Decode(rst)
		if c.input != u64 {
			t.Fatalf("%d-th: decode: input: %v; want: %v; actual: %v",
				i+1, c.input, c.input, u64)
		}
		if c.wantsize != n {
			t.Fatalf("%d-th: decoded size: input: %v; want: %v; actual: %v",
				i+1, c.input, c.wantsize, n)
		}
	}
}
