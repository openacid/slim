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
)

// run defines the variable inputs struct in one benchmark.
type run struct {
	keyCnt int64
	keyLen uint32
	valLen uint32
}

var runs = []run{
	{1, 1024, 2},
	{10, 1024, 2},
	{100, 1024, 2},
	{1000, 1024, 2},
	{1000, 512, 2},
	{1000, 256, 2},
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

// TrieSearchCost show the key search result with a constructed data.
// Used to transfer benchmark result currently.
// TrieSearchCost also defines the column titles when output to a chart.
type TrieSearchCost struct {
	KeyCnt                int64
	KeyLen                uint32
	ExsitingKeyNsPerOp    int64
	NonexsitentKeyNsPerOp int64
}

// TrieSearchTable defines the row content when output TrieSearchCost to a chart.
type TrieSearchTable struct {
	Header string
	Body   []*TrieSearchCost
}

// OutputToChart output the TrieSearchTable to a chart, and return a string result.
func (tb *TrieSearchTable) OutputToChart() string {

	body := table.Table(tb.Body)
	return fmt.Sprintf("%s:\n%s\n", tb.Header, body)
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

// MakeTrieSearchBench benchmark the trie search with existing and nonexistent key, return the
// result as a constructed data `TrieSearchTable`. Then you can get the benchmark result without
// `go test ...` conveniently.
func MakeTrieSearchBench() *TrieSearchTable {

	var trieSearchCostList = make([]*TrieSearchCost, len(runs))

	for i, r := range runs {
		testSrc := makeTestSrc(r.keyCnt, r.keyLen, r.valLen)

		tr := testSrc.root

		// existing key
		existingRst := testing.Benchmark(makeTrieBenchFunc(tr, testSrc.searchKey))

		// nonexistent key
		searchKey := fmt.Sprintf("%snot found", testSrc.searchKey)
		nonexistentRst := testing.Benchmark(makeTrieBenchFunc(tr, searchKey))

		trieSearchCostList[i] = &TrieSearchCost{
			KeyCnt:                r.keyCnt,
			KeyLen:                r.keyLen,
			ExsitingKeyNsPerOp:    existingRst.NsPerOp(),
			NonexsitentKeyNsPerOp: nonexistentRst.NsPerOp(),
		}
	}

	return &TrieSearchTable{
		Header: "cost of trie search with existing & existent key",
		Body:   trieSearchCostList,
	}
}
