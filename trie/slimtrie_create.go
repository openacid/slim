package trie

import (
	"bytes"
	"math/bits"
	"reflect"
	"sort"

	"github.com/openacid/errors"
	"github.com/openacid/low/bitmap"
	"github.com/openacid/low/bitstr"
	"github.com/openacid/low/bmtree"
	"github.com/openacid/low/sigbits"
	"github.com/openacid/must"
	"github.com/openacid/slim/encode"
)

// subset of keys: keys[keyStart:keyEnd].
// fromKeyBit specifies what bits to use.
type subset struct {
	keyStart   int32
	keyEnd     int32
	fromKeyBit int32
	level      int32
}

type creator struct {

	// creator status

	// Whether to create a big inner node or normal inner node for next
	// call to addInner().
	//
	// In SlimTrie the first several inner nodes are big, of size 256+1.
	// The following inner nodes are normal, of size 16+1.
	// Big nodes for speeding up query.
	// Normal nodes for reducing memory cost.
	isBig   bool
	bigCnt  int32
	nodeCnt int32

	leafCnt int32

	// the max level ever seen by creator.
	maxLevel int32

	withLeaves bool

	// options

	option *Opt

	// data

	// len: nr of all nodes

	innerIndexes []int32
	innerBMs     [][]int32
	innerSizes   []int32

	// len: nr of inner nodes

	prefixIndexes  []int32
	prefix4BitLens []uint8
	prefixByteLens []int32
	prefixes       []byte

	// len: nr of leaf nodes = len(nodes) - len(inner_nodes)
	leaves [][]byte

	// leafIndexes[i] is the index of the ith leaf
	leafIndexes []int32

	leafPrefixIndexes []int32
	leafPrefixLens    []int32
	leafPrefixes      []byte
	prefs             map[string]int32

	// stats those affect creating

	innerBMCnt []map[uint64]int32

	clusteredStarts  []uint32
	clusteredOffsets []uint32
	clusteredBytes   []byte
}

func newCreator(n int, withLeaves bool, opt *Opt) *creator {

	c := &creator{

		isBig:   true,
		bigCnt:  0,
		nodeCnt: 0,
		leafCnt: 0,

		maxLevel: 0,

		withLeaves: withLeaves,

		option: opt,

		innerIndexes: make([]int32, 0, n),
		innerBMs:     make([][]int32, 0, n),
		innerSizes:   make([]int32, 0, n),

		prefixIndexes:  make([]int32, 0, n),
		prefix4BitLens: make([]uint8, 0, n),
		prefixByteLens: make([]int32, 0, n),
		prefixes:       make([]byte, 0, n),

		leaves: make([][]byte, 0, n),

		leafIndexes: make([]int32, 0, n),

		leafPrefixIndexes: make([]int32, 0, n),
		leafPrefixLens:    make([]int32, 0, n),
		leafPrefixes:      make([]byte, 0, n),
		prefs:             make(map[string]int32),

		innerBMCnt: make([]map[uint64]int32, maxShortSize+1),

		clusteredStarts:  make([]uint32, 0),
		clusteredOffsets: make([]uint32, 0),
		clusteredBytes:   make([]byte, 0),
	}

	for i := int32(0); i < maxShortSize+1; i++ {
		c.innerBMCnt[i] = make(map[uint64]int32)
	}

	return c
}

// addNode get done the common job of addLeaf and addInner.
func (c *creator) addNode(nid, level int32) {

	must.Be.Equal(c.nodeCnt, nid)

	c.nodeCnt++

	if c.maxLevel < level {
		c.maxLevel = level
	}
}

// addInner adds an inner node of node id "nid".
func (c *creator) addInner(nid, level int32, bmindex []int32, bmsize int32, prefixBitFrom, prefixBitTo int32, key string) {

	c.addNode(nid, level)

	if c.isBig {
		c.bigCnt++
	} else {

		// Only index non-big inner node.

		nbit := int32(len(bmindex))
		if nbit < maxShortSize+1 {
			bm := get17bitmap(bmindex)
			c.innerBMCnt[nbit][bm]++
		}
	}

	c.innerIndexes = append(c.innerIndexes, nid)
	c.innerSizes = append(c.innerSizes, bmsize)
	c.innerBMs = append(c.innerBMs, bmindex)

	c.setPrefix(nid, prefixBitFrom, prefixBitTo, key)
}

