package trie

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/kr/pretty"
	"github.com/openacid/errors"
	"github.com/openacid/slim/encode"
	"github.com/openacid/testkeys"
	"github.com/stretchr/testify/require"
)

type searchRst struct {
	lVal  interface{}
	eqVal interface{}
	rVal  interface{}
}

type searchCase struct {
	key  string
	want searchRst
}

type slimCase struct {
	keys     []string
	values   []int
	searches []searchCase
}

// from8bit create string from 8bit words
func from8bit(x ...byte) string {
	return string(x)
}

func TestNewSlimTrie(t *testing.T) {

	ta := require.New(t)

	st, err := NewSlimTrie(encode.Int{}, []string{"ab", "cd"}, []int{1, 2})
	if err != nil {
		t.Fatalf("expect no error but: %v", err)
	}

	v, found := st.Get("ab")
	ta.True(found)
	ta.Equal(1, v)

	v, found = st.Get("cd")
	ta.True(found)
	ta.Equal(2, v)
}

func TestNewSlimTrie_Error(t *testing.T) {

	ta := require.New(t)

	cases := []struct {
		keys    []string
		values  []int
		wanterr error
	}{
		{
			[]string{"a", "a"},
			[]int{1, 2},
			ErrKeyOutOfOrder,
		},
		{
			[]string{"ab", "a"},
			[]int{1, 2},
			ErrKeyOutOfOrder,
		},
		{
			[]string{"ab", "aa"},
			[]int{1, 2},
			ErrKeyOutOfOrder,
		},
		{
			[]string{"ab", "aaa"},
			[]int{1, 2},
			ErrKeyOutOfOrder,
		},
	}

	for i, c := range cases {
		st, err := NewSlimTrie(encode.Int{}, c.keys, c.values)
		ta.Equal(c.wanterr, errors.Cause(err), "%d-th: input: keys: %v; vals: %v; wanterr: %v; actual: %v",
			i+1, c.keys, c.values, c.wanterr, err)

		if err == nil && len(c.keys) > 0 {
			v, found := st.Get(c.keys[0])
			if !found {
				t.Fatalf("%d-th: should be found but not. key=%q",
					i+1, c.keys[0])
			}

			if v == nil {
				t.Fatalf("%d-th: should be found but not. key=%q",
					i+1, c.keys[0])
			}
		}
	}
}

func TestNewSlimTrie_empty(t *testing.T) {

	ta := require.New(t)

	ks := []string{}
	vs := []int32{}

	st, err := NewSlimTrie(encode.I32{}, ks, vs)
	ta.NoError(err)
	testUnknownKeysGRS(t, st, randVStrings(100, 0, 10))

	// marshal
	buf1, err := proto.Marshal(st)
	ta.NoError(err)

	st2, _ := NewSlimTrie(encode.I32{}, nil, nil)
	err = proto.Unmarshal(buf1, st2)
	ta.NoError(err)
	slimtrieEqual(st, st2, t)
	testUnknownKeysGRS(t, st2, randVStrings(100, 0, 10))
}

