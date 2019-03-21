package marshal

// Bytes converts a byte slice into fixed length slice.
// Result slice length is defined by Bytes.Size .
type Bytes struct {
	Size int
}

// Marshal converts byte slice to byte slice.
func (c Bytes) Marshal(d interface{}) []byte {
	return d.([]byte)
}

// Unmarshal copies fixed length slice out of source byte slice.
// The returned bytes are NOT copied.
func (c Bytes) Unmarshal(b []byte) (int, interface{}) {
	s := b[:c.Size]
	return c.Size, s
}

// GetSize returns the length: c.Size.
func (c Bytes) GetSize(d interface{}) int {
	return c.Size
}

// GetMarshaledSize returns c.Size
func (c Bytes) GetMarshaledSize(b []byte) int {
	return c.Size
}
