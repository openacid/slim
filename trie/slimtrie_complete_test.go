package trie

import (
	"testing"

	"github.com/openacid/slim/encode"
	"github.com/openacid/testkeys"
	"github.com/stretchr/testify/require"
)

func TestSlimTrie_Complete_GRS_0_tiny(t *testing.T) {

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
	values := makeI32s(len(keys))

	st, err := NewSlimTrie(encode.I32{}, keys, values, Opt{Complete: Bool(true)})
	ta.NoError(err)

	// prefix len of #007 is 8 because WithPrefixContent stores prefix aligned
	// to 8 bits
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

	dd(st.content())
	dd(st)
	dd(st.inner)

	ta.Equal(wantstr, st.String())

	testAbsentKeysGRS(t, st, keys)
	testPresentKeysGRS(t, st, keys, values)
}

func TestSlimTrie_Complete_GRS_1_empty_string_branch(t *testing.T) {

	// In Get() the loop must end after i reaches lenWords
	// or it can not find the first key "b".
	//
	// In this case it creates a slimtrie node that just ends at the 8-th bit.

	ta := require.New(t)

	keys := []string{
		"b",
		"ba",
		"cc",
		"dc",
		"pc",
	}
	values := makeI32s(len(keys))

	st, err := NewSlimTrie(encode.I32{}, keys, values, Opt{Complete: Bool(true)})
	ta.NoError(err)

	wantstr := trim(`
#000*2
    -0110->#001*3
               -0010->#003*2
                          -->#006=0
                          -0110->#007=1
               -0011->#004=2
               -0100->#005=3
    -0111->#002=4
`)

	dd(st)

	ta.Equal(wantstr, st.String())

	testAbsentKeysGRS(t, st, keys)
	testPresentKeysGRS(t, st, keys, values)
}

func TestSlimTrie_Complete_GRS_1_zerokeys(t *testing.T) {

	ta := require.New(t)

	keys := []string{}
	values := makeI32s(len(keys))

	st, err := NewSlimTrie(encode.I32{}, keys, values, Opt{Complete: Bool(true)})
	ta.NoError(err)

	ta.Equal("", st.String())

	ks := randVStrings(500, 0, 10)
	for i, key := range ks {
		_ = i

		nid := st.GetID(key)
		ta.Equal(int32(-1), nid)

		v, found := st.Get(key)
		ta.Nil(v)
		ta.False(found)

		v, found = st.RangeGet(key)
		ta.Nil(v)
		ta.False(found)

		lid, eid, rid := st.searchID(key)
		ta.Equal(int32(-1), lid)
		ta.Equal(int32(-1), eid)
		ta.Equal(int32(-1), rid)

		l, e, r := st.Search(key)
		ta.Nil(l)
		ta.Nil(e)
		ta.Nil(r)
	}
}

func TestSlimTrie_Complete_GRS_1_onekey(t *testing.T) {

	ta := require.New(t)

	keys := []string{
		"abc",
	}
	values := makeI32s(len(keys))

	st, err := NewSlimTrie(encode.I32{}, keys, values, Opt{Complete: Bool(true)})
	ta.NoError(err)

	wantstr := trim(`
#000=0
`)
	ta.Equal(wantstr, st.String())

	testAbsentKeysGRS(t, st, keys)
	testPresentKeysGRS(t, st, keys, values)
}

func TestSlimTrie_Complete_GRS_1_twokeys_emptysteps(t *testing.T) {

	ta := require.New(t)

	// the first bit diffs, thus no step needed
	keys := []string{
		"abc",
		"\x80bc",
	}
	values := makeI32s(len(keys))

	st, err := NewSlimTrie(encode.I32{}, keys, values, Opt{Complete: Bool(true)})
	ta.NoError(err)

	wantstr := trim(`
#000*2
    -0110->#001=0
    -1000->#002=1
`)
	dd(st)

	ta.Equal(wantstr, st.String())

	testAbsentKeysGRS(t, st, keys)
	testPresentKeysGRS(t, st, keys, values)
}

func TestSlimTrie_Complete_GRS_2_small_keyset(t *testing.T) {

	ta := require.New(t)

	for _, typ := range testkeys.AssetNames() {

		keys := getKeys(typ)

		if len(keys) >= 1000 {
			continue
		}

		dd("small keyset: %s", typ)

		values := makeI32s(len(keys))
		st, err := NewSlimTrie(encode.I32{}, keys, values, Opt{Complete: Bool(true)})
		ta.NoError(err)

		// dd(st)

		testAbsentKeysGRS(t, st, keys)
		testPresentKeysGet(t, st, keys, values)
	}
}

func TestSlimTrie_Complete_GRS_3_bigInner_300(t *testing.T) {

	ta := require.New(t)
	keys := getKeys("300vl50")
	values := makeI32s(len(keys))

	st, err := NewSlimTrie(encode.I32{}, keys, values, Opt{Complete: Bool(true)})
	ta.NoError(err)

	// dd(st)

	ta.True(st.inner.BigInnerCnt > 0)

	testAbsentKeysGRS(t, st, keys)
	testPresentKeysGRS(t, st, keys, values)
}

func TestSlimTrie_Complete_GRS_3_bigInner_a2t(t *testing.T) {

	ta := require.New(t)
	keys := []string{
		"a",
		"b",
		"c",
		"d",
		"e",
		"f",
		"g",
		"h",
		"i",
		"j",
		"k",
		"l",
		"m",
		"n",
		"o",
		"p",
		"q",
		"r",
		"s",
		"t",
	}
	values := makeI32s(len(keys))

	st, err := NewSlimTrie(encode.I32{}, keys, values, Opt{Complete: Bool(true)})
	ta.NoError(err)

	dd(st)

	ta.True(st.inner.BigInnerCnt > 0)

	testAbsentKeysGRS(t, st, keys)
	testPresentKeysGRS(t, st, keys, values)
}

func TestSlimTrie_Complete_GRS_4_20kvlen10(t *testing.T) {

	iambig(t)

	ta := require.New(t)

	keys := getKeys("20kvl10")
	values := makeI32s(len(keys))
	st, err := NewSlimTrie(encode.I32{}, keys, values, Opt{Complete: Bool(true)})
	ta.NoError(err)

	// dd(st)

	testAbsentKeysGRS(t, st, keys)
	testPresentKeysGRS(t, st, keys, values)
}

func TestSlimTrie_Complete_GRS_9_allkeyset(t *testing.T) {
	testBigKeySet(t, func(t *testing.T, keys []string) {
		ta := require.New(t)
		values := makeI32s(len(keys))
		st, err := NewSlimTrie(encode.I32{}, keys, values, Opt{Complete: Bool(true)})
		ta.NoError(err)

		testAbsentKeysGRS(t, st, keys)
		testPresentKeysGRS(t, st, keys, values)

	})
}
