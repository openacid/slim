package trie

import (
	"fmt"

	"github.com/openacid/slim/encode"
)

func ExampleSlimTrie_RangeGet() {

	// To index a map of key range to value with SlimTrie is very simple:
	//
	// Gives a set of key the same value, and use RangeGet() instead of Get().
	// SlimTrie does not store branches for adjacent leaves with the same value.

	keys := []string{
		"abc",
		"abcd",

		"bc",

		"bcd",
		"bce",
	}
	values := []int{
		1, 1,
		2,
		3, 3,
	}
	st, err := NewSlimTrie(encode.Int{}, keys, values)
	if err != nil {
		panic(err)
	}

	cases := []struct {
		key string
		msg string
	}{
		{"ab", "out of range"},

		{"abc", "in range"},
		{"abc1", "FALSE POSITIVE"},
		{"abc2", "FALSE POSITIVE"},
		{"abcd", "in range"},

		{"abcde", "FALSE POSITIVE: a suffix of abcd"},

		{"acc", "FALSE POSITIVE"},

		{"bc", "in single key range [bc]"},
		{"bc1", "FALSE POSITIVE"},

		{"bcd1", "FALSE POSITIVE"},

		// {"def", "FALSE POSITIVE"},
	}

	for _, c := range cases {
		v, found := st.RangeGet(c.key)
		fmt.Printf("%-10s %-5v %-5t: %s\n", c.key, v, found, c.msg)
	}

	// Output:
	// ab         <nil> false: out of range
	// abc        1     true : in range
	// abc1       1     true : FALSE POSITIVE
	// abc2       1     true : FALSE POSITIVE
	// abcd       1     true : in range
	// abcde      1     true : FALSE POSITIVE: a suffix of abcd
	// acc        1     true : FALSE POSITIVE
	// bc         2     true : in single key range [bc]
	// bc1        2     true : FALSE POSITIVE
	// bcd1       3     true : FALSE POSITIVE
}
