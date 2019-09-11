package trie

import (
	"fmt"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/openacid/slim/encode"
	"github.com/openacid/slimcompatible/testkeys"
)

func MakeMarshaledData(fn string, keys []string) {

	if len(keys) == 0 {
		keys = testkeys.Keys["50kl10"]
	}

	n := len(keys)
	values := make([]int32, n)
	for i := 0; i < n; i++ {
		values[i] = int32(i)
	}

	st, err := NewSlimTrie(encode.I32{}, keys, values)
	if err != nil {
		panic(err)
	}

	b, err := proto.Marshal(st)
	if err != nil {
		panic(err)
	}

	fn = fmt.Sprintf(fn, st.GetVersion())
	f := newFile(fn)
	defer f.Close()

	_, err = f.Write(b)
	if err != nil {
		panic(err)
	}
}

func newFile(fn string) *os.File {
	f, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	err = f.Truncate(0)
	if err != nil {
		panic(err)
	}
	return f
}
