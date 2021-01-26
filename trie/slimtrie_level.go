package trie

import (
	"fmt"

	"github.com/openacid/low/bitmap"
)

// levelInfo records node count upto every level(inclusive).
// Slim has a slice []levelInfo to track node counts.
// These nodes count info helps to speed up in finding out the original position
// of a key in the creating key value array, with the help with another bitmap
// NodeTypeBM, in which a "1" indicates an inner node.
//
// The 0-th elt is always 0,0,0. An empty slim has only one level.
//
// The 1th elt describe root node in slim.
// The 1th elt is:
// - 1,0,1 if the slim has only one(must be a leaf) node.
// - 1,1,0 if the slim has more than one nodes.
//
//
// With a slim in the following structure:
//
//   node-0(root)
//        +--a--> node-1
//        |            +--x--> node-4(leaf)
//        |            `--y--> node-5(leaf)
//        |
//        +--b--> node-2(leaf)
//        |
//        `--c--> node-3
//                     +--u--> node-6(leaf)
//                     `--v--> node-7(leaf)
//
// The node count of every level is:
//
//   [[0, 0, 0],
//    [1, 1, 0],
//    [4, 3, 1],
//    [8, 2, 5],
//   ]
//
// The NodeTypeBM is:
//   1 101 0000
//
// To find out the position of node-2(nodeId=2, level=2):
//
//   rank0(NodeTypeBM, nodeId=2) - levels[2].leaf // the count of leaf at level 2.
//  +rank0(NodeTypeBM, nodeId=6) - levels[3].leaf // the count of leaf at level 3.
//
// E.g., at every level, count the leaf nodes and sum them.
// When reaching a leaf, find the next inner node at this level(in our case
// node-3) and walks to its first child.
//
// Since 0.5.12
type levelInfo struct {
	// total number of nodes
	// number of inner nodes
	// number of leaf nodes
	// 	total = inner + leaf
	total, inner, leaf int32

	cache []innerCache
}

type innerCache struct {
	nodeId int32
	// N.O. leaves from this node or from preceding node at the same level
	leafCount int32
}

// levelStr builds a slice of string for every level in form of:
//	<i>: <total> = <inner> + <leaf>  <total'> = <inner'> + <leaf'>
//
// Since 0.5.12
func levelsStr(l []levelInfo) []string {
	lineFmt := "%2d: %8d =%8d + %-8d  %8d =%8d + %-8d"

	rst := make([]string, 0, len(l))
	rst = append(rst, " 0:    total =   inner + leaf         total'=  inner' + leaf'")

	for i := 1; i < len(l); i++ {
		ll := l[i]
		prev := l[i-1]
		rst = append(rst, fmt.Sprintf(lineFmt,
			i,
			ll.total,
			ll.inner,
			ll.leaf,
			ll.total-prev.total,
			ll.inner-prev.inner,
			ll.leaf-prev.leaf,
		))
	}

	return rst
}

// initLevels builds the levelInfo slice.
//
// Since 0.5.12
func (st *SlimTrie) initLevels() {
	ns := st.inner
	ntyps := ns.NodeTypeBM

	if ntyps == nil {
		st.levels = []levelInfo{{0, 0, 0, nil}}
		return
	}

	st.levels = make([]levelInfo, 0)

	totalInner, b := bitmap.Rank64(ntyps.Words, ntyps.RankIndex, int32(len(ntyps.Words)*64-1))
	totalInner += b

	// single leaf slim
	total := int32(1)
	if totalInner > 0 {
		var b int32
		total, b = bitmap.Rank128(ns.Inners.Words, ns.Inners.RankIndex, int32(len(ns.Inners.Words)*64-1))
		total += b + 1
	}

	// From root node, walks to the first node at next level, until there is no
	// inner node at next level.

	currId := int32(0)

	qr := &querySession{}
	for {
		// currId is the first node id at current level

		nextInnerIdx, _ := bitmap.Rank64(ntyps.Words, ntyps.RankIndex, currId)

		// update prev level
		st.levels = append(st.levels, levelInfo{total: currId, inner: nextInnerIdx, leaf: currId - nextInnerIdx})

		if nextInnerIdx == totalInner {
			// no more inner node at this level, this is the bottom level
			break
		}

		st.getIthInnerFrom(nextInnerIdx, qr)

		leftMostChild, _ := bitmap.Rank128(ns.Inners.Words, ns.Inners.RankIndex, qr.from)
		currId = leftMostChild + 1
	}

	st.levels = append(st.levels, levelInfo{total: total, inner: totalInner, leaf: total - totalInner})
}
