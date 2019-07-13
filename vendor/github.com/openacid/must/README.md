# must

`must` is a "design by contract" implementation in golang,
for addressing silent bug that is caused by unexpected input etc.

"Design by contract" requires that some conditions **MUST** be satisfied for the input of a
function, thus the function does not need to do a lot check on arguments.

It is the responsibility of `must` to enable these checking in a test env and to
disable them in production env(for performance concern).

[![Travis-CI](https://api.travis-ci.org/openacid/must.svg?branch=master)](https://travis-ci.org/openacid/must)
[![GoDoc](https://godoc.org/github.com/openacid/must?status.svg)](http://godoc.org/github.com/openacid/must)
[![Report card](https://goreportcard.com/badge/github.com/openacid/must)](https://goreportcard.com/report/github.com/openacid/must)
[![GolangCI](https://golangci.com/badges/github.com/openacid/must.svg)](https://golangci.com/r/github.com/openacid/must)
[![Sourcegraph](https://sourcegraph.com/github.com/openacid/must/-/badge.svg)](https://sourcegraph.com/github.com/openacid/must?badge)
![stability-stable](https://img.shields.io/badge/stability-stable-green.svg)

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
 

- [Usage](#usage)
  - [Efficiency](#efficiency)
- [API](#api)
- [Examples](#examples)
- [Install](#install)
- [Customize tags](#customize-tags)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Usage

To enable expectation check in test environment, use `go build|test -tags debug`.
To disable it for a release, just `go build`.

```go
package main

import (
	"fmt"
	"math/bits"

	"github.com/openacid/must"
)

func rshift(a, b int) int {

	// "go build" emits a single No-op instruction.
	// "go build -tags debug" will call the function and to the checking.
	must.Be.OK(func() {
		must.Be.NotZero(b)
		must.Be.True(bits.TrailingZeros(uint(a)) > 2,
			"a must be multiple of 8")
	})

	return a >> uint(b)
}

func main() {
	// panic at line 19 with "go run -tags debug"
	fmt.Println(rshift(0xf, 1))
}
```

With the above code:

**Enable check** with `go run -tags debug .`.
It would panic because `a` does not satisfy the input
expectation:

```
panic:
        ...
        Error:          Should be true
        Messages:       a must be multiple of 8
```

**Disable check** with `go run .`
It just silently ignores the expectation and print the result:

```
7
```

## Efficiency

**With debug**, there are checking statement instructions generated:

```
> go build -tags debug -o bin-debug .
> go tool objdump -S bin-debug

TEXT main.rshift(SB) github.com/openacid/must/examples/rshift/composite/main.go
func rshift(a, b int) int {
  ...
        mustbe.OK(func() {
  0x12481fd             0f57c0                  XORPS X0, X0
  0x1248200             0f110424                MOVUPS X0, 0(SP)
  ...
        f()
  0x124822d             488b1c24                MOVQ 0(SP), BX
  ...
```

**Without debug**, there is only a `NOPL` instruction:

```
> go build -o bin-release .
> go tool objdump -S bin-release

TEXT main.rshift(SB) github.com/openacid/must/examples/rshift/composite/main.go
        mustbe.OK(func() {
  0x1246030             90                      NOPL
  0x1246031             488b4c2410              MOVQ 0x10(SP), CX
        return a >> uint(b)
  0x1246036             4883f940                CMPQ $0x40, CX
  ...
  0x1246050             c3                      RET
```


# API

`must` uses the popular [testify](https://github.com/stretchr/testify) as underlying
asserting engine, thus there is a corresponding function defined for every
`testify` assertion function.

And `must.Be.OK(f func())` should be the entry of a set of checks:

```
must.Be.OK(f func())

must.Be.Condition(comp assert.Comparison, msgAndArgs ...interface{})
must.Be.Contains(s, contains interface{}, msgAndArgs ...interface{})
must.Be.DirExists(path string, msgAndArgs ...interface{})
must.Be.ElementsMatch(listA, listB interface{}, msgAndArgs ...interface{})
must.Be.Empty(object interface{}, msgAndArgs ...interface{})
must.Be.Equal(expected, actual interface{}, msgAndArgs ...interface{})
must.Be.EqualError(theError error, errString string, msgAndArgs ...interface{})
must.Be.EqualValues(expected, actual interface{}, msgAndArgs ...interface{})
must.Be.Error(err error, msgAndArgs ...interface{})
must.Be.Exactly(expected, actual interface{}, msgAndArgs ...interface{})
must.Be.Fail(failureMessage string, msgAndArgs ...interface{})
must.Be.FailNow(failureMessage string, msgAndArgs ...interface{})
must.Be.False(value bool, msgAndArgs ...interface{})
must.Be.FileExists(path string, msgAndArgs ...interface{})
must.Be.Implements(interfaceObject interface{}, object interface{}, msgAndArgs ...interface{})
must.Be.InDelta(expected, actual interface{}, delta float64, msgAndArgs ...interface{})
must.Be.InDeltaMapValues(expected, actual interface{}, delta float64, msgAndArgs ...interface{})
must.Be.InDeltaSlice(expected, actual interface{}, delta float64, msgAndArgs ...interface{})
must.Be.InEpsilon(expected, actual interface{}, epsilon float64, msgAndArgs ...interface{})
must.Be.InEpsilonSlice(expected, actual interface{}, epsilon float64, msgAndArgs ...interface{})
must.Be.IsType(expectedType interface{}, object interface{}, msgAndArgs ...interface{})
must.Be.JSONEq(expected string, actual string, msgAndArgs ...interface{})
must.Be.Len(object interface{}, length int, msgAndArgs ...interface{})
must.Be.Nil(object interface{}, msgAndArgs ...interface{})
must.Be.NoError(err error, msgAndArgs ...interface{})
must.Be.NotContains(s, contains interface{}, msgAndArgs ...interface{})
must.Be.NotEmpty(object interface{}, msgAndArgs ...interface{})
must.Be.NotEqual(expected, actual interface{}, msgAndArgs ...interface{})
must.Be.NotNil(object interface{}, msgAndArgs ...interface{})
must.Be.NotPanics(f assert.PanicTestFunc, msgAndArgs ...interface{})
must.Be.NotRegexp(rx interface{}, str interface{}, msgAndArgs ...interface{})
must.Be.NotSubset(list, subset interface{}, msgAndArgs ...interface{})
must.Be.NotZero(i interface{}, msgAndArgs ...interface{})
must.Be.Panics(f assert.PanicTestFunc, msgAndArgs ...interface{})
must.Be.PanicsWithValue(expected interface{}, f assert.PanicTestFunc, msgAndArgs ...interface{})
must.Be.Regexp(rx interface{}, str interface{}, msgAndArgs ...interface{})
must.Be.Subset(list, subset interface{}, msgAndArgs ...interface{})
must.Be.True(value bool, msgAndArgs ...interface{})
must.Be.WithinDuration(expected, actual time.Time, delta time.Duration, msgAndArgs ...interface{})
must.Be.Zero(i interface{}, msgAndArgs ...interface{})
```

See: [assert-functions](https://godoc.org/github.com/stretchr/testify/assert)


# Examples

[rshift-simple-check](examples/rshift/simple)

[rshift-composite-check](examples/rshift/composite)


# Install

```
go get -u github.com/openacid/must
```


# Customize tags

`must` provides with a default tag "debug" to enable expectation check.
It is also very easy to define your own tags.
To do this, create two files in one of your package, such as `mymust_debug.go` and `mymust_release.go`, like following:

`mymust_debug.go`:

```go
// +build mydebug

package your_package
import "github.com/openacid/must/enabled"
var mymustBe = enabled.Be
```

`mymust_release.go`:

```go
// +build !mydebug

package your_package
import "github.com/openacid/must/disabled"
var mymustBe = disabled.Be
```

And replace `must.Be` with `mymustBe` in your source codes.
Then your could enable it just by:

```
go build -tags mydebug
```

See more: [go-build-constraints](https://golang.org/pkg/go/build/#hdr-Build_Constraints)