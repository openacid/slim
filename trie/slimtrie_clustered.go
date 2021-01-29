package trie

import (
	"bytes"
	"unsafe"
)

type clusteredInner struct {
	FirstLeafId int32
	Offsets     []uint32
	Bytes       []byte
}

// initClusteredInner initiates a clusteredInner struct.
//
// Since 0.5.12
func (st *SlimTrie) initClusteredInner(
	ithClusteredInner int32,
	cl *clusteredInner,
) {
	ns := st.inner

	i := ithClusteredInner
	s, e := ns.Clustered.Starts[i], ns.Clustered.Starts[i+1]
	offsets := ns.Clustered.Offsets[s : e+1]

	bottom := len(st.levels) - 1

	cl.FirstLeafId = st.levels[bottom-1].total + int32(s)
	cl.Offsets = offsets
	cl.Bytes = ns.Clustered.Bytes
}

func newClusteredInner(first int32, keys []string, prefixLen int32) *clusteredInner {
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

	r := &clusteredInner{
		FirstLeafId: first,
		Offsets:     offsets,
		Bytes:       records,
	}
	return r
}

func (cl *clusteredInner) keyCnt() int {
	// There are n + 1 offsets. the last one is len(cl.clusteredInner)
	return len(cl.Offsets) - 1
}

func (cl *clusteredInner) keys() [][]byte {
	n := cl.keyCnt()
	keys := make([][]byte, 0, n)

	for i := 0; i < n; i++ {
		k := cl.Bytes[cl.Offsets[i]:cl.Offsets[i+1]]
		keys = append(keys, k)
	}

	return keys
}

func (cl *clusteredInner) get(key string) int32 {

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
func (cl *clusteredInner) firstLeafId() int32 {
	return cl.FirstLeafId
}

// lastLeafId returns the node id of the last record.
func (cl *clusteredInner) lastLeafId() int32 {
	return cl.FirstLeafId + int32(len(cl.Offsets)) - 2
}

func (cl *clusteredInner) search(key string) (int32, int32, int32) {

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
