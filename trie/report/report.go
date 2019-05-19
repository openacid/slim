// This app runs trie search benchmark.
// Without `go test -bench`, this app show the benchmark result with a chart, which shows a
// better view and is more convenient to compare the cost change with key length and key count.

package main

import (
	"github.com/openacid/slim/benchhelper"
	"github.com/openacid/slim/trie/benchmark"
)

var flg *benchhelper.ReportCmdFlag

func main() {
	flg = benchhelper.InitCmdFlag()
	// getPresent()
	// getAbsent()

	memOverhead()

	// fprGet()

}

func getPresent() {
	if flg.Bench {
		keyCounts := []int{1, 10, 100, 1000, 2000, 5000, 10000, 20000}
		results := benchmark.GetPresent(keyCounts)
		benchhelper.WriteTableFiles("report/bench_get_present", results)
	}

	if flg.Plot {
		script := `
fn = "report/bench_get_present.data"
set yr [50:250]
set xlabel 'key-count (n)'
set ylabel 'Get() present key ns/op' offset 1,0
`
		script += benchhelper.Fformat.JPGHistogramMid
		script += benchhelper.LineStyles.Yellow
		script += benchhelper.Plot.Histogram

		benchhelper.Fplot("report/bench_get_present.jpg", script)
	}
}

func getAbsent() {
	if flg.Bench {
		keyCounts := []int{1, 10, 100, 1000, 2000, 5000, 10000, 20000}
		results := benchmark.GetAbsent(keyCounts)
		benchhelper.WriteTableFiles("report/bench_get_absent", results)
	}

	if flg.Plot {
		script := `
fn = "report/bench_get_absent.data"
set yr [50:300]
set xlabel 'key-count (n)'
set ylabel 'Get() present key ns/op' offset 1,0
`
		script += benchhelper.Fformat.JPGHistogramMid
		script += benchhelper.LineStyles.Orange
		script += benchhelper.Plot.Histogram

		benchhelper.Fplot("report/bench_get_absent.jpg", script)
	}
}

func memOverhead() {
	if flg.BenchMem {
		keyCounts := []int{1000, 2000, 5000}
		results := benchmark.Mem(keyCounts)
		benchhelper.WriteTableFiles("report/mem_usage", results)
	}

	if flg.Plot {
		script := `
fn = "report/mem_usage.data"
set yr [0:50]
set xlabel 'key-count (n)'
set ylabel 'bits/key' offset 1,0
`
		script += benchhelper.Fformat.JPGHistogramTiny
		script += benchhelper.LineStyles.Green
		script += benchhelper.Plot.Histogram

		benchhelper.Fplot("report/mem_usage.jpg", script)
	}
}

func fprGet() {

	if flg.FPR {
		rsts := benchmark.GetFPR([]int{1000, 10000, 20000})
		benchhelper.WriteTableFiles("report/fpr_get", rsts)
	}

	if flg.Plot {
		script := `
fn = "report/fpr_get.data"
set yr [0:100]
set format y "%g%%"
set xlabel 'key-count (n)'
set ylabel 'false positive'
`
		script += benchhelper.Fformat.JPGHistogramTiny
		script += benchhelper.LineStyles.Green
		script += benchhelper.Plot.Histogram

		benchhelper.Fplot("report/fpr_get.jpg", script)
	}
}
