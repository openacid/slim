package trie

import (
	"testing"

	"github.com/openacid/slim/encode"
	"github.com/stretchr/testify/require"
)

func fibhash64(i uint64) uint64 {
	return i * 11400714819323198485
}

func TestSlimTrie_GetI8(t *testing.T) {

	ta := require.New(t)

	keys := getKeys("20kvl10")
	values := make([]int8, len(keys))
	for i := 0; i < len(keys); i++ {
		values[i] = int8(fibhash64(uint64(i)))
	}

	st, err := NewSlimTrie(encode.I8{}, keys, values)
	ta.NoError(err)

	testUnknownKeysGRS(t, st, randVStrings(len(keys)*5, 0, 10))

	for i, key := range keys {

		// dd("test Get: present: %s", key)

		v, found := st.GetI8(key)
		ta.True(found, "Get:%v", key)
		ta.Equal(values[i], v, "Get:%v", key)
	}
}

func TestSlimTrie_GetI16(t *testing.T) {

	ta := require.New(t)

	keys := getKeys("20kvl10")
	values := make([]int16, len(keys))
	for i := 0; i < len(keys); i++ {
		values[i] = int16(fibhash64(uint64(i)))
	}

	st, err := NewSlimTrie(encode.I16{}, keys, values)
	ta.NoError(err)

	testUnknownKeysGRS(t, st, randVStrings(len(keys)*5, 0, 10))

	for i, key := range keys {

		// dd("test Get: present: %s", key)

		v, found := st.GetI16(key)
		ta.True(found, "Get:%v", key)
		ta.Equal(values[i], v, "Get:%v", key)
	}
}

func TestSlimTrie_GetI32(t *testing.T) {

	ta := require.New(t)

	keys := getKeys("20kvl10")
	values := make([]int32, len(keys))
	for i := 0; i < len(keys); i++ {
		values[i] = int32(fibhash64(uint64(i)))
	}

	st, err := NewSlimTrie(encode.I32{}, keys, values)
	ta.NoError(err)

	testUnknownKeysGRS(t, st, randVStrings(len(keys)*5, 0, 10))

	for i, key := range keys {

		// dd("test Get: present: %s", key)

		v, found := st.GetI32(key)
		ta.True(found, "Get:%v", key)
		ta.Equal(values[i], v, "Get:%v", key)
	}
}

func TestSlimTrie_GetI64(t *testing.T) {

	ta := require.New(t)

	keys := getKeys("20kvl10")
	values := make([]int64, len(keys))
	for i := 0; i < len(keys); i++ {
		values[i] = int64(fibhash64(uint64(i)))
	}

	st, err := NewSlimTrie(encode.I64{}, keys, values)
	ta.NoError(err)

	testUnknownKeysGRS(t, st, randVStrings(len(keys)*5, 0, 10))

	for i, key := range keys {

		// dd("test Get: present: %s", key)

		v, found := st.GetI64(key)
		ta.True(found, "Get:%v", key)
		ta.Equal(values[i], v, "Get:%v", key)
	}
}
