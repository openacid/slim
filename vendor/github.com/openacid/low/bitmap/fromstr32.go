package bitmap

// FromStr32 returns a bit array from string and put them in the least
// significant "to-from" bits of a uint64.
// "to-from" must be smaller than or equal 32.
//
// It returns actual number of bits used from the string, and a uint64.
//
// Since 0.1.9
func FromStr32(s string, frombit, tobit int32) (int32, uint64) {

	size := tobit - frombit
	spanSize := tobit - (frombit & ^7)

	blen := int32(len(s)<<3) - frombit

	if blen > size {
		blen = size
	}

	if blen <= 0 {
		return 0, 0
	}

	l := int32(len(s))
	toByte := (tobit + 7) >> 3
	if l > toByte {
		l = toByte
	}

	i := frombit >> 3
	b := uint64(0)

	if i < l {
		b |= uint64(s[i]) << 32
		i++
		if i < l {
			b |= uint64(s[i]) << 24
			i++
			if i < l {
				b |= uint64(s[i]) << 16
				i++
				if i < l {
					b |= uint64(s[i]) << 8
					i++
					if i < l {
						b |= uint64(s[i])
					}
				}
			}
		}
	}

	return blen, (b >> uint(40-spanSize)) & Mask[size]
}
