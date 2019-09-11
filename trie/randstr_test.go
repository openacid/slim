package trie

import (
	"math/rand"
	"sort"
)

var runes = []rune("~!@#$%^&*()_+`-=[]{};:<>?,./abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

// var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func randVStrings(n, minLen, maxLen int, ops ...interface{}) []string {

	var from []rune

	if len(ops) > 0 {
		from = ops[0].([]rune)
	}

	if from == nil {
		from = runes
	}

	rlen := len(from)

	mp := make(map[string]bool)

	for i := 0; i < n; i++ {
		l := rand.Intn((maxLen-minLen)+1) + minLen
		b := make([]rune, l)
		for j := 0; j < l; j++ {
			k := rand.Intn(rlen)

			b[j] = from[k]
		}
		s := string(b)
		if _, ok := mp[s]; ok {
			i--
		} else {
			mp[s] = true
		}
	}

	rst := make([]string, 0, n)
	for k := range mp {
		rst = append(rst, k)
	}

	sort.Strings(rst)
	return rst
}
