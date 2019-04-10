// This app runs trie search benchmark.
// Without `go test -bench`, this app show the benchmark result with a chart, which shows a
// better view and is more convenient to compare the cost change with key length and key count.

package main

import (
	"fmt"
	"os"

	"github.com/openacid/slim/benchhelper"
	"github.com/openacid/slim/trie/benchmark"
)

func main() {
	keyCntLenSearch()
	keyGetPresent()
}

func keyCntLenSearch() {
	var runs = []benchmark.Config{
		{KeyCnt: 1, KeyLen: 1024, ValLen: 2},
		{KeyCnt: 10, KeyLen: 1024, ValLen: 2},
		{KeyCnt: 100, KeyLen: 1024, ValLen: 2},
		{KeyCnt: 1000, KeyLen: 1024, ValLen: 2},
		{KeyCnt: 1000, KeyLen: 512, ValLen: 2},
		{KeyCnt: 1000, KeyLen: 256, ValLen: 2},
	}

	trieCostTb := benchmark.MakeTrieSearchBench(runs)

	chart := benchmark.OutputToChart(
		"cost of trie search with existing & inexistent key",
		trieCostTb)

	fmt.Println(chart)

	resultFn := "trie_search_cost.chart"
	writeToFile(resultFn, chart)
	fmt.Println("result also wirte to", resultFn)
}

func keyGetPresent() {
	// TODO add absent key search bench

	keyCounts := []int{1, 10, 100, 1000, 2000, 5000, 10000, 20000}
	trieCostTb := benchmark.GetPresent(keyCounts)

	benchhelper.WriteMDFile("bench_get_present.md", trieCostTb)
	benchhelper.WriteDataFile("bench_get_present.data",
		[]string{"key-count", "k=64", "k=128", "k=256"},
		trieCostTb)

	scriptGetPresent := `
fn = "bench_get_present.data"
set yr [50:250]
set xlabel 'key-count (n)'
set ylabel 'Get() present key ns/op' offset 1,0
`
	scriptGetPresent += benchhelper.Fformat.JPGHistogramMid
	scriptGetPresent += benchhelper.LineStyles.Green
	scriptGetPresent += benchhelper.Plot.Histogram

	benchhelper.Fplot("bench_get_present.jpg", scriptGetPresent)

}

// writeToFile wirte a string body to file.
func writeToFile(resultFn, body string) {

	f, err := os.Create(resultFn)
	if err != nil {
		fmt.Printf("failed to open result file: %s", resultFn)
	}
	defer f.Close()

	fmt.Fprintln(f, body)
}
