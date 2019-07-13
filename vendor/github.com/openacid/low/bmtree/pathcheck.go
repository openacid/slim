package bmtree

import (
	"math/bits"

	"github.com/openacid/must"
)

func pathCheck(path uint64) {

	// path only has at most 30 bits
	must.Be.Equal(uint64(0), path&0xc0000000c0000000)

	if uint32(path) == 0 {
		return
	}

	// path mask must be consecutive "1"s

	// fill trailing 0 to 1
	extended := uint32(path | (path - 1))
	pheight := 32 - bits.LeadingZeros32(uint32(path))
	must.Be.Equal(pheight, bits.OnesCount32(extended))

	// path bits must be shorter than mask

	pmask := uint32(path)
	pbits := uint32(path >> 32)

	must.Be.Equal(uint32(0), ^pmask&pbits)

}
