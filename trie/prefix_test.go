package trie

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewPrefix(t *testing.T) {

	ta := require.New(t)

	cases := []struct {
		str       string
		s, e      int32
		wantbytes string
	}{
		{"abc", 0, 0, "\x00"},
		{"abc", 0, 1, "\x01\x40"},
		{"abc", 0, 2, "\x01\x60"},
		{"abc", 0, 3, "\x01\x70"},
		{"abc", 0, 4, "\x01\x68"},
		{"abc", 0, 5, "\x01\x64"},
		{"abc", 0, 6, "\x01\x62"},
		{"abc", 0, 7, "\x01\x61"},
		{"abc", 0, 8, "\x00\x61"},
		{"abc", 0, 9, "\x01\x61\x40"},
		{"abc", 0, 10, "\x01\x61\x60"},
		{"abc", 0, 11, "\x01\x61\x70"},
		{"abc", 0, 12, "\x01\x61\x68"},
		{"abc", 0, 13, "\x01\x61\x64"},
		{"abc", 0, 14, "\x01\x61\x62"},
		{"abc", 0, 15, "\x01\x61\x63"},
		{"abc", 0, 16, "\x00\x61\x62"},

		{"abc", 8, 8, "\x00"},
		{"abc", 8, 9, "\x01\x40"},
		{"abc", 8, 10, "\x01\x60"},
		{"abc", 8, 11, "\x01\x70"},
		{"abc", 8, 12, "\x01\x68"},
		{"abc", 8, 13, "\x01\x64"},
		{"abc", 8, 14, "\x01\x62"},
		{"abc", 8, 15, "\x01\x63"},
		{"abc", 8, 16, "\x00\x62"},
	}

	for i, c := range cases {
		gotbytes := newPrefix(c.str, c.s, c.e)
		if len(c.wantbytes) > 0 {
			dd("%08b %08b", c.wantbytes[0], gotbytes[0])
		}
		ta.Equal(c.wantbytes, string(gotbytes), "%d-th: case: %+v", i+1, c)

		// ta.Equal(0, prefixCompare([]byte(c.str[c.s>>3:]), gotbytes))
		ta.Equal(0, prefixCompare(c.str[c.s>>3:], gotbytes))

		if (c.e - c.s) > 0 {
			ta.Equal(c.e-c.s, prefixLen(gotbytes))
		}
	}
}

