package array

// ArrayU16 is an implementation of Array with uint16 element
type ArrayU16 struct {
	Array32Index
	Data []uint16
}

// NewArrayU16 creates a ArrayU16
func NewArrayU16(index []uint32, elts []uint16) (a *ArrayU16, err error) {

	a = &ArrayU16{Data: elts}

	err = a.InitIndexBitmap(index)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Get2 returns value indexed by `idx` and a bool indicate if the value is
// found.
func (a *ArrayU16) Get2(idx uint32) (uint16, bool) {
	i, ok := a.GetEltIndex(idx)
	if ok {
		return a.Data[i], true
	}

	return 0, false
}
