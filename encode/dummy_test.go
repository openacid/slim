package encode_test

import (
	"testing"

	"github.com/openacid/slim/encode"
	"github.com/stretchr/testify/require"
)

func TestDummy(t *testing.T) {

	ta := require.New(t)

	m := encode.Dummy{}

	cases := []struct {
		input interface{}
		want  int
	}{
		{nil, 0},
		{true, 0},
		{1, 0},
		{int32(1), 0},
		{[]byte(""), 0},
		{[]byte("a"), 0},
		{[]byte("abc"), 0},
	}

	for i, c := range cases {
		ta.Equal(c.want, m.GetSize(c.input), "%d-th: case: %+v", i+1, c)
		ta.Equal(c.want, m.GetEncodedSize([]byte{1, 2, 3}), "%d-th: case: %+v", i+1, c)
		ta.Equal([]byte{}, m.Encode(c.input), "%d-th: case: %+v", i+1, c)
		n, got := m.Decode([]byte{})
		ta.Equal(0, n, "%d-th: case: %+v", i+1, c)
		ta.Nil(got, "%d-th: case: %+v", i+1, c)
	}
}
