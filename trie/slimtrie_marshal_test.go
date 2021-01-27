package trie

import (
	"bytes"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/openacid/errors"
	"github.com/openacid/low/pbcmpl"
	"github.com/openacid/low/vers"
	"github.com/openacid/slim/encode"
	"github.com/openacid/testutil"
	"github.com/stretchr/testify/require"
)

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
			// {"ab", searchRst{nil, nil, 0}},
			{"abc", searchRst{nil, 0, 1}},
			{"abd", searchRst{1, 2, 3}},
			{"ac", searchRst{nil, nil, 0}},
			{"adc", searchRst{nil, 0, 1}},
			{"bcd", searchRst{4, 5, 6}},
			{"bce", searchRst{4, 5, 6}},
			{"cde", searchRst{6, 7, nil}},
		},
	}
)

func TestSlimTrie_Unmarshal_incompatible(t *testing.T) {

	ta := require.New(t)

	st1, err := NewSlimTrie(encode.Int{}, marshalCase.keys, marshalCase.values)
	ta.NoError(err)

	buf, err := st1.Marshal()
	ta.NoError(err)

	st2, err := NewSlimTrie(encode.Int{}, nil, nil)
	ta.NoError(err)

	cases := []struct {
		input string
		want  error
	}{
		{slimtrieVersion, nil},
		{"0.5.13", ErrIncompatible},
		{"0.6.0", ErrIncompatible},
		{"0.9.9", ErrIncompatible},
		{"1.0.1", ErrIncompatible},
	}

	for i, c := range cases {

		dd("load from: %s", c.input)

		bad := make([]byte, len(buf))
		copy(bad, buf)

		// clear buf for version
		for i := 0; i < 16; i++ {
			bad[i] = 0
		}
		copy(bad, []byte(c.input))

		err := proto.Unmarshal(bad, st2)
		ta.Equal(c.want, errors.Cause(err), "%d-th: case: %+v", i+1, c)
	}
}

func TestSlimTrie_MarshalUnmarshal(t *testing.T) {

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
	values := []int{0, 1, 2, 3, 4, 5, 6, 7}

	st1, err := NewSlimTrie(encode.Int{}, keys, values)
	ta.NoError(err)

	// marshal

	marshalSize := proto.Size(st1)

	buf, err := st1.Marshal()
	ta.NoError(err)
	ta.Equal(len(buf), marshalSize)

	// marshal twice

	buf1, err := proto.Marshal(st1)
	ta.NoError(err)
	ta.Equal(buf, buf1)

	// check version
	r := bytes.NewBuffer(buf)
	n, h, err := pbcmpl.ReadHeader(r)
	ta.NoError(err)
	ta.Equal(int64(32), n)
	ta.Equal(slimtrieVersion, h.GetVersion())

	// unmarshal

	st2, _ := NewSlimTrie(encode.Int{}, nil, nil)
	err = proto.Unmarshal(buf, st2)
	ta.NoError(err)
	slimtrieEqual(st1, st2, t)

	// proto.Unmarshal twice

	err = proto.Unmarshal(buf, st2)
	ta.NoError(err)
	slimtrieEqual(st1, st2, t)

	// Reset()

	st2.Reset()
	empty := &SlimTrie{
		encoder: encode.Int{},
		vars:    nil,
		levels:  []levelInfo{{0, 0, 0, nil}},
		inner:   &Slim{},
	}
	ta.Equal(empty, st2, "reset")

	// ensure slimtrie.String()

	_ = st1.String()
}

