package array

import (
	"github.com/openacid/errors"
)

// ErrIndexNotAscending indicates that the indexes to initialize an Array is not
// in ascending order.
var ErrIndexNotAscending = errors.New("index must be an ascending ordered slice")

// ErrIndexLen indicates that the number of indexes does not equal the number of
// elements, when initializing an Array.
var ErrIndexLen = errors.New("the length of indexes and elts must be equal")

// ErrUnknownSize indicates that the element is not a fixed-type.
var ErrUnknownSize = errors.New("the size of array size is unknown")

// ErrDifferentEltSize indicates that two elements have different size.
var ErrDifferentEltSize = errors.New("elements have different size")
