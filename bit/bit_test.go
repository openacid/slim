package bit

import "testing"

func TestCnt1Before(t *testing.T) {
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
		actN := Cnt1Before(n, iBit)
		if actN != expN {
			t.Fatalf("failed, case:%+v, actN:%d", c, actN)
		}
	}
}
