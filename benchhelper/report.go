package benchhelper

import (
	"flag"
)

type ReportCmdFlag struct {
	Bench    bool
	BenchMem bool
	Plot     bool
}

func InitCmdFlag() *ReportCmdFlag {
	f := &ReportCmdFlag{}
	flag.BoolVar(&f.Bench, "bench", true, "whether to re-benchmark")
	flag.BoolVar(&f.BenchMem, "benchmem", true, "whether to re-benchmark memory usage")
	flag.BoolVar(&f.Plot, "plot", true, "whether to generate plot picture")
	flag.Parse()
	return f
}
