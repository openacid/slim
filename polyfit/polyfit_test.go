package polyfit_test

import (
	"testing"

	. "github.com/openacid/slim/polyfit"
	"github.com/stretchr/testify/assert"
)

func TestNewFitting(t *testing.T) {

	ta := assert.New(t)

	var f *Fitting
	var want string

	f = NewFitting([]float64{1}, []float64{1}, 3)
	want = `
n=1 degree=3
1.000 1.000 1.000 1.000
1.000 1.000 1.000 1.000
1.000 1.000 1.000 1.000
1.000 1.000 1.000 1.000

1.000
1.000
1.000
1.000`[1:]
	ta.Equal(want, f.String())

	f = NewFitting([]float64{1, 1}, []float64{1, 1}, 3)
	want = `
n=2 degree=3
2.000 2.000 2.000 2.000
2.000 2.000 2.000 2.000
2.000 2.000 2.000 2.000
2.000 2.000 2.000 2.000

2.000
2.000
2.000
2.000`[1:]
	ta.Equal(want, f.String())

	f = NewFitting([]float64{1, 2}, []float64{1, 2}, 3)
	want = `
n=2 degree=3
2.000 3.000 5.000 9.000
3.000 5.000 9.000 17.000
5.000 9.000 17.000 33.000
9.000 17.000 33.000 65.000

3.000
5.000
9.000
17.000`[1:]
	ta.Equal(want, f.String())

}

func TestFitting_Add(t *testing.T) {

	ta := assert.New(t)

	xs := []float64{1, 2, 3, 4}
	ys := []float64{1, 2, 3, 4}

	f := NewFitting([]float64{}, []float64{}, 3)
	ta.Equal(0, f.N)

	for i, x := range xs {
		f.Add(x, ys[i])
		ta.Equal(i+1, f.N)
	}
}

func TestFitting_Merge(t *testing.T) {

	ta := assert.New(t)

	xs := []float64{1, 2, 3, 4}
	ys := []float64{1, 2, 3, 4}

	f := NewFitting(xs, ys, 3)

	fa := NewFitting(xs[:2], ys[:2], 3)
	fb := NewFitting(xs[2:], ys[2:], 3)

	fa.Merge(fb)

	ta.Equal(f.String(), fa.String())
}

func TestFitting_Solve(t *testing.T) {

	ta := assert.New(t)

	xs := []float64{1, 2, 3, 4}
	ys := []float64{6, 5, 7, 10}

	cases := []struct {
		degree int
		want   []float64
	}{
		{1, []float64{3.5, 1.4}},
		{2, []float64{8.5, -3.6, 1}},
		{3, []float64{12, -9.1666666, 3.5, -0.33333}},
		{4, []float64{12, -9.1666666, 3.5, -0.33333, 0}},
		{5, []float64{12, -9.1666666, 3.5, -0.33333, 0, 0}},
		{6, []float64{12, -9.1666666, 3.5, -0.33333, 0, 0, 0}},
	}

	for i, c := range cases {
		f := NewFitting(xs, ys, c.degree)
		for _, minimize := range []bool{false, true} {

			poly := f.Solve(minimize)
			ta.Equal(c.degree+1, len(poly))

			if minimize {
				ta.InDeltaSlice(c.want, poly, 0.0001,
					"%d-th: input: %#v; want: %#v; actual: %#v",
					i+1, c.degree, c.want, poly)
			}

			if c.degree >= len(xs)-1 {
				// curve pass every point

				for j, x := range xs {
					v := eval(poly, x)

					ta.InDelta(ys[j], v, 0.0001,
						"%d-th: input: %#v; want: %#v; actual: %#v",
						i+1, c.degree, ys[j], v)
				}
			}
		}

	}
}

func eval(poly []float64, x float64) float64 {
	rst := float64(0)
	pow := float64(1)
	for _, coef := range poly {
		rst += coef * pow
		pow *= x
	}

	return rst
}
