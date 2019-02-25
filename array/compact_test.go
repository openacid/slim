package array

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"runtime"
	"testing"
	"time"

	proto "github.com/golang/protobuf/proto"
)

func TestNewErrorArgments(t *testing.T) {
	var index []uint32
	eltsData := []uint32{12, 15, 19, 120, 300}

	var err error

	index = []uint32{1, 5, 9, 203}
	_, err = NewU32(index, eltsData)
	if err == nil {
		t.Fatalf("new with wrong index length must error")
	}

	index = []uint32{1, 5, 5, 203, 400}
	_, err = NewU32(index, eltsData)
	if err == nil {
		t.Fatalf("new with unsorted index must error")
	}
}

func TestNew(t *testing.T) {
	var cases = []struct {
		index    []uint32
		eltsData []uint32
	}{
		{
			[]uint32{}, []uint32{},
		},
		{
			[]uint32{0, 5, 9, 203, 400}, []uint32{12, 15, 19, 120, 300},
		},
	}

	for _, c := range cases {
		index, eltsData := c.index, c.eltsData
		cnt := uint32(len(index))

		ca, err := NewU32(index, eltsData)
		if err != nil {
			t.Fatalf("failed new compacted array, err: %s", err)
		}

		if ca.Cnt != cnt {
			t.Fatalf("cnt is not equal expect: %d, act: %d", cnt, ca.Cnt)
		}

		buf := new(bytes.Buffer)
		binary.Write(buf, binary.LittleEndian, eltsData)

		expElts := buf.Bytes()
		if expElts == nil {
			expElts = []byte{}
		}

		if !reflect.DeepEqual(ca.Elts, expElts) {
			t.Fatalf("elts is not equal expect: %d, act: %d", expElts, ca.Elts)
		}
	}

}

func TestGet(t *testing.T) {
	index, eltsData := []uint32{}, []uint32{}
	rnd := rand.New(rand.NewSource(time.Now().Unix()))

	keysMap := map[uint32]bool{}
	num, idx, cnt := uint32(0), uint32(0), uint32(1024)
	for {
		if rnd.Intn(2) == 1 {
			index = append(index, idx)
			eltsData = append(eltsData, rnd.Uint32())
			num++
			keysMap[idx] = true
		}
		idx++
		if num == cnt {
			break
		}
	}

	ca, err := NewU32(index, eltsData)
	if err != nil {
		t.Fatalf("failed new compacted array, err: %s", err)
	}

	dataIdx := uint32(0)
	for ii := uint32(0); ii < idx; ii++ {

		actByte := ca.Get(ii)
		if _, ok := keysMap[ii]; ok {
			act := actByte.(uint32)
			if eltsData[dataIdx] != act {
				t.Fatalf("Get i:%d is not equal expect: %d, act: %d", ii, eltsData[dataIdx], act)
			}
		} else {
			if actByte != nil {
				t.Fatalf("Get i:%d is not nil expect: nil, act:%d", ii, actByte)
			}
		}

		// test Get2

		actByte, found := ca.Get2(ii)
		_, present := keysMap[ii]
		if found != present {
			t.Fatalf("Get i:%d present:%t but:%t", ii, present, found)
		}

		if found {
			if eltsData[dataIdx] != actByte.(uint32) {
				t.Fatalf("Get i:%d is not equal expect: %d, act: %d", ii, eltsData[dataIdx], actByte.(uint32))
			}
		}

		if _, ok := keysMap[ii]; ok {
			dataIdx++
		}
	}
}

