package trie_test

import (
	"testing"

	"github.com/openacid/slim/encode"
	"github.com/openacid/slim/trie"
)

var Output int

func BenchmarkNewSlimTrie(b *testing.B) {
	keys := words2
	values := make([]uint32, len(keys))
	for i := 0; i < len(keys); i++ {
		values[i] = uint32(i)
	}
	b.ResetTimer()
	var s int
	for i := 0; i < b.N; i++ {
		st, err := trie.NewSlimTrie(encode.U32{}, keys, values)
		if err != nil {
			panic(err)
		}
		s += int(st.Children.Cnt)
	}

	Output = s
}
