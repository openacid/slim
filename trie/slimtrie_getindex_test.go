package trie

import (
	"sort"
	"testing"

	"github.com/openacid/low/mathext/zipf"
	"github.com/openacid/slim/encode"
	"github.com/stretchr/testify/require"
)

var OutputIndex int32

type indexCase struct {
	keys    []string
	slimStr string
}

var indexCases = map[string]indexCase{
	"empty": {
		keys:    []string{},
		slimStr: trim(""),
	},
	"singleKey": {
		keys:    []string{"foo"},
		slimStr: trim("#000=0"),
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
	},
	"earlyLeaf": {
		keys: []string{
			"abcd1",
			"abcd2x1",
			"abcd2x2",
			"abcd2y",
			"abd",
			"abde",
			"bc",
			"bcd",
			"bcde",
			"cde",
			"def",
		},
		slimStr: trim(`
#000+4*4
    -0001->#001+12*2
               -0011->#005+12*2
                          -0001->#009=0
                          -0010->#010+4*2
                                     -1000->#015+4*2
                                                -0001->#017=1
                                                -0010->#018=2
                                     -1001->#016=3
               -0100->#006*2
                          -->#011=4
                          -0110->#012=5
    -0010->#002+8*2
               -->#007=6
               -0110->#008+8*2
                          -->#013=7
                          -0110->#014=8
    -0011->#003=9
    -0100->#004=10
`),
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
	},
}

func TestSlimTrie_getIndex(t *testing.T) {

	for name, c := range indexCases {
		t.Run(name, func(t *testing.T) {

			values := makeI32s(len(c.keys))

			t.Run("complete", func(t *testing.T) {

				ta := require.New(t)

				st, err := NewSlimTrie(encode.I32{}, c.keys, values, Opt{Complete: Bool(true)})
				ta.NoError(err)

				dd(st)
				ta.Equal(c.slimStr, st.String())

				absentN := len(c.keys) * 50

				subGetIndex(t, c, st)
				subGetIndexAbsent(t, c, st, absentN)

				subGetLRIndex(t, c, st)
				subGetLRIndexAbsent(t, c, st, absentN)
			})

			t.Run("minimal", func(t *testing.T) {

				ta := require.New(t)

				st, err := NewSlimTrie(encode.I32{}, c.keys, values)
				ta.NoError(err)

				subGetIndex(t, c, st)
				subGetLRIndex(t, c, st)
			})
		})
	}
}

func TestSlimTrie_getIndex_large(t *testing.T) {

	testBigKeySet(t, func(t *testing.T, typ string, keys []string) {

		ta := require.New(t)

		c := indexCase{
			keys:    keys,
			slimStr: "",
		}

		values := makeI32s(len(keys))

		st, err := NewSlimTrie(encode.I32{}, keys, values, Opt{Complete: Bool(true)})
		ta.NoError(err)

		absentN := clap(len(keys)*50, 50, len(keys)*2)

		subGetIndex(t, c, st)
		subGetIndexAbsent(t, c, st, absentN)
		subGetLRIndex(t, c, st)
		subGetLRIndexAbsent(t, c, st, absentN)
	})
}

func subGetIndex(t *testing.T, c indexCase, st *SlimTrie) {
	t.Run("GetIndex", func(t *testing.T) {
		ta := require.New(t)

		for i, k := range c.keys {
			ta.Equal(int32(i), st.GetIndex(k), "key: %s", k)
		}
	})
}

func subGetIndexAbsent(t *testing.T, c indexCase, st *SlimTrie, n int) {
	t.Run("GetIndexAbsent", func(t *testing.T) {
		ta := require.New(t)

		absentKeys := makeAbsentKeys(c.keys, n, 0, 20)
		for i, key := range absentKeys {
			_ = i

			v := st.GetIndex(key)
			ta.Equal(int32(-1), v, "absent key: %s", key)
		}
	})
}

func subGetLRIndex(t *testing.T, c indexCase, st *SlimTrie) {
	t.Run("GetLRIndex", func(t *testing.T) {
		ta := require.New(t)

		for i, k := range c.keys {
			l, e := st.GetLRIndex(k)
			ta.Equal(int32(i), l, "key: %s", k)
			ta.Equal(int32(i), e, "key: %s", k)
		}
	})
}

func subGetLRIndexAbsent(t *testing.T, c indexCase, st *SlimTrie, n int) {
	t.Run("GetLRIndexAbsent", func(t *testing.T) {
		ta := require.New(t)

		absentKeys := makeAbsentKeys(c.keys, n, 0, 20)
		for i, k := range absentKeys {
			_ = i
			idx := sort.SearchStrings(c.keys, k)
			l, e := st.GetLRIndex(k)
			ta.Equal(int32(idx-1), l, "key: %s, %v", k, []byte(k))
			ta.Equal(int32(idx), e, "key: %s, %v", k, []byte(k))
		}
	})
}

func BenchmarkSlimTrie_GetIndex_big(b *testing.B) {

	benchBigKeySet(b, func(b *testing.B, typ string, keys []string) {

		st, err := NewSlimTrie(encode.I32{}, keys, nil, Opt{InnerPrefix: Bool(true)})
		if err != nil {
			panic("˙∆˙...")
		}

		accesses := zipf.Accesses(2, 1.5, len(keys), b.N, nil)

		b.ResetTimer()

		var id int32
		for i := 0; i < b.N; i++ {
			id += st.GetIndex(keys[accesses[i]])
		}
		OutputIndex = id

	})
}
