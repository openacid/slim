package serialize

import (
	"bytes"
	"encoding/binary"
	"math"
	"os"
	"testing"
	"xec/sparse"
	"xec/version"
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

	if header.dataSize != dataSize {
		t.Fatalf("wrong data size")
	}

	if header.headerSize != headerSize {
		t.Fatalf("wrong header size")
	}

	verStr, err := bytesToString(header.version[:], 0)
	if err != nil || verStr != ver {
		t.Fatalf("wrong version: %v, %s, expect: %s", err, verStr, ver)
	}

	header = makeDefaultDataHeader(dataSize)
	if header.dataSize != dataSize {
		t.Fatalf("wrong data size")
	}

	// sizeof(uint64) * 2 + version.MAXLEN
	if header.headerSize != 32 {
		t.Fatalf("wrong header size: %v", header.headerSize)
	}

	if len(header.version) != version.MAXLEN {
		t.Fatalf("wrong version length: %v", len(header.version))
	}

	verStr, err = bytesToString(header.version[:], 0)
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

	t.Logf("header to marshal: %v", sHeader)
	writer.Close()

	// unmarshal
	reader, err := os.OpenFile(testDataFn, os.O_RDONLY, 0755)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	defer reader.Close()

	rSHeader, err := unmarshalHeader(reader)
	if err != nil {
		t.Fatalf("failed to unmarshalHeader: %v", err)
	}

	t.Logf("header unmarshaled: %v", rSHeader)

	if rSHeader.dataSize != sHeader.dataSize {
		t.Fatalf("wrong data size: %v, %v", rSHeader.dataSize, sHeader.dataSize)
	}

	if rSHeader.headerSize != sHeader.headerSize {
		t.Fatalf("wrong header size: %v, %v",
			rSHeader.headerSize, sHeader.headerSize)
	}

	for idx, sByte := range sHeader.version {
		rByte := rSHeader.version[idx]
		if rByte != sByte {
			t.Fatalf("wrong byte in version: %v, %v, %v", idx, sByte, rByte)
		}
	}

	rVersion, err := bytesToString(rSHeader.version[:], 0)
	if err != nil {
		t.Fatalf("failed to restore version string")
	}

	if rVersion != version.VERSION {
		t.Fatalf("wrong version string: %v, %v", rVersion, version.VERSION)
	}
}

func TestMarshalUnMarshal(t *testing.T) {
	index := []uint32{10, 20, 30, 40, 50, 60}

	sArray := &sparse.SparseArray{EltConverter: sparse.U32Conv{}}
	err := sArray.Init(index, index)
	if err != nil {
		t.Fatalf("failed to init sparse array")
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
		t.Fatalf("failed to store sparse array: %v", err)
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

	rSArray := &sparse.SparseArray{EltConverter: sparse.U32Conv{}}

	err = Unmarshal(reader, rSArray)
	if err != nil {
		t.Fatalf("failed to load data: %v", err)
	}

	// check sparse array
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
