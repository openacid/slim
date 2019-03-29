package array_test

// Do NOT edit. re-generate this file with "go generate ./..."

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"

	proto "github.com/golang/protobuf/proto"
	"github.com/kr/pretty"
	"github.com/openacid/slim/array"
)

func TestU16NewErrorArgments(t *testing.T) {
	var index []int32
	eltsData := []uint16{12, 15, 19, 120, 300}

	var err error

	index = []int32{1, 5, 9, 203}
	_, err = array.NewU16(index, eltsData)
	if err == nil {
		t.Fatalf("new with wrong index length must error")
	}

	index = []int32{1, 5, 5, 203, 400}
	_, err = array.NewU16(index, eltsData)
	if err == nil {
		t.Fatalf("new with unsorted index must error")
	}
}

func TestU16New(t *testing.T) {
	var cases = []struct {
		index    []int32
		eltsData []uint16
	}{
		{
			[]int32{}, []uint16{},
		},
		{
			[]int32{0, 5, 9, 203, 400}, []uint16{12, 15, 19, 120, 300},
		},
	}

	for _, c := range cases {
		index, eltsData := c.index, c.eltsData
		cnt := int32(len(index))

		ca, err := array.NewU16(index, eltsData)
		if err != nil {
			t.Fatalf("failed new compacted array, err: %s", err)
		}

		if ca.Cnt != cnt {
			t.Fatalf("cnt is not equal expect: %d, act: %d", cnt, ca.Cnt)
		}

		buf := new(bytes.Buffer)
		_ = binary.Write(buf, binary.LittleEndian, eltsData)

		expElts := buf.Bytes()
		if expElts == nil {
			expElts = []byte{}
		}

		if !reflect.DeepEqual(ca.Elts, expElts) && len(ca.Elts) != 0 && len(expElts) != 0 {
			fmt.Println(pretty.Diff(ca.Elts, expElts))
			t.Fatalf("elts is not equal expect: %d, act: %d", expElts, ca.Elts)
		}
	}

}

func TestU16Get(t *testing.T) {
	index, eltsData := []int32{}, []uint16{}
	rnd := rand.New(rand.NewSource(time.Now().Unix()))

	keysMap := map[int32]bool{}
	num, idx, cnt := int32(0), int32(0), int32(1024)
	for {
		if rnd.Intn(2) == 1 {
			index = append(index, idx)
			eltsData = append(eltsData, uint16(rnd.Uint64()))
			num++
			keysMap[idx] = true
		}
		idx++
		if num == cnt {
			break
		}
	}

	ca, err := array.NewU16(index, eltsData)
	if err != nil {
		t.Fatalf("failed new compacted array, err: %s", err)
	}

	dataIdx := int32(0)
	for ii := int32(0); ii < idx; ii++ {

		actByte, found := ca.Get(ii)
		_, present := keysMap[ii]
		if found != present {
			t.Fatalf("Get i:%d present:%t but:%t", ii, present, found)
		}

		if found {
			if eltsData[dataIdx] != actByte {
				t.Fatalf("Get i:%d is not equal expect: %d, act: %d", ii, eltsData[dataIdx], actByte)
			}
		}

		if _, ok := keysMap[ii]; ok {
			dataIdx++
		}
	}
}

