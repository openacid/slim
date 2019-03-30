package array_test

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/kr/pretty"
	"github.com/openacid/errors"
	"github.com/openacid/slim/array"
	"github.com/openacid/slim/encode"
)

func randIndexes(cnt int32) []int32 {

	indexes := []int32{}
	rnd := rand.New(rand.NewSource(time.Now().Unix()))

	for i := int32(0); cnt > 0; i++ {
		if rnd.Intn(2) == 1 {
			indexes = append(indexes, i)
			cnt--
		}
	}

	return indexes
}

func makeIndexMap(idx []int32) map[int32]bool {

	indexMap := map[int32]bool{}
	for _, i := range idx {
		indexMap[i] = true
	}

	return indexMap
}

func TestBaseInitIndex(t *testing.T) {

	cases := []struct {
		input       []int32
		wantCnt     int32
		wantBitmaps []uint64
		wantOffsets []int32
		wanterr     error
	}{
		{
			[]int32{},
			0,
			[]uint64{},
			[]int32{},
			nil,
		},
		{
			[]int32{0},
			1,
			[]uint64{1},
			[]int32{0},
			nil,
		},
		{
			[]int32{1, 0},
			0,
			[]uint64{},
			[]int32{},
			array.ErrIndexNotAscending,
		},
		{
			[]int32{1, 3},
			2,
			[]uint64{0x0a},
			[]int32{0},
			nil,
		},
		{
			[]int32{1, 65},
			2,
			[]uint64{0x02, 0x02},
			[]int32{0, 1},
			nil,
		},
	}

	for i, c := range cases {
		ab := &array.Base{}
		err := ab.InitIndex(c.input)

		if errors.Cause(err) != c.wanterr {
			t.Fatalf("expect error: %v but: %v", c.wanterr, err)
		}

		if err == nil {

			if c.wantCnt != ab.Cnt {
				t.Fatalf("%d-th: input: %v; want: %v; actual: %v",
					i+1, c.input, c.wantCnt, ab.Cnt)
			}

			if !reflect.DeepEqual(c.wantOffsets, ab.Offsets) {
				t.Fatalf("%d-th: input: %v; want: %v; actual: %v",
					i+1, c.input, c.wantOffsets, ab.Offsets)
			}

			if !reflect.DeepEqual(c.wantBitmaps, ab.Bitmaps) {
				t.Fatalf("%d-th: input: %v; want: %v; actual: %v",
					i+1, c.input, c.wantBitmaps, ab.Bitmaps)
			}
		}
	}
}

func testPanic(t *testing.T, f func(), msg string) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic: %s", msg)
		}
	}()

	f()
}

func TestBaseInit(t *testing.T) {

	ab := &array.Base{}
	testPanic(t, func() { _ = ab.Init([]int32{}, 1) }, "elts int")
	testPanic(t, func() { _ = ab.Init([]int32{1, 2}, []interface{}{uint32(1), uint64(1)}) }, "different type elts")

	cases := []struct {
		input       []int32
		elts        interface{}
		wantBitmaps []uint64
		wantOffsets []int32
		wantElts    []byte
		wanterr     error
	}{
		{
			[]int32{}, []int{},
			[]uint64{},
			[]int32{},
			[]byte{},
			nil,
		},
		{
			[]int32{}, []int{1},
			[]uint64{},
			[]int32{},
			[]byte{},
			array.ErrIndexLen,
		},
		{
			[]int32{1, 0}, []int{1, 1},
			[]uint64{},
			[]int32{},
			[]byte{},
			array.ErrIndexNotAscending,
		},
		{
			[]int32{0}, []int{1},
			[]uint64{},
			[]int32{},
			[]byte{},
			encode.ErrNotFixedSize,
		},
		{
			[]int32{0}, []byte{1},
			[]uint64{1},
			[]int32{0},
			[]byte{1},
			nil,
		},
		{
			[]int32{1, 3}, []int16{1, 2},
			[]uint64{0x0a},
			[]int32{0},
			[]byte{1, 0, 2, 0},
			nil,
		},
		{
			[]int32{1, 65}, []int32{2, 3},
			[]uint64{0x02, 0x02},
			[]int32{0, 1},
			[]byte{2, 0, 0, 0, 3, 0, 0, 0},
			nil,
		},
		{
			[]int32{1, 65}, []interface{}{uint32(2), uint32(3)},
			[]uint64{0x02, 0x02},
			[]int32{0, 1},
			[]byte{2, 0, 0, 0, 3, 0, 0, 0},
			nil,
		},
	}

	for i, c := range cases {
		ab := &array.Base{}
		err := ab.Init(c.input, c.elts)

		if errors.Cause(err) != c.wanterr {
			t.Fatalf("%d-th, expect error: %v but: %v", i+1, c.wanterr, err)
		}

		if err == nil {

			if !reflect.DeepEqual(c.wantOffsets, ab.Offsets) {
				t.Fatalf("%d-th: input: %v; want: %v; actual: %v",
					i+1, c.input, c.wantOffsets, ab.Offsets)
			}

			if !reflect.DeepEqual(c.wantBitmaps, ab.Bitmaps) {
				t.Fatalf("%d-th: input: %v; want: %v; actual: %v",
					i+1, c.input, c.wantBitmaps, ab.Bitmaps)
			}

			if len(c.wantElts) != 0 && len(ab.Elts) != 0 && !reflect.DeepEqual(c.wantElts, ab.Elts) {
				t.Fatalf("%d-th: input: %v; want: %v; actual: %v",
					i+1, c.input, c.wantElts, ab.Elts)
			}
		}
	}
}

