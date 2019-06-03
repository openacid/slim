package trie_test

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/openacid/slim/trie/benchmark"
)

var runs = []benchmark.Config{
	// {KeyCnt: 10, KeyLen: 8, ValLen: 2},

	{KeyCnt: 100, KeyLen: 32, ValLen: 2},
	{KeyCnt: 100, KeyLen: 64, ValLen: 2},
	{KeyCnt: 100, KeyLen: 128, ValLen: 2},
	{KeyCnt: 1000, KeyLen: 32, ValLen: 2},
	{KeyCnt: 1000, KeyLen: 64, ValLen: 2},
	{KeyCnt: 1000, KeyLen: 128, ValLen: 2},
	{KeyCnt: 10 * 1000, KeyLen: 32, ValLen: 2},
	{KeyCnt: 10 * 1000, KeyLen: 64, ValLen: 2},
	{KeyCnt: 10 * 1000, KeyLen: 128, ValLen: 2},
	{KeyCnt: 100 * 1000, KeyLen: 32, ValLen: 2},
	{KeyCnt: 100 * 1000, KeyLen: 64, ValLen: 2},
	{KeyCnt: 100 * 1000, KeyLen: 128, ValLen: 2},
	{KeyCnt: 1000 * 1000, KeyLen: 32, ValLen: 2},
	{KeyCnt: 1000 * 1000, KeyLen: 64, ValLen: 2},
	{KeyCnt: 1000 * 1000, KeyLen: 128, ValLen: 2},
}

var OutputBench int32 = 0

func Benchmark_Get_slim_btree_map_array(b *testing.B) {

	v := int32(0)
	for _, r := range runs {

		gst := benchmark.NewGetSetting(r.KeyCnt, r.KeyLen)
		mask := 1
		n := len(gst.Keys)
		for ; (mask << 1) <= n; mask <<= 1 {
		}

		nk := fmt.Sprintf("n=%d k=%d", r.KeyCnt, r.KeyLen)

		{
			name := "SlimTrie: " + nk

			b.Run(name+": present", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					v += gst.SlimKV.Get(gst.Keys[i&mask])
				}
			})

			b.Run(name+": absent", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					v += gst.SlimKV.Get(gst.AbsentKeys[i&mask])
				}
			})
		}

		{
			name := "Map: " + nk

			b.Run(name+": present", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					v += gst.Map[gst.Keys[i&mask]]
				}
			})

			b.Run(name+": absent", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					v += gst.Map[gst.AbsentKeys[i&mask]]
				}
			})
		}

		{
			name := "Btree: " + nk

			b.Run(name+": present", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					itm := &benchmark.KVElt{Key: gst.Keys[i&mask], Val: gst.Values[i&mask]}
					ee := gst.Btree.Get(itm)
					// if ee == nil {
					//     panic("not found")
					// }
					v += ee.(*benchmark.KVElt).Val
				}
			})

			b.Run(name+": absent", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					itm := &benchmark.KVElt{Key: gst.AbsentKeys[i&mask], Val: -1}
					ee := gst.Btree.Get(itm)
					if ee == nil {
						v++
					}
				}
			})
		}

		{
			name := "Array: " + nk

			b.Run(name+": present", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					v += sortedArraySearch(gst.Keys, gst.Values, gst.Keys[i&mask])
				}
			})

			b.Run(name+": absent", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					v += sortedArraySearch(gst.Keys, gst.Values, gst.AbsentKeys[i&mask])
				}
			})
		}
	}

	OutputBench += v
}

func sortedArraySearch(keys []string, values []int32, searchKey string) int32 {

	n := len(keys)

	idx := sort.Search(
		n,
		func(i int) bool {
			return strings.Compare(keys[i], searchKey) >= 0
		},
	)

	if idx < n && strings.Compare(keys[idx], searchKey) == 0 {
		return values[idx]
	}

	return -1
}
