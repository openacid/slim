package marshal

import (
	"testing"
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

	m := String16{}

	for i, c := range cases {
		rst := m.Marshal(c.input)
		if len(rst) != c.want {
			t.Fatalf("%d-th: marshalled len: input: %v; want: %v; actual: %v",
				i+1, c.input, c.want, len(rst))
		}

		l := m.GetMarshaledSize(rst)
		if l != c.want {
			t.Fatalf("%d-th: marshaled size: input: %v; want: %v; actual: %v",
				i+1, c.input, c.want, l)
		}

		n, s := m.Unmarshal(rst)
		if c.want != n {
			t.Fatalf("%d-th: unmarshalled size: input: %v; want: %v; actual: %v",
				i+1, c.input, c.want, n)
		}
		if c.input != s {
			t.Fatalf("%d-th: unmarshal: input: %v; want: %v; actual: %v",
				i+1, c.input, c.input, s)
		}
	}
}

func TestU64(t *testing.T) {

	cases := []struct {
		input    uint64
		want     string
		wantsize int
	}{
		{0, "\x00\x00\x00\x00\x00\x00\x00\x00", 8},
		{1, "\x01\x00\x00\x00\x00\x00\x00\x00", 8},
		{0x1234, "\x34\x12\x00\x00\x00\x00\x00\x00", 8},
		{0xffffffffffffffff, "\xff\xff\xff\xff\xff\xff\xff\xff", 8},
	}

	m := U64{}

	for i, c := range cases {
		rst := m.Marshal(c.input)
		if string(rst) != c.want {
			t.Fatalf("%d-th: input: %v; want: %v; actual: %v",
				i+1, c.input, []byte(c.want), rst)
		}

		n := m.GetMarshaledSize(rst)
		if c.wantsize != n {
			t.Fatalf("%d-th: input: %v; wantsize: %v; actual: %v",
				i+1, c.input, c.wantsize, n)
		}

		n, u64 := m.Unmarshal(rst)
		if c.input != u64 {
			t.Fatalf("%d-th: unmarshal: input: %v; want: %v; actual: %v",
				i+1, c.input, c.input, u64)
		}
		if c.wantsize != n {
			t.Fatalf("%d-th: unmarshalled size: input: %v; want: %v; actual: %v",
				i+1, c.input, c.wantsize, n)
		}
	}
}

func TestU32(t *testing.T) {

	cases := []struct {
		input    uint32
		want     string
		wantsize int
	}{
		{0, "\x00\x00\x00\x00", 4},
		{1, "\x01\x00\x00\x00", 4},
		{0x1234, "\x34\x12\x00\x00", 4},
		{0xffffffff, "\xff\xff\xff\xff", 4},
	}

	m := U32{}

	for i, c := range cases {
		rst := m.Marshal(c.input)
		if string(rst) != c.want {
			t.Fatalf("%d-th: input: %v; want: %v; actual: %v",
				i+1, c.input, []byte(c.want), rst)
		}

		n := m.GetMarshaledSize(rst)
		if c.wantsize != n {
			t.Fatalf("%d-th: input: %v; wantsize: %v; actual: %v",
				i+1, c.input, c.wantsize, n)
		}

		n, u32 := m.Unmarshal(rst)
		if c.input != u32 {
			t.Fatalf("%d-th: unmarshal: input: %v; want: %v; actual: %v",
				i+1, c.input, c.input, u32)
		}
		if c.wantsize != n {
			t.Fatalf("%d-th: unmarshalled size: input: %v; want: %v; actual: %v",
				i+1, c.input, c.wantsize, n)
		}
	}
}

func TestU16(t *testing.T) {

	cases := []struct {
		input    uint16
		want     string
		wantsize int
	}{
		{0, "\x00\x00", 2},
		{1, "\x01\x00", 2},
		{0x1234, "\x34\x12", 2},
		{0xffff, "\xff\xff", 2},
	}

	m := U16{}

	for i, c := range cases {
		rst := m.Marshal(c.input)
		if string(rst) != c.want {
			t.Fatalf("%d-th: input: %v; want: %v; actual: %v",
				i+1, c.input, []byte(c.want), rst)
		}

		n := m.GetMarshaledSize(rst)
		if c.wantsize != n {
			t.Fatalf("%d-th: input: %v; wantsize: %v; actual: %v",
				i+1, c.input, c.wantsize, n)
		}

		n, u16 := m.Unmarshal(rst)
		if c.input != u16 {
			t.Fatalf("%d-th: unmarshal: input: %v; want: %v; actual: %v",
				i+1, c.input, c.input, u16)
		}
		if c.wantsize != n {
			t.Fatalf("%d-th: unmarshalled size: input: %v; want: %v; actual: %v",
				i+1, c.input, c.wantsize, n)
		}
	}
}

func TestGetMarshaller(t *testing.T) {

	cases := []struct {
		input   interface{}
		want    Marshaller
		wanterr error
	}{
		{
			uint16(0),
			U16{},
			nil,
		},
		{
			uint32(0),
			U32{},
			nil,
		},
		{
			uint64(0),
			U64{},
			nil,
		},
		{
			[]int{},
			nil,
			ErrUnknownEltType,
		},
		{
			nil,
			nil,
			ErrUnknownEltType,
		},
	}

	for i, c := range cases {
		rst, err := GetMarshaller(c.input)
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

func TestGetSliceEltMarshaller(t *testing.T) {

	cases := []struct {
		input   interface{}
		want    Marshaller
		wanterr error
	}{
		{
			[]uint16{},
			U16{},
			nil,
		},
		{
			[]uint32{},
			U32{},
			nil,
		},
		{
			[]uint64{},
			U64{},
			nil,
		},
		{
			[]int{},
			nil,
			ErrUnknownEltType,
		},
		{
			int(1),
			nil,
			ErrNotSlice,
		},
	}

	for i, c := range cases {
		rst, err := GetSliceEltMarshaller(c.input)
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
