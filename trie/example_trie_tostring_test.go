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
	// 000(2):001(1):002(1):003(2):  $(0):
	//                             004(1):  $(0):
	//        002(1):003(2):  $(0):
	//                      004(2):  $(0):
	//                             005(1):  $(0):
}
