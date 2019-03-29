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
	"github.com/openacid/slim/marshal"
)

func makeRandArray(cnt int32) (idx int32, indexes []int32, keysMap map[int32]bool, ar *array.ArrayBase, err error) {
	arr := &array.ArrayBase{}

	indexes = []int32{}
	rnd := rand.New(rand.NewSource(time.Now().Unix()))

	keysMap = map[int32]bool{}
	num, idx := int32(0), int32(0)
	for ; num < cnt; idx++ {
		if rnd.Intn(2) == 1 {
			indexes = append(indexes, idx)
			num++
			keysMap[idx] = true
		}
	}

	err = arr.InitIndex(indexes)
	return idx, indexes, keysMap, arr, err
}

func TestArrayBaseInitIndex(t *testing.T) {

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
		ab := &array.ArrayBase{}
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

func TestArrayBaseInit(t *testing.T) {

	ab := &array.ArrayBase{}
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
			marshal.ErrNotFixedSize,
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
		ab := &array.ArrayBase{}
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

type IntMarshaler struct{}

func (m IntMarshaler) Marshal(v interface{}) []byte {
	return []byte{byte(v.(int))}
}
func (m IntMarshaler) Unmarshal(d []byte) (int, interface{}) {
	return 1, int(d[0])
}
func (m IntMarshaler) GetSize(v interface{}) int {
	return 1
}
func (m IntMarshaler) GetMarshaledSize(d []byte) int {
	return 1
}

func TestArrayBaseInitWithMarshaler(t *testing.T) {
	ab := &array.ArrayBase{}
	ab.EltMarshaler = IntMarshaler{}
	err := ab.Init([]int32{1, 2}, []int{3, 4})
	if err != nil {
		t.Fatalf("expected no error but: %#v", err)
	}

	wantelts := []byte{3, 4}
	if !reflect.DeepEqual(wantelts, ab.Elts) {
		t.Fatalf("not equal %v", pretty.Diff(wantelts, ab.Elts))
	}
}

func TestArrayBaseHasAndGetEltIndex(t *testing.T) {

	maxIndex, indexes, keysMap, arr, err := makeRandArray(1024)
	if err != nil {
		t.Fatalf("expect no err but: %s", err)
	}

	for i := int32(0); i < maxIndex+128; i++ {
		if _, ok := keysMap[i]; ok {
			if !arr.Has(i) {
				t.Fatalf("expect has but not: %d", i)
			}
			eltIndex, found := arr.GetEltIndex(i)
			if !found {
				t.Fatalf("should found but not: %d", i)
			}
			if indexes[eltIndex] != i {
				t.Fatalf("i=%d should be at %d", i, eltIndex)
			}
		} else {
			if arr.Has(i) {
				t.Fatalf("expect not has but has: %d", i)
			}
			_, found := arr.GetEltIndex(i)
			if found {
				t.Fatalf("should not found but found: %d", i)
			}
		}
	}
}

func TestArrayBaseGetTo(t *testing.T) {

	indexes := []int32{1, 3, 100}
	elts := []uint16{1, 3, 100}

	ab := &array.ArrayBase{}

	// test empty array
	v := uint16(0)
	found := ab.GetTo(1, &v)
	if found {
		t.Fatalf("expect not found but: %v", found)
	}

	err := ab.Init(indexes, elts)
	if err != nil {
		t.Fatalf("expected no error but: %#v", err)
	}

	testPanic(t, func() { ii := int(0); ab.GetTo(1, &ii) }, "int is not size fixed")

	for i, idx := range indexes {
		found := ab.GetTo(idx, &v)
		if !found {
			t.Fatalf("%d-th, expect found but: %d %v", i+1, idx, found)
		}
		if v != elts[i] {
			t.Fatalf("%d-th, %d expect %d found but: %v", i+1, idx, elts[i], v)
		}

		v = 0
		found = ab.GetTo(idx+1, &v)
		if found {
			t.Fatalf("%d-th, expect not found but: %d %v", i+1, idx, found)
		}
		if v != 0 {
			t.Fatalf("%d-th, %d expect nil found but: %v", i+1, idx, v)
		}
	}
}

func TestArrayBaseGet(t *testing.T) {

	indexes := []int32{1, 3, 100}
	elts := []uint16{1, 3, 100}

	ab := &array.ArrayBase{}
	ab.EltMarshaler = marshal.U16{}

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

func BenchmarkHasAndGetEltIndex(b *testing.B) {

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
		maxIndex, _, _, arr, _ := makeRandArray(r.cnt)

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

func BenchmarkArrayBaseGetTo(b *testing.B) {

	indexes := []int32{1, 3, 100}
	elts := []uint16{1, 3, 100}

	ab := &array.ArrayBase{}
	err := ab.Init(indexes, elts)
	if err != nil {
		panic(err)
	}

	v := uint16(0)
	for i := 0; i < b.N; i++ {
		ab.GetTo(1, &v)
	}
}

func BenchmarkArrayBaseGet(b *testing.B) {

	indexes := []int32{1, 3, 100}
	elts := []uint16{1, 3, 100}

	ab := &array.ArrayBase{}
	err := ab.Init(indexes, elts)
	if err != nil {
		panic(err)
	}

	ab.EltMarshaler = marshal.U16{}

	for i := 0; i < b.N; i++ {
		ab.Get(1)
	}
}

func BenchmarkArrayBaseGetBytes(b *testing.B) {

	indexes := []int32{1, 3, 100}
	elts := []uint16{1, 3, 100}

	ab := &array.ArrayBase{}
	err := ab.Init(indexes, elts)
	if err != nil {
		panic(err)
	}

	ab.EltMarshaler = marshal.U16{}

	for i := 0; i < b.N; i++ {
		ab.GetBytes(1, 2)
	}
}
