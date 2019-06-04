// Package serialize provide general serialize and de-serialize functions.
//
// It adds a header for every object to serialize. The header contains version,
// header length and data length.
// Thus we de-serializing, total size can be accquired at once, without reading
// the entire data.
//
// Deprecated: use github.com/openacid/low/pbcmpl; This module will be removed
// in 1.0.0 .
package serialize

import (
	"bytes"
	"encoding/binary"
	"io"
	"unsafe"

	"github.com/openacid/low/iohelper"
	"github.com/openacid/slim/version"

	proto "github.com/golang/protobuf/proto"
)

const (
	// MaxMarshalledSize defines the max marshaled size for an object.
	MaxMarshalledSize int64 = 1024 * 1024 * 1024
)

// DataHeader defines the header format of a serialized byte stream.
//
// It contains version, header size(size of this struct) and data size(size of the user data).
//
// To ensure Compatibility:
//
//    - do NOT change type of fields
//    - do NOT reuse any ever existing names
//    - do NOT adjust fields order
//    - only append fields
//      - only use fixed-size type, e.g. not int, use int32 or int64
//      - test Every version of dataHeader ever existed
//
type DataHeader struct {
	// Version of this serialized data, for compatibility check.
	// It is in form of <major>.<minor>.<release>;.
	// As long as version is a string, its max size is 16.
	//
	// See: https://semver.org/
	Version [version.MAXLEN]byte

	// HeaderSize is the serialized size in byte for DataHeader.
	HeaderSize uint64

	// DataSize is the size in byte of user data.
	DataSize uint64
}

func bytesToString(buf []byte, delimter byte) string {
	delimPos := bytes.IndexByte(buf, delimter)
	if delimPos == -1 {
		delimPos = len(buf)
	}

	return string(buf[:delimPos])
}

func makeDataHeader(verStr string, headerSize uint64, dataSize uint64) *DataHeader {
	if len(verStr) >= version.MAXLEN {
		panic("version length overflow")
	}

	if verStr > version.VERSION {
		panic("forward compatibility is not supported")
	}

	header := DataHeader{
		HeaderSize: headerSize,
		DataSize:   dataSize,
	}

	copy(header.Version[:], verStr)

	return &header
}

func makeDefaultDataHeader(dataSize uint64) *DataHeader {
	headerSize := GetMarshalHeaderSize()

	return makeDataHeader(version.VERSION, uint64(headerSize), dataSize)
}

// UnmarshalHeader reads just enough bytes from reader and load the data into a
// DataHeader object.
func UnmarshalHeader(reader io.Reader) (header *DataHeader, err error) {
	verBuf := make([]byte, version.MAXLEN)

	// io.ReadFull returns err:
	//     EOF:              means n = 0
	//     ErrUnexpectedEOF: means n < len(buf)  underlaying Reader returns EOF
	//     nil:              means n == len(buf)
	if _, err := io.ReadFull(reader, verBuf); err != nil {
		return nil, err
	}

	verStr := bytesToString(verBuf, 0)

	var headerSize uint64
	err = binary.Read(reader, binary.LittleEndian, &headerSize)
	if err != nil {
		return nil, err
	}

	toRead := headerSize - version.MAXLEN - uint64(unsafe.Sizeof(headerSize))
	buf := make([]byte, toRead)

	if _, err := io.ReadFull(reader, buf); err != nil {
		return nil, err
	}

	var dataSize uint64
	restReader := bytes.NewReader(buf)
	err = binary.Read(restReader, binary.LittleEndian, &dataSize)
	if err != nil {
		return nil, err
	}

	return makeDataHeader(verStr, headerSize, dataSize), nil
}

func marshalHeader(writer io.Writer, header *DataHeader) (err error) {
	return binary.Write(writer, binary.LittleEndian, header)
}

// Marshal serializes a protobuf object into a io.Writer .
//
// It returns number of bytes actually written, and encountered error.
//
// The content written to writer may be wrong if there were error during Marshal().
// So make a temp copy, and copy it to destination if everything is ok.
func Marshal(writer io.Writer, obj proto.Message) (cnt int64, err error) {
	marshaledData, err := proto.Marshal(obj)
	if err != nil {
		return 0, err
	}

	dataSize := uint64(len(marshaledData))
	dataHeader := makeDefaultDataHeader(dataSize)

	// write to headerBuf to get cnt
	headerBuf := new(bytes.Buffer)
	err = marshalHeader(headerBuf, dataHeader)
	if err != nil {
		return 0, err
	}

	nHeader, err := writer.Write(headerBuf.Bytes())
	if err != nil {
		return int64(nHeader), err
	}

	nData, err := writer.Write(marshaledData)

	return int64(nHeader + nData), err
}

// MarshalAt is similar to Marshal except it writes data into io.WriterAt
// interface instead of io.Writer .
func MarshalAt(writer io.WriterAt, offset int64, obj proto.Message) (cnt int64, err error) {

	w := iohelper.NewSectionWriter(writer, offset, MaxMarshalledSize)
	return Marshal(w, obj)
}

// Unmarshal deserialize data from an io.Reader and load the data into a
// protobuf object.
//
// One must specifies the type of the object before Unmarshal it.
// TODO: return the number of byte read. Since all other [un]marhshal[at]
// functions return the number of byte written or read.
func Unmarshal(reader io.Reader, obj proto.Message) (err error) {
	dataHeader, err := UnmarshalHeader(reader)
	if err != nil {
		return err
	}

	dataBuf := make([]byte, dataHeader.DataSize)

	// Repeat reader.Read until encounting an error or read full
	//
	// io.Reader:Read() does not guarantee to read all
	// len(dataBuf)
	if _, err := io.ReadFull(reader, dataBuf); err != nil {
		return err
	}

	return proto.Unmarshal(dataBuf, obj)
}

// UnmarshalAt is similar to Unmarshal except it reads from io.ReaderAt thus it
// is able to specify where to start to read.
func UnmarshalAt(reader io.ReaderAt, offset int64, obj proto.Message) (n int64, err error) {

	// Wrap io.ReaderAt with a offset-self-maintained io.Reader
	// The 3rd argument specifies right boundary. It is not buffer size related
	// thus we just give it a big enough value.
	r := io.NewSectionReader(reader, offset, MaxMarshalledSize)

	err = Unmarshal(r, obj)
	n, seekErr := r.Seek(0, io.SeekCurrent)
	if seekErr != nil {
		// It must be a programming error.
		// seekErr is not nil only when:
		// - whence is invalid
		// - or return value would be a negative int.
		panic("seekErr must be nil")
	}
	return n, err

}

// GetMarshalHeaderSize returns the serialized size of a DataHeader struct.
func GetMarshalHeaderSize() int64 {
	return int64(unsafe.Sizeof(uint64(0))*2 + version.MAXLEN)
}

// GetMarshalSize calculates the total size for a serialized object.
func GetMarshalSize(obj proto.Message) int64 {
	return GetMarshalHeaderSize() + int64(proto.Size(obj))
}
