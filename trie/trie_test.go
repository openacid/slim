package trie

import (
	"fmt"
	"reflect"
	"testing"
	"xec/xerrors"
)

func TestTrie(t *testing.T) {

	type ExpectType struct {
		key []byte
		rst int
	}

	var cases = []struct {
		key      [][]byte
		value    interface{}
		expected []ExpectType
	}{
		{
			key: [][]byte{
				{'a', 'b', 'c'},
				{'a', 'b', 'd'},
				{'b', 'c', 'd'},
				{'b', 'c', 'e'},
				{'c', 'd', 'e'},
			},
			value: []int{
				0,
				1,
				2,
				3,
				4,
			},
			expected: []ExpectType{
				ExpectType{[]byte{'a', 'b', 'c'}, 0},
				ExpectType{[]byte{'a', 'b', 'd'}, 1},
				ExpectType{[]byte{'b', 'c', 'd'}, 2},
				ExpectType{[]byte{'b', 'c', 'e'}, 3},
				ExpectType{[]byte{'c', 'd', 'e'}, 4},
			},
		},
		{
			key: [][]byte{
				{'a', 'b', 'c'},
				{'a', 'b', 'c', 'd'},
				{'b', 'c'},
				{'b', 'c', 'd'},
				{'b', 'c', 'd', 'e'},
			},
			value: []int{
				0,
				1,
				2,
				3,
				4,
			},
			expected: []ExpectType{
				ExpectType{[]byte{'a', 'b', 'c'}, 0},
				ExpectType{[]byte{'a', 'b', 'c', 'd'}, 1},
				ExpectType{[]byte{'b', 'c'}, 2},
				ExpectType{[]byte{'b', 'c', 'd'}, 3},
				ExpectType{[]byte{'b', 'c', 'd', 'e'}, 4},
			},
		},
	}

	for _, c := range cases {

		trie, _ := New(c.key, c.value)
		for _, ex := range c.expected {
			_, rst, _ := trie.Search(ex.key)

			if !reflect.DeepEqual(ex.rst, rst) {
				t.Error("ks: ", ex.key, "expected value: ", ex.rst, "rst: ", rst)
			}
		}

		trie.Squash()
		for _, ex := range c.expected {
			_, rst, _ := trie.Search(ex.key)

			if !reflect.DeepEqual(ex.rst, rst) {
				t.Error("ks: ", ex.key, "expected value: ", ex.rst, "rst: ", rst)
			}
		}

		trie.Squash()
		trie.Squash()
		for _, ex := range c.expected {
			_, rst, _ := trie.Search(ex.key)

			if !reflect.DeepEqual(ex.rst, rst) {
				t.Error("ks: ", ex.key, "expected value: ", ex.rst, "rst: ", rst)
			}
		}
	}
}

