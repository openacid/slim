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
	"fmt"

	"github.com/openacid/low/bitmap"
	"github.com/openacid/slim/encode"
)

const (
	// MaxNodeCnt is the max number of node. Node id in SlimTrie is int32.
	MaxNodeCnt = (1 << 31) - 1

	// NodeTypeInner represents an inner node.
	//
	// Since 0.5.10
	NodeTypeInner = int32(1)

	// NodeTypeLeaf represents a leaf node.
	//
	// Since 0.5.10
	NodeTypeLeaf = int32(0)

	// minPrefix is the minimal prefix to create.
	// If a sub set keys have common prefix but prefix length is smaller than
	// minPrefix, it creates an inner node instead of a step.
	minPrefix = int32(0)

	// 16 4-bit words and a 0-bit word
	wordSize  = int32(4)
	innerSize = int32(1)<<uint(wordSize) + 1

	// 256 8-bit words and a 0-bit word
	bigWordSize  = int32(8)
	bigInnerSize = int32(1)<<uint(bigWordSize) + 1

	// maxShortSize is the max bits a short node can have.
	// The number of bits of short node is decided during creating.
	maxShortSize = int32(10)

	// maxWordSize is the longest bit to look forward when creating.
	maxWordSize = int32(24)
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
//
// Since 0.2.0
type SlimTrie struct {
	nodes   *Nodes
	encoder encode.Encoder
}

// Opt specifies options for creating a SlimTrie.
//
// By default SlimTrie remove unnecessary information for locating a PRESENT
// key, such as trie branch content.
// And it introduces false positives.
// In this case it is the responsibility of upper level to ensure whether a query result
// is absolutely correct.
//
// But SlimTrie can also store complete key information thus let it always
// returns correct query result, without any false positive.
//
// Since 0.5.10
type Opt struct {
	// CompleteInner tells SlimTrie to store text on a trie branch to inner
	// node(not to leaf node), instead of storing only branch length.
	// With this option SlimTrie costs more space but reduces false positive
	// rate.
	//
	// Default false.
	//
	// Since 0.5.10
	CompleteInner bool

	// CompleteLeaf tells SlimTrie to store text on a branch to leaf node.
	// With this option SlimTrie costs more space but reduces false positive
	// rate.
	//
	// Default false.
	//
	// Since 0.5.10
	CompleteLeaf bool

	// Complete tells SlimTrie to store complete keys content.
	// This option implies "CompleteInner" and "CompleteLeaf".
	// With this option there is no false positive and SlimTrie works just like
	// a static key-value map.
	//
	// Default false.
	//
	// Since 0.5.10
	Complete bool
}

func (va *Nodes) GetVersion() string {
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
func NewSlimTrie(e encode.Encoder, keys []string, values interface{}, opts ...Opt) (*SlimTrie, error) {

	opt := Opt{}

	if len(opts) > 0 {
		opt = opts[0]
	}

	if opt.Complete {
		opt.CompleteInner = true
		opt.CompleteLeaf = true
	}

	return newSlimTrie(e, keys, values, opt)
}

// func (st *SlimTrie) GetStat() map[string]float64 {
//     return st.nodes.Stat
// }

func (st *SlimTrie) content() []string {
	rst := []string{}
	ns := st.nodes
	rst = append(rst, fmt.Sprintf(`InnerBM: %+v`, bitmap.ToArray(ns.NodeTypeBM.Words)))
	rst = append(rst, fmt.Sprintf(`StepBM: %+v`, bitmap.ToArray(ns.InnerPrefixes.PresenceBM.Words)))
	if ns.InnerPrefixes.PositionBM != nil {
		rst = append(rst, fmt.Sprintf(`PrefixStarts: %+v`, bitmap.ToArray(ns.InnerPrefixes.PositionBM.Words)))
	} else {
		rst = append(rst, fmt.Sprintf(`Steps: %+v`, ns.InnerPrefixes.Bytes))
	}

	return rst
}
