package trie

import "errors"

var KeyNotExist = errors.New("Key not Exist")

type TrieNode struct {
	Children map[byte]*TrieNode
	Skip     uint16
	IsLeaf   bool
	Elt      []byte
}

func New(key, value [][]byte) (root *TrieNode) {

	root = &TrieNode{Children: make(map[byte]*TrieNode), Skip: 1}

	for i, k := range key {

		var node = root
		var n_pref int

		for j, b := range k {
			n_pref = j
			if node.Children[b] != nil {
				node = node.Children[b]
			} else {
				break
			}
		}

		for _, b := range k[n_pref:] {
			n := &TrieNode{Children: make(map[byte]*TrieNode), Skip: 1}
			node.Children[b] = n
			node = n
		}

		node.IsLeaf = true
		node.Elt = value[i]
	}

	return root
}

func (root *TrieNode) Squash() {

	for k, n := range root.Children {
		n.Squash()

		if n.IsLeaf {
			continue
		}

		if len(n.Children) == 1 {
			for _, child := range n.Children {
				child.Skip += 1
				root.Children[k] = child
			}
		}
	}
}

func (root *TrieNode) Search(key []byte) (value []byte, err error) {

	if root.IsLeaf {
		value = root.Elt
		return
	}

	node := root
	for i := 0; i < len(key); {

		if node.Children[key[i]] == nil {
			err = KeyNotExist
			return
		}

		node = node.Children[key[i]]
		i += int(node.Skip)

		if node.IsLeaf {
			if i == len(key) {
				value = node.Elt
				return
			}
		}

	}
	err = KeyNotExist
	return
}
