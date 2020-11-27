package trie

import (
	"github.com/openacid/slim/encode"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSlimTrie_Opt_DedupValue(t *testing.T) {

	ta := require.New(t)

	// case: a key not to keep and is a leaf: "Al"
	keys := []string{
		"Aaron",
		"Agatha",
		"Al",
		"Albert",

		"Alexander",
		"Alison",
	}
	values := []int32{
		0, 0, 0, 0,
		1, 1,
	}

	wantDedupedStr := trim(`
#000+12*2
    -0001->#001=0
    -1100->#002
               -0110->#003
                          -0101->#004=1
`)
	wantNonDedupedStr := trim(`
#000+12*3
    -0001->#001=0
    -0111->#002=0
    -1100->#003*2
               -->#004=0
               -0110->#005*3
                          -0010->#006=0
                          -0101->#007=1
                          -1001->#008=1
`)

	{ // default: dedup
		st, err := NewSlimTrie(encode.I32{}, keys, values, Opt{})
		ta.NoError(err)

		dd(keys)
		dd(st)

		ta.Equal(wantDedupedStr, st.String())

		testPresentKeysRangeGet(t, st, keys, values)
	}
	{ // dedup
		st, err := NewSlimTrie(encode.I32{}, keys, values, Opt{DedupValue: Bool(true)})
		ta.NoError(err)

		dd(keys)
		dd(st)

		ta.Equal(wantDedupedStr, st.String())

		testPresentKeysRangeGet(t, st, keys, values)
	}
	{ // no-dedup
		st, err := NewSlimTrie(encode.I32{}, keys, values, Opt{DedupValue: Bool(false)})
		ta.NoError(err)

		dd(keys)
		dd(st)

		ta.Equal(wantNonDedupedStr, st.String())

		testPresentKeysGRS(t, st, keys, values)
	}
}
