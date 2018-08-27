package bit

const m1 uint64 = 0x5555555555555555  // binary: 0101...
const m2 uint64 = 0x3333333333333333  // binary: 00110011...
const m4 uint64 = 0x0f0f0f0f0f0f0f0f  // binary: 0000111100001111...
const h01 uint64 = 0x0101010101010101 // the sum of 256 to the power of 0,1,2,3...

// this algorithm has more introduction in `https://en.wikipedia.org/wiki/Hamming_weight`
func Cnt1Before(n uint64, iBit uint32) uint32 {
	n = n & ((1 << iBit) - 1)

	n -= (n >> 1) & m1             // put count of each 2 bits into those 2 bits
	n = (n & m2) + ((n >> 2) & m2) // put count of each 4 bits into those 4 bits
	n = (n + (n >> 4)) & m4        // put count of each 8 bits into thoes 8 bits

	return uint32((n * h01) >> 56) // returns left 8 bits of x + (x << 8) + (x << 16) + (x<<24) + ...
}
