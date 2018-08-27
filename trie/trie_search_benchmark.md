### trie search 性能测试的步骤

1. `go test -cpuprofile prof.cpu -benchmem -bench=Trie -run=none`

> 会运行 BenchmarkTrieSearc 和 BenchmarkTrieSearcha 两个函数，同时生成 prof.cpu 文件其中保留着此次测
> 试的cpu 时间采样信息。

2. `go-torch -b prof.cpu -f prof.cpu.svg`

> 依据步骤 1 中生成的 prof.cpu 文件，生成时间采样的火焰图 prof.cpu.svg。这个图更直观的展示了各部分所占
> 用的时间的分布。

3. `go tool pprof prof.cpu`

> 依据步骤 1 中生成的 prof.cpu 文件，可以详细查看每个函数内部运行的时间占用。
> 例如： `list SearchString` 则展示出 SearchString 这个函数内部的每部分占用时间。

### 生成火焰图需要的依赖

1. 安装 `FlameGraph`

`git clone https://github.com/brendangregg/FlameGraph.git`

2. 把 `FlameGraph` 下的 `flamegraph.pl` 拷贝到 $PATH 路径中，例如：
`cp flamegraph.pl /usr/local/bin`

3. 安装 `go-torch`
`go get -v github.com/uber/go-torch`


### test search exsisted key between trie, map, array and btree:

| key length | key count | trie cost (ns) | map cost (ns) | array cost (ns) | btree cost (ns) |
| ---:       | ---:      | ---:           | ---:          | ---:            | ---:            |
| 1024       | 1         | 86.367         | 5.037         | 19.767          | 36.933          |
| 1024       | 10        | 90.733         | 59.767        | 53.067          | 99.533          |
| 1024       | 100       | 123.333        | 60.100        | 98.200          | 240.667         |
| 1024       | 1000      | 157.000        | 63.567        | 146.667         | 389.667         |
| 512        | 1000      | 152.667        | 40.033        | 149.667         | 363.000         |
| 256        | 1000      | 152.333        | 28.833        | 141.333         | 332.333         |

### test search not exsisted key between trie, map, array and btree:

| key length | key count | trie cost (ns) | map cost (ns) | array cost (ns) | btree cost (ns) |
| ---:       | ---:      | ---:           | ---:          | ---:            | ---:            |
| 1024       | 1         | 60.267         | 10.900        | 41.267          | 92.567          |
| 1024       | 10        | 63.667         | 68.733        | 67.967          | 178.000         |
| 1024       | 100       | 100.066        | 71.699        | 103.000         | 297.333         |
| 1024       | 1000      | 134.667        | 68.833        | 129.000         | 441.667         |
| 512        | 1000      | 131.333        | 49.300        | 117.667         | 389.667         |
| 256        | 1000      | 136.333        | 38.300        | 108.667         | 362.667         |

### ratio of search cost between trie, map, array and btree:

| key length | key count | trie cost / map cost | trie cost / array cost | trie cost / btree cost |
| ---:       | ---:      | ---:                 | ---:                   | ---:                   |
| 1024       | 1         | 1714.7 %             | 436.9 %                | 233.8 %                |
| 1024       | 10        | 151.8 %              | 170.9 %                | 91.2 %                 |
| 1024       | 100       | 205.2 %              | 125.6 %                | 51.2 %                 |
| 1024       | 1000      | 246.9 %              | 107.1 %                | 42.3 %                 |
| 512        | 1000      | 381.4 %              | 102.0 %                | 42.1 %                 |
| 256        | 1000      | 528.3 %              | 107.9 %                | 45.8 %                 |

### ratio of search not exsisted key cost between trie, map, array and btree

| key length | key count | trie cost / map cost | trie cost / array cost | trie cost / btree cost |
| ---:       | ---:      | ---:                 | ---:                   | ---:                   |
| 1024       | 1         | 552.9 %              | 146.0 %                | 65.1 %                 |
| 1024       | 10        | 92.6 %               | 93.7 %                 | 35.8 %                 |
| 1024       | 100       | 139.6 %              | 97.2 %                 | 33.7 %                 |
| 1024       | 1000      | 195.6 %              | 104.4 %                | 30.5 %                 |
| 512        | 1000      | 266.4 %              | 111.6 %                | 33.7 %                 |
| 256        | 1000      | 355.9 %              | 125.5 %                | 37.6 %                 |

### test search lt, eq, gt key of trie search

| key length | key count | exsisted key cost (ns) | not exsisted key cost (ns) |
| ---:       | ---:      | ---:                   | ---:                       |
| 1024       | 1         | 119.000                | 106.667                    |
| 1024       | 10        | 182.333                | 179.000                    |
| 1024       | 100       | 246.000                | 218.000                    |
| 1024       | 1000      | 260.667                | 238.000                    |
| 512        | 1000      | 262.000                | 238.333                    |
| 256        | 1000      | 248.667                | 207.667                    |

