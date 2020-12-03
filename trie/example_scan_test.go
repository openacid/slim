package trie

import (
	"fmt"

	"github.com/openacid/slim/encode"
)

func ExampleSlimTrie_ScanFrom() {
	var keys = []string{
		"",
		"`",
		"a",
		"ab",
		"abc",
		"abca",
		"abcd",
		"abcd1",
		"abce",
		"be",
		"c",
		"cde0",
		"d",
	}
	values := makeI32s(len(keys))

	codec := encode.I32{}
	st, _ := NewSlimTrie(codec, keys, values, Opt{
		Complete: Bool(true),
	})

	// untilD stops when encountering "d".
	untilD := func(k, v []byte) bool {
		if string(k) == "d" {
			return false
		}

		_, i32 := codec.Decode(v)
		fmt.Println(string(k), i32)
		return true
	}

	fmt.Println("scan (ab, +∞):")
	st.ScanFrom("ab", false, true, untilD)

	fmt.Println()
	fmt.Println("scan [be, +∞):")
	st.ScanFrom("be", true, true, untilD)

	fmt.Println()
	fmt.Println("scan (ab, be):")
	st.ScanFromTo(
		"ab", false,
		"be", false,
		true, untilD)

	// Output:
	//
	// scan (ab, +∞):
	// abc 4
	// abca 5
	// abcd 6
	// abcd1 7
	// abce 8
	// be 9
	// c 10
	// cde0 11
	//
	// scan [be, +∞):
	// be 9
	// c 10
	// cde0 11
	//
	// scan (ab, be):
	// abc 4
	// abca 5
	// abcd 6
	// abcd1 7
	// abce 8
}
