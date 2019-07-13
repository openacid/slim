package bmtree

// Decode a encoded bitmap, retrieves all paths in it.
//
// Since 0.1.9
func Decode(bitmapSize int32, bm []uint64) []uint64 {

	rst := make([]uint64, 0)

	paths := AllPaths(bitmapSize, 0, 1<<63)

	for _, p := range paths {
		idx := PathToIndex(bitmapSize, p)

		wordI := idx >> 6

		if int32(len(bm)) > wordI && bm[wordI]&(1<<uint(idx&63)) != 0 {
			rst = append(rst, p)
		}
	}
	return rst
}
