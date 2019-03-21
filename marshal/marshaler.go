// Package marshal provides marshaling API definition and with several commonly
// used Marshaler suchas uint32 and uint64 etc.
package marshal

import (
	"errors"
	"reflect"
)

var (
	// ErrNotSlice indicates it expects a slice type but not
	ErrNotSlice = errors.New("it is not a slice")
	// ErrUnknownEltType indicates a type this package does not support.
	ErrUnknownEltType = errors.New("element type is unknown")

	// ErrNotFixedSize indicates the size of value of a type can not be
	// determined by its type.
	// Such slice of interface.
	ErrNotFixedSize = errors.New("element type is not fixed size")
)

// A Marshaler converts one element between serialized byte stream
// and in-memory data structure.
type Marshaler interface {
	// Convert into serialised byte stream.
	Marshal(interface{}) []byte

	// Read byte stream and convert it back to typed data.
	Unmarshal([]byte) (int, interface{})

	// GetSize returns the size in byte after marshaling v.
	// If v is of type this marshaler can not marshal, it panics.
	GetSize(v interface{}) int

	// GetMarshaledSize returns size of the marshaled value.
	// Marshaled element may be var-length.
	// This function is used to determine element size without the need of
	// unmarshaling it.
	GetMarshaledSize([]byte) int
}

// GetMarshaler returns a `Marshaler` implementation for type of `e`
func GetMarshaler(e interface{}) (Marshaler, error) {
	k := reflect.ValueOf(e).Kind()
	return getMarshalerByKind(k)
}

// GetSliceEltMarshaler creates a `Marshaler` for type of element in slice `s`
func GetSliceEltMarshaler(s interface{}) (Marshaler, error) {
	sl := reflect.ValueOf(s)
	if sl.Kind() != reflect.Slice {
		return nil, ErrNotSlice
	}

	eltKind := reflect.TypeOf(s).Elem().Kind()

	return getMarshalerByKind(eltKind)
}

func getMarshalerByKind(k reflect.Kind) (Marshaler, error) {
	var m Marshaler
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
func (s String16) Marshal(d interface{}) []byte {
	ss := d.(string)
	l := len(ss)
	rst := make([]byte, 2, 2+l)
	rst[0] = byte(l >> 8)
	rst[1] = byte(l)
	return append(rst, []byte(ss)...)
}

// Unmarshal converts slice of 2 bytes to uint16.
// It returns number bytes consumed and an uint16.
func (s String16) Unmarshal(b []byte) (int, interface{}) {
	l := int(b[0])<<8 + int(b[1])
	ss := string(b[2 : 2+l])
	return 2 + l, ss
}

// GetSize returns number of byte required to marshal a string.
// It is len(str) + 2;
func (s String16) GetSize(d interface{}) int {
	ss := d.(string)
	l := len(ss)
	return 2 + l
}

// GetMarshaledSize returned size of marshaled data.
func (s String16) GetMarshaledSize(b []byte) int {
	l := int(b[0])<<8 + int(b[1])
	return 2 + l
}
