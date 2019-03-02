package array

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func makeRandArray(cnt uint32) (idx uint32, indexes []uint32, keysMap map[uint32]bool, ar *Array32Index, err error) {
	arr := &Array32Index{}

	indexes = []uint32{}
	rnd := rand.New(rand.NewSource(time.Now().Unix()))

	keysMap = map[uint32]bool{}
	num, idx := uint32(0), uint32(0)
	for ; num < cnt; idx++ {
		if rnd.Intn(2) == 1 {
			indexes = append(indexes, idx)
			num++
			keysMap[idx] = true
		}
	}

	err = arr.InitIndexBitmap(indexes)
	return idx, indexes, keysMap, arr, err
}

func TestHasAndGetEltIndex(t *testing.T) {
	maxIndex, indexes, keysMap, arr, err := makeRandArray(1024)
	if err != nil {
		t.Fatalf("expect no err but: %s", err)
	}

	for i := uint32(0); i < maxIndex; i++ {
		if _, ok := keysMap[i]; ok {
			if !arr.Has(i) {
				t.Fatalf("expect has but not: %d", i)
			}
			eltIndex, found := arr.GetEltIndex(i)
			if !found {
				t.Fatalf("should found but not: %d", i)
			}
			if indexes[eltIndex] != i {
				t.Fatalf("i=%d should be at %d", i, eltIndex)
			}
		} else {
			if arr.Has(i) {
				t.Fatalf("expect not has but has: %d", i)
			}
			_, found := arr.GetEltIndex(i)
			if found {
				t.Fatalf("should not found but found: %d", i)
			}
		}
	}
}

func BenchmarkHasAndGetEltIndex(b *testing.B) {

	var name string
	runs := []struct{ cnt uint32 }{
		{32},
		{256},
		{1024},
		{10240},
		{102400},
	}
	for _, r := range runs {
		maxIndex, _, _, arr, _ := makeRandArray(r.cnt)

		name = fmt.Sprintf("Has-%d", r.cnt)
		b.Run(name, func(b *testing.B) {

			j := uint32(0)
			for i := 0; i < b.N; i++ {
				arr.Has(j)
				j = (j + 1) % (maxIndex)
			}
		})

		name = fmt.Sprintf("GetEltIndex-%d", r.cnt)
		b.Run(name, func(b *testing.B) {

			j := uint32(0)
			for i := 0; i < b.N; i++ {
				arr.GetEltIndex(j)
				j = (j + 1) % (maxIndex)
			}
		})

	}
}
