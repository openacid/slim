package trie

import (
	"fmt"
	"sort"
	"strings"

	"github.com/openacid/errors"
	"github.com/openacid/slim/typehelper"
)

// ErrDuplicateKeys indicates two keys are identical.
var ErrDuplicateKeys = errors.New("keys can not be duplicate")

// ErrValuesNotSlice is returned if a value to create Trie is not slice type.
var ErrValuesNotSlice = errors.New("values must be in a slice")

// ErrKVLenNotMatch means the keys and values to create Trie has different
// number of elements.
var ErrKVLenNotMatch = errors.New("length of keys and values not equal")

// ErrKeyOutOfOrder means keys to create Trie are not ascendingly ordered.
var ErrKeyOutOfOrder = errors.New("keys not ascending sorted")

// Node is a Trie node.
type Node struct {
	Children map[int]*Node
	Branches []int
	Step     uint16
	Value    interface{}

	isStartLeaf bool

	NodeCnt int
}

const leafBranch = -1

// NewTrie creates a Trie from a serial of ascendingly ordered keys and corresponding values.
//
// `values` must be a slice.
func NewTrie(keys [][]byte, values interface{}) (root *Node, err error) {

	root = &Node{Children: make(map[int]*Node), Step: 1}

	if keys == nil {
		return
	}

	valSlice, ok := typehelper.ToSlice(values)
	if !ok {
		err = ErrValuesNotSlice
		return
	}

	if len(keys) != len(valSlice) {
		err = ErrKVLenNotMatch
		return
	}

	for i := 0; i < len(keys); i++ {
		key := keys[i]
		_, err = root.Append(key, valSlice[i], false, false)
		if err != nil {
			err = errors.Wrapf(err, "trie failed to add kvs")
			return
		}
	}

	return
}

// toStrings convert a Trie to human readalble representation.
// TODO add example.
func (r *Node) toStrings(cc int) []string {

	var line string
	if cc == leafBranch {
		line = fmt.Sprintf("  $(%d):", len(r.Branches))
	} else {
		line = fmt.Sprintf("%03d(%d):", cc, len(r.Branches))
	}

	rst := make([]string, 0, 64)

	if len(r.Branches) > 0 {

		for _, b := range r.Branches {
			subtrie := r.Children[b].toStrings(b)
			indent := strings.Repeat(" ", len(line))
			for _, s := range subtrie {
				if len(rst) == 0 {
					rst = append(rst, line+s)
				} else {
					rst = append(rst, indent+s)
				}
			}
		}

	} else {
		rst = append(rst, line)
	}

	return rst
}

// Squash compresses a Trie by removing single-branch nodes.
func (r *Node) Squash() int {

	var cnt int

	for k, n := range r.Children {
		cnt += n.Squash()

		if len(n.Branches) == 1 {
			if n.Branches[0] == leafBranch {
				continue
			}
			child := n.Children[n.Branches[0]]
			child.Step++
			r.Children[k] = child
			cnt++
		}
	}

	return cnt
}

// Search for `key` in a Trie.
//
// It returns 3 values of:
// The left sibling value.
// The matching value.
// The right sibling value.
//
// Any of them could be nil.
func (r *Node) Search(key []byte) (ltValue, eqValue, gtValue interface{}) {

	var eqNode = r
	var ltNode *Node
	var gtNode *Node

	for i := 0; ; {

		var br int
		if len(key) == i {
			br = leafBranch
		} else {
			br = int(key[i])
		}

		li, ri := neighborBranches(eqNode.Branches, br)
		if li >= 0 {
			ltNode = eqNode.Children[eqNode.Branches[li]]
		}
		if ri >= 0 {
			gtNode = eqNode.Children[eqNode.Branches[ri]]
		}

		eqNode = eqNode.Children[br]

		if eqNode == nil {
			break
		}

		if br == leafBranch {
			break
		}

		i += int(eqNode.Step)

		if i > len(key) {
			gtNode = eqNode
			eqNode = nil
			break
		}
	}

	if ltNode != nil {
		ltValue = ltNode.rightMost().Value
	}
	if gtNode != nil {
		gtValue = gtNode.leftMost().Value
	}
	if eqNode != nil {
		eqValue = eqNode.Value
	}

	return
}

