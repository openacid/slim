package array

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

func newByteArray32(eSize int, index []uint32, elts [][]byte) (*Array32, error) {
	return New(ByteConv{EltSize: eSize}, index, elts)
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

func calcMem(caCnt int, indexes []uint32, eltSize int, elts [][]byte) uint64 {
	rss1 := readRss()

	arr := []*Array32{}

	for i := 0; i < caCnt; i++ {
		a32, _ := newByteArray32(eltSize, indexes, elts)
		arr = append(arr, a32)
	}

	rss2 := readRss()
	var _ []uint64 = arr[0].Bitmaps

	return rss2 - rss1
}

func makeIndexes(maxIdx uint32, factor float64) []uint32 {
	rnd := rand.New(rand.NewSource(time.Now().Unix()))

	indexes := make([]uint32, 0)

	for i := uint32(0); i < maxIdx; i++ {
		if rnd.Float64() < factor {
			indexes = append(indexes, i)
		}
	}

	return indexes
}

func calcMemOverHead(factor float64, maxIdx uint32, eltSize int) (uint32, float64) {
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

func benchMemOverHead(eltSize int, maxIdx uint32) func(*testing.B) {
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

func BenchmarkMemOverhead(b *testing.B) {
	var cases = []struct {
		eltSize int
		maxIdx  uint32
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