func TestSerialize(t *testing.T) {
	// TODO add Marshal function for package array, to Marshal multi versioned array

	// fmt.Printf("%#v\n",  data )
	serialized := []byte{
		0x8, 0x4, 0x12, 0x6, 0xa2, 0x4, 0x0, 0x0,
		0x80, 0x10, 0x1a, 0x4, 0x0, 0x0, 0x0, 0x3,
		0x22, 0x10, 0xc, 0x0, 0x0, 0x0, 0xf, 0x0,
		0x0, 0x0, 0x13, 0x0, 0x0, 0x0, 0x78, 0x0,
		0x0, 0x0,
	}

	index := []uint32{1, 5, 9, 203}
	eltsData := []uint32{12, 15, 19, 120}

	arr, err := NewU32(index, eltsData)
	if err != nil {
		t.Fatalf("create array failure: %s", err)
	}

	data, err := proto.Marshal(arr)
	if err != nil {
		t.Fatalf("proto.Marshal: %s", err)
	}

	if !reflect.DeepEqual(serialized, data) {
		fmt.Println(serialized)
		fmt.Println(data)
		t.Fatalf("serialized data incorrect")
	}

	loaded := &Array32{
		Converter: U32Conv{},
	}

	err = proto.Unmarshal(data, loaded)
	if err != nil {
		t.Fatalf("proto.Unmarshal: %s", err)
	}

	second, err := proto.Marshal(loaded)
	if err != nil {
		t.Fatalf("proto.Marshal: %s", err)
	}

	if !reflect.DeepEqual(serialized, second) {
		fmt.Println(serialized)
		fmt.Println(second)
		t.Fatalf("second serialized data incorrect")
	}
}

func BenchmarkInit(b *testing.B) {

	n := 10240
	index := make([]uint32, n)
	elts := make([]uint32, n)

	for i := 0; i < n; i++ {
		index = append(index, uint32(i))
		elts = append(elts, uint32(i))
	}

	for i := 0; i < b.N; i++ {
		NewU32(index, elts)
	}

}

func newByte(eSize int, index []uint32, elts [][]byte) (*Array32, error) {
	ca, err := New32(ByteConv{EltSize: eSize}, index, elts)
	return ca, err
}

func readRss() uint64 {
	var stats runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&stats)
	return stats.Alloc
}

func makeTestData(eltSize int, cnt uint32) [][]byte {
	eltsData := make([][]byte, cnt)

	for i := uint32(0); i < cnt; i++ {
		eltsData[i] = make([]byte, eltSize)
	}

	return eltsData
}

func makeTestIndex(maxIdx, idxDis uint32) []uint32 {
	index := make([]uint32, 0, maxIdx)

	for i := uint32(0); i < maxIdx; i++ {
		if i%idxDis == 0 {
			index = append(index, i)

		}
	}

	return index
}

func BenchmarkMemOverhead(b *testing.B) {
	var cases = []struct {
		eltSize int
		maxIdx  uint32
	}{
		{1, 1 << 16},
		{2, 1 << 16},
		{4, 1 << 16},
		{8, 1 << 16},
	}

	var sca []*Array32
	fmt.Printf("%-10s%-10s%-10s%-10s%-12s%-12s%-12s%-10s\n",
		"eltSize", "eltCount", "idxDis", "caCnt", "totalSize", "caAvgSize", "dataAvgSize", "Overhead")

	for _, c := range cases {
		eltSize, maxIdx := c.eltSize, c.maxIdx

		for i := uint32(1); i < 1<<16; i++ {
			idxDis := uint32(math.Pow((float64(1)+math.Sqrt(5))/2, float64(i)))
			if idxDis >= maxIdx {
				break
			}

			sca = []*Array32{}

			index := makeTestIndex(maxIdx, idxDis)
			eltCnt := uint32(len(index))
			elts := makeTestData(eltSize, eltCnt)

			rss1 := readRss()

			caCnt := 1024
			var ca *Array32
			for i := 0; i < caCnt; i++ {
				ca, _ = newByte(eltSize, index, elts)
				sca = append(sca, ca)
			}
			ca = nil

			rss2 := readRss()
			var _ []uint64 = sca[0].Bitmaps

			totalSize := rss2 - rss1
			dataAvgSize := uint64(eltSize) * uint64(eltCnt)
			caAvgSize := totalSize / uint64(caCnt)
			overhead := float64(caAvgSize) / float64(dataAvgSize)

			fmt.Printf("%-10d%-10d%-10d%-10d%-12d%-12d%-12d%-10.3f\n",
				eltSize, eltCnt, idxDis, caCnt, totalSize, caAvgSize, dataAvgSize, overhead)
		}
	}
}