func TestPrefixCompare(t *testing.T) {

	ta := require.New(t)

	cases := []struct {
		key  string
		src  string
		s, e int32
		want int
	}{
		// original str
		{"abc", "abc", 0, 0, 0},
		{"abc", "abc", 0, 1, 0},
		{"abc", "abc", 0, 2, 0},
		{"abc", "abc", 0, 3, 0},
		{"abc", "abc", 0, 4, 0},
		{"abc", "abc", 0, 5, 0},
		{"abc", "abc", 0, 6, 0},
		{"abc", "abc", 0, 7, 0},
		{"abc", "abc", 0, 8, 0},
		{"abc", "abc", 0, 9, 0},
		{"abc", "abc", 0, 10, 0},
		{"abc", "abc", 0, 11, 0},
		{"abc", "abc", 0, 12, 0},
		{"abc", "abc", 0, 13, 0},
		{"abc", "abc", 0, 14, 0},
		{"abc", "abc", 0, 15, 0},
		{"abc", "abc", 0, 16, 0},

		{"bc", "abc", 8, 8, 0},
		{"bc", "abc", 8, 9, 0},
		{"bc", "abc", 8, 10, 0},
		{"bc", "abc", 8, 11, 0},
		{"bc", "abc", 8, 12, 0},
		{"bc", "abc", 8, 13, 0},
		{"bc", "abc", 8, 14, 0},
		{"bc", "abc", 8, 15, 0},
		{"bc", "abc", 8, 16, 0},

		// empty str
		{"", "abc", 0, 0, 0},
		{"", "abc", 0, 1, -1},
		{"", "abc", 0, 2, -1},
		{"", "abc", 0, 3, -1},
		{"", "abc", 0, 4, -1},
		{"", "abc", 0, 5, -1},
		{"", "abc", 0, 6, -1},
		{"", "abc", 0, 7, -1},
		{"", "abc", 0, 8, -1},
		{"", "abc", 0, 9, -1},
		{"", "abc", 0, 10, -1},
		{"", "abc", 0, 11, -1},
		{"", "abc", 0, 12, -1},
		{"", "abc", 0, 13, -1},
		{"", "abc", 0, 14, -1},
		{"", "abc", 0, 15, -1},
		{"", "abc", 0, 16, -1},

		{"", "abc", 8, 8, 0},
		{"", "abc", 8, 9, -1},
		{"", "abc", 8, 10, -1},
		{"", "abc", 8, 11, -1},
		{"", "abc", 8, 12, -1},
		{"", "abc", 8, 13, -1},
		{"", "abc", 8, 14, -1},
		{"", "abc", 8, 15, -1},
		{"", "abc", 8, 16, -1},

		// smaller str
		{"abc", "bcd", 0, 0, 0},
		{"abc", "bcd", 0, 1, 0},
		{"abc", "bcd", 0, 2, 0},
		{"abc", "bcd", 0, 3, 0},
		{"abc", "bcd", 0, 4, 0},
		{"abc", "bcd", 0, 5, 0},
		{"abc", "bcd", 0, 6, 0},
		{"abc", "bcd", 0, 7, -1},
		{"abc", "bcd", 0, 8, -1},
		{"abc", "bcd", 0, 9, -1},
		{"abc", "bcd", 0, 10, -1},
		{"abc", "bcd", 0, 11, -1},
		{"abc", "bcd", 0, 12, -1},
		{"abc", "bcd", 0, 13, -1},
		{"abc", "bcd", 0, 14, -1},
		{"abc", "bcd", 0, 15, -1},
		{"abc", "bcd", 0, 16, -1},

		{"bc", "bcd", 8, 8, 0},
		{"bc", "bcd", 8, 9, 0},
		{"bc", "bcd", 8, 10, 0},
		{"bc", "bcd", 8, 11, 0},
		{"bc", "bcd", 8, 12, 0},
		{"bc", "bcd", 8, 13, 0},
		{"bc", "bcd", 8, 14, 0},
		{"bc", "bcd", 8, 15, 0},
		{"bc", "bcd", 8, 16, -1},

		// greater str
		{"bcd", "abc", 0, 0, 0},
		{"bcd", "abc", 0, 1, 0},
		{"bcd", "abc", 0, 2, 0},
		{"bcd", "abc", 0, 3, 0},
		{"bcd", "abc", 0, 4, 0},
		{"bcd", "abc", 0, 5, 0},
		{"bcd", "abc", 0, 6, 0},
		{"bcd", "abc", 0, 7, 1},
		{"bcd", "abc", 0, 8, 1},
		{"bcd", "abc", 0, 9, 1},
		{"bcd", "abc", 0, 10, 1},
		{"bcd", "abc", 0, 11, 1},
		{"bcd", "abc", 0, 12, 1},
		{"bcd", "abc", 0, 13, 1},
		{"bcd", "abc", 0, 14, 1},
		{"bcd", "abc", 0, 15, 1},
		{"bcd", "abc", 0, 16, 1},

		{"cde", "abc", 8, 8, 0},
		{"cde", "abc", 8, 9, 0},
		{"cde", "abc", 8, 10, 0},
		{"cde", "abc", 8, 11, 0},
		{"cde", "abc", 8, 12, 0},
		{"cde", "abc", 8, 13, 0},
		{"cde", "abc", 8, 14, 0},
		{"cde", "abc", 8, 15, 0},
		{"cde", "abc", 8, 16, 1},
	}

	for i, c := range cases {
		pref := newPrefix(c.src, c.s, c.e)
		// got := prefixCompare([]byte(c.key), pref, conf)
		got := prefixCompare(c.key, pref)
		ta.Equal(c.want, got, "%d-th: case: %+v", i+1, c)
	}
}

var OutputPrefixCompare int

func BenchmarkPrefixCompare_EqLen_3(b *testing.B) {

	var s int
	pref := newPrefix("abcdef", 0, 22)

	bs := "abc"

	for i := 0; i < b.N; i++ {
		s += prefixCompare(bs, pref)
	}

	OutputPrefixCompare = s
}

func BenchmarkPrefixCompare_EqLen_8(b *testing.B) {

	var s int
	pref := newPrefix("abcdefghijk", 0, 60)

	bs := "abcdefghi"

	for i := 0; i < b.N; i++ {
		s += prefixCompare(bs, pref)
	}

	OutputPrefixCompare = s
}

func BenchmarkPrefixCompare_EqLen_17(b *testing.B) {

	var s int
	pref := newPrefix("abcdefghijkabcdefghijk", 0, 161)

	bs := "abcdefghijkabcdef"

	for i := 0; i < b.N; i++ {
		s += prefixCompare(bs, pref)
	}

	OutputPrefixCompare = s
}

func BenchmarkPrefixCompare_LTLen_5(b *testing.B) {

	var s int
	pref := newPrefix("abcdefghijk", 0, 60)

	bs := "abcde"

	for i := 0; i < b.N; i++ {
		s += prefixCompare(bs, pref)
	}

	OutputPrefixCompare = s
}
