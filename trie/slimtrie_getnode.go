package trie

import (
	"math/bits"

	"github.com/openacid/low/bitmap"
)

func (st *SlimTrie) getIthInner(ithInner int32, qr *querySession) {
	ns := st.inner

	innWordI := ithInner >> 6
	innBitI := ithInner & 63

	if ithInner < ns.BigInnerCnt {
		qr.wordSize = bigWordSize
		qr.from = ithInner * bigInnerSize
		qr.to = qr.from + bigInnerSize
	} else {
		qr.wordSize = wordSize

		ithShort := ns.ShortBM.RankIndex[innWordI] + int32(bits.OnesCount64(ns.ShortBM.Words[innWordI]&bitmap.Mask[innBitI]))

		qr.from = ns.BigInnerOffset + innerSize*ithInner + ns.ShortMinusInner*ithShort

		// if this is a short node
		if ns.ShortBM.Words[innWordI]&bitmap.Bit[innBitI] != 0 {

			qr.to = qr.from + ns.ShortSize

			j := qr.from & 63
			w := ns.Inners.Words[qr.from>>6]

			var bm uint64

			if j <= 64-ns.ShortSize {
				bm = (w >> uint32(j)) & ns.ShortMask
			} else {
				w2 := ns.Inners.Words[qr.to>>6]
				bm = (w >> uint32(j)) | (w2 << uint(64-j) & ns.ShortMask)
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

	if ithInner < ns.BigInnerCnt {
		qr.from = ithInner * bigInnerSize
	} else {
		innWordI := ithInner >> 6

		ithShort := ns.ShortBM.RankIndex[innWordI] + int32(bits.OnesCount64(ns.ShortBM.Words[innWordI]&bitmap.Mask[ithInner&63]))

		qr.from = ns.BigInnerOffset + innerSize*ithInner + ns.ShortMinusInner*ithShort
	}
}
