package index_test

import (
	"fmt"
	"strings"

	"github.com/openacid/slim/index"
)

type RangeData string

func (d RangeData) Read(offset int64, key string) (string, bool) {
	for i := 0; i < 4; i++ {
		if int(offset) >= len(d) {
			break
		}

		kv := strings.Split(string(d)[offset:], ",")[0:2]
		if kv[0] == key {
			return kv[1], true
		}
		offset += int64(len(kv[0]) + len(kv[1]) + 2)

	}
	return "", false
}

func Example_indexRanges() {

	// Index ranges instead of keys:
	// In this example at most 4 keys shares one index item.

	data := RangeData("Aaron,1,Agatha,1,Al,2,Albert,3,Alexander,5,Alison,8")

	// keyOffsets is a prebuilt index that stores range start, range end and its offset.
	keyOffsets := []index.OffsetIndexItem{
		// Aaron  +--> 0
		// Agatha |
		// Al     |
		// Albert |

		// Alexander +--> 31
		// Alison    |

		{Key: "Aaron", Offset: 0},
		{Key: "Albert", Offset: 0},

		{Key: "Alexander", Offset: 31},
		{Key: "Alison", Offset: 31},
	}

	st, err := index.NewSlimIndex(keyOffsets, data)
	if err != nil {
		panic(err)
	}

	v, found := st.RangeGet("Aaron")
	fmt.Printf("key: %q\n  found: %t\n  value: %q\n", "Aaron", found, v)

	v, found = st.RangeGet("Al")
	fmt.Printf("key: %q\n  found: %t\n  value: %q\n", "Al", found, v)

	v, found = st.RangeGet("foo")
	fmt.Printf("key: %q\n  found: %t\n  value: %q\n", "foo", found, v)

	// Output:
	// key: "Aaron"
	//   found: true
	//   value: "1"
	// key: "Al"
	//   found: true
	//   value: "2"
	// key: "foo"
	//   found: false
	//   value: ""
}
