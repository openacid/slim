package trie

import (
	"testing"

	"github.com/openacid/slim/encode"
	"github.com/stretchr/testify/require"
)

func TestSlimTrie_stringValue(t *testing.T) {

	ta := require.New(t)

	keys := []string{
		"Aaron",
		"Agatha",
		"Al",
		"Albert",
	}
	values := []string{
		"abc",
		"def",
		"ghi",
		"jkl",
	}

	{
		st, err := NewSlimTrie(encode.String16{}, keys, values)
		ta.NoError(err)

		dd(keys)
		dd(st)

		for i, key := range keys {

			v, found := st.RangeGet(key)
			ta.True(found, "RangeGet:%v", key)
			ta.Equal(values[i], v, "RangeGet:%v", key)
		}
	}

}
