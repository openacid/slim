// Package tree provides abstract tree operations.
//
// Since 0.1.5
package tree

import (
	"fmt"
	"strings"
)

// Tree defines required functions to convert a tree-like data structure into
// string.
//
// A nil indicates root node.
//
// Since 0.1.5
type Tree interface {

	// Child returns the child of a node
	//
	// Since 0.1.5
	Child(node, label interface{}) interface{}

	// Labels returns the outgoing branch labels from a node.
	//
	// Since 0.1.5
	Labels(node interface{}) []interface{}

	// NodeID returns the a node id string.
	//
	// Since 0.1.5
	NodeID(node interface{}) string

	// LabelInfo returns a string representation of a branch label.
	//
	// Since 0.1.5
	LabelInfo(label interface{}) string

	// NodeInfo returns string describing the node.
	//
	// Since 0.1.5
	NodeInfo(node interface{}) string

	// LeafVal returns the value if the node is leaf.
	//
	// Since 0.1.5
	LeafVal(node interface{}) (interface{}, bool)
}

// String converts a Tree into a multiline string.
//
// A tree node is in form of
//   <income-label>-><node-id>+<step>*<fanout-count>=<value>
// E.g.:
//   000->#000+2*3
//            001->#001+4*2
//                     003->#004+1=0
//                              006->#007+2=1
//                     004->#005+1=2
//                              006->#008+2=3
//            002->#002+3=4
//                     006->#006+2=5
//                              006->#009+2=6
//            003->#003+5=7`[1:]
//
// Since 0.1.5
func String(t Tree) string {
	lines := toStrings(t, nil, nil)
	return strings.Join(lines, "\n")
}

// DepthFirst walks through a tree in a depth-first manner: process children in
// order then their parent.
//
// Since 0.1.5
func DepthFirst(t Tree, np func(t Tree, parent, label, node interface{})) {
	depthFirst(t, nil, nil, t.Child(nil, nil), np)
}

func depthFirst(t Tree, parent, label, node interface{}, np func(t Tree, parent, label, node interface{})) {

	for _, b := range t.Labels(node) {
		child := t.Child(node, b)
		depthFirst(t, node, b, child, np)
	}
	np(t, parent, label, node)
}

// toStrings converts a node and its subtree into multi strings.
//
// Since 0.1.5
func toStrings(t Tree, inbranch, node interface{}) []string {

	line, ind := nodeStr(t, inbranch, node)

	indent := strings.Repeat(" ", ind)

	rst := make([]string, 0, 64)
	rst = append(rst, line)

	for _, b := range t.Labels(node) {
		sub := toStrings(t, b, t.Child(node, b))
		for _, s := range sub {
			rst = append(rst, indent+s)
		}
	}
	return rst
}

// nodeStr returns a string representation of a node, and subtree indent.
//
// Since 0.1.5
func nodeStr(t Tree, inbranch, node interface{}) (string, int) {

	var line string

	if inbranch != nil {
		line += fmt.Sprintf("-%03v->", t.LabelInfo(inbranch))
	}

	nodeid := t.NodeID(node)
	if nodeid != "" {
		line += "#" + nodeid
	}

	indent := len(line)

	line += t.NodeInfo(node)

	brCnt := len(t.Labels(node))
	if brCnt > 1 {
		line += fmt.Sprintf("*%d", brCnt)
	}

	v, isLeaf := t.LeafVal(node)
	if isLeaf {
		line += fmt.Sprintf("=%v", v)
	}
	return line, indent
}
