// Package trie provides SlimTrie implementation.
//
// A SlimTrie is a static, compressed Trie implementation.
// It is created from a standard Trie by removing unnecessary Trie-node.
// And it internally uses 3 compacted array to store a Trie.
//
// SlimTrie memory overhead is about 6 bytes per key, or less.
//
// TODO benchmark
// TODO detail explain.
//
// Key value map or key-range value map
//
// SlimTrie is natively something like a key value map.
// Actually besides as a key value map,
// to index a map of key range to value with SlimTrie is also very simple:
//
// Just give two adjacent keys the same value, then SlimTrie
// knows these keys belong to a "range".
// These two keys are left and right boundaries of a range, and are both
// inclusive.
//
//     // a to g --> 1
//     // h      --> 2
//     st, err := NewSlimTrie(encode.Int{}, []string{"a", "g", "h"}, []int{1, 1, 2})
//
//     st.Get("a")      // 1,   true   A normal key-value Get()
//     st.Get("c")      // nil, false  A key-value Get() got nothing.
//     st.RangeGet("c") // 1,   true   A range get got 1
//     st.RangeGet("g") // 1,   true
//     st.RangeGet("h") // 2,   true
//
// See SlimTrie.RangeGet .
//
// False Positive
//
// Just like Bloomfilter, SlimTrie does not contain full information of keys,
// thus there could be a false positive return:
// It returns some value and "true" but the key is not in there.
package trie

import (
	"bytes"
	"fmt"
	"math/bits"
	"strings"

	"github.com/openacid/errors"
	"github.com/openacid/slim/array"
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
//
// Since 0.2.0
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
//
// Since 0.2.0
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
//
// Since 0.2.0
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

		// TODO leaf node does need to record step.
		// But this modification leads to too many changes.
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

	lID, eqID, rID := st.searchID(key)

	// an "equal" macth means key is a prefix of either start or end of a range.
	if eqID != -1 {
		return st.Leaves.Get(eqID)
	}

	// key is smaller than any range-start or range-end.
	if lID == -1 {
		return nil, false
	}

	// key is greater than any range-start or range-end.
	if rID == -1 {
		return nil, false
	}

	lVal, _ := st.Leaves.Get(lID)
	rVal, _ := st.Leaves.Get(rID)

	// If left-value != right-value, the key is between a range-end and next
	// range-start.
	if lVal != rVal {
		return nil, false
	}

	// If range[i].end == range[i+1].start, it is a false positive.
	// SlimTrie can not distinguish this from a positive match.
	//
	// Otherwise, lVal and rVal must be the start and end of a single range.
	return lVal, true
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
		lVal, _ = st.Leaves.Get(lID)
	}
	if eqID != -1 {
		eqVal, _ = st.Leaves.Get(eqID)
	}
	if rID != -1 {
		rVal, _ = st.Leaves.Get(rID)
	}

	return
}

// searchID searches for key and returns 3 leaf node id:
//
// The id of greatest key < `key`. It is -1 if `key` is the smallest.
// The id of `key`. It is -1 if there is not a matching.
// The id of smallest key > `key`. It is -1 if `key` is the greatest.
func (st *SlimTrie) searchID(key string) (lID, eqID, rID int32) {
	lID, eqID, rID = -1, 0, -1
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
		lCnt := bits.OnesCount16(bitmap & (wordBit - 1))

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
	}
	if rID != -1 {
		rID = int32(st.leftMost(uint16(rID)))
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
//
// Since 0.2.0
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

		bm, of, _ := st.getChild(uint16(eqIdx))

		// if not found, bm == 0, thus no need to check if child found.
		if (bm >> word & 1) == 1 {
			chNum := bits.OnesCount16(bm & ((uint16(1) << word) - 1))
			eqIdx = int32(of) + int32(chNum)
		} else {
			eqIdx = -1
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
		chNum := bits.OnesCount16(bm)
		idx = of + uint16(chNum-1)

	}
}

// toStrings convert a Trie to human readalble representation.
func (st *SlimTrie) toStrings(branch byte, id int32) []string {

	bitmap, child0ID, _ := st.getChild(uint16(id))
	step := st.getStep(uint16(id))
	v, isLeaf := st.Leaves.Get(id)

	brCnt := bits.OnesCount16(bitmap)
	line := fmt.Sprintf("%03d->#%03d", branch, id)
	indent := strings.Repeat(" ", len(line))
	line += fmt.Sprintf("+%d", step)
	if brCnt > 1 {
		line += fmt.Sprintf("*%d", brCnt)
	}
	if isLeaf {
		line += fmt.Sprintf("=%v", v)
	}

	rst := make([]string, 0, 64)
	rst = append(rst, line)

	if brCnt > 0 {

		nth := uint16(0)
		for b := byte(0); b < 16; b++ {
			if bitmap&(1<<b) == 0 {
				continue
			}
			childID := child0ID + nth
			sub := st.toStrings(b, int32(childID))
			for _, s := range sub {
				rst = append(rst, indent+s)
			}
			nth++
		}
	}
	return rst
}

// Marshal serializes it to byte stream.
//
// Since 0.4.3
func (st *SlimTrie) Marshal() ([]byte, error) {
	var buf []byte
	writer := bytes.NewBuffer(buf)

	if _, err := serialize.Marshal(writer, &st.Children); err != nil {
		return nil, errors.WithMessage(err, "failed to marshal children")
	}

	if _, err := serialize.Marshal(writer, &st.Steps); err != nil {
		return nil, errors.WithMessage(err, "failed to marshal steps")
	}

	if _, err := serialize.Marshal(writer, &st.Leaves); err != nil {
		return nil, errors.WithMessage(err, "failed to marshal leaves")
	}

	return writer.Bytes(), nil
}

// Unmarshal de-serializes and loads SlimTrie from a byte stream.
//
// Since 0.4.3
func (st *SlimTrie) Unmarshal(buf []byte) error {
	reader := bytes.NewReader(buf)

	if err := serialize.Unmarshal(reader, &st.Children); err != nil {
		return errors.WithMessage(err, "failed to unmarshal children")
	}

	if err := serialize.Unmarshal(reader, &st.Steps); err != nil {
		return errors.WithMessage(err, "failed to unmarshal steps")
	}

	if err := serialize.Unmarshal(reader, &st.Leaves); err != nil {
		return errors.WithMessage(err, "failed to unmarshal leaves")
	}

	return nil
}

// Reset implements proto.Message
//
// Since 0.4.3
func (st *SlimTrie) Reset() {
	st.Children.Array32.Reset()
	st.Steps.Array32.Reset()
	st.Leaves.Array32.Reset()
}

// String implements proto.Message and output human readable multiline
// representation.
//
// A node is in form of
//   <income-label>-><node-id>+<step>*<fanout-count>=<value>
// E.g.:
//   000->#000+2*3
//            001->#001+4*2
//                     003->#004+1=0
//                              006->#007+2=1
//                     004->#005+1=2
//                              006->#008+2=3
//            002->#002+3=4
//                     006->#006+2=5
//                              006->#009+2=6
//            003->#003+5=7`[1:]
//
// Since 0.4.3
func (st *SlimTrie) String() string {
	lines := st.toStrings(0, 0)
	return strings.Join(lines, "\n")
}

// ProtoMessage implements proto.Message
//
// Since 0.4.3
func (st *SlimTrie) ProtoMessage() {}