// addClusteredInner adds an inner node "nid" with clustered leaves.
func (c *creator) addClusteredInner(nodeId, level int32,
	prefixBitFrom, prefixBitTo int32,
	keys []string,
) {

	c.addNode(nodeId, level)

	// the prefix of this key set from the beginning.
	fullPrefixByteLen := prefixBitTo >> 3

	c.innerIndexes = append(c.innerIndexes, nodeId)
	c.setPrefix(nodeId, prefixBitFrom, prefixBitTo, keys[0])

	c.clusteredStarts = append(c.clusteredStarts,
		uint32(len(c.clusteredOffsets)))

	for _, k := range keys {
		c.clusteredOffsets = append(c.clusteredOffsets, uint32(len(c.clusteredBytes)))
		c.clusteredBytes = append(c.clusteredBytes, k[fullPrefixByteLen:]...)
	}
}

func get17bitmap(bmindex []int32) uint64 {
	if len(bmindex) == 0 {
		return 0
	}
	return bitmap.Of(bmindex)[0]
}

// setPrefix add prefix information to node "nid".
func (c *creator) setPrefix(nid int32, prefixBitFrom, prefixBitTo int32, key string) {

	must.Be.OK(func() {
		must.Be.Equal(c.nodeCnt-1, nid)
	})

	prefixBitLen := prefixBitTo - prefixBitFrom

	if prefixBitLen == 0 {
		return
	}

	if prefixBitLen&3 != 0 {
		panic("step not mod 3")
	}

	c.prefixIndexes = append(c.prefixIndexes, int32(len(c.innerIndexes)-1))

	if *c.option.InnerPrefix {

		prefix := bitstr.New(key, prefixBitFrom, prefixBitFrom+prefixBitLen)

		c.prefixByteLens = append(c.prefixByteLens, int32(len(prefix)))
		c.prefixes = append(c.prefixes, prefix...)

	} else {

		c.prefix4BitLens = append(c.prefix4BitLens, encStep(prefixBitLen)...)
	}
}

func encStep(step int32) []byte {
	step >>= 2
	return []byte{byte(step >> 8), byte(step & 0xff)}
}

func decStep(bs []byte) int32 {
	step := int32(bs[0])<<8 | int32(bs[1])
	return step << 2
}

func (c *creator) setLeafPrefix(nid int32, key string, keyidx int32) {

	must.Be.Equal(c.nodeCnt-1, nid)

	if *c.option.LeafPrefix {

		pref := key[keyidx>>3:]
		if len(pref) > 0 {

			leafCnt := c.nodeCnt - int32(len(c.innerIndexes))

			c.leafPrefixIndexes = append(c.leafPrefixIndexes, leafCnt-1)
			c.leafPrefixLens = append(c.leafPrefixLens, int32(len(pref)))
			c.leafPrefixes = append(c.leafPrefixes, pref...)
			c.prefs[pref]++

			// fmt.Printf("set for node %d %d-th leaf prefix: %q key: %s\n", nid, leafCnt-1, pref, key)
		}
	}
}

// addLeaf adds the content in []byte of a leaf.
func (c *creator) addLeaf(nid, level int32, v []byte) {

	c.addNode(nid, level)

	if c.withLeaves {
		c.leafCnt++
		c.leaves = append(c.leaves, v)
	}
}

// addLeaf adds the index of a leaf.
func (c *creator) addLeafIndex(nid, level int32, idx int32) {

	c.addNode(nid, level)

	if c.withLeaves {
		c.leafCnt++
		c.leafIndexes = append(c.leafIndexes, idx)
	}
}

// counterElt stores an at most 17 bit bitmap and how many times it is used.
type counterElt struct {
	bitmap17 uint64
	cnt      int32
}

func sortedBMCounts(bmCounts []map[uint64]int32) [][]counterElt {

	rst := make([][]counterElt, 0)

	for nbit, bmCnts := range bmCounts {
		_ = nbit

		ss := make([]counterElt, 0, 10)

		for k, v := range bmCnts {
			ss = append(ss, counterElt{k, v})
		}

		sort.Slice(ss, func(i, j int) bool {
			if ss[i].cnt == ss[j].cnt {
				return ss[i].bitmap17 > ss[j].bitmap17
			}
			return ss[i].cnt > ss[j].cnt
		})

		rst = append(rst, ss)
	}
	return rst
}

