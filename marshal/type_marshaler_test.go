package marshal_test

import (
	"encoding/binary"
	"reflect"
	"testing"

	"github.com/openacid/errors"
	"github.com/openacid/slim/marshal"
)

type typeXY struct {
	X int32
	Y int32
}

func testPanic(t *testing.T, f func(), msg string) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic: %s", msg)
		}
	}()

	f()
}

func TestNewTypeMarshaler(t *testing.T) {

	m, _ := marshal.NewTypeMarshaler(int32(1))
	if m.Endian != binary.LittleEndian {
		t.Fatalf("expect default endian is %#v but %#v", binary.LittleEndian, m.Endian)
	}

	ii := int32(1)

	cases := []struct {
		input   interface{}
		want    *marshal.TypeMarshaler
		wanterr error
	}{
		{
			int(1),
			nil,
			marshal.ErrNotFixedSize,
		},
		{
			[]int32{1},
			nil,
			marshal.ErrNotFixedSize,
		},
		{
			int32(1),
			&marshal.TypeMarshaler{
				Endian: binary.LittleEndian,
				Type:   reflect.ValueOf(int32(1)).Type(),
				Size:   4,
			},
			nil,
		},
		{
			&ii,
			&marshal.TypeMarshaler{
				Endian: binary.LittleEndian,
				Type:   reflect.ValueOf(int32(1)).Type(),
				Size:   4,
			},
			nil,
		},
		{
			typeXY{1, 2},
			&marshal.TypeMarshaler{
				Endian: binary.LittleEndian,
				Type:   reflect.ValueOf(typeXY{}).Type(),
				Size:   8,
			},
			nil,
		},
		{
			&typeXY{1, 2},
			&marshal.TypeMarshaler{
				Endian: binary.LittleEndian,
				Type:   reflect.ValueOf(typeXY{}).Type(),
				Size:   8,
			},
			nil,
		},
	}

	for i, c := range cases {
		rst, err := marshal.NewTypeMarshalerEndian(c.input, nil)
		if errors.Cause(err) != c.wanterr {
			t.Fatalf("%d-th: input: %#v; wanterr: %#v; actual: %#v",
				i+1, c.input, c.wanterr, err)
		}

		if !reflect.DeepEqual(c.want, rst) {
			t.Fatalf("%d-th: input: %#v; want: %#v; actual: %#v",
				i+1, c.input, c.want, rst)
		}

		m, err := marshal.NewTypeMarshalerEndianByType(
			reflect.Indirect(reflect.ValueOf(c.input)).Type(), nil)
		if errors.Cause(err) != c.wanterr {
			t.Fatalf("%d-th: input: %#v; wanterr: %#v; actual: %#v",
				i+1, c.input, c.wanterr, err)
		}

		if !reflect.DeepEqual(c.want, m) {
			t.Fatalf("%d-th: input: %#v; want: %#v; actual: %#v",
				i+1, c.input, c.want, m)
		}
	}
}

func TestTypeMarshalerMarshal(t *testing.T) {

	m, err := marshal.NewTypeMarshalerEndian(int32(1), nil)
	if err != nil {
		t.Fatalf("expected no error but: %v", err)
	}

	testPanic(t, func() { m.Marshal(uint32(1)) }, "int32: uint32")
	testPanic(t, func() { m.Marshal([]int32{1}) }, "int32: []int32")

	// indirect value results in no panic
	ii := int32(5)
	bs := m.Marshal(&ii)
	want := []byte{5, 0, 0, 0}
	if !reflect.DeepEqual(want, bs) {
		t.Fatalf("want: %#v, but: %#v", want, bs)
	}

	cases := []struct {
		input interface{}
		want  []byte
	}{
		{
			int32(1),
			[]byte{1, 0, 0, 0},
		},
		{
			byte(1),
			[]byte{1},
		},
		{
			typeXY{1, 2},
			[]byte{1, 0, 0, 0, 2, 0, 0, 0},
		},
		{
			&typeXY{1, 2},
			[]byte{1, 0, 0, 0, 2, 0, 0, 0},
		},
	}

	for i, c := range cases {
		m, err := marshal.NewTypeMarshalerEndian(c.input, nil)
		if err != nil {
			t.Fatalf("%d-th: expected no error but: %#v", i+1, err)
		}

		n := m.GetSize(c.input)
		if n != binary.Size(c.input) {
			t.Fatalf("expect n to be %d but %d", binary.Size(c.input), n)
		}

		bs := m.Marshal(c.input)
		if !reflect.DeepEqual(c.want, bs) {
			t.Fatalf("%d-th: input: %#v; want: %#v; actual: %#v",
				i+1, c.input, c.want, bs)
		}
	}
}

func TestTypeMarshalerUnmarshal(t *testing.T) {

	ii := int32(5)

	cases := []struct {
		input interface{}
		want  interface{}
	}{
		{
			int32(1),
			int32(1),
		},
		{
			byte(1),
			byte(1),
		},
		{
			&ii,
			int32(5),
		},
		{
			typeXY{1, 2},
			typeXY{1, 2},
		},
		{
			&typeXY{1, 2},
			typeXY{1, 2},
		},
	}

	for i, c := range cases {
		m, err := marshal.NewTypeMarshalerEndian(c.input, nil)
		if err != nil {
			t.Fatalf("%d-th: expected no error but: %#v", i+1, err)
		}

		bs := m.Marshal(c.input)
		n, v := m.Unmarshal(bs)

		if n != m.Size {
			t.Fatalf("expect n to b %d but %d", m.Size, n)
		}

		if n != m.GetMarshaledSize(bs) {
			t.Fatalf("expect n to b %d but %d", m.GetMarshaledSize(bs), n)
		}

		if !reflect.DeepEqual(c.want, v) {
			t.Fatalf("%d-th: input: %#v; want: %#v; actual: %#v",
				i+1, c.input, c.want, v)
		}
	}
}
