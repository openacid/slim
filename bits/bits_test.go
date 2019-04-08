package bits_test

import (
	gobits "math/bits"
	"testing"

	. "github.com/openacid/slim/bits"
)

// Exported (global) variable serving as input for some
// of the benchmarks to ensure side-effect free calls
// are not optimized away.
var Input uint64 = 0x03f79d71b4ca8b09

// Exported (global) variable to store function results
// during benchmarking to ensure side-effect free calls
// are not optimized away.
var Output int

func TestOnesCount64AndBefore(t *testing.T) {
	var cases = []struct {
		n    uint64
		iBit uint
		want int
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
		n, iBit, want := c.n, c.iBit, c.want

		rst := OnesCount64Before(n, iBit)
		if rst != want {
			t.Fatalf("failed, case:%+v, rst:%d", c, rst)
		}

		// uint test
		if UintSize == 64 {
			rst := OnesCountBefore(uint(n), iBit)
			if rst != want {
				t.Fatalf("failed, case:%+v, rst:%d", c, rst)
			}
		}
	}
}

func TestOnesCount32AndBefore(t *testing.T) {
	var cases = []struct {
		n    uint32
		iBit uint
		want int
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
		n, iBit, want := c.n, c.iBit, c.want

		rst := OnesCount32Before(n, iBit)
		if rst != want {
			t.Fatalf("failed, case:%+v, rst:%d", c, rst)
		}

		// uint test
		if UintSize == 32 {
			rst := OnesCountBefore(uint(n), iBit)
			if rst != want {
				t.Fatalf("failed, case:%+v, rst:%d", c, rst)
			}
		}
	}
}

func TestOnesCount16AndBefore(t *testing.T) {
	var cases = []struct {
		n    uint16
		iBit uint
		want int
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
		n, iBit, want := c.n, c.iBit, c.want

		rst := OnesCount16Before(n, iBit)
		if rst != want {
			t.Fatalf("failed, case:%+v, rst:%d", c, rst)
		}
	}
}

func TestOnesCount8AndBefore(t *testing.T) {
	var cases = []struct {
		n    uint8
		iBit uint
		want int
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
		n, iBit, want := c.n, c.iBit, c.want

		rst := OnesCount8Before(n, iBit)
		if rst != want {
			t.Fatalf("failed, case:%+v, rst:%d", c, rst)
		}
	}
}

func BenchmarkOnesCount64Before(b *testing.B) {

	var s int
	for i := 0; i < b.N; i++ {
		s += OnesCount64Before(Input, 15)
	}

	Output = s
}

func BenchmarkOnesCount32Before(b *testing.B) {

	var s int
	for i := 0; i < b.N; i++ {
		s += OnesCount32Before(uint32(Input), 15)
	}
	Output = s
}

func BenchmarkOnesCount16Before(b *testing.B) {

	var s int
	for i := 0; i < b.N; i++ {
		s += OnesCount16Before(uint16(Input), 15)
	}
	Output = s
}

func BenchmarkOnesCount8Before(b *testing.B) {

	var s int
	for i := 0; i < b.N; i++ {
		s += OnesCount8Before(uint8(Input), 6)
	}
	Output = s
}

func BenchmarkOnesCountBefore(b *testing.B) {

	var s int
	for i := 0; i < b.N; i++ {
		s += OnesCountBefore(uint(Input), 6)
	}
	Output = s
}

func BenchmarkGoOnesCount8(b *testing.B) {
	var s int
	for i := 0; i < b.N; i++ {
		s += gobits.OnesCount8(uint8(Input))
	}
	Output = s
}

func BenchmarkGoOnesCount16(b *testing.B) {
	var s int
	for i := 0; i < b.N; i++ {
		s += gobits.OnesCount16(uint16(Input))
	}
	Output = s
}

func BenchmarkGoOnesCount32(b *testing.B) {
	var s int
	for i := 0; i < b.N; i++ {
		s += gobits.OnesCount32(uint32(Input))
	}
	Output = s
}

func BenchmarkGoOnesCount64(b *testing.B) {
	var s int
	for i := 0; i < b.N; i++ {
		s += gobits.OnesCount64(Input)
	}
	Output = s
}
