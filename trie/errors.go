package trie

import "errors"

var (

	// ErrKeyOutOfOrder means keys to create Trie are not ascendingly ordered.
	ErrKeyOutOfOrder = errors.New("keys not ascending sorted")

	// ErrIncompatible means it is trying to unmarshal data from an incompatible
	// version.
	ErrIncompatible = errors.New("incompatible with marshaled data")
)
