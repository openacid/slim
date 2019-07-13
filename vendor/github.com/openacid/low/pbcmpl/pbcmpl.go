// Package pbcmpl provides a header for proto.Message
//
// One of the known protobuf issue is that user must control the size when
// unmarshaling.
// This package gives a solution to add a header for every proto.Message, in which there
// are: a semantic-version for checking compatibility, a header size and size of
// marshaled proto.Message.
//
// Since 0.1.6
package pbcmpl

import (
	"io"

	"github.com/openacid/errors"

	proto "github.com/golang/protobuf/proto"
)

// VersionedMessage must provide a "GetVersion" returning the version of a
// proto.Message .
// The version is a string in 16 bytes.
//
// Since 0.1.6
type VersionedMessage interface {
	GetVersion() string
	proto.Message
}

func marshal(msg proto.Message, ver string) ([]byte, []byte, error) {

	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, nil, err
	}

	h := newHeader(ver, uint64(len(data)))

	// header is tested and should never encounter error when marshaling it
	header, _ := proto.Marshal(h)

	return header, data, nil
}

// Marshal a proto.Message following a small header into an io.Writer.
// The header contains version(if msg is a VersionedMessage), header size and
// payload size.
//
// It returns the number of written bytes and an error.
//
// Since v0.1.6
func Marshal(w io.Writer, msg proto.Message) (int64, error) {

	ver := DefaultVer
	vmsg, ok := msg.(VersionedMessage)
	if ok {
		ver = vmsg.GetVersion()
	}

	h, d, err := marshal(msg, ver)
	if err != nil {
		return 0, err
	}

	n, err := w.Write(h)
	if err != nil {
		return int64(n), err
	}

	n2, err := w.Write(d)
	n += n2
	if err != nil {
		return int64(n), err
	}

	return int64(n), nil
}

// ReadHeader reads header info from a stream marshaled by this module.
// It returns number of bytes it has read, a Header interface for retreiving
// header info and an error.
// The number of bytes may be greater than 0 even if there is an error.
//
// Since 0.1.7
func ReadHeader(r io.Reader) (int64, Header, error) {

	b := make([]byte, fixedSize)

	// io.ReadFull returns err:
	//     EOF:              means n = 0
	//     ErrUnexpectedEOF: means n < len(buf)  underlaying Reader returns EOF
	//     nil:              means n == len(buf)
	n, err := io.ReadFull(r, b)
	if err != nil {
		return int64(n), nil, errors.WithStack(err)
	}

	h := &header{}
	// by tests, Unmarshaling a header should never return an error
	proto.Unmarshal(b, h)

	return int64(n), &headerInfo{h}, nil
}

// Unmarshal a header and following message form an io.Reader .
// It returns the number of read bytes, a version in string and an error.
//
// Since 0.1.6
func Unmarshal(r io.Reader, msg proto.Message) (int64, string, error) {

	n, hi, err := ReadHeader(r)
	if err != nil {
		return n, "", err
	}

	ver := hi.GetVersion()

	if hi.GetHeaderSize() != int64(fixedSize) {
		return n, ver, errors.WithStack(ErrInvalidHeaderSize)
	}

	b := make([]byte, hi.GetBodySize())
	nbody, err := io.ReadFull(r, b)
	n += int64(nbody)
	if err != nil {
		return n, ver, errors.WithStack(err)
	}

	err = proto.Unmarshal(b, msg)
	return n, ver, errors.WithStack(err)
}

// HeaderSize returns the marshaled size of the header for a proto.Message .
//
// Since 0.1.6
func HeaderSize(msg proto.Message) int {
	return fixedSize
}

// Size returns the size of the header and the marshaled message.
//
// Since 0.1.6
func Size(msg proto.Message) int {
	return HeaderSize(msg) + proto.Size(msg)
}

func verStr(buf []byte) string {
	var i int
	for i = len(buf) - 1; i >= 0 && buf[i] == 0; i-- {
	}
	return string(buf[:i+1])
}
