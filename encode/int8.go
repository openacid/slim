package encode

// I8 converts int8 to slice of 1 byte and back.
type I8 struct{}

// Encode converts int8 to slice of 1 byte.
func (c I8) Encode(d interface{}) []byte {
	return []byte{byte(d.(int8))}
}

// Decode converts slice of 1 byte to int8.
// It returns number bytes consumed and an int8.
func (c I8) Decode(b []byte) (int, interface{}) {
	return 1, int8(b[0])
}

// GetSize returns the size in byte after encoding v.
func (c I8) GetSize(d interface{}) int {
	return 1
}

// GetEncodedSize returns 2.
func (c I8) GetEncodedSize(b []byte) int {
	return 1
}
