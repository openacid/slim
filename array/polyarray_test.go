package array

import (
	"fmt"
	"testing"

	proto "github.com/golang/protobuf/proto"
	"github.com/openacid/slim/benchhelper"
	"github.com/stretchr/testify/assert"
)

var polyTestNums []int32 = []int32{
	0, 16, 32, 48, 64, 79, 95, 111, 126, 142, 158, 174, 190, 206, 222, 236,
	252, 268, 275, 278, 281, 283, 285, 289, 296, 301, 304, 307, 311, 313, 318,
	321, 325, 328, 335, 339, 344, 348, 353, 357, 360, 364, 369, 372, 377, 383,
	387, 393, 399, 404, 407, 410, 415, 418, 420, 422, 426, 430, 434, 439, 444,
	446, 448, 451, 456, 459, 462, 465, 470, 473, 479, 482, 488, 490, 494, 500,
	506, 509, 513, 519, 521, 528, 530, 534, 537, 540, 544, 546, 551, 556, 560,
	566, 568, 572, 574, 576, 580, 585, 588, 592, 594, 600, 603, 606, 608, 610,
	614, 620, 623, 628, 630, 632, 638, 644, 647, 653, 658, 660, 662, 665, 670,
	672, 676, 681, 683, 687, 689, 691, 693, 695, 697, 703, 706, 710, 715, 719,
	722, 726, 731, 735, 737, 741, 748, 750, 753, 757, 763, 766, 768, 775, 777,
	782, 785, 791, 795, 798, 800, 806, 811, 815, 818, 821, 824, 829, 832, 836,
	838, 842, 846, 850, 855, 860, 865, 870, 875, 878, 882, 886, 890, 895, 900,
	906, 910, 913, 916, 921, 925, 929, 932, 937, 940, 942, 944, 946, 952, 954,
	956, 958, 962, 966, 968, 971, 975, 979, 983, 987, 989, 994, 997, 1000,
	1003, 1008, 1014, 1017, 1024, 1028, 1032, 1034, 1036, 1040, 1044, 1048,
	1050, 1052, 1056, 1058, 1062, 1065, 1068, 1072, 1078, 1083, 1089, 1091,
	1094, 1097, 1101, 1104, 1106, 1110, 1115, 1117, 1119, 1121, 1126, 1129,
	1131, 1134, 1136, 1138, 1141, 1143, 1145, 1147, 1149, 1151, 1153, 1155,
	1157, 1159, 1161, 1164, 1166, 1168, 1170, 1172, 1174, 1176, 1178, 1180,
	1182, 1184, 1186, 1189, 1191, 1193, 1195, 1197, 1199, 1201, 1203, 1205,
	1208, 1210, 1212, 1214, 1217, 1219, 1221, 1223, 1225, 1227, 1229, 1231,
	1233, 1235, 1237, 1239, 1241, 1243, 1245, 1247, 1249, 1251, 1253, 1255,
	1257, 1259, 1261, 1263, 1265, 1268, 1270, 1272, 1274, 1276, 1278, 1280,
	1282, 1284, 1286, 1288, 1290, 1292, 1294, 1296, 1298, 1300, 1302, 1304,
	1306, 1308, 1310, 1312, 1314, 1316, 1318, 1320, 1322, 1324, 1326, 1328,
	1330, 1332, 1334, 1336, 1338, 1340, 1342, 1344, 1346, 1348, 1350, 1352}

func TestMarginWidth(t *testing.T) {

	ta := assert.New(t)

	cases := []struct {
		input int32
		want  uint32
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{3, 2},
		{4, 4},
		{15, 4},
		{16, 8},
		{255, 8},
		{256, 16},
		{65535, 16},
	}

	for i, c := range cases {
		got := marginWidth(c.input)
		ta.Equal(c.want, got,
			"%d-th: input: %#v; want: %#v; got: %#v",
			i+1, c.input, c.want, got)
	}

	ta.Panics(
		func() {
			marginWidth(int32(65536))
		})
}

