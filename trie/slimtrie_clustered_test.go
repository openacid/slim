package trie

import (
	"fmt"
	"testing"

	"github.com/openacid/slim/encode"
	"github.com/openacid/testutil"
	"github.com/stretchr/testify/require"
)

var (
	clusteredKeysSimple = []string{
		"abc",
		"abcd",
		"abd",
		"abde",
		"bc",
		"bcd",
		"bcde",
		"cde",
	}

	clusteredCases = map[string]struct {
		keys          []string
		maxLevel      int32
		slimStr       string
		wantLevels    []levelInfo
		wantClustered *Clustered
	}{
		"empty-10": {
			keys:     []string{},
			maxLevel: 10,
			slimStr:  trim(""),
			wantLevels: []levelInfo{
				{0, 0, 0, nil},
			},
			wantClustered: nil,
		},
		"empty-0": {
			keys:     []string{},
			maxLevel: 0,
			slimStr:  trim(""),
			wantLevels: []levelInfo{
				{0, 0, 0, nil},
			},
			wantClustered: nil,
		},
		"singleKey-10": {
			keys:     []string{"foo"},
			maxLevel: 10,
			slimStr:  trim("#000=0"),
			wantLevels: []levelInfo{
				{0, 0, 0, nil},
				{1, 0, 1, nil},
			},
			wantClustered: nil,
		},
		"singleKey-1": {
			keys:     []string{"foo"},
			maxLevel: 1,
			slimStr:  trim("#000=0"),
			wantLevels: []levelInfo{
				{0, 0, 0, nil},
				{1, 0, 1, nil},
			},
			wantClustered: nil,
		},
		"simple-1": {
			keys:     clusteredKeysSimple,
			maxLevel: 1,
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
			wantLevels: []levelInfo{
				{0, 0, 0, nil},
				{1, 1, 0, nil},
				{4, 3, 1, nil},
				{8, 6, 2, nil},
				{14, 6, 8, nil},
			},
			wantClustered: nil,
		},
		"simple-2": {
			keys:     clusteredKeysSimple,
			maxLevel: 2,
			slimStr: trim(`
#000+4*8
    -abc->#001=0
    -abcd->#002=1
    -abd->#003=2
    -abde->#004=3
    -bc->#005=4
    -bcd->#006=5
    -bcde->#007=6
    -cde->#008=7
`),
			wantLevels: []levelInfo{
				{0, 0, 0, nil},
				{1, 1, 0, nil},
				{9, 1, 8, nil},
			},
			wantClustered: &Clustered{
				Starts:  []uint32{0, 8},
				Offsets: []uint32{0, 3, 7, 10, 14, 16, 19, 23, 26},
				Bytes:   []byte("abcabcdabdabdebcbcdbcdecde"),
			},
		},
		"simple-3": {
			keys:     clusteredKeysSimple,
			maxLevel: 3,
			slimStr: trim(`
#000+4*3
    -0001->#001+12*4
               -c->#004=0
               -cd->#005=1
               -d->#006=2
               -de->#007=3
    -0010->#002+8*3
               -->#008=4
               -d->#009=5
               -de->#010=6
    -0011->#003=7
`),
			wantLevels: []levelInfo{
				{0, 0, 0, nil},
				{1, 1, 0, nil},
				{4, 3, 1, nil},
				{11, 3, 8, nil},
			},
			wantClustered: &Clustered{
				Starts:  []uint32{0, 4, 7},
				Offsets: []uint32{0, 1, 3, 4, 6, 6, 7, 9},
				Bytes:   []byte("ccdddedde"),
			},
		},
		"simple-4": {
			keys:     clusteredKeysSimple,
			maxLevel: 4,
			slimStr: trim(`
#000+4*3
    -0001->#001+12*2
               -0011->#004*2
                          -->#008=0
                          -d->#009=1
               -0100->#005*2
                          -->#010=2
                          -e->#011=3
    -0010->#002+8*2
               -->#006=4
               -0110->#007+8*2
                          -->#012=5
                          -e->#013=6
    -0011->#003=7
`),
			wantLevels: []levelInfo{
				{0, 0, 0, nil},
				{1, 1, 0, nil},
				{4, 3, 1, nil},
				{8, 6, 2, nil},
				{14, 6, 8, nil},
			},
			wantClustered: &Clustered{
				Starts:  []uint32{0, 2, 4, 6},
				Offsets: []uint32{0, 0, 1, 1, 2, 2, 3},
				Bytes:   []byte("dee"),
			},
		},
		"simple-5": {
			keys:     clusteredKeysSimple,
			maxLevel: 5,
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
			wantLevels: []levelInfo{
				{0, 0, 0, nil},
				{1, 1, 0, nil},
				{4, 3, 1, nil},
				{8, 6, 2, nil},
				{14, 6, 8, nil},
			},
			wantClustered: nil,
		},
		"simple-6": {
			keys:     clusteredKeysSimple,
			maxLevel: 6,
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
			wantLevels: []levelInfo{
				{0, 0, 0, nil},
				{1, 1, 0, nil},
				{4, 3, 1, nil},
				{8, 6, 2, nil},
				{14, 6, 8, nil},
			},
			wantClustered: nil,
		},
		"emptyKey-2": {
			keys: []string{
				"",
				"a",
				"abc",
				"abd",
				"bc",
				"bcd",
				"cde",
			},
			maxLevel: 2,
			slimStr: trim(`
#000*7
    -->#001=0
    -a->#002=1
    -abc->#003=2
    -abd->#004=3
    -bc->#005=4
    -bcd->#006=5
    -cde->#007=6
`),
			wantLevels: []levelInfo{
				{0, 0, 0, nil},
				{1, 1, 0, nil},
				{8, 1, 7, nil},
			},
			wantClustered: &Clustered{
				Starts:  []uint32{0, 7},
				Offsets: []uint32{0, 0, 1, 4, 7, 9, 12, 15},
				Bytes:   []byte("aabcabdbcbcdcde"),
			},
		},
		"emptyKey-5": {
			keys: []string{
				"",
				"a",
				"abc",
				"abd",
				"bc",
				"bcd",
				"cde",
			},
			maxLevel: 5,
			slimStr: trim(`
#000*2
    -->#001=0
    -0110->#002*3
               -0001->#003*2
                          -->#006=1
                          -0110->#007+12*2
                                     -c->#010=2
                                     -d->#011=3
               -0010->#004+8*2
                          -->#008=4
                          -0110->#009=5
               -0011->#005=6
`),
			wantLevels: []levelInfo{
				{0, 0, 0, nil},
				{1, 1, 0, nil},
				{3, 2, 1, nil},
				{6, 4, 2, nil},
				{10, 5, 5, nil},
				{12, 5, 7, nil},
			},
			wantClustered: &Clustered{
				Starts:  []uint32{0, 2},
				Offsets: []uint32{0, 1, 2},
				Bytes:   []byte("cd"),
			},
		},
		"emptyKey-6": {
			keys: []string{
				"",
				"a",
				"abc",
				"abd",
				"bc",
				"bcd",
				"cde",
			},
			maxLevel: 6,
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
			wantLevels: []levelInfo{
				{0, 0, 0, nil},
				{1, 1, 0, nil},
				{3, 2, 1, nil},
				{6, 4, 2, nil},
				{10, 5, 5, nil},
				{12, 5, 7, nil},
			},
			wantClustered: nil,
		},
	}
)