type intEncoder struct{}

func (m intEncoder) Encode(v interface{}) []byte {
	return []byte{byte(v.(int))}
}
func (m intEncoder) Decode(d []byte) (int, interface{}) {
	return 1, int(d[0])
}
func (m intEncoder) GetSize(v interface{}) int {
	return 1
}
func (m intEncoder) GetEncodedSize(d []byte) int {
	return 1
}

func TestBaseInitWithEncoder(t *testing.T) {
	ab := &array.Base{}
	ab.EltEncoder = intEncoder{}
	err := ab.Init([]int32{1, 2}, []int{3, 4})
	if err != nil {
		t.Fatalf("expected no error but: %#v", err)
	}

	wantelts := []byte{3, 4}
	if !reflect.DeepEqual(wantelts, ab.Elts) {
		t.Fatalf("not equal %v", pretty.Diff(wantelts, ab.Elts))
	}
}

func TestBaseHasAndGetEltIndex(t *testing.T) {

	cnt := int32(024)
	indexes := randIndexes(cnt)
	indexMap := makeIndexMap(indexes)

	maxIndex := indexes[len(indexes)-1]

	arr := &array.Base{}
	err := arr.InitIndex(indexes)
	if err != nil {
		t.Fatalf("expect no err but: %s", err)
	}

	for i := int32(0); i < maxIndex+128; i++ {

		_, inMap := indexMap[i]

		has := arr.Has(i)
		if inMap != has {
			t.Fatalf("Has(%d) expect: %v; but: %v", i, inMap, has)
		}

		eltIndex, inElts := arr.GetEltIndex(i)
		if inMap != inElts {
			t.Fatalf("Has(%d) expect: %v; but: %v", i, inMap, inElts)
		}

		if inMap && indexes[eltIndex] != i {
			t.Fatalf("i=%d should be at %d", i, eltIndex)
		}
	}
}

func TestBaseGet(t *testing.T) {

	indexes := []int32{1, 3, 100}
	elts := []uint16{1, 3, 100}

	ab := &array.Base{}
	ab.EltEncoder = encode.U16{}

	// test empty array
	v, found := ab.Get(1)
	if found {
		t.Fatalf("expect not found but: %v", found)
	}
	if v != nil {
		t.Fatalf("expect nil but: %v", v)
	}

	err := ab.Init(indexes, elts)
	if err != nil {
		t.Fatalf("expected no error but: %#v", err)
	}

	for i, idx := range indexes {
		v, found := ab.Get(idx)
		if !found {
			t.Fatalf("%d-th, expect found but: %d %v", i+1, idx, found)
		}
		if v != elts[i] {
			t.Fatalf("%d-th, %d expect %d found but: %v", i+1, idx, elts[i], v)
		}

		v, found = ab.Get(idx + 1)
		if found {
			t.Fatalf("%d-th, expect not found but: %d %v", i+1, idx, found)
		}
		if v != nil {
			t.Fatalf("%d-th, %d expect nil found but: %v", i+1, idx, v)
		}
	}
}

func BenchmarkBaseHasAndGetEltIndex(b *testing.B) {

	var name string
	runs := []struct{ cnt int32 }{
		{5},
		{32},
		{256},
		{1024},
		{10240},
		{102400},
	}
	for _, r := range runs {

		indexes := randIndexes(r.cnt)
		maxIndex := indexes[len(indexes)-1]

		arr := &array.Base{}
		err := arr.InitIndex(indexes)
		if err != nil {
			panic(err)
		}

		name = fmt.Sprintf("Has-%d", r.cnt)
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				arr.Has(int32(i) % maxIndex)
			}
		})

		name = fmt.Sprintf("GetEltIndex-%d", r.cnt)
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				arr.GetEltIndex(int32(i) % maxIndex)
			}
		})
	}
}

func BenchmarkBaseGet(b *testing.B) {

	indexes := []int32{1, 3, 100}
	elts := []uint16{1, 3, 100}

	ab := &array.Base{}
	err := ab.Init(indexes, elts)
	if err != nil {
		panic(err)
	}

	ab.EltEncoder = encode.U16{}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ab.Get(1)
	}
}

func BenchmarkBaseGetBytes(b *testing.B) {

	indexes := []int32{1, 3, 100}
	elts := []uint16{1, 3, 100}

	ab := &array.Base{}
	err := ab.Init(indexes, elts)
	if err != nil {
		panic(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ab.GetBytes(1, 2)
	}
}
