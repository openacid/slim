package trie

import (
	"bytes"
	"unsafe"
)

// newClusteredLeaves creates a ClusteredLeaves struct.
//
// Since 0.5.12
func newClusteredLeaves(firstLeafId int32, keys []string, prefixLen int32) *ClusteredLeaves {
	size := 0
	offsetCnt := int32(len(keys)) + 1
	for i := int32(0); i < int32(len(keys)); i++ {
		size += len(keys[i]) - int(prefixLen)
	}

	records := make([]byte, 0, size)
	offsets := make([]uint32, offsetCnt)

	for i := int32(0); i < int32(len(keys)); i++ {
		offsets[i] = uint32(len(records))
		records = append(records, keys[i][prefixLen:]...)
	}

	offsets[offsetCnt-1] = uint32(len(records))

	r := &ClusteredLeaves{
		FirstLeafId: firstLeafId,
		Offsets:     offsets,
		Bytes:       records,
	}
	return r
}

func (cl *ClusteredLeaves) keyCnt() int {
	// There are n + 1 offsets. the last one is len(cl.ClusteredLeaves)
	return len(cl.Offsets) - 1
}

func (cl *ClusteredLeaves) keys() [][]byte {
	n := cl.keyCnt()
	keys := make([][]byte, 0, n)

	for i := 0; i < n; i++ {
		k := cl.Bytes[cl.Offsets[i]:cl.Offsets[i+1]]
		keys = append(keys, k)
	}

	return keys
}

func (cl *ClusteredLeaves) get(key string) int32 {

	kBytes := *(*[]byte)(unsafe.Pointer(&key))
	n := cl.keyCnt()

	if n < 4 {
		for i := 0; i < n; i++ {
			k := cl.Bytes[cl.Offsets[i]:cl.Offsets[i+1]]
			if bytes.Compare(kBytes, k) == 0 {
				return cl.FirstLeafId + int32(i)
			}
		}
		return -1
	}

	s, e := 0, n
	for s < e {
		mid := (s + e) >> 1
		rst := bytes.Compare(
			kBytes,
			cl.Bytes[cl.Offsets[mid]:cl.Offsets[mid+1]])
		if rst == 0 {
			return cl.FirstLeafId + int32(mid)
		} else if rst > 0 {
			s = mid + 1
		} else {
			e = mid
		}
	}

	return -1
}

// firstLeafId returns the node id of the first record.
func (cl *ClusteredLeaves) firstLeafId() int32 {
	return cl.FirstLeafId
}

// lastLeafId returns the node id of the last record.
func (cl *ClusteredLeaves) lastLeafId() int32 {
	return cl.FirstLeafId + int32(len(cl.Offsets)) - 2
}

func (cl *ClusteredLeaves) search(key string) (int32, int32, int32) {

	kBytes := *(*[]byte)(unsafe.Pointer(&key))
	n := cl.keyCnt()

	s, e := 0, n
	for s < e {
		mid := (s + e) >> 1
		rst := bytes.Compare(
			kBytes,
			cl.Bytes[cl.Offsets[mid]:cl.Offsets[mid+1]])
		if rst == 0 {
			lId := int32(-1)
			eqId := cl.FirstLeafId + int32(mid)
			rId := int32(-1)

			if mid > 0 {
				lId = eqId - 1
			}

			if mid < n-1 {
				rId = eqId + 1
			}
			return lId, eqId, rId
		} else if rst > 0 {
			s = mid + 1
		} else {
			e = mid
		}
	}

	lId := int32(-1)
	if e > 0 {
		lId = cl.FirstLeafId + int32(e) - 1
	}
	eqId := int32(-1)
	rId := int32(-1)
	if e < n {
		rId = cl.FirstLeafId + int32(e)
	}
	return lId, eqId, rId
}