func TestSlimTrie_clustered_small(t *testing.T) {

	for name, c := range clusteredCases {
		t.Run(name, func(t *testing.T) {
			ta := require.New(t)

			values := makeI32s(len(c.keys))
			st, err := NewSlimTrie(encode.I32{}, c.keys, values,
				Opt{Complete: Bool(true),
					MaxLevel: I32(c.maxLevel)})

			ta.NoError(err)

			dd(st)

			ta.Equal(c.wantLevels, st.levels)
			ta.Equal(c.wantClustered, st.inner.Clustered)

			ta.Equal(c.slimStr, st.String())

			testPresentKeysGRS(t, st, c.keys, values)
			testAbsentKeysGRS(t, st, c.keys)
		})
	}
}

func TestSlimTrie_clustered_big(t *testing.T) {

	testBigKeySet(t, func(t *testing.T, typ string, keys []string) {

		ta := require.New(t)

		values := makeI32s(len(keys))
		st, err := NewSlimTrie(encode.I32{}, keys, values,
			Opt{Complete: Bool(true),
				MaxLevel: I32(3)})

		ta.NoError(err)

		testPresentKeysGRS(t, st, keys, values)
		testAbsentKeysGRS(t, st, keys)

	})
}
func TestSlimTrie_initClusteredInner(t *testing.T) {

	ta := require.New(t)

	keys := clusteredKeysSimple
	values := makeI32s(len(keys))

	t.Run("simple-2", func(t *testing.T) {
		level := int32(2)
		st, err := NewSlimTrie(encode.I32{}, keys, values,
			Opt{Complete: Bool(true),
				MaxLevel: I32(level)})

		ta.NoError(err)

		cl := clusteredInner{}
		st.initClusteredInner(0, &cl)

		ta.Equal(clusteredInner{
			FirstLeafId: 1,
			Offsets:     []uint32{0, 3, 7, 10, 14, 16, 19, 23, 26},
			Bytes:       []byte("abcabcdabdabdebcbcdbcdecde"),
		}, cl)

	})

	t.Run("simple-3", func(t *testing.T) {
		level := int32(3)
		st, err := NewSlimTrie(encode.I32{}, keys, values,
			Opt{Complete: Bool(true),
				MaxLevel: I32(level)})

		ta.NoError(err)

		// 0-th clustered inner

		cl := clusteredInner{}
		st.initClusteredInner(0, &cl)
		ta.Equal(clusteredInner{
			FirstLeafId: 4,
			Offsets:     []uint32{0, 1, 3, 4, 6},
			Bytes:       []byte("ccdddedde"),
		}, cl)

		// 1-th clustered inner

		st.initClusteredInner(1, &cl)

		ta.Equal(clusteredInner{
			FirstLeafId: 8,
			Offsets:     []uint32{6, 6, 7, 9},
			Bytes:       []byte("ccdddedde"),
		}, cl)

	})

}

