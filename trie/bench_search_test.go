package trie_test

import (
	"fmt"
	"testing"

	"github.com/openacid/slim/trie/benchmark"
)

var runs = []benchmark.Config{
	{KeyCnt: 1, KeyLen: 1024, ValLen: 2},
	{KeyCnt: 10, KeyLen: 1024, ValLen: 2},
	{KeyCnt: 100, KeyLen: 1024, ValLen: 2},
	{KeyCnt: 1000, KeyLen: 1024, ValLen: 2},
	{KeyCnt: 1000, KeyLen: 512, ValLen: 2},
	{KeyCnt: 1000, KeyLen: 256, ValLen: 2},
}

func BenchmarkTrieSearch(b *testing.B) {

	for _, r := range runs {
		testSrc := benchmark.MakeTestSrc(r.KeyCnt, r.KeyLen, r.ValLen)

		tr := testSrc.Slim

		name := fmt.Sprintf("%d-keys-%d-length: slimstrie search existing", r.KeyCnt, r.KeyLen)
		b.Run(name, benchmark.MakeTrieBenchFunc(tr, testSrc.SearchKey))

		// trie search nonexistent key
		name = fmt.Sprintf("%d-keys-%d-length: slimstrie search nonexistent", r.KeyCnt, r.KeyLen)
		searchKey := fmt.Sprintf("%snot found", testSrc.SearchKey)
		b.Run(name, benchmark.MakeTrieBenchFunc(tr, searchKey))
	}
}

func BenchmarkMapSearch(b *testing.B) {

	for _, r := range runs {
		testSrc := benchmark.MakeTestSrc(r.KeyCnt, r.KeyLen, r.ValLen)

		m := testSrc.Map

		name := fmt.Sprintf("%d-keys-%d-length: map search existing", r.KeyCnt, r.KeyLen)
		b.Run(name, benchmark.MakeMapBenchFunc(m, testSrc.SearchKey))

		name = fmt.Sprintf("%d-keys-%d-length: map search nonexistent", r.KeyCnt, r.KeyLen)
		searchKey := fmt.Sprintf("%snot found", testSrc.SearchKey)
		b.Run(name, benchmark.MakeMapBenchFunc(m, searchKey))
	}
}

func BenchmarkArraySearch(b *testing.B) {

	for _, r := range runs {
		testSrc := benchmark.MakeTestSrc(r.KeyCnt, r.KeyLen, r.ValLen)

		keys := testSrc.Keys
		values := testSrc.Values

		name := fmt.Sprintf("%d-keys-%d-length: array search existing", r.KeyCnt, r.KeyLen)
		b.Run(name, benchmark.MakeArrayBenchFunc(keys, values, testSrc.SearchKey))

		name = fmt.Sprintf("%d-keys-%d-length: array search nonexistent", r.KeyCnt, r.KeyLen)
		searchKey := fmt.Sprintf("%snot found", testSrc.SearchKey)
		b.Run(name, benchmark.MakeArrayBenchFunc(keys, values, searchKey))
	}
}

func BenchmarkBTreeSearch(b *testing.B) {

	for _, r := range runs {
		testSrc := benchmark.MakeTestSrc(r.KeyCnt, r.KeyLen, r.ValLen)
		bt := testSrc.Btree

		name := fmt.Sprintf("%d-keys-%d-length: btree search existing", r.KeyCnt, r.KeyLen)
		searchItem := &benchmark.TrieBenchKV{Key: testSrc.SearchKey, Val: testSrc.SearchValue}
		b.Run(name, benchmark.MakeBTreeBenchFunc(bt, searchItem))

		name = fmt.Sprintf("%d-keys-%d-length: btree search nonexistent", r.KeyCnt, r.KeyLen)
		searchKey := fmt.Sprintf("%snot found", testSrc.SearchKey)
		searchItem = &benchmark.TrieBenchKV{Key: searchKey, Val: testSrc.SearchValue}
		b.Run(name, benchmark.MakeBTreeBenchFunc(bt, searchItem))
	}
}
