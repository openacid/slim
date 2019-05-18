package benchhelper

import (
	"reflect"
	"unsafe"
)

// SizeOf returns the size in byte a value(not type) costs.
// TODO test
func SizeOf(data interface{}) int {
	return sizeof(reflect.ValueOf(data))
}

var (
	mapsize       = int(unsafe.Sizeof(map[int]int{}))
	slicesize     = int(unsafe.Sizeof([]int8{}))
	stringsize    = int(unsafe.Sizeof(""))
	pointersize   = int(unsafe.Sizeof(&mapsize))
	interfacesize = int(unsafe.Sizeof(interface{}(nil)))
)

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
