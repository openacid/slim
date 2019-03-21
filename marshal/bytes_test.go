package marshal_test

import (
	"reflect"
	"testing"

	"github.com/openacid/slim/marshal"
)

func TestBytes(t *testing.T) {

	var x marshal.Marshaler = marshal.Bytes{}
	_ = x

	cases := []struct {
		input []byte
		want  int
	}{
		{[]byte(""), 0},
		{[]byte("a"), 1},
		{[]byte("abc"), 3},
	}

	for i, c := range cases {
		m := marshal.Bytes{Size: c.want}
		l := m.GetSize(c.input)
		if l != c.want {
			t.Fatalf("%d-th: GetSize: input: %v; want: %v; actual: %v",
				i+1, c.input, c.want, l)
		}

		rst := m.Marshal(c.input)
		if len(rst) != c.want {
			t.Fatalf("%d-th: marshaled len: input: %v; want: %v; actual: %v",
				i+1, c.input, c.want, len(rst))
		}

		l = m.GetMarshaledSize(rst)
		if l != c.want {
			t.Fatalf("%d-th: marshaled size: input: %v; want: %v; actual: %v",
				i+1, c.input, c.want, l)
		}

		n, s := m.Unmarshal(rst)
		if c.want != n {
			t.Fatalf("%d-th: unmarshaled size: input: %v; want: %v; actual: %v",
				i+1, c.input, c.want, n)
		}
		if !reflect.DeepEqual(c.input, s) {
			t.Fatalf("%d-th: unmarshal: input: %v; want: %v; actual: %v",
				i+1, c.input, c.input, s)
		}

		if len(rst) > 0 {
			rst[0] = 'x'
			if s.([]byte)[0] != 'x' {
				t.Fatalf("should be not be copied.")
			}
		}

	}
}
