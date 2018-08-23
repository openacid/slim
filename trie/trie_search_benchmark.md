
### test search exsisted key between trie, map, array and btree:

| key length | key count | trie cost (ns) | map cost (ns) | array cost (ns) | btree cost (ns) |
| ---:       | ---:      | ---:           | ---:          | ---:            | ---:            |
| 1024       | 1         | 2860.333       | 6.177         | 19.500          | 39.467          |
| 1024       | 10        | 3946.000       | 63.267        | 55.367          | 100.567         |
| 1024       | 100       | 4519.000       | 60.367        | 108.167         | 200.333         |
| 1024       | 1000      | 5270.000       | 68.933        | 173.333         | 304.333         |
| 512        | 1000      | 2761.333       | 36.600        | 148.667         | 380.000         |
| 256        | 1000      | 2018.000       | 37.133        | 152.000         | 335.000         |

### test search not exsisted key between trie, map, array and btree:

| key length | key count | trie cost (ns) | map cost (ns) | array cost (ns) | btree cost (ns) |
| ---:       | ---:      | ---:           | ---:          | ---:            | ---:            |
| 1024       | 1         | 2814.667       | 14.700        | 40.999          | 86.933          |
| 1024       | 10        | 3137.667       | 71.567        | 71.600          | 167.333         |
| 1024       | 100       | 4051.000       | 72.800        | 98.267          | 292.667         |
| 1024       | 1000      | 4593.667       | 78.267        | 142.000         | 397.000         |
| 512        | 1000      | 2512.333       | 49.867        | 134.333         | 406.000         |
| 256        | 1000      | 1854.333       | 43.733        | 115.667         | 363.333         |

### ratio of search cost between trie, map, array and btree:

| key length | key count | trie cost / map cost | trie cost / array cost | trie cost / btree cost |
| ---:       | ---:      | ---:                 | ---:                   | ---:                   |
| 1024       | 1         | 463.061              | 146.684                | 72.474                 |
| 1024       | 10        | 62.371               | 71.269                 | 39.238                 |
| 1024       | 100       | 74.858               | 41.778                 | 22.557                 |
| 1024       | 1000      | 76.451               | 30.404                 | 17.316                 |
| 512        | 1000      | 75.446               | 18.574                 | 7.267                  |
| 256        | 1000      | 54.345               | 13.276                 | 6.024                  |

### ratio of search not exsisted key cost between trie, map, array and btree

| key length | key count | trie cost / map cost | trie cost / array cost | trie cost / btree cost |
| ---:       | ---:      | ---:                 | ---:                   | ---:                   |
| 1024       | 1         | 191.474              | 68.652                 | 32.377                 |
| 1024       | 10        | 43.842               | 43.822                 | 18.751                 |
| 1024       | 100       | 55.646               | 41.224                 | 13.842                 |
| 1024       | 1000      | 58.692               | 32.349                 | 11.571                 |
| 512        | 1000      | 50.381               | 18.702                 | 6.188                  |
| 256        | 1000      | 42.401               | 16.032                 | 5.104                  |
