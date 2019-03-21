package marshal_test

import (
	"testing"

	"github.com/openacid/slim/marshal"
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

	m := marshal.String16{}

	for i, c := range cases {
		rst := m.Marshal(c.input)
		if len(rst) != c.want {
			t.Fatalf("%d-th: marshaled len: input: %v; want: %v; actual: %v",
				i+1, c.input, c.want, len(rst))
		}

		l := m.GetMarshaledSize(rst)
		if l != c.want {
			t.Fatalf("%d-th: marshaled size: input: %v; want: %v; actual: %v",
				i+1, c.input, c.want, l)
		}

		n, s := m.Unmarshal(rst)
		if c.want != n {
			t.Fatalf("%d-th: unmarshaled size: input: %v; want: %v; actual: %v",
				i+1, c.input, c.want, n)
		}
		if c.input != s {
			t.Fatalf("%d-th: unmarshal: input: %v; want: %v; actual: %v",
				i+1, c.input, c.input, s)
		}
	}
}

func TestGetMarshaler(t *testing.T) {

	cases := []struct {
		input   interface{}
		want    marshal.Marshaler
		wanterr error
	}{
		{
			uint16(0),
			marshal.U16{},
			nil,
		},
		{
			uint32(0),
			marshal.U32{},
			nil,
		},
		{
			uint64(0),
			marshal.U64{},
			nil,
		},
		{
			[]int{},
			nil,
			marshal.ErrUnknownEltType,
		},
		{
			nil,
			nil,
			marshal.ErrUnknownEltType,
		},
	}

	for i, c := range cases {
		rst, err := marshal.GetMarshaler(c.input)
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

func TestGetSliceEltMarshaler(t *testing.T) {

	cases := []struct {
		input   interface{}
		want    marshal.Marshaler
		wanterr error
	}{
		{
			[]uint16{},
			marshal.U16{},
			nil,
		},
		{
			[]uint32{},
			marshal.U32{},
			nil,
		},
		{
			[]uint64{},
			marshal.U64{},
			nil,
		},
		{
			[]int{},
			nil,
			marshal.ErrUnknownEltType,
		},
		{
			int(1),
			nil,
			marshal.ErrNotSlice,
		},
	}

	for i, c := range cases {
		rst, err := marshal.GetSliceEltMarshaler(c.input)
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
