package typehelper

import (
	"reflect"
	"testing"
)

func TestToSlice(t *testing.T) {

	cases := []struct {
		input  interface{}
		want   []interface{}
		wantok bool
	}{{
		[]int{1, 2, 3},
		[]interface{}{1, 2, 3},
		true,
	}, {
		int(1),
		nil,
		false,
	}, {
		nil,
		nil,
		false,
	}}

	for i, c := range cases {
		rst, ok := ToSlice(c.input)
		if !reflect.DeepEqual(c.want, rst) || c.wantok != ok {
			t.Fatalf("%d-th: input: %v; want: %v %v; actual: %v %v",
				i+1, c.input, c.want, c.wantok, rst, ok)
		}
	}
}
