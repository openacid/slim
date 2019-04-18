// Package strhelper provides string operations functions
package strhelper

var (
	// mask for 1, 2, 4, 8 bit word
	wordMask = []byte{
		// 1, 2, 3, 4, 5, 6, 7, 8
		0, 1, 3, 0, 15, 0, 0, 0, 255,
	}
)

// ToBitWords split a string into a slice of byte.
// A char in string is split into 8/`n` `n`-bit words
// Value of every byte is in range [0, 2^n-1].
// `n` must be a one of [1, 2, 4, 8].
//
// Significant bits in a char is place at left.
// Thus the result byte slice keeps order with the original string.
func ToBitWords(s string, n int) []byte {
	if wordMask[n] == 0 {
		panic("n must be one of 1, 2, 4, 8")
	}

	mask := wordMask[n]

	// number of words per char
	m := 8 / n
	lenSrc := len(s)
	words := make([]byte, lenSrc*m)

	for i := 0; i < lenSrc; i++ {
		b := s[i]

		for j := 0; j < m; j++ {
			words[i*m+j] = (b >> uint(8-n*j-n)) & mask
		}
	}
	return words
}

// SliceToBitWords converts a `[]string` to a n-bit word `[][]byte`.
func SliceToBitWords(strs []string, n int) [][]byte {
	rst := make([][]byte, len(strs))
	for i, s := range strs {
		rst[i] = ToBitWords(s, n)
	}
	return rst
}

// FromBitWords is the reverse of ToBitWords.
// It composes a string of which each byte is formed from 8/n words from bs.
func FromBitWords(bs []byte, n int) string {
	if wordMask[n] == 0 {
		panic("n must be one of 1, 2, 4, 8")
	}

	// number of words per char
	m := 8 / n
	sz := (len(bs) + m - 1) / m
	strbs := make([]byte, sz)

	var b byte
	for i := 0; i < len(strbs); i++ {
		b = 0
		for j := 0; j < m; j++ {
			if i*m+j < len(bs) {
				b = (b << uint(n)) + bs[i*m+j]
			} else {
				b = b << uint(n)
			}
		}
		strbs[i] = b
	}

	return string(strbs)
}

// SliceFromBitWords converts a `[][]byte` back to a `[]string`.
func SliceFromBitWords(bytesslice [][]byte, n int) []string {
	rst := make([]string, len(bytesslice))
	for i, s := range bytesslice {
		rst[i] = FromBitWords(s, n)
	}
	return rst
}
