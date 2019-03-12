package trie

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"sort"
	"strings"
	"testing"
	"time"
	"unsafe"

	"github.com/google/btree"
	"github.com/modood/table"
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

// testSrcType defines benchmark data source.
type testSrcType struct {
	srcKeys   []string
	srcValues [][]byte

	root *SlimTrie
	m    map[string][]byte
	bt   *btree.BTree

	searchKey   string
	searchValue []byte
}

// testKV defines a key-value struct to be used as a value in SlimTrie in test.
type testKV struct {
	key string
	val []byte
}

// Less is used to implements google/btree.Item
func (kv testKV) Less(than btree.Item) bool {
	anotherKV := than.(*testKV)

	return kv.key < anotherKV.key
}

// testKVConv implements array.Converter to be a converter of testKV.
type testKVConv struct {
	keySize uint32
	valSize uint32
}

func (c testKVConv) Marshal(d interface{}) []byte {

	elt := d.(*testKV)

	p := unsafe.Pointer(&elt)

	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, *(*uint64)(p))

	return b
}

func (c testKVConv) Unmarshal(b []byte) (int, interface{}) {

	size := 8
	s := b[:size]

	buf := binary.LittleEndian.Uint64(s)

	// addr of uint64 == addr of elt pointer
	p := unsafe.Pointer(&buf)

	// convter pointer
	covP := *(*unsafe.Pointer)(p)

	// addr of *testKV
	eltP := (*testKV)(covP)

	return 8, eltP
}

func (c testKVConv) GetMarshaledSize(b []byte) int {
	return 8
}

// OutputToChart output the benchmark result to a chart.
func OutputToChart(header string, body []*benchmark.SearchResult) string {

	b := table.Table(body)
	return fmt.Sprintf("%s:\n%s\n", header, b)
}

func makeStrings(cnt, leng int64) ([]string, error) {
	srcs, err := makeByteSlices(cnt, leng)
	if err != nil {
		return nil, err
	}

	rsts := make([]string, cnt)

	for i := int64(0); i < cnt; i++ {
		rsts[i] = string(srcs[i])
	}

	sort.Strings(rsts)
	return rsts, nil
}

func makeByteSlices(cnt, leng int64) ([][]byte, error) {
	rsts := make([][]byte, cnt)

	for i := int64(0); i < cnt; i++ {
		bs := make([]byte, leng)

		if _, err := io.ReadFull(crand.Reader, bs); err != nil {
			return nil, err
		}

		rsts[i] = bs
	}

	return rsts, nil
}

func makeKVElts(srcKeys []string, srcVals [][]byte) []*testKV {
	vals := make([]*testKV, len(srcKeys))
	for i, k := range srcKeys {
		vals[i] = &testKV{key: k, val: srcVals[i]}
	}
	return vals
}

func splitStringTo4BitWords(s string) []byte {

	lenSrc := len(s)
	words := make([]byte, lenSrc*2)

	for i := 0; i < lenSrc; i++ {
		b := s[i]
		words[2*i] = (b & 0xf0) >> 4
		words[2*i+1] = b & 0x0f
	}
	return words
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

	t, err := NewTrie(keys, vals)
	if err != nil {
		panic(fmt.Sprintf("build trie failed: %v", err))
	}

	t.Squash()

	ct, _ := NewSlimTrie(testKVConv{keySize: keyLen, valSize: valLen}, nil, nil)
	err = ct.LoadTrie(t)
	if err != nil {
		panic(fmt.Sprintf("build compacted trie failed: %v", err))
	}

	// make test map
	m := make(map[string][]byte, cnt)
	for i := 0; i < len(srcKeys); i++ {
		m[srcKeys[i]] = srcVals[i]
	}

	// make test btree
	bt := btree.New(2)

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

func trieSearchTestKV(ct *SlimTrie, key string) []byte {

	//_, eq, _ := ct.SearchString(key)
	eq := ct.Get(key)
	if eq == nil {
		return nil
	}

	val := eq.(*testKV)

	if strings.Compare(val.key, key) != 0 {
		return nil
	}

	return val.val
}

func makeTrieBenchFunc(tr *SlimTrie, searchKey string) func(*testing.B) {

	return func(b *testing.B) {

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_ = trieSearchTestKV(tr, searchKey)
		}

		b.StopTimer()
	}
}

func makeMapBenchFunc(m map[string][]byte, searchKey string) func(*testing.B) {

	return func(b *testing.B) {

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_ = m[searchKey]
		}

		b.StopTimer()
	}
}

func makeArrayBenchFunc(keys []string, values [][]byte, searchKey string) func(*testing.B) {

	return func(b *testing.B) {

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_ = sortedArraySearch(keys, values, searchKey)
		}

		b.StopTimer()
	}
}

func makeBTreeBenchFunc(bt *btree.BTree, searchItem *testKV) func(*testing.B) {
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

// MakeTrieSearchBench benchmark the trie search with existing and nonexistent
// key, return a slice of `TrieSearchCost`.
// `go test ...` conveniently.
func MakeTrieSearchBench(runs []benchmark.Config) []*benchmark.SearchResult {

	var spents = make([]*benchmark.SearchResult, len(runs))

	for i, r := range runs {
		testSrc := makeTestSrc(r.KeyCnt, r.KeyLen, r.ValLen)

		tr := testSrc.root

		// existing key
		existingRst := testing.Benchmark(makeTrieBenchFunc(tr, testSrc.searchKey))

		// nonexistent key
		searchKey := fmt.Sprintf("%snot found", testSrc.searchKey)
		nonexistentRst := testing.Benchmark(makeTrieBenchFunc(tr, searchKey))

		spents[i] = &benchmark.SearchResult{
			KeyCnt:                r.KeyCnt,
			KeyLen:                r.KeyLen,
			ExsitingKeyNsPerOp:    existingRst.NsPerOp(),
			NonexsitentKeyNsPerOp: nonexistentRst.NsPerOp(),
		}
	}

	return spents
}
