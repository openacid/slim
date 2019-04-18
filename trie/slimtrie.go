// Package trie provides slimtrie implementation.
//
// A slimtrie is a static, compressed Trie implementation.
// It is created from a standard Trie by removing unnecessary Trie-node.
// And it uses SlimArray to store Trie.
//
// slimtrie memory overhead is about 6 byte per key, or less.
//
// TODO benchmark
// TODO detail explain.
package trie

import (
	"io"
	gobits "math/bits"

	"github.com/openacid/errors"
	"github.com/openacid/slim/array"
	"github.com/openacid/slim/bits"
	"github.com/openacid/slim/encode"
	"github.com/openacid/slim/serialize"
	"github.com/openacid/slim/strhelper"
)

const (
	// WordMask is bit mask for extracting a 4-bit word.
	WordMask = 0xf
	// LeafWord is a special value to indicate a leaf node in a Trie.
	LeafWord = byte(0x10)
	// MaxNodeCnt is the max number of node(including leaf and inner node).
	MaxNodeCnt = 65536
)

// SlimTrie is a space efficient Trie index.
//
// The space overhead is about 6 byte per key and is irrelevant to key length.
//
// It does not store full key information, but only just enough info for
// locating a record.
// That's why an end user must re-validate the record key after reading it from
// other storage.
//
// It stores three parts of information in three SlimArray:
//
// `Children` stores node branches and children position.
// `Steps` stores the number of words to skip between a node and its parent.
// `Leaves` stores user data.
//
// TODO add scenario.
type SlimTrie struct {
	Children array.U32
	Steps    array.U16
	Leaves   array.Array
}

var (
	// ErrTooManyTrieNodes indicates the number of trie nodes(not number of
	// keys) exceeded.
	ErrTooManyTrieNodes = errors.New("compacted trie exceeds max node count=65536")
	// ErrTrieBranchValueOverflow indicate input key consists of a word greater
	// than the max 4-bit word(0x0f).
	ErrTrieBranchValueOverflow = errors.New("compacted trie branch value must <=0x0f")
)

// NewSlimTrie create an empty SlimTrie.
// Argument m implements a encode.Encoder to convert user data to serialized
// bytes and back.
// Leave it to nil if element in values are size fixed type.
//	   int is not of fixed size.
//	   struct { X int64; Y int32; } hax fixed size.
func NewSlimTrie(m encode.Encoder, keys []string, values interface{}) (*SlimTrie, error) {
	st := &SlimTrie{
		Children: array.U32{},
		Steps:    array.U16{},
		Leaves:   array.Array{},
	}
	st.Leaves.EltEncoder = m

	if keys != nil {
		return st, st.load(keys, values)
	}

	return st, nil
}

// load Loads keys and values and builds a SlimTrie.
//
// values must be a slice of data-type of fixed size or compatible with
// SlimTrie.Leaves.Encoder.
func (st *SlimTrie) load(keys []string, values interface{}) (err error) {
	ks := strhelper.SliceToBitWords(keys, 4)
	return st.loadBytes(ks, values)
}

func (st *SlimTrie) loadBytes(keys [][]byte, values interface{}) (err error) {

	trie, err := NewTrie(keys, values, true)
	if err != nil {
		return err
	}

	err = st.LoadTrie(trie)
	return err
}

