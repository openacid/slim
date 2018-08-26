package trie

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"testing"
	"time"
	"xec/testutil"

	"github.com/google/btree"
)

type testSrcType struct {
	srcKeys   []string
	srcValues [][]byte

	root *CompactedTrie
	m    map[string][]byte
	bt   *btree.BTree

	searchKey   string
	searchValue []byte
}

var (
	cnt    = int64(1000)
	keyLen = uint32(1024)
	valLen = uint32(2)

	btreeDegree = 2

	testSrc *testSrcType
)

func BenchmarkTrieSearch(b *testing.B) {

	if testSrc == nil {
		testSrc = makeTestSrc(cnt, keyLen, valLen)
	}

	tr := testSrc.root
	//searchKey := splitStringTo4BitWords(testSrc.searchKey)
	searchKey := testSrc.searchKey

	var val []byte

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val = trieSearchTestKV(tr, searchKey)
		//_, _, _ = tr.Search(searchKey)
	}

	if !testutil.ByteSliceEqual(val, testSrc.searchValue) {
		b.Errorf("search value is wrong.")
	}
}

func BenchmarkTrieSearchNotFound(b *testing.B) {

	if testSrc == nil {
		testSrc = makeTestSrc(cnt, keyLen, valLen)
	}

	tr := testSrc.root
	searchKey := fmt.Sprintf("%snot found", testSrc.searchKey)

	var val []byte

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val = trieSearchTestKV(tr, searchKey)
		//_, _, _ = tr.Search(searchKey)
	}

	if val != nil {
		b.Errorf("search not exsisted value failed.")
	}
}

func BenchmarkMapSearch(b *testing.B) {

	if testSrc == nil {
		testSrc = makeTestSrc(cnt, keyLen, valLen)
	}

	m := testSrc.m
	searchKey := testSrc.searchKey

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m[searchKey]
	}
}

func BenchmarkArraySearch(b *testing.B) {

	if testSrc == nil {
		testSrc = makeTestSrc(cnt, keyLen, valLen)
	}

	keys := testSrc.srcKeys
	values := testSrc.srcValues

	searchKey := testSrc.searchKey

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sortedArraySearch(keys, values, searchKey)
	}
}

func BenchmarkArraySearchNotFound(b *testing.B) {

	if testSrc == nil {
		testSrc = makeTestSrc(cnt, keyLen, valLen)
	}

	keys := testSrc.srcKeys
	values := testSrc.srcValues

	searchKey := fmt.Sprintf("%snot found", testSrc.searchKey)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sortedArraySearch(keys, values, searchKey)
	}
}

func BenchmarkMapSearchNotFound(b *testing.B) {

	if testSrc == nil {
		testSrc = makeTestSrc(cnt, keyLen, valLen)
	}

	m := testSrc.m
	searchKey := fmt.Sprintf("%snot found", testSrc.searchKey)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m[searchKey]
	}
}

func BenchmarkBTreeSearch(b *testing.B) {

	if testSrc == nil {
		testSrc = makeTestSrc(cnt, keyLen, valLen)
	}

	bt := testSrc.bt
	searchItem := &testKV{key: testSrc.searchKey, val: testSrc.searchValue}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = bt.Get(searchItem)
	}
}

func BenchmarkBTreeSearchNotFound(b *testing.B) {

	if testSrc == nil {
		testSrc = makeTestSrc(cnt, keyLen, valLen)
	}

	bt := testSrc.bt
	searchKey := fmt.Sprintf("%snot found", testSrc.searchKey)
	searchItem := &testKV{key: searchKey, val: testSrc.searchValue}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = bt.Get(searchItem)
	}
}

func makeTestSrc(cnt int64, keyLen, valLen uint32) *testSrcType {

	srcKeys, err := makeStrings(cnt, int64(keyLen))
	if err != nil {
		panic(fmt.Sprintf("make source keys error: %v", err))
	}

	srcVals, err := makeByteSlices(cnt, int64(valLen))
	if err != nil {
		panic(fmt.Sprintf("make source value error: %v", err))
	}

	// make test trie
	keys := make([][]byte, cnt)
	for i, k := range srcKeys {
		keys[i] = splitStringTo4BitWords(k)
	}

	vals := makeKVElts(srcKeys, srcVals)

	t, err := New(keys, vals)
	if err != nil {
		panic(fmt.Sprintf("build trie failed: %v", err))
	}

	t.Squash()

	ct := NewCompactedTrie(testKVConv{keySize: keyLen, valSize: valLen})
	err = ct.Compact(t)
	if err != nil {
		panic(fmt.Sprintf("build compacted trie failed: %v", err))
	}

	// make test map
	m := make(map[string][]byte, cnt)
	for i := 0; i < len(srcKeys); i++ {
		m[srcKeys[i]] = srcVals[i]
	}

	// make test btree
	bt := btree.New(btreeDegree)

	for _, v := range vals {
		bt.ReplaceOrInsert(v)
	}

	// get search key
	r := rand.New(rand.NewSource(time.Now().Unix()))
	idx := r.Int63n(cnt)

	searchKey := srcKeys[idx]
	searchVal := srcVals[idx]

	return &testSrcType{
		srcKeys:   srcKeys,
		srcValues: srcVals,

		root: ct,
		m:    m,
		bt:   bt,

		searchKey:   searchKey,
		searchValue: searchVal,
	}
}

func trieSearchTestKV(ct *CompactedTrie, key string) []byte {
	eq := ct.SearchStringEqual(key)
	if eq == nil {
		return nil
	}

	val := eq.(*testKV)

	if strings.Compare(val.key, key) != 0 {
		return nil
	}

	return val.val
}

func sortedArraySearch(keys []string, values [][]byte, searchKey string) []byte {
	keyCnt := len(keys)

	idx := sort.Search(
		keyCnt,
		func(i int) bool {
			return strings.Compare(keys[i], searchKey) >= 0
		},
	)

	if idx < keyCnt && strings.Compare(keys[idx], searchKey) == 0 {
		return values[idx]
	}

	return nil
}
