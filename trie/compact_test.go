package trie

import (
	"bytes"
	"encoding/binary"
	"os"
	"reflect"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/openacid/errors"
	"github.com/openacid/slim/array"
	"github.com/openacid/slim/strhelper"
)

var testDataFn = "data"

type searchRst struct {
	ltVal interface{}
	eqVal interface{}
	gtVal interface{}
}

/* this covertor is not secure,
just to make it easier to test it as a number of uint64
*/
type TestIntConv struct{}

func (c TestIntConv) Marshal(d interface{}) []byte {
	b := make([]byte, 8)
	v := uint64(d.(int))
	binary.LittleEndian.PutUint64(b, v)
	return b
}

func (c TestIntConv) Unmarshal(b []byte) (int, interface{}) {

	size := 8
	s := b[:size]

	d := binary.LittleEndian.Uint64(s)
	return size, int(d)
}

func (c TestIntConv) GetMarshaledSize(b []byte) int {
	return 8
}

func wordKey(key []byte) []byte {
	w := make([]byte, len(key)*2)

	for i, k := range key {
		w[i*2] = (k & 0xf0) >> 4
		w[i*2+1] = k & 0x0f
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

	trie, err := NewTrie(keys, values)
	if err != nil {
		t.Fatalf("create new trie")
	}

	trie.Squash()

	ctrie, _ := NewSlimTrie(array.U16Conv{}, nil, nil)
	err = ctrie.LoadTrie(trie)
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

	trie, err := NewTrie(keys, values)
	if err != nil {
		t.Fatalf("create new trie: %v", err)
	}

	trie.Squash()

	ctrie, _ := NewSlimTrie(array.U16Conv{}, nil, nil)
	err = ctrie.LoadTrie(trie)
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	if ctrie.Children.Cnt != uint32(mx-1) {
		t.Fatalf("children cnt should be %d, but: %d", mx-1, ctrie.Children.Cnt)
	}
	if ctrie.Steps.Cnt != uint32(0) {
		t.Fatalf("Steps cnt should be %d", mx)
	}
	if ctrie.Leaves.Cnt != uint32(mx) {
		t.Fatalf("leaves cnt should be %d", mx)
	}
}

func TestSlimTrie(t *testing.T) {

	type ExpectKeyType struct {
		key []byte
		rst searchRst
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
				{[]byte{1, 2, 3}, searchRst{nil, 0, 1}},
				{[]byte{1, 2, 4}, searchRst{0, 1, 2}},
				{[]byte{2, 3, 4}, searchRst{1, 2, 3}},
				{[]byte{2, 3, 5}, searchRst{2, 3, 4}},
				{[]byte{3, 4, 5}, searchRst{3, 4, nil}},
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
				{[]byte{1, 2, 3}, searchRst{nil, 0, 1}},
				{[]byte{1, 2, 3, 4}, searchRst{0, 1, 2}},
				{[]byte{2, 3}, searchRst{1, 2, 3}},
				{[]byte{2, 3, 0}, searchRst{2, 3, 4}},
				{[]byte{2, 3, 4}, searchRst{3, 4, 5}},
				{[]byte{2, 3, 4, 5}, searchRst{4, 5, 6}},
				{[]byte{2, 3, 15}, searchRst{5, 6, nil}},
			},
		},
	}

	for _, c := range cases {

		for i, k := range c.key {
			c.key[i] = wordKey(k)
		}

		trie, _ := NewTrie(c.key, c.value)
		for _, ex := range c.expected {
			lt, eq, gt := trie.Search(wordKey(ex.key))
			rst := searchRst{lt, eq, gt}

			if !reflect.DeepEqual(ex.rst, rst) {
				t.Fatal("key: ", wordKey(ex.key), "expected value: ", ex.rst, "rst: ", rst)
			}
		}

		ctrie, _ := NewSlimTrie(TestIntConv{}, nil, nil)
		err := ctrie.LoadTrie(trie)
		if err != nil {
			t.Error("compact trie error:", err)
		}

		for _, ex := range c.expected {
			lt, eq, gt := ctrie.searchWords(wordKey(ex.key))
			rst := searchRst{lt, eq, gt}

			if !reflect.DeepEqual(ex.rst, rst) {
				t.Fatal("key: ", wordKey(ex.key), "expected value: ", ex.rst, "rst: ", rst)
			}
		}

		trie.Squash()
		trie.Squash()
		for _, ex := range c.expected {
			lt, eq, gt := trie.Search(wordKey(ex.key))
			rst := searchRst{lt, eq, gt}

			if !reflect.DeepEqual(ex.rst, rst) {
				t.Fatal("key: ", ex.key, "expected value: ", ex.rst, "rst: ", rst)
			}
		}

		ctrie, _ = NewSlimTrie(TestIntConv{}, nil, nil)
		err = ctrie.LoadTrie(trie)
		if err != nil {
			t.Error("compact trie error:", err)
		}

		for _, ex := range c.expected {
			lt, eq, gt := ctrie.searchWords(wordKey(ex.key))
			rst := searchRst{lt, eq, gt}

			if !reflect.DeepEqual(ex.rst, rst) {
				t.Fatal("key: ", ex.key, "expected value: ", ex.rst, "rst: ", rst)
			}
		}

		trie.Squash()
		trie.Squash()
		for _, ex := range c.expected {
			lt, eq, gt := trie.Search(wordKey(ex.key))
			rst := searchRst{lt, eq, gt}

			if !reflect.DeepEqual(ex.rst, rst) {
				t.Fatal("key: ", ex.key, "expected value: ", ex.rst, "rst: ", rst)
			}
		}

		ctrie, _ = NewSlimTrie(TestIntConv{}, nil, nil)
		err = ctrie.LoadTrie(trie)
		if err != nil {
			t.Error("compact trie error:", err)
		}

		for _, ex := range c.expected {
			lt, eq, gt := ctrie.searchWords(wordKey(ex.key))
			rst := searchRst{lt, eq, gt}

			if !reflect.DeepEqual(ex.rst, rst) {
				t.Fatal("key: ", ex.key, "expected value: ", ex.rst, "rst: ", rst)
			}
		}
	}
}

