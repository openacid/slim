// Package index provides a data index structure, contains a SlimTrie
// instance as index and a data provider `DataReader`.
//
// The purpose of `index` is to accelerate accessing external data such as a bunch
// of key-values on disk.
//
// SlimIndex is an index only, but not a full functional kv-map.
// The relationship of SlimIndex and its data just like that of B+ tree internal
// nodes and B+ tree leaf nodes:
// In a B+ tree only leaf nodes store record.  Internal nodes only help on
// locating a leaf node then a record.
package index

import (
	"github.com/openacid/slim/encode"
	"github.com/openacid/slim/trie"
)

// DataReader defines interface to let SlimIndex access the data it indexes.
type DataReader interface {
	// Read value at `offset`, of `key`.
	// Because the internal SlimTrie does not store complete info of a key(to
	// reduce memory consumption).
	// Thus the offset SlimTrie returns might not be correct for an abscent
	// record.
	// It is data providers' responsibility to check if the record at `offset`
	// has the exact `key`.
	Read(offset int64, key string) (string, bool)
}

// OffsetIndexItem defines data types for a offset-based index, such as an index
// of on-disk records.
type OffsetIndexItem struct {
	// Key is the record identity.
	Key string
	// Offset is the position of this record in its storage, E.g. the file offset
	// where this record is.
	Offset int64
}

// SlimIndex contains a SlimTrie instance as index and a data provider
// `DataReader`.
type SlimIndex struct {
	trie.SlimTrie
	DataReader
}

// NewSlimIndex creates SlimIndex instance.
//
// The keys in `index` must be in ascending order.
func NewSlimIndex(index []OffsetIndexItem, dr DataReader) (*SlimIndex, error) {

	l := len(index)
	keys := make([]string, 0, l)
	offsets := make([]int64, 0, l)
	for i := 0; i < l; i++ {
		keys = append(keys, index[i].Key)
		offsets = append(offsets, index[i].Offset)
	}

	st, err := trie.NewSlimTrie(encode.I64{}, keys, offsets)
	if err != nil {
		return nil, err
	}

	return &SlimIndex{*st, dr}, nil
}

// Get returns the value of `key` which is found by `SlimIndex.DataReader`, and
// a bool value indicating if the `key` is found or not.
func (si *SlimIndex) Get(key string) (string, bool) {
	o := si.SlimTrie.Get(key)
	if o == nil {
		return "", false
	}

	offset := o.(int64)

	return si.DataReader.Read(offset, key)
}
