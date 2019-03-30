// Code generated 'by go generate ./...'; DO NOT EDIT.

package encode

import "encoding/binary"

// U16 converts uint16 to slice of 2 bytes and back.
type U16 struct{}

// Encode converts uint16 to slice of 2 bytes.
func (c U16) Encode(d interface{}) []byte {
	b := make([]byte, 2)
	v := d.(uint16)
	binary.LittleEndian.PutUint16(b, v)
	return b
}

// Decode converts slice of 2 bytes to uint16.
// It returns number bytes consumed and an uint16.
func (c U16) Decode(b []byte) (int, interface{}) {

	size := int(2)
	s := b[:size]

	d := binary.LittleEndian.Uint16(s)
	return size, d
}

// GetSize returns the size in byte after encoding v.
func (c U16) GetSize(d interface{}) int {
	return 2
}

// GetEncodedSize returns 2.
func (c U16) GetEncodedSize(b []byte) int {
	return 2
}

// U32 converts uint32 to slice of 4 bytes and back.
type U32 struct{}

// Encode converts uint32 to slice of 4 bytes.
func (c U32) Encode(d interface{}) []byte {
	b := make([]byte, 4)
	v := d.(uint32)
	binary.LittleEndian.PutUint32(b, v)
	return b
}

// Decode converts slice of 4 bytes to uint32.
// It returns number bytes consumed and an uint32.
func (c U32) Decode(b []byte) (int, interface{}) {

	size := int(4)
	s := b[:size]

	d := binary.LittleEndian.Uint32(s)
	return size, d
}

// GetSize returns the size in byte after encoding v.
func (c U32) GetSize(d interface{}) int {
	return 4
}

// GetEncodedSize returns 4.
func (c U32) GetEncodedSize(b []byte) int {
	return 4
}

// U64 converts uint64 to slice of 8 bytes and back.
type U64 struct{}

// Encode converts uint64 to slice of 8 bytes.
func (c U64) Encode(d interface{}) []byte {
	b := make([]byte, 8)
	v := d.(uint64)
	binary.LittleEndian.PutUint64(b, v)
	return b
}

// Decode converts slice of 8 bytes to uint64.
// It returns number bytes consumed and an uint64.
func (c U64) Decode(b []byte) (int, interface{}) {

	size := int(8)
	s := b[:size]

	d := binary.LittleEndian.Uint64(s)
	return size, d
}

// GetSize returns the size in byte after encoding v.
func (c U64) GetSize(d interface{}) int {
	return 8
}

// GetEncodedSize returns 8.
func (c U64) GetEncodedSize(b []byte) int {
	return 8
}

// I16 converts int16 to slice of 2 bytes and back.
type I16 struct{}

// Encode converts int16 to slice of 2 bytes.
func (c I16) Encode(d interface{}) []byte {
	b := make([]byte, 2)
	v := uint16(d.(int16))
	binary.LittleEndian.PutUint16(b, v)
	return b
}

// Decode converts slice of 2 bytes to int16.
// It returns number bytes consumed and an int16.
func (c I16) Decode(b []byte) (int, interface{}) {

	size := int(2)
	s := b[:size]

	d := int16(binary.LittleEndian.Uint16(s))
	return size, d
}

// GetSize returns the size in byte after encoding v.
func (c I16) GetSize(d interface{}) int {
	return 2
}

// GetEncodedSize returns 2.
func (c I16) GetEncodedSize(b []byte) int {
	return 2
}

// I32 converts int32 to slice of 4 bytes and back.
type I32 struct{}

// Encode converts int32 to slice of 4 bytes.
func (c I32) Encode(d interface{}) []byte {
	b := make([]byte, 4)
	v := uint32(d.(int32))
	binary.LittleEndian.PutUint32(b, v)
	return b
}

// Decode converts slice of 4 bytes to int32.
// It returns number bytes consumed and an int32.
func (c I32) Decode(b []byte) (int, interface{}) {

	size := int(4)
	s := b[:size]

	d := int32(binary.LittleEndian.Uint32(s))
	return size, d
}

// GetSize returns the size in byte after encoding v.
func (c I32) GetSize(d interface{}) int {
	return 4
}

// GetEncodedSize returns 4.
func (c I32) GetEncodedSize(b []byte) int {
	return 4
}

// I64 converts int64 to slice of 8 bytes and back.
type I64 struct{}

// Encode converts int64 to slice of 8 bytes.
func (c I64) Encode(d interface{}) []byte {
	b := make([]byte, 8)
	v := uint64(d.(int64))
	binary.LittleEndian.PutUint64(b, v)
	return b
}

// Decode converts slice of 8 bytes to int64.
// It returns number bytes consumed and an int64.
func (c I64) Decode(b []byte) (int, interface{}) {

	size := int(8)
	s := b[:size]

	d := int64(binary.LittleEndian.Uint64(s))
	return size, d
}

// GetSize returns the size in byte after encoding v.
func (c I64) GetSize(d interface{}) int {
	return 8
}

// GetEncodedSize returns 8.
func (c I64) GetEncodedSize(b []byte) int {
	return 8
}
