package pbcmpl

import "errors"

var (
	// ErrInvalidHeaderSize indicates the size in the header is incorrect
	ErrInvalidHeaderSize = errors.New("headersize is incorrect")
)
