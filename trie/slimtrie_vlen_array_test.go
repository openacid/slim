package trie

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVLenArray_get(t *testing.T) {

	ta := require.New(t)

	// Fixed size
	{
		elts := [][]byte{
			{'a', 'b'},
			{}, // empty
			{},
			{'c', 'd'},
			{'e', 'f'},
			{},
		}

		va := newVLenArray(elts)

		ta.Equal(int32(6), va.N)
		ta.Equal(int32(3), va.EltCnt)
		ta.Equal(int32(2), va.FixedSize)

		ta.Equal([]byte{'a', 'b'}, va.get(0))
		ta.Equal([]byte{}, va.get(1))
		ta.Equal([]byte{}, va.get(2))
		ta.Equal([]byte{'c', 'd'}, va.get(3))
		ta.Equal([]byte{'e', 'f'}, va.get(4))
		ta.Equal([]byte{}, va.get(5))
	}

	// Var-len size
	{
		elts := [][]byte{
			{'a', 'b', 'c'},
			{}, // empty
			{},
			{'c', 'd'},
			{'e', 'f'},
			{},
		}

		va := newVLenArray(elts)

		ta.Equal(int32(6), va.N)
		ta.Equal(int32(3), va.EltCnt)
		ta.Equal(int32(0), va.FixedSize)

		ta.Equal([]byte{'a', 'b', 'c'}, va.get(0))
		ta.Equal([]byte{}, va.get(1))
		ta.Equal([]byte{}, va.get(2))
		ta.Equal([]byte{'c', 'd'}, va.get(3))
		ta.Equal([]byte{'e', 'f'}, va.get(4))
		ta.Equal([]byte{}, va.get(5))
	}

	// dd(st)
}
