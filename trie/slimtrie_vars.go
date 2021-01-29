package trie

import "github.com/openacid/low/bitmap"

// slimVars stores several internally used variables by slim, to speed up calculation
// during querying.
//
// Since 0.5.12
type slimVars struct {
	// BigInnerOffset is the offset caused by "BigInner" nodes:
	//
	// Supposing that the i-th inner node is the j-th short inner node(an inner
	// node can be a short).
	//
	// The offset of this node in "Inners" is
	//
	//     257 * BigInnerCnt +
	//     17 * (i-BigInnerCnt-j) +
	//     ShortSize * j
	//
	// Thus we could create 2 variables to reduce offset calculation time:
	//
	//     BigInnerOffset = (257 - 17) * BigInnerCnt
	//     ShortMinusInner = ShortSize - 17
	//
	// The the offset is:
	//
	//     BigInnerOffset + 17 * i + ShortMinusInner * j
	//
	// Since 0.5.12
	BigInnerOffset int32

	// ShortMinusInner is ShortSize minus 17.
	// See BigInnerOffset.
	//
	// Since 0.5.12
	ShortMinusInner int32

	// ShortMask has the lower ShortSize bit set to 1.
	//
	// Since 0.5.12
	ShortMask uint64
}

// initNodeLocatingVars initialize internal st.vars that are related to locating an inner node
//
// Since 0.5.12
func (st *SlimTrie) initNodeLocatingVars() {

	ns := st.inner

	st.vars = &slimVars{
		BigInnerOffset:  (bigInnerSize - innerSize) * ns.BigInnerCnt,
		ShortMinusInner: ns.ShortSize - innerSize,
		ShortMask:       bitmap.Mask[ns.ShortSize],
	}
}
