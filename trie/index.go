package trie

import "github.com/openacid/slim/marshal"

// DataReader defines methods for index data types.
//
// A data letting SlimIndex to index it should implement this interface.
type DataReader interface {
	// Read value at offset, of `key`.
	// Because SlimIndex does not store complete info of a key(to reduce memory
	// consumption).
	// Thus the offset SlimIndex returns might not be correct for a abscent key.
	// It is data providers' responsibility to check if the record at `offset`
	// has the exact `key`.
	Read(offset int64, key string) (string, bool)
}

// OffsetIndexItem defines data types for a offset-based index, such as an index
// of on-disk records.
type OffsetIndexItem struct {
	// Key is the item identifier to identify a record.
	Key string
	// Offset is the position of this record in some data storage, E.g. the
	// offset in a file where this record is stored.
	Offset int64
}

// SlimIndex provides a commonly used index structure, it contains an SlimTrie
// instance as index and data provider `DataReader`.
type SlimIndex struct {
	SlimTrie
	DataReader
}

// NewSlimIndex provides a handy index implmentation.
//
// It creates a memory efficient index with SlimTrie.
func NewSlimIndex(index []OffsetIndexItem, dr DataReader) (*SlimIndex, error) {

	l := len(index)
	keys := make([]string, 0, l)
	offsets := make([]int64, 0, l)
	for i := 0; i < l; i++ {
		keys = append(keys, index[i].Key)
		offsets = append(offsets, index[i].Offset)
	}

	st, err := NewSlimTrie(marshal.I64{}, keys, offsets)
	if err != nil {
		return nil, err
	}

	return &SlimIndex{*st, dr}, nil
}

// Get returns the value of `key` which is found in `SlimIndex.DataReader`, and
// a bool value indicate if the `key` is found or not.
func (si *SlimIndex) Get(key string) (string, bool) {
	o := si.SlimTrie.Get(key)
	if o == nil {
		return "", false
	}

	offset := o.(int64)

	return si.DataReader.Read(offset, key)
}