func TestDense_New(t *testing.T) {
	ta := assert.New(t)

	cases := [][]int32{
		{},
		{0},
		{-1},
		{-1, -2},
		polyTestNums[:10],
		polyTestNums[:50],
		polyTestNums[:200],
		polyTestNums,
	}

	for _, nums := range cases {

		a := NewPolyArray(nums)
		for i, n := range nums {
			r := a.Get(int32(i))
			ta.Equal(n, r, "i=%d expect: %v; but: %v", i, n, r)
		}

		ta.Equal(len(nums), a.Len())

		// Stat() should work
		_ = a.Stat()
	}
}

func TestNewDense_eltWidthSmall(t *testing.T) {

	ta := assert.New(t)

	n := 500
	nums := make([]int32, n)
	for i := 0; i < n; i++ {
		nums[i] = int32(15 * i)
	}

	a := NewPolyArray(nums[:2])
	ta.Equal(uint32(0), a.Segments[0].Info[0])

	a = NewPolyArray(nums)
	ta.True(a.Stat()["bits/elt"] <= 2)

}

func TestNewDense_default(t *testing.T) {

	ta := assert.New(t)

	a := NewPolyArray(polyTestNums)
	ta.Equal(int32(len(polyTestNums)), a.N)

	fmt.Println(a.Stat())
	st := a.Stat()
	ta.Equal(int32(3), st["elt_width"])

}

func TestNewDense_big(t *testing.T) {

	ta := assert.New(t)

	n := int32(1024 * 1024)
	step := int32(64)
	ns := benchhelper.RandI32Slice(0, n, step)

	a := NewPolyArray(ns)

	for i, n := range ns {
		r := a.Get(int32(i))
		ta.Equal(n, r, "i=%d ", i)
	}
}

func TestNewDense_largenum(t *testing.T) {

	ta := assert.New(t)

	n := int32(1024 * 1024)
	step := int32(64)
	ns := benchhelper.RandI32Slice(1<<30, n, step)

	for i := 0; i < len(ns); i++ {
		if ns[i] < 0 {
			panic("<0")
		}
	}

	a := NewPolyArray(ns)

	for i, n := range ns {
		r := a.Get(int32(i))
		ta.Equal(n, r, "i=%d ", i)
	}
}

func TestDense_Get_panic(t *testing.T) {
	ta := assert.New(t)

	a := NewPolyArray(polyTestNums)
	ta.Panics(func() {
		a.Get(int32(len(polyTestNums)))
	})
	ta.Panics(func() {
		a.Get(int32(-1))
	})
}

func TestDense_Stat(t *testing.T) {
	ta := assert.New(t)

	a := NewPolyArray(polyTestNums)

	st := a.Stat()
	want := map[string]int32{
		"seg_cnt":   1,
		"elt_width": 3,
		"mem_elts":  224,
		"mem_total": st["mem_total"], // do not compare this
		"polys/seg": 3,
		"bits/elt":  7,
	}

	ta.Equal(want, st)
}

func TestDense_marshalUnmarshal(t *testing.T) {
	ta := assert.New(t)

	a := NewPolyArray(polyTestNums)

	bytes, err := proto.Marshal(a)
	ta.Nil(err, "want no error but: %+v", err)

	b := &PolyArray{}

	err = proto.Unmarshal(bytes, b)
	ta.Nil(err, "want no error but: %+v", err)

	for i, n := range polyTestNums {
		r := b.Get(int32(i))
		ta.Equal(n, r, "i=%d ", i)
	}
}

var Output int

func BenchmarkDense_Get(b *testing.B) {

	n := int32(1024 * 1024)
	mask := int(n - 1)
	step := int32(128)
	ns := benchhelper.RandI32Slice(0, n, step)

	s := int32(0)

	a := NewPolyArray(ns)
	fmt.Println(a.Stat())

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s += a.Get(int32(i & mask))
	}

	Output = int(s)
}

func BenchmarkNewDense(b *testing.B) {

	n := int32(1024 * 1024)
	step := int32(128)
	ns := benchhelper.RandI32Slice(0, n, step)

	s := int32(0)

	b.ResetTimer()
	var a *PolyArray
	for i := 0; i < b.N; i++ {
		a = NewPolyArray(ns)
		s += a.Get(int32(0))
	}

	fmt.Println(a.Stat())

	Output = int(s)
}
