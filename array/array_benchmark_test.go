package array_test

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"

	"github.com/openacid/slim/array"
	"github.com/openacid/slim/encode"
)

func newByteArray32(eSize int, index []int32, elts [][]byte) (*array.Array, error) {
	a := &array.Array{}
	a.EltEncoder = encode.Bytes{Size: eSize}
	err := a.Init(index, elts)
	return a, err
}

func readRss() uint64 {
	var stats runtime.MemStats

	runtime.GC()
	runtime.ReadMemStats(&stats)

	return stats.Alloc
}

func makeData(eltSize int, cnt uint32) [][]byte {
	eltsData := make([][]byte, cnt)

	for i := uint32(0); i < cnt; i++ {
		eltsData[i] = make([]byte, eltSize)
	}

	return eltsData
}

func calcMem(cnt int, indexes []int32, eltSize int, elts [][]byte) uint64 {
	rss1 := readRss()

	arr := []*array.Array{}

	for i := 0; i < cnt; i++ {
		a32, err := newByteArray32(eltSize, indexes, elts)
		if err != nil {
			panic(err)
		}
		arr = append(arr, a32)
	}

	rss2 := readRss()
	var _ []uint64 = arr[0].Bitmaps

	return rss2 - rss1
}

func makeIndexes(maxIdx int32, factor float64) []int32 {
	rnd := rand.New(rand.NewSource(time.Now().Unix()))

	indexes := make([]int32, 0)

	for i := int32(0); i < maxIdx; i++ {
		if rnd.Float64() < factor {
			indexes = append(indexes, i)
		}
	}

	return indexes
}

func calcMemOverHead(factor float64, maxIdx int32, eltSize int) (uint32, float64) {
	cnt := 1024

	indexes := makeIndexes(maxIdx, factor)
	eltCnt := uint32(len(indexes))

	elts := makeData(eltSize, eltCnt)
	actSize := calcMem(cnt, indexes, eltSize, elts)

	dataAvgSize := uint64(eltSize) * uint64(eltCnt)
	actAvgSize := actSize / uint64(cnt)

	overHead := float64(actAvgSize)/float64(dataAvgSize) - 1

	return eltCnt, overHead
}

func benchMemOverHead(eltSize int, maxIdx int32) func(*testing.B) {
	return func(B *testing.B) {
		factor := []float64{1.0, 0.5, 0.2, 0.1, 0.005, 0.001}

		fmt.Printf("%12s%12s%12s%12s\n", "eltSize", "eltCount", "loadFactor", "Overhead")

		for _, f := range factor {

			eltCnt, overHead := calcMemOverHead(f, maxIdx, eltSize)

			oh := fmt.Sprintf("+%d", int(overHead*100))

			fmt.Printf("%12d%12d%12.3f%12s\n", eltSize, eltCnt, f, oh)
		}
	}
}

func BenchmarkArrayMemOverhead(b *testing.B) {
	var cases = []struct {
		eltSize int
		maxIdx  int32
	}{
		{1, 1 << 16},
		{2, 1 << 16},
		{4, 1 << 16},
		{8, 1 << 16},
	}

	for _, c := range cases {
		b.Run("", benchMemOverHead(c.eltSize, c.maxIdx))
	}
}

func BenchmarkArrayGet(b *testing.B) {
	indexes := []int32{0, 5, 9, 203, 400}
	elts := []uint32{12, 15, 19, 120, 300}
	a, _ := array.New(indexes, elts)

	for i := 0; i < b.N; i++ {
		a.Get(5)
	}
}
