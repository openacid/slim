## Memory overhead

-   Random string, fixed length, default mode, no label is store if possible:

    **Bits/key**: memory or disk-space in bits a key consumed in average.
    It does not change when key-length(`k`) becomes larger!

    ![](trie/report/mem_usage.jpg)


-   1 million var-length string, 10 to 20 byte in different mode SlimTrie:

    | -          | size  | gzip-size |
    | :--        | --:   | --:       |
    | Original   | 15.0M | 14.0M     |
    | Complete   | 14.0M | 10.0M     |
    | InnerLabel |  1.3M |  0.9M     |
    | NoLabel    |  1.3M |  0.8M     |

    Raw string list and serialized slim is stored in:
    https://github.com/openacid/testkeys/tree/master/assets

    -   Original: raw string lines in a text file.

    -   Complete: `NewSlimTrie(..., Opt{Complete:Bool(true)})`: lossless SlimTrie,
        stores complete info of every string. This mode provides accurate query.

    -   InnerLabel: `NewSlimTrie(..., Opt{InnerPrefix:Bool(true)})` SlimTrie stores
        only label strings of inner nodes(but not label to a leaf). There is false positive in this mode.

    -   NoLabel: No label info is stored. False positive rate is higher.


## Performance

Time(in nano second) spent on a `Get()` with golang-map, SlimTrie, array and [btree][] by google.

- **3.3 times faster** than the [btree][].
- **2.3 times faster** than binary search.

![](trie/report/bench_msab_present_zipf.jpg)


Time(in nano second) spent on a `Get()` with different key count(`n`) and key length(`k`):

![](trie/report/bench_get_present_zipf.jpg)


## False Positive Rate

![](trie/report/fpr_get.jpg)

> Bloom filter requires about 9 bits/key to archieve less than 1% FPR.

See: [trie/report/](trie/report/)
