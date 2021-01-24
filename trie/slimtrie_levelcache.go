package trie

import (
	"unsafe"

	"github.com/openacid/low/bitmap"
	"github.com/openacid/low/size"
)

// initLevelCache builds the cache for selected level the count of leaf nodes.
//
// Since 0.5.12
func (st *SlimTrie) initLevelCache() {
	ns := st.inner

	sz := size.Of(st)
	mem := sz * 5 / 100

	cachePolicy := st.findLevelCaches(int64(mem))
	for _, lvl := range cachePolicy.levels {
		innerCnt := st.levels[lvl].inner - st.levels[lvl-1].inner
		cache := make([]innerCache, 0, innerCnt)
		for i := st.levels[lvl-1].total; i < st.levels[lvl].total; i++ {
			w := ns.NodeTypeBM.Words[i>>6]
			if w&bitmap.Bit[i&63] != 0 {
				// this is an inner node
				v := st.countReachableLeaves(i)
				cache = append(cache, innerCache{nodeId: i, leafCount: v})
			}
		}

		st.levels[lvl].cache = cache
	}
}

// countReachableLeaves count the number of reachable leaves by walking from an
// inner node or from a preceding node at this level.
//
//   I0 --- I1 --- L5
//      --- L2
//      --- L3
//          I4 --- L6
//             --- L7
//
// In the above example of a trie:
//  q(I4) = 3 // L5, L2, L3
//  q(L3) = 2 // L5, L2
//  q(L6) = 2 // L5, L6
func (st *SlimTrie) countReachableLeaves(nodeId int32) int32 {

	// find which level this node is at:
	nodeLevel := int32(0)
	for lvl, lvlInfo := range st.levels {
		if lvlInfo.total > nodeId {
			nodeLevel = int32(lvl)
			break
		}
	}

	cursor := &walkingCursor{
		id:         nodeId + 1,
		smallerCnt: 0,
		lvl:        nodeLevel,
	}

	leafIndex := st.cursorLeafIndex(cursor, false)
	return leafIndex
}

// findLevelCaches finds the levels to add cache to, in order to maximize steps
// reduction within specified memory limit: `maxMem` in bytes.
//
// Because the performance of Get() is level related, walking down to a level
// takes about 8 ns. Thus one of the ways to reduce latency of Get() is to
// reduce the levels it has to walk.
//
// With this leaves pre-count cache on inner nodes, it does not need to walk
// down to the bottom level to find the original position of a record.
//
// The bottom level is `bottom`
// With several non-bottom level with pre-count leaves count,
//
// The total number of steps to walk all leaves is the sum of area of all rectangles:
//
//       N.O. leaves
//       ^
//       |
//  L(b) |------------------------------------------+
//       |....................................******|
//  L(x) |----------------------------------+-------|
//       |........................**********|       |
//       |.........**************...........|       |
//  L(y) |--------+-------------------------|       |
//       |......*.|                         |       |
//       |...***..|                         |       |
//       |..*.....|                         |       |
//       |.*......|                         |       |
//       +-------------------------------------------> level
//       0        y                         x       b
//
// Since 0.5.12
func (st *SlimTrie) findLevelCaches(maxMem int64) levelCachePolicy {

	cacheItemSize := int64(unsafe.Sizeof(innerCache{}))
	bottom := int64(len(st.levels) - 1)

	// The total steps without walking down to bottom(area of dots):
	//
	//       N.O. leaves
	//       ^
	//       |
	//       |....................................******|
	//       |..................................*       |
	//       |........................**********|       |
	//       |.........**************           |       |
	//       |........*                         |       |
	//       |......* |                         |       |
	//       |...***  |                         |       |
	//       |..*     |                         |       |
	//       |.*      |                         |       |
	//       +-------------------------------------------> level
	//       0        y                         x       b
	minSteps := int64(0)
	for i := 1; i < len(st.levels); i++ {
		minSteps += int64(st.levels[i].leaf-st.levels[i-1].leaf) * int64(i)
	}

	// from every level, the additional steps introduced.
	// That is the max steps that can be saved by adding a level cache.
	maxSave := make([]int64, bottom+1)
	for i := bottom - 1; i > 0; i-- {
		lvl := st.levels[i]
		prev := st.levels[i-1]
		maxSave[i] = (bottom-i)*int64(lvl.leaf-prev.leaf) + maxSave[i+1]
	}

	best := levelCachePolicy{levels: make([]int32, 0)}
	curr := levelCachePolicy{levels: make([]int32, 0)}

	trials := 0

	// Adding a level cache at y, it reduces steps that falls in the area of "=".
	// And when adding more level cache after y, the max possible reduced steps
	// are maxSave[y].
	//       N.O. leaves
	//       ^
	//       |
	//       |....................................******|
	//       |..................................*       |
	//       |........................**********        |
	//       |.........**************                   |
	//       |........*=================================|
	//       |......* |                                 |
	//       |...***  |                                 |
	//       |..*     |                                 |
	//       |.*      |                                 |
	//       +-------------------------------------------> level
	//       0        y                                 b
	var dfs func(int64, int64)
	dfs = func(left int64, mem int64) {
		// dfs recursively evaluates the steps reduction at every level in (left, bottom). (exclude left and bottom)

		// N.O. leaves upto level left
		frmLeavesCnt := st.levels[left].leaf

		for i := left + 1; i < bottom; i++ {

			leavesCnt := st.levels[i].leaf
			if leavesCnt == 0 {
				// no leaf at or before this level, no need to cache
				continue
			}

			levelReduced := (bottom - i) * int64(leavesCnt-frmLeavesCnt)
			maxPossibleReduce := levelReduced + maxSave[i+1]

			memCost := int64(st.levels[i].inner-st.levels[i-1].inner) * cacheItemSize
			trials++

			if curr.reduced+maxPossibleReduce < best.reduced {
				continue
			}
			if memCost > mem {
				continue
			}

			curr.reduced += levelReduced
			curr.levels = append(curr.levels, int32(i))

			if curr.reduced > best.reduced {
				best.reduced = curr.reduced
				best.levels = append(best.levels[0:0], curr.levels...)
			}

			dfs(i, mem-memCost)

			// popup
			curr.levels = curr.levels[:len(curr.levels)-1]
			curr.reduced -= levelReduced
		}
	}

	dfs(0, maxMem)

	best.steps = stepsWithCache(st.levels, best.levels)
	return best
}

func stepsWithCache(levels []levelInfo, chosen []int32) int64 {

	bottom := int64(len(levels)) - 1

	// number of steps to walk all keys.
	steps := int64(0)

	prev := int64(0)
	for _, chosenLvl := range chosen {
		lvlInfo := levels[chosenLvl]
		steps += int64(chosenLvl) * (int64(lvlInfo.leaf) - prev)
		prev = int64(lvlInfo.leaf)
	}
	steps += bottom * (int64(levels[bottom].leaf) - prev)

	return steps
}
