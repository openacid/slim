package trie

import (
	"fmt"
	"testing"

	proto "github.com/golang/protobuf/proto"
	"github.com/openacid/low/mathext/zipf"
	"github.com/openacid/low/size"
	"github.com/openacid/slim/encode"
)

var Outputxxx int32

func BenchmarkSlimTrie_GetXXX(b *testing.B) {
	benchBigKeySet(b, func(b *testing.B, typ string, keys []string) {
		n := len(keys)
		values := makeI32s(n)

		for _, c := range []struct {
			name string
			opt  Opt
		}{
			{name: "complete", opt: Opt{Complete: Bool(true)}},
			{name: "innerPrefix", opt: Opt{InnerPrefix: Bool(true)}},
		} {
			b.Run(c.name, func(b *testing.B) {

				for _, maxLevel := range []int32{3, 4, 5, 6, 7, 8, 9, 15} {

					opt := c.opt
					opt.MaxLevel = I32(maxLevel)

					st, _ := NewSlimTrie(encode.I32{}, keys, values, opt)

					protosz := proto.Size(st)
					sz := size.Of(st)
					n := len(keys)
					avg := sz/n - 4 // exclude value size
					fmt.Println("protosize:", protosz/1024, "kb")
					fmt.Println("total size: ", sz/1024, "kb; sz/n:", avg)
					// fmt.Println(size.Stat(st, 10, 3))

					// mx := uint32(0)
					// ofs := st.inner.Clustered.Offsets
					// pp := 8
					// for i := pp; i < len(ofs); i++ {
					//     diff := ofs[i] - ofs[i-pp]
					//     if diff > mx {
					//         mx = diff
					//     }

					// }
					// fmt.Println("max offset diff(32):", mx)

					// panic("foo")

					b.Run(fmt.Sprintf("Lvl:%d", maxLevel),
						func(b *testing.B) {
							subBenGetID(b, st, keys)
							subBenGet(b, st, keys)
						})
				}
			})
		}
	})
}

func subBenGetID(b *testing.B, st *SlimTrie, keys []string) {
	n := len(keys)
	b.Run("GetID", func(b *testing.B) {
		load := zipf.Accesses(2, 1.5, n, b.N, nil)
		var id int32

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			id += st.GetID(keys[load[i]])
		}
		Outputxxx = id
	})
}

func subBenGet(b *testing.B, st *SlimTrie, keys []string) {
	n := len(keys)
	b.Run("Get", func(b *testing.B) {
		load := zipf.Accesses(2, 1.5, n, b.N, nil)
		var id int32

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			v, _ := st.Get(keys[load[i]])
			id += v.(int32)
		}
		Outputxxx = id
	})
}

func BenchmarkSlimTrie_withPrefixContent_GetID_20k_vlen10(b *testing.B) {

	keys := getKeys("20kvl10")
	values := makeI32s(len(keys))
	st, _ := NewSlimTrie(encode.I32{}, keys, values,
		Opt{InnerPrefix: Bool(true)},
	)

	var id int32

	b.ResetTimer()

	i := b.N
	for {
		for _, k := range keys {
			id += st.GetID(k)

			i--
			if i == 0 {
				Outputxxx = id
				return
			}
		}
	}
}
