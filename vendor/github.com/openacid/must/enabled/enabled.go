// Package enabled implements checking functions.
//
// Since 0.1.0
package enabled

import (
	"fmt"
	"strings"
	"time"

	"github.com/stretchr/testify/assert"
)

type beTyp struct{}

type tTyp struct {
	msg []string
}

func (t *tTyp) Errorf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	t.msg = append(t.msg, s)
}

func (t *tTyp) chk(rst bool) {
	if !rst {
		panic(strings.Join(t.msg, "\n"))
	}
}

var (
	// Be is the container of all checking APIs, such as "must.Be.Equal(a, b)".
	//
	// Since 0.1.0
	Be = &beTyp{}
)

func (be *beTyp) OK(f func()) {
	f()
}

// wrappers of testify/assert

func (be *beTyp) Condition(comp assert.Comparison, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.Condition(t, comp, msgAndArgs...))
}
func (be *beTyp) Contains(s, contains interface{}, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.Contains(t, s, contains, msgAndArgs...))
}
func (be *beTyp) DirExists(path string, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.DirExists(t, path, msgAndArgs...))
}
func (be *beTyp) ElementsMatch(listA, listB interface{}, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.ElementsMatch(t, listA, listB, msgAndArgs...))
}
func (be *beTyp) Empty(object interface{}, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.Empty(t, object, msgAndArgs...))
}
func (be *beTyp) Equal(expected, actual interface{}, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.Equal(t, expected, actual, msgAndArgs...))
}
func (be *beTyp) EqualError(theError error, errString string, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.EqualError(t, theError, errString, msgAndArgs...))
}
func (be *beTyp) EqualValues(expected, actual interface{}, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.EqualValues(t, expected, actual, msgAndArgs...))
}
func (be *beTyp) Error(err error, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.Error(t, err, msgAndArgs...))
}
func (be *beTyp) Exactly(expected, actual interface{}, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.Exactly(t, expected, actual, msgAndArgs...))
}
func (be *beTyp) Fail(failureMessage string, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.Fail(t, failureMessage, msgAndArgs...))
}
func (be *beTyp) FailNow(failureMessage string, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.FailNow(t, failureMessage, msgAndArgs...))
}
func (be *beTyp) False(value bool, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.False(t, value, msgAndArgs...))
}
func (be *beTyp) FileExists(path string, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.FileExists(t, path, msgAndArgs...))
}
func (be *beTyp) Implements(interfaceObject interface{}, object interface{}, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.Implements(t, interfaceObject, object, msgAndArgs...))
}
func (be *beTyp) InDelta(expected, actual interface{}, delta float64, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.InDelta(t, expected, actual, delta, msgAndArgs...))
}
func (be *beTyp) InDeltaMapValues(expected, actual interface{}, delta float64, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.InDeltaMapValues(t, expected, actual, delta, msgAndArgs...))
}
func (be *beTyp) InDeltaSlice(expected, actual interface{}, delta float64, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.InDeltaSlice(t, expected, actual, delta, msgAndArgs...))
}
func (be *beTyp) InEpsilon(expected, actual interface{}, epsilon float64, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.InEpsilon(t, expected, actual, epsilon, msgAndArgs...))
}
func (be *beTyp) InEpsilonSlice(expected, actual interface{}, epsilon float64, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.InEpsilonSlice(t, expected, actual, epsilon, msgAndArgs...))
}
func (be *beTyp) IsType(expectedType interface{}, object interface{}, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.IsType(t, expectedType, object, msgAndArgs...))
}
func (be *beTyp) JSONEq(expected string, actual string, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.JSONEq(t, expected, actual, msgAndArgs...))
}
func (be *beTyp) Len(object interface{}, length int, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.Len(t, object, length, msgAndArgs...))
}
func (be *beTyp) Nil(object interface{}, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.Nil(t, object, msgAndArgs...))
}
func (be *beTyp) NoError(err error, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.NoError(t, err, msgAndArgs...))
}
func (be *beTyp) NotContains(s, contains interface{}, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.NotContains(t, s, contains, msgAndArgs...))
}
func (be *beTyp) NotEmpty(object interface{}, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.NotEmpty(t, object, msgAndArgs...))
}
func (be *beTyp) NotEqual(expected, actual interface{}, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.NotEqual(t, expected, actual, msgAndArgs...))
}
func (be *beTyp) NotNil(object interface{}, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.NotNil(t, object, msgAndArgs...))
}
func (be *beTyp) NotPanics(f assert.PanicTestFunc, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.NotPanics(t, f, msgAndArgs...))
}
func (be *beTyp) NotRegexp(rx interface{}, str interface{}, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.NotRegexp(t, rx, str, msgAndArgs...))
}
func (be *beTyp) NotSubset(list, subset interface{}, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.NotSubset(t, list, subset, msgAndArgs...))
}
func (be *beTyp) NotZero(i interface{}, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.NotZero(t, i, msgAndArgs...))
}
func (be *beTyp) Panics(f assert.PanicTestFunc, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.Panics(t, f, msgAndArgs...))
}
func (be *beTyp) PanicsWithValue(expected interface{}, f assert.PanicTestFunc, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.PanicsWithValue(t, expected, f, msgAndArgs...))
}
func (be *beTyp) Regexp(rx interface{}, str interface{}, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.Regexp(t, rx, str, msgAndArgs...))
}
func (be *beTyp) Subset(list, subset interface{}, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.Subset(t, list, subset, msgAndArgs...))
}
func (be *beTyp) True(value bool, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.True(t, value, msgAndArgs...))
}
func (be *beTyp) WithinDuration(expected, actual time.Time, delta time.Duration, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.WithinDuration(t, expected, actual, delta, msgAndArgs...))
}
func (be *beTyp) Zero(i interface{}, msgAndArgs ...interface{}) {
	t := &tTyp{}
	t.chk(assert.Zero(t, i, msgAndArgs...))
}
