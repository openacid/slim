package trie

import (
	"math/bits"

	"github.com/openacid/low/bitmap"
)

// newVLenArray builds a VLenArray from a slice of []byte.
//
// It returns nil if no need to build at all, i.e., all elements are empty.
func newVLenArray(elts [][]byte) *VLenArray {
	allEqual := true
	prevSize := -1
	totalSize := 0
	sizes := make([]int32, 0, len(elts))
	nonEmptyIndexes := make([]int32, 0, len(elts))

	for i, elt := range elts {
		eltSize := len(elt)
		sizes = append(sizes, int32(eltSize))
		totalSize += eltSize

		if eltSize > 0 {
			nonEmptyIndexes = append(nonEmptyIndexes, int32(i))
			if prevSize != -1 && prevSize != eltSize {
				allEqual = false
			}
			prevSize = eltSize
		}
	}

	if totalSize == 0 {
		return nil
	}

	// Pack all []byte into one buffer.

	vlenArray := &VLenArray{}
	buf := make([]byte, 0, totalSize)
	for i := 0; i < len(elts); i++ {
		buf = append(buf, elts[i]...)
	}
	vlenArray.Bytes = buf

	vlenArray.N = int32(len(sizes))
	vlenArray.EltCnt = int32(len(nonEmptyIndexes))
	vlenArray.PresenceBM = newBM(nonEmptyIndexes, int32(len(elts)), "r64")

	if allEqual {
		// All non-empty elements are of the same size.
		// Build a fixed size array
		vlenArray.FixedSize = int32(prevSize)
	} else {
		// Build a var-length array
		vlenArray.PositionBM = newBM(stepToPos(sizes, 0), 0, "s32")
	}

	return vlenArray
}

// get returns the `index`-th element.
func (va *VLenArray) get(index int32) []byte {
	if index >= va.N {
		panic("out of bound")
	}
	wordI := index >> 6
	bitI := index & 63

	presence := va.PresenceBM

	if presence.Words[wordI]&bitmap.Bit[bitI] == 0 {
		return []byte{}
	}

	ithElt := presence.RankIndex[wordI] + int32(bits.OnesCount64(presence.Words[wordI]&bitmap.Mask[bitI]))

	positions := va.PositionBM

	if positions == nil {
		// Fixed size elements
		from := ithElt * va.FixedSize
		return va.Bytes[from : from+va.FixedSize]
	}

	// Var-len element

	from, to := bitmap.Select32R64(positions.Words, positions.SelectIndex, positions.RankIndex, ithElt)
	return va.Bytes[from:to]

}
