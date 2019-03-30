package encode

// Bytes converts a byte slice into fixed length slice.
// Result slice length is defined by Bytes.Size .
type Bytes struct {
	Size int
}

// Encode converts byte slice to byte slice.
func (c Bytes) Encode(d interface{}) []byte {
	return d.([]byte)
}

// Decode copies fixed length slice out of source byte slice.
// The returned bytes are NOT copied.
func (c Bytes) Decode(b []byte) (int, interface{}) {
	s := b[:c.Size]
	return c.Size, s
}

// GetSize returns the length: c.Size.
func (c Bytes) GetSize(d interface{}) int {
	return c.Size
}

// GetEncodedSize returns c.Size
func (c Bytes) GetEncodedSize(b []byte) int {
	return c.Size
}