var (
	searchKeys = []string{
		"abc",
		"abcd",
		"abd",
		"abde",
		"bc",
		"bcd",
		"bcde",
		"cde",
	}
	searchValues = []int{0, 1, 2, 3, 4, 5, 6, 7}
)

func TestNewSlimTrieWithKVs(t *testing.T) {

	st, err := NewSlimTrie(TestIntConv{}, []string{"ab", "cd"}, []int{1, 2})
	if err != nil {
		t.Error("expect no error but:", err)
	}

	v := st.Get("ab")
	if v == nil {
		t.Fatalf("%q should be found", "ab")
	}

	if v.(int) != 1 {
		t.Fatalf("v should be 2, but: %v", v)
	}
}

func TestNewSlimTrie(t *testing.T) {

	ctrie, _ := NewSlimTrie(TestIntConv{}, nil, nil)
	err := ctrie.load(searchKeys, searchValues)
	if err != nil {
		t.Error("compact trie error:", err)
	}

	var cases = []struct {
		key      string
		expected searchRst
	}{
		{"ab", searchRst{nil, nil, 0}},
		{"abc", searchRst{nil, 0, 1}},
		{"abcde", searchRst{1, nil, 2}},
		{"abd", searchRst{1, 2, 3}},
		{"ac", searchRst{nil, nil, 0}},
		{"acb", searchRst{nil, nil, 0}},
		{"acd", searchRst{1, 2, 3}},
		{"adc", searchRst{nil, 0, 1}},
		{"bcd", searchRst{4, 5, 6}},
		{"bce", searchRst{4, 5, 6}},
		{"c", searchRst{6, nil, 7}},
		{"cde", searchRst{6, 7, nil}},
		{"cfe", searchRst{6, 7, nil}},
		{"cff", searchRst{6, 7, nil}},
	}

	for _, c := range cases {

		lt, eq, gt := ctrie.Search(c.key)
		rst := searchRst{lt, eq, gt}
		if !reflect.DeepEqual(c.expected, rst) {
			t.Fatal("key: ", c.key, "expected value: ", c.expected, "rst: ", rst)
		}

		kk := strhelper.ToBitWords(c.key, 4)
		lt, eq, gt = ctrie.searchWords(kk)
		rst = searchRst{lt, eq, gt}
		if !reflect.DeepEqual(c.expected, rst) {
			t.Fatal("key: ", kk, "expected value: ", c.expected, "rst: ", rst)
		}
	}
}

func TestSlimTrieLoad(t *testing.T) {

	cases := []struct {
		keys    []string
		vals    interface{}
		wanterr error
	}{
		{
			[]string{"a", "a"},
			[]int{1},
			ErrKVLenNotMatch,
		},
		{
			[]string{"a", "a"},
			[]int{1, 2},
			ErrDuplicateKeys,
		},
		{
			[]string{"ab", "a"},
			[]int{1, 2},
			ErrKeyOutOfOrder,
		},
		{
			[]string{"ab", "aa"},
			[]int{1, 2},
			ErrKeyOutOfOrder,
		},
		{
			[]string{"ab", "aaa"},
			[]int{1, 2},
			ErrKeyOutOfOrder,
		},
	}

	for i, c := range cases {
		st, err := NewSlimTrie(TestIntConv{}, c.keys, c.vals)

		if c.wanterr != errors.Cause(err) {
			t.Fatalf("%d-th: input: keys: %v; vals: %v; wanterr: %v; actual: %v",
				i+1, c.keys, c.vals, c.wanterr, err)
		}

		if err == nil && len(c.keys) > 0 {
			v := st.Get(c.keys[0])
			if v == nil {
				t.Fatalf("%d-th: should be found but not. key=%q",
					i+1, c.keys[0])
			}
		}
	}
}

