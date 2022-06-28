package trie

import (
	fmt "fmt"

	"github.com/openacid/slim/encode"
)

func Example_stringValue() {
	keys := []string{"foo", "hello"}
	values := []string{"bar", "world"}

	st, err := NewSlimTrie(encode.String16{}, keys, values)
	if err != nil {
		panic(err)
	}

	fmt.Println(st.Get("hello"))
	fmt.Println(st.Get("foo"))
	fmt.Println(st.Get("unknown"))

	// Output:
	//
	// world true
	// bar true
	// <nil> false
}
