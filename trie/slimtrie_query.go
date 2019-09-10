package trie

import (
	"bytes"
	"math/bits"

	"github.com/openacid/low/bitmap"
	"github.com/openacid/low/bmtree"
)

type querySession struct {
	// Inner node bit range
	from, to int32

	// Extracted bitmap for most-used node
	bm uint64

	// The size in bit of a inner node, such as 4-bit or 8-bit.
	wordSize int32

	// Whether current node is an inner node or leaf node.
	isInner bool

	keyBitLen int32
	key       string

	// Whether an inner node has common prefix.
	// It may stores only length of prefix in prefixBitLen, or extact prefix
	// string in prefix.
	hasPrefixContent bool

	// Number of bit of a prefix.
	prefixLen int32

	// Prefix string.
	prefix []byte

	ithLeaf       int32
	hasLeafPrefix bool
	leafPrefix    []byte
}

// Get the value of the specified key from SlimTrie.
//
// If the key exist in SlimTrie, it returns the correct value.
// If the key does NOT exist in SlimTrie, it could also return some value.
//
// Because SlimTrie is a "index" but not a "kv-map", it does not stores complete
// info of all keys.
// SlimTrie tell you "WHERE IT POSSIBLY BE", rather than "IT IS JUST THERE".
//
// Since 0.2.0
func (st *SlimTrie) Get(key string) (interface{}, bool) {

	eqID := st.GetID(key)

	if eqID == -1 {
		return nil, false
	}

	v := st.getLeaf(eqID)
	return v, true
}

// GetI32 is same as Get() except it optimized for int32.
//
// Since 0.5.10
func (st *SlimTrie) GetI32(key string) (int32, bool) {

	// TODO test and doc

	eqID := st.GetID(key)

	if eqID == -1 {
		return 0, false
	}

	ith, _ := st.getLeafIndex(eqID)
	stIdx := ith << 2

	b := st.nodes.LeafBytes[stIdx : stIdx+4]

	v := int32(b[0]) | int32(b[1])<<8 | int32(b[2])<<16 | int32(b[3])<<24

	return v, true
}

// RangeGet look for a range that contains a key in SlimTrie.
//
// A range that contains a key means range-start <= key <= range-end.
//
// It returns the value the range maps to, and a bool indicate if a range is
// found.
//
// A positive return value does not mean the range absolutely exists, which in
// this case, is a "false positive".
//
// Since 0.4.3
func (st *SlimTrie) RangeGet(key string) (interface{}, bool) {

	lID, eqID, _ := st.searchID(key)

	// an "equal" macth means key is a prefix of either start or end of a range.
	if eqID != -1 {
		// TODO eqID must be a leaf if it is not -1
		return st.getLeaf(eqID), true
	}

	// key is smaller than any range-start or range-end.
	if lID == -1 {
		return nil, false
	}

	// Preceding value is the start of this range.
	// It might be a false-positive

	return st.getLeaf(lID), true
}

// Search for a key in SlimTrie.
//
// It returns values of 3 values:
// The value of greatest key < `key`. It is nil if `key` is the smallest.
// The value of `key`. It is nil if there is not a matching.
// The value of smallest key > `key`. It is nil if `key` is the greatest.
//
// A non-nil return value does not mean the `key` exists.
// An in-existent `key` also could matches partial info stored in SlimTrie.
//
// Since 0.2.0
func (st *SlimTrie) Search(key string) (lVal, eqVal, rVal interface{}) {

	lID, eqID, rID := st.searchID(key)

	if lID != -1 {
		lVal = st.getLeaf(lID)
	}
	if eqID != -1 {
		eqVal = st.getLeaf(eqID)
	}
	if rID != -1 {
		rVal = st.getLeaf(rID)
	}

	return
}

// GetID looks up for key and return the node id.
// It should only be used to create a user-defined, type specific SlimTrie.
//
// Since 0.5.10
func (st *SlimTrie) GetID(key string) int32 {

	eqID := int32(0)

	if st.nodes.NodeTypeBM == nil {
		return -1
	}

	l := int32(8 * len(key))
	qr := &querySession{
		keyBitLen: l,
		key:       key,
	}

	i := int32(0)

	for {

		qr.isInner = false
		qr.prefixLen = 0
		qr.hasPrefixContent = false

		st.getInner(eqID, qr)
		if !qr.isInner {
			// leaf
			break
		}

		if qr.hasPrefixContent {
			r := prefixCompare(key[i>>3:], qr.prefix)
			if r != 0 {
				return -1
			}
			i = i&(^7) + qr.prefixLen
		} else {
			i += qr.prefixLen
		}

		if i > l {
			return -1
		}

		lchID, has := st.getLEChildID(qr, i)
		if has == 0 {
			// no such branch of label
			return -1
		}
		eqID = lchID + 1

		if i == l {
			// must be a leaf
			break
		}

		i += qr.wordSize
	}

	// eqID must not be -1

	if st.nodes.WithLeafPrefix {
		if i == l {
			if qr.hasLeafPrefix {
				return -1
			} else {
				return eqID
			}
		} else {
			if !qr.hasLeafPrefix {
				return -1
			} else {
				if !bytes.Equal(qr.leafPrefix, []byte(key[i>>3:])) {
					return -1
				}
			}
		}
	}

	return eqID
}

