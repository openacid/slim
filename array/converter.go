package array

import (
	"encoding/binary"
)

// A Converter is used to convert one element between serialized byte stream
// and in-memory data structure.
type Converter interface {
	// Convert into serialised byte stream.
	Marshal(interface{}) []byte

	// Read byte stream and convert it back to typed data.
	Unmarshal([]byte) (int, interface{})

	// Marshaled element may be var-length.
	// This function is used to determine element size without the need of
	// unmarshaling it.
	GetMarshaledSize([]byte) int
}

// U16Conv converts uint16 to slice of 2 bytes and back.
type U16Conv struct{}

// Marshal converts uint16 to slice of 2 bytes.
func (c U16Conv) Marshal(d interface{}) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, d.(uint16))
	return b
}

// Unmarshal converts slice of 2 bytes to uint16.
// It returns number bytes consumed and an uint16.
func (c U16Conv) Unmarshal(b []byte) (int, interface{}) {
	elt := binary.LittleEndian.Uint16(b[:2])

	return 2, elt
}

// GetMarshaledSize returns 2.
func (c U16Conv) GetMarshaledSize(b []byte) int {
	return 2
}

// U32Conv converts uint32 to slice of 4 bytes and back.
type U32Conv struct{}

// Marshal converts uint32 to slice of 4 bytes.
func (c U32Conv) Marshal(d interface{}) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, d.(uint32))
	return b
}

// Unmarshal converts slice of 4 bytes to uint32.
// It returns number bytes consumed and an uint32.
func (c U32Conv) Unmarshal(b []byte) (int, interface{}) {

	size := int(4)
	s := b[:size]

	d := binary.LittleEndian.Uint32(s)
	return size, d
}

// GetMarshaledSize returns 4.
func (c U32Conv) GetMarshaledSize(b []byte) int {
	return 4
}

// ByteConv converts a byte slice into fixed length slice.
// Result slice length is defined by ByteConv.EltSize .
type ByteConv struct {
	EltSize int
}

// Marshal converts byte slice to byte slice.
func (c ByteConv) Marshal(d interface{}) []byte {
	return d.([]byte)
}

// Unmarshal copies fixed length slice out of source byte slice.
func (c ByteConv) Unmarshal(b []byte) (int, interface{}) {
	// TODO: converted value should be copied. referencing source data is dangerous.
	s := b[:c.EltSize]
	return c.EltSize, s
}

// GetMarshaledSize returns c.EltSize
func (c ByteConv) GetMarshaledSize(b []byte) int {
	return c.EltSize
}

// U32to3ByteConv converts uint32 to a slice of 3 bytes and back.
//
// Thus the int passed in should be smaller than 2^24.
type U32to3ByteConv struct{}

// Marshal converts uint32 to a slice of 3 bytes.
func (c U32to3ByteConv) Marshal(d interface{}) []byte {
	size := 4
	b := make([]byte, size)
	binary.LittleEndian.PutUint32(b, d.(uint32))
	return b[:3]
}

// Unmarshal converts slice of 3 bytes to a uint32.
func (c U32to3ByteConv) Unmarshal(b []byte) (int, interface{}) {
	size := 4
	s := make([]byte, size)
	copy(s[:3], b)

	d := binary.LittleEndian.Uint32(s)
	return size, d
}

// GetMarshaledSize returns 3
func (c U32to3ByteConv) GetMarshaledSize(b []byte) int {
	return int(3)
}
