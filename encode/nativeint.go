package encode

import (
	"encoding/binary"
	"math/bits"
)

// Int converts int to slice of bytes and back.
type Int struct{}

// Encode converts int to slice of bytes.
func (c Int) Encode(d interface{}) []byte {
	size := bits.UintSize / 8
	b := make([]byte, size)
	v := d.(int)
	if size == 4 {
		binary.LittleEndian.PutUint32(b, uint32(v))
	} else if size == 8 {
		binary.LittleEndian.PutUint64(b, uint64(v))
	} else {
		panic("unknown int size")
	}

	return b
}

// Decode converts slice of bytes to int.
// It returns number bytes consumed and an int.
func (c Int) Decode(b []byte) (int, interface{}) {

	size := bits.UintSize / 8
	s := b[:size]

	var d int
	if size == 4 {
		d = int(binary.LittleEndian.Uint32(s))
	} else if size == 8 {
		d = int(binary.LittleEndian.Uint64(s))
	} else {
		panic("unknown int size")
	}
	return size, d
}

// GetSize returns native int size in byte after encoding v.
func (c Int) GetSize(d interface{}) int {
	return bits.UintSize / 8
}

// GetEncodedSize returns native int size.
func (c Int) GetEncodedSize(b []byte) int {
	return bits.UintSize / 8
}
