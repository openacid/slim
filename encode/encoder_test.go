package encode_test

import (
	"testing"

	"github.com/openacid/slim/encode"
)

func TestString16(t *testing.T) {

	cases := []struct {
		input string
		want  int
	}{
		{"", 2},
		{"a", 3},
		{"abc", 5},
	}

	m := encode.String16{}

	for i, c := range cases {
		rst := m.Encode(c.input)
		if len(rst) != c.want {
			t.Fatalf("%d-th: encoded len: input: %v; want: %v; actual: %v",
				i+1, c.input, c.want, len(rst))
		}

		l := m.GetEncodedSize(rst)
		if l != c.want {
			t.Fatalf("%d-th: encoded size: input: %v; want: %v; actual: %v",
				i+1, c.input, c.want, l)
		}

		n, s := m.Decode(rst)
		if c.want != n {
			t.Fatalf("%d-th: decoded size: input: %v; want: %v; actual: %v",
				i+1, c.input, c.want, n)
		}
		if c.input != s {
			t.Fatalf("%d-th: decode: input: %v; want: %v; actual: %v",
				i+1, c.input, c.input, s)
		}
	}
}

func TestGetEncoder(t *testing.T) {

	cases := []struct {
		input   interface{}
		want    encode.Encoder
		wanterr error
	}{
		{
			uint16(0),
			encode.U16{},
			nil,
		},
		{
			uint32(0),
			encode.U32{},
			nil,
		},
		{
			uint64(0),
			encode.U64{},
			nil,
		},
		{
			[]int{},
			nil,
			encode.ErrUnknownEltType,
		},
		{
			nil,
			nil,
			encode.ErrUnknownEltType,
		},
	}

	for i, c := range cases {
		rst, err := encode.EncoderOf(c.input)
		if rst != c.want {
			t.Fatalf("%d-th: input: %v; want: %v; actual: %v",
				i+1, c.input, c.want, rst)
		}
		if err != c.wanterr {
			t.Fatalf("%d-th: input: %v; wanterr: %v; actual: %v",
				i+1, c.input, c.wanterr, err)
		}
	}
}

func TestGetSliceEltEncoder(t *testing.T) {

	cases := []struct {
		input   interface{}
		want    encode.Encoder
		wanterr error
	}{
		{
			[]uint16{},
			encode.U16{},
			nil,
		},
		{
			[]uint32{},
			encode.U32{},
			nil,
		},
		{
			[]uint64{},
			encode.U64{},
			nil,
		},
		{
			[]int{},
			nil,
			encode.ErrUnknownEltType,
		},
		{
			int(1),
			nil,
			encode.ErrNotSlice,
		},
	}

	for i, c := range cases {
		rst, err := encode.GetSliceEltEncoder(c.input)
		if rst != c.want {
			t.Fatalf("%d-th: input: %v; want: %v; actual: %v",
				i+1, c.input, c.want, rst)
		}
		if err != c.wanterr {
			t.Fatalf("%d-th: input: %v; wanterr: %v; actual: %v",
				i+1, c.input, c.wanterr, err)
		}
	}
}
