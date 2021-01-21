package trie

import (
	"testing"

	"github.com/openacid/slim/encode"
)

var OutputNewSlimTrie int

func BenchmarkNewSlimTrie(b *testing.B) {

	benchBigKeySet(b, func(b *testing.B, typ string, keys []string) {

		n := len(keys)
		values := makeI32s(len(keys))

		b.ResetTimer()
		var s int
		for i := 0; i < b.N/n; i++ {
			st, err := NewSlimTrie(encode.I32{}, keys, values)
			if err != nil {
				panic(err)
			}
			s += int(st.inner.NodeTypeBM.Words[0])
		}

		OutputNewSlimTrie = s
	})
}
