package array

// Do NOT edit. re-generate this file with "go generate ./..."

// ArrayU16 is an implementation of Array32Index with uint16 element
type ArrayU16 struct {
	ArrayBase
}

// NewU16 creates a ArrayU16
func NewU16(index []int32, elts []uint16) (a *ArrayU16, err error) {
	a = &ArrayU16{}
	err = a.Init(index, elts)
	if err != nil {
		a = nil
	}
	return a, err
}

// Get returns value indexed by idx and a bool indicate if the value is
// found.
func (a *ArrayU16) Get(idx int32) (uint16, bool) {
	bs, ok := a.GetBytes(idx, 2)
	if ok {
		return endian.Uint16(bs), true
	}

	return 0, false
}

// ArrayU32 is an implementation of Array32Index with uint32 element
type ArrayU32 struct {
	ArrayBase
}

// NewU32 creates a ArrayU32
func NewU32(index []int32, elts []uint32) (a *ArrayU32, err error) {
	a = &ArrayU32{}
	err = a.Init(index, elts)
	if err != nil {
		a = nil
	}
	return a, err
}

// Get returns value indexed by idx and a bool indicate if the value is
// found.
func (a *ArrayU32) Get(idx int32) (uint32, bool) {
	bs, ok := a.GetBytes(idx, 4)
	if ok {
		return endian.Uint32(bs), true
	}

	return 0, false
}

// ArrayU64 is an implementation of Array32Index with uint64 element
type ArrayU64 struct {
	ArrayBase
}

// NewU64 creates a ArrayU64
func NewU64(index []int32, elts []uint64) (a *ArrayU64, err error) {
	a = &ArrayU64{}
	err = a.Init(index, elts)
	if err != nil {
		a = nil
	}
	return a, err
}

// Get returns value indexed by idx and a bool indicate if the value is
// found.
func (a *ArrayU64) Get(idx int32) (uint64, bool) {
	bs, ok := a.GetBytes(idx, 8)
	if ok {
		return endian.Uint64(bs), true
	}

	return 0, false
}
