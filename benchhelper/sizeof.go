package benchhelper

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

// SizeOf returns the size in byte a value(not type) costs.
// TODO test
func SizeOf(data interface{}) int {
	return sizeof(reflect.ValueOf(data))
}

func SizeStat(v interface{}, depth int) string {
	lines := sizestat(reflect.ValueOf(v), depth)
	return strings.Join(lines, "\n")
}

var (
	mapsize       = int(unsafe.Sizeof(map[int]int{}))
	slicesize     = int(unsafe.Sizeof([]int8{}))
	stringsize    = int(unsafe.Sizeof(""))
	pointersize   = int(unsafe.Sizeof(&mapsize))
	interfacesize = int(unsafe.Sizeof(interface{}(nil)))
)

func sizestat(v reflect.Value, depth int) []string {
	header := fmt.Sprintf("%s: %d", v.Type(), sizeof(v))
	if depth == 0 {
		return []string{header}
	}

	depth -= 1

	lines := []string{header}

	switch v.Kind() {
	case reflect.Map:
		keys := v.MapKeys()
		for i := 0; i < len(keys); i++ {
			mapkey := keys[i]
			subs := sizestat(v.MapIndex(mapkey), depth)
			subs[0] = fmt.Sprintf("%s: ", mapkey) + subs[0]

			lines = append(lines, subs...)
		}
	case reflect.Slice, reflect.Array:
		for i, n := 0, v.Len(); i < n; i++ {
			subs := sizestat(v.Index(i), depth)
			subs[0] = fmt.Sprintf("%d: ", i) + subs[0]
			lines = append(lines, subs...)
		}

	case reflect.Ptr, reflect.Interface:
		p := (*[]byte)(unsafe.Pointer(v.Pointer()))
		if p != nil {
			lines = append(lines, sizestat(v.Elem(), depth)...)
		}
	case reflect.Struct:
		for i, n := 0, v.NumField(); i < n; i++ {
			subs := sizestat(v.Field(i), depth)
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
	sum := 0
	switch v.Kind() {
	case reflect.Map:
		keys := v.MapKeys()
		for i := 0; i < len(keys); i++ {
			mapkey := keys[i]
			s := sizeof(mapkey)
			if s < 0 {
				return -1
			}
			sum += s
			s = sizeof(v.MapIndex(mapkey))
			if s < 0 {
				return -1
			}
			sum += s
		}
	case reflect.Slice, reflect.Array:
		for i, n := 0, v.Len(); i < n; i++ {
			s := sizeof(v.Index(i))
			if s < 0 {
				return -1
			}
			sum += s
		}

	case reflect.String:
		for i, n := 0, v.Len(); i < n; i++ {
			s := sizeof(v.Index(i))
			if s < 0 {
				return -1
			}
			sum += s
		}

	case reflect.Ptr, reflect.Interface:
		p := (*[]byte)(unsafe.Pointer(v.Pointer()))
		if p == nil {
			sum = 0
		} else {
			sum = sizeof(v.Elem())
		}
	case reflect.Struct:
		for i, n := 0, v.NumField(); i < n; i++ {
			s := sizeof(v.Field(i))
			if s < 0 {
				return -1
			}
			sum += s
		}

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
		reflect.Int:
		sum = int(v.Type().Size())

	default:
		return -1
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
