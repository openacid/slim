package array

import (
	"fmt"
	"math/bits"

	"github.com/openacid/slim/array/pb"
	"github.com/openacid/slim/benchhelper"
	"gonum.org/v1/gonum/mat"
)

const (
	polyDegree  = 3
	polyCoefCnt = polyDegree + 1
)

var (
	// capByWidth maps EltWidth to log(word capacity, 2).
	capByWidth = map[uint32]uint32{
		1:  6,
		2:  5,
		4:  4,
		8:  3,
		16: 2,
	}
)

// Dense is an array that uses one or more polynomial to compress and store a
// series of int64.
//
// Since 0.5.2
type Dense struct {
	pb.PolyArray
}

// NewDense creates a "Dense" array from a slice of int64.
// A "Dense" array uses several polynomial curves to compress data.
// "segmentSize" specifies into how many segments to split the whole data set.
// "eltWidth" specifies the size in bit for every element and must be one of [1, 2, 4, 8, 16].
// Set "segmentSize" or "eltWidth" to 0 to use default the value 1024 and 8.
//
// It is very efficient to store a serias integers with a overall trend, such as
// a sorted array.
//
// For large set of data, the total memory spent is about eltWidth * len(nums).
//
// Smaller eltWidth requires more Polynomials thus results in more memory
// overhead.
// Larger eltWidth reduces Polynomials count but it requires more space for
// every element.
// Normally 4 or 8 bits for eltWidth is an appropriate choice.
//
// Since 0.5.2
func NewDense(nums []int64, segmentSize int, eltWidth uint) *Dense {

	if segmentSize <= 0 {
		segmentSize = 1024
	}

	if eltWidth == 0 {
		eltWidth = 8
	}

	width := uint32(eltWidth)

	segWidth := uint32(63 - bits.LeadingZeros64(uint64(segmentSize)))
	segmentSize = 1 << segWidth

	pa := &Dense{
		pb.PolyArray{
			EltWidth:     width,
			EltMask:      mask(width),
			WordCapWidth: capByWidth[width],
			WordCapMask:  mask(capByWidth[width]),
			SegWidth:     segWidth,
			SegMask:      mask(segWidth),
			N:            int32(len(nums)),
		},
	}

	for {
		if len(nums) > segmentSize {
			pa.newSeg(nums[:segmentSize])
			nums = nums[segmentSize:]
		} else {
			pa.newSeg(nums)
			break
		}
	}

	segs := make([]*pb.Segment, len(pa.Segments))
	copy(segs, pa.Segments)
	pa.Segments = segs

	return pa
}

// Get returns the uncompressed int64 value.
//
// A Get() costs about 11 ns
//
// Since 0.5.2
func (m *Dense) Get(i int32) int64 {
	iSeg := i >> m.SegWidth
	i = i & m.SegMask

	seg := m.Segments[iSeg]

	var j, l int = 0, len(seg.Starts)
	for ; j < l && i >= seg.Starts[j]; j++ {
	}

	// If i >= l, here is a panic
	j = (j - 1) * polyCoefCnt
	r := seg.Polynomials[j : j+polyCoefCnt]
	v := int64(eval3(r, float64(i)))

	d := seg.Words[i>>m.WordCapWidth]

	inWordIndex := i & m.WordCapMask
	d = d >> (uint(inWordIndex) * uint(m.EltWidth))

	return v + d&int64(m.EltMask)
}

// Len returns number of elements.
//
// Since 0.5.2
func (m *Dense) Len() int {
	return int(m.N)
}

// Stat returns a map describing memory usage.
//
//    elt_width :8
//    seg_size  :1024
//    seg_cnt   :512
//    polys/seg :7
//    mem_elts  :1048576
//    mem_total :1195245
//    bits/elt  :9
//
// Since 0.5.2
func (m *Dense) Stat() map[string]int64 {
	nseg := len(m.Segments)
	totalmem := benchhelper.SizeOf(m)

	total := 0
	memWords := 0
	for _, seg := range m.Segments {
		total += len(seg.Polynomials)
		memWords += benchhelper.SizeOf(seg.Words)
	}

	n := m.Len()
	if n == 0 {
		n = 1
	}

	st := map[string]int64{
		"seg_size":  int64(1 << m.SegWidth),
		"seg_cnt":   int64(nseg),
		"elt_width": int64(m.EltWidth),
		"mem_total": int64(totalmem),
		"mem_elts":  int64(memWords),
		"bits/elt":  int64(totalmem * 8 / n),
		"polys/seg": int64(total / nseg / polyCoefCnt),
	}

	return st
}

// newSeg creates a new segment with nums
//
// Since 0.5.2
func (m *Dense) newSeg(nums []int64) {

	n := int32(len(nums))
	xs := make([]float64, n)
	ys := make([]float64, n)

	for i, v := range nums {
		xs[i] = float64(i)
		ys[i] = float64(v)
	}

	margin := (1 << m.EltWidth) - 1

	seg := &pb.Segment{
		Words: m.makeWords(int(n)),
	}

	start := int32(0)
	Starts := make([]int32, 0)
	polys := make([]float64, 0)

	for start < n {
		poly, nn := polyFitByMargin(xs[start:], ys[start:], polyDegree, float64(margin))

		Starts = append(Starts, start)
		polys = append(polys, poly...)

		for i := int32(0); i < nn; i++ {
			j := start + i

			v := eval3(poly, xs[j])

			d := nums[j] - int64(v)
			if d > int64(margin) || d < 0 {
				panic(fmt.Sprintf("d=%d must smaller than %d and > 0", d, margin))
			}

			inWordIndex := j & m.WordCapMask
			seg.Words[j>>m.WordCapWidth] |= d << uint(int32(m.EltWidth)*inWordIndex)
		}

		start += nn
	}

	Starts = append(Starts, start)
	seg.Starts = make([]int32, len(Starts))
	copy(seg.Starts, Starts)

	seg.Polynomials = make([]float64, len(polys))
	copy(seg.Polynomials, polys)

	m.Segments = append(m.Segments, seg)

}

