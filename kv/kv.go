package kv

import (
	"github.com/openacid/slim/encode"
	"github.com/openacid/slim/trie"
)

type KV struct {
	st     *trie.SlimTrie
	keys   []string
	values []int32
}

func NewKV(keys []string, values []int32) *KV {
	// st, err := trie.NewSlimTrie(encode.I32{}, keys, nil, trie.Opt{Complete: trie.Bool(true)})
	// st, err := trie.NewSlimTrie(encode.I32{}, keys, nil, trie.Opt{InnerPrefix: trie.Bool(true)})
	st, err := trie.NewSlimTrie(encode.I32{}, keys, nil)
	if err != nil {
		panic("fff")
	}

	return &KV{
		st:     st,
		keys:   keys,
		values: values,
	}

}

func (kv *KV) Get(key string) (int32, bool) {
	r := kv.st.GetIndex(key)
	if r == -1 {
		return 0, false
	}
	if kv.keys[r] != key {
		return 0, false
	}
	return kv.values[r], true

}