func TestTrieSearch(t *testing.T) {

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
	var value = []int{
		0,
		1,
		2,
		3,
		4,
		5,
		6,
		7,
	}

	var trie, _ = New(key, value)

	type ExpectType struct {
		ltVal interface{}
		eqVal interface{}
		gtVal interface{}
	}

	var cases = []struct {
		key      []byte
		expected ExpectType
	}{
		{
			[]byte{'a', 'b', 'c'},
			ExpectType{nil, 0, 1},
		},
		{
			[]byte{'a', 'b', 'd'},
			ExpectType{1, 2, 3},
		},
		{
			[]byte{'b', 'c', 'd'},
			ExpectType{4, 5, 6},
		},
		{
			[]byte{'b', 'c', 'e'},
			ExpectType{6, nil, 7},
		},
		{
			[]byte{'c', 'd', 'e'},
			ExpectType{6, 7, nil},
		},
		{
			[]byte{'a', 'c', 'b'},
			ExpectType{3, nil, 4},
		},
		{
			[]byte{90, 'a', 'v'},
			ExpectType{nil, nil, 0},
		},
		{
			[]byte{'a', 'b'},
			ExpectType{nil, nil, 0},
		},
		{
			[]byte{'a', 'c'},
			ExpectType{3, nil, 4},
		},
		{
			[]byte{90},
			ExpectType{nil, nil, 0},
		},
		{
			[]byte{'a', 'b', 'c', 'd', 'e'},
			ExpectType{1, nil, 2},
		},
	}

	for _, c := range cases {

		ltVal, eqVal, gtVal := trie.Search(c.key)

		if !reflect.DeepEqual(c.expected.ltVal, ltVal) {
			t.Error("key: ", c.key, "expected lt value: ", c.expected.ltVal, "rst: ", ltVal)
		}
		if !reflect.DeepEqual(c.expected.eqVal, eqVal) {
			t.Error("key: ", c.key, "expected eq value: ", c.expected.eqVal, "rst: ", eqVal)
		}
		if !reflect.DeepEqual(c.expected.gtVal, gtVal) {
			t.Error("key: ", c.key, "expected gt value: ", c.expected.gtVal, "rst: ", gtVal)
		}
	}

	var squashedCases = []struct {
		key      []byte
		expected ExpectType
	}{
		{
			[]byte{'a', 'b'},
			ExpectType{nil, nil, 0},
		},
		{
			[]byte{'a', 'b', 'c'},
			ExpectType{nil, 0, 1},
		},
		{
			[]byte{'a', 'd', 'c'},
			ExpectType{nil, 0, 1},
		},
		{
			[]byte{'a', 'b', 'd'},
			ExpectType{1, 2, 3},
		},
		{
			[]byte{'a', 'c', 'd'},
			ExpectType{1, 2, 3},
		},
		{
			[]byte{'c', 'd', 'e'},
			ExpectType{6, 7, nil},
		},
		{
			[]byte{'c', 'f', 'e'},
			ExpectType{6, 7, nil},
		},
		{
			[]byte{'c', 'f', 'f'},
			ExpectType{6, 7, nil},
		},
		{
			[]byte{'c'},
			ExpectType{6, nil, 7},
		},
		{
			[]byte{'a', 'c'},
			ExpectType{nil, nil, 0},
		},
		{
			[]byte{'a', 'b', 'c', 'd', 'e'},
			ExpectType{1, nil, 2},
		},
	}

	trie.Squash()
	for _, c := range squashedCases {

		ltVal, eqVal, gtVal := trie.Search(c.key)

		if !reflect.DeepEqual(c.expected.ltVal, ltVal) {
			t.Error("key: ", c.key, "expected lt value: ", c.expected.ltVal, "rst: ", ltVal)
		}
		if !reflect.DeepEqual(c.expected.eqVal, eqVal) {
			t.Error("key: ", c.key, "expected eq value: ", c.expected.eqVal, "rst: ", eqVal)
		}
		if !reflect.DeepEqual(c.expected.gtVal, gtVal) {
			t.Error("key: ", c.key, "expected gt value: ", c.expected.gtVal, "rst: ", gtVal)
		}
	}
}

func TestTrieNew(t *testing.T) {

	var cases = []struct {
		keys        [][]byte
		values      interface{}
		expectedErr error
	}{
		{
			keys: [][]byte{
				{'a', 'b', 'c'},
				{'b', 'c', 'd'},
				{'c', 'd', 'e'},
			},
			values:      []int{0, 1, 2},
			expectedErr: nil,
		},
		{
			keys: [][]byte{
				{'a', 'b', 'c'},
				{'a', 'b', 'c'},
				{'c', 'd', 'e'},
			},
			values:      []int{0, 1, 2},
			expectedErr: xerrors.New(ErrDuplicateKeys, fmt.Sprintf("key: abc")),
		},
		{
			keys: [][]byte{
				{'a', 'b', 'c'},
				{'b', 'c', 'd'},
				{'c', 'd', 'e'},
			},
			values: map[string]int{
				"abc": 0,
				"bcd": 1,
				"cde": 2,
			},
			expectedErr: ErrValuesNotSlice,
		},
	}

	for _, c := range cases {
		_, err := New(c.keys, c.values)
		if !reflect.DeepEqual(err, c.expectedErr) {
			t.Fatalf("new trie: expectedErr: %v, got: %v", c.expectedErr, err)
		}
	}
}

