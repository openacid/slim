package array

import (
	"bytes"
	"encoding/binary"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/openacid/slim/prototype"
)

type storage interface {
	GetStorage() *prototype.Array32Storage
}

// TODO: proto.Message.Size() might report wrong size because user-data is not
// stored in proto.Message.

// getUserEltsField returns `arr.Data`.
// `arr.Data` is recognized by Marshal() and Unmarshal() as a user data
// container.
// It must be a slice of fixed size type.
// Fixed size type means its space can be determined by its type.
// Fixed size type: `struct {int32, uint32}`
// Non-fixed size type: `struct {[]int32, uint32}`.
func getUserEltsField(arr interface{}) reflect.Value {
	v := reflect.Indirect(reflect.ValueOf(arr))
	return v.FieldByName("Data")
}

// Marshal array `arr`.
// `arr` can be either with or without user-defined field `Data`.
// If there is a `arr.Data` field, it first serialize `arr.Data` into the
// underlaying field `Elts`. Then marshal `arr` except `Data`.
func Marshal(arr storage) ([]byte, error) {

	var bb []byte
	var err error
	sto := arr.GetStorage()

	f := getUserEltsField(arr)
	if f.IsValid() {
		uelts := f.Interface()

		b := &bytes.Buffer{}
		err = binary.Write(b, binary.LittleEndian, uelts)
		if err != nil {
			return nil, err
		}
		sto.Elts = b.Bytes()

		bb, err = proto.Marshal(sto)

		// sto.Elts is only a intermedia space for marshal/unmarshal.
		// For a array with user-defined elements, clear it to save space.
		sto.Elts = nil
	} else {
		bb, err = proto.Marshal(sto)
	}

	return bb, err
}

// Unmarshal array `arr` from `raw`.
// `arr` can be either with or without user-defined field `Data`.
// If there is a `arr.Data` field, it fills up `arr.Data`.
func Unmarshal(arr storage, raw []byte) (int, error) {

	sto := arr.GetStorage()

	err := proto.Unmarshal(raw, sto)
	if err != nil {
		return 0, err
	}

	nread := proto.Size(sto)

	f := getUserEltsField(arr)

	if f.IsValid() {

		n := int(sto.Cnt)
		emp := reflect.MakeSlice(f.Type(), n, n)
		f.Set(emp)

		uelts := f.Interface()

		r := bytes.NewBuffer(sto.Elts)
		err = binary.Read(r, binary.LittleEndian, uelts)
		sto.Elts = nil
	}

	return nread, err
}
