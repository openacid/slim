package array

import (
	"fmt"
	"math"
	"math/bits"

	"github.com/openacid/slim/benchhelper"
	"github.com/openacid/slim/polyfit"
)

const (
	maxEltWidth    = 16
	maxBitIndex    = 1 << 16
	maxSegmentSize = maxBitIndex / maxEltWidth

	// allows max eltWidth=64 thus "start" is under 2^16
	segmentSize = 1024
	segWidth    = uint(10)
	segMask     = int32(1024 - 1)
)

const polyDegree = 2
const polyCoefCnt = polyDegree + 1

// evalpoly is a quick path of eval.
//
// Since 0.5.2
func evalpoly(poly []float64, x float64) float64 {
	// return poly[0] + poly[1]*x + poly[2]*x*x + poly[3]*x*x*x
	return poly[0] + poly[1]*x + poly[2]*x*x
}

// NewPolyArray creates a "PolyArray" array from a slice of int32.
// A "PolyArray" array uses several polynomial curves to compress data.
//
// It is very efficient to store a serias integers with a overall trend, such as
// a sorted array.
//
// Since 0.5.2
func NewPolyArray(nums []int32) *PolyArray {

	pa := &PolyArray{
		N: int32(len(nums)),
	}

	for {
		if len(nums) > segmentSize {
			pa.addSeg(nums[:segmentSize])
			nums = nums[segmentSize:]
		} else {
			pa.addSeg(nums)
			break
		}
	}

	segs := make([]*Segment, len(pa.Segments))
	copy(segs, pa.Segments)
	pa.Segments = segs

	return pa
}

// Get returns the uncompressed int32 value.
//
// A Get() costs about 15 ns
//
// Since 0.5.2
func (m *PolyArray) Get(i int32) int32 {

	if i >= m.N {
		panic(fmt.Sprintf("i=%d out of boundary N=%d", i, m.N))
	}

	iSeg := i >> segWidth
	i = i & segMask
	x := float64(i)

	seg := m.Segments[iSeg]

	iPoly := i >> seg.PolySpanWidth
	i = i & (seg.PolySpan - 1)

	j := iPoly * polyCoefCnt
	p := seg.Polynomials

	// evalpoly(r, x)
	v := int32(p[j] + p[j+1]*x + p[j+2]*x*x)

	// start, eltWidth := unpackInfo(seg.Info[iPoly])
	info := seg.Info[iPoly]
	start, eltWidth := uint16(info>>8), uint8(info)

	if eltWidth == 0 {
		return v
	}

	eltMask := (1 << eltWidth) - 1

	ibit := int32(start) + i*int32(eltWidth)

	d := seg.Words[ibit>>6]

	d = d >> (uint(ibit & 63))

	return v + int32(d&int64(eltMask))
}

// Len returns number of elements.
//
// Since 0.5.2
func (m *PolyArray) Len() int {
	return int(m.N)
}

// Stat returns a map describing memory usage.
//
//    elt_width :8
//    seg_cnt   :512
//    polys/seg :7
//    mem_elts  :1048576
//    mem_total :1195245
//    bits/elt  :9
//
// Since 0.5.2
func (m *PolyArray) Stat() map[string]int32 {
	nseg := len(m.Segments)
	totalmem := benchhelper.SizeOf(m)

	polyCnt := 0
	memWords := 0
	widthAvg := 0
	for _, seg := range m.Segments {
		polyCnt += len(seg.Polynomials)
		memWords += benchhelper.SizeOf(seg.Words)

		width := 0
		for _, inf := range seg.Info {
			_, w := unpackInfo(inf)
			width += int(w)
		}
		widthAvg += width / len(seg.Info)
	}

	n := m.Len()
	if n == 0 {
		n = 1
	}

	st := map[string]int32{
		"seg_cnt":   int32(nseg),
		"elt_width": int32(widthAvg / nseg),
		"mem_total": int32(totalmem),
		"mem_elts":  int32(memWords),
		"bits/elt":  int32(totalmem * 8 / n),
		"polys/seg": int32(polyCnt / nseg / polyCoefCnt),
	}

	return st
}

func packInfo(start uint16, width uint8) uint32 {
	return (uint32(start) << 8) + (uint32(width))
}

func unpackInfo(info uint32) (uint16, uint8) {
	return uint16(info >> 8), uint8(info)
}

