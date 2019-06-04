package trie

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/kr/pretty"
	"github.com/openacid/errors"
	"github.com/openacid/low/bitword"
	"github.com/openacid/low/pbcmpl"
	"github.com/openacid/slim/array"
	"github.com/openacid/slim/encode"
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

// from8bit create string from 4bit words
func from4bit(x ...byte) string {
	return bw4.ToStr(x)
}

var (
	// a squashed case and also as case for marshaling test
	marshalCase = slimCase{
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
			{"abcde", searchRst{0, 1, 2}}, // false positive
			{"abd", searchRst{1, 2, 3}},
			{"ac", searchRst{nil, nil, 0}},
			{"acb", searchRst{nil, nil, 0}},
			{"acd", searchRst{1, 2, 3}},
			{"adc", searchRst{nil, 0, 1}},
			{"bcd", searchRst{4, 5, 6}},
			{"bce", searchRst{4, 5, 6}},
			{"c", searchRst{6, 7, nil}}, // false positive
			{"cde", searchRst{6, 7, nil}},
			{"cfe", searchRst{6, 7, nil}},
			{"cff", searchRst{6, 7, nil}},
		},
	}
)

func TestMaxKeys(t *testing.T) {

	ta := require.New(t)

	nn := 256
	// a milllion keys
	mx := nn * nn * 16

	keys := make([]string, 0, mx)
	values := make([]int32, 0, mx)

	for i := 0; i < nn; i++ {
		for j := 0; j < nn; j++ {
			for k := 0; k < 16; k++ {
				key := string([]byte{byte(i), byte(j), byte(k << 4)})
				keys = append(keys, key)

				value := i*nn*nn + j*nn + k
				values = append(values, int32(value))
			}
		}
	}

	st, err := NewSlimTrie(encode.I32{}, keys, values)
	ta.Nil(err)

	ta.Equal(int32(1+16+256+4096+65536), st.Children.Cnt)
	ta.Equal(int32(0), st.Steps.Cnt)
	ta.Equal(int32(mx), st.Leaves.Cnt)

	for i := 0; i < nn; i++ {
		for j := 0; j < nn; j++ {
			for k := 0; k < 16; k++ {

				key := string([]byte{byte(i), byte(j), byte(k << 4)})

				v, _ := st.Get(key)
				ta.Equal(values[i*nn*16+j*16+k], v)
			}
		}
	}
}

func TestMaxNode(t *testing.T) {

	ta := require.New(t)

	mx := 32768

	keys := make([]string, 0, mx)
	values := make([]int, 0, mx)

	for i := 0; i < mx; i++ {

		key := []byte{
			byte((i >> 14) & 0x01),
			byte((i >> 13) & 0x01),
			byte((i >> 12) & 0x01),
			byte((i >> 11) & 0x01),
			byte((i >> 10) & 0x01),
			byte((i >> 9) & 0x01),
			byte((i >> 8) & 0x01),
			byte((i >> 7) & 0x01),
			byte((i >> 6) & 0x01),
			byte((i >> 5) & 0x01),
			byte((i >> 4) & 0x01),
			byte((i >> 3) & 0x01),
			byte((i >> 2) & 0x01),
			byte((i >> 1) & 0x01),
			byte(i & 0x01),
		}

		keys = append(keys, bitword.BitWord[4].ToStr(key))
		values = append(values, i)
	}

	sl, err := NewSlimTrie(encode.Int{}, keys, values)
	ta.Nil(err)

	if sl.Children.Cnt != int32(mx-1) {
		t.Fatalf("children cnt should be %d, but: %d", mx-1, sl.Children.Cnt)
	}
	if sl.Steps.Cnt != int32(0) {
		t.Fatalf("Steps cnt should be %d", mx)
	}
	if sl.Leaves.Cnt != int32(mx) {
		t.Fatalf("leaves cnt should be %d", mx)
	}
}

