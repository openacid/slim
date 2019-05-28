package trie

import "fmt"

// trieStringly is a wrapper that implements tree.Tree .
// It is a helper to convert trie to string.
//
// Since 0.5.1
type trieStringly struct {
	tnode *Node
}

// Child implements tree.Tree
//
// Since 0.5.1
func (s *trieStringly) Child(node, branch interface{}) interface{} {

	if node == nil && branch == nil {
		// root
		return s.tnode
	}

	n := s.trieNode(node)
	b := branch.(int)
	return n.Children[b]
}

// Labels implements tree.Tree
//
// Since 0.5.1
func (s *trieStringly) Labels(node interface{}) []interface{} {
	n := s.trieNode(node)
	rst := []interface{}{}
	for _, b := range n.Branches {
		rst = append(rst, b)
	}
	return rst
}

// NodeID implements tree.Tree
//
// Since 0.5.1
func (s *trieStringly) NodeID(node interface{}) string {
	return ""
}

// LabelInfo implements tree.Tree
//
// Since 0.5.5
func (s *trieStringly) LabelInfo(label interface{}) string {
	l := label.(int)
	if l == -1 {
		return "$"
	}
	return fmt.Sprintf("%d", l)
}

// NodeInfo implements tree.Tree
//
// Since 0.5.1
func (s *trieStringly) NodeInfo(node interface{}) string {
	n := s.trieNode(node)
	if n.Step <= 1 {
		return ""
	}
	return fmt.Sprintf("+%d", n.Step)
}

// LeafVal implements tree.Tree
//
// Since 0.5.1
func (s *trieStringly) LeafVal(node interface{}) (interface{}, bool) {
	n := s.trieNode(node)
	return n.Value, n.Value != nil
}

// trieNode convert a interface to SlimTrie node id.
func (s *trieStringly) trieNode(node interface{}) *Node {
	if node == nil {
		return s.tnode
	}
	return node.(*Node)
}
