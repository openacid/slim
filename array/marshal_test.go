package array_test

import (
	"reflect"
	"testing"

	proto "github.com/golang/protobuf/proto"
	"github.com/kr/pretty"
	"github.com/openacid/slim/array"
	"github.com/openacid/slim/prototype"
)

func TestMarshalUnmarshal(t *testing.T) {

	indexes := []int32{1, 5, 9, 203}
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

		a, err := array.NewU16(indexes[:c.n], elts[:c.n])
		if err != nil {
			t.Errorf("expect no error but: %s", err)
		}

		rst, err := proto.Marshal(a)
		if err != nil {
			t.Errorf("expect no error but: %s", err)
		}

		if !reflect.DeepEqual(rst, c.want) {
			t.Fatalf("%d-th: n: %v; want: %v; actual: %v",
				i+1, c.n, c.want, rst)
		}

		// Unmarshal

		b := &array.U16{}
		err = proto.Unmarshal(rst, b)

		if err != nil {
			t.Errorf("expect no error but: %s", err)
		}

		if !reflect.DeepEqual(a.Elts, b.Elts) {
			t.Fatalf("%d-th: n: %v; compare Elts: a: %v; b: %v",
				i+1, c.n, a.Elts, b.Elts)
		}

		// protobuf handles empty structure specially.
		if c.n == 0 {
			continue
		}

		// ignore proto's field when compare
		a.XXX_sizecache = 0

		if !reflect.DeepEqual(a, b) {
			t.Fatalf("%d-th: n: %v; compare a b: %v",
				i+1, c.n, pretty.Diff(a, b))
		}

	}
}

func TestMarshalUnmarshalBit(t *testing.T) {

	n := 102400
	step := 2
	indexes := []int32{}
	elts := []uint16{}

	for i := 0; i < n; i += step {
		indexes = append(indexes, int32(i))
		elts = append(elts, uint16(i))
	}

	a, err := array.NewU16(indexes, elts)
	if err != nil {
		t.Errorf("expect no error but: %s", err)
	}

	rst, err := proto.Marshal(a)
	if err != nil {
		t.Errorf("expect no error but: %s", err)
	}

	b := &array.U16{}
	err = proto.Unmarshal(rst, b)
	if err != nil {
		t.Errorf("expect no error but: %s", err)
	}

	// proto pollute this field
	a.XXX_sizecache = 0
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("compare: a b: %v", pretty.Diff(a, b))
	}

}

func TestProtobufMarshalUnmarshal2Types(t *testing.T) {

	indexes := []int32{1, 5, 9, 203}
	elts := []uint16{12, 15, 19, 120}

	udArr, _ := array.NewU16(indexes, elts)
	convArr, _ := array.NewEmpty(uint16(0))

	// ud to conv

	buf, _ := proto.Marshal(udArr)
	err := proto.Unmarshal(buf, convArr)
	if err != nil {
		t.Fatalf("expect no err, but: %v", err)
	}

	checkArrayEqual(t, udArr, convArr, indexes)

	// conv to ud

	buf, err = proto.Marshal(convArr)
	if err != nil {
		t.Fatalf("expect no err, but: %v", err)
	}

	err = proto.Unmarshal(buf, udArr)
	if err != nil {
		t.Fatalf("expect no err, but: %v", err)
	}

	checkArrayEqual(t, udArr, convArr, indexes)
}

func checkArrayEqual(t *testing.T, udArr *array.U16, convArr *array.Array, indexes []int32) {
	for _, i := range indexes {
		av, afound := udArr.Get(i)
		bv, bfound := convArr.Get(i)

		if av != bv || afound != bfound {
			t.Fatalf("expect same result i=%d, %d %t %d %t", i, av, afound, bv, bfound)
		}
	}
}

func TestMigrateToSignedCntAndOffsets(t *testing.T) {
	// marshaled data from previous prototype.Array with uint32 Cnt and uint32 Offsets
	// message Array32 {
	//     uint32 Cnt              = 1;
	//     repeated uint64 Bitmaps = 2;
	//     repeated uint32 Offsets = 3;
	//     bytes  Elts             = 4;
	// }
	// prototype.Array32{
	//     Cnt: 0xffffffff,
	//     Bitmaps: []uint64{0},
	//     Offsets: []uint32{0xffffffff},
	//     Elts: []byte{},
	// }
	prevMarshalled := []byte{
		0x8, 0xff, 0xff, 0xff, 0xff, 0xf, 0x12, 0x1,
		0x0, 0x1a, 0x5, 0xff, 0xff, 0xff, 0xff, 0xf,
	}

	b := &prototype.Array32{}
	err := proto.Unmarshal(prevMarshalled, b)
	if err != nil {
		t.Errorf("expect no error but: %s", err)
	}

	if b.Cnt != -1 {
		t.Fatalf("expect -1 but: %v", b.Cnt)
	}

	if b.Offsets[0] != -1 {
		t.Fatalf("expect -1 but: %v", b.Offsets)
	}
}