func TestNewClusteredInner(t *testing.T) {

	ta := require.New(t)

	firstLeafId := int32(10)
	keys := []string{
		"abc",
		"de",
	}

	cl := newClusteredInner(firstLeafId, keys, 1)
	ta.Equal(int32(10), cl.FirstLeafId)
	ta.Equal([]uint32{0, 2, 3}, cl.Offsets)
	ta.Equal([]byte("bce"), cl.Bytes)

	ta.Equal(2, cl.keyCnt(), "keyCnt()")
	ta.Equal([][]byte{[]byte("bc"), []byte("e")}, cl.keys(), "keys()")

	ta.Equal(int32(-1), cl.get("b"))
	ta.Equal(int32(-1), cl.get("c"))
	ta.Equal(int32(10), cl.get("bc"))
	ta.Equal(int32(11), cl.get("e"))

	ta.Equal(int32(10), cl.firstLeafId())
	ta.Equal(int32(11), cl.lastLeafId())

	l, e, r := cl.search("bc")
	ta.Equal([]int32{-1, 10, 11}, []int32{l, e, r})

	l, e, r = cl.search("bd")
	ta.Equal([]int32{10, -1, 11}, []int32{l, e, r})

	l, e, r = cl.search("e")
	ta.Equal([]int32{10, 11, -1}, []int32{l, e, r})

	l, e, r = cl.search("ee")
	ta.Equal([]int32{11, -1, -1}, []int32{l, e, r})
}

var OutputClusteredLeavesGet int

func BenchmarkClusteredInner_get(b *testing.B) {
	for _, n := range []int{2, 4, 16, 32, 64} {

		keys := testutil.RandStrSlice(n, 5, 10)
		mask := n - 1

		b.Run(fmt.Sprintf("keyLen:5-10/keyCnt:%d", n),
			func(b *testing.B) {

				cl := newClusteredInner(0, keys, 0)

				s := int32(0)

				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					s += cl.get(keys[i&mask])
				}

				OutputClusteredLeavesGet = int(s)
			})
	}
}
