// Package sigbits extracts significant bits from a sorted list of strings.
// Significant bits are the minimal bit set to distinguish all of the strings.
// E.g. we have 3 strings:
//
//     ab  // bin: 0110 0001 0110 0010
//     ac  // bin: 0110 0001 0110 0011
//     b   // bin: 0110 0010
//
// The significant bits are [15, 7], which means the first different bits of
// "ab" and "ac" is the 15-th bit, and the 7-th bit for "ac" and "b".
//
// Since 0.1.9
package sigbits

// SigBits is a helper to get information about how many bits are required to
// distinguish a set of keys.
//
// Since 0.1.9
type SigBits struct {
	keys    []string
	sigbits []int32
}

// New creates a SigBits struct.
//
// Since 0.1.9
func New(keys []string) *SigBits {
	sb := &SigBits{
		keys:    keys,
		sigbits: FirstDiffBits(keys),
	}

	return sb
}