// memIncrOfShortSize calculates how many bits more memory will be used if
// it uses a ShortSize-bit short node.
// Normally the return value is negative, which means memory usage drops.
func memIncrOfShortSize(sorted [][]counterElt, shortSize int32) (int32, int32) {

	nShort := int32(1) << uint(shortSize)

	// Mem increment by lookup table.
	mem := nShort * 64
	shortCnt := int32(0)

	// how many different n-bit bitmap is already used to reduce memory.
	nbitIth := make([]int32, shortSize+1)

	for short := int32(0); short < 1<<uint(shortSize); short++ {

		nbit := int32(bits.OnesCount64(uint64(short)))

		if nbitIth[nbit] < int32(len(sorted[nbit])) {
			bmcnt := sorted[nbit][nbitIth[nbit]]

			// Mem reduction by replace all 17-bit inner with shortSize-bit short
			red := (innerSize - shortSize) * bmcnt.cnt
			mem -= red
			shortCnt += bmcnt.cnt

			nbitIth[nbit]++
		}
	}
	return mem, shortCnt
}

func findMinShortSize(sorted [][]counterElt) (int32, int32) {

	sz := int32(0)
	minCost, shortCnt := memIncrOfShortSize(sorted, 0)

	for shortSize := int32(1); shortSize < maxShortSize+1; shortSize++ {
		incr, s := memIncrOfShortSize(sorted, shortSize)
		if incr < minCost {
			minCost = incr
			sz = shortSize
			shortCnt = s
		}
	}

	return sz, shortCnt
}

func (c *creator) build() *Slim {

	sorted := sortedBMCounts(c.innerBMCnt)
	shortSize, shortCnt := findMinShortSize(sorted)

	_ = shortCnt

	ns := &Slim{
		ShortSize:   shortSize,
		BigInnerCnt: c.bigCnt,
	}

	if len(c.clusteredOffsets) > 0 {
		c.clusteredStarts = append(c.clusteredStarts, uint32(len(c.clusteredOffsets)))
		c.clusteredOffsets = append(c.clusteredOffsets, uint32(len(c.clusteredBytes)))
		ns.Clustered = &Clustered{
			Starts:  c.clusteredStarts,
			Offsets: c.clusteredOffsets,
			Bytes:   c.clusteredBytes,
		}
	}

	// Mapping most used 17-bit bitmap inner node to short inner node.
	//
	// SlimTrie tries to replace some most used 17-bit bitmap with a shorter bitmap,
	// with the guarantee that the number of "1" does not change, thus the Rank
	// operation still work.
	//
	mostUsed := map[uint64]int32{}

	for short := int32(0); short < 1<<uint(ns.ShortSize); short++ {

		nbit := int32(bits.OnesCount64(uint64(short)))

		bms := sorted[nbit]

		if len(bms) > 0 {
			bm := bms[0].bitmap17
			mostUsed[bm] = short

			ns.ShortTable = append(ns.ShortTable, uint32(bm))

			sorted[nbit] = sorted[nbit][1:]

		} else {
			ns.ShortTable = append(ns.ShortTable, 0)
		}
	}

	// convert most used node bitmap to short

	// shortIndex is a bitmap indicating which inner node is replaced with a
	// short one.
	shortIndex := make([]int32, 0, c.nodeCnt)

	for innerI := c.bigCnt; innerI < int32(len(c.innerBMs)); innerI++ {
		bmindex := c.innerBMs[innerI]

		bm := bitmap.Of(bmindex)[0]
		short, has := mostUsed[bm]
		if has {
			idx := bitmap.ToArray([]uint64{uint64(short)})
			c.innerBMs[innerI] = idx
			c.innerSizes[innerI] = ns.ShortSize

			// index by inner index
			shortIndex = append(shortIndex, innerI)

		}
	}

	innerCnt := int32(len(c.innerIndexes))

	ns.ShortBM = newBM(shortIndex, innerCnt, "r64")

	// If it is empty, do not create NodeTypeBM. Query funcs check this field to
	// to determine if it is empty.
	if c.nodeCnt > 0 {
		// Extend to avoid index out of bound panic.
		ns.NodeTypeBM = newBM(c.innerIndexes, c.nodeCnt, "r64")
	}
	ns.Inners = &Bitmap{
		Words: bitmap.OfMany(c.innerBMs, c.innerSizes),
	}
	ns.Inners.indexit("r128")

	ns.InnerPrefixes = &VLenArray{}
	ns.InnerPrefixes.EltCnt = int32(len(c.prefixIndexes))
	ns.InnerPrefixes.PresenceBM = newBM(c.prefixIndexes, innerCnt, "r128")
	if *c.option.InnerPrefix {
		ns.InnerPrefixes.PositionBM = newBM(stepToPos(c.prefixByteLens, 0), 0, "s32")
		ns.InnerPrefixes.Bytes = c.prefixes

	} else {
		ns.InnerPrefixes.FixedSize = 2
		ns.InnerPrefixes.Bytes = c.prefix4BitLens
	}

	if *c.option.LeafPrefix {
		ns.LeafPrefixes = &VLenArray{}
		ns.LeafPrefixes.PresenceBM = newBM(c.leafPrefixIndexes, c.leafCnt, "r64")
		ns.LeafPrefixes.PositionBM = newBM(stepToPos(c.leafPrefixLens, 0), 0, "s32")
		ns.LeafPrefixes.Bytes = c.leafPrefixes
	}

	return ns
}