func TestUnsquashedSearch(t *testing.T) {

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
				{"abcde", searchRst{1, nil, 2}},
				{"abd", searchRst{1, 2, 3}},
				{"ac", searchRst{3, nil, 4}},
				{"acb", searchRst{3, nil, 4}},
				{"acd", searchRst{3, nil, 4}},
				{"adc", searchRst{3, nil, 4}},
				{"bcd", searchRst{4, 5, 6}},
				{"bce", searchRst{6, nil, 7}},
				{"c", searchRst{6, nil, 7}},
				{"cde", searchRst{6, 7, nil}},
				{"cfe", searchRst{7, nil, nil}},
				{"cff", searchRst{7, nil, nil}},
			},
		},
	}

	for _, c := range cases {

		keys := bw4.FromStrs(c.keys)

		// Unsquashed Trie

		trie, err := NewTrie(keys, c.values, false)
		if err != nil {
			t.Fatalf("expected no error but: %+v", err)
		}

		for _, ex := range c.searches {
			lt, eq, gt := trie.Search(bw4.FromStr(ex.key))
			rst := searchRst{lt, eq, gt}

			if !reflect.DeepEqual(ex.want, rst) {
				fmt.Println(trie)
				fmt.Println("search:", bw4.FromStr(ex.key))
				t.Fatal("key: ", bw4.FromStr(ex.key), "expected value: ", ex.want, "rst: ", rst)
			}
		}
	}
}

func TestSquashedTrieSearch(t *testing.T) {

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
				{"abcde", searchRst{1, nil, 2}},
				{"abd", searchRst{1, 2, 3}},
				{"ac", searchRst{nil, nil, 0}},
				{"acb", searchRst{nil, nil, 0}},
				{"acd", searchRst{1, 2, 3}},
				{"adc", searchRst{nil, 0, 1}},
				{"bcd", searchRst{4, 5, 6}},
				{"bce", searchRst{4, 5, 6}},
				{"c", searchRst{6, nil, 7}},
				{"cde", searchRst{6, 7, nil}},
				{"cfe", searchRst{6, 7, nil}},
				{"cff", searchRst{6, 7, nil}},
			},
		},
	}

	for _, c := range cases {

		keys := bw4.FromStrs(c.keys)

		// Squashed Trie

		trie, err := NewTrie(keys, c.values, true)
		if err != nil {
			t.Fatalf("expected no error but: %+v", err)
		}

		for _, ex := range c.searches {
			lt, eq, gt := trie.Search(bw4.FromStr(ex.key))
			rst := searchRst{lt, eq, gt}

			if !reflect.DeepEqual(ex.want, rst) {
				t.Fatal("key: ", ex.key, "expected value: ", ex.want, "rst: ", rst)
			}
		}

		// Squashed twice Trie

		trie.Squash()
		for _, ex := range c.searches {
			lt, eq, gt := trie.Search(bw4.FromStr(ex.key))
			rst := searchRst{lt, eq, gt}

			if !reflect.DeepEqual(ex.want, rst) {
				t.Fatal("key: ", ex.key, "expected value: ", ex.want, "rst: ", rst)
			}
		}
	}
}

func TestSlimTrieSearch(t *testing.T) {

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
		marshalCase,
	}
	for _, c := range cases {

		st, err := NewSlimTrie(encode.Int{}, c.keys, c.values)
		if err != nil {
			t.Fatalf("expected no error but: %+v", err)
		}

		for _, ex := range c.searches {
			lt, eq, gt := st.Search(ex.key)
			rst := searchRst{lt, eq, gt}

			if !reflect.DeepEqual(ex.want, rst) {
				fmt.Println(c.keys)
				fmt.Println(st)
				t.Fatal("key: ", ex.key, "expected value: ", ex.want, "rst: ", rst)
			}
		}
	}
}

func TestRangeGet_search(t *testing.T) {

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
	values := []int{
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
		{"abc", 0, true},
		{"abcde", 0, true}, // false positive
		{"abd", 1, true},
		{"ac", nil, false},
		{"acb", nil, false},
		{"acd", 1, true}, // false positive
		{"adc", 0, true}, // false positive
		{"bcd", 3, true},
		{"bce", 3, true},
		{"c", 4, true}, // false positive
		{"cde", 4, true},
		{"cfe", 4, true}, // false positive
		{"cff", 4, true}, // false positive
		{"def", 4, true}, // false positive
	}

	st, err := NewSlimTrie(encode.Int{}, keys, values)
	ta.Nil(err)

	wantstr := `
#000+1*3
    -001->#001+3*2
              -003->#004=0
              -004->#005=1
    -002->#002+2=2
              -006->#006
                        -004->#007=3
    -003->#003=4`[1:]
	ta.Equal(wantstr, st.String())

	for i, c := range searches {
		rst, found := st.RangeGet(c.key)
		if c.want != rst {
			t.Fatalf("%d-th key: %s expect: %v; but: %v", i+1, c.key, c.want, rst)
		}
		if c.wantfound != found {
			t.Fatalf("%d-th key: %s expect: %v; but: %v", i+1, c.key, c.wantfound, found)
		}
	}
}
func TestSlimTrie_RangeGet_leafNotToKeep(t *testing.T) {

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
	ta.Nil(err)

	wantstr := `
#000+3*2
    -001->#001=0
    -012->#002
              -006->#003
                        -005->#004=1`[1:]

	ta.Equal(wantstr, st.String())

	for i, c := range keys {
		rst, found := st.RangeGet(c)
		ta.Equal(values[i], rst, "%d-th: search: %+v", i+1, c)
		ta.Equal(true, found, "%d-th: search: %+v", i+1, c)
	}
}

