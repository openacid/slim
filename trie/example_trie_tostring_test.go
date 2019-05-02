package trie

import "fmt"

func ExampleNode_String() {
	keys := [][]byte{
		{1, 2, 3},
		{1, 2, 3, 4},
		{2, 3},
		{2, 3, 4},
		{2, 3, 4, 5},
	}
	values := []int{0, 1, 2, 3, 4}
	trie, err := NewTrie(keys, values, false)
	if err != nil {
		panic(err)
	}

	fmt.Println(trie)

	// Output:
	// *2
	// -001->
	//       -002->
	//             -003->*2
	//                   -00$->=0
	//                   -004->
	//                         -00$->=1
	// -002->
	//       -003->*2
	//             -00$->=2
	//             -004->*2
	//                   -00$->=3
	//                   -005->
	//                         -00$->=4
}
