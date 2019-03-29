package array

import (
	"github.com/openacid/errors"
)

// ErrIndexNotAscending indicates that the indexes to initialize an Array is not
// in ascending order.
//
// Since 0.2.0
var ErrIndexNotAscending = errors.New("index must be an ascending ordered slice")

// ErrIndexLen indicates that the number of indexes does not equal the number of
// elements, when initializing an Array.
//
// Since 0.2.0
var ErrIndexLen = errors.New("the length of indexes and elts must be equal")

// ErrUnknownSize indicates that the element is not a fixed-type.
//
// Since 0.2.0
var ErrUnknownSize = errors.New("the size of array size is unknown")

// ErrDifferentEltSize indicates that two elements have different size.
//
// Since 0.2.0
var ErrDifferentEltSize = errors.New("elements have different size")
