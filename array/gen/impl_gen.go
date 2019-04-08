package main

import (
	"github.com/openacid/slim/genhelper"
)

var implHead = `package array
import (
	"github.com/openacid/slim/bits"
)
`

var implTemplate = `
// {{.Name}} is an implementation of Base with {{.ValType}} element
//
// Since 0.2.0
type {{.Name}} struct {
	Base
}

// New{{.Name}} creates a {{.Name}}
//
// Since 0.2.0
func New{{.Name}}(index []int32, elts []{{.ValType}}) (a *{{.Name}}, err error) {
	a = &{{.Name}}{}
	err = a.Init(index, elts)
	if err != nil {
		a = nil
	}
	return a, err
}

// Get returns value at "idx" and a bool indicating if the value is
// found.
//
// Since 0.2.0
func (a *{{.Name}}) Get(idx int32) ({{.ValType}}, bool) {

	iBm, iBit := bmBit(idx)

	if iBm >= int32(len(a.Bitmaps)) {
		return 0, false
	}

	var n = a.Bitmaps[iBm]

	if ((n >> uint(iBit)) & 1) == 0 {
		return 0, false
	}

	cnt1 := bits.OnesCount64Before(n, uint(iBit))

	stIdx := a.Offsets[iBm]*{{.ValLen}} + int32(cnt1)*{{.ValLen}}

	return {{.ValType}}(endian.{{.Decoder}}(a.Elts[stIdx:])), true
}
`

var testHead = `package array_test

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
`

var testTemplate = `
func Test{{.Name}}NewErrorArgments(t *testing.T) {
	var index []int32
	eltsData := []{{.ValType}}{12, 15, 19, 120, 300}

	var err error

	index = []int32{1, 5, 9, 203}
	_, err = array.New{{.Name}}(index, eltsData)
	if err == nil {
		t.Fatalf("new with wrong index length must error")
	}

	index = []int32{1, 5, 5, 203, 400}
	_, err = array.New{{.Name}}(index, eltsData)
	if err == nil {
		t.Fatalf("new with unsorted index must error")
	}
}

func Test{{.Name}}New(t *testing.T) {
	var cases = []struct {
		index    []int32
		eltsData []{{.ValType}}
	}{
		{
			[]int32{}, []{{.ValType}}{},
		},
		{
			[]int32{0, 5, 9, 203, 400}, []{{.ValType}}{12, 15, 19, 120, 300},
		},
	}

	for _, c := range cases {
		index, eltsData := c.index, c.eltsData
		cnt := int32(len(index))

		a, err := array.New{{.Name}}(index, eltsData)
		if err != nil {
			t.Fatalf("failed new compacted array, err: %s", err)
		}

		if a.Cnt != cnt {
			t.Fatalf("cnt is not equal expect: %d, act: %d", cnt, a.Cnt)
		}

		buf := new(bytes.Buffer)
		_ = binary.Write(buf, binary.LittleEndian, eltsData)

		expElts := buf.Bytes()
		if expElts == nil {
			expElts = []byte{}
		}

		if !reflect.DeepEqual(a.Elts, expElts) && len(a.Elts) != 0 && len(expElts) != 0 {
			fmt.Println(pretty.Diff(a.Elts, expElts))
			t.Fatalf("elts is not equal expect: %d, act: %d", expElts, a.Elts)
		}
	}

}

func Test{{.Name}}Get(t *testing.T) {
	index, eltsData := []int32{}, []{{.ValType}}{}
	rnd := rand.New(rand.NewSource(time.Now().Unix()))

	keysMap := map[int32]bool{}
	num, idx, cnt := int32(0), int32(0), int32(1024)
	for {
		if rnd.Intn(2) == 1 {
			index = append(index, idx)
			eltsData = append(eltsData, {{.ValType}}(rnd.Uint64()))
			num++
			keysMap[idx] = true
		}
		idx++
		if num == cnt {
			break
		}
	}

	a, err := array.New{{.Name}}(index, eltsData)
	if err != nil {
		t.Fatalf("failed new compacted array, err: %s", err)
	}

	dataIdx := int32(0)
	for ii := int32(0); ii < idx; ii++ {

		actByte, found := a.Get(ii)
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

func Test{{.Name}}EncodeDecode(t *testing.T) {

	indexes := []int32{1, 5, 9, 203}
	elts := []{{.ValType}}{12, 15, 19, 120}

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

		a, err := array.New{{.Name}}(indexes[:c.n], elts[:c.n])
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
			want = append(c.want, byte(c.n*{{.ValLen}}))
			for i := 0; i < c.n; i++ {
				b := make([]byte, {{.ValLen}})
				binary.LittleEndian.Put{{.Decoder}}(b, {{.EncodeCast}}(elts[i]))
				want = append(want, b...)
			}
		}

		if !reflect.DeepEqual(rst, want) {
			t.Fatalf("%d-th: n: %v; want: %v; actual: %v",
				i+1, c.n, want, rst)
		}

		// Decode

		b := &array.{{.Name}}{}
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

func Test{{.Name}}EncodeDecodeBig(t *testing.T) {

	n := 102400
	step := 2
	indexes := []int32{}
	elts := []{{.ValType}}{}

	for i := 0; i < n; i += step {
		indexes = append(indexes, int32(i))
		elts = append(elts, {{.ValType}}(i))
	}

	a, err := array.New{{.Name}}(indexes, elts)
	if err != nil {
		t.Errorf("expect no error but: %s", err)
	}

	rst, err := proto.Marshal(a)
	if err != nil {
		t.Errorf("expect no error but: %s", err)
	}

	b := &array.{{.Name}}{}
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

func Benchmark{{.Name}}Get(b *testing.B) {
	a, err := array.New{{.Name}}([]int32{1, 2, 3}, []{{.ValType}}{1, 2, 3})
	if err != nil {
		panic(err)
	}
	for i := 0; i < b.N; i++ {
		a.Get(2)
	}
}
`

func main() {

	pref := "int"
	implfn := pref + ".go"
	testfn := pref + "_test.go"

	impls := []interface{}{
		genhelper.IntConfig{Name: "U16", ValType: "uint16", ValLen: 2, Decoder: "Uint16", EncodeCast: "uint16"},
		genhelper.IntConfig{Name: "U32", ValType: "uint32", ValLen: 4, Decoder: "Uint32", EncodeCast: "uint32"},
		genhelper.IntConfig{Name: "U64", ValType: "uint64", ValLen: 8, Decoder: "Uint64", EncodeCast: "uint64"},
		genhelper.IntConfig{Name: "I16", ValType: "int16", ValLen: 2, Decoder: "Uint16", EncodeCast: "uint16"},
		genhelper.IntConfig{Name: "I32", ValType: "int32", ValLen: 4, Decoder: "Uint32", EncodeCast: "uint32"},
		genhelper.IntConfig{Name: "I64", ValType: "int64", ValLen: 8, Decoder: "Uint64", EncodeCast: "uint64"},
	}

	genhelper.Render(implfn, implHead, implTemplate, impls, []string{"gofmt", "unconvert"})
	genhelper.Render(testfn, testHead, testTemplate, impls, []string{"gofmt", "unconvert"})
}
