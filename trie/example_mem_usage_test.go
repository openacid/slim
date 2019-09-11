package trie

import (
	"fmt"

	"github.com/openacid/low/size"
	"github.com/openacid/slim/encode"
)

var (
	m = encode.String16{}
)

func Example_memoryUsage() {

	data := make([]byte, 0)

	keys := getKeys("50kl10")
	offsets := []uint16{}

	for i, k := range keys {
		v := fmt.Sprintf("is the %d-th word", i)

		offsets = append(offsets, uint16(len(data)))

		data = append(data, m.Encode(k)...)
		data = append(data, m.Encode(v)...)

	}

	vsize := 2.0

	st, _ := NewSlimTrie(encode.U16{}, keys, offsets)
	ksize := size.Of(keys)

	sz := size.Of(st)

	avgIdxLen := float64(sz)/float64(len(keys)) - vsize
	avgKeyLen := float64(ksize) / float64(len(keys))

	ratio := avgIdxLen / avgKeyLen * 100

	fmt.Printf(
		"Orignal:: %.1f byte/key --> SlimTrie index: %.1f byte/index\n"+
			"Saved %.1f%%",
		avgKeyLen,
		avgIdxLen,
		100-ratio,
	)

}
