package trie

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"unsafe"

	"github.com/openacid/slim/array"
	"github.com/openacid/slim/bit"
	"github.com/openacid/slim/serialize"
)

const (
	// WordMask is bit mask for extracting a 4-bit word.
	WordMask = 0xf
	// LeafWord is a special value to indicate a leaf node in a Trie.
	LeafWord = byte(0x10)
	// MaxNodeCnt is the max number of node(including leaf and inner node).
	MaxNodeCnt = 65536
)

// CompactedTrie is a space efficient Trie index.
//
// The space overhead is about 6 byte per key and is irrelevant to key length.
//
// It does not store full key information, but only just enough info for
// locating a record.
// That's why user must re-validate the record key after reading it from other
// storage.
//
// It stores three parts of information in three CompactedArray:
//
// `Children` stores node branches and children position.
// `Steps` stores the number of words to skip between a node and its parent.
// `Leaves` stores user data.
//
// TODO add scenario.
type CompactedTrie struct {
	Children array.CompactedArray
	Steps    array.CompactedArray
	Leaves   array.CompactedArray
}

type children struct {
	Bitmap uint16
	Offset uint16
}

var (
	// ErrTooManyTrieNodes indicates the number of trie nodes(not number of
	// keys) exceeded.
	ErrTooManyTrieNodes = errors.New("compacted trie exceeds max node count=65536")
	// ErrTrieBranchValueOverflow indicate input key consists of a word greater
	// than the max 4-bit word(0x0f).
	ErrTrieBranchValueOverflow = errors.New("compacted trie branch value must <=0x0f")
)

// childConv implements array.Converter and is the CompactedArray adaptor for
// CompactedTrie.Children .
type childConv struct {
	child *children
}

// Marshal
func (c childConv) Marshal(d interface{}) []byte {
	child := d.(*children)

	b := make([]byte, 4)
	binary.LittleEndian.PutUint16(b[:2], child.Bitmap)
	binary.LittleEndian.PutUint16(b[2:4], child.Offset)

	return b
}

func (c childConv) Unmarshal(b []byte) (int, interface{}) {

	c.child.Bitmap = binary.LittleEndian.Uint16(b[:2])
	c.child.Offset = binary.LittleEndian.Uint16(b[2:4])

	return 4, c.child
}

func (c childConv) GetMarshaledSize(b []byte) int {
	return 4
}

type stepConv struct {
	step *uint16
}

func (c stepConv) Marshal(d interface{}) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, *(d.(*uint16)))
	return b
}

func (c stepConv) Unmarshal(b []byte) (int, interface{}) {
	*c.step = binary.LittleEndian.Uint16(b[:2])
	return 2, c.step
}

func (c stepConv) GetMarshaledSize(b []byte) int {
	return 2
}

// NewCompactedTrie create an empty CompactedTrie.
// Argument c implements a array.Converter to convert user data to serialized
// bytes and back.
func NewCompactedTrie(c array.Converter) *CompactedTrie {
	var step uint16
	ct := &CompactedTrie{
		Children: array.CompactedArray{Converter: childConv{child: &children{}}},
		Steps:    array.CompactedArray{Converter: stepConv{step: &step}},
		Leaves:   array.CompactedArray{Converter: c},
	}

	return ct
}

