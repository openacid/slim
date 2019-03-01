// Package bit provides efficient bitwise operations on integer numbers.
package bit

const (
	m64_1 uint64 = 0x5555555555555555 // binary: 0101...
	m64_2 uint64 = 0x3333333333333333 // binary: 00110011...
	m64_4 uint64 = 0x0f0f0f0f0f0f0f0f // binary: 0000111100001111...
	h64_8 uint64 = 0x0101010101010101 // binary: 0000000100000001...

	m32_1 uint32 = 0x55555555 // binary: 0101...
	m32_2 uint32 = 0x33333333 // binary: 00110011...
	m32_4 uint32 = 0x0f0f0f0f // binary: 0000111100001111...
	h32_8 uint32 = 0x01010101 // binary: 0000000100000001...

	m16_1 uint16 = 0x5555 // binary: 0101...
	m16_2 uint16 = 0x3333 // binary: 00110011...
	h16_4 uint16 = 0x0111 // binary: 0000000100010001

	m8_1 uint8 = 0x55 // binary: 01010101
	m8_2 uint8 = 0x33 // binary: 00110011
	h8_4 uint8 = 0x11 // binary: 00010001
)

// PopCnt64Before counts the number of "1" before specified bit position `iBit`.
//
// "pop-cnt" means population( of "1" ) count.
//
// E.g.:
//		(Significant bits on left)
//
//		3 = ...011	PopCnt64Before(3, 0) == 0
//		          	PopCnt64Before(3, 1) == 1
//		          	PopCnt64Before(3, 2) == 2
//		          	PopCnt64Before(3, 3) == 2
//
// This algorithm has more introduction in:
// https://en.wikipedia.org/wiki/Hamming_weight#Efficient_implementation
func PopCnt64Before(n uint64, iBit uint32) uint32 {
	n = n & ((uint64(1) << iBit) - 1)

	n -= (n >> 1) & m64_1                // put count of each 2 bits into those 2 bits
	n = (n & m64_2) + ((n >> 2) & m64_2) // put count of each 4 bits into those 4 bits
	n = (n + (n >> 4)) & m64_4           // put count of each 8 bits into thoes 8 bits

	return uint32((n * h64_8) >> 56) // returns left 8 bits of x + (x << 8) + (x << 16) + (x<<24) + ...
}

// PopCnt32Before is similar to PopCnt64Before, except it accepts a uint32 `n`
func PopCnt32Before(n uint32, iBit uint32) int {
	n = n & ((uint32(1) << iBit) - 1)

	n -= (n >> 1) & m32_1
	n = (n & m32_2) + ((n >> 2) & m32_2)
	n = (n + (n >> 4)) & m32_4

	return int((n * h32_8) >> 24)
}

// PopCnt16Before is similar to PopCnt64Before, except it accepts a uint16 `n`
// instead of a uint64 `n`.
// It is about 2.5% faster than `PopCnt64Before(uint16(i), iBit)`.
func PopCnt16Before(n uint16, iBit uint32) int {
	n = n & ((uint16(1) << iBit) - 1)

	n -= (n >> 1) & m16_1
	n = (n & m16_2) + ((n >> 2) & m16_2)

	// Because 4 4-bit segments have at most 16 ones.
	// It is not enought for the left most 4 bit to store the count.
	return int(((n * h16_4) >> 12) + (n & 0x0f))
}

// PopCnt8Before is similar to PopCnt64Before, except it accepts a uint8 `n`
// It is about 5% faster than `PopCnt64Before(uint16(i), iBit)`.
func PopCnt8Before(n uint8, iBit uint32) int {
	n = n & ((uint8(1) << iBit) - 1)

	n -= (n >> 1) & m8_1
	n = (n & m8_2) + ((n >> 2) & m8_2)

	return int((n * h8_4) >> 4)
}
