package trie

import (
	"bytes"
	"encoding/binary"
	"os"
	"reflect"
	"testing"
	"xec/array"

	proto "github.com/golang/protobuf/proto"
)

var testDataFn = "data"

type CompactedExpectType struct {
	ltVal interface{}
	eqVal interface{}
	gtVal interface{}
}

/* this covertor is not secure,
just to make it easier to test it as a number of uint64
*/
type TestIntConv struct{}

func (c TestIntConv) MarshalElt(d interface{}) []byte {
	b := make([]byte, 8)
	v := uint64(d.(int))
	binary.LittleEndian.PutUint64(b, v)
	return b
}

func (c TestIntConv) UnmarshalElt(b []byte) (uint32, interface{}) {

	size := uint32(8)
	s := b[:size]

	d := binary.LittleEndian.Uint64(s)
	return size, int(d)
}

func (c TestIntConv) GetMarshaledEltSize(b []byte) uint32 {
	return 8
}

func wordKey(key []byte) []byte {
	w := make([]byte, len(key)*2)

	for i, k := range key {
		w[i*2] = byte((k & 0xf0) >> 4)
		w[i*2+1] = byte(k & 0x0f)
	}

	return w
}

func TestMaxKeys(t *testing.T) {

	nn := 16
	mx := 32768

	keys := make([][]byte, 0, mx)
	values := make([]interface{}, 0, mx)

	for i := 0; i < nn; i++ {
		for j := 0; j < nn; j++ {
			for k := 0; k < nn; k++ {
				for l := 0; l < 8; l++ {
					key := []byte{byte(i), byte(j), byte(k), byte(l)}
					keys = append(keys, key)

					value := i*nn*nn*nn + j*nn*nn + k*nn + l
					values = append(values, uint16(value))

				}
			}

		}
	}

	trie, err := New(keys, values)
	if err != nil {
		t.Fatalf("create new trie")
	}

	trie.Squash()

	ctrie := NewCompactedTrie(array.U16Conv{})
	err = ctrie.Compact(trie)
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	if ctrie.Children.Cnt != 1+16+256+4096 {
		t.Fatalf("children cnt should be %d", 1+16+256+4096)
	}
	if ctrie.Steps.Cnt != uint32(0) {
		t.Fatalf("Steps cnt should be %d", mx)
	}
	if ctrie.Leaves.Cnt != uint32(mx) {
		t.Fatalf("leaves cnt should be %d", mx)
	}
}

func TestMaxNode(t *testing.T) {

	mx := 32768

	keys := make([][]byte, 0, mx)
	values := make([]interface{}, 0, mx)

	for i := 0; i < mx; i++ {

		key := []byte{
			byte((i >> 14) & 0x01),
			byte((i >> 13) & 0x01),
			byte((i >> 12) & 0x01),
			byte((i >> 11) & 0x01),
			byte((i >> 10) & 0x01),
			byte((i >> 9) & 0x01),
			byte((i >> 8) & 0x01),
			byte((i >> 7) & 0x01),
			byte((i >> 6) & 0x01),
			byte((i >> 5) & 0x01),
			byte((i >> 4) & 0x01),
			byte((i >> 3) & 0x01),
			byte((i >> 2) & 0x01),
			byte((i >> 1) & 0x01),
			byte(i & 0x01),
		}

		keys = append(keys, key)
		values = append(values, uint16(i))
	}

	trie, err := New(keys, values)
	if err != nil {
		t.Fatalf("create new trie: %v", err)
	}

	trie.Squash()

	ctrie := NewCompactedTrie(array.U16Conv{})
	err = ctrie.Compact(trie)
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	if ctrie.Children.Cnt != uint32(mx-1) {
		t.Fatalf("children cnt should be %d", mx-1)
	}
	if ctrie.Steps.Cnt != uint32(0) {
		t.Fatalf("Steps cnt should be %d", mx)
	}
	if ctrie.Leaves.Cnt != uint32(mx) {
		t.Fatalf("leaves cnt should be %d", mx)
	}
}