func (st *SlimTrie) cmpLeafPrefix(tail string, qr *querySession) int32 {

	if st.nodes.WithLeafPrefix {
		var leafPrefix []byte
		if qr.hasLeafPrefix {
			leafPrefix = qr.leafPrefix
		} else {
			leafPrefix = []byte{}
		}
		return int32(bytes.Compare([]byte(tail), leafPrefix))
	}

	return 0
}

// searchID searches for key and returns 3 leaf node id:
//
// The id of greatest key < `key`. It is -1 if `key` is the smallest.
// The id of `key`. It is -1 if there is not a matching.
// The id of smallest key > `key`. It is -1 if `key` is the greatest.
func (st *SlimTrie) searchID(key string) (lID, eqID, rID int32) {

	if st.nodes.NodeTypeBM == nil {
		return -1, -1, -1
	}

	lID, eqID, rID = -1, 0, -1
	l := int32(8 * len(key))
	ns := st.nodes

	qr := &querySession{
		keyBitLen: l,
		key:       key,
	}

	i := int32(0)

	for {

		qr.isInner = false
		qr.prefixLen = 0
		qr.hasPrefixContent = false

		st.getInner(eqID, qr)
		if !qr.isInner {
			// leaf
			break
		}

		if qr.hasPrefixContent {
			r := prefixCompare(key[i>>3:], qr.prefix)
			if r == 0 {
				i = i&(^7) + qr.prefixLen
			} else if r < 0 {
				rID = eqID
				eqID = -1
				break
			} else {
				lID = eqID
				eqID = -1
				break
			}

		} else {
			i += qr.prefixLen
			if i > l {
				rID = eqID
				eqID = -1
				break
			}
		}

		lchID, has := st.getLEChildID(qr, i)
		chID := lchID + has
		rchID := chID + 1

		chll, _ := bitmap.Rank128(ns.Inners.Words, ns.Inners.RankIndex, qr.from)
		chll++
		chrr, bit := bitmap.Rank128(ns.Inners.Words, ns.Inners.RankIndex, qr.to-1)
		chrr += bit

		if lchID >= chll && lchID <= chrr {
			lID = lchID
		}
		if rchID >= chll && rchID <= chrr {
			rID = rchID
		}

		if has == 0 {
			eqID = -1
			break
		}
		eqID = chID

		if i == l {
			// must be a leaf
			break
		}

		i += qr.wordSize
	}

	if eqID != -1 {
		tail := key[i>>3:]
		r := st.cmpLeafPrefix(tail, qr)
		if r == -1 {
			rID = eqID
			eqID = -1
		} else if r == 1 {
			lID = eqID
			eqID = -1
		}
	}

	if lID != -1 {
		lID = st.rightMost(lID)
	}
	if rID != -1 {
		rID = st.leftMost(rID)
	}

	return
}

func (st *SlimTrie) leftMost(idx int32) int32 {

	ns := st.nodes

	for {

		qr := &querySession{}

		st.getInner(idx, qr)
		if !qr.isInner {
			break
		}

		r0, _ := bitmap.Rank128(ns.Inners.Words, ns.Inners.RankIndex, qr.from)
		idx = r0 + 1
	}
	return idx
}

func (st *SlimTrie) rightMost(idx int32) int32 {

	ns := st.nodes

	for {
		qr := &querySession{}
		st.getInner(idx, qr)
		if !qr.isInner {
			break
		}

		r0, bit := bitmap.Rank128(ns.Inners.Words, ns.Inners.RankIndex, qr.to-1)
		idx = r0 + bit
		// index out of range with this:
		// r0, _ := bitmap.Rank128(ns.Inners.Words, ns.Inners.RankIndex, n.to)
		// idx = r0
	}
	return idx
}

