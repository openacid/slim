package main

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/openacid/slim/array/benchmark"
)

func main() {
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

	fn := "report/memusage.md"
	f, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = f.Truncate(0)
	if err != nil {
		panic(err)
	}

	table := tablewriter.NewWriter(f)
	table.SetHeader([]string{"Elt-Size", "Elt-Count", "Load-Factor", "Overhead%"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetColumnAlignment([]int{
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
		tablewriter.ALIGN_RIGHT,
	})
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
