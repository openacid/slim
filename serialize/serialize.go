package serialize

import (
	"bytes"
	"encoding/binary"
	proto "github.com/golang/protobuf/proto"
	"io"
	"unsafe"
	"xec/version"
)

/**
 * Compatiblity gurantee:
 *     - do NOT change type of fields
 *     - do NOT reuse any ever existing names
 *     - do NOT adjust fields order
 *     - only append fields
 *	   - only use fixed-size type, e.g. not int, use int32 or int64
 *	   - test Every version of dataHeader ever existed
 */
type DataHeader struct {
	Version    [version.MAXLEN]byte // version.VERSION, major.minor.release
	HeaderSize uint64               // the length in bytes of dataHeader size
	DataSize   uint64               // the length in bytes of serialized data size
}

func bytesToString(buf []byte, delimter byte) (str string, err error) {
	delimPos := bytes.IndexByte(buf, delimter)
	if delimPos == -1 {
		delimPos = len(buf)
	}

	return string(buf[:delimPos]), err
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

func readFull(reader io.Reader, buf []byte) (cnt int, err error) {
	n, cnt, toRead := 0, 0, len(buf)

	for n < toRead {
		cnt, err = reader.Read(buf[n:])
		n += cnt

		if err != nil {
			break
		}
	}

	return n, err
}

func UnmarshalHeader(reader io.Reader) (header *DataHeader, err error) {
	verBuf := make([]byte, version.MAXLEN)

	if _, err := readFull(reader, verBuf); err != nil {
		return nil, err
	}

	verStr, err := bytesToString(verBuf, 0)
	if err != nil {
		return nil, err
	}

	var headerSize uint64
	err = binary.Read(reader, binary.LittleEndian, &headerSize)
	if err != nil {
		return nil, err
	}

	toRead := headerSize - version.MAXLEN - uint64(unsafe.Sizeof(headerSize))
	buf := make([]byte, toRead)

	if _, err := readFull(reader, buf); err != nil {
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

/**
 * the content written to writer may be wrong if there were error during Marshal()
 * So make a temp copy, and copy it to destination if everything is ok
 */
func Marshal(writer io.Writer, obj proto.Message) (cnt uint64, err error) {
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
		return uint64(nHeader), err
	}

	nData, err := writer.Write(marshaledData)

	return uint64(nHeader + nData), err
}

func Unmarshal(reader io.Reader, obj proto.Message) (err error) {
	dataHeader, err := UnmarshalHeader(reader)
	if err != nil {
		return err
	}

	dataBuf := make([]byte, dataHeader.DataSize)

	if _, err := reader.Read(dataBuf); err != nil {
		return err
	}

	if err := proto.Unmarshal(dataBuf, obj); err != nil {
		return err
	}

	return nil
}

func GetMarshalHeaderSize() int64 {
	return int64(unsafe.Sizeof(uint64(0))*2 + version.MAXLEN)
}

func GetMarshalSize(obj proto.Message) int64 {
	return GetMarshalHeaderSize() + int64(proto.Size(obj))
}
