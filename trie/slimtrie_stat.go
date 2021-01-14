package trie

type Stat struct {
	LevelCnt int32
	Levels   []struct {
		Total, Inner, Leaf int32
	}
	KeyCnt  int32
	NodeCnt int32
}

// Stat() returns a struct `Stat` describing SlimTrie internal stats. E.g.:
//
//     &trie.Stat{
//         LevelCnt: 5,
//         Levels:   {
//             {},
//             {Total:1, Inner:1, Leaf:0},
//             {Total:4, Inner:3, Leaf:1},
//             {Total:8, Inner:6, Leaf:2},
//             {Total:14, Inner:6, Leaf:8},
//         },
//         KeyCnt:  8,
//         NodeCnt: 14,
//     }
//
// Since 0.5.12
func (st *SlimTrie) Stat() *Stat {

	ns := st.inner

	rst := &Stat{}

	// The 0-th level is a trivial level with {0, 0, 0}
	level_cnt := len(st.levels)
	rst.LevelCnt = int32(level_cnt)
	for i := 0; i < level_cnt; i++ {
		l := st.levels[i]
		rst.Levels = append(rst.Levels, struct{ Total, Inner, Leaf int32 }{
			l.total, l.inner, l.leaf,
		})
	}

	if ns.NodeTypeBM == nil {
		rst.KeyCnt = 0
	} else {
		rst.KeyCnt = st.levels[level_cnt-1].leaf
	}

	rst.NodeCnt = st.levels[level_cnt-1].total

	return rst
}
