
### test search exsisted key between trie, map, array and btree:

| key length | key count | trie cost (ns) | map cost (ns) | array cost (ns) | btree cost (ns) |
| ---:       | ---:      | ---:           | ---:          | ---:            | ---:            |
| 1024       | 1         | 233.333        | 6.177         | 19.500          | 39.467          |
| 1024       | 10        | 270.000        | 63.267        | 55.367          | 100.567         |
| 1024       | 100       | 357.333        | 60.367        | 108.167         | 200.333         |
| 1024       | 1000      | 548.667        | 68.933        | 173.333         | 304.333         |
| 512        | 1000      | 458.000        | 36.600        | 148.667         | 380.000         |
| 256        | 1000      | 419.333        | 37.133        | 152.000         | 335.000         |

### test search not exsisted key between trie, map, array and btree:

| key length | key count | trie cost (ns) | map cost (ns) | array cost (ns) | btree cost (ns) |
| ---:       | ---:      | ---:           | ---:          | ---:            | ---:            |
| 1024       | 1         | 120.333        | 14.700        | 40.999          | 86.933          |
| 1024       | 10        | 129.000        | 71.567        | 71.600          | 167.333         |
| 1024       | 100       | 235.667        | 72.800        | 98.267          | 292.667         |
| 1024       | 1000      | 436.667        | 78.267        | 142.000         | 397.000         |
| 512        | 1000      | 334.333        | 49.867        | 134.333         | 406.000         |
| 256        | 1000      | 305.333        | 43.733        | 115.667         | 363.333         |

### ratio of search cost between trie, map, array and btree:

| key length | key count | trie cost / map cost | trie cost / array cost | trie cost / btree cost |
| ---:       | ---:      | ---:                 | ---:                   | ---:                   |
| 1024       | 1         | 37.774               | 11.966                 | 5.912                  |
| 1024       | 10        | 4.267                | 4.876                  | 2.685                  |
| 1024       | 100       | 5.919                | 3.303                  | 1.784                  |
| 1024       | 1000      | 7.959                | 3.165                  | 1.803                  |
| 512        | 1000      | 12.514               | 3.081                  | 1.205                  |
| 256        | 1000      | 11.293               | 2.759                  | 1.252                  |

### ratio of search not exsisted key cost between trie, map, array and btree

| key length | key count | trie cost / map cost | trie cost / array cost | trie cost / btree cost |
| ---:       | ---:      | ---:                 | ---:                   | ---:                   |
| 1024       | 1         | 8.186                | 2.935                  | 1.384                  |
| 1024       | 10        | 1.803                | 1.802                  | 0.771                  |
| 1024       | 100       | 3.237                | 2.398                  | 0.805                  |
| 1024       | 1000      | 5.579                | 3.075                  | 1.099                  |
| 512        | 1000      | 6.704                | 2.489                  | 0.823                  |
| 256        | 1000      | 6.982                | 2.639                  | 0.840                  |

### test search lt, eq, gt key of trie search

| key length | key count | exsisted  key cost (ns) | not exsisted key cost (ns) |
| ---:       | ---:      | ---:                    | ---:                       |
| 1024       | 1         | 272.333                 | 246.333                    |
| 1024       | 10        | 552.000                 | 495.000                    |
| 1024       | 100       | 823.333                 | 723.333                    |
| 1024       | 1000      | 936.667                 | 762.000                    |
| 512        | 1000      | 866.000                 | 700.667                    |
| 256        | 1000      | 844.333                 | 527.333                    |
