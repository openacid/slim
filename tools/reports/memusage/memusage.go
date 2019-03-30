package main

import (
	crand "crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	mrand "math/rand"
	"runtime"
	"sort"

	"github.com/openacid/slim/encode"
	"github.com/openacid/slim/trie"
)

type testKV struct {
	key string
	val uint16
}

type testKVConv struct {
	EltSize uint32
}

func (c testKVConv) Encode(d interface{}) []byte {

	elt := d.(*testKV)

	kSize := len(elt.key)
	vSize := 2

	b := make([]byte, kSize+vSize)
	for i, k := range elt.key {
		b[i] = byte(k)
	}

	binary.LittleEndian.PutUint16(b[kSize:], elt.val)

	return b
}

func (c testKVConv) Decode(b []byte) (int, interface{}) {

	elt := testKV{}
	vSize := uint32(2) // an uint16
	kSize := c.EltSize - vSize

	elt.key = string(b[:kSize])
	elt.val = binary.LittleEndian.Uint16(b[kSize:])

	return int(c.EltSize), elt
}

func (c testKVConv) GetSize(d interface{}) int {
	return int(c.EltSize)
}
func (c testKVConv) GetEncodedSize(b []byte) int {
	return int(c.EltSize)
}

func readMem() int64 {
	var stats runtime.MemStats
	for i := 0; i < 10; i++ {
		runtime.GC()
	}

	runtime.ReadMemStats(&stats)
	return int64(stats.Alloc)
}

func main() {
	compareTrieMapMemUse()
}

func compareTrieMapMemUse() {

	writeTableHeader()

	for _, cnt := range []int64{1000, 2000, 5000} {
		for _, l := range []int64{32, 64, 256, 512} {

			trieSize, err := getTrieMem(cnt, l)
			if err != nil {
				fmt.Printf("failed to get trie size: %v", err)
			}

			mapSize := getMapMem(cnt, l)
			if err != nil {
				fmt.Printf("failed to get map size: %v", err)
			}

			// make key + value as a value in trie
			kvTrieSize, err := getKVTrieMem(cnt, l)
			if err != nil {
				fmt.Printf("failed to get KV trie size: %v", err)
			}

			mapAvg := float64(mapSize) / float64(cnt)
			trieAvg := float64(trieSize) / float64(cnt)
			kvTrieAvg := float64(kvTrieSize) / float64(cnt)

			writeTableRow(cnt, l, 2, trieAvg, mapAvg, kvTrieAvg)
		}
	}
}

func writeTableHeader() {

	fmt.Printf("| %s | %s | %s | %s | %s | %s |\n",
		"Key Count", "Key Length", "Value Size", "Trie Size (Byte/key)", "Map Size (Byte/key)",
		"KV Trie Size (Byte/key)")

	fmt.Printf("| --- | --- | --- | --- | --- | --- |\n")
}

func writeTableRow(cnt, kLen, vLen int64, trieAvg, mapAvg, kvTrieAvg float64) {

	fmt.Printf("| %5d | %5d | %5d | %6.1f | %6.1f | %6.1f |\n",
		cnt, kLen, vLen, trieAvg, mapAvg, kvTrieAvg)
}

func makeKeys(kCnt, kLen int64) []string {
	keys := make([]string, kCnt)
	for i := int64(0); i < kCnt; i++ {
		key := make([]byte, kLen)
		if _, err := io.ReadFull(crand.Reader, key); err != nil {
			panic("making random keys" + err.Error())
		}
		keys[i] = string(key)
	}

	sort.Strings(keys)
	return keys
}

func makeVals(cnt int64) []uint16 {
	vals := make([]uint16, cnt)
	for i := int64(0); i < cnt; i++ {
		vals[i] = uint16(mrand.Intn(1 << 15))
	}
	return vals
}

func makeKVs(keys []string, vals []uint16) (kvs []*testKV) {
	kvs = make([]*testKV, len(keys))
	for i, k := range keys {
		kvs[i] = &testKV{key: k, val: vals[i]}
	}
	return
}

func getTrieMem(keyCnt, keyLen int64) (size int64, err error) {

	memStart := readMem()

	keys := makeKeys(keyCnt, keyLen)
	vals := makeVals(keyCnt)

	t, err := trie.NewSlimTrie(encode.U16{}, keys, vals)
	if err != nil {
		return
	}

	keys = nil
	vals = nil

	memEnd := readMem()

	size = memEnd - memStart

	// reference them or memory is freed
	_ = t.Children
	_ = t.Steps
	_ = t.Leaves

	_ = keys
	_ = vals

	return
}

func getKVTrieMem(keyCnt, keyLen int64) (size int64, err error) {
	// make key + value as a value in trie

	memStart := readMem()

	keys := makeKeys(keyCnt, keyLen)
	vals := makeVals(keyCnt)
	kvs := makeKVs(keys, vals)

	t, err := trie.NewSlimTrie(testKVConv{EltSize: uint32(keyLen + 2)}, keys, kvs)
	if err != nil {
		return
	}

	keys = nil
	vals = nil
	kvs = nil

	memEnd := readMem()

	size = memEnd - memStart

	// reference them or memory is freed
	_ = t.Children
	_ = t.Steps
	_ = t.Leaves

	_ = keys
	_ = vals
	_ = kvs

	return
}

func getMapMem(keyCnt, keyLen int64) int64 {

	memStart := readMem()

	keys := makeKeys(keyCnt, keyLen)
	vals := makeVals(keyCnt)

	m := make(map[string]uint16, len(keys))

	for i := 0; i < len(keys); i++ {
		m[keys[i]] = vals[i]
	}

	keys = nil
	vals = nil

	memEnd := readMem()

	size := memEnd - memStart

	_ = keys
	_ = vals

	return size
}
