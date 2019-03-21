package marshal

import (
	"bytes"
	"encoding/binary"
	"reflect"

	"github.com/openacid/errors"
)

// defaultEndian is default endian
var defaultEndian = binary.LittleEndian

// TypeMarshaler provides marshaling for fixed size types.
// Such as int32 or struct { X int32; Y int64; }
//
// "int" is not a fixed size type: int on different platform has different size,
// 4 or 8 bytes.
//
// "[]int32" is not a fixed size type: the data size is also defined by the
// number of elements.
type TypeMarshaler struct {
	// Endian defines the byte order to marshal a value.
	// By default it is binary.LittleEndian
	Endian binary.ByteOrder
	// Type is the data type to marshal.
	Type reflect.Type
	// Size is the marshaled size of this type.
	Size int
}

// NewTypeMarshaler creates a *TypeMarshaler by a value.
// The value "zero" defines what type this Marshaler can deal with and must be a
// fixed size type.
func NewTypeMarshaler(zero interface{}) (*TypeMarshaler, error) {
	return NewTypeMarshalerEndian(zero, nil)
}

// NewTypeMarshalerEndian creates a *TypeMarshaler with a specified byte order.
//
// "endian" could be binary.LittleEndian or binary.BigEndian.
func NewTypeMarshalerEndian(zero interface{}, endian binary.ByteOrder) (*TypeMarshaler, error) {
	if endian == nil {
		endian = defaultEndian
	}
	m := &TypeMarshaler{
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

// NewTypeMarshalerEndianByType creates a *TypeMarshaler for specified type and with a specified byte order.
//
// "endian" could be binary.LittleEndian or binary.BigEndian.
func NewTypeMarshalerEndianByType(t reflect.Type, endian binary.ByteOrder) (*TypeMarshaler, error) {
	v := reflect.New(t)
	return NewTypeMarshalerEndian(v.Interface(), endian)
}

// Marshal converts a m.Type value to byte slice.
// If a different type value from the one used with NewTypeMarshaler passed in,
// it panics.
func (m *TypeMarshaler) Marshal(d interface{}) []byte {
	if reflect.Indirect(reflect.ValueOf(d)).Type() != m.Type {
		panic("different type from TypeMarshaler.Type")
	}

	b := bytes.NewBuffer(make([]byte, 0, m.Size))
	err := binary.Write(b, m.Endian, d)
	if err != nil {
		// there should not be any error if type is fixed size
		panic(err)
	}
	return b.Bytes()
}

// Unmarshal converts byte slice to a pointer to Type value.
// It returns number bytes consumed and an Type value in interface{}.
func (m *TypeMarshaler) Unmarshal(b []byte) (int, interface{}) {

	b = b[0:m.Size]
	v := reflect.New(m.Type)
	err := binary.Read(bytes.NewBuffer(b), m.Endian, v.Interface())
	if err != nil {
		panic(err)
	}
	return m.Size, reflect.Indirect(v).Interface()
}

// GetSize returns m.Size.
func (m *TypeMarshaler) GetSize(d interface{}) int {
	return m.Size
}

// GetMarshaledSize returns m.Size.
func (m *TypeMarshaler) GetMarshaledSize(b []byte) int {
	return m.Size
}