// makeWords creates a just long enough int64 slice for all elements.
//
// Since 0.5.2
func (m *Dense) makeWords(n int) []int64 {
	eltPerWord := int(64 / m.EltWidth)
	nWords := (n + eltPerWord - 1) / eltPerWord

	return make([]int64, nWords)
}

// polyFitByMargin finds a polynomial curve that covers as many points as possible so that
// their distant to the curve smaller than margin.
//
// It returns the coeffecients of the curve and how many points is covered.
//
// Since 0.5.2
func polyFitByMargin(xs, ys []float64, degree int, margin float64) ([]float64, int32) {

	l, r := int32(0), int32(len(xs)+1)

	for {
		for l < r-1 {
			mid := (l + r) / 2
			xx, yy := xs[:mid], ys[:mid]

			poly := polyFit(xx, yy, degree)
			max, min := maxminResiduals(poly, xx, yy)
			if max-min <= margin {
				l = mid
			} else {
				r = mid
			}
		}

		xs, ys = xs[:l], ys[:l]
		poly := polyFit(xs, ys, degree)
		max, min := maxminResiduals(poly, xs, ys)

		// max-min are not guaranteed to be incremental.
		// Thus if max-min exceed margin, reset r to l and re-run binary search.
		if max-min > margin {
			l, r = 0, l
			continue
		} else {
			// Makes every point be above the curve
			poly[0] += min
			return poly, l
		}
	}
}

// polyFit models a polynomial y from sample points xs and ys, to minimizes the squared residuals.
// It returns coefficients of the polynomial y:
//
//    y = β₁ + β₂x + β₃x² + ...
//
// It use linear regression, which assumes y is in form of:
//        m
//    y = ∑ βⱼ Φⱼ(x)
//        j=1
//
// In our case:
//    Φⱼ(x) = x^(j-1)
//
// Then
//    (Xᵀ X) βⱼ = Xᵀ Y
//    Xᵢⱼ = [ Φⱼ(xᵢ) ]
//
// See https://en.wikipedia.org/wiki/Least_squares#Linear_least_squares
//
// Since 0.5.2
func polyFit(xs, ys []float64, degree int) []float64 {

	// Number of sample points
	n := len(xs)
	deg := degree

	// We do not need degree-4 curve for 5 point or less
	if deg > n-1 {
		deg = n - 1
	}

	// Number of βⱼ is degree+1
	m := deg + 1

	// build matrix: Xᵢⱼ = [ Φⱼ(xᵢ) ]
	d := make([]float64, n*m)
	for i := 0; i < n; i++ {
		x := xs[i]
		v := float64(1)
		for j := 0; j < m; j++ {
			d[i*m+j] = v
			v *= x
		}
	}

	mtr := mat.NewDense(n, m, d)

	var right mat.Dense
	var coef mat.Dense
	var beta mat.Dense

	// coef * beta = right
	coef.Mul(mtr.T(), mtr)
	right.Mul(mtr.T(), mat.NewDense(n, 1, ys))
	beta.Solve(&coef, &right)

	rst := make([]float64, degree+1)
	for i := 0; i < m; i++ {
		rst[i] = beta.At(i, 0)
	}

	for i := m; i < degree+1; i++ {
		rst[i] = 0
	}

	return rst
}

// maxminResiduals finds max and min offset along a curve.
//
// Since 0.5.2
func maxminResiduals(poly, xs, ys []float64) (float64, float64) {

	max, min := float64(0), float64(0)

	for i, x := range xs {
		v := eval3(poly, x)
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

// // eval1 is a quick path of eval.
// //
// // Since 0.5.2
// func eval1(poly []float64, x float64) float64 {
//     return poly[0] + poly[1]*x
// }

// // eval2 is a quick path of eval.
// //
// // Since 0.5.2
// func eval2(poly []float64, x float64) float64 {
//     return poly[0] + poly[1]*x + poly[2]*x*x
// }

// eval3 is a quick path of eval.
//
// Since 0.5.2
func eval3(poly []float64, x float64) float64 {
	return poly[0] + poly[1]*x + poly[2]*x*x + poly[3]*x*x*x
}

// // eval4 is a quick path of eval.
// //
// // Since 0.5.2
// func eval4(poly []float64, x float64) float64 {
//     return poly[0] + poly[1]*x + poly[2]*x*x + poly[3]*x*x*x + poly[4]*x*x*x*x
// }

// eval evaluates polynomial at x
//
// Since 0.5.2
func eval(poly []float64, x float64) float64 {
	rst := float64(0)
	pow := float64(1)
	for _, coef := range poly {
		rst += coef * pow
		pow *= x
	}

	return rst
}

func mask(width uint32) int32 {
	return (1 << width) - 1
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
