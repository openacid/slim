package trie

import (
	"fmt"
	"testing"

	"github.com/openacid/low/size"
	"github.com/openacid/slim/encode"
	"github.com/openacid/testkeys"
	"github.com/stretchr/testify/require"
)

type reachableCase struct {
	nodeId int32
	cnt    int32
}

type cacheCase struct {
	mem  int64
	want levelCachePolicy
}

var (
	levelcacheCases = map[string]struct {
		keys      []string
		slimStr   string
		levels    []levelInfo
		reachable []reachableCase
		caches    []cacheCase
	}{
		"empty": {
			keys:    []string{},
			slimStr: trim(""),
			levels:  []levelInfo{{0, 0, 0, nil}},
			caches: []cacheCase{
				{0, levelCachePolicy{0, 0, []int32{}}},
				{1024, levelCachePolicy{0, 0, []int32{}}},
			},
		},
		"singleKey": {
			keys:    []string{"foo"},
			slimStr: trim("#000=0"),
			levels: []levelInfo{
				{0, 0, 0, nil},
				{1, 0, 1, nil},
			},
			reachable: []reachableCase{
				{0, 1},
			},
			caches: []cacheCase{
				{0, levelCachePolicy{1, 0, []int32{}}},
				{1024, levelCachePolicy{1, 0, []int32{}}},
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
				{4, 3, 1, []innerCache{{1, 4}, {2, 7}}},
				{8, 6, 2, []innerCache{{4, 2}, {5, 4}, {7, 7}}},
				{14, 6, 8, nil},
			},
			reachable: []reachableCase{
				{0, 8},
				{1, 4},
				{2, 7},
				{3, 8},
				{4, 2},
				{5, 4},
				{6, 5},
				{7, 7},
				{8, 1},
				{9, 2},
				{10, 3},
				{11, 4},
				{12, 5},
				{13, 6},
			},
			caches: []cacheCase{
				{0, levelCachePolicy{32, 0, []int32{}}},
				{4, levelCachePolicy{32, 0, []int32{}}},
				{8, levelCachePolicy{32, 0, []int32{}}},
				{16, levelCachePolicy{30, 2, []int32{2}}},
				{20, levelCachePolicy{30, 2, []int32{2}}},
				{24, levelCachePolicy{30, 2, []int32{2}}},
				{28, levelCachePolicy{30, 2, []int32{2}}},
				{32, levelCachePolicy{30, 2, []int32{2}}},
				{36, levelCachePolicy{30, 2, []int32{2}}},
				{40, levelCachePolicy{29, 3, []int32{2, 3}}},
				{1024, levelCachePolicy{29, 3, []int32{2, 3}}},
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
				{3, 2, 1, []innerCache{{2, 7}}},
				{6, 4, 2, []innerCache{{3, 3}, {4, 5}}},
				{10, 5, 5, []innerCache{{7, 3}}},
				{12, 5, 7, nil},
			},
			reachable: []reachableCase{
				{0, 7},
				{1, 1},
				{2, 7},
				{3, 3},
				{4, 5},
				{5, 6},
				{6, 1},
				{7, 3},
				{8, 4},
				{9, 5},
				{10, 1},
				{11, 2},
			},
			caches: []cacheCase{
				{0, levelCachePolicy{35, 0, []int32{}}},
				{4, levelCachePolicy{35, 0, []int32{}}},
				{8, levelCachePolicy{30, 5, []int32{4}}},
				{12, levelCachePolicy{30, 5, []int32{4}}},
				{16, levelCachePolicy{28, 7, []int32{2, 4}}},
				{20, levelCachePolicy{28, 7, []int32{2, 4}}},
				{24, levelCachePolicy{28, 7, []int32{2, 4}}},
				{28, levelCachePolicy{28, 7, []int32{2, 4}}},
				{32, levelCachePolicy{27, 8, []int32{2, 3, 4}}},
				{36, levelCachePolicy{27, 8, []int32{2, 3, 4}}},
				{40, levelCachePolicy{27, 8, []int32{2, 3, 4}}},
				{1024, levelCachePolicy{27, 8, []int32{2, 3, 4}}},
			},
		},
	}
)

func TestSlimTrie_countReacheableLeaves(t *testing.T) {

	for name, c := range levelcacheCases {
		t.Run(name, func(t *testing.T) {
			ta := require.New(t)

			values := makeI32s(len(c.keys))
			st, err := NewSlimTrie(encode.I32{}, c.keys, values, Opt{Complete: Bool(true)})
			ta.NoError(err)

			dd(st)
			ta.Equal(c.slimStr, st.String())

			for _, rc := range c.reachable {
				got := st.countReachableLeaves(rc.nodeId)
				ta.Equal(rc.cnt, got, "node id:%d", rc.nodeId)
			}
		})
	}
}
func TestSlimTrie_findLevelCaches(t *testing.T) {

	for name, c := range levelcacheCases {
		t.Run(name, func(t *testing.T) {
			ta := require.New(t)

			values := makeI32s(len(c.keys))

			st, err := NewSlimTrie(encode.I32{}, c.keys, values, Opt{Complete: Bool(true)})
			ta.NoError(err)

			dd(st)
			ta.Equal(c.slimStr, st.String())

			for _, cache := range c.caches {
				t.Run(fmt.Sprintf("%d", cache.mem), func(t *testing.T) {

					ta := require.New(t)

					got := st.findLevelCaches(cache.mem)
					ta.Equal(cache.want, got)
				})
			}
		})
	}
}

