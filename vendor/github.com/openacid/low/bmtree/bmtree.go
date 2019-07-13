// Package bmtree encode a binary tree into a bitmap.
// Present nodes are set to 1 in bitmap, absent nodes are set to 0.
//
// A node is identified with its searching path:
// E.g. a path 0110 means going from root node and then choose left, right,
// right and left branch.
//
// "bmtree" supports encoding a binary tree of up to 30 levels and the result
// bitmap size is 0~2^31-1.
// E.g. the following three nodes "0", "01", "10" can be encoded into a 7 bits
// bitmap: "0101010".
// The mapping of all nodes in a 2 level tree are:
//
//    path          index in bitmap
//    ""            0
//    "0"           1
//    "00"          2
//    "01"          3
//    "1"           4
//    "10"          5
//    "11"          6
//
// How it works
//
// If there is a binary tree:
// an empty bit array is mapped to root node,
// "01" is mapped to node 01, etc.
//
//               root             h=2
//              /    \
//             0      1           h=1
//           /  \    /  \
//          00   01 10   11       h=0
//
// To assign a node to a index in bitmap,
// We  places them in "pre-order"(or NLR: parent node comes before left child,
// then right child),
// thus "0" comes before "00" in the bitmap:
//
//   Put a node at index x,
//   then start from x+1 to put its left-subtree,
//   then right-subtree, recursively.
//
// Root node is at the 0-th bit.
//
// Path:
//
// A path is for selecting a node.
// E.g.: a "0" in a path for selecting left child, a "1" for right child.
//
// The representation of a path is a uint64,
// the higher 32 bits are path mask,
// the lower 32 bits are searching bits:
//	0...011111000 0...0xxxxx000
//	<-- significant bits
//
// The Left most bit in path-mask indicates the tree height.
//
// A full binary tree has 2^(h+1)-1 nodes thus the result bitmap has 2^(h+1)-1
// bits.
//
// Since 0.1.9
package bmtree
