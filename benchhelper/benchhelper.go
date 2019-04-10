// Package benchhelper provides utilities for large data set memory or cpu
// benchmark.
package benchhelper

import (
	crand "crypto/rand"
	"math/rand"
	"os"
	"runtime"
	"sort"
	"time"

	"github.com/openacid/tablewriter"
)

// Allocated returns the in-use heap in bytes.
func Allocated() int64 {
	for i := 0; i < 10; i++ {
		runtime.GC()
	}

	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	return int64(stats.Alloc)
}

func NewBytesSlices(eltSize int, n int) [][]byte {
	slices := make([][]byte, n)

	for i := 0; i < n; i++ {
		slices[i] = make([]byte, eltSize)
	}

	return slices
}

func RandI32SliceBetween(min int32, max int32, factor float64) []int32 {
	rnd := rand.New(rand.NewSource(time.Now().Unix()))

	indexes := make([]int32, 0)

	for i := min; i < max; i++ {
		if rnd.Float64() < factor {
			indexes = append(indexes, i)
		}
	}

	return indexes
}

func RandSortedStrings(cnt, leng int) []string {
	rsts := make([]string, cnt)

	for i := 0; i < cnt; i++ {
		rsts[i] = RandString(leng)
	}

	sort.Strings(rsts)
	return rsts
}

func RandByteSlices(cnt, leng int) [][]byte {
	rsts := make([][]byte, cnt)

	for i := int(0); i < cnt; i++ {
		rsts[i] = RandBytes(leng)
	}

	return rsts
}

func RandString(leng int) string {
	return string(RandBytes(leng))
}

func RandBytes(leng int) []byte {
	bs := make([]byte, leng)
	n, err := crand.Read(bs)
	if err != nil {
		panic(err)
	}
	if n != leng {
		panic("not read enough")
	}
	return bs
}

func newFile(fn string) *os.File {
	f, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	err = f.Truncate(0)
	if err != nil {
		panic(err)
	}
	return f
}

func NewMDFileTable(fn string) (*os.File, *tablewriter.Table) {

	f := newFile(fn)
	tb := tablewriter.NewWriter(f)
	tb.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	tb.SetCenterSeparator("|")

	return f, tb
}

func NewDataFileTable(fn string) (*os.File, *tablewriter.Table) {

	f := newFile(fn)
	tb := tablewriter.NewWriter(f)
	tb.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	tb.SetCenterSeparator("")
	tb.SetColumnSeparator("")

	return f, tb
}

func WriteMDFile(fn string, content interface{}) {
	f, tb := NewMDFileTable(fn)
	defer f.Close()

	tb.SetContent(content)
	tb.Render()
}

func WriteDataFile(fn string, headers []string, content interface{}) {
	f, tb := NewDataFileTable(fn)
	defer f.Close()

	tb.SetContent(content)
	tb.ClearHeader()
	tb.SetHeader(headers)
	tb.SetHeaderLine(false)
	tb.Render()
}