func TestRangeTrie(t *testing.T) {

	var starts = [][]byte{
		[]byte("0000001118"),
		[]byte("0000001128"),
		[]byte("0000001138"),
		[]byte("0000001148"),
	}

	var ends = [][]byte{
		[]byte("0000001127"),
		[]byte("0000001137"),
		[]byte("0000001146"),
		[]byte("0000001148"),
	}

	var values = []int{
		0,
		1,
		2,
		3,
	}

	rt, err := NewRangeTrie(starts, ends, values)
	if err != nil {
		t.Fatal("new range trie failed")
	}

	type ExpectType struct {
		ltVal interface{}
		eqVal interface{}
		gtVal interface{}
	}

	var cases = []struct {
		key      []byte
		expected ExpectType
	}{
		{
			[]byte("0000001118"),
			ExpectType{nil, 0, 1},
		},
		{
			[]byte("0000001120"),
			ExpectType{0, nil, 1},
		},
		{
			[]byte("0000001128"),
			ExpectType{0, 1, 2},
		},
		{
			[]byte("0000001147"),
			ExpectType{2, nil, 3},
		},
		{
			[]byte("0000001148"),
			ExpectType{2, 3, nil},
		},
	}

	rt.RemoveEndLeaves()
	for _, c := range cases {
		ltVal, eqVal, gtVal := rt.Search(c.key)

		if !reflect.DeepEqual(c.expected.ltVal, ltVal) {
			t.Error("key: ", c.key, "expected lt value: ", c.expected.ltVal, "rst: ", ltVal)
		}
		if !reflect.DeepEqual(c.expected.eqVal, eqVal) {
			t.Error("key: ", c.key, "expected eq value: ", c.expected.eqVal, "rst: ", eqVal)
		}
		if !reflect.DeepEqual(c.expected.gtVal, gtVal) {
			t.Error("key: ", c.key, "expected gt value: ", c.expected.gtVal, "rst: ", gtVal)
		}
	}
}

func TestRangeTrieSquash(t *testing.T) {

	var starts = [][]byte{
		[]byte{'a', 'b', 'c', 'd'},
		[]byte{'a', 'b', 'd', 'f'},
		[]byte{'a', 'b', 'e', 'h'},
		[]byte{'a', 'b', 'f', 'a'},
	}

	var ends = [][]byte{
		[]byte{'a', 'b', 'd', 'e'},
		[]byte{'a', 'b', 'e', 'f'},
		[]byte{'a', 'b', 'e', 'h'},
		[]byte{'a', 'b', 'f', 'j'},
	}

	var values = []int{
		0,
		1,
		2,
		3,
	}

	rt, _ := NewRangeTrie(starts, ends, values)
	rt.Squash()
	rt.RemoveEndLeaves()

	type ExpectType struct {
		ltVal interface{}
		eqVal interface{}
		gtVal interface{}
	}

	var cases = []struct {
		key      []byte
		expected ExpectType
	}{
		{
			[]byte{'a', 'b', 'c', 'd'},
			ExpectType{nil, 0, 1},
		},
		{
			[]byte{'a', 'b', 'd', 'e'},
			ExpectType{0, nil, 1},
		},
		{
			[]byte{'a', 'b', 'e', 'h'},
			ExpectType{1, 2, 3},
		},
		{
			[]byte{'a', 'b', 'f', 'h'},
			ExpectType{3, nil, nil},
		},
		{
			[]byte{'a', 'b', 'f', 'j'},
			ExpectType{3, nil, nil},
		},
	}

	for _, c := range cases {
		ltVal, eqVal, gtVal := rt.Search(c.key)

		if !reflect.DeepEqual(c.expected.ltVal, ltVal) {
			t.Error("key: ", c.key, "expected lt value: ", c.expected.ltVal, "rst: ", ltVal)
		}
		if !reflect.DeepEqual(c.expected.eqVal, eqVal) {
			t.Error("key: ", c.key, "expected eq value: ", c.expected.eqVal, "rst: ", eqVal)
		}
		if !reflect.DeepEqual(c.expected.gtVal, gtVal) {
			t.Error("key: ", c.key, "expected gt value: ", c.expected.gtVal, "rst: ", gtVal)
		}
	}
}
