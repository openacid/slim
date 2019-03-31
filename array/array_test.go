package array_test

import (
	"fmt"
	"reflect"
	"testing"

	proto "github.com/golang/protobuf/proto"
	"github.com/openacid/slim/array"
)

func TestArrayNewEmpty(t *testing.T) {

	type D struct {
		X int32
		Y int16
	}

	a, err := array.NewEmpty(D{})
	if err != nil {
		t.Fatalf("expect no err but: %v", err)
	}

	_ = a

	// TODO test marshal

}

func TestArrayOfU32(t *testing.T) {

	a := &array.Array{}

	err := a.Init([]int32{1, 2, 3}, []uint32{4, 5, 6})
	if err != nil {
		t.Fatalf("expect no err but: %v", err)
	}

	v, found := a.Get(1)
	if !found {
		t.Fatalf("expect: %v; but: %v", true, false)
	}
	if v == nil {
		t.Fatalf("v should not be nil expect: %v; but: %v", "not nil", v)
	}
}

func TestArrayOfStruct(t *testing.T) {

	a := &array.Array{}

	err := a.Init([]int32{10, 12, 13},
		[]struct {
			X int32
			Y uint16
		}{
			{1, 2},
			{3, 4},
			{5, 6},
		})
	if err != nil {
		t.Fatalf("expect no err but: %v", err)
	}

	v, found := a.Get(12)
	if !found {
		t.Fatalf("expect: %v; but: %v", true, false)
	}
	if v == nil {
		t.Fatalf("v should not be nil expect: %v; but: %v", "not nil", v)
	}
}

func TestArrayAndU32InterMarshal(t *testing.T) {

	// created with fmt.Printf("%#v\n",  data )
	serialized := []byte{
		0x8, 0x4, 0x12, 0x6, 0xa2, 0x4, 0x0, 0x0,
		0x80, 0x10, 0x1a, 0x4, 0x0, 0x0, 0x0, 0x3,
		0x22, 0x10, 0xc, 0x0, 0x0, 0x0, 0xf, 0x0,
		0x0, 0x0, 0x13, 0x0, 0x0, 0x0, 0x78, 0x0,
		0x0, 0x0,
	}

	index := []int32{1, 5, 9, 203}
	eltsData := []uint32{12, 15, 19, 120}

	arr, err := array.NewU32(index, eltsData)
	if err != nil {
		t.Fatalf("create array failure: %s", err)
	}

	data, err := proto.Marshal(arr)
	if err != nil {
		t.Fatalf("proto.Marshal: %s", err)
	}

	if !reflect.DeepEqual(serialized, data) {
		fmt.Println(serialized)
		fmt.Println(data)
		t.Fatalf("serialized data incorrect")
	}

	loaded, err := array.NewEmpty(uint32(0))
	if err != nil {
		t.Fatalf("expected no error but: %+v", err)
	}

	err = proto.Unmarshal(data, loaded)
	if err != nil {
		t.Fatalf("proto.Unmarshal: %+v", err)
	}

	second, err := proto.Marshal(loaded)
	if err != nil {
		t.Fatalf("proto.Marshal: %+v", err)
	}

	if !reflect.DeepEqual(serialized, second) {
		fmt.Println(serialized)
		fmt.Println(second)
		t.Fatalf("second serialized data incorrect")
	}
}

func BenchmarkArrayGet(b *testing.B) {
	indexes := []int32{0, 5, 9, 203, 400}
	elts := []uint32{12, 15, 19, 120, 300}
	a, _ := array.New(indexes, elts)

	for i := 0; i < b.N; i++ {
		a.Get(5)
	}
}
