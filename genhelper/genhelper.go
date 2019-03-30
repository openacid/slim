// Package genhelper provides with utilities to generate codes
package genhelper

import (
	"fmt"
	"os"
	"os/exec"
	"text/template"
)

// IntConfig defines a integer type based template redner config.
type IntConfig struct {
	// Name of the data type.
	Name string
	// ValType is the actual underlying int type, such as int32.
	ValType string
	// ValLen specifies the length of ValType
	ValLen int
	// Decoder defines the name of function to decode raw bytes into ValType.
	Decoder string
	// EncodeCast defines a cast type/function to convert values before encode.
	// Because sometimes encoder does not provides a exact type.
	EncodeCast string
}

// Render generate a file "fn".
// File content is defined by a "header", a repeated body template
// "tmpl" and slice of "data" to render the body template.
// Additianlly some linters can be specified to run after generating.
// Supported linters are "gofmt" and "unconvert".
func Render(fn string, header string, tmpl string, datas []interface{}, linters []string) {

	f, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	err = f.Truncate(0)
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(f, "// Code generated 'by go generate ./...'; DO NOT EDIT.")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, header)

	t, err := template.New("foo").Parse(tmpl)
	if err != nil {
		panic(err)
	}

	for _, d := range datas {
		err = t.Execute(f, d)
		if err != nil {
			panic(err)
		}
	}
	err = f.Sync()
	if err != nil {
		panic(err)
	}

	err = f.Close()
	if err != nil {
		panic(err)
	}

	for _, linter := range linters {
		var cmds []string
		switch linter {
		case "gofmt":
			cmds = []string{"gofmt", "-s", "-w", fn}
		case "unconvert":
			cmds = []string{"unconvert", "-v", "-apply", "./"}
		default:
			panic("unknown linter:" + linter)
		}

		out, err := exec.Command(cmds[0], cmds[1:]...).CombinedOutput()
		if err != nil {
			fmt.Println(cmds)
			fmt.Println(string(out))
			fmt.Println(err)
			panic(err)
		}
	}
}
