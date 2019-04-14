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
	getPresent()
	getAbsent()
}

func getPresent() {
	if flg.Bench {
		keyCounts := []int{1, 10, 100, 1000, 2000, 5000, 10000, 20000}
		results := benchmark.GetPresent(keyCounts)

		benchhelper.WriteMDFile("report/bench_get_present.md", results)
		benchhelper.WriteDataFile("report/bench_get_present.data",
			[]string{"key-count", "k=64", "k=128", "k=256"},
			results)
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

		benchhelper.WriteMDFile("report/bench_get_absent.md", results)
		benchhelper.WriteDataFile("report/bench_get_absent.data",
			[]string{"key-count", "k=64", "k=128", "k=256"},
			results)
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
