// Package kvmap provides a readonly key-value map based on SlimTrie
package kvmap

import (
	"github.com/openacid/errors"
	"github.com/openacid/slim/index"
)

// Item defines a key-value struct to be used as a value in SlimTrie in test.
type Item struct {
	Key string
	Val string
}

// KVMap is a readonly key-value map struct
type KVMap struct {
	Items []Item
	Index *index.SlimIndex
}

// Less compares it with another Item
func (kv *Item) Less(than Item) bool {
	anotherKV := than

	return kv.Key < anotherKV.Key
}

// Read implements index.DataReader
func (kv *KVMap) Read(offset int64, key string) (string, bool) {
	i := int(offset)
	if kv.Items[i].Key == key {
		return kv.Items[i].Val, true
	}
	return "", false
}

// Get retrieve the value and a bool indicating if the found.
func (kv *KVMap) Get(key string) (string, bool) {
	return kv.Index.Get(key)
}

// NewKVMap create a *KVMap instance.
func NewKVMap(kvs []Item) (*KVMap, error) {

	n := len(kvs)

	kv := &KVMap{Items: kvs}

	iitems := make([]index.OffsetIndexItem, n)
	for i := 0; i < n; i++ {
		iitems[i] = index.OffsetIndexItem{Key: kvs[i].Key, Offset: int64(i)}
	}

	ii, err := index.NewSlimIndex(iitems, kv)
	if err != nil {
		return nil, errors.Wrapf(err, "failure creating KVMap")
	}
	kv.Index = ii

	return kv, nil
}