func TestU16MarshalUnmarshal(t *testing.T) {

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
			[]byte{8, 1, 18, 1, 2, 26, 1, 0, 34},
		},
		{
			2,
			[]byte{8, 2, 18, 1, 34, 26, 1, 0, 34},
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

		// build Elts part for template generated test codes
		var want []byte = c.want
		if c.n > 0 {
			want = append(c.want, byte(c.n*2))
			for i := 0; i < c.n; i++ {
				b := make([]byte, 2)
				binary.LittleEndian.PutUint16(b, elts[i])
				want = append(want, b...)
			}
		}

		if !reflect.DeepEqual(rst, want) {
			t.Fatalf("%d-th: n: %v; want: %v; actual: %v",
				i+1, c.n, want, rst)
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

func TestU16MarshalUnmarshalBig(t *testing.T) {

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

func BenchmarkU16Get(b *testing.B) {
	a, err := array.NewU16([]int32{1, 2, 3}, []uint16{1, 2, 3})
	if err != nil {
		panic(err)
	}
	for i := 0; i < b.N; i++ {
		a.Get(2)
	}
}

func TestU32NewErrorArgments(t *testing.T) {
	var index []int32
	eltsData := []uint32{12, 15, 19, 120, 300}

	var err error

	index = []int32{1, 5, 9, 203}
	_, err = array.NewU32(index, eltsData)
	if err == nil {
		t.Fatalf("new with wrong index length must error")
	}

	index = []int32{1, 5, 5, 203, 400}
	_, err = array.NewU32(index, eltsData)
	if err == nil {
		t.Fatalf("new with unsorted index must error")
	}
}

func TestU32New(t *testing.T) {
	var cases = []struct {
		index    []int32
		eltsData []uint32
	}{
		{
			[]int32{}, []uint32{},
		},
		{
			[]int32{0, 5, 9, 203, 400}, []uint32{12, 15, 19, 120, 300},
		},
	}

	for _, c := range cases {
		index, eltsData := c.index, c.eltsData
		cnt := int32(len(index))

		ca, err := array.NewU32(index, eltsData)
		if err != nil {
			t.Fatalf("failed new compacted array, err: %s", err)
		}

		if ca.Cnt != cnt {
			t.Fatalf("cnt is not equal expect: %d, act: %d", cnt, ca.Cnt)
		}

		buf := new(bytes.Buffer)
		_ = binary.Write(buf, binary.LittleEndian, eltsData)

		expElts := buf.Bytes()
		if expElts == nil {
			expElts = []byte{}
		}

		if !reflect.DeepEqual(ca.Elts, expElts) && len(ca.Elts) != 0 && len(expElts) != 0 {
			fmt.Println(pretty.Diff(ca.Elts, expElts))
			t.Fatalf("elts is not equal expect: %d, act: %d", expElts, ca.Elts)
		}
	}

}

func TestU32Get(t *testing.T) {
	index, eltsData := []int32{}, []uint32{}
	rnd := rand.New(rand.NewSource(time.Now().Unix()))

	keysMap := map[int32]bool{}
	num, idx, cnt := int32(0), int32(0), int32(1024)
	for {
		if rnd.Intn(2) == 1 {
			index = append(index, idx)
			eltsData = append(eltsData, uint32(rnd.Uint64()))
			num++
			keysMap[idx] = true
		}
		idx++
		if num == cnt {
			break
		}
	}

	ca, err := array.NewU32(index, eltsData)
	if err != nil {
		t.Fatalf("failed new compacted array, err: %s", err)
	}

	dataIdx := int32(0)
	for ii := int32(0); ii < idx; ii++ {

		actByte, found := ca.Get(ii)
		_, present := keysMap[ii]
		if found != present {
			t.Fatalf("Get i:%d present:%t but:%t", ii, present, found)
		}

		if found {
			if eltsData[dataIdx] != actByte {
				t.Fatalf("Get i:%d is not equal expect: %d, act: %d", ii, eltsData[dataIdx], actByte)
			}
		}

		if _, ok := keysMap[ii]; ok {
			dataIdx++
		}
	}
}

func TestU32MarshalUnmarshal(t *testing.T) {

	indexes := []int32{1, 5, 9, 203}
	elts := []uint32{12, 15, 19, 120}

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
			[]byte{8, 1, 18, 1, 2, 26, 1, 0, 34},
		},
		{
			2,
			[]byte{8, 2, 18, 1, 34, 26, 1, 0, 34},
		},
	}

	for i, c := range cases {

		a, err := array.NewU32(indexes[:c.n], elts[:c.n])
		if err != nil {
			t.Errorf("expect no error but: %s", err)
		}

		rst, err := proto.Marshal(a)
		if err != nil {
			t.Errorf("expect no error but: %s", err)
		}

		// build Elts part for template generated test codes
		var want []byte = c.want
		if c.n > 0 {
			want = append(c.want, byte(c.n*4))
			for i := 0; i < c.n; i++ {
				b := make([]byte, 4)
				binary.LittleEndian.PutUint32(b, elts[i])
				want = append(want, b...)
			}
		}

		if !reflect.DeepEqual(rst, want) {
			t.Fatalf("%d-th: n: %v; want: %v; actual: %v",
				i+1, c.n, want, rst)
		}

		// Unmarshal

		b := &array.U32{}
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

func TestU32MarshalUnmarshalBig(t *testing.T) {

	n := 102400
	step := 2
	indexes := []int32{}
	elts := []uint32{}

	for i := 0; i < n; i += step {
		indexes = append(indexes, int32(i))
		elts = append(elts, uint32(i))
	}

	a, err := array.NewU32(indexes, elts)
	if err != nil {
		t.Errorf("expect no error but: %s", err)
	}

	rst, err := proto.Marshal(a)
	if err != nil {
		t.Errorf("expect no error but: %s", err)
	}

	b := &array.U32{}
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

func BenchmarkU32Get(b *testing.B) {
	a, err := array.NewU32([]int32{1, 2, 3}, []uint32{1, 2, 3})
	if err != nil {
		panic(err)
	}
	for i := 0; i < b.N; i++ {
		a.Get(2)
	}
}

func TestU64NewErrorArgments(t *testing.T) {
	var index []int32
	eltsData := []uint64{12, 15, 19, 120, 300}

	var err error

	index = []int32{1, 5, 9, 203}
	_, err = array.NewU64(index, eltsData)
	if err == nil {
		t.Fatalf("new with wrong index length must error")
	}

	index = []int32{1, 5, 5, 203, 400}
	_, err = array.NewU64(index, eltsData)
	if err == nil {
		t.Fatalf("new with unsorted index must error")
	}
}

func TestU64New(t *testing.T) {
	var cases = []struct {
		index    []int32
		eltsData []uint64
	}{
		{
			[]int32{}, []uint64{},
		},
		{
			[]int32{0, 5, 9, 203, 400}, []uint64{12, 15, 19, 120, 300},
		},
	}

	for _, c := range cases {
		index, eltsData := c.index, c.eltsData
		cnt := int32(len(index))

		ca, err := array.NewU64(index, eltsData)
		if err != nil {
			t.Fatalf("failed new compacted array, err: %s", err)
		}

		if ca.Cnt != cnt {
			t.Fatalf("cnt is not equal expect: %d, act: %d", cnt, ca.Cnt)
		}

		buf := new(bytes.Buffer)
		_ = binary.Write(buf, binary.LittleEndian, eltsData)

		expElts := buf.Bytes()
		if expElts == nil {
			expElts = []byte{}
		}

		if !reflect.DeepEqual(ca.Elts, expElts) && len(ca.Elts) != 0 && len(expElts) != 0 {
			fmt.Println(pretty.Diff(ca.Elts, expElts))
			t.Fatalf("elts is not equal expect: %d, act: %d", expElts, ca.Elts)
		}
	}

}

func TestU64Get(t *testing.T) {
	index, eltsData := []int32{}, []uint64{}
	rnd := rand.New(rand.NewSource(time.Now().Unix()))

	keysMap := map[int32]bool{}
	num, idx, cnt := int32(0), int32(0), int32(1024)
	for {
		if rnd.Intn(2) == 1 {
			index = append(index, idx)
			eltsData = append(eltsData, rnd.Uint64())
			num++
			keysMap[idx] = true
		}
		idx++
		if num == cnt {
			break
		}
	}

	ca, err := array.NewU64(index, eltsData)
	if err != nil {
		t.Fatalf("failed new compacted array, err: %s", err)
	}

	dataIdx := int32(0)
	for ii := int32(0); ii < idx; ii++ {

		actByte, found := ca.Get(ii)
		_, present := keysMap[ii]
		if found != present {
			t.Fatalf("Get i:%d present:%t but:%t", ii, present, found)
		}

		if found {
			if eltsData[dataIdx] != actByte {
				t.Fatalf("Get i:%d is not equal expect: %d, act: %d", ii, eltsData[dataIdx], actByte)
			}
		}

		if _, ok := keysMap[ii]; ok {
			dataIdx++
		}
	}
}

func TestU64MarshalUnmarshal(t *testing.T) {

	indexes := []int32{1, 5, 9, 203}
	elts := []uint64{12, 15, 19, 120}

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
			[]byte{8, 1, 18, 1, 2, 26, 1, 0, 34},
		},
		{
			2,
			[]byte{8, 2, 18, 1, 34, 26, 1, 0, 34},
		},
	}

	for i, c := range cases {

		a, err := array.NewU64(indexes[:c.n], elts[:c.n])
		if err != nil {
			t.Errorf("expect no error but: %s", err)
		}

		rst, err := proto.Marshal(a)
		if err != nil {
			t.Errorf("expect no error but: %s", err)
		}

		// build Elts part for template generated test codes
		var want []byte = c.want
		if c.n > 0 {
			want = append(c.want, byte(c.n*8))
			for i := 0; i < c.n; i++ {
				b := make([]byte, 8)
				binary.LittleEndian.PutUint64(b, elts[i])
				want = append(want, b...)
			}
		}

		if !reflect.DeepEqual(rst, want) {
			t.Fatalf("%d-th: n: %v; want: %v; actual: %v",
				i+1, c.n, want, rst)
		}

		// Unmarshal

		b := &array.U64{}
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

func TestU64MarshalUnmarshalBig(t *testing.T) {

	n := 102400
	step := 2
	indexes := []int32{}
	elts := []uint64{}

	for i := 0; i < n; i += step {
		indexes = append(indexes, int32(i))
		elts = append(elts, uint64(i))
	}

	a, err := array.NewU64(indexes, elts)
	if err != nil {
		t.Errorf("expect no error but: %s", err)
	}

	rst, err := proto.Marshal(a)
	if err != nil {
		t.Errorf("expect no error but: %s", err)
	}

	b := &array.U64{}
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

func BenchmarkU64Get(b *testing.B) {
	a, err := array.NewU64([]int32{1, 2, 3}, []uint64{1, 2, 3})
	if err != nil {
		panic(err)
	}
	for i := 0; i < b.N; i++ {
		a.Get(2)
	}
}
