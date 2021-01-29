package trie

import (
	"fmt"
	"testing"

	"github.com/openacid/low/mathext/zipf"
	"github.com/openacid/low/size"
)

var OutputMap int

func Benchmark_map(b *testing.B) {

	benchBigKeySet(b, func(b *testing.B, typ string, keys []string) {

		n := len(keys)
		values := makeI32s(len(keys))

		m := make(map[string]int32, n)
		for i := 0; i < len(keys); i++ {
			m[keys[i]] = values[i]
		}

		sz := size.Of(m)
		fmt.Println(sz/1024, sz/len(keys))
		fmt.Println(size.Stat(m, 10, 10))

		accesses := zipf.Accesses(2, 1.5, len(keys), b.N, nil)

		b.ResetTimer()

		var id int32
		for i := 0; i < b.N; i++ {
			idx := accesses[i]
			id += m[keys[idx]]

		}
		OutputMap = int(id)
	})
}
