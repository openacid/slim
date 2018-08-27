package trie

import (
	crand "crypto/rand"
	"encoding/binary"
	"io"
	"sort"
	"unsafe"

	"github.com/google/btree"
)

// testKV defines a key-value struct to be used as a value in compactedTrie.
type testKV struct {
	key string
	val []byte
}

// testKVConv implements array.EltConverter to be a converter of testKV.
type testKVConv struct {
	keySize uint32
	valSize uint32
}

// Less is create to implements google/btree.Item
func (kv testKV) Less(than btree.Item) bool {
	anotherKV := than.(*testKV)

	if kv.key < anotherKV.key {
		return true
	}

	return false
}

func (c testKVConv) MarshalElt(d interface{}) []byte {

	elt := d.(*testKV)

	p := unsafe.Pointer(&elt)

	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, *(*uint64)(p))

	return b
}

func (c testKVConv) UnmarshalElt(b []byte) (uint32, interface{}) {

	size := uint32(8)
	s := b[:size]

	buf := binary.LittleEndian.Uint64(s)

	// addr of uint64 == addr of elt pointer
	p := unsafe.Pointer(&buf)

	// convter pointer
	covP := *(*unsafe.Pointer)(p)

	// addr of *testKV
	eltP := (*testKV)(covP)

	return uint32(8), eltP
}

func (c testKVConv) GetMarshaledEltSize(b []byte) uint32 {
	return uint32(8)
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
		b := byte(s[i])
		words[2*i] = (b & 0xf0) >> 4
		words[2*i+1] = b & 0x0f
	}
	return words
}
