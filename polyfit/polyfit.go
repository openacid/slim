// Package polyfit models a polynomial y from sample points xs and ys, to minimizes the squared residuals.
//
// See https://en.wikipedia.org/wiki/Least_squares#Linear_least_squares
//
// Since 0.5.4
package polyfit

import (
	"fmt"
	"strings"

	"gonum.org/v1/gonum/mat"
)

// Fitting models a polynomial y from sample points xs and ys, to minimizes the squared residuals.
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
//    (Xᵀ × X) βⱼ = Xᵀ × Y
//    Xᵢⱼ = [ Φⱼ(xᵢ) ]
//
// See https://en.wikipedia.org/wiki/Least_squares#Linear_least_squares
//
// Since 0.5.4
type Fitting struct {
	N      int
	Degree int

	// cache Xᵀ X
	xtx []float64
	// cache Xᵀ Y
	xty []float64
}

// NewFitting creates a new polynomial fitting context.
//
// Since 0.5.4
func NewFitting(xs, ys []float64, degree int) *Fitting {

	n := len(xs)

	m := degree + 1

	f := &Fitting{
		N:      0,
		Degree: degree,

		xtx: make([]float64, m*m),
		xty: make([]float64, m),
	}

	for i := 0; i < m*m; i++ {
		f.xtx[i] = 0
	}

	for i := 0; i < m; i++ {
		f.xty[i] = 0
	}

	for i := 0; i < n; i++ {
		f.Add(xs[i], ys[i])
	}

	return f
}

// Add a point(x, y) into this fitting.
//
// Since 0.5.4
func (f *Fitting) Add(x, y float64) {

	m := f.Degree + 1

	xpows := make([]float64, m)
	v := float64(1)
	for i := 0; i < m; i++ {
		xpows[i] = v
		v *= x
	}

	for i := 0; i < m; i++ {
		for j := 0; j < m; j++ {
			f.xtx[i*m+j] += xpows[i] * xpows[j]
		}
	}

	for i := 0; i < m; i++ {
		f.xty[i] += xpows[i] * y
	}

	f.N++
}

// Merge Combines two sets of sample data.
//
// This can be done because:
//    |X₁|ᵀ × |X₁| = X₁ᵀ × X₁ + X₂ᵀ × X₂
//    |X₂|    |X₂|
//
// Since 0.5.4
func (f *Fitting) Merge(b *Fitting) {

	if f.Degree != b.Degree {
		panic(fmt.Sprintf("different degree: %d %d", f.Degree, b.Degree))
	}

	f.N += b.N

	m := f.Degree + 1

	for i := 0; i < m; i++ {
		f.xty[i] += b.xty[i]
		for j := 0; j < m; j++ {
			f.xtx[i*m+j] += b.xtx[i*m+j]
		}
	}
}

// Solve the equation and returns coefficients of result polynomial.
// The number of coefficients is f.Degree + 1.
//
// Since 0.5.4
func (f *Fitting) Solve(minimizeDegree bool) []float64 {

	m := f.Degree + 1

	coef := mat.NewDense(m, m, f.xtx)
	right := mat.NewDense(m, 1, f.xty)

	if minimizeDegree && f.Degree+1 > f.N {

		m = f.N

		coef = coef.Slice(0, m, 0, m).(*mat.Dense)
		right = right.Slice(0, m, 0, 1).(*mat.Dense)
	}

	var beta mat.Dense
	beta.Solve(coef, right)

	rst := make([]float64, f.Degree+1)
	for i := 0; i < m; i++ {
		rst[i] = beta.At(i, 0)
	}

	for i := m; i < f.Degree+1; i++ {
		rst[i] = 0
	}

	return rst
}

// String prints human readable info of a fitting.
// It includes:
// n: the number of points.
// degree: expected degree of polynomial.
// and two matrix.
//
// Since 0.5.4
func (f *Fitting) String() string {

	m := f.Degree + 1
	ss := []string{}

	xtx := f.matrixStrings(f.xtx)

	ss = append(ss, fmt.Sprintf("n=%d degree=%d", f.N, f.Degree))
	ss = append(ss, xtx...)
	ss = append(ss, "")
	for i := 0; i < m; i++ {
		s := fmt.Sprintf("%3.3f", f.xty[i])
		ss = append(ss, s)
	}
	return strings.Join(ss, "\n")
}

func (f *Fitting) matrixStrings(mat []float64) []string {

	m := f.Degree + 1

	ss := []string{}

	for i := 0; i < m; i++ {
		line := []string{}
		for j := 0; j < m; j++ {
			s := fmt.Sprintf("%3.3f", mat[i*m+j])
			line = append(line, s)
		}

		linestr := strings.Join(line, " ")
		ss = append(ss, linestr)
	}

	return ss
}
