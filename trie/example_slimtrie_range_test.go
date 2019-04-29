package trie

import (
	"fmt"

	"github.com/openacid/slim/encode"
)

func ExampleSlimTrie_RangeGet() {

	// To index a map of key range to value with SlimTrie is very simple:
	//
	// Just give two adjacent keys the same value, then SlimTrie knows these
	// keys belong to a "range".
	// These two keys are left and right boundaries of a range, and are both
	// inclusive.
	//
	// In this example we:
	//
	//   map [abc, abcd] to 1
	//   map [bc, bc]    to 2 // this range has only one key in it.
	//   map [bcd, bce]  to 3
	//
	// With RangeGet() to get any key that "abc" <= key <= "abcd", such as
	// "abc1", "abc2"... should return "1".
	//
	// False Positive
	//
	// Just like Bloomfilter, SlimTrie does not contains full information of keys,
	// thus there could be a false positive return:
	// It returns some value and "true" but the key is not in there.

	keys := []string{
		"abc", "abcd",
		"bc",
		"bcd", "bce",
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
		{"ab", "smaller than any"},

		{"abc", "in range [abc, abcd]"},
		{"abc1", "in range [abc, abcd]"},
		{"abc2", "in range [abc, abcd]"},
		{"abcd", "in range [abc, abcd]"},

		{"abcde", "FALSE POSITIVE: a suffix of abcd"},

		{"acc", "FALSE POSITIVE: not in range [abc, abcd]"},

		{"bc", "in single key range [bc]"},
		{"bc1", "not in single key range [bc]"},

		{"bcd1", "in range [bcd, bce]"},

		{"def", "greater than any"},
	}

	for _, c := range cases {
		v, found := st.RangeGet(c.key)
		fmt.Printf("%-10s %-5v %-5t: %s\n", c.key, v, found, c.msg)
	}

	// Output:
	// ab         <nil> false: smaller than any
	// abc        1     true : in range [abc, abcd]
	// abc1       1     true : in range [abc, abcd]
	// abc2       1     true : in range [abc, abcd]
	// abcd       1     true : in range [abc, abcd]
	// abcde      1     true : FALSE POSITIVE: a suffix of abcd
	// acc        1     true : FALSE POSITIVE: not in range [abc, abcd]
	// bc         2     true : in single key range [bc]
	// bc1        <nil> false: not in single key range [bc]
	// bcd1       3     true : in range [bcd, bce]
	// def        <nil> false: greater than any
}
