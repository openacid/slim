package trie

import (
	"reflect"
	"testing"
)

func TestTrie(t *testing.T) {

	var cases = []struct {
		key      [][]byte
		value    [][]byte
		expected [][][]byte
	}{
		{
			key: [][]byte{
				{'a', 'b', 'c'},
				{'a', 'b', 'd'},
				{'b', 'c', 'd'},
				{'b', 'c', 'e'},
				{'c', 'd', 'e'},
			},
			value: [][]byte{
				{0},
				{1},
				{2},
				{3},
				{4},
			},
			expected: [][][]byte{
				{{'a', 'b', 'c'}, {0}},
				{{'a', 'b', 'd'}, {1}},
				{{'b', 'c', 'd'}, {2}},
				{{'b', 'c', 'e'}, {3}},
				{{'c', 'd', 'e'}, {4}},
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
			value: [][]byte{
				{0},
				{1},
				{2},
				{3},
				{4},
			},
			expected: [][][]byte{
				{{'a', 'b', 'c'}, {0}},
				{{'a', 'b', 'c', 'd'}, {1}},
				{{'b', 'c'}, {2}},
				{{'b', 'c', 'd'}, {3}},
				{{'b', 'c', 'd', 'e'}, {4}},
			},
		},
	}

	for _, c := range cases {

		trie := New(c.key, c.value)
		for _, kv := range c.expected {
			ks := kv[0]
			val := kv[1]
			rst, err := trie.Search(ks)

			if err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(val, rst) {
				t.Error("ks: ", ks, "expected value: ", val, "rst: ", rst)
			}
		}

		trie.Squash()
		for _, kv := range c.expected {
			ks := kv[0]
			val := kv[1]
			rst, err := trie.Search(ks)

			if err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(val, rst) {
				t.Error("ks: ", ks, "expected value: ", val, "rst: ", rst)
			}
		}
	}
}

func TestTrieSearch(t *testing.T) {

	var key = [][]byte{
		{'a', 'b', 'c'},
		{'a', 'b', 'd'},
		{'b', 'c', 'd'},
		{'b', 'c', 'e'},
		{'c', 'd', 'e'},
	}
	var value = [][]byte{
		{0},
		{1},
		{2},
		{3},
		{4},
	}

	var trie = New(key, value)

	var cases = []struct {
		key      []byte
		expected []byte
		err      error
	}{
		{
			[]byte{'a', 'b', 'c'},
			[]byte{0},
			nil,
		},
		{
			[]byte{'a', 'b', 'd'},
			[]byte{1},
			nil,
		},
		{
			[]byte{'b', 'c', 'd'},
			[]byte{2},
			nil,
		},
		{
			[]byte{'b', 'c', 'e'},
			[]byte{3},
			nil,
		},
		{
			[]byte{'c', 'd', 'e'},
			[]byte{4},
			nil,
		},
		{
			[]byte{'a', 'c', 'b'},
			nil,
			KeyNotExist,
		},
		{
			[]byte{'a', 'b'},
			nil,
			KeyNotExist,
		},
		{
			[]byte{'a', 'b', 'c', 'd'},
			nil,
			KeyNotExist,
		},
		{
			[]byte{90, 'a', 'v'},
			nil,
			KeyNotExist,
		},
	}

	for _, c := range cases {
		rst, err := trie.Search(c.key)
		if !reflect.DeepEqual(err, c.err) {
			t.Error("err not equal.", "expected: ", c.err, "got: ", err)
		}

		if !reflect.DeepEqual(c.expected, rst) {
			t.Error("key: ", c.key, "expected value: ", c.expected, "rst: ", rst)
		}
	}

	var squashedCases = []struct {
		key      []byte
		expected []byte
		err      error
	}{
		{
			[]byte{'a', 'b', 'c'},
			[]byte{0},
			nil,
		},
		{
			[]byte{'a', 'd', 'c'},
			[]byte{0},
			nil,
		},
		{
			[]byte{'a', 'b', 'd'},
			[]byte{1},
			nil,
		},
		{
			[]byte{'a', 'c', 'd'},
			[]byte{1},
			nil,
		},
		{
			[]byte{'b', 'c', 'd'},
			[]byte{2},
			nil,
		},
		{
			[]byte{'b', 'e', 'd'},
			[]byte{2},
			nil,
		},
		{
			[]byte{'b', 'c', 'e'},
			[]byte{3},
			nil,
		},
		{
			[]byte{'b', 'd', 'e'},
			[]byte{3},
			nil,
		},
		{
			[]byte{'c', 'd', 'e'},
			[]byte{4},
			nil,
		},
		{
			[]byte{'c', 'f', 'e'},
			[]byte{4},
			nil,
		},
		{
			[]byte{'c'},
			nil,
			KeyNotExist,
		},
		{
			[]byte{'a', 'c'},
			nil,
			KeyNotExist,
		},
		{
			[]byte{'a', 'b', 'c', 'd'},
			nil,
			KeyNotExist,
		},
	}

	trie.Squash()
	for _, c := range squashedCases {
		rst, err := trie.Search(c.key)
		if !reflect.DeepEqual(err, c.err) {
			t.Error("err not equal.", "expected: ", c.err, "got: ", err)
		}

		if !reflect.DeepEqual(c.expected, rst) {
			t.Error("key: ", c.key, "expected value: ", c.expected, "rst: ", rst)
		}
	}

}
