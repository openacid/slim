#!/bin/sh

go test . -run=none -bench=BenchmarkKV_Get -benchmem -cpuprofile prof.cpu -memprofile prof.mem

go tool pprof -output cpu.svg -svg prof.cpu
go tool pprof -output mem.svg -svg prof.mem
