package benchhelper

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
)

var Fformat = struct {
	JPGHistogramSmall string
	JPGHistogramMid   string
}{
	JPGHistogramSmall: `
set terminal jpeg size 400,300;
set boxwidth 0.8;
set style fill solid;
set grid ytics;
`,
	JPGHistogramMid: `
set terminal jpeg size 800,600;
set boxwidth 0.8;
set style fill solid;
set grid ytics;
`,
}

var LineStyles = struct {
	Green  string
	Yellow string
}{
	Green: `
set style line 1 lc rgb '#aaffa5' pt 1 ps 1 lt 1 lw 2;
set style line 2 lc rgb '#99f094' pt 6 ps 1 lt 1 lw 2;
set style line 3 lc rgb '#87e082' pt 6 ps 1 lt 1 lw 2;
set style line 4 lc rgb '#76d171' pt 6 ps 1 lt 1 lw 2;
set style line 5 lc rgb '#65c260' pt 6 ps 1 lt 1 lw 2;
set style line 6 lc rgb '#54b34f' pt 6 ps 1 lt 1 lw 2;
set style line 7 lc rgb '#42a33d' pt 6 ps 1 lt 1 lw 2;
set style line 8 lc rgb '#31942c' pt 6 ps 1 lt 1 lw 2;
`,
	Yellow: `
set style line 1  lc rgb '#fffc85' pt 1 ps 1 lt 1 lw 2;
set style line 2  lc rgb '#fdf476' pt 6 ps 1 lt 1 lw 2;
set style line 3  lc rgb '#fbed67' pt 6 ps 1 lt 1 lw 2;
set style line 4  lc rgb '#f8e559' pt 6 ps 1 lt 1 lw 2;
set style line 5  lc rgb '#f6dd4a' pt 6 ps 1 lt 1 lw 2;
set style line 6  lc rgb '#f4d63b' pt 6 ps 1 lt 1 lw 2;
set style line 7  lc rgb '#f2ce2c' pt 6 ps 1 lt 1 lw 2;
set style line 8  lc rgb '#efc61e' pt 6 ps 1 lt 1 lw 2;
set style line 9  lc rgb '#edbf0f' pt 6 ps 1 lt 1 lw 2;
set style line 10 lc rgb '#ebb700' pt 6 ps 1 lt 1 lw 2;
`,
}

var Plot = struct {
	Histogram string
}{
	Histogram: `
stats fn skip 1 nooutput
max_col = STATS_columns

plot for [col=2:max_col] fn \
using col:xtic(1) \
with histogram \
linestyle col-1 \
title columnheader
`,
}

// Plot a image by gnuplot script "script" and output it to "fn".
func Fplot(fn, script string) {
	gp := exec.Command("gnuplot")
	stdin, err := gp.StdinPipe()
	if err != nil {
		panic(err)
	}

	var stdout bytes.Buffer
	gp.Stdout = &stdout

	var stderr bytes.Buffer
	gp.Stderr = &stderr

	err = gp.Start()
	if err != nil {
		panic(err)
	}

	_, err = io.WriteString(stdin, script)
	if err != nil {
		panic(err)
	}
	err = stdin.Close()
	if err != nil {
		panic(err)
	}

	err = gp.Wait()
	if err != nil {
		fmt.Println(stdout.String())
		fmt.Println(stderr.String())
		panic(err)
	}

	err = ioutil.WriteFile(fn, stdout.Bytes(), 0777)
	if err != nil {
		panic(err)
	}
}
