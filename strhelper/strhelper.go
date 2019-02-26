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
