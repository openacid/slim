package trie

import (
	"fmt"
	"testing"

	"github.com/openacid/slim/encode"
	"github.com/openacid/testkeys"
	"github.com/stretchr/testify/require"
)

func TestSlimTrie_withPrefixContent_Get(t *testing.T) {

	ta := require.New(t)

	keys := []string{
		"abc",
		"abcd",
		"abd",
		"abde",
		"bc",
		"bcd",
		"bcde",
		"cde",
	}
	values := []int32{0, 1, 2, 3, 4, 5, 6, 7}
	searches := []struct {
		key  string
		want interface{}
	}{
		{"ab", nil},
		{"abc", int32(0)},
		{"abcde", int32(1)}, // false positive
		{"abd", int32(2)},
		{"ac", nil},
		{"acb", nil},
		{"acd", nil},
		{"adc", nil},
		{"bcd", int32(5)},
		{"bce", nil},
		{"c", int32(7)}, // false positive
		{"cde", int32(7)},
		{"cfe", int32(7)},
		{"cff", int32(7)},
	}

	st, err := NewSlimTrie(encode.I32{}, keys, values,
		Opt{InnerPrefix: Bool(true)})
	ta.NoError(err)

	fmt.Println(st.content())
	fmt.Println(st.String())

	wantstr := trim(`
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
`)

	_ = wantstr
	ta.Equal(wantstr, st.String())

	for _, c := range searches {
		fmt.Println("get:", c.key)
		v, found := st.Get(c.key)
		ta.Equal(c.want != nil, found, "case: %+v", c)
		ta.Equal(c.want, v, "case: %+v", c)
	}
}

func TestSlimTrie_withPrefixContent_Get_small_keyset(t *testing.T) {

	ta := require.New(t)

	for _, typ := range testkeys.AssetNames() {

		keys := getKeys(typ)

		if len(keys) >= 1000 {
			continue
		}

		dd("small keyset: %s", typ)

		values := makeI32s(len(keys))
		st, err := NewSlimTrie(encode.I32{}, keys, values,
			Opt{InnerPrefix: Bool(true)})
		ta.NoError(err)

		dd(st)

		testPresentKeysGRS(t, st, keys, values)
	}
}

func TestSlimTrie_withPrefixContent_GRS_all_keyset(t *testing.T) {

	testBigKeySet(t, func(t *testing.T, keys []string) {
		ta := require.New(t)
		values := makeI32s(len(keys))
		st, err := NewSlimTrie(encode.I32{}, keys, values,
			Opt{InnerPrefix: Bool(true)})
		ta.NoError(err)

		testPresentKeysGRS(t, st, keys, values)
	})
}
