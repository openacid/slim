// Package size provides value size operations.
package size

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

// const uintSize = 32 << (^uint(0) >> 32 & 1) // 32 or 64

// Of returns the size in byte a value(not type) costs.
//
// Since 0.1.0
func Of(data interface{}) int {
	if data == nil {
		return 0
	}
	return sizeof(reflect.ValueOf(data))
}

// Stat returns a multi line string describe size of every component of a
// vavlue, such as slice element, struct member:
//
//	   size_test.my: 48
//		   a: []int32: 36
//			   0: int32: 4
//			   1: int32: 4
//			   2: int32: 4
//		   b: [3]int32: 12
//			   0: int32: 4
//			   1: int32: 4
//			   2: int32: 4
//
// Since 0.1.0
func Stat(v interface{}, depth int, maxItem int) string {
	lines := stat(reflect.ValueOf(v), depth, maxItem)
	return strings.Join(lines, "\n")
}

var (
	mapsize       = int(unsafe.Sizeof(map[int]int{}))    // 8
	slicesize     = int(unsafe.Sizeof([]int8{}))         // 24
	stringsize    = int(unsafe.Sizeof(""))               // 16
	pointersize   = int(unsafe.Sizeof(&mapsize))         // 8
	interfacesize = int(unsafe.Sizeof(interface{}(nil))) // 16
)

func stat(v reflect.Value, depth int, maxItem int) []string {
	if !v.IsValid() {
		return []string{"<nil>"}
	}
	header := fmt.Sprintf("%s: %d", v.Type(), sizeof(v))
	if depth == 0 {
		return []string{header}
	}

	depth -= 1

	lines := []string{header}

	switch v.Kind() {
	case reflect.Map:
		keys := v.MapKeys()
		for i := 0; i < len(keys) && i < maxItem; i++ {
			mapkey := keys[i]
			subs := stat(v.MapIndex(mapkey), depth, maxItem)
			subs[0] = fmt.Sprintf("%s: ", mapkey) + subs[0]

			lines = append(lines, subs...)
		}
	case reflect.Slice, reflect.Array:
		for i, n := 0, v.Len(); i < n && i < maxItem; i++ {
			subs := stat(v.Index(i), depth, maxItem)
			subs[0] = fmt.Sprintf("%d: ", i) + subs[0]
			lines = append(lines, subs...)
		}

	case reflect.Ptr:
		p := (*[]byte)(unsafe.Pointer(v.Pointer()))
		if p != nil {
			lines = append(lines, stat(v.Elem(), depth, maxItem)...)
		}
	case reflect.Interface:
		lines = append(lines, stat(v.Elem(), depth, maxItem)...)
	case reflect.Struct:
		for i, n := 0, v.NumField(); i < n; i++ {
			subs := stat(v.Field(i), depth, maxItem)
			subs[0] = fmt.Sprintf("%s: ", v.Type().Field(i).Name) + subs[0]
			lines = append(lines, subs...)
		}
	}
	for i, s := range lines {
		if i > 0 {
			lines[i] = "    " + s
		}
	}

	return lines
}

func sizeof(v reflect.Value) int {
	if !v.IsValid() {
		return 0
	}

	sum := 0
	switch v.Kind() {
	case reflect.Map:
		keys := v.MapKeys()
		for i := 0; i < len(keys); i++ {
			mapkey := keys[i]
			s := sizeof(mapkey)
			sum += s
			s = sizeof(v.MapIndex(mapkey))
			sum += s
		}
	case reflect.Slice, reflect.Array:
		for i, n := 0, v.Len(); i < n; i++ {
			s := sizeof(v.Index(i))
			sum += s
		}

	case reflect.String:
		for i, n := 0, v.Len(); i < n; i++ {
			s := sizeof(v.Index(i))
			sum += s
		}

	case reflect.Ptr:
		p := (*[]byte)(unsafe.Pointer(v.Pointer()))
		if p == nil {
			sum = 0
		} else {
			sum = sizeof(v.Elem())
		}
	case reflect.Interface:
		sum = sizeof(v.Elem())
	case reflect.Struct:
		for i, n := 0, v.NumField(); i < n; i++ {
			s := sizeof(v.Field(i))
			sum += s
		}

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
		reflect.Int:
		sum = int(v.Type().Size())
	default:
		panic(fmt.Sprintf("unknown kind: %s", v.Kind()))
	}

	switch v.Kind() {
	case reflect.Map:
		sum += mapsize
	case reflect.Slice:
		sum += slicesize
	case reflect.String:
		sum += stringsize
	case reflect.Ptr:
		sum += pointersize
	case reflect.Interface:
		sum += interfacesize
	}
	return sum
}