func TestSlimTrie_Marshal_allkeys(t *testing.T) {

	testBigKeySet(t, func(t *testing.T, typ string, keys []string) {
		ta := require.New(t)

		values := makeI32s(len(keys))

		t.Run("Default", func(t *testing.T) {
			st, err := NewSlimTrie(encode.I32{}, keys, values)
			ta.NoError(err)

			testUnknownKeysGRS(t, st, testutil.RandStrSlice(100, 0, 10))
			testPresentKeysGet(t, st, keys, values)

			buf, err := proto.Marshal(st)
			ta.NoError(err)

			st2, err := NewSlimTrie(encode.I32{}, nil, nil)
			ta.NoError(err)

			err = proto.Unmarshal(buf, st2)
			ta.NoError(err)

			testUnknownKeysGRS(t, st2, testutil.RandStrSlice(100, 0, 10))
			testPresentKeysGet(t, st2, keys, values)
		})

		t.Run("Complete", func(t *testing.T) {
			st, err := NewSlimTrie(encode.I32{}, keys, values,
				Opt{Complete: Bool(true)})
			ta.NoError(err)

			testAbsentKeysGRS(t, st, keys)
			testPresentKeysGet(t, st, keys, values)

			buf, err := proto.Marshal(st)
			ta.NoError(err)

			st2, err := NewSlimTrie(encode.I32{}, nil, nil)
			ta.NoError(err)

			err = proto.Unmarshal(buf, st2)
			ta.NoError(err)

			testAbsentKeysGRS(t, st2, keys)
			testPresentKeysGet(t, st2, keys, values)
		})

	})
}

func TestSlimTrie_Unmarshal_old_data(t *testing.T) {

	testOldData(t,
		func(t *testing.T,
			dataSetName, dataOpt, ver string,
			keys []string,
			buf []byte) {

			ta := require.New(t)
			st, err := NewSlimTrie(encode.I32{}, nil, nil)
			ta.NoError(err)

			err = proto.Unmarshal(buf, st)
			ta.NoError(err)

			// < 0.5.10: slimtrie-data-10ll16k-0.5.9
			// => 0.5.10: slimtrie-data-10ll16k-allpref-0.5.10

			if vers.Check(ver, slimtrieVersion, ">=0.5.10") {
				switch dataOpt {
				case "nopref":
					testUnknownKeysGRS(t, st, testutil.RandStrSlice(100, 0, 10))
				case "innpref":
					testUnknownKeysGRS(t, st, testutil.RandStrSlice(100, 0, 10))
				case "allpref":
					// in all prefmode, there is no false positive
					testAbsentKeysGRS(t, st, keys)
				}

			} else {
				// < 0.5.10, slimtrie does not store complete info. There is
				// false positive
				testUnknownKeysGRS(t, st, testutil.RandStrSlice(100, 0, 10))
			}

			testPresentKeysGRS(t, st, keys, makeI32s(len(keys)))

			// test scan

			if dataOpt == "allpref" {
				// only slim with Opt{Complete:true} support scan and iter
				n := len(keys)
				frm := clap(n/5, 0, n)
				to := clap(frm+10, frm, n)
				subTestIter(t, st, keys, keys[frm:to])
				subTestIter(t, st, keys, testutil.RandStrSlice(clap(len(keys), 50, 1024), 0, 10))
				subTestScan(t, st, keys, keys[frm:to])
				subTestScan(t, st, keys, testutil.RandStrSlice(clap(len(keys), 50, 1024), 0, 10))
			}
		})
}

// Just keeps old test.
func TestSlimTrie_Unmarshal_0_5_0(t *testing.T) {

	ta := require.New(t)

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

	st, err := NewSlimTrie(encode.Int{}, nil, nil)
	ta.NoError(err)

	err = proto.Unmarshal(marshaled, st)
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
	_ = wantstr

	dd(st.String())

	ta.Equal(wantstr, st.String())

	for _, ex := range marshalCase.searches {
		lt, eq, gt := st.Search(ex.key)
		rst := searchRst{lt, eq, gt}

		ta.Equal(ex.want, rst, "search for %s", ex.key)
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

	st, err := NewSlimTrie(encode.Int{}, nil, nil)
	ta.NoError(err)

	ta.NoError(proto.Unmarshal(marshaled, st))

	for _, ex := range marshalCase.searches {
		lt, eq, gt := st.Search(ex.key)
		rst := searchRst{lt, eq, gt}

		ta.Equal(ex.want, rst)
	}
}
