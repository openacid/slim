package typehelper

import (
	"reflect"
	"testing"

	"github.com/openacid/slim/benchhelper"
)

func TestToSlice(t *testing.T) {

	benchhelper.WantPanic(t, func() { ToSlice(1) }, "int")
	benchhelper.WantPanic(t, func() { ToSlice(nil) }, "nil")

	cases := []struct {
		input interface{}
		want  []interface{}
	}{{
		[]int{1, 2, 3},
		[]interface{}{1, 2, 3},
	}}

	for i, c := range cases {
		rst := ToSlice(c.input)
		if !reflect.DeepEqual(c.want, rst) {
			t.Fatalf("%d-th: input: %v; want: %v; actual: %v",
				i+1, c.input, c.want, rst)
		}
	}
}
