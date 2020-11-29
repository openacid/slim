package trie

// GetI8 is same as Get() except it is optimized for int8.
//
// Since 0.5.10
func (st *SlimTrie) GetI8(key string) (int8, bool) {

	eqID := st.GetID(key)

	if eqID == -1 {
		return 0, false
	}

	ith, _ := st.getLeafIndex(eqID)

	v := int8(st.inner.Leaves.Bytes[ith])

	return v, true
}

// GetI16 is same as Get() except it is optimized for int16.
//
// Since 0.5.10
func (st *SlimTrie) GetI16(key string) (int16, bool) {

	eqID := st.GetID(key)

	if eqID == -1 {
		return 0, false
	}

	ith, _ := st.getLeafIndex(eqID)
	stIdx := ith << 1

	b := st.inner.Leaves.Bytes[stIdx : stIdx+2]

	v := int16(b[0]) | int16(b[1])<<8

	return v, true
}

// GetI32 is same as Get() except it is optimized for int32.
//
// Since 0.5.10
func (st *SlimTrie) GetI32(key string) (int32, bool) {

	eqID := st.GetID(key)

	if eqID == -1 {
		return 0, false
	}

	ith, _ := st.getLeafIndex(eqID)
	stIdx := ith << 2

	b := st.inner.Leaves.Bytes[stIdx : stIdx+4]

	v := int32(b[0]) | int32(b[1])<<8 | int32(b[2])<<16 | int32(b[3])<<24

	return v, true
}

// GetI64 is same as Get() except it is optimized for int64.
//
// Since 0.5.10
func (st *SlimTrie) GetI64(key string) (int64, bool) {

	eqID := st.GetID(key)

	if eqID == -1 {
		return 0, false
	}

	ith, _ := st.getLeafIndex(eqID)
	stIdx := ith << 3

	b := st.inner.Leaves.Bytes[stIdx : stIdx+8]

	v := int64(b[0]) | int64(b[1])<<8 | int64(b[2])<<16 | int64(b[3])<<24 | int64(b[4])<<32 | int64(b[5])<<40 | int64(b[6])<<48 | int64(b[7])<<56

	return v, true
}
