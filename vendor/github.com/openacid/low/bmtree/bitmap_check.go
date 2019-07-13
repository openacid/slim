package bmtree

import (
	"math/bits"

	"github.com/openacid/must"
)

func bitmapSizeCheck(bitmapSize int32) {

	height := int32(31 - bits.LeadingZeros32(uint32(bitmapSize)))

	must.Be.True(height <= 30)
	must.Be.NotEqual(int32(0), bitmapSize)
}

func bitmapMustHaveLevel(bitmapSize, l int32) {
	must.Be.Equal(int32(1), (bitmapSize>>uint(l))&1,
		"level[pathlen] must be stored by bitmap")
}
