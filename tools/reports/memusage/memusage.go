package main

import (
	"fmt"

	"github.com/openacid/slim/benchhelper"
	"github.com/openacid/slim/encode"
	"github.com/openacid/slim/trie"
)

var Output int
var _ = Output

func main() {
	compareTrieMapMemUse()
}

func compareTrieMapMemUse() {

	writeTableHeader()

	for _, n := range []int{1000, 2000, 5000} {
		for _, k := range []int{64, 256, 1024} {

			trieSize := getTrieMem(n, k)
			mapSize := getMapMem(n, k)
			// make key + value as a value in trie
			// kvTrieSize := getKVTrieMem(cnt, l)
			kvTrieSize := getKVTrieMem2(n, k)

			mapAvg := float64(mapSize) / float64(n)
			trieAvg := float64(trieSize) / float64(n)
			kvTrieAvg := float64(kvTrieSize) / float64(n)

			writeTableRow(n, k, 2, trieAvg, mapAvg, kvTrieAvg)
		}
	}
}

func writeTableHeader() {

	fmt.Printf("| %s | %s | %s | %s | %s | %s |\n",
		"Key Count", "Key Length", "Value Size", "Trie Size (Byte/key)", "Map Size (Byte/key)",
		"KV Trie Size (Byte/key)")

	fmt.Printf("| --- | --- | --- | --- | --- | --- |\n")
}

func writeTableRow(cnt, kLen, vLen int, trieAvg, mapAvg, kvTrieAvg float64) {

	fmt.Printf("| %5d | %5d | %5d | %6.1f | %6.1f | %6.1f |\n",
		cnt, kLen, vLen, trieAvg, mapAvg, kvTrieAvg)
}

func getTrieMem(keyCnt, keyLen int) int64 {

	memStart := benchhelper.Allocated()

	keys := benchhelper.RandSortedStrings(keyCnt, keyLen)
	vals := make([]uint16, keyCnt)

	t, err := trie.NewSlimTrie(encode.U16{}, keys, vals)
	if err != nil {
		panic(err)
	}

	keys = nil
	vals = nil

	memEnd := benchhelper.Allocated()

	size := memEnd - memStart

	_ = keys
	_ = vals

	// reference them or memory is freed
	_ = t.Children
	_ = t.Steps
	_ = t.Leaves

	return size
}

func getKVTrieMem2(keyCnt, keyLen int) int64 {
	// make key + value as a value in trie

	memStart := benchhelper.Allocated()

	keys := benchhelper.RandSortedStrings(keyCnt, keyLen)
	vals := make([]uint16, keyCnt)
	indexes := make([]uint32, keyCnt)
	for i := 0; i < len(keys); i++ {
		indexes[i] = uint32(i)
	}

	t, err := trie.NewSlimTrie(encode.U32{}, keys, indexes)
	if err != nil {
		panic(err)
	}

	memEnd := benchhelper.Allocated()

	size := memEnd - memStart

	var s int
	for _, k := range keys {
		i := t.Get(k)
		// in real world, need to compare keys[i] and k to ensure it is a positive match
		v := vals[i.(uint32)]
		s += int(v)
	}
	Output = len(keys) + len(vals) + len(indexes) + s

	// reference them or memory is freed
	_ = t.Children
	_ = t.Steps
	_ = t.Leaves

	return size
}

func getMapMem(keyCnt, keyLen int) int64 {

	memStart := benchhelper.Allocated()

	keys := benchhelper.RandSortedStrings(keyCnt, keyLen)
	vals := make([]uint16, keyCnt)

	m := make(map[string]uint16, len(keys))

	for i := 0; i < len(keys); i++ {
		m[keys[i]] = vals[i]
	}

	keys = nil
	vals = nil

	memEnd := benchhelper.Allocated()

	size := memEnd - memStart

	Output = len(m)

	_ = keys
	_ = vals

	return size
}
