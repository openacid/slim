package serialize

import (
	"bytes"
	"encoding/binary"
	"math"
	"os"
	"testing"

	"github.com/openacid/slim/array"
	"github.com/openacid/slim/version"
)

var testDataFn = "data"

func TestBinaryWriteUint64Length(t *testing.T) {
	expectSize := 8

	for _, u64 := range []uint64{0, math.MaxUint64 / 2, math.MaxUint64} {
		writer := new(bytes.Buffer)

		err := binary.Write(writer, binary.LittleEndian, u64)
		if err != nil {
			t.Fatalf("failed to write uint64: %v", err)
		}

		buf := writer.Bytes()

		if len(buf) != expectSize {
			t.Fatalf("uint64 does not take %d bytes", expectSize)
		}
	}
}

func TestBytesToString(t *testing.T) {
	str, err := bytesToString(nil, 0)
	if err != nil || str != "" {
		t.Fatalf("failed to handle nil: %s, %v", str, err)
	}

	str, err = bytesToString([]byte{}, 0)
	if err != nil || str != "" {
		t.Fatalf("failed to handle nil: %s, %v", str, err)
	}

	str, err = bytesToString([]byte{'a', 'b', 'c'}, 0)
	if err != nil || str != "abc" || len(str) != 3 {
		t.Fatalf("failed to handle abc: %s, %v", str, err)
	}

	str, err = bytesToString([]byte{'1', '.', '0', '.', '0', 0}, 0)
	if err != nil || str != "1.0.0" || len(str) != 5 {
		t.Fatalf("failed to handle 1.0.0'0': %s, %v", str, err)
	}

	bBuf := []byte{'1', '.', '0', '.', '0', 0}
	str, _ = bytesToString(bBuf, 0)

	bBuf[0] = '2'
	if str != "1.0.0" {
		t.Fatalf("wrong str value after modify byte buffer: %s", str)
	}
}

func TestMakeDataHeader(t *testing.T) {
	ver := "0.0.1"
	dataSize := uint64(1000)
	headerSize := uint64(100)
	header := makeDataHeader(ver, headerSize, dataSize)

	if header.DataSize != dataSize {
		t.Fatalf("wrong data size")
	}

	if header.HeaderSize != headerSize {
		t.Fatalf("wrong header size")
	}

	verStr, err := bytesToString(header.Version[:], 0)
	if err != nil || verStr != ver {
		t.Fatalf("wrong version: %v, %s, expect: %s", err, verStr, ver)
	}

	header = makeDefaultDataHeader(dataSize)
	if header.DataSize != dataSize {
		t.Fatalf("wrong data size")
	}

	// sizeof(uint64) * 2 + version.MAXLEN
	if header.HeaderSize != 32 {
		t.Fatalf("wrong header size: %v", header.HeaderSize)
	}

	if len(header.Version) != version.MAXLEN {
		t.Fatalf("wrong version length: %v", len(header.Version))
	}

	verStr, err = bytesToString(header.Version[:], 0)
	if err != nil || verStr != version.VERSION {
		t.Fatalf("wrong version: %v, %s, expect: %s", err, verStr, version.VERSION)
	}
}

