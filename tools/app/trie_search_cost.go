// This app runs trie search benchmark.
// Without `go test -bench`, this app show the benchmark result with a chart, which shows a
// better view and is more convenient to compare the cost change with key length and key count.

package main

import (
	"fmt"
	"os"

	"github.com/openacid/slim/trie"
)

func main() {
	trieCostTb := trie.MakeTrieSearchBench()

	chart := trieCostTb.OutputToChart()

	fmt.Println(chart)

	resultFn := "trie_search_cost.chart"
	writeToFile(resultFn, chart)
	fmt.Println("result also wirte to", resultFn)
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
