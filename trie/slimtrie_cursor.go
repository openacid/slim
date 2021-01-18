package trie

import "github.com/openacid/low/bitmap"

// walkingCursor tracks the state for walking a slim.
type walkingCursor struct {
	// id is the node id that the walking cursor is at
	id int32
	// smallerCnt is N.O. leaf before this level
	smallerCnt int32
	// lvl is current level where id is.
	lvl int32
}

// nextLevel updates states when walking to a child node at next level
func (cur *walkingCursor) nextLevel(ithInner int32, st *SlimTrie, childId int32) {
	cur.smallerCnt += cur.id - ithInner - st.levels[cur.lvl-1].leaf
	cur.lvl++
	cur.id = childId
}

func (st *SlimTrie) rightMostCursor(pos *walkingCursor) {

	ns := st.inner

	for {
		qr := &querySession{}
		st.getInnerTo(pos.id, qr)
		if qr.isInner == 0 {
			return
		}

		r0, bit := bitmap.Rank128(ns.Inners.Words, ns.Inners.RankIndex, qr.to-1)
		if pos != nil {
			pos.nextLevel(qr.ithInner, st, r0+bit)
		}
	}
}

// cursorLeafIndex finds out the position of a leaf node in the original
// key-value array for creating a slimtrie.
func (st *SlimTrie) cursorLeafIndex(cur *walkingCursor, useCache bool) int32 {
	ns := st.inner

	bottom := int32(len(st.levels) - 1)
	qr := &querySession{}

	for {
		nextInnerIdx, _ := bitmap.Rank64(ns.NodeTypeBM.Words, ns.NodeTypeBM.RankIndex, cur.id)

		if nextInnerIdx == st.levels[cur.lvl].inner {
			// Current node is the last inner node at this level.
			// Thus all leaves at higher level are before current leaf.
			cur.nextLevel(nextInnerIdx, st, -1)
			cur.smallerCnt += st.levels[bottom].leaf - st.levels[cur.lvl-1].leaf
			break
		}

		// if there is cached leaf count, use it
		if useCache && st.levels[cur.lvl].cache != nil {

			// in-level index:
			ii := nextInnerIdx - st.levels[cur.lvl-1].inner
			if ii == 0 {
				// all nodes at this level before current node are leaves.
				cur.smallerCnt += cur.id - st.levels[cur.lvl-1].total
				break
			}

			// find the closest inner node before it.
			cache := st.levels[cur.lvl].cache[ii-1]
			// leaves between previous inner and this node.
			cur.smallerCnt += cur.id - 1 - cache.nodeId
			cur.smallerCnt += cache.leafCount
			break
		}

		st.getIthInnerFrom(nextInnerIdx, qr)

		leftMostChild, _ := bitmap.Rank128(ns.Inners.Words, ns.Inners.RankIndex, qr.from)
		cur.nextLevel(nextInnerIdx, st, leftMostChild+1)
	}
	return cur.smallerCnt

}
