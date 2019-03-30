// Package benchmark provides internally used benchmark support
package benchmark

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"testing"
	"time"
	"unsafe"

	"github.com/google/btree"
	"github.com/modood/table"
	"github.com/openacid/slim/trie"
)

// Config defines the variable inputs struct in one benchmark.
type Config struct {
	KeyCnt int
	KeyLen int
	ValLen int
}

// SearchResult show the key search result with a constructed data.
// Used to transfer benchmark result currently.
// SearchResult also defines the column titles when output to a chart.
type SearchResult struct {
	KeyCnt                int
	KeyLen                int
	ExsitingKeyNsPerOp    int64
	NonexsitentKeyNsPerOp int64
}

// MakeTrieSearchBench benchmark the trie search with existing and nonexistent
// key, return a slice of `TrieSearchCost`.
// `go test ...` conveniently.
func MakeTrieSearchBench(runs []Config) []*SearchResult {

	var spents = make([]*SearchResult, len(runs))

	for i, r := range runs {
		testSrc := MakeTestSrc(r.KeyCnt, r.KeyLen, r.ValLen)

		tr := testSrc.Slim

		// existing key
		existingRst := testing.Benchmark(MakeTrieBenchFunc(tr, testSrc.SearchKey))

		// nonexistent key
		searchKey := fmt.Sprintf("%snot found", testSrc.SearchKey)
		nonexistentRst := testing.Benchmark(MakeTrieBenchFunc(tr, searchKey))

		spents[i] = &SearchResult{
			KeyCnt:                r.KeyCnt,
			KeyLen:                r.KeyLen,
			ExsitingKeyNsPerOp:    existingRst.NsPerOp(),
			NonexsitentKeyNsPerOp: nonexistentRst.NsPerOp(),
		}
	}

	return spents
}
func MakeTestSrc(cnt int, keyLen, valLen int) *testSrcType {

	keys := makeSortedStrings(cnt, keyLen)
	srcVals := makeByteSlices(cnt, valLen)

	vals := makeKVElts(keys, srcVals)

	st, err := trie.NewSlimTrie(testKVConv{keySize: keyLen, valSize: valLen}, keys, vals)
	if err != nil {
		panic(err)
	}

	// make test map
	m := make(map[string][]byte, cnt)
	for i := 0; i < len(keys); i++ {
		m[keys[i]] = srcVals[i]
	}

	// make test btree
	bt := btree.New(2)

	for _, v := range vals {
		bt.ReplaceOrInsert(v)
	}

	// get search key
	r := rand.New(rand.NewSource(time.Now().Unix()))
	idx := r.Int63n(int64(cnt))

	searchKey := keys[idx]
	searchVal := srcVals[idx]

	return &testSrcType{
		Keys:   keys,
		Values: srcVals,

		Slim:  st,
		Map:   m,
		Btree: bt,

		SearchKey:   searchKey,
		SearchValue: searchVal,
	}
}

// testSrcType defines benchmark data source.
type testSrcType struct {
	Keys   []string
	Values [][]byte

	Slim  *trie.SlimTrie
	Map   map[string][]byte
	Btree *btree.BTree

	SearchKey   string
	SearchValue []byte
}

// TrieBenchKV defines a key-value struct to be used as a value in SlimTrie in test.
type TrieBenchKV struct {
	Key string
	Val []byte
}

// Less is used to implements google/btree.Item
func (kv TrieBenchKV) Less(than btree.Item) bool {
	anotherKV := than.(*TrieBenchKV)

	return kv.Key < anotherKV.Key
}

// testKVConv implements array.Converter to be a converter of TrieBenchKV.
type testKVConv struct {
	keySize int
	valSize int
}

func (c testKVConv) Encode(d interface{}) []byte {

	elt := d.(*TrieBenchKV)

	p := unsafe.Pointer(&elt)

	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, *(*uint64)(p))

	return b
}

func (c testKVConv) Decode(b []byte) (int, interface{}) {

	size := 8
	s := b[:size]

	buf := binary.LittleEndian.Uint64(s)

	// addr of uint64 == addr of elt pointer
	p := unsafe.Pointer(&buf)

	// convter pointer
	covP := *(*unsafe.Pointer)(p)

	// addr of *TrieBenchKV
	eltP := (*TrieBenchKV)(covP)

	return 8, eltP
}

func (c testKVConv) GetSize(d interface{}) int {
	return 8
}

func (c testKVConv) GetEncodedSize(b []byte) int {
	return 8
}

// OutputToChart output the benchmark result to a chart.
func OutputToChart(header string, body []*SearchResult) string {

	b := table.Table(body)
	return fmt.Sprintf("%s:\n%s\n", header, b)
}

func makeSortedStrings(cnt, leng int) []string {
	rsts := make([]string, cnt)

	for i := 0; i < cnt; i++ {
		rsts[i] = RandomString(leng)
	}

	sort.Strings(rsts)
	return rsts
}

func makeByteSlices(cnt, leng int) [][]byte {
	rsts := make([][]byte, cnt)

	for i := int(0); i < cnt; i++ {
		rsts[i] = RandomBytes(leng)
	}

	return rsts
}

func RandomString(leng int) string {
	return string(RandomBytes(leng))
}

func RandomBytes(leng int) []byte {
	bs := make([]byte, leng)
	n, err := crand.Read(bs)
	if err != nil {
		panic(err)
	}
	if n != leng {
		panic("not read enough")
	}
	return bs
}

func makeKVElts(srcKeys []string, srcVals [][]byte) []*TrieBenchKV {
	vals := make([]*TrieBenchKV, len(srcKeys))
	for i, k := range srcKeys {
		vals[i] = &TrieBenchKV{Key: k, Val: srcVals[i]}
	}
	return vals
}

func trieSearchTestKV(ct *trie.SlimTrie, key string) []byte {

	eq := ct.Get(key)
	if eq == nil {
		return nil
	}

	val := eq.(*TrieBenchKV)

	if strings.Compare(val.Key, key) != 0 {
		return nil
	}

	return val.Val
}

func MakeTrieBenchFunc(st *trie.SlimTrie, searchKey string) func(*testing.B) {

	return func(b *testing.B) {

		for i := 0; i < b.N; i++ {
			_ = trieSearchTestKV(st, searchKey)
		}

	}
}

func MakeMapBenchFunc(m map[string][]byte, searchKey string) func(*testing.B) {

	return func(b *testing.B) {

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_ = m[searchKey]
		}

		b.StopTimer()
	}
}

func MakeArrayBenchFunc(keys []string, values [][]byte, searchKey string) func(*testing.B) {

	return func(b *testing.B) {

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_ = sortedArraySearch(keys, values, searchKey)
		}

		b.StopTimer()
	}
}

func MakeBTreeBenchFunc(bt *btree.BTree, searchItem *TrieBenchKV) func(*testing.B) {
	return func(b *testing.B) {

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_ = bt.Get(searchItem)
		}

		b.StopTimer()
	}
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
