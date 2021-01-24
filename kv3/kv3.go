package kv3

import (
	"fmt"
	"strings"

	"github.com/openacid/low/sigbits"
	"github.com/openacid/slim/trie"
)

type KV3 struct {
	slim     *trie.SlimTrie
	segments []trie.Records
}

func New(keys []string, values [][]byte) *KV3 {
	kv := &KV3{}
	return kv
}

type shard struct {
	prefix   string
	suffixes []string
}

func (s *shard) String() string {
	rst := make([]string, 0, len(s.suffixes)+1)
	rst = append(rst, "="+s.prefix)
	for i := 0; i < len(s.suffixes); i++ {
		rst = append(rst, "-"+strings.Repeat(" ", len(s.prefix))+s.suffixes[i])
	}

	return strings.Join(rst, "\n")
}

func newShard(prefixLen int32, keys []string) *shard {

	var prefix string
	if prefixLen == 0 {
		prefix = ""
	} else {
		prefix = keys[0][:prefixLen]
	}

	s := &shard{
		prefix:   prefix,
		suffixes: make([]string, 0, len(keys)),
	}
	for _, k := range keys {
		s.suffixes = append(s.suffixes, k[len(prefix):])
	}

	return s
}

func splitKeys(keys []string, thre int32) []*shard {

	firstDiffs := sigbits.FirstDiffBits(keys)

	n := int32(len(keys))

	// ends[i] is the end position of item i.
	ends := make([]int32, n)
	for i := int32(0); i < n; i++ {
		ends[i] = i
	}
	starts := make([]int32, n)
	for i := int32(0); i < n; i++ {
		starts[i] = i
	}
	prefSizes := make([]int32, n)
	for i := int32(0); i < n; i++ {
		prefSizes[i] = int32(len(keys[i]))
	}

	findMergeable := func(i int32) (int32, int32) {

		prevSize := int32(0)
		nextSize := int32(0)
		currSize := prefSizes[i]

		prev := int32(0)
		nxt := int32(0)

		newPrefLen := firstDiffs[i]

		// whether the prefix at i is less or equal the previous one and the
		// next one.
		if i > 0 {
			prev = starts[i-1]
			prevSize = prefSizes[prev]
		} else {
			prev = -1
			prevSize = 0
		}

		if ends[i] < n-1 {
			nxt = ends[i] + 1
			nextSize = prefSizes[nxt]
		} else {
			nxt = -1
			nextSize = 0
		}

		if currSize >= prevSize && currSize >= nextSize {
			if prevSize > nextSize {
				return prev, -1
			} else if prevSize < nextSize {
				return -1, nxt
			} else {
				return prev, nxt
			}

		}

		return -1, -1
	}

	i := int32(0)

	for i < n {

		mergeLeft, mergeRight := findMergeable(i)
		fmt.Println("i=", i, "merge:", mergeLeft, mergeRight, "key:", keys[i])
		if mergeLeft != -1 {

			fmt.Println("left:", keys[mergeLeft])

			if ends[i]-mergeLeft+1 <= thre {
				fmt.Println("merge:", mergeLeft, i)
				prefSizes[mergeLeft] = min(firstDiffs[ends[mergeLeft]]>>3, prefSizes[mergeLeft])

				ends[mergeLeft] = ends[i]
				starts[ends[i]] = mergeLeft
				fmt.Println(starts[ends[i]], ends[mergeLeft])

				i = mergeLeft

				// the one before mergeLeft now may be able to merge.
				if mergeLeft > 0 {
					i = starts[i-1]
				}
				fmt.Println("i:", i)
				continue
			}
		}

		if mergeRight != -1 {
			fmt.Println("right:", keys[mergeRight])
			if ends[mergeRight]-i+1 <= thre {
				prefSizes[i] = min(firstDiffs[ends[i]]>>3, prefSizes[mergeRight])

				ends[i] = ends[mergeRight]
				starts[ends[mergeRight]] = i

				continue
			}
		}

		i = ends[i] + 1
		fmt.Println("update i:", i)
	}

	rst := make([]*shard, 0)
	for i := int32(0); i < n; i = ends[i] + 1 {
		prefix := keys[i][:prefSizes[i]]

		s := &shard{
			prefix:   prefix,
			suffixes: make([]string, 0, ends[i]-starts[i]+1),
		}
		for j := starts[i]; j <= ends[i]; j++ {
			s.suffixes = append(s.suffixes, keys[j][len(prefix):])
		}
		rst = append(rst, s)
	}

	return rst
}

func splitKeys3(keys []string, thre0 int32) []*shard {
	fmt.Println("keys:", keys, "thre0", thre0)

	rst := make([]*shard, 0)

	stack := make([]int32, 1)

	updateStack := func(prefixLen int32) (bool, int32) {
		// return where to add new prefixLen

		var j int32
		l := int32(len(stack))
		if prefixLen <= stack[l-1] && prefixLen > stack[l-2] {
			// allow to update the last one.
			// the last one is the current shard
			fmt.Println("updated stack:", stack)
			stack[l-1] = prefixLen
			return true, l - 1
		}
		if prefixLen > stack[l-1] {
			return true, l - 1
		}

		fmt.Println("fail to update stack")
		return false, j
	}

	firstDiffs := sigbits.FirstDiffBits(keys)

	n := int32(len(keys))
	i := int32(0)

	for i < n {

		rangeStart := i
		thre := min(n-i, thre0)
		prefixLen := int32(len(keys[i]))

		if i > 0 {
			diff := firstDiffs[i-1] >> 3
			last := len(stack) - 1
			if diff < stack[last] {
				for len(stack) > 0 && stack[len(stack)-1] >= diff {
					stack = stack[:len(stack)-1]
				}
			}
			stack = append(stack, diff)
		}

		stack = append(stack, prefixLen)

		i++

		for {
			if i-rangeStart == thre {
				break
			}

			curr := firstDiffs[i-1] >> 3
			newPrefixLen := min(curr, prefixLen)

			ok, _ := updateStack(newPrefixLen)
			if ok {
				i++
				continue
			}

			// conflict: current prefix equals a previous one
			break
		}

		s := newShard(stack[len(stack)-1], keys[rangeStart:i])
		fmt.Println("shard:")
		fmt.Println(s)
		rst = append(rst, s)

		// fmt.Println("from key:", keys[i])
	}

	return rst
}

func min(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}
func min3(a, b, c int32) int32 {
	return min(a, min(b, c))
}
