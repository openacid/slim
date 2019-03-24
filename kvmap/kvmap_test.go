package kvmap_test

import (
	crand "crypto/rand"
	"fmt"
	"sort"
	"testing"

	"github.com/openacid/slim/kvmap"
)

func RandomString(leng int) string {
	return string(RandomBytes(leng))
}

func RandomBytes(leng int) []byte {
	bs := make([]byte, leng)
	n, err := crand.Read(bs)
	if err != nil {
		panic(err)
	}
	if n != leng {
		panic("not read enough")
	}
	return bs
}

func TestKVMapNew(t *testing.T) {
	n := 100
	keylen := 32
	ks := make([]string, n)
	for i := 0; i < n; i++ {
		ks[i] = RandomString(keylen)
	}

	sort.Strings(ks)

	items := make([]kvmap.Item, n)
	for i := 0; i < n; i++ {
		items[i] = kvmap.Item{Key: ks[i], Val: ks[i]}
	}

	kv, err := kvmap.NewKVMap(items)
	if err != nil {
		t.Fatalf("expected no error but: %+v", err)
	}

	k := items[0].Key
	fmt.Println(kv.Get(k))
}

func makeit(n, keylen int) (*kvmap.KVMap, []string) {
	ks := make([]string, n)
	for i := 0; i < n; i++ {
		ks[i] = RandomString(keylen)
	}

	sort.Strings(ks)

	items := make([]kvmap.Item, n)
	for i := 0; i < n; i++ {
		items[i] = kvmap.Item{Key: ks[i], Val: ks[i]}
	}

	kv, err := kvmap.NewKVMap(items)
	if err != nil {
		panic(err)
	}
	return kv, ks
}

func BenchmarkKVMapNew(b *testing.B) {

	n := 5000
	keylen := 32

	kv, ks := makeit(n, keylen)

	k := ks[50]

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		kv.Get(k)
	}

}
