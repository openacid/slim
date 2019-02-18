# Trie Search Performance

## Performance test of trie search

While doing benchmark of trie search, we also benchmark the serach of B-tree and Map in the same
environment. There is a [benchmark comparision](benchmark_result.md) report.

Benchmark uses the Go language built-in `map` and the [BTree implementation for Go](https://github.com/google/btree).

1. `go test -cpuprofile cpu.profiling -benchmem -bench=. -run=none`

> run this command in `slim/trie/` will run key search benchmark of trie, map and BTree.

> In addition, it will generate a `cpu.profiling`, which contains the cpu profile of this benchmark.


## Performance optimize

1. `go test -cpuprofile trie.cpu.profiling -benchmem -bench=Trie -run=none`

> this command use `-bench=Trie` to get the profiling of trie search only.

2. `go tool pprof -http [host]:[port] trie.cpu.profiling`

> get an interactive web interface at the specified `host:port` that can be used to navigate through
> various views of the profile `trie.cpu.profiling`.

> It is more convenient to compare and see the profiling details.

> It needs go language version is *Go 1.11* or the latest go tool `pprof`.


## Use Flamegraph

[go-torch](https://github.com/uber/go-torch) is a flamegraph visualization of profiling. And
flamegraph visualization is also available to `go tool pprof` in `Go 1.11`.

1. Install FlameGraph

```
git clone https://github.com/brendangregg/FlameGraph.git

# copy 'FlameGraph/flamegraph.pl' to your `$PATH`, like：
cp FlameGraph/flamegraph.pl /usr/local/bin/
```

2. Install go-torch
`go get -v github.com/uber/go-torch`

3. Install pprof

If your go language version is `Go 1.11` or higher, your `go tool pprof` is useful to view flamegraph.
Or you can get the latest `pprof` tool and use it:
```
go get -u github.com/google/pprof

pprof -http=":8088" <profile>
```
then chose `view>flamegraph` on your browser web view.

It needs [Graphviz](https://www.graphviz.org) to use `pprof` to get an interactive web views of profile,
you should ensure that Graphviz is usable before using `pprof`.

