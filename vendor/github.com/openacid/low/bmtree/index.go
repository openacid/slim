package bmtree

import (
	"math/bits"

	"github.com/openacid/low/bitmap"
	"github.com/openacid/must"
)

// PathToIndex convert a node searching path to a bit index in bitmap.
//
// Bitmap size
//
// The tree height is "h". The following is a tree of height 2.
//
//               root             h=2
//              /    \
//             0      1             1
//           /  \    /  \
//          00   01 10   11         0
//
// The bitmap size for a full binary tree is 2^(h+1) - 1, or
//
//  11111...111
//  `-- h+1 --'
//
// If we do not need to store a entire level of a tree, we do not need to assign
// bits for these node.
// Thus we reduce the bitmap size by 2^(h-i) for a tree without i-th level.
// E.g. if a tree of height 2 does not need level-1 nodes, the bitmap size is:
// 101 instead of 111.
//
// Thus the bitmap size "T" for a tree with some levels absent is:
//
//  T = t[h]t[h-1]...t[1]t[0] = 101110..011
//                              `-- h+1 --'
//
// In which t[h-i] == "1" indicates the bitmap stores i-th level nodes.
// And we see that the size of a subtree rooted at i-th level node is:
//  T>>(h-i)
//
// Mapping node to bitmap index
//
//
// If we start encode a subtree of height i from the x-th bit in bitmap,
// then:
//
// if the bitmap stores the i-th level nodes(t[h-i] == 1):
// then its left subtree starts from x + 1, otherwise its left subtree starts
// from x.
//
// Thus its left subtree starts from x + t[h-i].
// And its right subtree starts from x + t[h-i] + T>>(h-i+1).
//
// Because we starts with x=0, and i from h down to l:
//
// Thus the bit position of a path of length l is:
//
//  bitIndex =   t[0]   + p[h-1] * (T>>1)    // root at h
//             + t[1]   + p[h-2] * (T>>2)    // at h - 1
//             + t[2]   + p[h-3] * (T>>3)    // at h - 2
//             ...
//             + t[l-1] + p[0]   * (T>>l)
//
// Quick path
//
// If T is all "1", there is a quick path:
//
//                 h-1
//  bitIndex = l + sum(p[i] * 2^(i+1)) - rank1(p)
//                 i=0
//
//           = p * 2 + l - rank1(p)
//           = p * 2 + rank1(pathmask) - rank1(p)
//           = p * 2 + rank1(pathmask) - 32 + rank1(^p)
//
//           // "rank1" returns count of "1" in a int.
//
// Because we store the path-mask togeter with the path in uint64,
// we can use just one rank1 to get the index:
//
//	bitIndex = p * 2 + rank1(path^0xffffffff) - 32
//
// Since 0.1.9
func PathToIndex(bitmapSize int32, path uint64) int32 {

	must.Be.OK(func() {
		bitmapSizeCheck(bitmapSize)
		pathCheck(path)
		bitmapPathMustHaveEqualHeight(bitmapSize, path)
		bitmapMustHaveLevel(bitmapSize, PathLen(path))
	})

	height := Height(bitmapSize)
	sz := uint64(bitmapSize)

	if sz == bitmap.MaskUpto[height] {

		// sz is all "1", a full bitmap.

		return (int32(path>>32) << 1) + int32(bits.OnesCount64(path^0xffffffff00000000)) - 32

	} else if sz == bitmap.Bit[height] {

		// only leaf nodes

		return int32(path >> 32)

	} else {

		idx := shiftMulti(sz, path>>32, uint64(height))
		return int32(idx + uint64(bits.OnesCount64(sz&bitmap.Mask[PathLen(path)])))
	}
}

// PathToIndexLoose is same as PathToIndex except it allows the path to locate
// at a level the bitmap does not store.
// It returns the index and 1 if the level exist, otherwise the index and 0.
//
// Since 0.1.9
func PathToIndexLoose(bitmapSize int32, path uint64) (int32, int32) {

	must.Be.OK(func() {
		bitmapSizeCheck(bitmapSize)
		pathCheck(path)
		bitmapPathMustHaveEqualHeight(bitmapSize, path)
	})

	height := Height(bitmapSize)
	sz := uint64(bitmapSize)
	pl := PathLen(path)
	has := (bitmapSize >> uint(pl)) & 1

	if sz == bitmap.MaskUpto[height] {

		// sz is all "1", a full bitmap.

		return (int32(path>>32) << 1) + int32(bits.OnesCount64(path^0xffffffff00000000)) - 32, has

	} else if sz == bitmap.Bit[height] {

		// only leaf nodes

		return int32(path >> 32), has

	} else {

		idx := shiftMulti(path>>32, sz, uint64(height))
		return int32(idx + uint64(bits.OnesCount64(sz&bitmap.Mask[PathLen(path)]))), has
	}
}

