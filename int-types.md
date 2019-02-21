# Choose appropriate integer types

## Signed integer

Signed int is used as arithmetic data.
Such as size, length, array index etc.

**Do not use unsigned int if possible**.
Minus operation with unsigned int overflows.
E.g.:

```go
var a uint32 = 3
var b uint32 = 5

fmt.Println(a - b)
// Output: 4294967294
```

-   `int` for in-memory size, length.

-   `int64` for large size, offset etc.

## Unsigned integer

Unsigned int is used as non-arithmetic data, such as bitmap, bit mask etc.

-   `uint64` for bitmap etc.