func TestSlimTrie_GRS_1_zerokeys(t *testing.T) {

	ta := require.New(t)

	keys := []string{}
	values := makeI32s(len(keys))

	st, err := NewSlimTrie(encode.I32{}, keys, values)
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

func TestSlimTrie_GRS_1_zeroLengthValues(t *testing.T) {

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
	// values := make([]interface{}, len(keys))

	// TODO values should be optional for dummy mode.
	// TODO as a filter, the check a range.
	st, err := NewSlimTrie(encode.Dummy{}, keys, nil)
	ta.NoError(err)

	wantstr := trim(`
#000+4*3
    -0001->#001+12*2
               -0011->#004*2
                          -->#008=<nil>
                          -0110->#009=<nil>
               -0100->#005*2
                          -->#010=<nil>
                          -0110->#011=<nil>
    -0010->#002+8*2
               -->#006=<nil>
               -0110->#007+4*2
                          -->#012=<nil>
                          -0110->#013=<nil>
    -0011->#003=<nil>
`)
	dd(st)

	ta.Equal(wantstr, st.String())

	testUnknownKeysGRS(t, st, randVStrings(len(keys)*5, 0, 10))

	for _, key := range keys {

		dd("test Get: present: %s", key)

		v, found := st.Get(key)
		ta.True(found, "get %v", key)
		ta.Nil(v, "search for %v", key)
	}
	for _, key := range keys {

		dd("test RangeGet: present:%s", key)

		v, found := st.RangeGet(key)
		ta.True(found, "RangeGet:%v", key)
		ta.Nil(v, "RangeGet:%v", key)
	}
	for _, key := range keys {

		dd("test Search: present: %s", key)

		l, e, r := st.Search(key)
		ta.Nil(l, "Search:%v, l", key)
		ta.Nil(e, "Search:%v, e", key)
		ta.Nil(r, "Search:%v, r", key)
	}
}

func TestSlimTrie_GRS_1_nilValuesIgnoreEncoder(t *testing.T) {

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

	// use I32, but actually ignored.
	st, err := NewSlimTrie(encode.I32{}, keys, nil)
	ta.NoError(err)

	wantstr := trim(`
#000+4*3
    -0001->#001+12*2
               -0011->#004*2
                          -->#008=<nil>
                          -0110->#009=<nil>
               -0100->#005*2
                          -->#010=<nil>
                          -0110->#011=<nil>
    -0010->#002+8*2
               -->#006=<nil>
               -0110->#007+4*2
                          -->#012=<nil>
                          -0110->#013=<nil>
    -0011->#003=<nil>
`)
	dd(st)

	ta.Equal(wantstr, st.String())

	testUnknownKeysGRS(t, st, randVStrings(len(keys)*5, 0, 10))

	for _, key := range keys {

		dd("test Get: present: %s", key)

		v, found := st.Get(key)
		ta.True(found, "get %v", key)
		ta.Nil(v, "search for %v", key)
	}
	for _, key := range keys {

		dd("test RangeGet: present:%s", key)

		v, found := st.RangeGet(key)
		ta.True(found, "RangeGet:%v", key)
		ta.Nil(v, "RangeGet:%v", key)
	}
	for _, key := range keys {

		dd("test Search: present: %s", key)

		l, e, r := st.Search(key)
		ta.Nil(l, "Search:%v, l", key)
		ta.Nil(e, "Search:%v, e", key)
		ta.Nil(r, "Search:%v, r", key)
	}
}

func TestSlimTrie_GRS_1_onekey(t *testing.T) {

	ta := require.New(t)

	keys := []string{
		"abc",
	}
	values := makeI32s(len(keys))

	st, err := NewSlimTrie(encode.I32{}, keys, values)
	ta.NoError(err)

	wantstr := trim(`
#000=0
`)
	ta.Equal(wantstr, st.String())

	testUnknownKeysGRS(t, st, randVStrings(len(keys)*5, 0, 10))
	testPresentKeysGRS(t, st, keys, values)
}

func TestSlimTrie_GRS_1_twokeys_emptysteps(t *testing.T) {

	ta := require.New(t)

	// the first bit diffs, thus no step needed
	keys := []string{
		"abc",
		"\x80bc",
	}
	values := makeI32s(len(keys))

	st, err := NewSlimTrie(encode.I32{}, keys, values)
	ta.NoError(err)

	wantstr := trim(`
#000*2
    -0110->#001=0
    -1000->#002=1
`)
	dd(st)

	ta.Equal(wantstr, st.String())
	testUnknownKeysGRS(t, st, randVStrings(len(keys)*5, 0, 10))
	testPresentKeysGRS(t, st, keys, values)
}

func TestSlimTrie_RangeGet_0_tiny(t *testing.T) {

	ta := require.New(t)

	keys := []string{
		"abc",
		"abcd",

		"abd",
		"abde",

		"bc",

		"bcd",
		"bce",

		"cde",
	}
	values := []int32{
		0, 0,
		1, 1,
		2,
		3, 3,
		4,
	}
	searches := []struct {
		key       string
		want      interface{}
		wantfound bool
	}{
		{"ab", nil, false},
		{"abc", int32(0), true},
		{"abcde", int32(0), true},
		{"abd", int32(1), true},
		{"ac", nil, false},
		{"acd", int32(1), true},
		{"adc", int32(0), true},
		{"bcd", int32(3), true},
		{"bce", int32(3), true},
		{"c", int32(4), true},
		{"cde", int32(4), true},
		{"cfe", int32(4), true},
		{"cff", int32(4), true},
	}

	st, err := NewSlimTrie(encode.I32{}, keys, values)
	ta.NoError(err)

	wantstr := trim(`
#000+4*3
    -0001->#001+12*2
               -0011->#004
                          -->#008=0
               -0100->#005
                          -->#009=1
    -0010->#002+8*2
               -->#006=2
               -0110->#007
                          -0100->#010=3
    -0011->#003=4
`)
	dd(keys)
	dd(st)
	ta.Equal(wantstr, st.String())

	for _, c := range searches {

		dd("RangeGet: %s", c.key)

		rst, found := st.RangeGet(c.key)

		ta.Equal(c.want, rst, "search for %v", c.key)
		ta.Equal(c.wantfound, found, "get %v", c.key)
	}
}

func TestSlimTrie_RangeGet_1_leafNotToKeep(t *testing.T) {

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

	st, err := NewSlimTrie(encode.I32{}, keys, values)
	ta.NoError(err)

	dd(keys)
	dd(st)

	wantstr := trim(`
#000+12*2
    -0001->#001=0
    -1100->#002
               -0110->#003
                          -0101->#004=1
`)

	ta.Equal(wantstr, st.String())

	for i, c := range keys {

		rst, found := st.RangeGet(c)

		ta.Equal(values[i], rst, "%d-th: search: %+v", i+1, c)
		ta.Equal(true, found, "%d-th: search: %+v", i+1, c)
	}
}

func TestSlimTrie_RangeGet_1_rangeindex_bug_2019_05_21(t *testing.T) {

	// RangeGet has bug found by Liu Baohai:

	ta := require.New(t)

	keys := []string{
		"test/存界needleid00011end",

		"test/山我needleid00009end",
		"test/界世needleid00005end",
		"test/白我needleid00006end",

		"test/白测needleid00008end",
		"test/试世needleid00014end",
	}
	values := []int32{
		0,
		1, 1, 1,
		2, 2,
	}

	st, err := NewSlimTrie(encode.I32{}, keys, values)
	ta.NoError(err)

	dd(keys)
	dd(st)

	wantstr := trim(`
#000+44*2
    -0101->#001*2
               -1010->#003=0
               -1011->#004=1
    -0111->#002+4
               -1001->#005+16
                          -1011->#006=2
`)

	ta.Equal(wantstr, st.String())

	testPresentKeysRangeGet(t, st, keys, values)
}

func TestSlimTrie_GRS_1_u16step_bug_2019_05_29(t *testing.T) {

	// Reported by @aaaton
	// 2019 May 29
	//
	// When number of keys becomes greater than 50000,
	// SlimTrie.Get() returns negaitve for some existent keys.
	// Caused by SlimTrie.step has been using uint16 id, it should be int32.

	iambig(t)

	ta := require.New(t)

	keys := getKeys("50kl10")
	n := len(keys)
	values := makeI32s(n)

	st, err := NewSlimTrie(encode.I32{}, keys, values)
	ta.NoError(err)

	testUnknownKeysGRS(t, st, randVStrings(len(keys)*5, 0, 10))
	testPresentKeysGRS(t, st, keys, values)
}

func TestSlimTrie_GRS_0_tiny(t *testing.T) {

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

	st, err := NewSlimTrie(encode.I32{}, keys, values)
	ta.NoError(err)

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
               -0110->#007+4*2
                          -->#012=5
                          -0110->#013=6
    -0011->#003=7
`)
	dd(st.content())
	dd(st)

	ta.Equal(wantstr, st.String())

	testUnknownKeysGRS(t, st, randVStrings(len(keys)*5, 0, 10))
	testPresentKeysGRS(t, st, keys, values)
}

func TestSlimTrie_GRS_3_bigInner_a2t(t *testing.T) {

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

	st, err := NewSlimTrie(encode.I32{}, keys, values)
	ta.NoError(err)

	dd(st)

	ta.True(st.inner.BigInnerCnt > 0)

	testUnknownKeysGRS(t, st, randVStrings(len(keys)*5, 0, 20))
	testPresentKeysGRS(t, st, keys, values)
}

func TestSlimTrie_GRS_3_bigInner_300(t *testing.T) {

	ta := require.New(t)
	keys := getKeys("300vl50")
	values := makeI32s(len(keys))

	st, err := NewSlimTrie(encode.I32{}, keys, values)
	ta.NoError(err)

	dd(st)

	ta.True(st.inner.BigInnerCnt > 0)

	testUnknownKeysGRS(t, st, randVStrings(len(keys)*5, 0, 10))
	testPresentKeysGRS(t, st, keys, values)
}

func TestSlimTrie_GRS_9_allkeyset(t *testing.T) {

	for _, typ := range testkeys.AssetNames() {

		t.Run(fmt.Sprintf("keyset: %s", typ), func(t *testing.T) {
			ta := require.New(t)
			keys := getKeys(typ)
			if len(keys) >= 1000 {
				iambig(t)
			}

			values := makeI32s(len(keys))
			st, err := NewSlimTrie(encode.I32{}, keys, values)
			ta.NoError(err)

			testUnknownKeysGRS(t, st, randVStrings(len(keys)*5, 0, 20))
			testPresentKeysGRS(t, st, keys, values)
		})
	}
}

func TestSlimTrie_GRS_1_empty_string_branch(t *testing.T) {

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

	st, err := NewSlimTrie(encode.I32{}, keys, values)
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

	testUnknownKeysGRS(t, st, randVStrings(len(keys)*5, 0, 10))
	testPresentKeysGRS(t, st, keys, values)
}

func TestSlimTrie_Search_0_tiny(t *testing.T) {

	ta := require.New(t)

	cases := []slimCase{
		{
			keys: []string{
				from8bit(1, 2, 3),
				from8bit(1, 2, 4),
				from8bit(2, 3, 4),
				from8bit(2, 3, 5),
				from8bit(3, 4, 5),
			},
			values: []int{0, 1, 2, 3, 4},
			searches: []searchCase{
				{from8bit(1, 2, 3), searchRst{nil, 0, 1}},
				{from8bit(1, 2, 4), searchRst{0, 1, 2}},
				{from8bit(2, 3, 4), searchRst{1, 2, 3}},
				{from8bit(2, 3, 5), searchRst{2, 3, 4}},
				{from8bit(3, 4, 5), searchRst{3, 4, nil}},
			},
		},
		{
			keys: []string{
				from8bit(1, 2, 3),
				from8bit(1, 2, 3, 4),
				from8bit(2, 3),
				from8bit(2, 3, 0),
				from8bit(2, 3, 4),
				from8bit(2, 3, 4, 5),
				from8bit(2, 3, 15),
			},
			values: []int{0, 1, 2, 3, 4, 5, 6},
			searches: []searchCase{
				{from8bit(1, 2, 3), searchRst{nil, 0, 1}},
				{from8bit(1, 2, 3, 4), searchRst{0, 1, 2}},
				{from8bit(2, 3), searchRst{1, 2, 3}},
				{from8bit(2, 3, 0), searchRst{2, 3, 4}},
				{from8bit(2, 3, 4), searchRst{3, 4, 5}},
				{from8bit(2, 3, 4, 5), searchRst{4, 5, 6}},
				{from8bit(2, 3, 15), searchRst{5, 6, nil}},
			},
		},
		{
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
			values: []int{0, 1, 2, 3, 4, 5, 6, 7},
			searches: []searchCase{
				{"ab", searchRst{nil, nil, 0}},
				{"abc", searchRst{nil, 0, 1}},
				{"abcde", searchRst{0, 1, 2}},
				{"abd", searchRst{1, 2, 3}},
				{"ac", searchRst{nil, nil, 0}},
				{"acd", searchRst{1, 2, 3}},
				{"adc", searchRst{nil, 0, 1}},
				{"bcd", searchRst{4, 5, 6}},
				{"bce", searchRst{4, 5, 6}},
				{"c", searchRst{6, 7, nil}},
				{"cde", searchRst{6, 7, nil}},
				{"cfe", searchRst{6, 7, nil}},
				{"cff", searchRst{6, 7, nil}},
			},
		},
	}
	for i, c := range cases {

		st, err := NewSlimTrie(encode.Int{}, c.keys, c.values)
		ta.NoError(err)

		dd("keys: %#v", c.keys)
		dd(st)

		for _, ex := range c.searches {
			lt, eq, gt := st.Search(ex.key)
			rst := searchRst{lt, eq, gt}

			dd("search for %#v", ex.key)

			ta.Equal(ex.want, rst, "%d-th case: %+v search: %+v", i+1, c, ex)
		}
	}
}

func slimtrieEqual(st1, st2 *SlimTrie, t *testing.T) {
	if !proto.Equal((st1.inner), (st2.inner)) {
		fmt.Println(st1)
		fmt.Println(st2)
		fmt.Println(pretty.Diff(st1.inner, st2.inner))
		t.Fatalf("Children not the same")
	}
}

func makeI32s(n int) []int32 {
	values := make([]int32, n)
	for i := int32(0); i < int32(n); i++ {
		values[i] = i
	}
	return values

}
func trim(s string) string {
	return strings.Trim(s, "\n")
}

func testPresentKeysGRS(t *testing.T, st *SlimTrie, keys []string, values []int32) {

	testPresentKeysGet(t, st, keys, values)
	testPresentKeysRangeGet(t, st, keys, values)
	testPresentKeysSearch(t, st, keys, values)
}

func testAbsentKeysGRS(t *testing.T, st *SlimTrie, keys []string) {

	absentKeys := makeAbsentKeys(keys, len(keys)*5, 0, 20)

	testAbsentKeysGet(t, st, absentKeys)
	testAbsentKeysRangeGet(t, st, keys, absentKeys)
	testAbsentKeysSearch(t, st, keys, absentKeys)
}

func testUnknownKeysGRS(t *testing.T, st *SlimTrie, keys []string) {

	for i, key := range keys {
		_ = i

		dd("test unknown: %s", key)

		// There may be false positive result

		id := st.GetID(key)
		_ = id

		v, found := st.Get(key)
		_ = v
		_ = found

		v, found = st.RangeGet(key)
		_ = v
		_ = found

		lid, eid, rid := st.searchID(key)
		_ = lid
		_ = eid
		_ = rid

		l, e, r := st.Search(key)
		_ = l
		_ = e
		_ = r

	}
}

func makeAbsentKeys(keys []string, n, minLen, maxLen int) []string {

	rst := make([]string, 0, n)

	mp := make(map[string]bool)
	for _, k := range keys {
		mp[k] = true
	}

	for len(rst) < n {

		ks := randVStrings(100, minLen, maxLen)
		for _, k := range ks {
			if mp[k] {
				continue
			}
			mp[k] = true
			rst = append(rst, k)
		}
	}

	sort.Strings(rst)
	return rst
}

func testAbsentKeysGet(t *testing.T, st *SlimTrie, absentKeys []string) {

	ta := require.New(t)

	for i, key := range absentKeys {
		_ = i

		dd("test Get: absent: %s", key)

		v, found := st.Get(key)
		ta.Nil(v, "get absent %s", key)
		ta.False(found, "get absent %s", key)
	}
}

func testAbsentKeysRangeGet(t *testing.T, st *SlimTrie, keys, absentKeys []string) {

	ta := require.New(t)

	j := int32(0)
	prev := int32(-1)

	for i, key := range absentKeys {
		_ = i

		for ; j < int32(len(keys)) && keys[j] < key; j++ {
			prev = j
		}

		dd("test RangeGet: absent: %s", key)

		v, found := st.RangeGet(key)
		dd(v, found)

		if prev == -1 {
			ta.Nil(v, "range get %v", key)
			ta.False(found, "range get %v", key)
		} else {
			ta.Equal(prev, v, "range get %v", key)
			ta.True(found, "range get %v", key)
		}
	}
}

func testAbsentKeysSearch(t *testing.T, st *SlimTrie, keys, absentKeys []string) {

	ta := require.New(t)

	j := int32(0)
	prev := int32(-1)
	for i, key := range absentKeys {
		_ = i
		for ; j < int32(len(keys)) && keys[j] < key; j++ {
			prev = j
		}

		dd("test Search: absent: %s", key)

		l, e, r := st.Search(key)

		ta.Equal(nil, e, "Search:%v e", key)

		if prev == -1 {

			ta.Equal(nil, l, "Search:%v l", key)
			ta.Equal(int32(0), r, "Search:%v r", key)

		} else if prev == int32(len(keys))-1 {

			ta.Equal(prev, l, "Search:%v l", key)
			ta.Equal(nil, r, "Search:%v r", key)

		} else {
			ta.Equal(prev, l, "Search:%v l", key)
			ta.Equal(prev+1, r, "Search:%v r", key)
		}

	}
}

func testPresentKeysGet(t *testing.T, st *SlimTrie, keys []string, values []int32) {

	ta := require.New(t)

	for i, key := range keys {

		dd("test Get: present: %s", key)

		v, found := st.Get(key)
		ta.True(found, "get %v", key)
		ta.Equal(values[i], v, "search for %v", key)
	}
}

func testPresentKeysRangeGet(t *testing.T, st *SlimTrie, keys []string, values []int32) {

	ta := require.New(t)

	for i, key := range keys {

		dd("test RangeGet: present:%s", key)

		v, found := st.RangeGet(key)
		ta.True(found, "RangeGet:%v", key)
		ta.Equal(values[i], v, "RangeGet:%v", key)
	}
}

func testPresentKeysSearch(t *testing.T, st *SlimTrie, keys []string, values []int32) {

	ta := require.New(t)

	for i, key := range keys {

		dd("test Search: present: %s", key)

		l, e, r := st.Search(key)
		if i == 0 {
			ta.Equal(nil, l)
		} else {
			ta.Equal(values[i-1], l, "Search:%v, l", key)
		}

		ta.Equal(values[i], e, "Search:%v, eq", key)

		if i == len(keys)-1 {
			ta.Equal(nil, r)
		} else {
			ta.Equal(values[i+1], r, "Search:%v, r", key)
		}
	}
}

func testBigKeySet(t *testing.T, f func(t *testing.T, keys []string)) {

	for _, typ := range testkeys.AssetNames() {

		if typ == "1mvl5_10" {
			continue
		}

		t.Run(typ, func(t *testing.T) {
			keys := getKeys(typ)
			n := len(keys)
			if n >= 1000 {
				iambig(t)
			}

			f(t, keys)
		})
	}
}

func clap(n, min, max int) int {
	if n < min {
		n = min
	}
	if n > max {
		n = max
	}
	return n
}

func iambig(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
}

func dd(msgAndArgs ...interface{}) {
	if !testing.Verbose() {
		return
	}

	fmt.Println(fmtMsg(msgAndArgs...))
}

func fmtMsg(msgAndArgs ...interface{}) string {
	if len(msgAndArgs) == 0 || msgAndArgs == nil {
		return ""
	}
	if len(msgAndArgs) == 1 {
		msg := msgAndArgs[0]
		if msgAsStr, ok := msg.(string); ok {
			return msgAsStr
		}
		return fmt.Sprintf("%+v", msg)
	}
	if len(msgAndArgs) > 1 {
		msg := msgAndArgs[0]
		if msgAsStr, ok := msg.(string); ok {
			return fmt.Sprintf(msgAsStr, msgAndArgs[1:]...)
		}

		rst := []string{}
		for _, m := range msgAndArgs {
			rst = append(rst, fmt.Sprintf("%+v", m))
		}
		return strings.Join(rst, " ")
	}
	return ""
}

var cache map[string][]string = map[string][]string{}

func getKeys(fn string) []string {

	ss, ok := cache[fn]
	if ok {
		return ss
	}

	ks := testkeys.Load(fn)
	cache[fn] = ks
	return ks
}
