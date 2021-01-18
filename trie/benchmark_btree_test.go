package trie

import (
	"fmt"
	"testing"

	"github.com/google/btree"
	"github.com/openacid/testkeys"
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

	for _, typ := range testkeys.AssetNames() {

		if typ == "empty" {
			continue
		}

		// if typ != "200kweb2" {
		// if typ != "50kvl10" {
		//     continue
		// }

		b.Run(fmt.Sprintf("%s", typ), func(b *testing.B) {

			keys := getKeys(typ)
			values := makeI32s(len(keys))
			bt := btree.New(32)
			elts := makeKVElts(keys, values)
			for _, v := range elts {
				bt.ReplaceOrInsert(v)
			}

			n := len(keys)
			mask := 1
			for ; (mask << 1) <= n; mask <<= 1 {
			}
			mask--

			b.ResetTimer()

			var id int32
			i := b.N
			for {
				for j, k := range keys {

					itm := &KVElt{Key: k, Val: values[j&mask]}
					ee := bt.Get(itm)

					id += ee.(*KVElt).Val

					i--
					if i == 0 {
						OutputBtree = int(id)
						return
					}
				}
			}
		})
	}
}
