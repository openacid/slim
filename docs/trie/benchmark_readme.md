# SlimTrie Search Performance

We make benchmarks to test the performance of trie search, and, we also benchmark the serach of
B-tree and Map in the same environment to make a comparision.

Benchmarks search one-key in diffrent key length-count pair and get the one-key-search cost.
It can be considered that the less cost the better performance in slimtrie, map and Btree.

There is a [benchmark comparision](benchmark_result.md) report.
Benchmark uses the Go language built-in `map` and the [BTree implementation for Go](https://github.com/google/btree).

Benchmark and optimization steps are as follows:

## Do benchmark

1. `go test cpu.profiling -benchmem -bench=. -run=none`

> Run this command in `slim/trie/` will run key search benchmark of trie, map and BTree.
> It return the one-key search cost of trie, map and Btree.

> In addition, it will generate a `cpu.profiling`, which contains the cpu profile of this benchmark.


2. `go test -cpuprofile trie.cpu.profiling -benchmem -bench=Trie -run=none`

> Can also use `-bench=Trie` to run the benchmark of trie search only.
> It return the key search of trie and the profiling for optimization.

3. `go tool pprof -http [host]:[port] trie.cpu.profiling`

> get an interactive web interface at the specified `host:port` that can be used to navigate through
> various views of the profile `trie.cpu.profiling`.
> It is more convenient to compare and see the profiling details.

> This command needs golang version is *Go 1.11* or the latest go tool `pprof`.

> Flamegraph view is a greate and recommanded way to optimize. Details are in next section.

4. `go run trie_search_cost.go` (optional)

> SlimTrie offers a tool in `slim/tools/app/` to run key search benchmark and output a chart result in
> string. It show a better view than `go test -bench`.


## Profiling with flamegraph

Flamegraph shows a visualization of profiling, which is convenient to observe the profiling details.
[go-torch](https://github.com/uber/go-torch) is a tool to create flamegraph with profiling. And
flamegraph visualization is also available to `go tool pprof` in `Go 1.11`.

*steps to use flamegraph*

1. Install FlameGraph

```
git clone https://github.com/brendangregg/FlameGraph.git

# copy 'FlameGraph/flamegraph.pl' to your `$PATH`, like：
cp FlameGraph/flamegraph.pl /usr/local/bin/
```

2. Install go-torch (Optional)

`go get -v github.com/uber/go-torch`

Step 2 and step 3 are alternative. `go-torch` makes a '.svg' file, if you want to save the 'flamegraph',
this step is recommanded.

3. Install pprof (Optional)

Step 2 and step 3 are alternative. `pprof` is recommanded because it is a official tool, so it will
be more stable in future.

If your golang version is `Go 1.11` or higher, your `go tool pprof` is useful to view flamegraph.
Or you can get the latest `pprof` tool and use it:
```
go get -u github.com/google/pprof

pprof -http=":8088" <profile>
```
then chose `view>flamegraph` on your browser web view.

*pprof Dependence*

It needs [Graphviz](https://www.graphviz.org) to use `pprof` to get an interactive web views of profile,
you should ensure that Graphviz is usable before using `pprof`.