// Compact compress a standard Trie and store compressed data in it.
func (st *CompactedTrie) Compact(root *Node) (err error) {
	if root == nil {
		return
	}

	childIndex, childData := []uint32{}, []*children{}
	stepIndex, stepData := []uint32{}, []*uint16{}
	leafIndex, leafData := []uint32{}, []interface{}{}

	tq := make([]*Node, 0, 256)
	tq = append(tq, root)

	for nID := uint32(0); ; {
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
			stepData = append(stepData, &node.Step)

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

			ch := &children{
				Bitmap: bitmap,
				Offset: offset,
			}

			childData = append(childData, ch)
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

	err = st.Steps.Init(stepIndex, stepData)
	if err != nil {
		return err
	}

	err = st.Leaves.Init(leafIndex, leafData)
	if err != nil {
		return err
	}

	return nil
}

// Search for a key in CompactedTrie.
//
// `key` is slice of 4-bit word each stored in a byte.
// the higher 4 bit in byte is removed.
//
// It returns values of 3 keys:
// The value of greatest key < `key`. It is nil if `key` is the smallest.
// The value of `key`. It is nil if there is not a matching.
// The value of smallest key > `key`. It is nil if `key` is the greatest.
//
// A non-nil return value does not mean the `key` exists.
// An in-existent `key` also could matches partial info stored in CompactedTrie.
func (st *CompactedTrie) Search(key []byte) (ltVal, eqVal, gtVal interface{}) {
	eqIdx, ltIdx, gtIdx := int32(0), int32(-1), int32(-1)
	ltLeaf := false

	for idx := uint16(0); ; {
		var word byte
		if uint16(len(key)) == idx {
			word = LeafWord
		} else {
			word = (uint8(key[idx]) & WordMask)
		}

		li, ei, ri, leaf := st.neighborBranches(uint16(eqIdx), word)
		if li >= 0 {
			ltIdx = li
			ltLeaf = leaf
		}

		if ri >= 0 {
			gtIdx = ri
		}

		eqIdx = ei
		if eqIdx == -1 {
			break
		}

		if word == LeafWord {
			break
		}

		idx += st.getStep(uint16(eqIdx))

		if idx > uint16(len(key)) {
			gtIdx = eqIdx
			eqIdx = -1
			break
		}
	}

	if ltIdx != -1 {
		if ltLeaf {
			ltVal = st.Leaves.Get(uint32(ltIdx))
		} else {
			rmIdx := st.rightMost(uint16(ltIdx))
			ltVal = st.Leaves.Get(uint32(rmIdx))
		}
	}
	if gtIdx != -1 {
		fmIdx := st.leftMost(uint16(gtIdx))
		gtVal = st.Leaves.Get(uint32(fmIdx))
	}
	if eqIdx != -1 {
		eqVal = st.Leaves.Get(uint32(eqIdx))
	}

	return
}

// SearchString is similar to Search, except it receives a string to search.
// A char in the string is split into 2 4-bit word.
func (st *CompactedTrie) SearchString(key string) (ltVal, eqVal, gtVal interface{}) {
	eqIdx, ltIdx, gtIdx := int32(0), int32(-1), int32(-1)
	ltLeaf := false

	// string to 4-bit words
	lenWords := 2 * uint16(len(key))

	for idx := uint16(0); ; {
		var word byte
		if lenWords == idx {
			word = LeafWord
		} else {
			if idx&uint16(1) == uint16(1) {
				word = (byte(key[idx>>1]) & 0x0f)
			} else {
				word = (byte(key[idx>>1]) & 0xf0) >> 4
			}
		}

		li, ei, ri, leaf := st.neighborBranches(uint16(eqIdx), word)
		if li >= 0 {
			ltIdx = li
			ltLeaf = leaf
		}

		if ri >= 0 {
			gtIdx = ri
		}

		eqIdx = ei
		if eqIdx == -1 {
			break
		}

		if word == LeafWord {
			break
		}

		idx += st.getStep(uint16(eqIdx))

		if idx > lenWords {
			gtIdx = eqIdx
			eqIdx = -1
			break
		}
	}

	if ltIdx != -1 {
		if ltLeaf {
			ltVal = st.Leaves.Get(uint32(ltIdx))
		} else {
			rmIdx := st.rightMost(uint16(ltIdx))
			ltVal = st.Leaves.Get(uint32(rmIdx))
		}
	}
	if gtIdx != -1 {
		fmIdx := st.leftMost(uint16(gtIdx))
		gtVal = st.Leaves.Get(uint32(fmIdx))
	}
	if eqIdx != -1 {
		eqVal = st.Leaves.Get(uint32(eqIdx))
	}

	return
}

// just return equal value for trie.Search benchmark

// SearchStringEqual is similar to SearchString, except it returns only 1 value
// of the matching key, without the left and right value.
func (st *CompactedTrie) SearchStringEqual(key string) (eqVal interface{}) {
	eqIdx := int32(0)

	// string to 4-bit words
	lenWords := 2 * uint16(len(key))

	for idx := uint16(0); ; {
		var word byte
		if lenWords == idx {
			word = LeafWord
		} else {
			if idx&uint16(1) == uint16(1) {
				word = (byte(key[idx>>1]) & 0x0f)
			} else {
				word = (byte(key[idx>>1]) & 0xf0) >> 4
			}
		}

		ei := st.nextBranch(uint16(eqIdx), word)

		eqIdx = ei
		if eqIdx == -1 {
			break
		}

		if word == LeafWord {
			break
		}

		idx += st.getStep(uint16(eqIdx))

		if idx > lenWords {
			eqIdx = -1
			break
		}
	}

	if eqIdx != -1 {
		eqVal = st.Leaves.Get(uint32(eqIdx))
	}

	return
}

func (st *CompactedTrie) getChild(idx uint16) *children {
	if !st.Children.Has(uint32(idx)) {
		return nil
	}

	cval := st.Children.Get(uint32(idx))
	ch := cval.(*children)

	return ch
}

func (st *CompactedTrie) getStep(idx uint16) uint16 {
	step := st.Steps.Get(uint32(idx))
	if step == nil {
		return uint16(1)
	}
	return *(step.(*uint16))
}

func getChildIdx(ch *children, offset uint16) uint16 {
	chNum := bit.PopCnt64Before(uint64(ch.Bitmap), uint32(offset))
	return ch.Offset + uint16(chNum-1)
}

func (st *CompactedTrie) neighborBranches(idx uint16, word byte) (ltIdx, eqIdx, rtIdx int32, ltLeaf bool) {
	ltIdx, eqIdx, rtIdx = int32(-1), int32(-1), int32(-1)
	ltLeaf = false

	isLeaf := st.Leaves.Has(uint32(idx))

	if word == LeafWord {
		if isLeaf {
			eqIdx = int32(idx)
		}
	} else {
		if isLeaf {
			ltIdx = int32(idx)
			ltLeaf = true
		}
	}

	ch := st.getChild(idx)
	if ch == nil {
		return
	}

	if (ch.Bitmap >> word & 1) == 1 {
		eqIdx = int32(getChildIdx(ch, uint16(word+1)))
	}

	ltStart := uint8(word) & WordMask
	for i := int8(ltStart) - 1; i >= 0; i-- {
		if (ch.Bitmap >> uint8(i) & 1) == 1 {
			ltIdx = int32(getChildIdx(ch, uint16(i+1)))
			ltLeaf = false
			break
		}
	}

	rtStart := word + 1
	if word == LeafWord {
		rtStart = uint8(0)
	}

	for i := rtStart; i < LeafWord; i++ {
		if (ch.Bitmap >> i & 1) == 1 {
			rtIdx = int32(getChildIdx(ch, uint16(i+1)))
			break
		}
	}

	return
}

func (st *CompactedTrie) nextBranch(idx uint16, word byte) (eqIdx int32) {
	eqIdx = int32(-1)

	isLeaf := st.Leaves.Has(uint32(idx))

	if word == LeafWord {
		if isLeaf {
			eqIdx = int32(idx)
		}
	}

	ch := st.getChild(idx)
	if ch == nil {
		return
	}

	if (ch.Bitmap >> word & 1) == 1 {
		eqIdx = int32(getChildIdx(ch, uint16(word+1)))
	}

	return
}

func (st *CompactedTrie) leftMost(idx uint16) uint16 {
	for {
		if st.Leaves.Has(uint32(idx)) {
			return idx
		}

		ch := st.getChild(idx)
		idx = ch.Offset
	}
}

func (st *CompactedTrie) rightMost(idx uint16) uint16 {
	offset := uint16(unsafe.Sizeof(uint16(0)) * 8)
	for {
		if !st.Children.Has(uint32(idx)) {
			return idx
		}

		ch := st.getChild(idx)
		idx = getChildIdx(ch, offset)
	}
}

// GetMarshalSize returns the serialized length in byte of a CompactedTrie.
func (ct *CompactedTrie) GetMarshalSize() int64 {
	cSize := serialize.GetMarshalSize(&ct.Children)
	sSize := serialize.GetMarshalSize(&ct.Steps)
	lSize := serialize.GetMarshalSize(&ct.Leaves)

	return cSize + sSize + lSize
}

// Marshal serializes it to byte stream.
func (ct *CompactedTrie) Marshal(writer io.Writer) (cnt int64, err error) {
	var n int64

	if n, err = serialize.Marshal(writer, &ct.Children); err != nil {
		return 0, err
	}
	cnt += n

	if n, err = serialize.Marshal(writer, &ct.Steps); err != nil {
		return 0, err
	}
	cnt += n

	if n, err = serialize.Marshal(writer, &ct.Leaves); err != nil {
		return 0, err
	}
	cnt += n

	return cnt, nil
}

// MarshalAt serializes it to byte stream and write the stream at specified
// offset.
// TODO change to io.WriterAt
func (ct *CompactedTrie) MarshalAt(f *os.File, offset int64) (cnt int64, err error) {

	buf := new(bytes.Buffer)
	if cnt, err = ct.Marshal(buf); err != nil {
		return 0, err
	}

	if _, err = f.WriteAt(buf.Bytes(), offset); err != nil {
		return 0, err
	}

	return cnt, nil
}

// Unmarshal de-serializes and loads CompactedTrie from a byte stream.
func (ct *CompactedTrie) Unmarshal(reader io.Reader) error {
	if err := serialize.Unmarshal(reader, &ct.Children); err != nil {
		return err
	}

	if err := serialize.Unmarshal(reader, &ct.Steps); err != nil {
		return err
	}

	if err := serialize.Unmarshal(reader, &ct.Leaves); err != nil {
		return err
	}

	return nil
}

// Unmarshal de-serializes and loads CompactedTrie from a byte stream at
// specified offset.
// TODO change to io.ReaderAt
func (ct *CompactedTrie) UnmarshalAt(f *os.File, offset int64) (n int64, err error) {
	childrenSize, err := serialize.UnmarshalAt(f, offset, &ct.Children)
	if err != nil {
		return n, err
	}
	offset += childrenSize

	stepsSize, err := serialize.UnmarshalAt(f, offset, &ct.Steps)
	if err != nil {
		return n, err
	}
	offset += stepsSize

	leavesSize, err := serialize.UnmarshalAt(f, offset, &ct.Leaves)
	if err != nil {
		return n, err
	}

	n = childrenSize + stepsSize + leavesSize
	return n, nil
}
