package trie

import (
	"math/bits"

	"github.com/openacid/low/bitmap"
)

func (st *SlimTrie) getIthInner(ithInner int32, qr *querySession) {
	ns := st.inner
	vars := st.vars

	innWordI := ithInner >> 6
	innBitI := ithInner & 63

	if ithInner < ns.BigInnerCnt {
		qr.wordSize = bigWordSize
		qr.from = ithInner * bigInnerSize
		qr.to = qr.from + bigInnerSize
	} else {
		qr.wordSize = wordSize

		ithShort := ns.ShortBM.RankIndex[innWordI] + int32(bits.OnesCount64(ns.ShortBM.Words[innWordI]&bitmap.Mask[innBitI]))

		qr.from = vars.BigInnerOffset + innerSize*ithInner + vars.ShortMinusInner*ithShort

		// if this is a short node
		if ns.ShortBM.Words[innWordI]&bitmap.Bit[innBitI] != 0 {

			qr.to = qr.from + ns.ShortSize

			j := qr.from & 63
			w := ns.Inners.Words[qr.from>>6]

			var bm uint64

			if j <= 64-ns.ShortSize {
				bm = (w >> uint32(j)) & vars.ShortMask
			} else {
				w2 := ns.Inners.Words[qr.to>>6]
				bm = (w >> uint32(j)) | (w2 << uint(64-j) & vars.ShortMask)
			}

			qr.bm = uint64(ns.ShortTable[bm])

		} else {
			qr.to = qr.from + innerSize
		}
	}
}

// getIthInnerFrom finds out the start position of the label bitmap of the ith inner node(not by node id).
func (st *SlimTrie) getIthInnerFrom(ithInner int32, qr *querySession) {
	ns := st.inner
	vars := st.vars

	if ithInner < ns.BigInnerCnt {
		qr.from = ithInner * bigInnerSize
	} else {
		innWordI := ithInner >> 6

		ithShort := ns.ShortBM.RankIndex[innWordI] + int32(bits.OnesCount64(ns.ShortBM.Words[innWordI]&bitmap.Mask[ithInner&63]))

		qr.from = vars.BigInnerOffset + innerSize*ithInner + vars.ShortMinusInner*ithShort
	}
}

// getIthInnerFromTo fills in qr.from and qr.to
func (st *SlimTrie) getIthInnerFromTo(ithInner int32, qr *querySession) {

	ns := st.inner
	vars := st.vars

	if ithInner < ns.BigInnerCnt {
		qr.from = ithInner * bigInnerSize
		qr.to = qr.from + bigInnerSize
	} else {

		ithShort, isShort := bitmap.Rank64(ns.ShortBM.Words, ns.ShortBM.RankIndex, ithInner)

		qr.from = vars.BigInnerOffset + innerSize*ithInner + vars.ShortMinusInner*ithShort

		// if this is a short node
		if isShort != 0 {
			qr.to = qr.from + ns.ShortSize
		} else {
			qr.to = qr.from + innerSize
		}
	}
}

// getIthInnerChildren returns the first child id and the last child id of the ith-inner node(not node id)
// An node i(first  <= i <= last) is also a child of the ith-inner node.
func (st *SlimTrie) getIthInnerChildren(idx int32) (int32, int32) {
	ns := st.inner
	qr := &querySession{}

	st.getIthInnerFromTo(idx, qr)
	firstChildId, _ := bitmap.Rank128(ns.Inners.Words, ns.Inners.RankIndex, qr.from)
	firstChildId = firstChildId + 1

	lastChildId, b := bitmap.Rank128(ns.Inners.Words, ns.Inners.RankIndex, qr.to-1)
	lastChildId = lastChildId + b

	return firstChildId, lastChildId
}
