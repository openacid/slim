package trie

import (
	"testing"

	"github.com/google/btree"
	"github.com/openacid/low/mathext/zipf"
)

type KVElt struct {
	Key string
	Val int32
}

func (kv *KVElt) Less(than btree.Item) bool {
	o := than.(*KVElt)
	return kv.Key < o.Key
}

func makeKVElts(srcKeys []string, srcVals []int32) []*KVElt {
	elts := make([]*KVElt, len(srcKeys))
	for i, k := range srcKeys {
		elts[i] = &KVElt{Key: k, Val: srcVals[i]}
	}
	return elts
}

var OutputBtree int

func Benchmark_btree(b *testing.B) {

	benchBigKeySet(b, func(b *testing.B, typ string, keys []string) {

		values := makeI32s(len(keys))
		bt := btree.New(32)
		elts := makeKVElts(keys, values)
		for _, v := range elts {
			bt.ReplaceOrInsert(v)
		}

		// sz := size.Of(bt)
		// fmt.Println(sz/1024, sz/len(keys))

		accesses := zipf.Accesses(2, 1.5, len(keys), b.N, nil)

		b.ResetTimer()

		var id int32
		for i := 0; i < b.N; i++ {
			idx := accesses[i]
			itm := &KVElt{Key: keys[idx], Val: values[idx]}
			ee := bt.Get(itm)
			id += ee.(*KVElt).Val

		}
		OutputBtree = int(id)
	})
}
