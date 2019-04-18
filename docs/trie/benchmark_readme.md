# SlimTrie Search Performance

We make benchmarks to test the performance of trie search, and, we also benchmark the search of
B-tree and Map in the same environment to make a comparison.

Benchmarks search one-key in different key length-count pair and get the one-key-search cost.
It can be considered that the less cost the better performance in slimtrie, map and Btree.

There is a [benchmark comparison](benchmark_result.md) report.
Benchmark uses the Go language built-in `map` and the [BTree implementation for Go](https://github.com/google/btree).

Benchmark and optimization steps are as follows:


## Do benchmark

1. `go test -benchmem -bench=. -run=none`

> Run this command in `slim/trie/` will run key search benchmark of trie, map and BTree.
> It return the one-key search cost of trie, map and Btree.

2. `go test -cpuprofile trie.cpu.profiling -benchmem -bench=Trie -run=none`

> Can also use `-bench=Trie` to run the benchmark of trie search only.
> It return the key search cost of trie.
> And use `-cpuprofile` to get the cpu profiling for optimization.


## Profiling with flamegraph

You must need a tool to view the profile and optimize.
Flamegraph shows a visualization of profile, which is convenient to observe the profile details.

*steps to use flamegraph*

1. Install FlameGraph

```
git clone https://github.com/brendangregg/FlameGraph.git

# copy 'FlameGraph/flamegraph.pl' to your `$PATH`, like：
cp FlameGraph/flamegraph.pl /usr/local/bin/
```

2. Install a visualization tool to use flamegraph

There are 2 recommanded tools, you can choose *one of those two*.

[go-torch](https://github.com/uber/go-torch) is a tool to create flamegraph with profiling.
And flamegraph visualization is also available to `go tool pprof` in `Go 1.11`.

*pprof*

`pprof` is recommanded because it is an official tool, so it will be more stable in future.
And, if your golang version is `Go 1.11` or higher, your `go tool pprof` is useful to view flamegraph.
Or you can get the latest `pprof` tool:
```
go get -u github.com/google/pprof
```

It also needs [Graphviz](https://www.graphviz.org) to use `pprof` to get an interactive web views of profile,
you should ensure that Graphviz is usable before using `pprof`.

use `pprof` as:

```
# if your Golang version >= 1.11
go tool pprof -http=":8088" <profile>

# if you use pprof tool
pprof -http=":8088" <profile>
```
you get an interactive web interface that can be used to navigate through various views,
then chose `view > flamegraph` on your browser web view to use flamegraph.

*go-torch*

`go-torch` makes a '.svg' file, if you want to save the 'flamegraph', install it:

```
go get -v github.com/uber/go-torch
```

and use it:

```
go-torch -b <profile> -f <profile.svg>
```
you get the flamegraph of profile, then:
```
open <profile.svg>
```
to view flamegraph.


## Trie search cost

SlimTrie offers a tool in `slim/tools/app/` to run key search benchmark and output a chart result in
string. It show a better view than `go test -bench`.

This cammand `go run tools/app/trie_search_cost.go` runs the slimtrie key search benchmark, and then
get a result like:

```
cost of trie search with existing & existent key:

┌────────┬────────┬────────────────────┬───────────────────────┐
│ KeyCnt │ KeyLen │ ExsitingKeyNsPerOp │ NonexsitentKeyNsPerOp │
├────────┼────────┼────────────────────┼───────────────────────┤
│ 1      │ 1024   │ 116                │ 85                    │
│ 10     │ 1024   │ 119                │ 96                    │
│ 100    │ 1024   │ 224                │ 197                   │
│ 1000   │ 1024   │ 222                │ 206                   │
│ 1000   │ 512    │ 272                │ 241                   │
│ 1000   │ 256    │ 222                │ 190                   │
└────────┴────────┴────────────────────┴───────────────────────┘
```
