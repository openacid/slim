package trie

import (
	"sort"
	"testing"

	"github.com/openacid/slim/encode"
	"github.com/openacid/testutil"
	"github.com/stretchr/testify/require"
)

var defaultScan = []string{
	"",
	"`",
	"a",
	"ab",
	"abc",
	"abca",
	"abcd",
	"abcd1",
	"abce",
	"be",
	"c",
	"cde0",
	"d",
}
var iterCases = map[string]struct {
	keys         []string
	slimStr      string
	paths        [][]int32
	scanFromKeys []string
}{
	"empty": {
		keys:         []string{},
		slimStr:      trim(""),
		paths:        [][]int32{{}},
		scanFromKeys: defaultScan,
	},
	"single": {
		keys:         []string{"ab"},
		slimStr:      trim(`#000=0`),
		paths:        [][]int32{{0}, {}},
		scanFromKeys: defaultScan,
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
		paths: [][]int32{
			{0, 1, 4, 8},
			{0, 1, 4, 9},
			{0, 1, 5, 10},
			{0, 1, 5, 11},
			{0, 2, 6},
			{0, 2, 7, 12},
			{0, 2, 7, 13},
			{0, 3},
			{}, // path seeking from after the last key
		},
		scanFromKeys: defaultScan,
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
		paths: [][]int32{
			{0, 1},
			{0, 2, 3, 6},
			{0, 2, 3, 7, 10},
			{0, 2, 3, 7, 11},
			{0, 2, 4, 8},
			{0, 2, 4, 9},
			{0, 2, 5},
			{}, // path seeking from after the last key
		},
		scanFromKeys: defaultScan,
	},
}

func TestSlimTrie_Iter(t *testing.T) {

	for name, c := range iterCases {
		t.Run(name, func(t *testing.T) {

			ta := require.New(t)

			values := makeI32s(len(c.keys))

			st, err := NewSlimTrie(encode.I32{}, c.keys, values, Opt{Complete: Bool(true)})
			ta.NoError(err)

			dd(st)
			ta.Equal(c.slimStr, st.String())

			subTestPath(t, st, c.keys, c.paths, c.keys)
			subTestPath(t, st, c.keys, c.paths, c.scanFromKeys)
			subTestIter(t, st, c.keys, c.keys)
			subTestIter(t, st, c.keys, c.scanFromKeys)
			subTestIter(t, st, c.keys, testutil.RandStrSlice(len(c.keys)*5, 0, 10))

			subTestScan(t, st, c.keys, c.keys)
			subTestScan(t, st, c.keys, c.scanFromKeys)
			subTestScan(t, st, c.keys, testutil.RandStrSlice(len(c.keys)*5, 0, 10))
		})
	}
}

func TestSlimTrie_NewIter_panic(t *testing.T) {

	ta := require.New(t)
	keys := iterCases["simple"].keys
	values := makeI32s(len(keys))

	ta.Panics(func() {
		st, err := NewSlimTrie(encode.I32{}, keys, values, Opt{InnerPrefix: Bool(true)})
		ta.NoError(err)
		st.NewIter("abc", true, true)
	}, "without leaf prefix")

	ta.Panics(func() {
		st, err := NewSlimTrie(encode.I32{}, keys, values)
		ta.NoError(err)
		st.NewIter("abc", true, true)
	}, "without inner prefix")
}

func TestSlimTrie_NewIter_slimWithoutValue(t *testing.T) {

	ta := require.New(t)

	c := iterCases["simple"]
	keys := c.keys

	st, err := NewSlimTrie(encode.I32{}, keys, nil, Opt{Complete: Bool(true)})
	ta.NoError(err)

	for _, sk := range c.scanFromKeys {
		idx := sort.SearchStrings(keys, sk)
		nxt := st.NewIter(sk, true, true)

		for i := int32(idx); i < int32(len(keys)); i++ {
			key := keys[i]
			gotKey, gotVal := nxt()
			ta.Equal([]byte(key), gotKey, "newIter from: %s %v, idx: %d", sk, []byte(sk), idx)
			ta.Nil(gotVal, "newIter from: %s %v, idx: %d", sk, []byte(sk), idx)
		}
		gotKey, gotVal := nxt()
		ta.Nil(gotKey)
		ta.Nil(gotVal)
	}
}

func TestSlimTrie_Iter_large(t *testing.T) {

	testBigKeySet(t, func(t *testing.T, typ string, keys []string) {
		ta := require.New(t)

		values := makeI32s(len(keys))

		st, err := NewSlimTrie(encode.I32{}, keys, values, Opt{Complete: Bool(true)})
		ta.NoError(err)

		subTestIter(t, st, keys, testutil.RandStrSlice(clap(len(keys), 50, 10*1024), 0, 10))
		subTestScan(t, st, keys, testutil.RandStrSlice(clap(len(keys), 50, 10*1024), 0, 10))
	})
}

var OutputIter int

func BenchmarkSlimTrie_Iter(b *testing.B) {

	typ := "1mvl5_10"

	keys := getKeys(typ)
	n := len(keys)
	values := makeI32s(n)

	st, err := NewSlimTrie(encode.I32{}, keys, values, Opt{Complete: Bool(true)})
	if err != nil {
		panic("err:" + err.Error())
	}

	scanN := 1024 * 100

	b.ResetTimer()

	s := 0

	for i := 0; i < b.N/scanN; i++ {
		nxt := st.NewIter("`", true, true)
		for j := 0; j < scanN; j++ {
			b, _ := nxt()
			s += int(b[0])

		}
	}
	OutputIter = s
}

