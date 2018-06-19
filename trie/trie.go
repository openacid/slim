package trie

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var ErrDuplicateKeys = errors.New("keys can not be duplicate")
var ErrValuesNotSlice = errors.New("values must be in a slice")

type Node struct {
	Children map[int]*Node
	Branches []int
	Step     uint16
	Value    interface{}
}

const leafBranch = -1

func New(keys [][]byte, values interface{}) (root *Node, err error) {

	valSlice, ok := toSlice(values)
	if !ok {
		err = ErrValuesNotSlice
		return
	}

	root = &Node{Children: make(map[int]*Node), Step: 1}

	for i, key := range keys {

		var node = root
		var j int

		for j = 0; j < len(key); j++ {
			br := int(key[j])
			if node.Children[br] == nil {
				break
			}
			node = node.Children[br]
		}

		for _, b := range key[j:] {
			br := int(b)
			n := &Node{Children: make(map[int]*Node), Step: 1}

			node.Children[br] = n
			node.Branches = append(node.Branches, br)
			node = n
		}

		if node.Children[leafBranch] != nil {
			err = ErrDuplicateKeys
			return
		}

		leaf := &Node{Value: valSlice[i]}

		node.Children[leafBranch] = leaf
		node.Branches = append(node.Branches, leafBranch)
	}
	return
}

func (r *Node) ToStrings(cc int) []string {

	var line string
	if cc == leafBranch {
		line = fmt.Sprintf("  $(%d):", len(r.Branches))
	} else {
		line = fmt.Sprintf("%03d(%d):", cc, len(r.Branches))
	}

	rst := make([]string, 0, 64)

	if len(r.Branches) > 0 {

		for _, b := range r.Branches {
			subtrie := r.Children[b].ToStrings(b)
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

func (r *Node) Squash() {

	for k, n := range r.Children {
		n.Squash()

		if len(n.Branches) == 1 {
			if n.Branches[0] == leafBranch {
				continue
			}
			child := n.Children[n.Branches[0]]
			child.Step += 1
			r.Children[k] = child
		}
	}
}

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

func toSlice(arg interface{}) (rst []interface{}, ok bool) {
	s := reflect.ValueOf(arg)
	if s.Kind() != reflect.Slice {
		return
	}
	l := s.Len()
	rst = make([]interface{}, l)
	for i := 0; i < l; i++ {
		rst[i] = s.Index(i).Interface()
	}
	ok = true
	return
}
