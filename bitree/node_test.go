package bitree_test

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/kr/pretty"
	"github.com/openacid/slim/bitree"
)

func bin(b string) byte {
	i, err := strconv.ParseInt(b, 2, 64)
	if err != nil {
		panic(err)
	}
	return byte(i)
}

func TestNode64New(t *testing.T) {

	cases := []struct {
		positions    []int32
		significants []byte
		elts         [][]byte
		want         []byte
	}{
		{
			positions: []int32{100, 102, 108, 109, 110, 110 + 200},
			significants: []byte{
				0, 9, 33, 35,
			},
			elts: [][]byte{
				{0, 1},
				{1, 2},
				{2, 3},
				{3, 4},
			},
			want: []byte{
				// flag
				0,
				// positions
				100, 2, 6, 1, 1, 200, 1,
				// significants
				bin("00000001"),
				bin("00000010"),
				bin("00000000"),
				bin("00000000"),
				bin("00001010"),
				bin("00000000"),
				bin("00000000"),
				bin("00000000"),
				// elts
				0, 1, 1, 2, 2, 3, 3, 4,
			},
		},
	}

	for i, c := range cases {
		rst := bitree.New(c.positions, c.significants, c.elts)
		if !reflect.DeepEqual(c.want, rst) {
			fmt.Println(pretty.Diff(c.want, rst))
			t.Fatalf("%d-th: input: %#v %#v %#v; want: %#v; actual: %#v",
				i+1, c.positions, c.significants, c.elts, c.want, rst)
		}

		ps := bitree.Node64ReadPositions(rst)
		if !reflect.DeepEqual(c.positions, ps) {
			fmt.Println(pretty.Diff(c.positions, ps))
			t.Fatalf("%d-th expect: %v; but: %v",
				i+1, c.positions, ps)
		}

		val := bitree.Node64Get(rst, 0)
		if !reflect.DeepEqual(c.elts[0], val) {
			t.Fatalf("0-th should be found expect: %v; but: %v", c.elts[0], val)
		}
		val = bitree.Node64Get(rst, 1)
		if val != nil {
			t.Fatalf("should be nli at 1 expect: %v; but: %v", nil, val)
		}

	}
}
