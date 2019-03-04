// Package marshal provides marshalling API definition and with several commonly
// used Marshaller suchas uint32 and uint64 etc.
package marshal

import (
	"encoding/binary"
	"errors"
	"reflect"
)

var (
	// ErrNotSlice indicates it expects a slice type but not
	ErrNotSlice = errors.New("it is not a slice")
	// ErrUnknownEltType indicates a type this package does not support.
	ErrUnknownEltType = errors.New("element type is unknown")
)

// A Marshaller converts one element between serialized byte stream
// and in-memory data structure.
type Marshaller interface {
	// Convert into serialised byte stream.
	Marshal(interface{}) []byte

	// Read byte stream and convert it back to typed data.
	Unmarshal([]byte) (int, interface{})

	// Marshaled element may be var-length.
	// This function is used to determine element size without the need of
	// unmarshaling it.
	GetMarshaledSize([]byte) int
}

// GetMarshaller returns a `Marshaller` implementation for type of `e`
func GetMarshaller(e interface{}) (Marshaller, error) {
	k := reflect.ValueOf(e).Kind()
	return getMarshallerByKind(k)
}

// GetMarshaller returns a `Marshaller` implementation for element type of slice `s`
func GetSliceEltMarshaller(s interface{}) (Marshaller, error) {
	sl := reflect.ValueOf(s)
	if sl.Kind() != reflect.Slice {
		return nil, ErrNotSlice
	}

	eltKind := reflect.TypeOf(s).Elem().Kind()

	return getMarshallerByKind(eltKind)
}

func getMarshallerByKind(k reflect.Kind) (Marshaller, error) {
	var m Marshaller
	switch k {
	case reflect.Uint16:
		m = U16{}
	case reflect.Uint32:
		m = U32{}
	case reflect.Uint64:
		m = U64{}
	default:
		return nil, ErrUnknownEltType
	}

	return m, nil
}

// String16 converts uint16 to slice of 2 bytes and back.
type String16 struct{}

// Marshal converts uint16 to slice of 2 bytes.
func (c String16) Marshal(d interface{}) []byte {
	s := d.(string)
	l := len(s)
	rst := make([]byte, 2, 2+l)
	rst[0] = byte(l >> 8)
	rst[1] = byte(l)
	return append(rst, []byte(s)...)
}

// Unmarshal converts slice of 2 bytes to uint16.
// It returns number bytes consumed and an uint16.
func (c String16) Unmarshal(b []byte) (int, interface{}) {
	l := int(b[0])<<8 + int(b[1])
	s := string(b[2 : 2+l])
	return 2 + l, s
}

// GetMarshaledSize returned size of marshaled data.
func (c String16) GetMarshaledSize(b []byte) int {
	l := int(b[0])<<8 + int(b[1])
	return 2 + l
}

// U64 converts uint64 to slice of 4 bytes and back.
type U64 struct{}

// Marshal converts uint64 to slice of 8 bytes.
func (c U64) Marshal(d interface{}) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, d.(uint64))
	return b
}

// Unmarshal converts slice of 4 bytes to uint64.
// It returns number bytes consumed and an uint64.
func (c U64) Unmarshal(b []byte) (int, interface{}) {

	size := int(8)
	s := b[:size]

	d := binary.LittleEndian.Uint64(s)
	return size, d
}

// GetMarshaledSize returns 8.
func (c U64) GetMarshaledSize(b []byte) int {
	return 8
}

// U32 converts uint32 to slice of 4 bytes and back.
type U32 struct{}

// Marshal converts uint32 to slice of 4 bytes.
func (c U32) Marshal(d interface{}) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, d.(uint32))
	return b
}

// Unmarshal converts slice of 4 bytes to uint32.
// It returns number bytes consumed and an uint32.
func (c U32) Unmarshal(b []byte) (int, interface{}) {

	size := int(4)
	s := b[:size]

	d := binary.LittleEndian.Uint32(s)
	return size, d
}

// GetMarshaledSize returns 4.
func (c U32) GetMarshaledSize(b []byte) int {
	return 4
}

// U16 converts uint16 to slice of 4 bytes and back.
type U16 struct{}

// Marshal converts uint16 to slice of 4 bytes.
func (c U16) Marshal(d interface{}) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, d.(uint16))
	return b
}

// Unmarshal converts slice of 2 bytes to uint16.
// It returns number bytes consumed and an uint16.
func (c U16) Unmarshal(b []byte) (int, interface{}) {

	size := int(2)
	s := b[:size]

	d := binary.LittleEndian.Uint16(s)
	return size, d
}

// GetMarshaledSize returns 2.
func (c U16) GetMarshaledSize(b []byte) int {
	return 2
}
