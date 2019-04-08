// Package bits provides efficient bitwise operations on integer numbers.
//
// Performance note
//
// Benchmarks shows that one counting costs 1 ns/op.
package bits

import (
	gobits "math/bits"
)

// UintSize is the size of a uint in bits.
const UintSize = 32 << (^uint(0) >> 32 & 1) // 32 or 64

// OnesCount64Before counts the number of "1"(population count) in a uint64 before specified bit position `iBit`.
//
// "math/bits" implements OnesCount64 with "popcnt" instruction thus it is very
// fast.
//
// E.g.:
//		(Significant bits on left)
//
//		3 = ...011	OnesCount64Before(3, 0) == 0
//		          	OnesCount64Before(3, 1) == 1
//		          	OnesCount64Before(3, 2) == 2
//		          	OnesCount64Before(3, 3) == 2
//
// This algorithm has more introduction in:
// https://en.wikipedia.org/wiki/Hamming_weight#Efficient_implementation
func OnesCount64Before(n uint64, iBit uint) int {
	return gobits.OnesCount64(n & ((uint64(1) << iBit) - 1))
}

// OnesCount32Before counts the number of "1"(population count) in a uint32 before specified bit position `iBit`.
func OnesCount32Before(n uint32, iBit uint) int {
	return gobits.OnesCount32(n & ((uint32(1) << iBit) - 1))
}

// OnesCount16Before counts the number of "1"(population count) in a uint16 before specified bit position `iBit`.
func OnesCount16Before(n uint16, iBit uint) int {
	return gobits.OnesCount16(n & ((uint16(1) << iBit) - 1))
}

// OnesCount8Before counts the number of "1"(population count) in a uint8 before specified bit position `iBit`.
func OnesCount8Before(n uint8, iBit uint) int {
	return gobits.OnesCount8(n & ((uint8(1) << iBit) - 1))
}

// OnesCount8Before counts the number of "1"(population count) in a uint before specified bit position `iBit`.
func OnesCountBefore(n uint, iBit uint) int {
	if UintSize == 32 {
		return gobits.OnesCount32(uint32(n) & ((uint32(1) << iBit) - 1))
	}
	return gobits.OnesCount64(uint64(n) & ((uint64(1) << iBit) - 1))
}
