package trie

// levelCachePolicy describes what to cache and the expected benefits.
//
// Since 0.5.12
type levelCachePolicy struct {
	// the steps it takes to walk all keys.
	steps int64

	// total number of steps reduced by cache
	reduced int64

	// the level indexes to cache
	levels []int32
}