func (st *SlimTrie) getLeafPrefix(nodeid int32, qr *querySession) {

	qr.ithLeaf, _ = st.getLeafIndex(nodeid)

	qr.hasLeafPrefix = false

	if st.nodes.WithLeafPrefix {

		ns := st.nodes

		wordI := qr.ithLeaf >> 6
		bitI := uint32(qr.ithLeaf & 63)

		if ns.LeafPrefixBM.Words[wordI]&bitmap.Bit[bitI] != 0 {
			ithPref := ns.LeafPrefixBM.RankIndex[wordI] + int32(bits.OnesCount64(ns.LeafPrefixBM.Words[wordI]&bitmap.Mask[bitI]))
			ps := ns.LeafPrefixStartBM
			from, to := ps.select32(ithPref)

			qr.hasLeafPrefix = true
			qr.leafPrefix = ns.LeafPrefixes[from:to]

		}
	}
}

func (st *SlimTrie) getInner(nodeid int32, qr *querySession) {

	var bm uint64

	ns := st.nodes

	wordI := nodeid >> 6
	bitI := uint32(nodeid & 63)

	if ns.NodeTypeBM.Words[wordI]&bitmap.Bit[bitI] == 0 {
		st.getLeafPrefix(nodeid, qr)
		return
	}
	qr.isInner = true

	ithInner := ns.NodeTypeBM.RankIndex[wordI] + int32(bits.OnesCount64(ns.NodeTypeBM.Words[wordI]&bitmap.Mask[bitI]))

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

	// if this node has prefix
	// TODO no prefix mode when create
	if ns.InnerPrefixCnt > 0 && ns.InnerPrefixBM.Words[innWordI]&bitmap.Bit[innBitI] != 0 {

		ithPref := rank128(ns.InnerPrefixBM.Words, ns.InnerPrefixBM.RankIndex, ithInner)

		if ns.WithPrefixContent {

			// stored actual prefix of a node.
			ps := ns.InnerPrefixStartBM
			from, to := ps.select32(ithPref)

			qr.prefix = ns.InnerPrefixes[from:to]
			qr.prefixLen = prefixLen(qr.prefix)
			qr.hasPrefixContent = true

		} else {
			qr.prefixLen = decStep(ns.InnerPrefixLens[ithPref<<1:])
		}
	}
}

func (st *SlimTrie) getLEChildID(qr *querySession, ki int32) (int32, int32) {

	ns := st.nodes

	ithBit := int32(0)

	// if ki > n.keyBitLen {
	//     panic("xx")
	// }

	if ki < qr.keyBitLen {

		if qr.wordSize == bigWordSize {
			ithBit = 1 + int32(qr.key[ki>>3])
		} else {

			b := qr.key[ki>>3]

			if ki&7 < 4 {
				b >>= 4
			}
			b &= 0xf

			ithBit = 1 + int32(b)
		}
	}

	if qr.to-qr.from == ns.ShortSize {

		r0 := rank128(ns.Inners.Words, ns.Inners.RankIndex, qr.from)
		r0 += int32(bits.OnesCount64(qr.bm & bitmap.Mask[ithBit]))
		return r0, int32(qr.bm >> uint(ithBit) & 1)

	} else {
		return bitmap.Rank128(ns.Inners.Words, ns.Inners.RankIndex, qr.from+ithBit)
	}

}

func rank128(words []uint64, rindex []int32, i int32) int32 {

	wordI := i >> 6
	j := uint32(i & 63)
	atRight := wordI & 1

	n := rindex[(i+64)>>7]
	w := words[wordI]

	cnt1 := int32(bits.OnesCount64(w))
	return n - atRight*cnt1 + int32(bits.OnesCount64(w&bitmap.Mask[j]))
}

// the second return value being 0 indicates it is a leaf
func (st *SlimTrie) getLeafIndex(nodeid int32) (int32, int32) {
	ns := st.nodes
	r, ith := bitmap.Rank64(ns.NodeTypeBM.Words, ns.NodeTypeBM.RankIndex, nodeid)
	return nodeid - r, ith
}

func (st *SlimTrie) getLeaf(nodeid int32) interface{} {
	leafI, nodeType := st.getLeafIndex(nodeid)
	if nodeType == 1 {
		panic("impossible!!")
	}

	return st.getIthLeaf(leafI)
}

func (st *SlimTrie) getIthLeaf(ith int32) interface{} {

	if !st.nodes.WithLeaves {
		return nil
	}

	eltsize := st.encoder.GetEncodedSize(nil)
	stIdx := ith * int32(eltsize)

	bs := st.nodes.LeafBytes[stIdx:]

	_, v := st.encoder.Decode(bs)
	return v
}

func (st *SlimTrie) getLabels(qr *querySession) []uint64 {

	ns := st.nodes

	if qr.to-qr.from == ns.ShortSize {
		return bmtree.Decode(innerSize, []uint64{qr.bm})
	}

	bm := bitmap.Slice(ns.Inners.Words, qr.from, qr.to)
	return bmtree.Decode(qr.to-qr.from, bm)
}
