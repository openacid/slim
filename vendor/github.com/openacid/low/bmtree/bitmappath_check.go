package bmtree

import (
	"github.com/openacid/must"
)

func bitmapPathMustHaveEqualHeight(bitmapSize int32, path uint64) {

	// TODO remove these, add to Height() and PathHeight()
	bitmapSizeCheck(bitmapSize)
	pathCheck(path)
	if uint32(path) != 0 {
		must.Be.Equal(Height(bitmapSize), PathHeight(path),
			"bitmap height == path height")
	}
}
