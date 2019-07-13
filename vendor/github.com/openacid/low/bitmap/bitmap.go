// Package bitmap provides basic bitmap operations.
// A bitmap uses []uint64 as storage.
//
// Since 0.1.8
package bitmap

func init() {
	initMasks()
	initSelectLookup()
}
