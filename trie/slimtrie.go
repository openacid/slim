// Package trie provides SlimTrie implementation.
//
// A SlimTrie is a static, compressed Trie data structure.
// It removes unnecessary trie-node(single branch node etc).
// And it internally uses 3 compacted array to store a trie.
//
// SlimTrie memory overhead is about 14 bits per key(without value), or less.
//
// Key value map or key-range value map
//
// SlimTrie is natively something like a key value map.
// Actually besides as a key value map,
// to index a map of key range to value with SlimTrie is also very simple:
//
// Gives a set of key the same value, and use RangeGet() instead of Get().
// SlimTrie does not store branches for adjacent leaves with the same value.
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
	"encoding/binary"
	"fmt"
	"math/bits"
	"reflect"
	"strings"

	"github.com/blang/semver"
	"github.com/openacid/errors"
	"github.com/openacid/low/bitword"
	"github.com/openacid/low/pbcmpl"
	"github.com/openacid/low/tree"
	"github.com/openacid/low/vers"
	"github.com/openacid/slim/array"
	"github.com/openacid/slim/encode"
)

var (
	bw4 = bitword.BitWord[4]
)

const (
	// MaxNodeCnt is the max number of node. Node id in SlimTrie is int32.
	MaxNodeCnt = (1 << 31) - 1
)

// SlimTrie is a space efficient Trie index.
//
// The space overhead is about 14 bits per key and is irrelevant to key length.
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
// Since 0.2.0
type SlimTrie struct {
	Children array.Bitmap16
	Steps    array.U16
	Leaves   array.Array
}

type versionedArray struct {
	*array.Base
}

func (va *versionedArray) GetVersion() string {
	return slimtrieVersion
}

func (st *SlimTrie) GetVersion() string {
	return slimtrieVersion
}

func (st *SlimTrie) compatibleVersions() []string {
	return []string{
		"==1.0.0", // before 0.5.8 it is "1.0.0" for historical reason.
		"==0.5.8",
		"==0.5.9",
		"==" + slimtrieVersion,
	}
}

// NewSlimTrie create an SlimTrie.
// Argument e implements a encode.Encoder to convert user data to serialized
// bytes and back.
// Leave it nil if element in values are size fixed type and you do not really
// care about performance.
//	   int is not of fixed size.
//	   struct { X int64; Y int32; } hax fixed size.
//
// Since 0.2.0
func NewSlimTrie(e encode.Encoder, keys []string, values interface{}) (*SlimTrie, error) {
	return newSlimTrie(e, keys, values)
}

type subset struct {
	keyStart  int
	keyEnd    int
	fromIndex int
}