// LoadTrie compress a standard Trie and store compressed data in it.
func (st *SlimTrie) LoadTrie(root *Node) (err error) {
	if root == nil {
		return
	}

	if root.NodeCnt > MaxNodeCnt {
		return ErrTooManyTrieNodes
	}

	childIndex, childData := []int32{}, []uint32{}
	stepIndex := []int32{}
	stepElts := []uint16{}
	leafIndex, leafData := []int32{}, []interface{}{}

	tq := make([]*Node, 0, 256)
	tq = append(tq, root)

	for nID := int32(0); ; {
		if len(tq) == 0 {
			break
		}

		node := tq[0]
		tq = tq[1:]

		if len(node.Branches) == 0 {
			continue
		}

		brs := node.Branches

		if brs[0] == leafBranch {
			leafIndex = append(leafIndex, nID)
			leafData = append(leafData, node.Children[brs[0]].Value)

			brs = brs[1:]
		}

		if node.Step > 1 {
			stepIndex = append(stepIndex, nID)
			stepElts = append(stepElts, node.Step)

		}

		if len(brs) > 0 {
			childIndex = append(childIndex, nID)
			offset := uint16(nID) + uint16(len(tq)) + uint16(1)

			bitmap := uint16(0)
			for _, b := range brs {
				if b&WordMask != b {
					return ErrTrieBranchValueOverflow
				}
				bitmap |= uint16(1) << (uint16(b) & WordMask)
			}

			childData = append(childData, (uint32(offset)<<16)+uint32(bitmap))
		}

		for _, b := range brs {
			tq = append(tq, node.Children[b])
		}

		nID++
		if nID > MaxNodeCnt {
			return ErrTooManyTrieNodes
		}
	}

	err = st.Children.Init(childIndex, childData)
	if err != nil {
		return err
	}

	err = st.Steps.Init(stepIndex, stepElts)
	if err != nil {
		return err
	}

	err = st.Leaves.Init(leafIndex, leafData)
	if err != nil {
		return errors.Wrapf(err, "failure init leaves")
	}

	return nil
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
func (st *SlimTrie) Search(key string) (lVal, eqVal, rVal interface{}) {
	eqID, lID, rID := int32(0), int32(-1), int32(-1)
	lIsLeaf := false

	// string to 4-bit words
	lenWords := 2 * len(key)

	for i := -1; ; {
		i += int(st.getStep(uint16(eqID)))
		bitmap, child0Id, hasChildren := st.getChild(uint16(eqID))

		if lenWords < i {
			rID = eqID
			eqID = -1
			break
		}

		if lenWords == i {
			if hasChildren {
				rID = int32(child0Id)
			}
			break
		}

		shift := 4 - (i&1)*4
		word := ((key[i>>1] >> uint(shift)) & 0x0f)
		wordBit := uint16(1) << uint(word)
		branchBit := bitmap & wordBit

		// This is a inner node at eqIdx,
		// update eq, and left right closest node

		hasLeaf := st.Leaves.Has(eqID)

		if hasLeaf {
			lID = eqID
			lIsLeaf = true
		}

		// Count of branch on the left
		lCnt := gobits.OnesCount16(bitmap & (wordBit - 1))

		if branchBit != 0 {
			eqID = int32(child0Id) + int32(lCnt)
		} else {
			eqID = -1
		}

		if lCnt > 0 {
			lID = int32(child0Id) + int32(lCnt) - 1
			lIsLeaf = false
		}
		// Might overflow but ok
		if bitmap > (wordBit<<1)-1 {
			rID = int32(child0Id) + int32(lCnt) + int32(branchBit>>uint16(word))
		}

		if branchBit == 0 {
			break
		}
	}

	if lID != -1 {
		if !lIsLeaf {
			lID = int32(st.rightMost(uint16(lID)))
		}
		lVal, _ = st.Leaves.Get(lID)
	}
	if rID != -1 {
		rID = int32(st.leftMost(uint16(rID)))
		rVal, _ = st.Leaves.Get(rID)
	}
	if eqID != -1 {
		eqVal, _ = st.Leaves.Get(eqID)
	}

	return
}

// just return equal value for trie.Search benchmark

// Get the value of the specified key from SlimTrie.
//
// If the key exist in SlimTrie, it returns the correct value.
// If the key does NOT exist in SlimTrie, it could also return some value.
//
// Because SlimTrie is a "index" but not a "kv-map", it does not stores complete
// info of all keys.
// SlimTrie tell you "WHERE IT POSSIBLY BE", rather than "IT IS JUST THERE".
func (st *SlimTrie) Get(key string) (eqVal interface{}, found bool) {

	var word byte
	found = false
	eqIdx := int32(0)

	// string to 4-bit words
	lenWords := 2 * len(key)

	for idx := -1; ; {
		idx += int(st.getStep(uint16(eqIdx)))
		if lenWords < idx {
			eqIdx = -1
			break
		}

		if lenWords == idx {
			break
		}

		// Get a 4-bit word from 8-bit words.
		// Use arithmetic to avoid branch missing.
		shift := 4 - (idx&1)*4
		word = ((key[idx>>1] >> uint(shift)) & 0x0f)

		eqIdx = st.nextBranch(uint16(eqIdx), word)
		if eqIdx == -1 {
			break
		}
	}

	if eqIdx != -1 {
		eqVal, found = st.Leaves.Get(eqIdx)
	}

	return
}

func (st *SlimTrie) getChild(idx uint16) (bitmap uint16, offset uint16, found bool) {
	cval, found := st.Children.Get(int32(idx))
	if found {
		return uint16(cval), uint16(cval >> 16), true
	}
	return 0, 0, false
}

func (st *SlimTrie) getStep(idx uint16) uint16 {
	step, found := st.Steps.Get(int32(idx))
	if found {
		return step
	}
	return uint16(1)
}

// getChildIdx returns the id of the specified child.
// This function does not check if the specified child `offset` exists or not.
func getChildIdx(bm uint16, of uint16, word uint16) uint16 {
	chNum := bits.OnesCount64Before(uint64(bm), uint(word))
	return of + uint16(chNum)
}

func (st *SlimTrie) nextBranch(idx uint16, word byte) int32 {

	bm, of, found := st.getChild(idx)
	if !found {
		return -1
	}

	if (bm >> word & 1) == 1 {
		return int32(getChildIdx(bm, of, uint16(word)))
	}

	return -1
}

func (st *SlimTrie) leftMost(idx uint16) uint16 {
	for {
		if st.Leaves.Has(int32(idx)) {
			return idx
		}

		_, idx, _ = st.getChild(idx)
	}
}

func (st *SlimTrie) rightMost(idx uint16) uint16 {
	for {
		bm, of, found := st.getChild(idx)
		if !found {
			return idx
		}

		// count number of all children
		// TODO use bits.PopCntXX without before.
		chNum := bits.OnesCount64Before(uint64(bm), 64)
		idx = of + uint16(chNum-1)

	}
}

// getMarshalSize returns the serialized length in byte of a SlimTrie.
func (st *SlimTrie) getMarshalSize() int64 {
	cSize := serialize.GetMarshalSize(&st.Children)
	sSize := serialize.GetMarshalSize(&st.Steps)
	lSize := serialize.GetMarshalSize(&st.Leaves)

	return cSize + sSize + lSize
}

// marshal serializes it to byte stream.
func (st *SlimTrie) marshal(writer io.Writer) (cnt int64, err error) {
	var n int64

	if n, err = serialize.Marshal(writer, &st.Children); err != nil {
		return 0, err
	}
	cnt += n

	if n, err = serialize.Marshal(writer, &st.Steps); err != nil {
		return 0, err
	}
	cnt += n

	if n, err = serialize.Marshal(writer, &st.Leaves); err != nil {
		return 0, err
	}
	cnt += n

	return cnt, nil
}

// unmarshal de-serializes and loads SlimTrie from a byte stream.
func (st *SlimTrie) unmarshal(reader io.Reader) error {
	if err := serialize.Unmarshal(reader, &st.Children); err != nil {
		return err
	}

	if err := serialize.Unmarshal(reader, &st.Steps); err != nil {
		return err
	}

	if err := serialize.Unmarshal(reader, &st.Leaves); err != nil {
		return err
	}

	return nil
}