func (c *creator) buildLeaves(bytesValues [][]byte) *VLenArray {

	// Since 0.5.12
	// when creating, creator only records the value indexes;
	// when unmarshal and rebuild old version < 0.5.10, it appends leaves one by one.

	if !c.withLeaves {
		return nil
	}

	leaves := &VLenArray{}

	if len(c.leafIndexes) > 0 {
		sz := 0
		for _, idx := range c.leafIndexes {
			sz += len(bytesValues[idx])
		}
		lb := make([]byte, 0, sz)
		for _, idx := range c.leafIndexes {
			lb = append(lb, bytesValues[idx]...)
		}
		leaves.Bytes = lb

	} else {

		// maybe an empty slim, e.g., c.leaves is empty, or a slim with leaves filled
		n := len(c.leaves)
		sz := 0
		for _, elt := range c.leaves {
			sz += len(elt)
		}

		lb := make([]byte, 0, sz)
		for i := 0; i < n; i++ {
			lb = append(lb, c.leaves[i]...)
		}
		leaves.Bytes = lb
	}
	return leaves

}

func newSlim(keys []string, bytesValues [][]byte, opt *Opt) (*Slim, error) {

	maxLevel := int32(0x7fffffff)
	if opt.MaxLevel != nil {
		maxLevel = *opt.MaxLevel
	}

	n := len(keys)
	if n == 0 {
		return &Slim{}, nil
	}

	for i := 0; i < n-1; i++ {
		if keys[i] >= keys[i+1] {
			return nil, errors.Wrapf(ErrKeyOutOfOrder,
				"keys[%d] >= keys[%d] %s %s", i, i+1, keys[i], keys[i+1])
		}
	}

	tokeep := newToKeep(n, bytesValues, opt)

	sb := sigbits.New(keys)
	c := newCreator(n, bytesValues != nil, opt)

	queue := make([]subset, 0, n*2)
	queue = append(queue, subset{0, int32(n), 0, 1})

	for i := 0; i < len(queue); i++ {
		nid := int32(i)
		o := queue[i]
		s, e := o.keyStart, o.keyEnd

		// single key, it is a leaf
		if e-s == 1 {
			must.Be.True(tokeep[s])
			c.addLeafIndex(nid, o.level, s)
			c.setLeafPrefix(nid, keys[s], o.fromKeyBit)
			continue
		}

		// create an inner node

		wordStart, prefCounts := sb.CountPrefixes(s, e, maxWordSize)
		_ = prefCounts

		var wordsize int32
		var bitmapSize int32

		if c.isBig {

			prefCnt := prefCounts[8-(wordStart&7)]

			if prefCnt > 10 {
				// create big inner node with 257 bits
				must.Be.Equal(int32(0), o.fromKeyBit&7)
				wordStart &= ^7
				wordsize = bigWordSize
				bitmapSize = bigInnerSize

				prefLen := (wordStart - o.fromKeyBit) / bigWordSize
				if prefLen < minPrefix {
					wordStart = o.fromKeyBit
				}
			} else {
				// too small, stop creatting big node
				c.isBig = false
			}
		}

		if !c.isBig {
			must.Be.Equal(int32(0), o.fromKeyBit&3)
			wordStart &= ^3
			wordsize = wordSize
			bitmapSize = innerSize

			prefLen := (wordStart - o.fromKeyBit) / wordSize
			if prefLen < minPrefix {
				wordStart = o.fromKeyBit
			}
		}

		if wordStart < o.fromKeyBit {
			panic("wordStart smaller than o.fromKeyBit")
		}

		ks := make([]string, 0)
		for i := s; i < e; i++ {
			if tokeep[i] {
				ks = append(ks, keys[i])
			}
		}

		if o.level == maxLevel-1 {
			// too many levels, build clustered leaves

			// Manually set the leaf node level.
			// Because clustered leaves does not need creator.
			c.maxLevel = o.level + 1
			c.addClusteredInner(nid, o.level, o.fromKeyBit, wordStart, ks)

			// add all lower level nodes as leaves
			if bytesValues != nil {
				for i := s; i < e; i++ {
					if tokeep[i] {
						queue = append(queue, subset{
							keyStart: i,
							keyEnd:   i + 1,
							level:    o.level + 1,
						})
					}
				}
			}
			continue
		}

		// A label is a word with 0, 4 or 8 bits.
		// A path is an encoded representation of both the length and the bits.
		labelPaths := bmtree.PathsOf(ks, wordStart, wordsize, true)
		must.Be.True(len(labelPaths) > 0)

		idxs := make([]int32, len(labelPaths))
		for i, p := range labelPaths {
			idxs[i] = bmtree.PathToIndex(bitmapSize, p)
		}

		// Without the bits of label word at parent node
		c.addInner(nid, o.level, idxs, bitmapSize, o.fromKeyBit, wordStart, keys[s])

		// put keys with the same starting word to queue.

		for _, pth := range labelPaths {

			// Find the first key starting with label
			for ; s < e; s++ {
				kpath := bmtree.PathOf(keys[s], wordStart, wordsize)
				if kpath == pth {
					break
				}
			}

			// Continue looking for the first key not starting with label
			var j int32
			for j = s + 1; j < e; j++ {
				kpath := bmtree.PathOf(keys[j], wordStart, wordsize)
				if kpath != pth {
					break
				}
			}

			p := subset{
				keyStart: s,
				keyEnd:   j,

				// skip the label word
				fromKeyBit: wordStart + bmtree.PathLen(pth),

				level: o.level + 1,
			}
			queue = append(queue, p)
			s = j
		}
	}

	slim := c.build()
	slim.Leaves = c.buildLeaves(bytesValues)

	return slim, nil
}

