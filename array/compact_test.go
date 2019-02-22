package array

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"reflect"
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
