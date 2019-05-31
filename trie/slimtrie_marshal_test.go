package trie

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/openacid/errors"
	"github.com/openacid/slim/encode"
	"github.com/stretchr/testify/require"
)

func TestSlimTrie_Unmarshal_incompatible(t *testing.T) {

	ta := require.New(t)

	st1, err := NewSlimTrie(encode.Int{}, marshalCase.keys, marshalCase.values)
	ta.Nil(err)

	buf, err := st1.Marshal()
	ta.Nil(err)

	st2, _ := NewSlimTrie(encode.Int{}, nil, nil)

	cases := []struct {
		input string
		want  error
	}{
		{"1.0.0", nil},
		{"0.5.8", nil},
		{"0.5.9", ErrIncompatible},
		{"0.9.9", ErrIncompatible},
		{"1.0.1", ErrIncompatible},
	}

	for i, c := range cases {
		bad := make([]byte, len(buf))
		copy(bad, buf)
		copy(buf, []byte(c.input))
		err := proto.Unmarshal(buf, st2)
		ta.Equal(c.want, errors.Cause(err), "%d-th: case: %+v", i+1, c)
	}
}

func TestSlimTrie_Unmarshal_old_data(t *testing.T) {

	ta := require.New(t)

	folder := "testdata/"
	finfos, err := ioutil.ReadDir(folder)
	ta.Nil(err)

	for _, finfo := range finfos {

		fn := finfo.Name()

		if !strings.HasPrefix(fn, "slimtrie-data-") {
			continue
		}

		path := filepath.Join(folder, fn)
		b, err := ioutil.ReadFile(path)
		ta.Nil(err)

		st, err := NewSlimTrie(encode.I32{}, nil, nil)
		ta.Nil(err)

		err = proto.Unmarshal(b, st)
		ta.Nil(err)

		keys := keys50k
		for i, key := range keys {
			v, found := st.Get(key)
			ta.True(found)
			ta.Equal(int32(i), v)
		}
	}
}