func TestCompactedTrie(t *testing.T) {

	type ExpectKeyType struct {
		key []byte
		rst CompactedExpectType
	}

	var cases = []struct {
		key      [][]byte
		value    []interface{}
		expected []ExpectKeyType
	}{
		{
			key: [][]byte{
				{1, 2, 3},
				{1, 2, 4},
				{2, 3, 4},
				{2, 3, 5},
				{3, 4, 5},
			},
			value: []interface{}{
				0,
				1,
				2,
				3,
				4,
			},
			expected: []ExpectKeyType{
				ExpectKeyType{[]byte{1, 2, 3}, CompactedExpectType{nil, 0, 1}},
				ExpectKeyType{[]byte{1, 2, 4}, CompactedExpectType{0, 1, 2}},
				ExpectKeyType{[]byte{2, 3, 4}, CompactedExpectType{1, 2, 3}},
				ExpectKeyType{[]byte{2, 3, 5}, CompactedExpectType{2, 3, 4}},
				ExpectKeyType{[]byte{3, 4, 5}, CompactedExpectType{3, 4, nil}},
			},
		},
		{
			key: [][]byte{
				{1, 2, 3},
				{1, 2, 3, 4},
				{2, 3},
				{2, 3, 0},
				{2, 3, 4},
				{2, 3, 4, 5},
				{2, 3, 15},
			},
			value: []interface{}{
				0,
				1,
				2,
				3,
				4,
				5,
				6,
			},
			expected: []ExpectKeyType{
				ExpectKeyType{[]byte{1, 2, 3}, CompactedExpectType{nil, 0, 1}},
				ExpectKeyType{[]byte{1, 2, 3, 4}, CompactedExpectType{0, 1, 2}},
				ExpectKeyType{[]byte{2, 3}, CompactedExpectType{1, 2, 3}},
				ExpectKeyType{[]byte{2, 3, 0}, CompactedExpectType{2, 3, 4}},
				ExpectKeyType{[]byte{2, 3, 4}, CompactedExpectType{3, 4, 5}},
				ExpectKeyType{[]byte{2, 3, 4, 5}, CompactedExpectType{4, 5, 6}},
				ExpectKeyType{[]byte{2, 3, 15}, CompactedExpectType{5, 6, nil}},
			},
		},
	}

	for _, c := range cases {

		for i, k := range c.key {
			c.key[i] = wordKey(k)
		}

		trie, _ := New(c.key, c.value)
		for _, ex := range c.expected {
			lt, eq, gt := trie.Search(wordKey(ex.key))
			rst := CompactedExpectType{lt, eq, gt}

			if !reflect.DeepEqual(ex.rst, rst) {
				t.Fatal("key: ", wordKey(ex.key), "expected value: ", ex.rst, "rst: ", rst)
			}
		}

		ctrie := NewCompactedTrie(TestIntConv{})
		err := ctrie.Compact(trie)
		if err != nil {
			t.Error("compact trie error:", err)
		}

		for _, ex := range c.expected {
			lt, eq, gt := ctrie.Search(wordKey(ex.key))
			rst := CompactedExpectType{lt, eq, gt}

			if !reflect.DeepEqual(ex.rst, rst) {
				t.Fatal("key: ", wordKey(ex.key), "expected value: ", ex.rst, "rst: ", rst)
			}
		}

		trie.Squash()
		trie.Squash()
		for _, ex := range c.expected {
			lt, eq, gt := trie.Search(wordKey(ex.key))
			rst := CompactedExpectType{lt, eq, gt}

			if !reflect.DeepEqual(ex.rst, rst) {
				t.Fatal("key: ", ex.key, "expected value: ", ex.rst, "rst: ", rst)
			}
		}

		ctrie = NewCompactedTrie(TestIntConv{})
		err = ctrie.Compact(trie)
		if err != nil {
			t.Error("compact trie error:", err)
		}

		for _, ex := range c.expected {
			lt, eq, gt := ctrie.Search(wordKey(ex.key))
			rst := CompactedExpectType{lt, eq, gt}

			if !reflect.DeepEqual(ex.rst, rst) {
				t.Fatal("key: ", ex.key, "expected value: ", ex.rst, "rst: ", rst)
			}
		}

		trie.Squash()
		trie.Squash()
		for _, ex := range c.expected {
			lt, eq, gt := trie.Search(wordKey(ex.key))
			rst := CompactedExpectType{lt, eq, gt}

			if !reflect.DeepEqual(ex.rst, rst) {
				t.Fatal("key: ", ex.key, "expected value: ", ex.rst, "rst: ", rst)
			}
		}

		ctrie = NewCompactedTrie(TestIntConv{})
		err = ctrie.Compact(trie)
		if err != nil {
			t.Error("compact trie error:", err)
		}

		for _, ex := range c.expected {
			lt, eq, gt := ctrie.Search(wordKey(ex.key))
			rst := CompactedExpectType{lt, eq, gt}

			if !reflect.DeepEqual(ex.rst, rst) {
				t.Fatal("key: ", ex.key, "expected value: ", ex.rst, "rst: ", rst)
			}
		}
	}
}

