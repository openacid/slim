package array

import (
	"encoding/binary"
)

// A Converter is used to convert one element between serialized byte stream
// and in-memory data structure
type Converter interface {
	Marshal(interface{}) []byte
	Unmarshal([]byte) (uint32, interface{})

	// Marshaled element may be var-length.
	// This function is used to determine element size without the need of
	// unmarshaling it.
	GetMarshaledSize([]byte) uint32
}

type U16Conv struct{}

func (c U16Conv) Marshal(d interface{}) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, d.(uint16))
	return b
}

func (c U16Conv) Unmarshal(b []byte) (uint32, interface{}) {
	elt := binary.LittleEndian.Uint16(b[:2])

	return 2, elt
}

func (c U16Conv) GetMarshaledSize(b []byte) uint32 {
	return 2
}

type U32Conv struct{}

func (c U32Conv) Marshal(d interface{}) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, d.(uint32))
	return b
}

func (c U32Conv) Unmarshal(b []byte) (uint32, interface{}) {

	size := uint32(4)
	s := b[:size]

	d := binary.LittleEndian.Uint32(s)
	return size, d
}

func (c U32Conv) GetMarshaledSize(b []byte) uint32 {
	return 4
}

type ByteConv struct {
	EltSize uint32
}

func (c ByteConv) Marshal(d interface{}) []byte {
	return d.([]byte)
}

func (c ByteConv) Unmarshal(b []byte) (uint32, interface{}) {
	s := b[:c.EltSize]
	return c.EltSize, s
}

func (c ByteConv) GetMarshaledSize(b []byte) uint32 {
	return c.EltSize
}

type U32to3ByteConv struct{}

func (c U32to3ByteConv) Marshal(d interface{}) []byte {
	size := 4
	b := make([]byte, size)
	binary.LittleEndian.PutUint32(b, d.(uint32))
	return b[:3]
}

func (c U32to3ByteConv) Unmarshal(b []byte) (uint32, interface{}) {
	size := uint32(4)
	s := make([]byte, size)
	copy(s[:3], b)

	d := binary.LittleEndian.Uint32(s)
	return size, d
}

func (c U32to3ByteConv) GetMarshaledSize(b []byte) uint32 {
	return uint32(3)
}
