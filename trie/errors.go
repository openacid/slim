package trie

import "github.com/openacid/errors"

var (
	// ErrTooManyTrieNodes indicates the number of trie nodes(not number of
	// keys) exceeded.
	ErrTooManyTrieNodes = errors.New("exceeds max node count=65536")

	// ErrTrieBranchValueOverflow indicate input key consists of a word greater
	// than the max 4-bit word(0x0f).
	ErrTrieBranchValueOverflow = errors.New("branch value must <=0x0f")

	// ErrDuplicateKeys indicates two keys are identical.
	ErrDuplicateKeys = errors.New("keys can not be duplicate")

	// ErrKVLenNotMatch means the keys and values to create Trie has different
	// number of elements.
	ErrKVLenNotMatch = errors.New("length of keys and values not equal")

	// ErrKeyOutOfOrder means keys to create Trie are not ascendingly ordered.
	ErrKeyOutOfOrder = errors.New("keys not ascending sorted")
)