// IndexToPath returns the path encoded to index.
//
// It determines every bit in path one by one:
//
// If the highest bit in path, p[h-1] == 0, then index is in range [1, 2^h);
// If p[h-1] == 1, then index is in range [2^h, 2^h+2^h-1).
// This way we narrow down the query into a subtree.
// Repeat it until index becomes to 0:
//
//	index[h] == 0: p[h-1] = 0, index -= 1
//	index[h] == 1: p[h-1] = 1, index -= 2^h
//
// In our implementation we have two optimization:
// First, the last 4 levels query with a lookup table.
// Second, before bit-by-bit parsing, it tries to find the longest common prefix
// of index - treeheight and index:
//
// Because index = p*2 + rank0(p),
// then we have index - width <= p*2 <= index.
// Thus common prefix of index - width and index is also a prefix of p*2
//
// Since 0.1.9
func IndexToPath(treeheight int32, index int32) uint64 {

	// calculate p*2.
	p2 := uint64(0)

	// Mask both the pathmask part and path part in the path*2.
	mask := uint64(0x0100000001) << uint(treeheight)

	if treeheight > 4 {

		// may overflow but ok
		i1 := index - treeheight
		i2 := index

		diffbits := 32 - int32(bits.LeadingZeros32(uint32(i1^i2)))
		// +1 because we are calculating path*2
		fixed := treeheight + 1 - diffbits
		var m uint64
		if fixed > 0 {
			// mask of fixed bits
			m = (mask << 1) - (uint64(0x0100000001) << uint(diffbits))

			// copy fixed bits from index to p2 directly
			p2 = ((uint64(index) << 32) | 0x00000000ffffffff) & m

			// If path = Pfixed Pleft:
			//   aaabbb00..
			//   ---==
			//   pa pb
			// index = path * 2 + rank0(path)
			//       = Pfixed * 2 + rank0(Pfixed)
			//       + Pleft * 2 + rank0(Pleft)
			// Thus
			//   index - Pfixed - rank0(Pfixed) = Pleft*2 + rank0(Pleft)
			index = index&int32(^m) - fixed + int32(bits.OnesCount32(uint32(index&int32(m))))
			mask >>= uint(fixed)
		}
	}

	// Loop ends if:
	// 1. index is 0: no more query needed.
	// 2. or reached the last 3 level, then use lookup table for query.
	for mask&15 == 0 && index > 0 {

		// if the index >= 2^h, the (h-1)-th bit in path must be 1
		maskAndPathBit := ((uint64(index) << 32) | 0x00000000ffffffff) & mask
		p2 |= maskAndPathBit

		// locate on left subtree(bit==0) or right subtree(bit==1) and continue with next level .
		if int32(maskAndPathBit>>32) == 0 {
			index--
		} else {
			index -= int32(maskAndPathBit >> 32)
		}
		mask = mask >> 1
	}

	p2 = (p2 >> 1) | idxToPath[mask&15][index]

	return p2
}

var (
	idxToPath = [][]uint64{
		// index becomes 0 before reaching the last 4 bits
		0: {
			0,
		},
		1: {
			(0x00000000 << 32) + 0x00000000, // 0   ""
		},
		2: {
			(0x00000000 << 32) + 0x00000000, // 0   ""
			(0x00000000 << 32) + 0x00000001, // 1   0
			(0x00000001 << 32) + 0x00000001, // 2   1
		},
		4: {
			(0x00000000 << 32) + 0x00000000, // 0   ""
			(0x00000000 << 32) + 0x00000002, // 1   0
			(0x00000000 << 32) + 0x00000003, // 2   00
			(0x00000001 << 32) + 0x00000003, // 3   01
			(0x00000002 << 32) + 0x00000002, // 4   1
			(0x00000002 << 32) + 0x00000003, // 5   10
			(0x00000003 << 32) + 0x00000003, // 6   11
		},
		8: {
			(0x00000000 << 32) + 0x00000000, // 0   ""
			(0x00000000 << 32) + 0x00000004, // 1   0
			(0x00000000 << 32) + 0x00000006, // 2   00
			(0x00000000 << 32) + 0x00000007, // 3   000
			(0x00000001 << 32) + 0x00000007, // 4   001
			(0x00000002 << 32) + 0x00000006, // 5   01
			(0x00000002 << 32) + 0x00000007, // 6   010
			(0x00000003 << 32) + 0x00000007, // 7   011
			(0x00000004 << 32) + 0x00000004, // 8   1
			(0x00000004 << 32) + 0x00000006, // 9   10
			(0x00000004 << 32) + 0x00000007, // 10  100
			(0x00000005 << 32) + 0x00000007, // 11  101
			(0x00000006 << 32) + 0x00000006, // 12  11
			(0x00000006 << 32) + 0x00000007, // 13  110
			(0x00000007 << 32) + 0x00000007, // 14  111
		},
	}
)