func TestMarshalUnMarshalHeader(t *testing.T) {
	// marshal
	wOFlags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	writer, err := os.OpenFile(testDataFn, wOFlags, 0755)
	if err != nil {
		t.Fatalf("failed to create file: %v", err)
	}
	defer os.Remove(testDataFn)

	sHeader := makeDefaultDataHeader(1000)

	gHeaderSize := GetMarshalHeaderSize()
	if gHeaderSize != 32 {
		t.Fatalf("wrong header size: 32, %d", gHeaderSize)
	}

	err = marshalHeader(writer, sHeader)
	if err != nil {
		t.Fatalf("failed to marshalHeader: %v", err)
	}

	writer.Close()

	// unmarshal
	reader, err := os.OpenFile(testDataFn, os.O_RDONLY, 0755)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	defer reader.Close()

	rSHeader, err := UnmarshalHeader(reader)
	if err != nil {
		t.Fatalf("failed to unmarshalHeader: %v", err)
	}

	if rSHeader.DataSize != sHeader.DataSize {
		t.Fatalf("wrong data size: %v, %v", rSHeader.DataSize, sHeader.DataSize)
	}

	if rSHeader.HeaderSize != sHeader.HeaderSize {
		t.Fatalf("wrong header size: %v, %v",
			rSHeader.HeaderSize, sHeader.HeaderSize)
	}

	for idx, sByte := range sHeader.Version {
		rByte := rSHeader.Version[idx]
		if rByte != sByte {
			t.Fatalf("wrong byte in version: %v, %v, %v", idx, sByte, rByte)
		}
	}

	rVersion, err := bytesToString(rSHeader.Version[:], 0)
	if err != nil {
		t.Fatalf("failed to restore version string")
	}

	if rVersion != version.VERSION {
		t.Fatalf("wrong version string: %v, %v", rVersion, version.VERSION)
	}
}

func TestMarshalUnMarshal(t *testing.T) {
	index := []uint32{10, 20, 30, 40, 50, 60}

	sArray := &array.CompactedArray{EltConverter: array.U32Conv{}}
	err := sArray.Init(index, index)
	if err != nil {
		t.Fatalf("failed to init compacted array")
	}

	marshalSize := GetMarshalSize(sArray)

	// marshal
	wOFlags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	writer, err := os.OpenFile(testDataFn, wOFlags, 0755)
	if err != nil {
		t.Fatalf("failed to create file: %s, %v", testDataFn, err)
	}
	defer os.Remove(testDataFn)

	cnt, err := Marshal(writer, sArray)
	if err != nil {
		t.Fatalf("failed to store compacted array: %v", err)
	}

	writer.Close()

	fInfo, err := os.Stat(testDataFn)
	if err != nil {
		t.Fatalf("failed to get file info: %s, %v", testDataFn, err)
	}

	if fInfo.Size() != int64(cnt) {
		t.Fatalf("wrong file size: %d, %d", fInfo.Size(), cnt)
	}

	if fInfo.Size() != marshalSize {
		t.Fatalf("wrong marshal size: %d, %d", fInfo.Size(), marshalSize)
	}

	// unmarshal
	reader, err := os.OpenFile(testDataFn, os.O_RDONLY, 0755)
	if err != nil {
		t.Fatalf("failed to read file: %s, %v", testDataFn, err)
	}
	defer reader.Close()

	rSArray := &array.CompactedArray{EltConverter: array.U32Conv{}}

	err = Unmarshal(reader, rSArray)
	if err != nil {
		t.Fatalf("failed to load data: %v", err)
	}

	// check compacted array
	checkCompactedArray(index, rSArray, sArray, t)
}

func TestMarshalAtUnMarshalAt(t *testing.T) {
	index1 := []uint32{10, 20, 30, 40, 50, 60}
	index2 := []uint32{15, 25, 35, 45, 55, 65}

	sArray1 := &array.CompactedArray{EltConverter: array.U32Conv{}}
	err := sArray1.Init(index1, index1)
	if err != nil {
		t.Fatalf("failed to init compacted array")
	}

	sArray2 := &array.CompactedArray{EltConverter: array.U32Conv{}}
	err = sArray2.Init(index2, index2)
	if err != nil {
		t.Fatalf("failed to init compacted array")
	}

	// marshalat
	wOFlags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	writer, err := os.OpenFile(testDataFn, wOFlags, 0755)
	if err != nil {
		t.Fatalf("failed to create file: %s, %v", testDataFn, err)
	}
	defer os.Remove(testDataFn)

	offset1 := int64(0)
	cnt, err := MarshalAt(writer, offset1, sArray1)
	if err != nil {
		t.Fatalf("failed to store compacted array: %v", err)
	}

	offset2 := offset1 + cnt
	_, err = MarshalAt(writer, offset2, sArray2)
	if err != nil {
		t.Fatalf("failed to store compacted array: %v", err)
	}

	writer.Close()

	// unmarshalat
	reader, err := os.OpenFile(testDataFn, os.O_RDONLY, 0755)
	if err != nil {
		t.Fatalf("failed to read file: %s, %v", testDataFn, err)
	}
	defer reader.Close()

	rSArray1 := &array.CompactedArray{EltConverter: array.U32Conv{}}
	_, err = UnmarshalAt(reader, offset1, rSArray1)
	if err != nil {
		t.Fatalf("failed to load data: %v", err)
	}

	checkCompactedArray(index1, rSArray1, sArray1, t)

	rSArray2 := &array.CompactedArray{EltConverter: array.U32Conv{}}
	_, err = UnmarshalAt(reader, offset2, rSArray2)
	if err != nil {
		t.Fatalf("failed to load data: %v", err)
	}

	checkCompactedArray(index2, rSArray2, sArray2, t)
}

