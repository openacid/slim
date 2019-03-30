// Package encode provides encoding API definition and with several commonly
// used Encoder suchas uint32 and uint64 etc.
package encode

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

// A Encoder converts one element between serialized byte stream
// and in-memory data structure.
type Encoder interface {
	// Convert into serialized byte stream.
	Encode(interface{}) []byte

	// Read byte stream and convert it back to typed data.
	Decode([]byte) (int, interface{})

	// GetSize returns the size in byte after encoding v.
	// If v is of type this encoder can not encode, it panics.
	GetSize(v interface{}) int

	// GetEncodedSize returns size of the encoded value.
	// Encoded element may be var-length.
	// This function is used to determine element size without the need of
	// encoding it.
	GetEncodedSize([]byte) int
}

// EncoderOf returns a `Encoder` implementation for type of `e`
func EncoderOf(e interface{}) (Encoder, error) {
	k := reflect.ValueOf(e).Kind()
	return EncoderByKind(k)
}

// GetSliceEltEncoder creates a `Encoder` for type of element in slice `s`
func GetSliceEltEncoder(s interface{}) (Encoder, error) {
	sl := reflect.ValueOf(s)
	if sl.Kind() != reflect.Slice {
		return nil, ErrNotSlice
	}

	eltKind := reflect.TypeOf(s).Elem().Kind()

	return EncoderByKind(eltKind)
}

func EncoderByKind(k reflect.Kind) (Encoder, error) {
	var m Encoder
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

// Encode converts uint16 to slice of 2 bytes.
func (s String16) Encode(d interface{}) []byte {
	ss := d.(string)
	l := len(ss)
	rst := make([]byte, 2, 2+l)
	rst[0] = byte(l >> 8)
	rst[1] = byte(l)
	return append(rst, []byte(ss)...)
}

// Decode converts slice of 2 bytes to uint16.
// It returns number bytes consumed and an uint16.
func (s String16) Decode(b []byte) (int, interface{}) {
	l := int(b[0])<<8 + int(b[1])
	ss := string(b[2 : 2+l])
	return 2 + l, ss
}

// GetSize returns number of byte required to encode a string.
// It is len(str) + 2;
func (s String16) GetSize(d interface{}) int {
	ss := d.(string)
	l := len(ss)
	return 2 + l
}

// GetEncodedSize returned size of encoded data.
func (s String16) GetEncodedSize(b []byte) int {
	l := int(b[0])<<8 + int(b[1])
	return 2 + l
}
