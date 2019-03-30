package main

import "github.com/openacid/slim/genhelper"

var implHead = `package marshal

import "encoding/binary"
`

var implTemplate = `
// {{.Name}} converts {{.ValType}} to slice of {{.ValLen}} bytes and back.
type {{.Name}} struct{}

// Marshal converts {{.ValType}} to slice of {{.ValLen}} bytes.
func (c {{.Name}}) Marshal(d interface{}) []byte {
	b := make([]byte, {{.ValLen}})
	v := {{.EncodeCast}}(d.({{.ValType}}))
	binary.LittleEndian.Put{{.Decoder}}(b, v)
	return b
}

// Unmarshal converts slice of {{.ValLen}} bytes to {{.ValType}}.
// It returns number bytes consumed and an {{.ValType}}.
func (c {{.Name}}) Unmarshal(b []byte) (int, interface{}) {

	size := int({{.ValLen}})
	s := b[:size]

	d := {{.ValType}}(binary.LittleEndian.{{.Decoder}}(s))
	return size, d
}

// GetSize returns the size in byte after marshaling v.
func (c {{.Name}}) GetSize(d interface{}) int {
	return {{.ValLen}}
}

// GetMarshaledSize returns {{.ValLen}}.
func (c {{.Name}}) GetMarshaledSize(b []byte) int {
	return {{.ValLen}}
}
`

var testHead = `package marshal_test

import (
	"testing"

	"github.com/openacid/slim/marshal"
)
`

var testTemplate = `
func Test{{.Name}}(t *testing.T) {

	v0 := [8]byte{}
	v1 := [8]byte{1}
	v1234 := [8]byte{0x34, 0x12}
	vneg := [8]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

	cases := []struct {
		input    {{.ValType}}
		want     string
		wantsize int
	}{
		{0, string(v0[:{{.ValLen}}]), {{.ValLen}}},
		{1, string(v1[:{{.ValLen}}]), {{.ValLen}}},
		{0x1234, string(v1234[:{{.ValLen}}]), {{.ValLen}}},
		{^{{.ValType}}(0), string(vneg[:{{.ValLen}}]), {{.ValLen}}},
	}

	m := marshal.{{.Name}}{}

	for i, c := range cases {
		rst := m.Marshal(c.input)
		if string(rst) != c.want {
			t.Fatalf("%d-th: input: %v; want: %v; actual: %v",
				i+1, c.input, []byte(c.want), rst)
		}

		n := m.GetSize(c.input)
		if c.wantsize != n {
			t.Fatalf("%d-th: input: %v; wantsize: %v; actual: %v",
				i+1, c.input, c.wantsize, n)
		}

		n = m.GetMarshaledSize(rst)
		if c.wantsize != n {
			t.Fatalf("%d-th: input: %v; wantsize: %v; actual: %v",
				i+1, c.input, c.wantsize, n)
		}

		n, u64 := m.Unmarshal(rst)
		if c.input != u64 {
			t.Fatalf("%d-th: unmarshal: input: %v; want: %v; actual: %v",
				i+1, c.input, c.input, u64)
		}
		if c.wantsize != n {
			t.Fatalf("%d-th: unmarshaled size: input: %v; want: %v; actual: %v",
				i+1, c.input, c.wantsize, n)
		}
	}
}
`

func main() {

	pref := "int"
	implfn := pref + ".go"
	testfn := pref + "_test.go"

	impls := []interface{}{
		genhelper.IntConfig{Name: "U16", ValType: "uint16", ValLen: 2, Decoder: "Uint16", EncodeCast: "uint16"},
		genhelper.IntConfig{Name: "U32", ValType: "uint32", ValLen: 4, Decoder: "Uint32", EncodeCast: "uint32"},
		genhelper.IntConfig{Name: "U64", ValType: "uint64", ValLen: 8, Decoder: "Uint64", EncodeCast: "uint64"},
		genhelper.IntConfig{Name: "I16", ValType: "int16", ValLen: 2, Decoder: "Uint16", EncodeCast: "uint16"},
		genhelper.IntConfig{Name: "I32", ValType: "int32", ValLen: 4, Decoder: "Uint32", EncodeCast: "uint32"},
		genhelper.IntConfig{Name: "I64", ValType: "int64", ValLen: 8, Decoder: "Uint64", EncodeCast: "uint64"},
	}

	genhelper.Render(implfn, implHead, implTemplate, impls, []string{"gofmt", "unconvert"})
	genhelper.Render(testfn, testHead, testTemplate, impls, []string{"gofmt", "unconvert"})
}
