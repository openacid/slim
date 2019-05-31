package trie

import (
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