func newSlimTrie(e encode.Encoder, keys []string, values interface{}) (*SlimTrie, error) {

	n := len(keys)
	if n == 0 {
		return emptySlimTrie(e), nil
	}

	for i := 0; i < len(keys)-1; i++ {
		if keys[i] >= keys[i+1] {
			return nil, errors.Wrapf(ErrKeyOutOfOrder,
				"keys[%d] >= keys[%d] %s %s", i, i+1, keys[i], keys[i+1])
		}
	}

	rvals := checkValues(reflect.ValueOf(values), n)
	tokeep := newValueToKeep(rvals)

	childi := make([]int32, 0, n)
	childv := make([]uint64, 0, n)

	stepi := make([]int32, 0, n)
	stepv := make([]uint16, 0, n)

	leavesi := make([]int32, 0, n)
	leavesv := make([]interface{}, 0, n)

	queue := make([]subset, 0, n*2)
	queue = append(queue, subset{0, n, 0})

	for i := 0; i < len(queue); i++ {
		nid := int32(i)
		o := queue[i]
		s, e := o.keyStart, o.keyEnd

		// single key, it is a leaf
		if e-s == 1 {
			if tokeep[s] {
				leavesi = append(leavesi, nid)
				leavesv = append(leavesv, getV(rvals, s))
			}
			continue
		}

		// need to create an inner node

		prefI := prefixIndex(keys[s:e], o.fromIndex)

		// the first key is a prefix of all other keys, which makes it a leaf.
		isFirstKeyALeaf := len(keys[s])*8/4 == prefI
		if isFirstKeyALeaf {
			if tokeep[s] {
				leavesi = append(leavesi, nid)
				leavesv = append(leavesv, getV(rvals, s))
			}
			s += 1
		}

		// create inner node from following keys

		labels, labelBitmap := getLabels(keys[s:e], prefI, tokeep[s:e])

		hasChildren := len(labels) > 0

		if hasChildren {
			childi = append(childi, nid)
			childv = append(childv, uint64(labelBitmap))

			// put keys with the same starting word to queue.

			for _, label := range labels {

				// Find the first key starting with label
				for ; s < e; s++ {
					word := bw4.Get(keys[s], prefI)
					if word == label {
						break
					}
				}

				// Continue looking for the first key not starting with label
				var j int
				for j = s + 1; j < e; j++ {
					word := bw4.Get(keys[j], prefI)
					if word != label {
						break
					}
				}

				p := subset{
					keyStart:  s,
					keyEnd:    j,
					fromIndex: prefI + 1, // skip the label word
				}
				queue = append(queue, p)
				s = j
			}

			// Exclude the label word at parent node
			step := (prefI - o.fromIndex)
			if step > 0xffff {
				panic(fmt.Sprintf("step=%d is too large. must < 2^16", step))
			}

			// By default to move 1 step forward, thus no need to store 1
			hasStep := step > 0
			if hasStep {
				stepi = append(stepi, nid)
				stepv = append(stepv, uint16(step))
			}
		}
	}

	nodeCnt := int32(len(queue))

	ch, err := array.NewBitmap16(childi, childv, 16)
	if err != nil {
		return nil, err
	}

	steps, err := array.NewU16(stepi, stepv)
	if err != nil {
		return nil, err
	}

	leaves := array.Array{}
	leaves.EltEncoder = e

	err = leaves.Init(leavesi, leavesv)
	if err != nil {
		return nil, errors.Wrapf(err, "failure init leaves")
	}

	// Avoid panic of slice index out of bound.
	ch.ExtendIndex(nodeCnt)
	steps.ExtendIndex(nodeCnt)
	leaves.ExtendIndex(nodeCnt)

	st := &SlimTrie{
		Children: *ch,
		Steps:    *steps,
		Leaves:   leaves,
	}
	return st, nil
}

func checkValues(rvals reflect.Value, n int) reflect.Value {

	if rvals.Kind() != reflect.Slice {
		panic("values is not a slice")
	}

	valn := rvals.Len()

	if n != valn {
		panic(fmt.Sprintf("len(keys) != len(values): %d, %d", n, valn))
	}
	return rvals

}

// newValueToKeep creates a slice indicating which key to keep.
// Value of key[i+1] with the same value with key[i] do not need to keep.
func newValueToKeep(rvals reflect.Value) []bool {

	n := rvals.Len()

	tokeep := make([]bool, n)
	tokeep[0] = true

	for i := 0; i < n-1; i++ {
		tokeep[i+1] = getV(rvals, i+1) != getV(rvals, i)
	}
	return tokeep
}

func getV(reflectSlice reflect.Value, i int) interface{} {
	return reflectSlice.Index(i).Interface()
}

func emptySlimTrie(e encode.Encoder) *SlimTrie {
	st := &SlimTrie{}
	st.Leaves.EltEncoder = e
	return st
}

func prefixIndex(keys []string, from int) int {
	if len(keys) == 1 {
		return len(keys[0])
	}

	n := len(keys)

	end := bw4.FirstDiff(keys[0], keys[n-1], from, -1)
	return end
}

