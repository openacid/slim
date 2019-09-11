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
	// -1->
	//     -2->
	//         -3->*2
	//             -$->=0
	//             -4->
	//                 -$->=1
	// -2->
	//     -3->*2
	//         -$->=2
	//         -4->*2
	//             -$->=3
	//             -5->
	//                 -$->=4
}