func checkCompactedArray(index []uint32, rSArray, sArray *array.CompactedArray, t *testing.T) {
	if rSArray.Cnt != sArray.Cnt {
		t.Fatalf("wrong Cnt: %d, %d", rSArray.Cnt, sArray.Cnt)
	}

	if len(sArray.Bitmaps) != len(rSArray.Bitmaps) {
		t.Fatalf("wrong bitmap len: %d, %d", len(rSArray.Bitmaps), len(sArray.Bitmaps))
	}

	for idx, elt := range sArray.Bitmaps {
		if rSArray.Bitmaps[idx] != elt {
			t.Fatalf("wrong bitmap value: %v, %v", rSArray.Bitmaps[idx], elt)
		}
	}

	if len(sArray.Offsets) != len(rSArray.Offsets) {
		t.Fatalf("wrong offset len: %v, %v", rSArray.Offsets, sArray.Offsets)
	}

	for idx, elt := range sArray.Offsets {
		if rSArray.Offsets[idx] != elt {
			t.Fatalf("wrong offsets value: %v, %v", rSArray.Offsets[idx], elt)
		}
	}

	if len(sArray.Elts) != len(rSArray.Elts) {
		t.Fatalf("wrong Elts len: %v, %v", rSArray.Elts, sArray.Elts)
	}

	for _, idx := range index {
		sVal := sArray.Get(idx).(uint32)
		rsVal := rSArray.Get(idx).(uint32)

		if sVal != rsVal || sVal != idx {
			t.Fatalf("wrong Elts value: %v, %v, %v", sVal, rsVal, idx)
		}
	}
}

type testWriterReader struct {
	b [512]byte
}

func (t *testWriterReader) WriteAt(b []byte, off int64) (n int, err error) {
	length := len(b)
	for i := 0; i < length; i++ {
		t.b[int64(i)+off] = b[i]
	}
	return length, nil
}

func (t *testWriterReader) ReadAt(b []byte, off int64) (n int, err error) {
	length := len(b)
	copy(b, t.b[off:off+int64(length)])
	return length, nil
}

func TestWriteAtReadAt(t *testing.T) {
	rw := &testWriterReader{}
	index1 := []uint32{10, 20, 30, 40, 50, 60}

	wArr := &array.CompactedArray{EltConverter: array.U32Conv{}}
	err := wArr.Init(index1, index1)
	if err != nil {
		t.Fatalf("failed to init compacted array")
	}

	_, err = MarshalAt(rw, 10, wArr)
	if err != nil {
		t.Fatalf("failed to store compacted array: %v", err)
	}

	rArr := &array.CompactedArray{EltConverter: array.U32Conv{}}
	_, err = UnmarshalAt(rw, 10, rArr)
	if err != nil {
		t.Fatalf("failed to load data: %v", err)
	}

	checkCompactedArray(index1, rArr, wArr, t)
}