func getLabels(keys []string, from int, tokeep []bool) ([]byte, uint16) {
	labels := make([]byte, 0, 1<<4)
	bitmap := uint16(0)

	for i, k := range keys {

		if !tokeep[i] {
			continue
		}

		word := bw4.Get(k, from)
		b := uint16(1) << word
		if bitmap&b == 0 {
			labels = append(labels, word)
			bitmap |= b
		}

	}
	return labels, bitmap
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
		v, found := st.Leaves.Get(eqID)
		if found {
			return v, found
		}

		// else: maybe matched at a inner node.
	}

	// key is smaller than any range-start or range-end.
	if lID == -1 {
		return nil, false
	}

	// Preceding value is the start of this range.
	// It might be a false-positive

	lVal, _ := st.Leaves.Get(lID)
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
		bitmap, child0Id, hasChildren := st.getChild(eqID)
		if !hasChildren {
			break
		}

		i += int(st.getStep(eqID))

		if lenWords < i {
			rID = eqID
			eqID = -1
			break
		}

		if lenWords == i {
			rID = child0Id
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
			eqID = child0Id + int32(lCnt)
		} else {
			eqID = -1
		}

		if lCnt > 0 {
			lID = child0Id + int32(lCnt) - 1
			lIsLeaf = false
		}
		// Might overflow but ok
		if bitmap > (wordBit<<1)-1 {
			rID = child0Id + int32(lCnt) + int32(branchBit>>uint16(word))
		}

		if branchBit == 0 {
			break
		}
	}

	if lID != -1 {
		if !lIsLeaf {
			lID = st.rightMost(lID)
		}
	}
	if rID != -1 {
		rID = st.leftMost(rID)
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
	eqID := int32(0)

	// string to 4-bit words
	lenWords := 2 * len(key)

	for idx := -1; ; {

		bm, rank, hasInner := st.Children.GetWithRank(eqID)
		if !hasInner {
			// maybe a leaf
			break
		}

		step, _ := st.Steps.Get(eqID)
		idx += 1 + int(step)

		if lenWords < idx {
			eqID = -1
			break
		}

		if lenWords == idx {
			break
		}

		// Get a 4-bit word from 8-bit words.
		// Use arithmetic to avoid branch missing.
		shift := 4 - (idx&1)*4
		word = ((key[idx>>1] >> uint(shift)) & 0x0f)

		bb := uint64(1) << word
		if bm&bb != 0 {
			chNum := bits.OnesCount64(bm & (bb - 1))
			eqID = rank + 1 + int32(chNum)
		} else {
			eqID = -1
			break
		}
	}

	if eqID != -1 {
		eqVal, found = st.Leaves.Get(eqID)
	}

	return
}

func (st *SlimTrie) getChild(idx int32) (bitmap uint16, offset int32, found bool) {
	bm, rank, found := st.Children.GetWithRank(idx)
	if found {
		return uint16(bm), rank + 1, true
	}
	return 0, 0, false
}

func (st *SlimTrie) getStep(idx int32) uint16 {
	step, _ := st.Steps.Get(idx)
	// 1 for the label word at parent node.
	return step + 1
}

func (st *SlimTrie) leftMost(idx int32) int32 {
	for {
		if st.Leaves.Has(idx) {
			return idx
		}

		_, idx, _ = st.getChild(idx)
	}
}

func (st *SlimTrie) rightMost(idx int32) int32 {
	for {
		bm, of, found := st.getChild(idx)
		if !found {
			return idx
		}

		// count number of all children
		chNum := bits.OnesCount16(bm)
		idx = of + int32(chNum-1)

	}
}

// Marshal serializes it to byte stream.
//
// Since 0.4.3
func (st *SlimTrie) Marshal() ([]byte, error) {
	var buf []byte
	writer := bytes.NewBuffer(buf)

	_, err := pbcmpl.Marshal(writer, &versionedArray{&st.Children.Base})
	if err != nil {
		return nil, errors.WithMessage(err, "failed to marshal children")
	}

	_, err = pbcmpl.Marshal(writer, &versionedArray{&st.Steps.Base})
	if err != nil {
		return nil, errors.WithMessage(err, "failed to marshal steps")
	}

	_, err = pbcmpl.Marshal(writer, &versionedArray{&st.Leaves.Base})
	if err != nil {
		return nil, errors.WithMessage(err, "failed to marshal leaves")
	}

	return writer.Bytes(), nil
}

