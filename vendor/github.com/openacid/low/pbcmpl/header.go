package pbcmpl

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	versionLen = 16

	// DefaultVer is the default version if the message to marshal does not
	// provide version.
	//
	// Since v0.1.6
	DefaultVer = "1.0.0"
)

var (
	fixedSize = binary.Size(&header{})
	endian    = binary.LittleEndian
)

// Header defines header info retrieving APIs.
//
// Since 0.1.7
type Header interface {
	GetVersion() string
	GetHeaderSize() int64
	GetBodySize() int64
}

type headerInfo struct {
	*header
}

func (hi *headerInfo) GetVersion() string {
	return verStr(hi.Version[:])
}

func (hi *headerInfo) GetHeaderSize() int64 {
	return int64(hi.HeaderSize)
}

func (hi *headerInfo) GetBodySize() int64 {
	return int64(hi.BodySize)
}

// header is a fixed-size structure that defines the header format of a
// marshaled byte stream.
//
// It contains version, header size(size of this struct) and data size(size of
// the user data).
//
// NEVER change this structure.
//
// Since v0.1.6
type header struct {
	// version of this marshaled data, for compatibility check.
	// It is a semantic version in form of <major>.<minor>.<release>;.
	// As long as version is a string, its max size is 16.
	//
	// See: https://semver.org/
	Version [versionLen]byte

	// HeaderSize is the marshaled size in byte for header.
	HeaderSize uint64

	// BodySize is the size in byte of user data.
	BodySize uint64
}

// Marshal header implement proto.Message
//
// Since 0.1.6
func (h *header) Marshal() ([]byte, error) {
	b := &bytes.Buffer{}
	err := binary.Write(b, endian, h)
	return b.Bytes(), err
}

// Unmarshal header implement proto.Message
//
// Since 0.1.6
func (h *header) Unmarshal(buf []byte) error {
	r := bytes.NewReader(buf)
	return binary.Read(r, endian, h)
}

// Reset header implement proto.Message
//
// Since 0.1.6
func (h *header) Reset() {
	*h = header{}
}

// String header implement proto.Message
//
// Since 0.1.6
func (h *header) String() string {
	return fmt.Sprintf("%s %d %d",
		verStr(h.Version[:]),
		h.HeaderSize,
		h.BodySize)
}

// ProtoMessage header implement proto.Message
//
// Since 0.1.6
func (h *header) ProtoMessage() {}

func newHeader(ver string, bodysize uint64) *header {

	if len(ver) > versionLen {
		panic("version length overflow")
	}

	h := &header{
		HeaderSize: uint64(fixedSize),
		BodySize:   bodysize,
	}
	copy(h.Version[:], ver)
	return h
}
