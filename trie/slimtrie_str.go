package trie

import (
	"fmt"
	"math/bits"
)

// slimTrieStringly is a wrapper that implements tree.Tree .
// It is a helper to convert SlimTrie to string.
//
// Since 0.5.1
type slimTrieStringly struct {
	st *SlimTrie
}

// Child implements tree.Tree
//
// Since 0.5.1
func (s *slimTrieStringly) Child(node, branch interface{}) interface{} {

	b := branch.(int)

	bitmap, child0ID, _ := s.st.getChild(stNodeID(node))
	if bitmap&(1<<uint(b)) != 0 {
		nth := bits.OnesCount16(bitmap << (16 - uint16(b)))
		return child0ID + int32(nth)
	} else {
		return nil
	}
}

// Labels implements tree.Tree
//
// Since 0.5.1
func (s *slimTrieStringly) Labels(node interface{}) []interface{} {

	rst := []interface{}{}

	bitmap, _, _ := s.st.getChild(stNodeID(node))

	for b := uint(0); b < 16; b++ {
		if bitmap&(1<<b) > 0 {
			rst = append(rst, int(b))
		}
	}

	return rst
}

// NodeID implements tree.Tree
//
// Since 0.5.1
func (s *slimTrieStringly) NodeID(node interface{}) string {
	return fmt.Sprintf("%03d", stNodeID(node))
}

// LabelInfo implements tree.Tree
//
// Since 0.5.1
func (s *slimTrieStringly) LabelInfo(label interface{}) string {
	return fmt.Sprintf("%d", label)
}

// NodeInfo implements tree.Tree
//
// Since 0.5.1
func (s *slimTrieStringly) NodeInfo(node interface{}) string {
	step, hasStep := s.st.Steps.Get(stNodeID(node))
	if hasStep {
		return fmt.Sprintf("+%d", step)
	}
	return ""
}

// LeafVal implements tree.Tree
//
// Since 0.5.1
func (s *slimTrieStringly) LeafVal(node interface{}) (interface{}, bool) {
	return s.st.Leaves.Get(stNodeID(node))
}

// stNodeID convert a interface to SlimTrie node id.
func stNodeID(node interface{}) int32 {
	if node == nil {
		node = int32(0)
	}
	return node.(int32)
}
