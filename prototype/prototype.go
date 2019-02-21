// Package prototype provides generated Protocol Buffers for cache serialization.
//
// Usage:
//		go generate ./...
//
// It just read the comment line starts with "//go:generate" and run.
package prototype

//go:generate protoc --proto_path=./proto --go_out=. proto/compacted_array.proto
