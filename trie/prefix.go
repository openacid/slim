package trie

import (
	"bytes"
	"math/bits"

	"github.com/openacid/low/bitmap"
	"github.com/openacid/must"
)

// newPrefix create a prefix.
// Our prefix represents a non-8-aligned bit serias.
// The first byte is a control byte.
// The following are content bytes.
func newPrefix(str string, s, e int32) []byte {

	prefixTail := byte(0)

	se := (e + 7) >> 3
	// a control byte and strings
	pref := append([]byte{0}, []byte(str[s>>3:se])...)

	res := e & 7

	if res != 0 {
		pref[0] = 1
		prefixTail = 1 << uint(7-res)

		last := pref[len(pref)-1]
		last = last&byte(^bitmap.MaskUpto[7-res]) | prefixTail
		pref[len(pref)-1] = last
	}

	return pref
}

// prefixLen retrives the length in bit of a prefix.
func prefixLen(pref []byte) int32 {

	must.Be.OK(func() {
		must.Be.True(len(pref) > 0)
	})

	pl := int32(len(pref)) - 1
	if pref[0]&1 == 0 {
		return pl << 3
	}

	plast := pref[pl]
	nZero := int32(bits.TrailingZeros8(plast))
	return pl<<3 - nZero - 1
}

// prefixCompare compare n bits of a key with prefix.
//
// Performance
// ~ 15 ns/op for > 8 bytes string.
// ~ 10 ns/op for <= 8 bytes string with config.
// ~ 8 ns/op for <= 8 bytes string without config or kl < pl.
//
func prefixCompare(key string, cPref []byte) int {

	kl := int32(len(key))
	pref := cPref[1:]
	pl := int32(len(pref))

	if kl > pl {
		kl = pl
	}

	if cPref[0]&1 == 0 || kl < pl {
		if pl > 8 {
			return bytes.Compare([]byte(key[:kl]), pref)
		} else {
			var i int32
			for i = 0; i < kl; i++ {
				if key[i] < pref[i] {
					return -1
				} else if key[i] > pref[i] {
					return 1
				}
			}

			if i < pl {
				return -1
			}
			return 0

		}
	}

	pl--

	// sub := key[:pl]
	// sbytes := []byte(key[:pl])

	if pl > 8 {
		// rst := bytes.Compare(sbytes, pref[:pl])
		rst := bytes.Compare([]byte(key[:pl]), pref[:pl])
		if rst != 0 {
			return rst
		}
	} else {
		var i int32
		for i = 0; i < pl && key[i] == pref[i]; i++ {
		}

		if i < pl {
			if key[i] < pref[i] {
				return -1
			} else if key[i] > pref[i] {
				return 1
			}
		}
	}

	plast := pref[pl]

	mask := ^(plast ^ (plast - 1))

	klast := key[pl] & mask
	plast = plast & mask

	if klast > plast {
		return 1
	} else if klast < plast {
		return -1
	}
	return 0
}
