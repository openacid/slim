package trie

import "github.com/openacid/low/bitmap"

func newBM(indexes []int32, capa int32, opts ...string) *Bitmap {
	bb := &Bitmap{
		Words: bitmap.Of(indexes, capa),
	}
	bb.indexit(opts...)
	return bb
}

func (b *Bitmap) indexit(opts ...string) {
	for _, opt := range opts {
		switch opt {
		case "r64":
			b.RankIndex = bitmap.IndexRank64(b.Words)
		case "r128":
			b.RankIndex = bitmap.IndexRank128(b.Words)
		case "s32":
			// select32 also requires rank index to locate a bit
			b.SelectIndex, b.RankIndex = bitmap.IndexSelect32R64(b.Words)
		default:
			panic("unknown " + opt)
		}
	}
}
