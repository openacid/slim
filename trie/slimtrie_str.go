package trie

import (
	"fmt"
	"sort"

	"github.com/openacid/low/bitmap"
	"github.com/openacid/low/bmtree"
	"github.com/openacid/low/tree"
)

// String implements proto.Message and output human readable multiline
// representation.
//
// A node is in form of
//   <income-label>-><node-id>+<step>*<fanOut-count>=<value>
// E.g.:
//   000->#000+4*3
//            001->#001+4*2
//                     003->#004+8=0
//                              006->#007+4=1
//                     004->#005+1=2
//                              006->#008+4=3
//            002->#002+4=4
//                     006->#006+4=5
//                              006->#009+8=6
//            003->#003+8=7`[1:]
//
// Since 0.4.3
func (st *SlimTrie) String() string {

	// empty SlimTrie
	if st.inner.NodeTypeBM == nil {
		return ""
	}

	s := &slimTrieStringly{
		st:     st,
		inners: bitmap.ToArray(st.inner.NodeTypeBM.Words),
		labels: make(map[int32]map[string]int32),
	}

	ch := st.inner
	n := &querySession{}
	emp := querySession{}

	for _, nid := range s.inners {
		*n = emp

		s.labels[nid] = make(map[string]int32)

		st.getNode(nid, n)

		paths := st.getLabels(n)

		leftChildId, _ := bitmap.Rank128(ch.Inners.Words, ch.Inners.RankIndex, n.from)

		for i, l := range paths {
			lstr := bmtree.PathStr(l)
			s.labels[nid][lstr] = leftChildId + 1 + int32(i)
		}
	}

	return tree.String(s)
}

// slimTrieStringly is a wrapper that implements tree.Tree .
// It is a helper to convert SlimTrie to string.
//
// Since 0.5.1
type slimTrieStringly struct {
	st     *SlimTrie
	inners []int32
	// node, label(varbits), node
	labels map[int32]map[string]int32
}

// Child implements tree.Tree
//
// Since 0.5.1
func (s *slimTrieStringly) Child(node, branch interface{}) interface{} {

	n := stNodeID(node)
	b := branch.(string)
	return s.labels[n][b]
}

// Labels implements tree.Tree
//
// Since 0.5.1
func (s *slimTrieStringly) Labels(node interface{}) []interface{} {

	n := stNodeID(node)

	rst := make([]string, 0, n)
	labels := s.labels[n]
	for l := range labels {
		rst = append(rst, l)
	}
	sort.Strings(rst)

	r := []interface{}{}
	for _, x := range rst {
		r = append(r, x)
	}
	return r
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
	return label.(string)
}

// NodeInfo implements tree.Tree
//
// Since 0.5.1
func (s *slimTrieStringly) NodeInfo(node interface{}) string {
	nid := stNodeID(node)
	n := &querySession{}
	emp := querySession{}

	if bitmap.Get(s.st.inner.NodeTypeBM.Words, nid) != 0 {

		*n = emp

		s.st.getNode(nid, n)
		step := n.innerPrefixLen
		if step > 0 {
			return fmt.Sprintf("+%d", step)
		}
	}
	return ""
}

// LeafVal implements tree.Tree
//
// Since 0.5.1
func (s *slimTrieStringly) LeafVal(node interface{}) (interface{}, bool) {
	leafI, nodeType := s.st.getLeafIndex(stNodeID(node))
	if nodeType == 1 {
		return nil, false
	}

	v := s.st.getIthLeaf(leafI)
	return v, true
}

// stNodeID convert a interface to SlimTrie node id.
func stNodeID(node interface{}) int32 {
	if node == nil {
		node = int32(0)
	}
	return node.(int32)
}
