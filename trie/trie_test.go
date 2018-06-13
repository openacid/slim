package trie

import (
	"reflect"
	"testing"
)

func TestTrie(t *testing.T) {

	type ExpectType struct {
		key []byte
		rst int
	}

	var cases = []struct {
		key      [][]byte
		value    []interface{}
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
			value: []interface{}{
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
			value: []interface{}{
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

		trie := New(c.key, c.value)
		for _, ex := range c.expected {
			rst := trie.Search(ex.key, EQ)

			if !reflect.DeepEqual(ex.rst, rst) {
				t.Error("ks: ", ex.key, "expected value: ", ex.rst, "rst: ", rst)
			}
		}

		trie.Squash()
		for _, ex := range c.expected {
			rst := trie.Search(ex.key, EQ)

			if !reflect.DeepEqual(ex.rst, rst) {
				t.Error("ks: ", ex.key, "expected value: ", ex.rst, "rst: ", rst)
			}
		}

		trie.Squash()
		trie.Squash()
		for _, ex := range c.expected {
			rst := trie.Search(ex.key, EQ)

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

	var trie = New(key, value)

	type ExpectType struct {
		mode Mode
		rst  interface{}
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
			[]byte{90, 'a', 'v'},
			[]ExpectType{
				{EQ, nil},
				{LT | EQ, nil},
				{LT, nil},
				{GT | EQ, []byte{0}},
				{GT, []byte{0}},
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
			[]byte{90},
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

	for _, c := range cases {
		for _, ex := range c.expected {

			rst := trie.Search(c.key, ex.mode)
			if !reflect.DeepEqual(ex.rst, rst) {
				t.Error("key: ", c.key, "expected value: ", ex.rst, "rst: ", rst, "mode: ", ex.mode)
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
	for _, c := range squashedCases {
		for _, ex := range c.expected {
			rst := trie.Search(c.key, ex.mode)
			if !reflect.DeepEqual(ex.rst, rst) {
				t.Error("key: ", c.key, "expected value: ", ex.rst, "rst: ", rst)
			}
		}
	}
}
