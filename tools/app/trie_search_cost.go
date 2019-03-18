// This app runs trie search benchmark.
// Without `go test -bench`, this app show the benchmark result with a chart, which shows a
// better view and is more convenient to compare the cost change with key length and key count.

package main

import (
	"fmt"
	"os"

	"github.com/openacid/slim/trie/benchmark"
)

func main() {
	keyCntLenSearch()
	keyCntIncrSearch()
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

func keyCntIncrSearch() {
	kl := int(256)
	vl := int(2)
	var runs = []benchmark.Config{
		{KeyCnt: 1, KeyLen: kl, ValLen: vl},
		{KeyCnt: 10, KeyLen: kl, ValLen: vl},
		{KeyCnt: 100, KeyLen: kl, ValLen: vl},
		{KeyCnt: 1000, KeyLen: kl, ValLen: vl},
		{KeyCnt: 2000, KeyLen: kl, ValLen: vl},
		{KeyCnt: 5000, KeyLen: kl, ValLen: vl},
		{KeyCnt: 10000, KeyLen: kl, ValLen: vl},
		{KeyCnt: 15000, KeyLen: kl, ValLen: vl},
		{KeyCnt: 20000, KeyLen: kl, ValLen: vl},
	}
	trieCostTb := benchmark.MakeTrieSearchBench(runs)
	chart := benchmark.OutputToChart(
		"search benchmark - key count",
		trieCostTb)

	fmt.Println(chart)
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
