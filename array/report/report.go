package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/openacid/slim/array"
	"github.com/openacid/slim/array/benchmark"
	"github.com/openacid/slim/benchhelper"
)

func memUsage() {
	var cases = []struct {
		eltSize int
		maxIdx  int32
	}{
		{4, 1 << 16},
		{8, 1 << 16},
	}

	factor := []float64{1.0, 0.5, 0.2, 0.1, 0.005, 0.001}

	usages := []*benchmark.MemoryUsage{}

	for _, c := range cases {
		for _, f := range factor {
			usage := benchmark.CollectMemoryUsage(f, c.maxIdx, c.eltSize)
			usages = append(usages, usage)
		}
	}

	fn := "report/mem_usage.md"
	f, table := benchhelper.NewMDFileTable(fn)
	defer f.Close()

	table.SetHeader([]string{"Elt-Size", "Elt-Count", "Load-Factor", "Overhead%"})
	table.SetColumnAlignment(">>>>")
	for _, u := range usages {
		row := []string{
			fmt.Sprintf("%d", u.EltSize),
			fmt.Sprintf("%d", u.EltCnt),
			fmt.Sprintf("%.1f%%", u.LoadFactor*100),
			fmt.Sprintf("+%.1f%%", u.Overhead*100),
		}

		table.Append(row)
	}
	table.Render()
}

var Output int

func benGet() {

	rows := [][]string{}

	ns := []int{1, 256, 65536}

	for _, n := range ns {

		indexes := make([]int32, n)
		for i := 0; i < n; i++ {
			indexes[i] = int32(i)
		}

		row := []string{fmt.Sprintf("%d", n)}

		{
			a, err := array.NewU16(indexes, make([]uint16, n))
			if err != nil {
				panic(err)
			}
			rst := testing.Benchmark(func(b *testing.B) {
				var s uint16
				for i := 0; i < b.N; i++ {
					v, _ := a.Get(int32(i % n))
					s += v
				}
				Output = int(s)
			})
			row = append(row, fmt.Sprintf("%d", rst.NsPerOp()))
		}

		{
			a, err := array.NewU32(indexes, make([]uint32, n))
			if err != nil {
				panic(err)
			}
			rst := testing.Benchmark(func(b *testing.B) {
				var s uint32
				for i := 0; i < b.N; i++ {
					v, _ := a.Get(int32(i % n))
					s += v
				}
				Output = int(s)
			})

			row = append(row, fmt.Sprintf("%d", rst.NsPerOp()))
		}

		{
			a, err := array.NewU64(indexes, make([]uint64, n))
			if err != nil {
				panic(err)
			}
			rst := testing.Benchmark(func(b *testing.B) {
				var s uint64
				for i := 0; i < b.N; i++ {
					v, _ := a.Get(int32(i % n))
					s += v
				}
				Output = int(s)
			})

			row = append(row, fmt.Sprintf("%d", rst.NsPerOp()))
		}

		rows = append(rows, row)

	}

	fn := "report/bench_get.md"
	f, table := benchhelper.NewMDFileTable(fn)
	defer f.Close()

	table.SetHeader([]string{"key-count", "u16", "u32", "u64"})
	table.SetColumnAlignment(">>>>")
	table.AppendBulk(rows)
	table.Render()

	lines := []string{}
	for _, row := range rows {
		line := strings.Join(row, " ")
		lines = append(lines, line)
	}

	cont := strings.Join(lines, "\n")

	fn = "report/bench_get.data"
	err := ioutil.WriteFile(fn, []byte(cont), 0777)
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("generating mem usage benchmark...")
	memUsage()

	fmt.Println("generating get() benchmark...")
	benGet()

	fmt.Println("generating get() benchmark chart...")
	scriptGet := benchhelper.Fformat.JPGHistogram +
		benchhelper.LineStyles.Green + `
set yr [5:20]
set xlabel 'key-count (n)'
set ylabel 'Get() ns/op' offset 1,0

plot "report/bench_get.data" using 2:xtic(1) with histogram linestyle 1 title "u16", \
     ''						 using 3:xtic(1) with histogram linestyle 2 title "u32", \
     ''						 using 4:xtic(1) with histogram linestyle 3 title "u64"
`

	benchhelper.Fplot("report/bench_get.jpg", scriptGet)
}
