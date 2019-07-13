package sigbits

// CountPrefies counts the number of n-bit prefix in a subset of keys [keyStart,
// keyEnd):
//
// The first return value is the index of the first bit where there is a
// different bit among all keys.
// The second return values is a []int32 of maxitem + 1 elements.
// The ith element is the number of i-bits long prefix.
//
// Since 0.1.9
func (sb *SigBits) CountPrefixes(keyStart, keyEnd, maxitem int32) (int32, []int32) {
	return countPrefixes(sb.sigbits[keyStart:keyEnd-1], maxitem)
}
