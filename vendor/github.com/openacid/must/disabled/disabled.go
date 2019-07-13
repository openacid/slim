// Package disabled implements disabled checking functions.
//
// Since 0.1.0
package disabled

import (
	"time"

	"github.com/stretchr/testify/assert"
)

type foo struct{}

var (
	// Be is the container of all checking APIs, such as "must.Be.Equal(a, b)".
	//
	// Since 0.1.0
	Be = &foo{}
)

func (fl *foo) OK(f func()) {}

func (fl *foo) Condition(comp assert.Comparison, msgAndArgs ...interface{})            {}
func (fl *foo) Contains(s, contains interface{}, msgAndArgs ...interface{})            {}
func (fl *foo) DirExists(path string, msgAndArgs ...interface{})                       {}
func (fl *foo) ElementsMatch(listA, listB interface{}, msgAndArgs ...interface{})      {}
func (fl *foo) Empty(object interface{}, msgAndArgs ...interface{})                    {}
func (fl *foo) Equal(expected, actual interface{}, msgAndArgs ...interface{})          {}
func (fl *foo) EqualError(theError error, errString string, msgAndArgs ...interface{}) {}
func (fl *foo) EqualValues(expected, actual interface{}, msgAndArgs ...interface{})    {}
func (fl *foo) Error(err error, msgAndArgs ...interface{})                             {}
func (fl *foo) Exactly(expected, actual interface{}, msgAndArgs ...interface{})        {}
func (fl *foo) Fail(failureMessage string, msgAndArgs ...interface{})                  {}
func (fl *foo) FailNow(failureMessage string, msgAndArgs ...interface{})               {}
func (fl *foo) False(value bool, msgAndArgs ...interface{})                            {}
func (fl *foo) FileExists(path string, msgAndArgs ...interface{})                      {}
func (fl *foo) Implements(interfaceObject interface{}, object interface{}, msgAndArgs ...interface{}) {
}
func (fl *foo) InDelta(expected, actual interface{}, delta float64, msgAndArgs ...interface{}) {}
func (fl *foo) InDeltaMapValues(expected, actual interface{}, delta float64, msgAndArgs ...interface{}) {
}
func (fl *foo) InDeltaSlice(expected, actual interface{}, delta float64, msgAndArgs ...interface{}) {
}
func (fl *foo) InEpsilon(expected, actual interface{}, epsilon float64, msgAndArgs ...interface{}) {
}
func (fl *foo) InEpsilonSlice(expected, actual interface{}, epsilon float64, msgAndArgs ...interface{}) {
}
func (fl *foo) IsType(expectedType interface{}, object interface{}, msgAndArgs ...interface{}) {}
func (fl *foo) JSONEq(expected string, actual string, msgAndArgs ...interface{})               {}
func (fl *foo) Len(object interface{}, length int, msgAndArgs ...interface{})                  {}
func (fl *foo) Nil(object interface{}, msgAndArgs ...interface{})                              {}
func (fl *foo) NoError(err error, msgAndArgs ...interface{})                                   {}
func (fl *foo) NotContains(s, contains interface{}, msgAndArgs ...interface{})                 {}
func (fl *foo) NotEmpty(object interface{}, msgAndArgs ...interface{})                         {}
func (fl *foo) NotEqual(expected, actual interface{}, msgAndArgs ...interface{})               {}
func (fl *foo) NotNil(object interface{}, msgAndArgs ...interface{})                           {}
func (fl *foo) NotPanics(f assert.PanicTestFunc, msgAndArgs ...interface{})                    {}
func (fl *foo) NotRegexp(rx interface{}, str interface{}, msgAndArgs ...interface{})           {}
func (fl *foo) NotSubset(list, subset interface{}, msgAndArgs ...interface{})                  {}
func (fl *foo) NotZero(i interface{}, msgAndArgs ...interface{})                               {}
func (fl *foo) Panics(f assert.PanicTestFunc, msgAndArgs ...interface{})                       {}
func (fl *foo) PanicsWithValue(expected interface{}, f assert.PanicTestFunc, msgAndArgs ...interface{}) {
}
func (fl *foo) Regexp(rx interface{}, str interface{}, msgAndArgs ...interface{}) {}
func (fl *foo) Subset(list, subset interface{}, msgAndArgs ...interface{})        {}
func (fl *foo) True(value bool, msgAndArgs ...interface{})                        {}
func (fl *foo) WithinDuration(expected, actual time.Time, delta time.Duration, msgAndArgs ...interface{}) {
}
func (fl *foo) Zero(i interface{}, msgAndArgs ...interface{}) {}
