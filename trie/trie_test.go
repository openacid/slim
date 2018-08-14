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

	var srcs = []struct {
		key     []byte
		value   int
		isStart bool
	}{
		{
			[]byte("0000001118"),
			0,
			true,
		},
		{
			[]byte("0000001128"),
			1,
			false,
		},
		{
			[]byte("0000001138"),
			2,
			false,
		},
		{
			[]byte("0000001148"),
			3,
			false,
		},

		{
			[]byte("0000001158"),
			10,
			true,
		},
		{
			[]byte("0000002128"),
			11,
			false,
		},
		{
			[]byte("0000003138"),
			12,
			false,
		},
		{
			[]byte("0000004148"),
			13,
			false,
		},

		{
			[]byte("0000004158"),
			20,
			true,
		},
		{
			[]byte("0000012128"),
			21,
			false,
		},
		{
			[]byte("0000023138"),
			22,
			false,
		},
		{
			[]byte("0000024148"),
			23,
			false,
		},
	}

	rt := NewRangeTrie()

	for _, s := range srcs {
		_, err := rt.AddKV(s.key, s.value, s.isStart, true)
		if err == ErrTooManyTrieNodes {
			fmt.Printf("warn: %v\n", err)
			break
		} else if err != nil {
			t.Fatalf("failed to new range trie. err: %v\n", err)
		}
	}

	rt.Squash()
	rt.RemoveUselessLeaves()

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
			ExpectType{nil, 0, 10},
		},
		{
			[]byte("0000001138"),
			ExpectType{0, nil, 10},
		},
		{
			[]byte("0000001158"),
			ExpectType{0, 10, 20},
		},
		{
			[]byte("0000002128"),
			ExpectType{10, nil, 20},
		},
		{
			[]byte("0000004148"),
			ExpectType{10, nil, 20},
		},
		{
			[]byte("0000004158"),
			ExpectType{10, 20, nil},
		},
		{
			[]byte("0000024148"),
			ExpectType{20, nil, nil},
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
