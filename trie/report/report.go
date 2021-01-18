// This app runs trie search benchmark.
// Without `go test -bench`, this app show the benchmark result with a chart, which shows a
// better view and is more convenient to compare the cost change with key length and key count.

package main

import (
	"fmt"
	"strings"

	"github.com/openacid/slim/benchhelper"
	"github.com/openacid/slim/trie/benchmark"
)

var (
	// different key counts for benchmark.
	keyCounts = []int{
		1000,
		10 * 1000,
		100 * 1000,
		1000 * 1000,
	}
)

var flg *benchhelper.ReportCmdFlag

func main() {
	flg = benchhelper.InitCmdFlag()

	for _, workload := range []string{
		"zipf", "scan",
	} {
		getPresent(workload)
		getAbsent(workload)
		getMapSlimArrayBtree(workload)

	}
	memOverhead()
	fprGet()
}

func getPresent(workload string) {

	fn := fmt.Sprintf("report/bench_get_present_%s", workload)

	if flg.Bench {
		results := benchmark.BenchGet(keyCounts, "present", workload)
		benchhelper.WriteTableFiles(fn, results)
	}

	if flg.Plot {
		script := strings.Join([]string{
			fmt.Sprintf(`fn = "%s.data"`, fn),
			`set yr [0:300]`,
			`set xlabel 'key-count: n'`,
			`set ylabel 'ns/Get() present key' offset 1,0`,
			``,
		}, "\n")
		script += benchhelper.Fformat.JPGHistogramTiny
		script += benchhelper.LineStyles.Green
		script += benchhelper.Plot.Histogram

		benchhelper.Fplot(fn+".jpg", script)
	}
}

func getAbsent(workload string) {

	fn := fmt.Sprintf("report/bench_get_absent_%s", workload)

	if flg.Bench {
		results := benchmark.BenchGet(keyCounts, "absent", workload)
		benchhelper.WriteTableFiles(fn, results)
	}

	if flg.Plot {
		script := strings.Join([]string{
			fmt.Sprintf(`fn = "%s.data"`, fn),
			`set yr [0:300]`,
			`set xlabel 'key-count: n'`,
			`set ylabel 'ns/Get() absent key' offset 1,0`,
			``,
		}, "\n")
		script += benchhelper.Fformat.JPGHistogramTiny
		script += benchhelper.LineStyles.Green
		script += benchhelper.Plot.Histogram

		benchhelper.Fplot(fn+".jpg", script)
	}
}

func getMapSlimArrayBtree(workload string) {

	fn := fmt.Sprintf("report/bench_msab_present_%s", workload)

	if flg.Bench {
		results := benchmark.GetMapSlimArrayBtree(keyCounts, workload)
		benchhelper.WriteTableFiles(fn, results)
	}

	if flg.Plot {
		script := strings.Join([]string{
			fmt.Sprintf(`fn = "%s.data"`, fn),
			`set yr [0:1000]`,
			`set xlabel 'key-count: n'`,
			`set ylabel 'ns/Get()' offset 1,0`,
			``,
		}, "\n")
		script += benchhelper.Fformat.JPGHistogramTiny
		script += benchhelper.LineStyles.Blue
		script += benchhelper.Plot.Histogram

		benchhelper.Fplot(fn+".jpg", script)
	}
}

func memOverhead() {
	if flg.BenchMem {
		keyCounts := []int{1000, 10 * 1000, 100 * 1000, 1000 * 1000}
		results := benchmark.Mem(keyCounts)
		benchhelper.WriteTableFiles("report/mem_usage", results)
	}

	if flg.Plot {
		script := `
fn = "report/mem_usage.data"
set yr [0:30]
set xlabel 'key-count: n'
set ylabel 'bits/key' offset 1,0
`
		script += benchhelper.Fformat.JPGHistogramTiny
		script += benchhelper.LineStyles.Yellow
		script += benchhelper.Plot.Histogram

		benchhelper.Fplot("report/mem_usage.jpg", script)
	}
}

func fprGet() {

	if flg.FPR {
		rsts := benchmark.GetFPR(keyCounts)
		benchhelper.WriteTableFiles("report/fpr_get", rsts)
	}

	if flg.Plot {
		script := `
fn = "report/fpr_get.data"
set yr [0:0.05]
set format y "%g%%"
set xlabel 'key-count: n'
set ylabel 'false positive'
`
		script += benchhelper.Fformat.JPGHistogramTiny
		script += benchhelper.LineStyles.Purple
		script += benchhelper.Plot.Histogram

		benchhelper.Fplot("report/fpr_get.jpg", script)
	}
}
