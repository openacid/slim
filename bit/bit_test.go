package bit

import "testing"

func TestPopCnt64Before(t *testing.T) {
	var cases = []struct {
		n    uint64
		iBit uint32
		expN uint32
	}{
		{0, 0, 0},   // 0b00
		{1, 0, 0},   // 0b01
		{2, 1, 0},   // 0b10
		{3, 1, 1},   // 0b11
		{3, 2, 2},   // 0b11
		{3, 4, 2},   // 0b11
		{32, 10, 1}, // 0b10000
	}
	for _, c := range cases {
		n, iBit, expN := c.n, c.iBit, c.expN
		actN := PopCnt64Before(n, iBit)
		if actN != expN {
			t.Fatalf("failed, case:%+v, actN:%d", c, actN)
		}
	}
}

func BenchmarkPopCnt64Before(b *testing.B) {

	var n uint64 = 12334567890

	for i := 0; i < b.N; i++ {
		PopCnt64Before(n+uint64(i), uint32(i)%64)
	}
}
