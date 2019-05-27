// Package bitword provides string operations functions
package bitword

// Interface defines operations for n-bit word.
//
// Since 0.1.4
type Interface interface {
	// FromStr split a string into a slice of `n`-bit words in byte.
	//
	// Since 0.1.4
	FromStr(string) []byte
	// FromStrs converts a `[]string` to a n-bit word `[][]byte`.
	//
	// Since 0.1.4
	FromStrs([]string) [][]byte
	// ToStr is the reverse of FromStr.
	//
	// Since 0.1.4
	ToStr([]byte) string
	// ToStrs converts a `[][]byte` back to a `[]string`.
	//
	// Since 0.1.4
	ToStrs([][]byte) []string
	// Get returns i-th n-bit word from a string.
	//
	// Since 0.1.4
	Get(string, int) byte
	// FirstDiff returns the index of the first different n-bit word,
	// from "start" upto "end".
	//
	// Since 0.1.4
	FirstDiff(a, b string, start, end int) int
}

// bitWord implements `n`-bit word operations.
//
// Since 0.1.4
type bitWord struct {
	// width is the word width.
	width int
	// byteCap defines how many words a byte contains.
	// It is 8 / width.
	byteCap int
	// wordMask is bit mask.
	// It sets the least siginicant "width" bit to 1.
	wordMask byte
}

func newBW(n int) Interface {
	return &bitWord{
		width:    n,
		byteCap:  8 / n,
		wordMask: (1 << uint(n)) - 1,
	}
}

var (
	// BitWord pre-defines n-bit words operations for 1, 2, 4, 8
	BitWord = map[int]Interface{
		1: newBW(1),
		2: newBW(2),
		4: newBW(4),
		8: newBW(8),
	}
)

// FromStr split a string into a slice of byte.
// A byte in string is split into 8/`n` `n`-bit words
// Value of every byte is in range [0, 2^n-1].
//
// Significant bits in a byte is place at left.
// Thus the result byte slice keeps order with the original string.
//
// Since 0.1.4
func (w *bitWord) FromStr(s string) []byte {

	// number of words per byte
	m := w.byteCap
	lenSrc := len(s)
	words := make([]byte, lenSrc*m)

	for i := 0; i < lenSrc; i++ {
		b := s[i]

		for j := 0; j < m; j++ {
			words[i*m+j] = (b >> uint(8-w.width*j-w.width)) & w.wordMask
		}
	}
	return words
}

// FromStrs converts a `[]string` to a n-bit word `[][]byte`.
//
// Since 0.1.4
func (w *bitWord) FromStrs(strs []string) [][]byte {
	rst := make([][]byte, len(strs))
	for i, s := range strs {
		rst[i] = w.FromStr(s)
	}
	return rst
}

// ToStr is the reverse of FromStr.
//
// Since 0.1.4
func (w *bitWord) ToStr(bs []byte) string {

	// number of words per byte
	m := w.byteCap
	sz := (len(bs) + m - 1) / m
	strbs := make([]byte, sz)

	var b byte
	for i := 0; i < len(strbs); i++ {
		b = 0
		for j := 0; j < m; j++ {
			if i*m+j < len(bs) {
				b = (b << uint(w.width)) + bs[i*m+j]
			} else {
				b = b << uint(w.width)
			}
		}
		strbs[i] = b
	}

	return string(strbs)
}

// ToStrs converts a `[][]byte` back to a `[]string`.
//
// Since 0.1.4
func (w *bitWord) ToStrs(bytesslice [][]byte) []string {
	rst := make([]string, len(bytesslice))
	for i, s := range bytesslice {
		rst[i] = w.ToStr(s)
	}
	return rst
}

// Get returns i-th n-bit word from a string.
//
// Since 0.1.4
func (w *bitWord) Get(s string, ith int) byte {
	i := w.width * ith
	end := (i + w.width - 1) & 7

	word := s[i>>3]
	return (word >> uint(7-end)) & w.wordMask
}

// FirstDiff returns the index of the first different n-bit word,
// that ge "from" and lt "end".
// If "end" is -1 it means to look up upto end of a or b.
//
// Since 0.1.4
func (w *bitWord) FirstDiff(a, b string, from, end int) int {
	la, lb := len(a)*w.byteCap, len(b)*w.byteCap

	if end == -1 {
		end = la
	}

	if end > la {
		end = la
	}

	if end > lb {
		end = lb
	}

	for i := from; i < end; i++ {
		if w.Get(a, i) != w.Get(b, i) {
			return i
		}
	}
	return end
}
