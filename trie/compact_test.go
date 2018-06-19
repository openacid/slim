package trie

import (
	"bytes"
	"reflect"
	"testing"
	"xec/sparse"

	proto "github.com/golang/protobuf/proto"
)

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
		t.Fatalf("create new trie: %v", err)
	}

	trie.Squash()

	ctrie := NewCompactedTrie(sparse.U16Conv{})
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

	ctrie := NewCompactedTrie(sparse.U16Conv{})
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

	type ExpectType struct {
		key  []byte
		rst  uint16
		mode Mode
	}

	var cases = []struct {
		key      [][]byte
		value    []interface{}
		expected []ExpectType
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
				uint16(0),
				uint16(1),
				uint16(2),
				uint16(3),
				uint16(4),
			},
			expected: []ExpectType{
				ExpectType{[]byte{1, 2, 3}, 0, EQ},
				ExpectType{[]byte{1, 2, 4}, 1, EQ},
				ExpectType{[]byte{2, 3, 4}, 2, EQ},
				ExpectType{[]byte{2, 3, 5}, 3, EQ},
				ExpectType{[]byte{3, 4, 5}, 4, EQ},
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
				uint16(0),
				uint16(1),
				uint16(2),
				uint16(3),
				uint16(4),
				uint16(5),
				uint16(6),
			},
			expected: []ExpectType{
				ExpectType{[]byte{1, 2, 3}, 0, EQ},
				ExpectType{[]byte{1, 2, 3, 4}, 1, EQ},
				ExpectType{[]byte{2, 3}, 2, EQ},
				ExpectType{[]byte{2, 3, 0}, 3, EQ},
				ExpectType{[]byte{2, 3, 4}, 4, EQ},
				ExpectType{[]byte{2, 3, 4, 5}, 5, EQ},
				ExpectType{[]byte{2, 3, 15}, 6, EQ},
				ExpectType{[]byte{2, 3, 4}, 3, LT},
				ExpectType{[]byte{2, 3, 0}, 2, LT},
				ExpectType{[]byte{2, 3, 0}, 3, LT | EQ},
				ExpectType{[]byte{2, 3, 4}, 5, GT},
				ExpectType{[]byte{2, 3, 6}, 6, GT},
			},
		},
	}

	for _, c := range cases {

		for i, k := range c.key {
			c.key[i] = wordKey(k)
		}

		trie, _ := New(c.key, c.value)

		ctrie := NewCompactedTrie(sparse.U16Conv{})
		err := ctrie.Compact(trie)
		if err != nil {
			t.Error("compact trie error:", err)
		}

		for _, ex := range c.expected {
			rst := ctrie.Search(wordKey(ex.key), ex.mode)

			if !reflect.DeepEqual(ex.rst, rst) {
				t.Fatal("ks: ", wordKey(ex.key), "expected value: ", ex.rst, "rst: ", rst, ex.mode)
			}
		}

		trie.Squash()
		trie.Squash()

		ctrie = NewCompactedTrie(sparse.U16Conv{})
		err = ctrie.Compact(trie)
		if err != nil {
			t.Error("compact trie error:", err)
		}

		for _, ex := range c.expected {
			rst := ctrie.Search(wordKey(ex.key), ex.mode)

			if !reflect.DeepEqual(ex.rst, rst) {
				t.Error("ks: ", ex.key, "expected value: ", ex.rst, "rst: ", rst)
			}
		}

		trie.Squash()
		trie.Squash()

		ctrie = NewCompactedTrie(sparse.U16Conv{})
		err = ctrie.Compact(trie)
		if err != nil {
			t.Error("compact trie error:", err)
		}

		for _, ex := range c.expected {
			rst := ctrie.Search(wordKey(ex.key), ex.mode)

			if !reflect.DeepEqual(ex.rst, rst) {
				t.Error("ks: ", ex.key, "expected value: ", ex.rst, "rst: ", rst)
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
		[]byte{0},
		[]byte{1},
		[]byte{2},
		[]byte{3},
		[]byte{4},
		[]byte{5},
		[]byte{6},
		[]byte{7},
	}

	type ExpectType struct {
		mode Mode
		rst  interface{}
	}

	for i, k := range key {
		key[i] = wordKey(k)
	}

	var trie, _ = New(key, value)

	ctrie := NewCompactedTrie(sparse.ByteConv{EltSize: uint32(1)})
	err := ctrie.Compact(trie)
	if err != nil {
		t.Error("compact trie error:", err)
	}

	var cases = []struct {
		key      []byte
		expected []ExpectType
	}{
		{
			[]byte{'a', 'b', 'c'},
			[]ExpectType{
				{EQ, []byte{0}},
				{LT | EQ, []byte{0}},
				{LT, nil},
				{GT | EQ, []byte{0}},
				{GT, []byte{1}},
			},
		},
		{
			[]byte{'a', 'b', 'd'},
			[]ExpectType{
				{EQ, []byte{2}},
				{LT | EQ, []byte{2}},
				{LT, []byte{1}},
				{GT | EQ, []byte{2}},
				{GT, []byte{3}},
			},
		},
		{
			[]byte{'b', 'c', 'd'},
			[]ExpectType{
				{EQ, []byte{5}},
				{LT | EQ, []byte{5}},
				{LT, []byte{4}},
				{GT | EQ, []byte{5}},
				{GT, []byte{6}},
			},
		},
		{
			[]byte{'b', 'c', 'e'},
			[]ExpectType{
				{EQ, nil},
				{LT | EQ, []byte{6}},
				{LT, []byte{6}},
				{GT | EQ, []byte{7}},
				{GT, []byte{7}},
			},
		},
		{
			[]byte{'c', 'd', 'e'},
			[]ExpectType{
				{EQ, []byte{7}},
				{LT | EQ, []byte{7}},
				{LT, []byte{6}},
				{GT | EQ, []byte{7}},
				{GT, nil},
			},
		},
		{
			[]byte{'a', 'c', 'b'},
			[]ExpectType{
				{EQ, nil},
				{LT | EQ, []byte{3}},
				{LT, []byte{3}},
				{GT | EQ, []byte{4}},
				{GT, []byte{4}},
			},
		},
		{
			[]byte{'a', 'b'},
			[]ExpectType{
				{EQ, nil},
				{LT | EQ, nil},
				{LT, nil},
				{GT | EQ, []byte{0}},
				{GT, []byte{0}},
			},
		},
		{
			[]byte{'a', 'c'},
			[]ExpectType{
				{EQ, nil},
				{LT | EQ, []byte{3}},
				{LT, []byte{3}},
				{GT | EQ, []byte{4}},
				{GT, []byte{4}},
			},
		},
		{
			[]byte{'a', 'b', 'c', 'd', 'e'},
			[]ExpectType{
				{EQ, nil},
				{LT | EQ, []byte{1}},
				{LT, []byte{1}},
				{GT | EQ, []byte{2}},
				{GT, []byte{2}},
			},
		},
	}

	for _, c := range cases {

		kk := wordKey(c.key)
		for _, ex := range c.expected {

			rst := ctrie.Search(kk, ex.mode)
			if !reflect.DeepEqual(ex.rst, rst) {
				t.Fatal("key: ", kk, "expected value: ", ex.rst, "rst: ", rst, "mode: ", ex.mode)
			}
		}
	}

	var squashedCases = []struct {
		key      []byte
		expected []ExpectType
	}{
		{
			[]byte{'a', 'b'},
			[]ExpectType{
				{EQ, nil},
				{LT | EQ, nil},
				{LT, nil},
				{GT | EQ, []byte{0}},
				{GT, []byte{0}},
			},
		},
		{
			[]byte{'a', 'b', 'c'},
			[]ExpectType{
				{EQ, []byte{0}},
				{LT | EQ, []byte{0}},
				{LT, nil},
				{GT | EQ, []byte{0}},
				{GT, []byte{1}},
			},
		},
		{
			[]byte{'a', 'd', 'c'},
			[]ExpectType{
				{EQ, []byte{0}},
				{LT | EQ, []byte{0}},
				{LT, nil},
				{GT | EQ, []byte{0}},
				{GT, []byte{1}},
			},
		},
		{
			[]byte{'a', 'b', 'd'},
			[]ExpectType{
				{EQ, []byte{2}},
				{LT | EQ, []byte{2}},
				{LT, []byte{1}},
				{GT | EQ, []byte{2}},
				{GT, []byte{3}},
			},
		},
		{
			[]byte{'a', 'c', 'd'},
			[]ExpectType{
				{EQ, []byte{2}},
				{LT | EQ, []byte{2}},
				{LT, []byte{1}},
				{GT | EQ, []byte{2}},
				{GT, []byte{3}},
			},
		},
		{
			[]byte{'c', 'd', 'e'},
			[]ExpectType{
				{EQ, []byte{7}},
				{LT | EQ, []byte{7}},
				{LT, []byte{6}},
				{GT | EQ, []byte{7}},
				{GT, nil},
			},
		},
		{
			[]byte{'c', 'f', 'e'},
			[]ExpectType{
				{EQ, []byte{7}},
				{LT | EQ, []byte{7}},
				{LT, []byte{6}},
				{GT | EQ, []byte{7}},
				{GT, nil},
			},
		},
		{
			[]byte{'c', 'f', 'f'},
			[]ExpectType{
				{EQ, []byte{7}},
				{LT | EQ, []byte{7}},
				{LT, []byte{6}},
				{GT | EQ, []byte{7}},
				{GT, nil},
			},
		},
		{
			[]byte{'c'},
			[]ExpectType{
				{EQ, nil},
				{LT | EQ, []byte{6}},
				{LT, []byte{6}},
				{GT | EQ, []byte{7}},
				{GT, []byte{7}},
			},
		},
		{
			[]byte{'a', 'c'},
			[]ExpectType{
				{EQ, nil},
				{LT | EQ, nil},
				{LT, nil},
				{GT | EQ, []byte{0}},
				{GT, []byte{0}},
			},
		},
		{
			[]byte{'a', 'b', 'c', 'd', 'e'},
			[]ExpectType{
				{EQ, nil},
				{LT | EQ, []byte{1}},
				{LT, []byte{1}},
				{GT | EQ, []byte{2}},
				{GT, []byte{2}},
			},
		},
	}

	trie.Squash()

	ctrie = NewCompactedTrie(sparse.ByteConv{EltSize: uint32(1)})
	err = ctrie.Compact(trie)
	if err != nil {
		t.Error("compact trie error:", err)
	}

	for _, c := range squashedCases {

		kk := wordKey(c.key)
		for _, ex := range c.expected {

			rst := ctrie.Search(kk, ex.mode)
			if !reflect.DeepEqual(ex.rst, rst) {
				t.Fatal("key: ", kk, "expected value: ", ex.rst, "rst: ", rst, "mode: ", ex.mode)
			}
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

	ctrie := NewCompactedTrie(sparse.U16Conv{})
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
	rCtrie := NewCompactedTrie(sparse.U16Conv{})
	err = rCtrie.Unmarshal(rw)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	// check
	if !proto.Equal(&ctrie.Children, &rCtrie.Children) {
		t.Fatalf("Children not the same")
	}

	if !proto.Equal(&ctrie.Steps, &rCtrie.Steps) {
		t.Fatalf("Step not the same")
	}

	if !proto.Equal(&ctrie.Leaves, &rCtrie.Leaves) {
		t.Fatalf("Leaves not the same")
	}
}
