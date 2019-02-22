package array

import (
	"fmt"
	"math"
	"runtime"
	"testing"
)

func newByteArray32(eSize int, index []uint32, elts [][]byte) (*Array32, error) {
	return New32(ByteConv{EltSize: eSize}, index, elts)
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

func makeDiscrete(max uint32) []uint32 {
	dis := make([]uint32, 0)

	for i := uint32(1); i < max; i++ {

		d := uint32(math.Pow((float64(1)+math.Sqrt(7))/2, float64(i)))
		if d >= max {
			break
		}

		dis = append(dis, d)
	}

	return dis
}

func makeDisIndexes(maxIdx, idxDis uint32) []uint32 {
	index := make([]uint32, 0, maxIdx)

	for i := uint32(0); i < maxIdx; i++ {
		if i%idxDis == 0 {
			index = append(index, i)

		}
	}

	return index
}

func calcMemOverHead(indexes []uint32, eltSize int) float64 {
	caCnt := 1024
	eltCnt := uint32(len(indexes))

	elts := makeData(eltSize, eltCnt)
	totalSize := calcMem(caCnt, indexes, eltSize, elts)

	dataAvgSize := uint64(eltSize) * uint64(eltCnt)
	caAvgSize := totalSize / uint64(caCnt)

	return float64(caAvgSize) / float64(dataAvgSize)
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

	fmt.Printf("%-10s%-10s%-10s%-10s\n", "eltSize", "eltCount", "idxDis", "Overhead")

	for _, c := range cases {
		eltSize, maxIdx := c.eltSize, c.maxIdx

		for _, idxDis := range makeDiscrete(maxIdx) {

			indexes := makeDisIndexes(maxIdx, idxDis)
			overhead := calcMemOverHead(indexes, eltSize)

			eltCnt := uint32(len(indexes))

			fmt.Printf("%-10d%-10d%-10d%-10.3f\n", eltSize, eltCnt, idxDis, overhead)
		}
	}
}