func encodeValues(n int, values interface{}, e encode.Encoder) [][]byte {
	if values == nil {
		return nil
	}

	vals := make([][]byte, 0, n)
	rvals := reflect.ValueOf(values)

	for i := 0; i < n; i++ {
		v := getV(rvals, int32(i))
		bs := e.Encode(v)
		vals = append(vals, bs)
	}
	return vals
}

// newToKeep creates a []bool about which record to keep in slim.
// If DedupValue is true, value[i+1] with the same value with value[i] do not need to keep.
func newToKeep(n int, values [][]byte, opt *Opt) []bool {

	tokeep := make([]bool, n)

	// If slim does not store value, it has to store all keys.
	if *opt.DedupValue && values != nil {
		tokeep[0] = true
		for i := 1; i < n; i++ {
			tokeep[i] = bytes.Compare(values[i-1], values[i]) != 0
		}
		return tokeep
	}

	for i := 0; i < n; i++ {
		tokeep[i] = true
	}

	return tokeep
}

func getV(reflectSlice reflect.Value, i int32) interface{} {
	if reflectSlice.IsNil() {
		return nil
	}
	return reflectSlice.Index(int(i)).Interface()
}

func stepToPos(steps []int32, shift int32) []int32 {

	mask := int32(bitmap.Mask[shift])

	n := int32(len(steps))
	ps := make([]int32, n+1)
	p := int32(0)
	for i := int32(0); i < n; i++ {
		ps[i] = p

		must.Be.Zero(steps[i] & mask)

		p += steps[i] >> uint(shift)
	}
	ps[n] = p
	return ps
}