func neighborBranches(branches []int, br int) (ltIndex, rtIndex int) {

	if len(branches) == 0 {
		return
	}

	var i int
	var b int

	for i, b = range branches {
		if b >= br {
			break
		}
	}

	if b == br {
		rtIndex = i + 1
		ltIndex = i - 1

		if rtIndex == len(branches) {
			rtIndex = -1
		}
		return
	}

	if b > br {
		rtIndex = i
		ltIndex = i - 1
		return
	}

	rtIndex = -1
	ltIndex = i

	return
}

func (r *Node) leftMost() *Node {

	node := r
	for {
		if len(node.Branches) == 0 {
			return node
		}

		firstBr := node.Branches[0]
		node = node.Children[firstBr]
	}
}

func (r *Node) rightMost() *Node {

	node := r
	for {
		if len(node.Branches) == 0 {
			return node
		}

		lastBr := node.Branches[len(node.Branches)-1]
		node = node.Children[lastBr]
	}
}

// Append adds a key-value pair into Trie.
//
// The key to add must be greater than any existent key in the Trie.
//
// `isStartLeaf` indicates if it should not be removed.
// `needSquash` indicates whether to compress the Trie after adding.
//
// It returns the leaf node representing the added key.
func (r *Node) Append(key []byte, value interface{}, isStartLeaf bool, needSquash bool) (leaf *Node, err error) {

	var node = r
	var j int

	for j = 0; j < len(key); j++ {
		br := int(key[j])
		if node.Children[br] == nil {
			l := len(node.Branches)
			if l > 0 && node.Branches[l-1] > br {
				err = errors.Wrapf(ErrKeyOutOfOrder, "append %q", key)
				return
			}
			break
		}
		node = node.Children[br]
	}

	if j == len(key) {
		leaf = node.Children[leafBranch]
		if leaf != nil {
			err = ErrDuplicateKeys
			return
		}

		if len(node.Branches) != 0 {
			// means this key is a prefix of an existed key, so key's adding order is not ascending.
			err = errors.Wrapf(ErrKeyOutOfOrder, "append %q is a prefix", key)
			return
		}
	}

	commonNode := node
	newBranch := int(key[j])

	var ltNode *Node
	numBr := len(commonNode.Branches)
	if numBr > 0 {
		ltNode = commonNode.Children[commonNode.Branches[numBr-1]]
	}

	for _, b := range key[j:] {
		br := int(b)
		n := &Node{Children: make(map[int]*Node), Step: 1}

		node.Children[br] = n
		node.Branches = append(node.Branches, br)
		node = n

		r.NodeCnt++
	}

	leaf = &Node{Value: value, isStartLeaf: isStartLeaf}

	node.Children[leafBranch] = leaf
	node.Branches = append(node.Branches, leafBranch)

	if needSquash {
		if ltNode != nil {
			r.NodeCnt -= ltNode.Squash()
		}

		if r.NodeCnt > MaxNodeCnt {
			delete(commonNode.Children, newBranch)
			commonNode.removeBranch(newBranch)
			err = errors.Wrapf(ErrTooManyTrieNodes, "append %q", key)
			return nil, err
		}
	}

	return
}

// newRangeTrie creates a range-serach Trie.
// TODO explain.
func newRangeTrie() *Node {
	return &Node{Children: make(map[int]*Node), Step: 1}
}

// removeNonboundaryLeaves removes leaf nodes and ancestor nodes that reach only
// these leaf nodes.
func (r *Node) removeNonboundaryLeaves() {
	// remove leaves which are not start leaf

	leaf := r.Children[leafBranch]
	if leaf != nil {

		if !leaf.isStartLeaf {

			delete(r.Children, leafBranch)
			r.removeBranch(leafBranch)
		}
	}

	for k, n := range r.Children {

		if len(n.Branches) == 0 {
			// leaf node do not need to remove any branch
			continue
		}

		n.removeNonboundaryLeaves()

		if len(n.Branches) == 0 {
			delete(r.Children, k)
			r.removeBranch(k)
		}
	}
}

func (r *Node) removeBranch(br int) {

	idx := sort.Search(
		len(r.Branches),
		func(i int) bool {
			return r.Branches[i] >= br
		},
	)

	if idx < len(r.Branches) && r.Branches[idx] == br {
		r.Branches = append(r.Branches[:idx], r.Branches[idx+1:]...)
	}
}
