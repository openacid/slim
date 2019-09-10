package benchhelper

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
)

var Fformat = struct {
	JPGHistogramTiny  string
	JPGHistogramSmall string
	JPGHistogramMid   string
}{
	JPGHistogramTiny: `
set terminal jpeg size 300,200;

# the bar width
set boxwidth 0.7 relative

set style fill solid border;

# show horizontal grid
set grid ytics;

set style line 101 lc rgb '#909090' lt 1 lw 1

# left and bottom border
set border 3 front ls 101

# plot title at the right top
set key font ",8"

# x-axis, the numbers
set tics font "Verdana,8"

# space between x-axis and x-tics label
set xtics offset 0,0.3,0

# cluster spacing is 1(bar-width)
set style histogram cluster gap 1
`,
	JPGHistogramSmall: `
set terminal jpeg size 400,300;
set boxwidth 0.8;
set style fill solid;
set grid ytics;

set style line 101 lc rgb '#909090' lt 1 lw 1
set border 3 front ls 101
`,
	JPGHistogramMid: `
set terminal jpeg size 600,400;
set boxwidth 0.7;
set style fill solid;
set grid ytics;

set style line 101 lc rgb '#909090' lt 1 lw 1
set border 3 front ls 101
`,
}

var LineStyles = struct {
	Colorful string
	Orange   string
	Yellow   string
	Green    string
	Cyan     string
	Blue     string
	Purple   string
}{
	Colorful: `
set style line 1 lc rgb '#4688F1' pt 1 ps 1 lt 1 lw 2;
set style line 2 lc rgb '#CA4E5D' pt 6 ps 1 lt 1 lw 2;
set style line 3 lc rgb '#79A2F1' pt 6 ps 1 lt 1 lw 2;
set style line 4 lc rgb '#8ED0F1' pt 6 ps 1 lt 1 lw 2;
set style line 5 lc rgb '#8AE7CC' pt 6 ps 1 lt 1 lw 2;

set style line 1 lc rgb '#4688F1' pt 1 ps 1 lt 1 lw 2;
set style line 2 lc rgb '#6CA8F3' pt 6 ps 1 lt 1 lw 2;
set style line 3 lc rgb '#79A2F1' pt 6 ps 1 lt 1 lw 2;
set style line 4 lc rgb '#8ED0F1' pt 6 ps 1 lt 1 lw 2;
set style line 5 lc rgb '#8AE7CC' pt 6 ps 1 lt 1 lw 2;
`,
	Orange: `
set style line 1  lc rgb '#edbe8a' pt 1 ps 1 lt 1 lw 2;
set style line 2  lc rgb '#e29543' pt 6 ps 1 lt 1 lw 2;
set style line 3  lc rgb '#da7409' pt 6 ps 1 lt 1 lw 2;
set style line 4  lc rgb '#c16400' pt 6 ps 1 lt 1 lw 2;
set style line 5  lc rgb '#ad5900' pt 6 ps 1 lt 1 lw 2;
`,
	Yellow: `
set style line 1  lc rgb '#e9d16c' pt 1 ps 1 lt 1 lw 2;
set style line 2  lc rgb '#e2c444' pt 6 ps 1 lt 1 lw 2;
set style line 3  lc rgb '#daaf08' pt 6 ps 1 lt 1 lw 2;
set style line 4  lc rgb '#cfb033' pt 6 ps 1 lt 1 lw 2;
set style line 5  lc rgb '#ad8a00' pt 6 ps 1 lt 1 lw 2;
`,
	Green: `
set style line 1 lc rgb '#a2e2b8' pt 1 ps 1 lt 1 lw 2;
set style line 2 lc rgb '#6ecd9b' pt 6 ps 1 lt 1 lw 2;
set style line 3 lc rgb '#5db191' pt 6 ps 1 lt 1 lw 2;
set style line 4 lc rgb '#519d7f' pt 6 ps 1 lt 1 lw 2;
set style line 5 lc rgb '#49856e' pt 6 ps 1 lt 1 lw 2;

`,
	Cyan: `
set style line 1 lc rgb '#adece1' pt 1 ps 1 lt 1 lw 2;
set style line 2 lc rgb '#5cd4d9' pt 6 ps 1 lt 1 lw 2;
set style line 3 lc rgb '#70bcca' pt 6 ps 1 lt 1 lw 2;
set style line 4 lc rgb '#4297a7' pt 6 ps 1 lt 1 lw 2;

`,
	Blue: `
set style line 1 lc rgb '#97c8d5' pt 1 ps 1 lt 1 lw 2;
set style line 2 lc rgb '#5ca6d9' pt 6 ps 1 lt 1 lw 2;
set style line 3 lc rgb '#4c80bc' pt 6 ps 1 lt 1 lw 2;
set style line 4 lc rgb '#4172a7' pt 6 ps 1 lt 1 lw 2;

`,
	Purple: `
set style line 1  lc rgb '#c4aecf' pt 1 ps 1 lt 1 lw 2;
set style line 2  lc rgb '#b674c0' pt 6 ps 1 lt 1 lw 2;
set style line 3  lc rgb '#a562a6' pt 6 ps 1 lt 1 lw 2;
set style line 4  lc rgb '#915593' pt 6 ps 1 lt 1 lw 2;
`,
}

var Plot = struct {
	Histogram string
}{
	Histogram: `
stats fn skip 1 nooutput
max_col = STATS_columns

plot for [col=2:max_col] fn \
using col:xticlabels(gprintf('10^{%T}',column(1)))           \
with histogram              \
linestyle col-1             \
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
		// fmt.Println(stdout.String())
		fmt.Println(stderr.String())
		panic(err)
	}

	err = ioutil.WriteFile(fn, stdout.Bytes(), 0777)
	if err != nil {
		panic(err)
	}
}
