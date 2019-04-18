// Package benchmark provides internally used benchmark support
package benchmark

import (
	"encoding/binary"
	"math/rand"
	"sort"
	"strings"
	"testing"
	"time"
	"unsafe"

	"github.com/google/btree"
	"github.com/openacid/slim/benchhelper"
	"github.com/openacid/slim/encode"
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

// GetResult represent the ns/Get() for virous key count and several predefined
// key length = 64, 128, 256
type GetResult struct {
	KeyCount int
	K64      int
	K128     int
	K256     int
}

// MemResult is a alias of GetResult
type MemResult GetResult

var Rec byte

// BenchGet benchmark the Get() of present key.
func GetPresent(keyCounts []int) []GetResult {

	var rst = make([]GetResult, 0, len(keyCounts))

	for _, n := range keyCounts {

		r := GetResult{
			KeyCount: n,
			K64:      benchGet(NewGetSetting(n, 64, 2), "present"),
			K128:     benchGet(NewGetSetting(n, 128, 2), "present"),
			K256:     benchGet(NewGetSetting(n, 256, 2), "present"),
		}

		rst = append(rst, r)
	}

	return rst
}

func GetAbsent(keyCounts []int) []GetResult {

	var rst = make([]GetResult, 0, len(keyCounts))

	for _, n := range keyCounts {

		r := GetResult{
			KeyCount: n,
			K64:      benchGet(NewGetSetting(n, 64, 2), "absent"),
			K128:     benchGet(NewGetSetting(n, 128, 2), "absent"),
			K256:     benchGet(NewGetSetting(n, 256, 2), "absent"),
		}

		rst = append(rst, r)
	}

	return rst
}

func Mem(keyCounts []int) []MemResult {

	rst := make([]MemResult, 0)
	for _, n := range keyCounts {
		r := MemResult{
			KeyCount: n,
			K64:      int(slimtrieMem(n, 64)) / n,
			K128:     int(slimtrieMem(n, 128)) / n,
			K256:     int(slimtrieMem(n, 256)) / n,
		}

		rst = append(rst, r)
	}
	return rst
}

func slimtrieMem(keyCnt, keyLen int) int64 {

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

func benchGet(setting *GetSetting, typ string) int {

	var keys []string

	if typ == "present" {
		keys = setting.Keys
	} else {
		keys = setting.AbsentKeys
	}

	rst := testing.Benchmark(
		func(b *testing.B) {

			st := setting.Slim
			n := len(setting.Keys)

			var rec byte

			for i := 0; i < b.N; i++ {
				v := getTestKV(st, keys[i%n])
				rec += v[0]
			}

			Rec = rec
		})

	return int(rst.NsPerOp())

}

func NewGetSetting(cnt int, keyLen, valLen int) *GetSetting {

	ks := benchhelper.RandSortedStrings(cnt*2, keyLen)

	keys := make([]string, cnt)
	absentkeys := make([]string, cnt)

	for i := 0; i < cnt; i++ {
		keys[i] = ks[i*2]
		absentkeys[i] = ks[i*2+1]
	}

	vals := benchhelper.RandByteSlices(cnt, valLen)

	elts := makeKVElts(keys, vals)

	st, err := trie.NewSlimTrie(testKVConv{keySize: keyLen, valSize: valLen}, keys, elts)
	if err != nil {
		panic(err)
	}

	// make test map
	m := make(map[string][]byte, cnt)
	for i := 0; i < len(keys); i++ {
		m[keys[i]] = vals[i]
	}

	// make test btree
	bt := btree.New(2)

	for _, v := range elts {
		bt.ReplaceOrInsert(v)
	}

	// get search key
	r := rand.New(rand.NewSource(time.Now().Unix()))
	idx := r.Int63n(int64(cnt))

	searchKey := keys[idx]
	searchVal := vals[idx]

	return &GetSetting{
		Keys:   keys,
		Values: vals,

		AbsentKeys: absentkeys,

		Slim:  st,
		Map:   m,
		Btree: bt,

		SearchKey:   searchKey,
		SearchValue: searchVal,
	}
}

// GetSetting defines benchmark data source.
type GetSetting struct {
	Keys   []string
	Values [][]byte

	AbsentKeys []string

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

func makeKVElts(srcKeys []string, srcVals [][]byte) []*TrieBenchKV {
	vals := make([]*TrieBenchKV, len(srcKeys))
	for i, k := range srcKeys {
		vals[i] = &TrieBenchKV{Key: k, Val: srcVals[i]}
	}
	return vals
}

func getTestKV(ct *trie.SlimTrie, key string) []byte {

	eq := ct.Get(key)
	if eq == nil {
		return []byte{0}
	}

	val := eq.(*TrieBenchKV)

	if strings.Compare(val.Key, key) != 0 {
		return []byte{0}
	}

	return val.Val
}

func MakeTrieBenchFunc(st *trie.SlimTrie, searchKey string) func(*testing.B) {

	return func(b *testing.B) {

		for i := 0; i < b.N; i++ {
			_ = getTestKV(st, searchKey)
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