// Unmarshal de-serializes and loads SlimTrie from a byte stream.
//
// Since 0.4.3
func (st *SlimTrie) Unmarshal(buf []byte) error {

	var ver string
	compatible := st.compatibleVersions()
	reader := bytes.NewReader(buf)

	_, ver, err := pbcmpl.Unmarshal(reader, &st.Children)
	if err != nil {
		return errors.WithMessage(err, "failed to unmarshal children")
	}

	if !vers.IsCompatible(ver, compatible) {
		return errors.Wrapf(ErrIncompatible,
			fmt.Sprintf(`version: "%s", compatible versions:"%s"`,
				ver,
				strings.Join(compatible, " || ")))
	}

	_, _, err = pbcmpl.Unmarshal(reader, &st.Steps)
	if err != nil {
		return errors.WithMessage(err, "failed to unmarshal steps")
	}

	_, _, err = pbcmpl.Unmarshal(reader, &st.Leaves)
	if err != nil {
		return errors.WithMessage(err, "failed to unmarshal leaves")
	}

	// backward compatible:

	before058ConvertToChildrenEltsToBMElts(st, ver)
	before059ExtendBitmapIndex(st, ver)
	before0510StepFrom0(st, ver)

	return nil
}

func checkVer(ver string, want string) bool {
	v, err := semver.Parse(ver)
	if err != nil {
		panic(err)
	}

	chk, err := semver.ParseRange(want)
	if err != nil {
		panic(err)
	}

	return chk(v)
}

func before058ConvertToChildrenEltsToBMElts(st *SlimTrie, ver string) {

	// 1.0.0 is the initial version.
	// From 0.5.8 it starts writing version to marshaled data.
	// In 0.5.4 it starts using Bitmap to store Children elements.
	// But 0.5.4 marshals data with version == 1.0.0

	// Convert u32 node to ranked bitmap node

	if checkVer(ver, "==1.0.0 || <0.5.8") {

		// There are two format with version 1.0.0:
		// Before 0.5.4 Children elements are in Elts
		// From 0.5.4 Children elements are in BMElts

		if st.Children.Flags&array.ArrayFlagIsBitmap == 0 {

			var endian = binary.LittleEndian

			b := &st.Children
			b.Flags |= array.ArrayFlagHasEltWidth | array.ArrayFlagIsBitmap
			b.EltWidth = 16

			indexes := b.Indexes()
			elts := make([]uint64, len(indexes))
			for i, idx := range indexes {
				eltIdx, found := b.GetEltIndex(idx)
				if !found {
					panic("not found index???")
				}
				v := endian.Uint32(b.Elts[eltIdx*4:])
				elts[i] = uint64(v & 0xffff)
			}
			b.BMElts = array.NewBitsJoin(elts, b.EltWidth, false).(*array.Bits)
			b.Elts = nil
		}
	}
}

func before059ExtendBitmapIndex(st *SlimTrie, ver string) {

	// From 0.5.9 it create aligned array bitmap.
	// Need to align array bitmap for previous versions.

	if checkVer(ver, "==1.0.0 || <0.5.9") {

		n := 0
		if n < len(st.Children.Bitmaps)*64 {
			n = len(st.Children.Bitmaps) * 64
		}
		if n < len(st.Steps.Bitmaps)*64 {
			n = len(st.Steps.Bitmaps) * 64
		}
		if n < len(st.Leaves.Bitmaps)*64 {
			n = len(st.Leaves.Bitmaps) * 64
		}

		nodeCnt := int32(n)
		st.Children.ExtendIndex(nodeCnt)
		st.Steps.ExtendIndex(nodeCnt)
		st.Leaves.ExtendIndex(nodeCnt)
	}
}

func before0510StepFrom0(st *SlimTrie, ver string) {

	// From 0.5.10 step does not include the count of the label word.

	if checkVer(ver, "==1.0.0 || <0.5.10") {

		ii := st.Steps.Indexes()
		for i, idx := range ii {
			v, found := st.Steps.Get(idx)
			if !found {
				panic("not found??")
			}

			v -= 1

			st.Steps.Elts[i*2] = byte(v)
			st.Steps.Elts[i*2+1] = byte(v >> 8)
		}
	}
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
	s := &slimTrieStringly{st: st}
	return tree.String(s)
}

// ProtoMessage implements proto.Message
//
// Since 0.4.3
func (st *SlimTrie) ProtoMessage() {}
