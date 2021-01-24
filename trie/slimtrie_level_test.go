package trie

import (
	"testing"

	"github.com/openacid/slim/encode"
	"github.com/stretchr/testify/require"
)

var (
	levelCases = map[string]struct {
		keys    []string
		slimStr string
		levels  []levelInfo
	}{
		"empty": {
			keys:    []string{},
			slimStr: trim(""),
			levels:  []levelInfo{{0, 0, 0, nil}},
		},
		"singleKey": {
			keys:    []string{"foo"},
			slimStr: trim("#000=0"),
			levels: []levelInfo{
				{0, 0, 0, nil},
				{1, 0, 1, nil},
			},
		},
		"simple": {
			keys: []string{
				"abc",
				"abcd",
				"abd",
				"abde",
				"bc",
				"bcd",
				"bcde",
				"cde",
			},
			slimStr: trim(`
#000+4*3
    -0001->#001+12*2
               -0011->#004*2
                          -->#008=0
                          -0110->#009=1
               -0100->#005*2
                          -->#010=2
                          -0110->#011=3
    -0010->#002+8*2
               -->#006=4
               -0110->#007+8*2
                          -->#012=5
                          -0110->#013=6
    -0011->#003=7
`),

			levels: []levelInfo{
				{0, 0, 0, nil},
				{1, 1, 0, nil},
				{4, 3, 1, nil},
				{8, 6, 2, nil},
				{14, 6, 8, nil},
			},
		},
		"emptyKey": {
			keys: []string{
				"",
				"a",
				"abc",
				"abd",
				"bc",
				"bcd",
				"cde",
			},
			slimStr: trim(`
#000*2
    -->#001=0
    -0110->#002*3
               -0001->#003*2
                          -->#006=1
                          -0110->#007+12*2
                                     -0011->#010=2
                                     -0100->#011=3
               -0010->#004+8*2
                          -->#008=4
                          -0110->#009=5
               -0011->#005=6
`),
			levels: []levelInfo{
				{0, 0, 0, nil},
				{1, 1, 0, nil},
				{3, 2, 1, nil},
				{6, 4, 2, nil},
				{10, 5, 5, nil},
				{12, 5, 7, nil},
			},
		},
	}
)

func TestSlimTrie_levels(t *testing.T) {

	for name, c := range levelCases {
		t.Run(name, func(t *testing.T) {
			ta := require.New(t)

			values := makeI32s(len(c.keys))
			st, err := NewSlimTrie(encode.I32{}, c.keys, values, Opt{Complete: Bool(true)})
			ta.NoError(err)

			dd(st)
			ta.Equal(c.slimStr, st.String())

			for i, lvl := range c.levels {
				ta.Equal(lvl.total, st.levels[i].total, "total: line %d", i)
				ta.Equal(lvl.inner, st.levels[i].inner, "inner: line %d", i)
				ta.Equal(lvl.leaf, st.levels[i].leaf, "leaf: line %d", i)
			}
		})
	}
}
