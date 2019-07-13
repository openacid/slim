// Package must provides "design-by-contract" support.
//
// With "go build -tags debug" it enables bug checking statement.
// If a bug is found it panics.
//
// Without "-tags debug", by default "go build" ignores bug checking statement
// by emitting an "NOP" instruction.
//
// Example: a function that right-shift an integer expect the input "a" is 8*n:
//
//     import (
//       "math/bits"
//       "github.com/openacid/must"
//     )
//
//     func rshift(a, b int) int {
//
//          must.Be.OK(func() {
//              must.Be.NotZero(b)
//              must.Be.True(bits.TrailingZeros(uint(a)) > 2)
//          })
//
//          return a >> uint(b)
//     }
//     func main() {
//         fmt.Println(rshift(0xf, 1))
//     }
//
// With the above code:
//
// `go run` just silently ignores the expectation and print the result.
//
// `go run -tags debug` would panic because `a` does not satisfy the input
// expectation.
//
// Since 0.1.0
package must
