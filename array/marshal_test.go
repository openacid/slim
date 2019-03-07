package array_test

import (
	"reflect"
	"testing"

	"github.com/openacid/slim/array"
)

func TestMarshalUnmarshal(t *testing.T) {

	indexes := []uint32{1, 5, 9, 203}
	elts := []uint16{12, 15, 19, 120}

	cases := []struct {
		n    int
		want []byte
	}{
		{
			0,
			[]byte{},
		},
		{
			1,
			[]byte{8, 1, 18, 1, 2, 26, 1, 0, 34, 2, 12, 0},
		},
		{
			2,
			[]byte{8, 2, 18, 1, 34, 26, 1, 0, 34, 4, 12, 0, 15, 0},
		},
	}

	for i, c := range cases {

		a, err := array.NewArrayU16(indexes[:c.n], elts[:c.n])
		if err != nil {
			t.Errorf("expect no error but: %s", err)
		}

		rst, err := array.Marshal(a)
		if err != nil {
			t.Errorf("expect no error but: %s", err)
		}

		if !reflect.DeepEqual(rst, c.want) {
			t.Fatalf("%d-th: n: %v; want: %v; actual: %v",
				i+1, c.n, c.want, rst)
		}

		// Unmarshal

		b := &array.ArrayU16{}
		nread, err := array.Unmarshal(b, rst)

		if err != nil {
			t.Errorf("expect no error but: %s", err)
		}

		if nread != len(rst) {
			t.Errorf("expcect to read %d but %d", len(rst), nread)
		}

		if !reflect.DeepEqual(a.Data, b.Data) {
			t.Fatalf("%d-th: n: %v; compare Elts: a: %v; b: %v",
				i+1, c.n, a.Data, b.Data)
		}

		// protobuf handles empty structure specially.
		if c.n == 0 {
			continue
		}

		if !reflect.DeepEqual(a, b) {
			t.Fatalf("%d-th: n: %v; compare Elts: a: %v; b: %v",
				i+1, c.n, a.Data, b.Data)
		}

	}
}

func TestMarshalUnmarshalBit(t *testing.T) {

	n := 102400
	step := 2
	indexes := []uint32{}
	elts := []uint16{}

	for i := 0; i < n; i += step {
		indexes = append(indexes, uint32(i))
		elts = append(elts, uint16(i))
	}

	a, err := array.NewArrayU16(indexes, elts)
	if err != nil {
		t.Errorf("expect no error but: %s", err)
	}

	rst, err := array.Marshal(a)
	if err != nil {
		t.Errorf("expect no error but: %s", err)
	}

	b := &array.ArrayU16{}
	nread, err := array.Unmarshal(b, rst)
	if err != nil {
		t.Errorf("expect no error but: %s", err)
	}

	if nread != len(rst) {
		t.Errorf("expcect to read %d but %d", len(rst), nread)
	}

	if !reflect.DeepEqual(a, b) {
		t.Fatalf("compare: a: %v; b: %v", a, b)
	}

}

func TestMarshalUnmarshal2Types(t *testing.T) {

	indexes := []uint32{1, 5, 9, 203}
	elts := []uint16{12, 15, 19, 120}

	a, _ := array.NewArrayU16(indexes, elts)
	rst, _ := array.Marshal(a)

	b, _ := array.New(array.U16Conv{}, []uint32{}, []uint16{})
	array.Unmarshal(b, rst)

	for _, i := range indexes {
		av, afound := a.Get2(i)
		bv, bfound := b.Get2(i)

		if av != bv || afound != bfound {
			t.Fatalf("expect same result i=%d, %d %t %d %t", i, av, afound, bv, bfound)
		}
	}

}
