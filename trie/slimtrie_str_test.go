package trie

import (
	"testing"

	"github.com/openacid/slim/encode"
	"github.com/openacid/testutil"
	"github.com/stretchr/testify/require"
)

func TestSlimTrie_String_empty(t *testing.T) {

	ta := require.New(t)

	keys := []string{}
	values := makeI32s(len(keys))
	st, err := NewSlimTrie(encode.I32{}, keys, values)
	ta.NoError(err)

	want := trim(`
`)
	dd(st)

	ta.Equal(want, st.String())

	testUnknownKeysGRS(t, st, testutil.RandStrSlice(100, 0, 10))
	testPresentKeysGet(t, st, keys, values)
}

func TestSlimTrie_String(t *testing.T) {

	ta := require.New(t)

	keys := []string{
		"abc",
		"abcd",
		"abcdx",
		"abcdy",
		"abcdz",
		"abd",
		"abde",
		"bc",
		"bcd",
		"bcde",
		"cde",
	}
	values := makeI32s(len(keys))
	st, err := NewSlimTrie(encode.I32{}, keys, values)
	ta.NoError(err)

	dd(st)

	want := trim(`
#000+4*3
    -0001->#001+12*2
               -0011->#004*2
                          -->#008=0
                          -0110->#009+4*2
                                     -->#014=1
                                     -0111->#015*3
                                                -1000->#016=2
                                                -1001->#017=3
                                                -1010->#018=4
               -0100->#005*2
                          -->#010=5
                          -0110->#011=6
    -0010->#002+8*2
               -->#006=7
               -0110->#007+4*2
                          -->#012=8
                          -0110->#013=9
    -0011->#003=10
`)

	ta.Equal(want, st.String())

	testPresentKeysGet(t, st, keys, values)
}
