package encode

import (
	"bytes"
	"encoding/binary"
	"reflect"

	"github.com/openacid/errors"
)

// defaultEndian is default endian
var defaultEndian = binary.LittleEndian

// TypeEncoder provides encoding for fixed size types.
// Such as int32 or struct { X int32; Y int64; }
//
// "int" is not a fixed size type: int on different platform has different size,
// 4 or 8 bytes.
//
// "[]int32" is not a fixed size type: the data size is also defined by the
// number of elements.
type TypeEncoder struct {
	// Endian defines the byte order to encode a value.
	// By default it is binary.LittleEndian
	Endian binary.ByteOrder
	// Type is the data type to encode.
	Type reflect.Type
	// Size is the encoded size of this type.
	Size int
}

// NewTypeEncoder creates a *TypeEncoder by a value.
// The value "zero" defines what type this Encoder can deal with and must be a
// fixed size type.
func NewTypeEncoder(zero interface{}) (*TypeEncoder, error) {
	return NewTypeEncoderEndian(zero, nil)
}

// NewTypeEncoderEndian creates a *TypeEncoder with a specified byte order.
//
// "endian" could be binary.LittleEndian or binary.BigEndian.
func NewTypeEncoderEndian(zero interface{}, endian binary.ByteOrder) (*TypeEncoder, error) {
	if endian == nil {
		endian = defaultEndian
	}
	m := &TypeEncoder{
		Endian: endian,
		Type:   reflect.Indirect(reflect.ValueOf(zero)).Type(),
		Size:   binary.Size(zero),
	}

	if m.Size == -1 {
		return nil, errors.Wrapf(ErrNotFixedSize, "type: %v", reflect.TypeOf(zero))
	}
	if m.Type.Kind() == reflect.Slice {
		return nil, errors.Wrapf(ErrNotFixedSize, "slice size is not fixed")
	}

	return m, nil
}

// NewTypeEncoderEndianByType creates a *TypeEncoder for specified type and with a specified byte order.
//
// "endian" could be binary.LittleEndian or binary.BigEndian.
func NewTypeEncoderEndianByType(t reflect.Type, endian binary.ByteOrder) (*TypeEncoder, error) {
	v := reflect.New(t)
	return NewTypeEncoderEndian(v.Interface(), endian)
}

// Encode converts a m.Type value to byte slice.
// If a different type value from the one used with NewTypeEncoder passed in,
// it panics.
func (m *TypeEncoder) Encode(d interface{}) []byte {
	if reflect.Indirect(reflect.ValueOf(d)).Type() != m.Type {
		panic("different type from TypeEncoder.Type")
	}

	b := bytes.NewBuffer(make([]byte, 0, m.Size))
	err := binary.Write(b, m.Endian, d)
	if err != nil {
		// there should not be any error if type is fixed size
		panic(err)
	}
	return b.Bytes()
}

// Decode converts byte slice to a pointer to Type value.
// It returns number bytes consumed and an Type value in interface{}.
func (m *TypeEncoder) Decode(b []byte) (int, interface{}) {

	b = b[0:m.Size]
	v := reflect.New(m.Type)
	err := binary.Read(bytes.NewBuffer(b), m.Endian, v.Interface())
	if err != nil {
		panic(err)
	}
	return m.Size, reflect.Indirect(v).Interface()
}

// GetSize returns m.Size.
func (m *TypeEncoder) GetSize(d interface{}) int {
	return m.Size
}

// GetEncodedSize returns m.Size.
func (m *TypeEncoder) GetEncodedSize(b []byte) int {
	return m.Size
}
