package trie_test

import (
	"fmt"
	"runtime"

	"github.com/openacid/slim/encode"
	"github.com/openacid/slim/trie"
)

var (
	m = encode.String16{}
)

func makeCopies(ncopy int, keys []string, offsets []uint16) []*trie.SlimTrie {
	copies := make([]*trie.SlimTrie, ncopy)

	for i := 0; i < len(copies); i++ {
		copies[i], _ = trie.NewSlimTrie(encode.U16{}, keys, offsets)
	}
	return copies
}

func Example_memoryUsage() {

	data := make([]byte, 0)

	keys := []string{}
	offsets := []uint16{}
	ksize := 0

	for i, k := range words2 {
		v := fmt.Sprintf("is the %d-th word", i)

		keys = append(keys, k)
		offsets = append(offsets, uint16(len(data)))

		data = append(data, m.Encode(k)...)
		data = append(data, m.Encode(v)...)

		ksize += len(k)
	}

	ncopy := 100
	vsize := 2.0

	rss1 := readRss()
	copies := makeCopies(ncopy, keys, offsets)
	rss2 := readRss()

	diff := float64(rss2 - rss1)

	avgIdxLen := diff/float64(ncopy)/float64(len(keys)) - vsize
	avgKeyLen := float64(ksize) / float64(len(keys))

	ratio := avgIdxLen / avgKeyLen * 100

	fmt.Printf(
		"Orignal:: %.1f byte/key --> SlimTrie index: %.1f byte/index\n"+
			"Saved %.1f%%",
		avgKeyLen,
		avgIdxLen,
		100-ratio,
	)

	for _, cc := range copies {
		_ = cc.Children
	}
}

func readRss() int64 {
	var stats runtime.MemStats

	runtime.GC()
	runtime.ReadMemStats(&stats)

	return int64(stats.Alloc)
}
