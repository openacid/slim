package bitmap

var (
	// Mask are pre-calculated width-indexed bit masks.
	// E.g. Mask[1] is 63 "0" and 1 "1": 000..01 .
	//
	// Since 0.1.9
	Mask [65]uint64

	// RMask are pre-calculated reverse mask of Mask.
	//
	// Since 0.1.9
	RMask [65]uint64

	// MaskUpto are mask with bits set upto i-th bit(include i-th bit).
	// E.g. MaskUpto[1] == Mask[2] == 000..011 .
	//
	// Since 0.1.9
	MaskUpto [64]uint64

	// RMaskUpto are reverse of MaskUpto.
	//
	// Since 0.1.9
	RMaskUpto [64]uint64

	// Bit set i-th bit to 1.
	//
	// Since 0.1.9
	Bit [64]uint64
)

func initMasks() {
	for i := 0; i < 65; i++ {
		Mask[i] = (1 << uint(i)) - 1
		RMask[i] = ^Mask[i]
	}

	for i := 0; i < 64; i++ {
		MaskUpto[i] = (1 << uint(i+1)) - 1
		RMaskUpto[i] = ^MaskUpto[i]
		Bit[i] = 1 << uint(i)
	}
}
