package bmtree

import "github.com/openacid/low/bitmap"

// NewPath creates a path, which is a uint64, from node searching path, path
// length and tree height.
//
// Since 0.1.9
func NewPath(searchingBits uint64, length, height int32) uint64 {
	return (searchingBits << 32) | (bitmap.Mask[length] << uint(height-length))
}

// PathOf creates a path, which is a uint64, from a string and sepcified
// starting bit position and height.
//
// Since 0.1.9
func PathOf(s string, frombit int32, height int32) uint64 {
	plen, path := bitmap.FromStr32(s, frombit, frombit+height)
	return NewPath(path, plen, height)
}

// PathsOf creates more than one node search paths at a time.
//
// Since 0.1.9
func PathsOf(keys []string, frombit int32, height int32, dedup bool) []uint64 {
	l := len(keys)
	rst := make([]uint64, 0, l)
	prev := ^uint64(0)
	for _, s := range keys {
		p := PathOf(s, frombit, height)
		if !dedup || p != prev {
			rst = append(rst, p)
		}
		prev = p
	}
	return rst
}
