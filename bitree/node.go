package bitree

import (
	"encoding/binary"

	"github.com/openacid/slim/bits"
)

func New(steps []int32, ps []byte, elts [][]byte) []byte {

	if len(steps) != 6 {
		panic("Node64 accept exact 6 steps: 2^6=64")
	}

	bs := make([]byte, 32)

	// flag
	bs[0] = 0

	n := 1
	prev := int32(0)
	for _, step := range steps {
		// The first is absolue position.
		// From the second is relative.
		if step <= prev {
			panic("steps not in ascending order")
		}
		n += binary.PutUvarint(bs[n:], uint64(step-prev))
		prev = step
	}

	if len(ps) != len(elts) {
		panic("keys count != elts count")
	}

	bitmap := uint64(0)

	// ps are 6-bit int
	for _, p := range ps {
		bitmap |= 1 << p
	}

	binary.LittleEndian.PutUint64(bs[n:], bitmap)
	n += 8

	elttotal := 0
	for _, elt := range elts {
		elttotal += len(elt)
	}

	node := make([]byte, 0, n+elttotal)
	node = append(node, bs[:n]...)
	for _, elt := range elts {
		node = append(node, elt...)
	}

	return node
}

func Node64ReadPositions(d []byte) []int32 {
	n := 1
	rst := make([]int32, 6)
	p0 := int32(0)
	for i := 0; i < 6; i++ {
		p, nn := binary.Uvarint(d[n:])
		if nn < 0 {
			panic("invalid varint")
		}
		p0 += int32(p)
		n += nn

		rst[i] = p0
	}

	return rst
}

func Node64Get(d []byte, key byte) []byte {

	n := 1
	for i := 0; i < 6; i++ {
		_, nn := binary.Uvarint(d[n:])
		if nn < 0 {
			panic("invalid varint")
		}
		n += nn
	}

	bitmap := binary.LittleEndian.Uint64(d[n:])
	if bitmap&(1<<key) == 0 {
		return nil
	}

	n += 8

	total := bits.OnesCount64Before(bitmap, 64)
	nth := bits.OnesCount64Before(bitmap, uint(key))

	eltsize := (len(d) - n) / total

	return d[n+eltsize*nth : n+eltsize*(nth+1)]
}