func TestSlimTrie_RangeGet_rangeindex_bug_2019_05_21(t *testing.T) {

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
	ta.Nil(err)

	wantstr := `
#000+11*2
    -005->#001*2
              -010->#003=0
              -011->#004=1
    -007->#002+1
              -009->#005+4
                        -011->#006=2`[1:]

	ta.Equal(wantstr, st.String())

	for i, c := range keys {
		rst, found := st.RangeGet(c)
		ta.Equal(values[i], rst, "%d-th: search: %+v", i+1, c)
		ta.Equal(true, found, "%d-th: search: %+v", i+1, c)
	}
}

func TestSlimTrie_u16step_bug_2019_05_29(t *testing.T) {

	// Reported by @aaaton
	// 2019 May 29
	//
	// When number of keys becomes greater than 50000,
	// SlimTrie.Get() returns negaitve for some existent keys.
	// Caused by SlimTrie.step has been using uint16 id, it should be int32.

	ta := require.New(t)

	keys := keys50k
	n := len(keys)
	values := make([]int32, n)
	for i := 0; i < n; i++ {
		values[i] = int32(i)
	}
	st, err := NewSlimTrie(encode.I32{}, keys, values)
	ta.Nil(err)

	for i, c := range keys {
		rst, found := st.Get(c)
		ta.Equal(values[i], rst, "%d-th: Get: %+v", i+1, c)
		ta.Equal(true, found, "%d-th: Get: %+v", i+1, c)
	}

}

func TestNewSlimTrie(t *testing.T) {

	ta := require.New(t)

	st, err := NewSlimTrie(encode.Int{}, []string{"ab", "cd"}, []int{1, 2})
	if err != nil {
		t.Fatalf("expect no error but: %v", err)
	}

	v, found := st.Get("ab")
	if !found {
		t.Fatalf("%q should be found", "ab")
	}

	ta.Equal(1, v)

	if v.(int) != 1 {
		t.Fatalf("v should be 2, but: %v", v)
	}
}

