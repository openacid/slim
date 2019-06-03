// Package benchhelper provides utilities for large data set memory or cpu
// benchmark.
package benchhelper

import (
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

	rst := make([]int32, 0)

	for i := min; i < max; i++ {
		if rnd.Float64() < factor {
			rst = append(rst, i)
		}
	}

	return rst
}

func RandI64Slice(min, n, step int64) []int64 {
	rnd := rand.New(rand.NewSource(time.Now().Unix()))

	rst := make([]int64, 0)

	p := min
	for i := 0; i < int(n); i++ {
		s := int64(rnd.Float64() * float64(step))
		p += s
		rst = append(rst, p)
	}

	return rst
}

func RandI32Slice(min, n, step int32) []int32 {
	rnd := rand.New(rand.NewSource(time.Now().Unix()))

	rst := make([]int32, 0)

	p := min
	for i := 0; i < int(n); i++ {
		s := int32(rnd.Float64() * float64(step))
		p += s
		rst = append(rst, p)
	}

	return rst
}

func RandSortedStrings(cnt, leng int, from []byte) []string {

	keys := make(map[string]bool, cnt)

	for i := 0; i < cnt; i++ {
		k := RandString(leng, from)
		if _, ok := keys[k]; ok {
			i--
		} else {
			keys[k] = true
		}
	}

	rsts := make([]string, cnt)
	j := 0
	for i := range keys {
		rsts[j] = i
		j++
	}

	sort.Strings(rsts)
	return rsts
}

func RandByteSlices(cnt, leng int) [][]byte {
	rsts := make([][]byte, cnt)

	for i := int(0); i < cnt; i++ {
		rsts[i] = RandBytes(leng, nil)
	}

	return rsts
}

func RandString(leng int, from []byte) string {
	return string(RandBytes(leng, from))
}

var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func RandBytes(leng int, from []byte) []byte {

	if from == nil {
		from = letters
	}

	b := make([]byte, leng)
	for i := range b {
		b[i] = from[rand.Intn(len(from))]
	}
	return b
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
	tb.SetAutoFormatHeaders(false)
	tb.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	tb.SetCenterSeparator("|")

	return f, tb
}

func NewDataFileTable(fn string) (*os.File, *tablewriter.Table) {

	f := newFile(fn)
	tb := tablewriter.NewWriter(f)
	tb.SetAutoFormatHeaders(false)
	tb.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	tb.SetCenterSeparator("")
	tb.SetColumnSeparator("")
	tb.SetHeaderLine(false)

	return f, tb
}

// WriteTableFiles write a .md file and a .data file
func WriteTableFiles(name string, content interface{}) {
	{
		f, tb := NewMDFileTable(name + ".md")
		defer f.Close()
		tb.SetContent(content)
		tb.Render()
	}

	{
		f, tb := NewDataFileTable(name + ".data")
		defer f.Close()
		tb.SetContent(content)
		tb.Render()
	}

}
