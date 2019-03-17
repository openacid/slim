package bits_test

import (
	gobits "math/bits"
	"testing"

	"github.com/openacid/slim/bits"
)

func TestPopCnt64Before(t *testing.T) {
	var cases = []struct {
		n    uint64
		iBit uint32
		expN uint32
	}{
		{0, 0, 0},                    // 0b00
		{1, 0, 0},                    // 0b01
		{2, 1, 0},                    // 0b10
		{3, 1, 1},                    // 0b11
		{3, 2, 2},                    // 0b11
		{3, 4, 2},                    // 0b11
		{32, 10, 1},                  // 0b10000
		{0xffffffffffffffff, 0, 0},   // 0b11....
		{0xffffffffffffffff, 1, 1},   // 0b11....
		{0xffffffffffffffff, 2, 2},   // 0b11....
		{0xffffffffffffffff, 63, 63}, // 0b11....
		{0xffffffffffffffff, 64, 64}, // 0b11....
		{0xffffffffffffffff, 65, 64}, // 0b11....
		{0xffffffffffffffff, 66, 64}, // 0b11....
	}
	for _, c := range cases {
		n, iBit, expN := c.n, c.iBit, c.expN
		actN := bits.PopCnt64Before(n, iBit)
		if actN != expN {
			t.Fatalf("failed, case:%+v, actN:%d", c, actN)
		}
	}
}

func TestPopCnt32Before(t *testing.T) {
	var cases = []struct {
		n    uint32
		iBit uint32
		expN int
	}{
		{0, 0, 0},            // 0b00
		{1, 0, 0},            // 0b01
		{2, 1, 0},            // 0b10
		{3, 1, 1},            // 0b11
		{3, 2, 2},            // 0b11
		{3, 4, 2},            // 0b11
		{32, 10, 1},          // 0b10000
		{0xffffffff, 0, 0},   // 0b11....
		{0xffffffff, 1, 1},   // 0b11....
		{0xffffffff, 2, 2},   // 0b11....
		{0xffffffff, 31, 31}, // 0b11....
		{0xffffffff, 32, 32}, // 0b11....
		{0xffffffff, 33, 32}, // 0b11....
		{0xffffffff, 33, 32}, // 0b11....
	}
	for _, c := range cases {
		n, iBit, expN := c.n, c.iBit, c.expN
		actN := bits.PopCnt32Before(n, iBit)
		if actN != expN {
			t.Fatalf("failed, case:%+v, actN:%d", c, actN)
		}
	}
}

func TestPopCnt16Before(t *testing.T) {
	var cases = []struct {
		n    uint16
		iBit uint32
		expN int
	}{
		{0, 0, 0},        // 0b00
		{1, 0, 0},        // 0b01
		{2, 1, 0},        // 0b10
		{3, 1, 1},        // 0b11
		{3, 2, 2},        // 0b11
		{3, 4, 2},        // 0b11
		{32, 10, 1},      // 0b10000
		{0xffff, 0, 0},   // 0b11....
		{0xffff, 1, 1},   // 0b11....
		{0xffff, 2, 2},   // 0b11....
		{0xffff, 15, 15}, // 0b11....
		{0xffff, 16, 16}, // 0b11....
		{0xffff, 17, 16}, // 0b11....
	}
	for _, c := range cases {
		n, iBit, expN := c.n, c.iBit, c.expN
		actN := bits.PopCnt16Before(n, iBit)
		if actN != expN {
			t.Fatalf("failed, case:%+v, actN:%d", c, actN)
		}
	}
}

func TestPopCnt8Before(t *testing.T) {
	var cases = []struct {
		n    uint8
		iBit uint32
		expN int
	}{
		{0, 0, 0},    // 0b00
		{1, 0, 0},    // 0b01
		{2, 1, 0},    // 0b10
		{3, 1, 1},    // 0b11
		{3, 2, 2},    // 0b11
		{3, 4, 2},    // 0b11
		{0xff, 0, 0}, // 0b11....
		{0xff, 1, 1}, // 0b11....
		{0xff, 2, 2}, // 0b11....
		{0xff, 7, 7}, // 0b11....
		{0xff, 8, 8}, // 0b11....
		{0xff, 9, 8}, // 0b11....
	}
	for _, c := range cases {
		n, iBit, expN := c.n, c.iBit, c.expN
		actN := bits.PopCnt8Before(n, iBit)
		if actN != expN {
			t.Fatalf("failed, case:%+v, actN:%d", c, actN)
		}
	}
}

func BenchmarkPopCnt64Before(b *testing.B) {

	var n uint64 = 12334567890

	for i := 0; i < b.N; i++ {
		bits.PopCnt64Before(n+uint64(i), uint32(i)%64)
	}
}

func BenchmarkPopCnt32Before(b *testing.B) {

	var n uint32 = 123345678

	for i := 0; i < b.N; i++ {
		bits.PopCnt32Before(n+uint32(i), uint32(i)%32)
	}
}

func BenchmarkPopCnt16(b *testing.B) {

	for i := 0; i < b.N; i++ {
		bits.PopCnt16Before(uint16(i), uint32(i)%16)
	}
}

func BenchmarkPopCnt8(b *testing.B) {

	for i := 0; i < b.N; i++ {
		bits.PopCnt8Before(uint8(i), uint32(i)%8)
	}
}

func BenchmarkGoPopCnt64Before(b *testing.B) {

	var n uint64 = 12334567890

	for i := 0; i < b.N; i++ {
		gobits.OnesCount64(n + uint64(i))
	}
}

func BenchmarkGoPopCnt32Before(b *testing.B) {

	var n uint32 = 123345678

	for i := 0; i < b.N; i++ {
		gobits.OnesCount32(n + uint32(i))
	}
}

func BenchmarkGoPopCnt16(b *testing.B) {

	for i := 0; i < b.N; i++ {
		gobits.OnesCount16(uint16(i))
	}
}

func BenchmarkGoPopCnt8(b *testing.B) {

	for i := 0; i < b.N; i++ {
		gobits.OnesCount8(uint8(i))
	}
}