func subTestPath(
	t *testing.T,
	st *SlimTrie,
	keys []string,
	paths [][]int32,
	scanFromKeys []string,
) {

	t.Run("getGEPath", func(t *testing.T) {

		ta := require.New(t)

		// searching from other keys should start from next present key.
		for _, sk := range scanFromKeys {
			idx := sort.SearchStrings(keys, sk)
			p, gotEqual := st.getGEPath(sk)
			ta.Equal(paths[idx], p, "key: %s %v", sk, []byte(sk))
			if idx == len(keys) {
				ta.False(gotEqual)
			} else {
				ta.Equal(keys[idx] == sk, gotEqual, "key: %s %v, gotEqual: %v", sk, []byte(sk), gotEqual)
			}
		}
	})

}
func subTestIter(
	t *testing.T,
	st *SlimTrie,
	keys []string,
	scanFromKeys []string,
) {
	t.Run("NewIter", func(t *testing.T) {

		ta := require.New(t)

		for _, sk := range scanFromKeys {
			idx := sort.SearchStrings(keys, sk)

			{ // test exclusive
				nxt := st.NewIter(sk, false, true)
				gotKey, _ := nxt()
				if idx == len(keys) {
					ta.Nil(gotKey, "newIter from: %s %v, idx: %d", sk, []byte(sk), idx)
				} else {
					if keys[idx] == sk {
						ta.NotEqual([]byte(keys[idx]), gotKey, "newIter from: %s %v, idx: %d", sk, []byte(sk), idx)
					} else {
						ta.Equal([]byte(keys[idx]), gotKey, "newIter from: %s %v, idx: %d", sk, []byte(sk), idx)
					}
				}
			}

			nxt := st.NewIter(sk, true, true)

			var i int32
			for i = int32(idx); i < int32(len(keys)) && i < int32(idx+200); i++ {
				key := keys[i]
				gotKey, gotVal := nxt()
				ta.Equal([]byte(key), gotKey, "newIter from: %s %v, idx: %d", sk, []byte(sk), idx)
				ta.Equal(st.encoder.Encode(i), gotVal, "newIter from: %s %v, idx: %d", sk, []byte(sk), idx)
			}
			if i == int32(len(keys)) {
				ta.Nil(nxt())
			}

			{ // newIter without yielding value
				nxt := st.NewIter(sk, true, false)
				_, gotVal := nxt()
				ta.Nil(gotVal)
				_, gotVal = nxt()
				ta.Nil(gotVal)
			}
		}
	})
}

func subTestScan(
	t *testing.T,
	st *SlimTrie,
	keys []string,
	scanFromKeys []string,
) {
	t.Run("NewIter", func(t *testing.T) {

		ta := require.New(t)

		for _, sk := range scanFromKeys {
			idx := sort.SearchStrings(keys, sk)

			{ // exclusive start
				st.ScanFrom(sk, false, true, func(gotKey, v []byte) bool {
					if idx == len(keys) {
						ta.Nil(gotKey, "Scan from: %s %v, idx: %d", sk, []byte(sk), idx)
					} else {
						if keys[idx] == sk {
							ta.NotEqual([]byte(keys[idx]), gotKey, "Scan from: %s %v, idx: %d", sk, []byte(sk), idx)
						} else {
							ta.Equal([]byte(keys[idx]), gotKey, "Scan from: %s %v, idx: %d", sk, []byte(sk), idx)
						}
					}

					// only check the first key
					return false
				})
			}

			{ // inclusive start

				i := int32(idx)
				st.ScanFrom(sk, true, true, func(gotKey, gotVal []byte) bool {
					key := keys[i]
					ta.Equal([]byte(key), gotKey, "Scan from: %s %v, idx: %d", sk, []byte(sk), idx)
					ta.Equal(st.encoder.Encode(i), gotVal, "Scan from: %s %v, idx: %d", sk, []byte(sk), idx)

					i++
					return i <= int32(idx+200)
				})
			}

			{ // scan from to

				i := int32(idx)
				endIdx := i + 50
				if endIdx >= int32(len(keys)) {
					endIdx = int32(len(keys)) - 1
				}

				var endKey string
				if endIdx == -1 {
					endKey = "foo"
				} else {
					endKey = keys[endIdx]
				}

				st.ScanFromTo(sk, true,
					endKey, false,
					true, func(gotKey, gotVal []byte) bool {
						key := keys[i]
						ta.Equal([]byte(key), gotKey, "Scan from: %s %v, idx: %d", sk, []byte(sk), idx)
						ta.Equal(st.encoder.Encode(i), gotVal, "Scan from: %s %v, idx: %d", sk, []byte(sk), idx)
						ta.True(string(gotKey) < endKey, "Scan from: %s %v, idx: %d", sk, []byte(sk), idx)

						i++
						return i <= int32(idx+200)
					})
			}

			{ // without yielding value
				st.ScanFrom(sk, true, false, func(k, v []byte) bool {
					ta.Nil(v)
					return false
				})
			}
		}
	})
}