func TestCompactedTrieSearch(t *testing.T) {

	var key = [][]byte{
		{'a', 'b', 'c'},
		{'a', 'b', 'c', 'd'},
		{'a', 'b', 'd'},
		{'a', 'b', 'd', 'e'},
		{'b', 'c'},
		{'b', 'c', 'd'},
		{'b', 'c', 'd', 'e'},
		{'c', 'd', 'e'},
	}
	var value = []interface{}{
		0,
		1,
		2,
		3,
		4,
		5,
		6,
		7,
	}

	for i, k := range key {
		key[i] = wordKey(k)
	}

	var trie, _ = New(key, value)

	ctrie := NewCompactedTrie(TestIntConv{})
	err := ctrie.Compact(trie)
	if err != nil {
		t.Error("compact trie error:", err)
	}

	var cases = []struct {
		key      []byte
		expected CompactedExpectType
	}{
		{
			[]byte{'a', 'b', 'c'},
			CompactedExpectType{nil, 0, 1},
		},
		{
			[]byte{'a', 'b', 'd'},
			CompactedExpectType{1, 2, 3},
		},
		{
			[]byte{'b', 'c', 'd'},
			CompactedExpectType{4, 5, 6},
		},
		{
			[]byte{'b', 'c', 'e'},
			CompactedExpectType{6, nil, 7},
		},
		{
			[]byte{'c', 'd', 'e'},
			CompactedExpectType{6, 7, nil},
		},
		{
			[]byte{'a', 'c', 'b'},
			CompactedExpectType{3, nil, 4},
		},
		{
			[]byte{'a', 'b'},
			CompactedExpectType{nil, nil, 0},
		},
		{
			[]byte{'a', 'c'},
			CompactedExpectType{3, nil, 4},
		},
		{
			[]byte{'a', 'b', 'c', 'd', 'e'},
			CompactedExpectType{1, nil, 2},
		},
		{
			[]byte{'a', 'b', 'c'},
			CompactedExpectType{nil, 0, 1},
		},
	}

	for _, c := range cases {

		kk := wordKey(c.key)
		lt, eq, gt := ctrie.Search(kk)
		rst := CompactedExpectType{lt, eq, gt}
		if !reflect.DeepEqual(c.expected, rst) {
			t.Fatal("key: ", kk, "expected value: ", c.expected, "rst: ", rst)
		}
	}

	var squashedCases = []struct {
		key      []byte
		expected CompactedExpectType
	}{
		{
			[]byte{'a', 'b'},
			CompactedExpectType{nil, nil, 0},
		},
		{
			[]byte{'a', 'b', 'c'},
			CompactedExpectType{nil, 0, 1},
		},
		{
			[]byte{'a', 'd', 'c'},
			CompactedExpectType{nil, 0, 1},
		},
		{
			[]byte{'a', 'b', 'd'},
			CompactedExpectType{1, 2, 3},
		},
		{
			[]byte{'a', 'c', 'd'},
			CompactedExpectType{1, 2, 3},
		},
		{
			[]byte{'c', 'd', 'e'},
			CompactedExpectType{6, 7, nil},
		},
		{
			[]byte{'c', 'f', 'e'},
			CompactedExpectType{6, 7, nil},
		},
		{
			[]byte{'c', 'f', 'f'},
			CompactedExpectType{6, 7, nil},
		},
		{
			[]byte{'c'},
			CompactedExpectType{6, nil, 7},
		},
		{
			[]byte{'a', 'c'},
			CompactedExpectType{nil, nil, 0},
		},
		{
			[]byte{'a', 'b', 'c', 'd', 'e'},
			CompactedExpectType{1, nil, 2},
		},
	}

	trie.Squash()

	ctrie = NewCompactedTrie(TestIntConv{})
	err = ctrie.Compact(trie)
	if err != nil {
		t.Error("compact trie error:", err)
	}

	for _, c := range squashedCases {

		kk := wordKey(c.key)
		lt, eq, gt := ctrie.Search(kk)
		rst := CompactedExpectType{lt, eq, gt}
		if !reflect.DeepEqual(c.expected, rst) {
			t.Fatal("key: ", kk, "expected value: ", c.expected, "rst: ", rst)
		}
	}
}