func (m *PolyArray) addSeg(nums []int32) {

	n := int32(len(nums))
	xs := make([]float64, n)
	ys := make([]float64, n)

	for i, v := range nums {
		xs[i] = float64(i)
		ys[i] = float64(v)
	}

	// min polyspan
	polyspan := int32(16)
	fts := initFittings(xs, ys, polyspan)

	polyspan, polys, widths, fts := findMinFittings(xs, ys, fts, polyspan)

	infos := make([]uint32, 0)
	words := make([]int64, n) // max size

	// where the first elt of a polynomial in words
	start := int32(0)

	var s, e int32
	for i := int32(0); i < int32(len(fts)); i++ {
		s = e
		e += int32(fts[i].N)

		poly := polys[i*polyCoefCnt : i*polyCoefCnt+polyCoefCnt]
		eltWidth := widths[i]
		margin := int32((1 << eltWidth) - 1)
		if eltWidth > 0 {
			start = (start + int32(eltWidth) - 1)
			start -= start % int32(eltWidth)
		}

		if start >= 65536 {
			panic(fmt.Sprintf("wordStart is too large:%d", start))
		}

		infos = append(infos, packInfo(uint16(start), uint8(eltWidth)))

		for j := s; j < e; j++ {

			v := evalpoly(poly, xs[j])

			d := int64(nums[j]) - int64(v)
			if d > int64(margin) || d < 0 {
				panic(fmt.Sprintf("d=%d must smaller than %d and > 0", d, margin))
			}
			iWord := start >> 6
			words[iWord] |= d << uint(start&63)
			start += int32(eltWidth)
		}
	}

	// last start is for len(nums)
	infos = append(infos, packInfo(uint16(start), uint8(0)))

	nWords := (start + 63) >> 6

	seg := &Segment{}
	seg.PolySpan = polyspan
	seg.PolySpanWidth = log2u64(uint64(polyspan))

	seg.Polynomials = append(polys[:0:0], polys...)
	seg.Words = append(words[:0:0], words[:nWords]...)
	seg.Info = append(infos[:0:0], infos...)

	m.Segments = append(m.Segments, seg)
}

func log2u64(i uint64) uint32 {

	if i == 0 {
		return 0
	}

	return uint32(63 - bits.LeadingZeros64(i))
}

func initFittings(xs, ys []float64, polysize int32) []*polyfit.Fitting {

	fts := make([]*polyfit.Fitting, 0)
	n := int32(len(xs))

	for i := int32(0); i < n; i += polysize {
		s := i
		e := s + polysize
		if e > n {
			e = n
		}

		xx := xs[s:e]
		yy := ys[s:e]
		ft := polyfit.NewFitting(xx, yy, polyDegree)
		fts = append(fts, ft)
	}
	return fts
}

func findMinFittings(xs, ys []float64, fts []*polyfit.Fitting, polysize int32) (int32, []float64, []uint32, []*polyfit.Fitting) {
	minMem := 1 << 30

	minPolys := make([]float64, 0)
	minWidths := make([]uint32, 0)
	minFts := []*polyfit.Fitting(nil)
	minPolySize := polysize
	for {
		polys := make([]float64, 0)
		widths := make([]uint32, 0)

		mem := 0

		var s, e int32
		for _, ft := range fts {
			s = e
			e += int32(ft.N)

			poly := ft.Solve(true)
			max, min := maxminResiduals(poly, xs[s:e], ys[s:e])
			margin := int32(math.Ceil(max - min))
			poly[0] += min

			eltWidth := marginWidth(margin)
			mem += memCost(poly, eltWidth, int32(ft.N))

			polys = append(polys, poly...)
			widths = append(widths, eltWidth)
		}

		if minMem > mem {
			minMem = mem
			minPolys = append(polys[:0:0], polys...)
			minWidths = append(widths[:0:0], widths...)
			minFts = fts
			minPolySize = polysize

			fts = mergeFittings(fts)
			polysize *= 2

		} else {
			return minPolySize, minPolys, minWidths, minFts
		}
	}

}

func mergeFittings(fts []*polyfit.Fitting) []*polyfit.Fitting {

	newFts := make([]*polyfit.Fitting, 0)
	for i := 0; i < len(fts)/2; i++ {
		f := polyfit.NewFitting(nil, nil, fts[i*2].Degree)
		f.Merge(fts[i*2])
		f.Merge(fts[i*2+1])
		newFts = append(newFts, f)
	}
	if len(fts)%2 == 1 {
		i := len(fts) - 1
		f := polyfit.NewFitting(nil, nil, fts[i].Degree)
		f.Merge(fts[i])
		newFts = append(newFts, f)
	}

	return newFts
}

func marginWidth(margin int32) uint32 {
	for _, width := range []uint32{0, 1, 2, 4, 8, 16} {
		if int32(1)<<width > margin {
			return width
		}
	}

	panic(fmt.Sprintf("margin is too large: %d >= 2^16", margin))
}

func memCost(poly []float64, eltWidth uint32, n int32) int {
	mm := 0
	mm += 64                             // PolySpan and PolySpanWidth
	mm += 64 * len(poly)                 // Polynomials
	mm += 32 * (len(poly) / polyCoefCnt) // Info
	mm += int(eltWidth) * int(n)         // Words
	return mm
}

// maxminResiduals finds max and min residuals along a curve.
//
// Since 0.5.2
func maxminResiduals(poly, xs, ys []float64) (float64, float64) {

	max, min := float64(0), float64(0)

	for i, x := range xs {
		v := evalpoly(poly, x)
		diff := ys[i] - v
		if diff > max {
			max = diff
		}
		if diff < min {
			min = diff
		}
	}

	return max, min
}

// // polyStr convert a polynomial to string for human
// //
// // Since 0.5.2
// func polyStr(poly []float64) string {

//     elts := []string{}

//     for i, coef := range poly {
//         if coef == 0 {
//             continue
//         }

//         cc := fmt.Sprintf("%.4f", coef)

//         if i > 0 {
//             cc += "x"
//         }
//         if i > 1 {
//             cc += fmt.Sprintf("^%d", i)
//         }

//         elts = append(elts, cc)
//     }

//     return strings.Replace(strings.Join(elts, " + "), "+ -", "- ", -1)
// }
