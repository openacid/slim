package marshal_test

// Do NOT edit. re-generate this file with "go generate ./..."

import (
	"testing"

	"github.com/openacid/slim/marshal"
)

func TestU16(t *testing.T) {

	v0 := [8]byte{}
	v1 := [8]byte{1}
	v1234 := [8]byte{0x34, 0x12}
	vneg := [8]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

	cases := []struct {
		input    uint16
		want     string
		wantsize int
	}{
		{0, string(v0[:2]), 2},
		{1, string(v1[:2]), 2},
		{0x1234, string(v1234[:2]), 2},
		{^uint16(0), string(vneg[:2]), 2},
	}

	m := marshal.U16{}

	for i, c := range cases {
		rst := m.Marshal(c.input)
		if string(rst) != c.want {
			t.Fatalf("%d-th: input: %v; want: %v; actual: %v",
				i+1, c.input, []byte(c.want), rst)
		}

		n := m.GetSize(c.input)
		if c.wantsize != n {
			t.Fatalf("%d-th: input: %v; wantsize: %v; actual: %v",
				i+1, c.input, c.wantsize, n)
		}

		n = m.GetMarshaledSize(rst)
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
			t.Fatalf("%d-th: unmarshaled size: input: %v; want: %v; actual: %v",
				i+1, c.input, c.wantsize, n)
		}
	}
}

func TestU32(t *testing.T) {

	v0 := [8]byte{}
	v1 := [8]byte{1}
	v1234 := [8]byte{0x34, 0x12}
	vneg := [8]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

	cases := []struct {
		input    uint32
		want     string
		wantsize int
	}{
		{0, string(v0[:4]), 4},
		{1, string(v1[:4]), 4},
		{0x1234, string(v1234[:4]), 4},
		{^uint32(0), string(vneg[:4]), 4},
	}

	m := marshal.U32{}

	for i, c := range cases {
		rst := m.Marshal(c.input)
		if string(rst) != c.want {
			t.Fatalf("%d-th: input: %v; want: %v; actual: %v",
				i+1, c.input, []byte(c.want), rst)
		}

		n := m.GetSize(c.input)
		if c.wantsize != n {
			t.Fatalf("%d-th: input: %v; wantsize: %v; actual: %v",
				i+1, c.input, c.wantsize, n)
		}

		n = m.GetMarshaledSize(rst)
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
			t.Fatalf("%d-th: unmarshaled size: input: %v; want: %v; actual: %v",
				i+1, c.input, c.wantsize, n)
		}
	}
}

func TestU64(t *testing.T) {

	v0 := [8]byte{}
	v1 := [8]byte{1}
	v1234 := [8]byte{0x34, 0x12}
	vneg := [8]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

	cases := []struct {
		input    uint64
		want     string
		wantsize int
	}{
		{0, string(v0[:8]), 8},
		{1, string(v1[:8]), 8},
		{0x1234, string(v1234[:8]), 8},
		{^uint64(0), string(vneg[:8]), 8},
	}

	m := marshal.U64{}

	for i, c := range cases {
		rst := m.Marshal(c.input)
		if string(rst) != c.want {
			t.Fatalf("%d-th: input: %v; want: %v; actual: %v",
				i+1, c.input, []byte(c.want), rst)
		}

		n := m.GetSize(c.input)
		if c.wantsize != n {
			t.Fatalf("%d-th: input: %v; wantsize: %v; actual: %v",
				i+1, c.input, c.wantsize, n)
		}

		n = m.GetMarshaledSize(rst)
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
			t.Fatalf("%d-th: unmarshaled size: input: %v; want: %v; actual: %v",
				i+1, c.input, c.wantsize, n)
		}
	}
}

func TestI16(t *testing.T) {

	v0 := [8]byte{}
	v1 := [8]byte{1}
	v1234 := [8]byte{0x34, 0x12}
	vneg := [8]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

	cases := []struct {
		input    int16
		want     string
		wantsize int
	}{
		{0, string(v0[:2]), 2},
		{1, string(v1[:2]), 2},
		{0x1234, string(v1234[:2]), 2},
		{^int16(0), string(vneg[:2]), 2},
	}

	m := marshal.I16{}

	for i, c := range cases {
		rst := m.Marshal(c.input)
		if string(rst) != c.want {
			t.Fatalf("%d-th: input: %v; want: %v; actual: %v",
				i+1, c.input, []byte(c.want), rst)
		}

		n := m.GetSize(c.input)
		if c.wantsize != n {
			t.Fatalf("%d-th: input: %v; wantsize: %v; actual: %v",
				i+1, c.input, c.wantsize, n)
		}

		n = m.GetMarshaledSize(rst)
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
			t.Fatalf("%d-th: unmarshaled size: input: %v; want: %v; actual: %v",
				i+1, c.input, c.wantsize, n)
		}
	}
}

func TestI32(t *testing.T) {

	v0 := [8]byte{}
	v1 := [8]byte{1}
	v1234 := [8]byte{0x34, 0x12}
	vneg := [8]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

	cases := []struct {
		input    int32
		want     string
		wantsize int
	}{
		{0, string(v0[:4]), 4},
		{1, string(v1[:4]), 4},
		{0x1234, string(v1234[:4]), 4},
		{^int32(0), string(vneg[:4]), 4},
	}

	m := marshal.I32{}

	for i, c := range cases {
		rst := m.Marshal(c.input)
		if string(rst) != c.want {
			t.Fatalf("%d-th: input: %v; want: %v; actual: %v",
				i+1, c.input, []byte(c.want), rst)
		}

		n := m.GetSize(c.input)
		if c.wantsize != n {
			t.Fatalf("%d-th: input: %v; wantsize: %v; actual: %v",
				i+1, c.input, c.wantsize, n)
		}

		n = m.GetMarshaledSize(rst)
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
			t.Fatalf("%d-th: unmarshaled size: input: %v; want: %v; actual: %v",
				i+1, c.input, c.wantsize, n)
		}
	}
}

func TestI64(t *testing.T) {

	v0 := [8]byte{}
	v1 := [8]byte{1}
	v1234 := [8]byte{0x34, 0x12}
	vneg := [8]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

	cases := []struct {
		input    int64
		want     string
		wantsize int
	}{
		{0, string(v0[:8]), 8},
		{1, string(v1[:8]), 8},
		{0x1234, string(v1234[:8]), 8},
		{^int64(0), string(vneg[:8]), 8},
	}

	m := marshal.I64{}

	for i, c := range cases {
		rst := m.Marshal(c.input)
		if string(rst) != c.want {
			t.Fatalf("%d-th: input: %v; want: %v; actual: %v",
				i+1, c.input, []byte(c.want), rst)
		}

		n := m.GetSize(c.input)
		if c.wantsize != n {
			t.Fatalf("%d-th: input: %v; wantsize: %v; actual: %v",
				i+1, c.input, c.wantsize, n)
		}

		n = m.GetMarshaledSize(rst)
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
			t.Fatalf("%d-th: unmarshaled size: input: %v; want: %v; actual: %v",
				i+1, c.input, c.wantsize, n)
		}
	}
}