func TestSlimTrie_findLevelCaches_big(t *testing.T) {

	iambig(t)

	cases := []struct {
		name string
		mem  int64
		want levelCachePolicy
	}{
		{"200kweb2", 0, levelCachePolicy{235886 * 28, 0, []int32{}}},
		{"200kweb2", 64, levelCachePolicy{235886*28 - 707639, 707639, []int32{25, 26, 27}}},
		{"200kweb2", 2 * 1024, levelCachePolicy{235886*28 - 1650282, 1650282, []int32{21, 22, 24, 25, 26, 27}}},
		// ~ 2%
		{"200kweb2", 40 * 1024, levelCachePolicy{235886*28 - 2570532, 2570532, []int32{3, 17, 19, 20, 21, 22, 23, 24, 25, 26, 27}}},
		// ~ 10%
		{"200kweb2", 200 * 1024, levelCachePolicy{235886*28 - 3243052, 3243052, []int32{3, 13, 16, 18, 21, 23, 24, 25, 26, 27}}},

		{"20kl10", 0, levelCachePolicy{20 * 1000 * 8, 0, []int32{}}},
		{"20kl10", 2 * 1024, levelCachePolicy{20*1000*8 - 39567, 39567, []int32{6, 7}}},
		{"20kl10", 20 * 1024, levelCachePolicy{20*1000*8 - 39567, 39567, []int32{6, 7}}},

		{"20kvl10", 0, levelCachePolicy{20 * 1000 * 7, 0, []int32{}}},
		{"20kvl10", 2 * 1024, levelCachePolicy{20*1000*7 - 20002, 20002, []int32{2, 6}}},
		{"20kvl10", 20 * 1024, levelCachePolicy{20*1000*7 - 39415, 39415, []int32{2, 5, 6}}},

		{"300vl50", 0, levelCachePolicy{327 * 8, 0, []int32{}}},
		{"300vl50", 2 * 1024, levelCachePolicy{327*8 - 815, 815, []int32{2, 3, 4, 5, 6, 7}}},
		{"300vl50", 20 * 1024, levelCachePolicy{327*8 - 815, 815, []int32{2, 3, 4, 5, 6, 7}}},

		{"50kl10", 0, levelCachePolicy{50 * 1000 * 8, 0, []int32{}}},
		{"50kl10", 2 * 1024, levelCachePolicy{50*1000*8 - 49992, 49992, []int32{7}}},
		{"50kl10", 20 * 1024, levelCachePolicy{50*1000*8 - 98841, 98841, []int32{6, 7}}},

		{"50kvl10", 0, levelCachePolicy{50 * 1000 * 7, 0, []int32{}}},
		{"50kvl10", 2 * 1024, levelCachePolicy{50*1000*7 - 49996, 49996, []int32{2, 6}}},
		{"50kvl10", 20 * 1024, levelCachePolicy{50*1000*7 - 98449, 98449, []int32{2, 5, 6}}},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%s/mem:%d", c.name, c.mem), func(t *testing.T) {

			ta := require.New(t)

			keys := getKeys(c.name)
			values := makeI32s(len(keys))

			st, err := NewSlimTrie(encode.I32{}, keys, values, Opt{InnerPrefix: Bool(true)})
			ta.NoError(err)

			got := st.findLevelCaches(c.mem)
			ta.Equal(c.want, got)
		})
	}
}

func TestSlimTrie_cursorWithLevelCache(t *testing.T) {

	for name, c := range levelcacheCases {
		t.Run(name, func(t *testing.T) {
			ta := require.New(t)

			values := makeI32s(len(c.keys))
			st, err := NewSlimTrie(encode.I32{}, c.keys, values, Opt{Complete: Bool(true)})
			ta.NoError(err)

			subCursorWithCache(t, st)
		})
	}
}

func TestSlimTrie_cursorWithLevelCache_big(t *testing.T) {

	testBigKeySet(t, func(t *testing.T, typ string, keys []string) {
		ta := require.New(t)

		values := makeI32s(len(keys))
		st, err := NewSlimTrie(encode.I32{}, keys, values, Opt{Complete: Bool(true)})
		ta.NoError(err)

		subCursorWithCache(t, st)

	})
}

func subCursorWithCache(t *testing.T, st *SlimTrie) {
	ta := require.New(t)
	bottom := len(st.levels) - 1
	totalNode := st.levels[bottom].total
	lvl := int32(0)
	for nid := int32(0); nid < totalNode; nid++ {
		if nid >= st.levels[lvl].total {
			lvl++
		}
		cur := &walkingCursor{
			id:         nid,
			smallerCnt: 0,
			lvl:        lvl,
		}
		want := st.cursorLeafIndex(cur, false)

		cur = &walkingCursor{
			id:         nid,
			smallerCnt: 0,
			lvl:        lvl,
		}
		got := st.cursorLeafIndex(cur, true)
		ta.Equal(want, got, "from nodeId: %d, level: %d", nid, lvl)
	}

}

var OutputLevelCache int64

func BenchmarkSlimTrie_findLevelCaches(b *testing.B) {

	for _, typ := range testkeys.AssetNames() {

		if typ == "1mvl5_10" {
			continue
		}

		for _, percentage := range []int64{1, 5} {

			b.Run(fmt.Sprintf("%s/%d%%", typ, percentage), func(b *testing.B) {
				keys := getKeys(typ)

				values := makeI32s(len(keys))

				st, err := NewSlimTrie(encode.I32{}, keys, values, Opt{InnerPrefix: Bool(true)})
				if err != nil {
					panic("˙∆˙...")
				}

				sz := size.Of(st)
				mem := int64(sz) * percentage / 100

				b.ResetTimer()

				var s int64
				for i := 0; i < b.N; i++ {
					got := st.findLevelCaches(mem)
					s += got.reduced
				}

				OutputLevelCache = s
			})
		}
	}
}
