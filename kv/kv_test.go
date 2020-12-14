package kv

import (
	"fmt"
	"testing"

	"github.com/openacid/low/size"
	"github.com/openacid/testkeys"
	"github.com/stretchr/testify/require"
)

func TestKV_Get(t *testing.T) {
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

	kv := NewKV(keys, values)

	for i, k := range keys {
		v, found := kv.Get(k)
		ta.Equal(v, values[i])
		ta.True(found)
	}

}

func makeI32s(n int) []int32 {

	values := make([]int32, n)
	for i := int32(0); i < int32(n); i++ {
		values[i] = i
	}
	return values

}

var Output int32

func BenchmarkKV_Get(b *testing.B) {
	keys := testkeys.Load("200kweb2")
	// keys := testkeys.Load("1mvl5_10")
	values := makeI32s(len(keys))
	kv := NewKV(keys, values)
	// fmt.Println((size.Of(kv) - 16001304 - 4*1024*1024) * 8 / len(keys))
	fmt.Println("index:", (size.Of(kv.st))*8/len(keys), "bit")
	fmt.Println("keys:", (size.Of(kv.keys))/len(keys), "byte")
	fmt.Println("values:", (size.Of(kv.values))/len(keys), "byte")
	fmt.Println("total:", (size.Of(kv))/len(keys), "byte")

	b.ResetTimer()

	s := int32(0)
	for i := 0; i < b.N; i++ {

		v, _ := kv.Get(keys[i&(128*1024-1)])
		s += v
	}

	Output = s
}
