package trie_test

import (
	"strings"
	"testing"

	"github.com/openacid/slim/trie"
)

type testIndexData string

func (d testIndexData) Read(offset int64, key string) (string, bool) {
	kv := strings.Split(string(d)[offset:], ",")[0:2]
	if kv[0] == key {
		return kv[1], true
	}
	return "", false
}

func TestSlimIndex(t *testing.T) {

	data := testIndexData("Aaron,1,Agatha,1,Al,2,Albert,3,Alexander,5,Alison,8")

	index := []trie.OffsetIndexItem{
		{Key: "Aaron", Offset: 0},
		{Key: "Agatha", Offset: 8},
		{Key: "Al", Offset: 17},
		{Key: "Albert", Offset: 22},
		{Key: "Alexander", Offset: 31},
		{Key: "Alison", Offset: 43},
	}

	st, err := trie.NewSlimIndex(index, data)
	if err != nil {
		t.Fatalf("expect no error but: %s", err)
	}

	cases := []struct {
		input     string
		want      string
		wantfound bool
	}{
		{"Aaron", "1", true},
		{"Agatha", "1", true},
		{"Al", "2", true},
		{"Albert", "3", true},
		{"Alexander", "5", true},
		{"Alison", "8", true},
		{"foo", "", false},
		{"Alexande", "", false},
		{"Alexander0", "", false},
		{"alexander", "", false},
	}

	for i, c := range cases {
		rst, found := st.Get2(c.input)
		if rst != c.want {
			t.Fatalf("%d-th: input: %v; want: %v; actual: %v",
				i+1, c.input, c.want, rst)
		}
		if found != c.wantfound {
			t.Fatalf("%d-th: input: %v; wantfound: %v; actual: %v",
				i+1, c.input, c.wantfound, found)
		}
	}

}
