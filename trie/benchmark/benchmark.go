// Package benchmark provides internally used benchmark support
package benchmark

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/google/btree"
	"github.com/openacid/low/size"
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
	KeyCount int `tw-title:"key-count"`
	K64      int `tw-title:"k=64"`
	K128     int `tw-title:"k=128"`
	K256     int `tw-title:"k=256"`
}

// FPRResult represent the false positive rate.
type FPRResult struct {
	KeyCount int     `tw-title:"key-count"`
	FPR      float64 `tw-title:"fpr" tw-fmt:"%.3f%%"`
}

// MemResult is a alias of GetResult
type MemResult GetResult

var Rec int32

// GetPresent benchmark the Get() of present key.
func GetPresent(keyCounts []int) []GetResult {

	var rst = make([]GetResult, 0, len(keyCounts))

	for _, n := range keyCounts {

		r := GetResult{
			KeyCount: n,
			K64:      benchGet(NewGetSetting(n, 64), "present"),
			K128:     benchGet(NewGetSetting(n, 128), "present"),
			K256:     benchGet(NewGetSetting(n, 256), "present"),
		}

		rst = append(rst, r)
	}

	return rst
}

// GetAbsent benchmark the Get() of absent key.
func GetAbsent(keyCounts []int) []GetResult {

	var rst = make([]GetResult, 0, len(keyCounts))

	for _, n := range keyCounts {

		r := GetResult{
			KeyCount: n,
			K64:      benchGet(NewGetSetting(n, 64), "absent"),
			K128:     benchGet(NewGetSetting(n, 128), "absent"),
			K256:     benchGet(NewGetSetting(n, 256), "absent"),
		}

		rst = append(rst, r)
	}

	return rst
}

// GetFPR estimate false positive rate(FPR) for Get.
func GetFPR(keyCounts []int) []FPRResult {

	var rst = make([]FPRResult, 0, len(keyCounts))

	keyLen := 64
	r := 100
	for _, n := range keyCounts {

		keys := benchhelper.RandSortedStrings(n, keyLen, nil)
		nAbsent := n * r

		present := map[string]bool{}
		for _, k := range keys {
			present[k] = true
		}

		vals := make([]uint16, n)
		st, err := trie.NewSlimTrie(encode.U16{}, keys, vals)
		if err != nil {
			panic(err)
		}

		fp := float64(0)

		for i := 0; i < nAbsent; {
			k := benchhelper.RandString(keyLen, nil)
			if _, ok := present[k]; ok {
				continue
			}

			_, found := st.Get(k)
			if found {
				fp++
			}
			i++
		}

		fpr := fp / float64(nAbsent)

		r := FPRResult{
			KeyCount: n,
			FPR:      fpr,
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
			K64:      int(slimtrieMem(n, 64)),
			K128:     int(slimtrieMem(n, 128)),
			K256:     int(slimtrieMem(n, 256)),
		}

		rst = append(rst, r)
	}
	return rst
}

func slimtrieMem(keyCnt, keyLen int) int64 {

	keys := benchhelper.RandSortedStrings(keyCnt, keyLen, nil)
	vals := make([]uint16, keyCnt)
	for i := 0; i < len(keys); i++ {
		vals[i] = uint16(i)
	}

	t, err := trie.NewSlimTrie(encode.U16{}, keys, vals)
	if err != nil {
		panic(err)
	}

	size := size.Of(t) - size.Of(vals)

	return int64(size) * 8 / int64(keyCnt)
}

func benchGet(setting *GetSetting, typ string) int {

	var keys []string

	if typ == "present" {
		keys = setting.Keys
	} else {
		keys = setting.AbsentKeys
	}

	st := setting.SlimKV
	n := len(setting.Keys)
	mask := 1
	for ; (mask << 1) <= n; mask <<= 1 {
	}

	var rec int32

	rst := testing.Benchmark(
		func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				v := st.Get(keys[i&mask])
				rec += v
			}
		})

	Rec = rec

	return int(rst.NsPerOp())
}

func NewGetSetting(cnt int, keyLen int) *GetSetting {

	ks := benchhelper.RandSortedStrings(cnt*2, keyLen, nil)

	keys := make([]string, cnt)
	absentkeys := make([]string, cnt)

	for i := 0; i < cnt; i++ {
		keys[i] = ks[i*2]
		absentkeys[i] = ks[i*2+1]
	}

	vals := make([]int32, cnt)
	for i := 0; i < cnt; i++ {
		vals[i] = int32(i)
	}

	elts := makeKVElts(keys, vals)

	st, err := trie.NewSlimTrie(encode.I32{}, keys, vals)
	if err != nil {
		panic(err)
	}

	// make test map
	m := make(map[string]int32, cnt)
	for i := 0; i < len(keys); i++ {
		m[keys[i]] = vals[i]
	}

	// make test btree
	bt := btree.New(32)

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

		SlimKV: &slimKV{Elts: elts, slim: st},
		Map:    m,
		Btree:  bt,

		SearchKey:   searchKey,
		SearchValue: searchVal,
	}
}

// GetSetting defines benchmark data source.
type GetSetting struct {
	Keys   []string
	Values []int32

	AbsentKeys []string

	SlimKV *slimKV
	Map    map[string]int32
	Btree  *btree.BTree

	SearchKey   string
	SearchValue int32
}

type slimKV struct {
	// SlimTrie as an index
	slim *trie.SlimTrie
	// full key-values
	Elts []*KVElt
}

func (s *slimKV) Get(key string) int32 {
	idx, found := s.slim.Get(key)
	if !found {
		return -1
	}

	i := idx.(int32)

	elt := s.Elts[i]
	if strings.Compare(elt.Key, key) != 0 {
		return -1
	}

	return elt.Val
}

// KVElt defines a key-value struct to be used as a value in SlimTrie in test.
type KVElt struct {
	Key string
	Val int32
}

// Less is used to implements google/btree.Item
func (kv *KVElt) Less(than btree.Item) bool {
	o := than.(*KVElt)
	return kv.Key < o.Key
}

func makeKVElts(srcKeys []string, srcVals []int32) []*KVElt {
	elts := make([]*KVElt, len(srcKeys))
	for i, k := range srcKeys {
		elts[i] = &KVElt{Key: k, Val: srcVals[i]}
	}
	return elts
}
