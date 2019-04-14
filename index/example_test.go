package index_test

import (
	"fmt"
	"strings"

	"github.com/openacid/slim/index"
)

type Data string

func (d Data) Read(offset int64, key string) (string, bool) {
	kv := strings.Split(string(d)[offset:], ",")[0:2]
	if kv[0] == key {
		return kv[1], true
	}
	return "", false
}

func Example() {

	// Accelerate external data accessing (in memory or on disk) by indexing
	// them with a SlimTrie:

	// `data` is a sample of some unindexed data. In our example it is a comma
	// separated key value series.
	//
	// In order to let SlimTrie be able to read data, `data` should have
	// a `Read` method:
	//     Read(offset int64, key string) (string, bool)
	data := Data("Aaron,1,Agatha,1,Al,2,Albert,3,Alexander,5,Alison,8")

	// keyOffsets is a prebuilt index that stores key and its offset in data accordingly.
	keyOffsets := []index.OffsetIndexItem{
		{Key: "Aaron", Offset: 0},
		{Key: "Agatha", Offset: 8},
		{Key: "Al", Offset: 17},
		{Key: "Albert", Offset: 22},
		{Key: "Alexander", Offset: 31},
		{Key: "Alison", Offset: 43},
	}

	// `SlimIndex` is simply a container of SlimTrie and its data.
	st, err := index.NewSlimIndex(keyOffsets, data)
	if err != nil {
		fmt.Println(err)
	}

	// Lookup
	v, found := st.Get("Alison")
	fmt.Printf("key: %q\n  found: %t\n  value: %q\n", "Alison", found, v)

	v, found = st.Get("foo")
	fmt.Printf("key: %q\n  found: %t\n  value: %q\n", "foo", found, v)

	// Output:
	// key: "Alison"
	//   found: true
	//   value: "8"
	// key: "foo"
	//   found: false
	//   value: ""
}
