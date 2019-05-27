// Package strhelper provides string operations functions
package strhelper

import (
	"fmt"
	"math/bits"
	"reflect"
	"strings"

	"github.com/openacid/low/bitword"
)

// ToBitWords split a string into a slice of byte.
// A char in string is split into 8/`n` `n`-bit words
// Value of every byte is in range [0, 2^n-1].
// `n` must be a one of [1, 2, 4, 8].
//
// Significant bits in a char is place at left.
// Thus the result byte slice keeps order with the original string.
//
// Deprecated: It will be removed since 1.0.0 .
// Use "github.com/openacid/low/bitword"
func ToBitWords(s string, n int) []byte {
	return bitword.BitWord[n].FromStr(s)
}

// SliceToBitWords converts a `[]string` to a n-bit word `[][]byte`.
//
// Deprecated: It will be removed since 1.0.0 .
// Use "github.com/openacid/low/bitword"
func SliceToBitWords(strs []string, n int) [][]byte {
	return bitword.BitWord[n].FromStrs(strs)
}

// FromBitWords is the reverse of ToBitWords.
// It composes a string of which each byte is formed from 8/n words from bs.
//
// Deprecated: It will be removed since 1.0.0 .
// Use "github.com/openacid/low/bitword"
func FromBitWords(bs []byte, n int) string {
	return bitword.BitWord[n].ToStr(bs)
}

// SliceFromBitWords converts a `[][]byte` back to a `[]string`.
//
// Deprecated: It will be removed since 1.0.0 .
// Use "github.com/openacid/low/bitword"
func SliceFromBitWords(bytesslice [][]byte, n int) []string {
	return bitword.BitWord[n].ToStrs(bytesslice)
}

// ToBin converts integer or slice of integer to binary format string.
// Significant bits are at right.
// E.g.:
//    int32(0x0102) --> 01000000 10000000
//
// Since 0.5.4
func ToBin(x interface{}) string {

	rst := []string{}

	v := reflect.ValueOf(x)
	if v.Kind() == reflect.Slice {
		n := v.Len()
		for i := 0; i < n; i++ {
			rst = append(rst, intToBin(v.Index(i).Interface()))
		}
		return strings.Join(rst, ",")

	} else {
		return intToBin(x)
	}
}

func intToBin(i interface{}) string {

	sz, v := intSize(i)

	rst := []string{}
	for i := 0; i < sz; i++ {
		b := uint8(v >> uint(i*8))
		s := fmt.Sprintf("%08b", bits.Reverse8(b))
		rst = append(rst, s)
	}
	return strings.Join(rst, " ")
}

func intSize(i interface{}) (int, uint64) {
	switch i := i.(type) {
	case int8:
		return 1, uint64(i)
	case uint8:
		return 1, uint64(i)
	case int16:
		return 2, uint64(i)
	case uint16:
		return 2, uint64(i)
	case int32:
		return 4, uint64(i)
	case uint32:
		return 4, uint64(i)
	case int64:
		return 8, uint64(i)
	case uint64:
		return 8, i
	}

	panic(fmt.Sprintf("not a int type: %v", i))
}
