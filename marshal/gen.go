package marshal

//go:generate go run gen/impl_gen.go
//go:generate gofmt -s -w int.go int_test.go
//go:generate unconvert -v -apply github.com/openacid/slim/marshal
