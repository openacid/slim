package trie

import (
	"testing"

	"github.com/openacid/slim/encode"
)

var Output int

func BenchmarkNewSlimTrie_300(b *testing.B) {
	keys := getKeys("300vl50")
	values := make([]uint32, len(keys))

	n := len(keys)

	for i := 0; i < n; i++ {
		values[i] = uint32(i)
	}

	b.ResetTimer()
	var s int
	for i := 0; i < b.N/n; i++ {
		st, err := NewSlimTrie(encode.U32{}, keys, values)
		if err != nil {
			panic(err)
		}
		s += int(st.nodes.NodeTypeBM.Words[0])
	}

	Output = s
}