func TestSlimTrieError(t *testing.T) {
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

		if c.wanterr != errors.Cause(err) {
			t.Fatalf("%d-th: input: keys: %v; vals: %v; wanterr: %v; actual: %v",
				i+1, c.keys, c.values, c.wanterr, err)
		}

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

func TestSlimTrie_MarshalUnmarshal(t *testing.T) {

	ta := require.New(t)

	st1, err := NewSlimTrie(encode.Int{}, marshalCase.keys, marshalCase.values)
	ta.Nil(err)

	// marshal

	marshalSize := proto.Size(st1)

	buf, err := st1.Marshal()
	ta.Nil(err)
	ta.Equal(len(buf), marshalSize)

	// marshal twice

	buf1, err := proto.Marshal(st1)
	ta.Nil(err)
	ta.Equal(buf, buf1)

	// check version
	r := bytes.NewBuffer(buf)
	n, h, err := pbcmpl.ReadHeader(r)
	ta.Nil(err)
	ta.Equal(int64(32), n)
	ta.Equal(slimtrieVersion, h.GetVersion())

	// unmarshal

	st2, _ := NewSlimTrie(encode.Int{}, nil, nil)
	err = proto.Unmarshal(buf, st2)
	ta.Nil(err)

	checkSlimTrie(st1, st2, t)

	// proto.Unmarshal twice

	err = proto.Unmarshal(buf, st2)
	ta.Nil(err)

	checkSlimTrie(st1, st2, t)

	// Reset()

	st2.Reset()
	empty := &SlimTrie{}
	empty.Leaves.EltEncoder = encode.Int{}
	if !reflect.DeepEqual(st2, empty) {
		t.Fatalf("reset slimtrie error")
	}

	// ensure slimtrie.String()

	_ = st1.String()
}

func TestSlimTrieString(t *testing.T) {

	st, err := NewSlimTrie(encode.Int{}, marshalCase.keys, marshalCase.values)
	if err != nil {
		t.Fatalf("expect no error but: %v", err)
	}

	want := `
#000+1*3
    -001->#001+3*2
              -003->#004=0
                        -006->#007=1
              -004->#005=2
                        -006->#008=3
    -002->#002+2=4
              -006->#006+1=5
                        -006->#009=6
    -003->#003=7`[1:]
	if want != st.String() {
		t.Fatalf("expect: \n%v\n; but: \n%v", want, st.String())
	}

}

func TestSlimTrie_Unmarshal_0_5_0(t *testing.T) {

	// Made with v0.5.0 from:
	//	   st, err := NewSlimTrie(encode.Int{}, marshalCase.keys, marshalCase.values)
	//	   b, err = proto.Marshal(st)
	//	   fmt.Printf("%#v\n", b)
	// Before v0.5.0 a leaf has "step" on it.
	marshaled := []byte{0x31, 0x2e, 0x30, 0x2e, 0x30, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x20, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x22, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x8, 0x6, 0x12, 0x1, 0x77, 0x1a, 0x1, 0x0, 0x22, 0x18,
		0xe, 0x0, 0x1, 0x0, 0x18, 0x0, 0x4, 0x0, 0x40, 0x0, 0x6, 0x0, 0x40, 0x0, 0x7,
		0x0, 0x40, 0x0, 0x8, 0x0, 0x40, 0x0, 0x9, 0x0, 0x31, 0x2e, 0x30, 0x2e, 0x30,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x20, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x1b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x8, 0x8, 0x12,
		0x2, 0xcf, 0x7, 0x1a, 0x1, 0x0, 0x22, 0x10, 0x2, 0x0, 0x4, 0x0, 0x3, 0x0, 0x5,
		0x0, 0x2, 0x0, 0x2, 0x0, 0x2, 0x0, 0x2, 0x0, 0x31, 0x2e, 0x30, 0x2e, 0x30, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x20, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x4b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x8, 0x8, 0x12, 0x2, 0xfc, 0x7,
		0x1a, 0x1, 0x0, 0x22, 0x40, 0x4, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x7, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x6, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}

	// An instance that loads data generated by an older version SlimTrie still works
	st, err := NewSlimTrie(encode.Int{}, nil, nil)
	if err != nil {
		t.Fatalf("expect no error but: %v", err)
	}
	err = proto.Unmarshal(marshaled, st)
	if err != nil {
		t.Fatalf("expect no error but: %v", err)
	}

	for _, ex := range marshalCase.searches {
		lt, eq, gt := st.Search(ex.key)
		rst := searchRst{lt, eq, gt}

		if !reflect.DeepEqual(ex.want, rst) {
			fmt.Println(marshalCase.keys)
			fmt.Println(st)
			t.Fatal("key: ", ex.key, "expected value: ", ex.want, "rst: ", rst)
		}
	}
}

func TestSlimTrie_Unmarshal_0_5_3(t *testing.T) {

	ta := require.New(t)

	// Made with v0.5.3 from:
	//	   st, err := NewSlimTrie(encode.Int{}, marshalCase.keys, marshalCase.values)
	//	   b, err := proto.Marshal(st)
	//	   fmt.Printf("%#v\n", b)
	// v0.5.3 or former uses array.U32 to store Chilldren.
	marshaled := []byte{
		0x31, 0x2e, 0x30, 0x2e, 0x30, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x20, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x22, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x8, 0x6, 0x12, 0x1, 0x77, 0x1a, 0x1, 0x0,
		0x22, 0x18, 0xe, 0x0, 0x1, 0x0, 0x18, 0x0, 0x4, 0x0, 0x40, 0x0, 0x6,
		0x0, 0x40, 0x0, 0x7, 0x0, 0x40, 0x0, 0x8, 0x0, 0x40, 0x0, 0x9, 0x0,
		0x31, 0x2e, 0x30, 0x2e, 0x30, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x20, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x12, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x8, 0x4, 0x12, 0x1, 0x47, 0x1a, 0x1, 0x0,
		0x22, 0x8, 0x2, 0x0, 0x4, 0x0, 0x3, 0x0, 0x2, 0x0, 0x31, 0x2e, 0x30,
		0x2e, 0x30, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x20,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x4b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x8, 0x8, 0x12, 0x2, 0xfc, 0x7, 0x1a, 0x1, 0x0, 0x22, 0x40, 0x4,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x7, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x6,
		0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	}

	// An instance that loads data generated by an older version SlimTrie still works
	st, err := NewSlimTrie(encode.Int{}, nil, nil)
	ta.Nil(err)

	ta.Nil(proto.Unmarshal(marshaled, st))

	for _, ex := range marshalCase.searches {
		lt, eq, gt := st.Search(ex.key)
		rst := searchRst{lt, eq, gt}

		ta.Equal(ex.want, rst)
	}
}

func TestSlimTrieInternalStructre(t *testing.T) {

	ta := require.New(t)

	type ExpectType struct {
		childIndex []int32
		childData  []uint64
		stepIndex  []int32
		stepElts   []uint16
		leafIndex  []int32
		leafData   []uint32
		flags      uint32
		eltWidth   int32
		nodeCnt    int32
	}

	cases := []struct {
		keys   []string
		values []uint32
		ExpectType
	}{
		{
			keys: []string{
				from4bit(1, 2, 3, 4, 0),
				from4bit(1, 2, 3, 4, 1),
				from4bit(1, 2, 3, 4, 2),
				from4bit(1, 2, 3, 4, 3),
				from4bit(1, 3, 3, 5, 4),
			},
			values: []uint32{
				0,
				1,
				2,
				3,
				4,
			},
			ExpectType: ExpectType{
				childIndex: []int32{0, 1},
				childData:  []uint64{12, 15},
				stepIndex:  []int32{0, 1},
				stepElts:   []uint16{1, 2},
				leafIndex:  []int32{2, 3, 4, 5, 6},
				leafData:   []uint32{4, 0, 1, 2, 3},
				flags:      3,
				eltWidth:   16,
				nodeCnt:    7,
			},
		},

		{
			keys: []string{
				from4bit(1, 2, 3),
			},
			values: []uint32{3},
			ExpectType: ExpectType{
				childIndex: []int32{},
				childData:  []uint64{},
				stepIndex:  []int32{},
				stepElts:   []uint16{},
				leafIndex:  []int32{0},
				leafData:   []uint32{3},
				flags:      3,
				eltWidth:   16,
				nodeCnt:    1,
			},
		},
	}

	for _, c := range cases {
		st, err := NewSlimTrie(encode.U32{}, c.keys, c.values)
		ta.Nil(err)

		expectedST, err := NewSlimTrie(encode.U32{}, nil, nil)
		ta.Nil(err)

		expectedST.Children.Flags = c.flags
		expectedST.Children.EltWidth = c.eltWidth

		ch, err := array.NewBitmap16(c.childIndex, c.childData, 16)
		ch.ExtendIndex(c.nodeCnt)
		ta.Nil(err)
		expectedST.Children = *ch

		err = expectedST.Steps.Init(c.stepIndex, c.stepElts)
		ta.Nil(err)
		expectedST.Steps.ExtendIndex(c.nodeCnt)

		err = expectedST.Leaves.Init(c.leafIndex, c.leafData)
		ta.Nil(err)
		expectedST.Leaves.ExtendIndex(c.nodeCnt)

		checkSlimTrie(expectedST, st, t)
	}
}

func checkSlimTrie(st1, st2 *SlimTrie, t *testing.T) {
	if !proto.Equal(&(st1.Children), &(st2.Children)) {
		fmt.Println(st1)
		fmt.Println(st2)
		fmt.Println(pretty.Diff(st1.Children, st2.Children))
		t.Fatalf("Children not the same")
	}

	if !proto.Equal(&(st1.Steps), &(st2.Steps)) {
		fmt.Println(st1)
		fmt.Println(st2)
		fmt.Println(pretty.Diff(st1.Steps, st2.Steps))
		t.Fatalf("Step not the same")
	}

	if !proto.Equal(&st1.Leaves, &st2.Leaves) {
		fmt.Println(st1)
		fmt.Println(st2)
		fmt.Println(pretty.Diff(st1.Leaves, st2.Leaves))
		t.Fatalf("Leaves not the same")
	}
}
