package trie

import (
	"xec/prototype"
)

/**
 * Keys                []string
 * L2IdxOffsets        []int64
 * NeedleOffsets       []int64
 * NeedleIDListOffsets []int64  offsets relative to the start addr of NeedleIDList
 */
type L1Index struct {
	prototype.L1Index
}

func NewL1Index() *L1Index {
	return &L1Index{
		L1Index: prototype.L1Index{
			Keys:                make([]string, 0),
			L2IdxOffsets:        make([]int64, 0),
			NeedleOffsets:       make([]int64, 0),
			NeedleIDListOffsets: make([]int64, 0),
		},
	}
}

func (l1Idx *L1Index) Add(key string, l2Offset, ndlOffset, ndlIDLOffset int64) {
	l1Idx.Keys = append(l1Idx.Keys, key)
	l1Idx.L2IdxOffsets = append(l1Idx.L2IdxOffsets, l2Offset)
	l1Idx.NeedleOffsets = append(l1Idx.NeedleOffsets, ndlOffset)
	l1Idx.NeedleIDListOffsets = append(l1Idx.NeedleIDListOffsets, ndlIDLOffset)
}

func (l1Idx *L1Index) RPop() (key string, l2Offset, ndlOffset, ndlIDLOffset int64, ok bool) {
	length := len(l1Idx.Keys)
	if length == 0 {
		return key, l2Offset, ndlOffset, ndlIDLOffset, false
	}

	key = l1Idx.Keys[length-1]
	l1Idx.Keys = l1Idx.Keys[:length-1]

	l2Offset = l1Idx.L2IdxOffsets[length-1]
	l1Idx.L2IdxOffsets = l1Idx.L2IdxOffsets[:length-1]

	ndlOffset = l1Idx.NeedleOffsets[length-1]
	l1Idx.NeedleOffsets = l1Idx.NeedleOffsets[:length-1]

	ndlIDLOffset = l1Idx.NeedleIDListOffsets[length-1]
	l1Idx.NeedleIDListOffsets = l1Idx.NeedleIDListOffsets[:length-1]

	return key, l2Offset, ndlOffset, ndlIDLOffset, true
}
