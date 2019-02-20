package array

import (
	"encoding/binary"
)

// A EltConverter is used to convert one element between serialized byte stream
// and in-memory data structure
type EltConverter interface {
	MarshalElt(interface{}) []byte
	UnmarshalElt([]byte) (uint32, interface{})

	// Marshaled element may be var-length.
	// This function is used to determine element size without the need of
	// unmarshaling it.
	GetMarshaledEltSize([]byte) uint32
}

type U16Conv struct{}

func (c U16Conv) MarshalElt(d interface{}) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, d.(uint16))
	return b
}

func (c U16Conv) UnmarshalElt(b []byte) (uint32, interface{}) {
	elt := binary.LittleEndian.Uint16(b[:2])

	return 2, elt
}

func (c U16Conv) GetMarshaledEltSize(b []byte) uint32 {
	return 2
}

type U32Conv struct{}

func (c U32Conv) MarshalElt(d interface{}) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, d.(uint32))
	return b
}

func (c U32Conv) UnmarshalElt(b []byte) (uint32, interface{}) {

	size := uint32(4)
	s := b[:size]

	d := binary.LittleEndian.Uint32(s)
	return size, d
}

func (c U32Conv) GetMarshaledEltSize(b []byte) uint32 {
	return 4
}

type ByteConv struct {
	EltSize uint32
}

func (c ByteConv) MarshalElt(d interface{}) []byte {
	return d.([]byte)
}

func (c ByteConv) UnmarshalElt(b []byte) (uint32, interface{}) {
	s := b[:c.EltSize]
	return c.EltSize, s
}

func (c ByteConv) GetMarshaledEltSize(b []byte) uint32 {
	return c.EltSize
}

type U32to3ByteConv struct{}

func (c U32to3ByteConv) MarshalElt(d interface{}) []byte {
	size := 4
	b := make([]byte, size)
	binary.LittleEndian.PutUint32(b, d.(uint32))
	return b[:3]
}

func (c U32to3ByteConv) UnmarshalElt(b []byte) (uint32, interface{}) {
	size := uint32(4)
	s := make([]byte, size)
	copy(s[:3], b)

	d := binary.LittleEndian.Uint32(s)
	return size, d
}

func (c U32to3ByteConv) GetMarshaledEltSize(b []byte) uint32 {
	return uint32(3)
}
