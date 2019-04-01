package main

import (
	"fmt"
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

func benGet() {

	var rst testing.BenchmarkResult

	fn := "report/bench_get.md"
	f, table := benchhelper.NewMDFileTable(fn)
	defer f.Close()

	table.SetHeader([]string{"Elt-Type", "Elt-Count", "ns/get"})
	table.SetColumnAlignment(">>>>")

	ns := []int{1, 256, 65536}

	for _, n := range ns {

		{
			indexes := make([]int32, n)
			elts := make([]uint16, n)
			for i := 0; i < n; i++ {
				indexes[i] = int32(i)
			}

			a, err := array.NewU16(indexes, elts)
			if err != nil {
				panic(err)
			}

			rst = testing.Benchmark(func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					a.Get(int32(i % n))
				}
			})
			row := []string{
				"uint16(2)",
				fmt.Sprintf("%d", n),
				fmt.Sprintf("%d", rst.NsPerOp()),
			}

			table.Append(row)
		}

		{
			indexes := make([]int32, n)
			elts := make([]uint64, n)
			for i := 0; i < n; i++ {
				indexes[i] = int32(i)
			}

			a, err := array.NewU64(indexes, elts)
			if err != nil {
				panic(err)
			}

			rst = testing.Benchmark(func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					a.Get(int32(i % n))
				}
			})

			row := []string{
				"uint64(8)",
				fmt.Sprintf("%d", n),
				fmt.Sprintf("%d", rst.NsPerOp()),
			}

			table.Append(row)
		}
	}

	table.Render()
}

func main() {
	memUsage()
	benGet()
}
