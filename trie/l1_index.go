package trie

import (
	"xec/prototype"
)

type L1Index struct {
	prototype.L1Index
}

func NewL1Index() *L1Index {
	return &L1Index{
		L1Index: prototype.L1Index{
			Keys:    make([]string, 0),
			Offsets: make([]int64, 0),
		},
	}
}

func (l1Idx *L1Index) Add(key string, offset int64) {
	l1Idx.Keys = append(l1Idx.Keys, key)
	l1Idx.Offsets = append(l1Idx.Offsets, offset)
}

func (l1Idx *L1Index) RPop() (key string, offset int64, ok bool) {
	length := len(l1Idx.Keys)
	if length == 0 {
		return key, offset, false
	}

	key = l1Idx.Keys[length-1]
	l1Idx.Keys = l1Idx.Keys[:length-1]

	offset = l1Idx.Offsets[length-1]
	l1Idx.Offsets = l1Idx.Offsets[:length-1]

	return key, offset, true
}
