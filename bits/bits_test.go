package bits_test

import (
	gobits "math/bits"
	"testing"

	"github.com/openacid/slim/bits"
)

func TestOnesCount64Before(t *testing.T) {
	var cases = []struct {
		n    uint64
		iBit uint
		expN int
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
		actN := bits.OnesCount64Before(n, iBit)
		if actN != expN {
			t.Fatalf("failed, case:%+v, actN:%d", c, actN)
		}
	}
}

func TestOnesCount32Before(t *testing.T) {
	var cases = []struct {
		n    uint32
		iBit uint
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
		actN := bits.OnesCount32Before(n, iBit)
		if actN != expN {
			t.Fatalf("failed, case:%+v, actN:%d", c, actN)
		}
	}
}

func TestOnesCount16Before(t *testing.T) {
	var cases = []struct {
		n    uint16
		iBit uint
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
		actN := bits.OnesCount16Before(n, iBit)
		if actN != expN {
			t.Fatalf("failed, case:%+v, actN:%d", c, actN)
		}
	}
}

func TestOnesCount8Before(t *testing.T) {
	var cases = []struct {
		n    uint8
		iBit uint
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
		actN := bits.OnesCount8Before(n, iBit)
		if actN != expN {
			t.Fatalf("failed, case:%+v, actN:%d", c, actN)
		}
	}
}

func BenchmarkOnesCount64Before(b *testing.B) {

	var n uint64 = 12334567890

	for i := 0; i < b.N; i++ {
		bits.OnesCount64Before(n+uint64(i), uint(i)%64)
	}
}

func BenchmarkOnesCount32Before(b *testing.B) {

	var n uint32 = 123345678

	for i := 0; i < b.N; i++ {
		bits.OnesCount32Before(n+uint32(i), uint(i)%32)
	}
}

func BenchmarkOnesCount16Before(b *testing.B) {

	for i := 0; i < b.N; i++ {
		bits.OnesCount16Before(uint16(i), uint(i)%16)
	}
}

func BenchmarkOnesCount8Before(b *testing.B) {

	for i := 0; i < b.N; i++ {
		bits.OnesCount8Before(uint8(i), uint(i)%8)
	}
}

func BenchmarkGoOnesCount64(b *testing.B) {

	var n uint64 = 12334567890

	for i := 0; i < b.N; i++ {
		gobits.OnesCount64(n + uint64(i))
	}
}

func BenchmarkGoOnesCount32(b *testing.B) {

	var n uint32 = 123345678

	for i := 0; i < b.N; i++ {
		gobits.OnesCount32(n + uint32(i))
	}
}

func BenchmarkGoOnesCount16(b *testing.B) {

	for i := 0; i < b.N; i++ {
		gobits.OnesCount16(uint16(i))
	}
}

func BenchmarkGoOnesCount8(b *testing.B) {

	for i := 0; i < b.N; i++ {
		gobits.OnesCount8(uint8(i))
	}
}
