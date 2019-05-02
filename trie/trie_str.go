package trie

import "fmt"

// trieStringly is a wrapper that implements treestr.Tree .
// It is a helper to convert trie to string.
//
// Since 0.5.1
type trieStringly struct {
	tnode *Node
}

type trieBranch int

func (tb trieBranch) String() string {
	if tb == -1 {
		return "$"
	}
	return fmt.Sprintf("%d", tb)
}

// Child implements treestr.Tree
//
// Since 0.5.1
func (s *trieStringly) Child(node, branch interface{}) interface{} {

	n := s.trieNode(node)
	b := branch.(trieBranch)
	return n.Children[int(b)]
}

// Branches implements treestr.Tree
//
// Since 0.5.1
func (s *trieStringly) Branches(node interface{}) []interface{} {
	n := s.trieNode(node)
	rst := []interface{}{}
	for _, b := range n.Branches {
		rst = append(rst, trieBranch(b))
	}
	return rst
}

// NodeID implements treestr.Tree
//
// Since 0.5.1
func (s *trieStringly) NodeID(node interface{}) string {
	return ""
}

// NodeInfo implements treestr.Tree
//
// Since 0.5.1
func (s *trieStringly) NodeInfo(node interface{}) string {
	n := s.trieNode(node)
	if n.Step <= 1 {
		return ""
	}
	return fmt.Sprintf("+%d", n.Step)
}

// LeafVal implements treestr.Tree
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
