package trie

import (
	"testing"

	"github.com/kr/pretty"
	"github.com/openacid/slim/encode"
	"github.com/stretchr/testify/require"
)

type statCase struct {
	keys    []string
	slimStr string
	stat    string
}

var statCases = map[string]statCase{
	"empty": {
		keys:    []string{},
		slimStr: trim(""),
		stat: trim(`
&trie.Stat{
    LevelCnt: 1,
    Levels:   {
        {},
    },
    KeyCnt:  0,
    NodeCnt: 0,
}
`),
	},
	"singleKey": {
		keys:    []string{"foo"},
		slimStr: trim("#000=0"),
		stat: trim(`
&trie.Stat{
    LevelCnt: 2,
    Levels:   {
        {},
        {Total:1, Inner:0, Leaf:1},
    },
    KeyCnt:  1,
    NodeCnt: 1,
}
`),
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
		stat: trim(`
&trie.Stat{
    LevelCnt: 5,
    Levels:   {
        {},
        {Total:1, Inner:1, Leaf:0},
        {Total:4, Inner:3, Leaf:1},
        {Total:8, Inner:6, Leaf:2},
        {Total:14, Inner:6, Leaf:8},
    },
    KeyCnt:  8,
    NodeCnt: 14,
}
`),
	},
}

func TestSlimTrie_Stat(t *testing.T) {

	for name, c := range statCases {
		t.Run(name, func(t *testing.T) {

			values := makeI32s(len(c.keys))

			t.Run("complete", func(t *testing.T) {

				ta := require.New(t)

				st, err := NewSlimTrie(encode.I32{}, c.keys, values, Opt{Complete: Bool(true)})
				ta.NoError(err)

				dd(st)
				ta.Equal(c.slimStr, st.String())
				ta.Equal(c.stat, pretty.Sprint(st.Stat()))
			})

			t.Run("minimal", func(t *testing.T) {

				ta := require.New(t)

				st, err := NewSlimTrie(encode.I32{}, c.keys, values)
				ta.NoError(err)

				ta.Equal(c.stat, pretty.Sprint(st.Stat()))
			})
		})
	}
}