func TestCompactedTrieMarshalUnmarshal(t *testing.T) {
	key := [][]byte{
		{1, 2, 3},
		{1, 2, 4},
		{2, 3, 4},
		{2, 3, 5},
		{3, 4, 5},
	}
	value := []interface{}{
		uint16(0),
		uint16(1),
		uint16(2),
		uint16(3),
		uint16(4),
	}

	trie, _ := New(key, value)

	ctrie := NewCompactedTrie(array.U16Conv{})
	err := ctrie.Compact(trie)
	if err != nil {
		t.Fatalf("compact trie error: %v", err)
	}

	rw := new(bytes.Buffer)

	size := ctrie.GetMarshalSize()

	n, err := ctrie.Marshal(rw)
	if err != nil {
		t.Fatalf("failed to marshal ctrie: %v", err)
	}

	if n != size || int64(rw.Len()) != size {
		t.Fatalf("wrong marshal size: %d, %d, %d", n, size, rw.Len())
	}

	// unmarshal
	rCtrie := NewCompactedTrie(array.U16Conv{})
	err = rCtrie.Unmarshal(rw)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	// check
	checkCompactedTrie(ctrie, rCtrie, t)
}

func TestCompactedTrieMarshalAtUnmarshalAt(t *testing.T) {
	key := [][]byte{
		{1, 2, 3},
		{1, 2, 4},
		{2, 3, 4},
		{2, 3, 5},
		{3, 4, 5},
	}

	value1 := []interface{}{
		uint16(0),
		uint16(1),
		uint16(2),
		uint16(3),
		uint16(4),
	}

	value2 := []interface{}{
		uint16(10),
		uint16(11),
		uint16(12),
		uint16(13),
		uint16(14),
	}

	trie1, _ := New(key, value1)
	trie2, _ := New(key, value2)

	ctrie1 := NewCompactedTrie(array.U16Conv{})
	err := ctrie1.Compact(trie1)
	if err != nil {
		t.Fatalf("compact trie error: %v", err)
	}

	ctrie2 := NewCompactedTrie(array.U16Conv{})
	err = ctrie2.Compact(trie2)
	if err != nil {
		t.Fatalf("compact trie error: %v", err)
	}

	// marshalat
	wOFlags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	writer, err := os.OpenFile(testDataFn, wOFlags, 0755)
	if err != nil {
		t.Fatalf("failed to create file: %s, %v", testDataFn, err)
	}
	defer os.Remove(testDataFn)

	offset1 := int64(0)
	n, err := ctrie1.MarshalAt(writer, offset1)
	if err != nil {
		t.Fatalf("failed to marshal ctrie: %v", err)
	}

	size := ctrie1.GetMarshalSize()
	if n != size {
		t.Fatalf("wrong marshal size: %d, %d", n, size)
	}

	offset2 := offset1 + n
	n, err = ctrie2.MarshalAt(writer, offset2)
	if err != nil {
		t.Fatalf("failed to marshal ctrie: %v", err)
	}
	size = ctrie1.GetMarshalSize()
	if n != size {
		t.Fatalf("wrong marshal size: %d, %d", n, size)
	}

	writer.Close()

	// unmarshalat
	reader, err := os.OpenFile(testDataFn, os.O_RDONLY, 0755)
	if err != nil {
		t.Fatalf("failed to read file: %s, %v", testDataFn, err)
	}
	defer reader.Close()

	rCtrie1 := NewCompactedTrie(array.U16Conv{})
	_, err = rCtrie1.UnmarshalAt(reader, offset1)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	checkCompactedTrie(ctrie1, rCtrie1, t)

	rCtrie2 := NewCompactedTrie(array.U16Conv{})
	_, err = rCtrie2.UnmarshalAt(reader, offset2)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	checkCompactedTrie(ctrie2, rCtrie2, t)
}

func checkCompactedTrie(ctrie, rCtrie *CompactedTrie, t *testing.T) {
	if !proto.Equal(&(ctrie.Children.CompactedArray), &(rCtrie.Children.CompactedArray)) {
		t.Fatalf("Children not the same")
	}

	if !proto.Equal(&(ctrie.Steps.CompactedArray), &(rCtrie.Steps.CompactedArray)) {
		t.Fatalf("Step not the same")
	}

	if !proto.Equal(&ctrie.Leaves, &rCtrie.Leaves) {
		t.Fatalf("Leaves not the same")
	}
}
