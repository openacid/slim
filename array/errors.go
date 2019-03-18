package array

import (
	"github.com/openacid/errors"
)

// ErrIndexNotAscending means indexes to initialize a Array must be in
// ascending order.
var ErrIndexNotAscending = errors.New("index must be an ascending ordered slice")

// ErrIndexLen is returned if number of indexes does not equal the number of
// datas, when initializing a Array.
var ErrIndexLen = errors.New("the length of indexes and elts must be equal")

// ErrUnknownSize indicates the element size can not be deside by its type.
var ErrUnknownSize = errors.New("the size of array size is unknown")

// ErrDifferentEltSize indicates two elements has different size when creating
// an array.
var ErrDifferentEltSize = errors.New("elements have different size")