func TestSlimTrieSearch(t *testing.T) {

	key := strhelper.SliceToBitWords(searchKeys, 4)

	var trie, _ = NewTrie(key, searchValues)

	ctrie, _ := NewSlimTrie(TestIntConv{}, nil, nil)
	err := ctrie.LoadTrie(trie)
	if err != nil {
		t.Error("compact trie error:", err)
	}

	var cases = []struct {
		key      string
		expected searchRst
	}{
		{"abc", searchRst{nil, 0, 1}},
		{"abd", searchRst{1, 2, 3}},
		{"bcd", searchRst{4, 5, 6}},
		{"bce", searchRst{6, nil, 7}},
		{"cde", searchRst{6, 7, nil}},
		{"acb", searchRst{3, nil, 4}},
		{"ab", searchRst{nil, nil, 0}},
		{"ac", searchRst{3, nil, 4}},
		{"abcde", searchRst{1, nil, 2}},
		{"abc", searchRst{nil, 0, 1}},
	}

	for _, c := range cases {

		kk := strhelper.ToBitWords(c.key, 4)
		lt, eq, gt := ctrie.searchWords(kk)
		rst := searchRst{lt, eq, gt}
		if !reflect.DeepEqual(c.expected, rst) {
			t.Fatal("key: ", kk, "expected value: ", c.expected, "rst: ", rst)
		}
	}
}

func TestSlimTrieMarshalUnmarshal(t *testing.T) {
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

	trie, _ := NewTrie(key, value)

	ctrie, _ := NewSlimTrie(array.U16Conv{}, nil, nil)
	err := ctrie.LoadTrie(trie)
	if err != nil {
		t.Fatalf("compact trie error: %v", err)
	}

	rw := new(bytes.Buffer)

	size := ctrie.getMarshalSize()

	n, err := ctrie.marshal(rw)
	if err != nil {
		t.Fatalf("failed to marshal ctrie: %v", err)
	}

	if n != size || int64(rw.Len()) != size {
		t.Fatalf("wrong marshal size: %d, %d, %d", n, size, rw.Len())
	}

	// unmarshal
	rCtrie, _ := NewSlimTrie(array.U16Conv{}, nil, nil)
	err = rCtrie.unmarshal(rw)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	// check
	checkSlimTrie(ctrie, rCtrie, t)
}

func TestSlimTrieMarshalAtUnmarshalAt(t *testing.T) {
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

	trie1, _ := NewTrie(key, value1)
	trie2, _ := NewTrie(key, value2)

	ctrie1, _ := NewSlimTrie(array.U16Conv{}, nil, nil)
	err := ctrie1.LoadTrie(trie1)
	if err != nil {
		t.Fatalf("compact trie error: %v", err)
	}

	ctrie2, _ := NewSlimTrie(array.U16Conv{}, nil, nil)
	err = ctrie2.LoadTrie(trie2)
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
	n, err := ctrie1.marshalAt(writer, offset1)
	if err != nil {
		t.Fatalf("failed to marshal ctrie: %v", err)
	}

	size := ctrie1.getMarshalSize()
	if n != size {
		t.Fatalf("wrong marshal size: %d, %d", n, size)
	}

	offset2 := offset1 + n
	n, err = ctrie2.marshalAt(writer, offset2)
	if err != nil {
		t.Fatalf("failed to marshal ctrie: %v", err)
	}
	size = ctrie1.getMarshalSize()
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

	rCtrie1, _ := NewSlimTrie(array.U16Conv{}, nil, nil)
	_, err = rCtrie1.unmarshalAt(reader, offset1)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	checkSlimTrie(ctrie1, rCtrie1, t)

	rCtrie2, _ := NewSlimTrie(array.U16Conv{}, nil, nil)
	_, err = rCtrie2.unmarshalAt(reader, offset2)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	checkSlimTrie(ctrie2, rCtrie2, t)
}

func checkSlimTrie(ctrie, rCtrie *SlimTrie, t *testing.T) {
	if !proto.Equal(&(ctrie.Children.Array32Storage), &(rCtrie.Children.Array32Storage)) {
		t.Fatalf("Children not the same")
	}

	if !proto.Equal(&(ctrie.Steps.Array32Storage), &(rCtrie.Steps.Array32Storage)) {
		t.Fatalf("Step not the same")
	}

	// TODO need to check non-Array32Storage fields, in future there is
	// user-defined underlaying data structure
	// if !proto.Equal(&ctrie.Leaves.Array32Storage, &rCtrie.Leaves.Array32Storage) {
	if !proto.Equal(&ctrie.Leaves, &rCtrie.Leaves) {
		t.Fatalf("Leaves not the same")
	}
}
