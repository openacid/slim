package trie

import (
	"fmt"
	"testing"
)

func BenchmarkTrieSearch(b *testing.B) {

	for _, r := range runs {
		testSrc := makeTestSrc(r.keyCnt, r.keyLen, r.valLen)

		tr := testSrc.root

		name := fmt.Sprintf("%d-keys-%d-length: slimstrie search existing", r.keyCnt, r.keyLen)
		b.Run(name, makeTrieBenchFunc(tr, testSrc.searchKey))

		// trie search nonexistent key
		name = fmt.Sprintf("%d-keys-%d-length: slimstrie search nonexistent", r.keyCnt, r.keyLen)
		searchKey := fmt.Sprintf("%snot found", testSrc.searchKey)
		b.Run(name, makeTrieBenchFunc(tr, searchKey))
	}
}

func BenchmarkMapSearch(b *testing.B) {

	for _, r := range runs {
		testSrc := makeTestSrc(r.keyCnt, r.keyLen, r.valLen)

		m := testSrc.m

		name := fmt.Sprintf("%d-keys-%d-length: map search existing", r.keyCnt, r.keyLen)
		b.Run(name, makeMapBenchFunc(m, testSrc.searchKey))

		name = fmt.Sprintf("%d-keys-%d-length: map search nonexistent", r.keyCnt, r.keyLen)
		searchKey := fmt.Sprintf("%snot found", testSrc.searchKey)
		b.Run(name, makeMapBenchFunc(m, searchKey))
	}
}

func BenchmarkArraySearch(b *testing.B) {

	for _, r := range runs {
		testSrc := makeTestSrc(r.keyCnt, r.keyLen, r.valLen)

		keys := testSrc.srcKeys
		values := testSrc.srcValues

		name := fmt.Sprintf("%d-keys-%d-length: array search existing", r.keyCnt, r.keyLen)
		b.Run(name, makeArrayBenchFunc(keys, values, testSrc.searchKey))

		name = fmt.Sprintf("%d-keys-%d-length: array search nonexistent", r.keyCnt, r.keyLen)
		searchKey := fmt.Sprintf("%snot found", testSrc.searchKey)
		b.Run(name, makeArrayBenchFunc(keys, values, searchKey))
	}
}

func BenchmarkBTreeSearch(b *testing.B) {

	for _, r := range runs {
		testSrc := makeTestSrc(r.keyCnt, r.keyLen, r.valLen)
		bt := testSrc.bt

		name := fmt.Sprintf("%d-keys-%d-length: btree search existing", r.keyCnt, r.keyLen)
		searchItem := &testKV{key: testSrc.searchKey, val: testSrc.searchValue}
		b.Run(name, makeBTreeBenchFunc(bt, searchItem))

		name = fmt.Sprintf("%d-keys-%d-length: btree search nonexistent", r.keyCnt, r.keyLen)
		searchKey := fmt.Sprintf("%snot found", testSrc.searchKey)
		searchItem = &testKV{key: searchKey, val: testSrc.searchValue}
		b.Run(name, makeBTreeBenchFunc(bt, searchItem))
	}
}
