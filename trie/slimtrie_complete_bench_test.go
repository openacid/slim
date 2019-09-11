package trie

import (
	"testing"

	"github.com/openacid/slim/encode"
)

var OutputCompleteGetID20kvl10 int32

func BenchmarkSlimTrie_Complete_GetID_20k_vlen10(b *testing.B) {

	keys := getKeys("20kvl10")
	values := makeI32s(len(keys))
	st, _ := NewSlimTrie(encode.I32{}, keys, values, Opt{Complete: true})

	var id int32

	b.ResetTimer()

	i := b.N
	for {
		for _, k := range keys {
			id += st.GetID(k)

			i--
			if i == 0 {
				OutputCompleteGetID20kvl10 = id
				return
			}
		}
	}
}
