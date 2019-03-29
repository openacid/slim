package array

// Do NOT edit. re-generate this file with "go generate ./..."

// U16 is an implementation of Base with uint16 element
type U16 struct {
	Base
}

// NewU16 creates a U16
func NewU16(index []int32, elts []uint16) (a *U16, err error) {
	a = &U16{}
	err = a.Init(index, elts)
	if err != nil {
		a = nil
	}
	return a, err
}

// Get returns value at "idx" and a bool indicating if the value is
// found.
func (a *U16) Get(idx int32) (uint16, bool) {
	bs, ok := a.GetBytes(idx, 2)
	if ok {
		return endian.Uint16(bs), true
	}

	return 0, false
}

// U32 is an implementation of Base with uint32 element
type U32 struct {
	Base
}

// NewU32 creates a U32
func NewU32(index []int32, elts []uint32) (a *U32, err error) {
	a = &U32{}
	err = a.Init(index, elts)
	if err != nil {
		a = nil
	}
	return a, err
}

// Get returns value at "idx" and a bool indicating if the value is
// found.
func (a *U32) Get(idx int32) (uint32, bool) {
	bs, ok := a.GetBytes(idx, 4)
	if ok {
		return endian.Uint32(bs), true
	}

	return 0, false
}

// U64 is an implementation of Base with uint64 element
type U64 struct {
	Base
}

// NewU64 creates a U64
func NewU64(index []int32, elts []uint64) (a *U64, err error) {
	a = &U64{}
	err = a.Init(index, elts)
	if err != nil {
		a = nil
	}
	return a, err
}

// Get returns value at "idx" and a bool indicating if the value is
// found.
func (a *U64) Get(idx int32) (uint64, bool) {
	bs, ok := a.GetBytes(idx, 8)
	if ok {
		return endian.Uint64(bs), true
	}

	return 0, false
}
