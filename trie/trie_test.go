package trie

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/openacid/errors"
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
				{[]byte{'a', 'b', 'c'}, 0},
				{[]byte{'a', 'b', 'd'}, 1},
				{[]byte{'b', 'c', 'd'}, 2},
				{[]byte{'b', 'c', 'e'}, 3},
				{[]byte{'c', 'd', 'e'}, 4},
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
				{[]byte{'a', 'b', 'c'}, 0},
				{[]byte{'a', 'b', 'c', 'd'}, 1},
				{[]byte{'b', 'c'}, 2},
				{[]byte{'b', 'c', 'd'}, 3},
				{[]byte{'b', 'c', 'd', 'e'}, 4},
			},
		},
	}

	for _, c := range cases {

		trie, _ := NewTrie(c.key, c.value)
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

	var trie, _ = NewTrie(key, value)

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
			expectedErr: errors.Wrapf(ErrDuplicateKeys, "key: abc"),
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
		_, err := NewTrie(c.keys, c.values)
		if errors.Cause(err) != errors.Cause(c.expectedErr) {
			t.Fatalf("new trie: expectedErr: %v, got: %v", c.expectedErr, err)
		}
	}
}

func TestNewError(t *testing.T) {

	cases := []struct {
		keys    [][]byte
		vals    interface{}
		wanterr error
	}{
		{nil, nil, nil},
		{[][]byte{}, nil, ErrValuesNotSlice},
		{nil, []int{}, nil},
		{nil, 1, nil},
		{[][]byte{{1}}, 1, ErrValuesNotSlice},
		{[][]byte{}, []int{1}, ErrKVLenNotMatch},
		{[][]byte{{1}}, []int{}, ErrKVLenNotMatch},
		{[][]byte{{1, 2}, {1}}, []int{1, 2}, ErrKeyOutOfOrder},
		{[][]byte{{1, 2}, {1, 1}}, []int{1, 2}, ErrKeyOutOfOrder},
		{[][]byte{{1, 2}, {1, 2}}, []int{1, 2}, ErrDuplicateKeys},
	}

	for i, c := range cases {
		_, err := NewTrie(c.keys, c.vals)
		if c.wanterr != errors.Cause(err) {
			t.Fatalf("%d-th: input: keys:%v; vals: %v; wanterr: %v; actual: %v",
				i+1, c.keys, c.vals, c.wanterr, err)
		}
	}
}

func TestAppend(t *testing.T) {

	tr, err := NewTrie([][]byte{{2, 3}, {2, 5}}, []int{1, 2})
	if err != nil {
		t.Fatalf("expect no error but: %v", err)
	}

	cases := []struct {
		keys    []byte
		wanterr error
	}{
		{[]byte{1}, ErrKeyOutOfOrder},
		{[]byte{1, 2}, ErrKeyOutOfOrder},
		{[]byte{2}, ErrKeyOutOfOrder},
		{[]byte{2, 2}, ErrKeyOutOfOrder},
		{[]byte{2, 3}, ErrDuplicateKeys},
		{[]byte{2, 4}, ErrKeyOutOfOrder},
		{[]byte{2, 5}, ErrDuplicateKeys},
		{[]byte{2, 6}, nil},
	}

	for i, c := range cases {
		_, err := tr.Append(c.keys, 1, false, false)
		if c.wanterr != errors.Cause(err) {
			t.Fatalf("%d-th: input: %v; want: %v; actual: %v",
				i+1, c.keys, c.wanterr, err)
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

	rt := newRangeTrie()

	for _, s := range srcs {
		_, err := rt.Append(s.key, s.value, s.isStart, true)
		if err == ErrTooManyTrieNodes {
			fmt.Printf("warn: %v\n", err)
			break
		} else if err != nil {
			t.Fatalf("failed to new range trie. err: %v\n", err)
		}
	}

	rt.Squash()
	rt.removeNonboundaryLeaves()

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

func TestToStrings(t *testing.T) {
	// TODO
	trie, err := NewTrie([][]byte{}, []int{})
	if err != nil {
		t.Fatalf("expect no err: %s", err)
	}

	fmt.Println(trie.toStrings(0))
}
