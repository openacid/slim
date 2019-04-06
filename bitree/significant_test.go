package bitree_test

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/openacid/slim/bitree"
)

func TestCalc(t *testing.T) {
	fn := "words"
	bytes, err := ioutil.ReadFile(fn)
	if err != nil {
		panic(err)
	}
	words := strings.Split(string(bytes), "\x0a")

	// l2 := []string{}

	// for _, nbit := range []int{3, 4, 5, 6, 7, 8} {
	for _, nbit := range []int{4} {
		KeyPerBit := float64(0)
		splits := 0
		w := words[:]
		for {
			ss, is := bitree.FindSignificantBits(w, nbit)
			KeyPerBit += float64(len(ss)) / float64(len(is))
			splits++

			steps := bitree.MakeSteps(is)

			fmt.Println(len(ss), steps)
			w = w[len(ss):]
			if len(w) == 0 {
				break
			}
		}

		KeyPerBit /= float64(splits)
		nkeys := int(float64(nbit) * KeyPerBit)
		bmsize := 1 << uint(nbit)
		fmt.Println(nbit, ":", KeyPerBit,
			"nkeys:", nkeys,
			"bit/k:", float64(bmsize)/float64(nkeys))
	}
}

func TestFindSignificantBits(t *testing.T) {

	fn := "words"
	bytes, err := ioutil.ReadFile(fn)
	if err != nil {
		panic(err)
	}
	words := strings.Split(string(bytes), "\x0a")

	cases := []struct {
		input     int
		wantStrs  []string
		wantIbits []int
	}{
		{
			1,
			words[:2],
			[]int{5},
		},
		{
			2,
			words[:3],
			[]int{5, 16},
		},
		{
			3,
			words[:4],
			[]int{5, 16, 32},
		},
		{
			4,
			words[:5],
			[]int{5, 16, 32, 48},
		},
		{
			5,
			words[:6],
			[]int{5, 16, 32, 47, 48},
		},
		{
			6,
			words[:9],
			[]int{5, 16, 32, 47, 48, 79},
		},
	}

	for i, c := range cases {
		rstStrs, rstIbits := bitree.FindSignificantBits(words, c.input)
		fmt.Println("strs:", len(rstStrs))
		for _, s := range rstStrs {
			fmt.Println(s)
		}
		for _, s := range rstStrs {
			fmt.Println(bitree.ToBin(s))
		}
		fmt.Println("ibits:")
		fmt.Println(rstIbits)

		if !reflect.DeepEqual(c.wantStrs, rstStrs) {
			t.Fatalf("%d-th: input: %#v; want: %#v; actual: %#v",
				i+1, c.input, c.wantStrs, rstStrs)
		}
		if !reflect.DeepEqual(c.wantIbits, rstIbits) {
			t.Fatalf("%d-th: input: %#v; want: %#v; actual: %#v",
				i+1, c.input, c.wantIbits, rstIbits)
		}
	}
}

func TestFindSignificantBits22(t *testing.T) {

	fn := "words"
	bytes, err := ioutil.ReadFile(fn)
	if err != nil {
		panic(err)
	}
	words := strings.Split(string(bytes), "\x0a")

	cases := []struct {
		input     int
		wantIbits []int
		// wantTrie  []bitree.SigTrie
	}{
		// {
			// 1,
			// []int{},
		// },
		{
			2,
			[]int{16},
		},
		{
			3,
			[]int{16, 39},
		},
		{
			4,
			[]int{16, 39, 80},
		},
		{
			5,
			[]int{16, 39, 80, 112},
		},
		{
			6,
			[]int{16, 39, 80, 112, 103},
		},
		{
			7,
			[]int{16, 39, 80, 112, 103, 121},
		},
	}

	for i, c := range cases {
		w := words[:c.input]
		rstIbits := bitree.FindSignificantBits222(words[:c.input])

		fmt.Println("strings:")
		for _, s := range words[:c.input] {
			fmt.Println(s)
		}
		for _, s := range words[:c.input] {
			fmt.Println(bitree.ToBin(s))
		}

		fmt.Println("ibits:")
		fmt.Println(rstIbits)

		if !reflect.DeepEqual(c.wantIbits, rstIbits) {
			t.Fatalf("%d-th: input: %#v; want: %#v; actual: %#v",
				i+1, c.input, c.wantIbits, rstIbits)
		}

		// // build trie from sig bits
		// trie := bitree.NewSigTrie(rstIbits)
		// fmt.Println( "trie:" )
		// fmt.Println(trie)

		bt:= bitree.NewBitrie(w)
		fmt.Printf("%+v\n", bt)

		// Get

		for i, k := range w {
			n := bt.Get(k)
			if i != n {
				trie := bitree.NewSigTrie(rstIbits)
				fmt.Println( "trie:" )
				fmt.Println(trie)
				t.Fatalf("in %v search %v expect: %v; but: %v", w, k, i, n)
			}
		}
	
	}

	// rstIbits := bitree.FindSignificantBits222(words[:32])
	// trie := bitree.NewSigTrie(rstIbits)
	// fmt.Println( "trie:" )
	// fmt.Println(trie)

	// bt:= bitree.NewBitrie(words[:32])
	// fmt.Println( len(bt) )
	// fmt.Printf("%+v\n", bt)
}

func TestFirstDiff(t *testing.T) {

	cases := []struct {
		a    string
		b    string
		want int
	}{
		{"", "a", 0},
		{"a", "b", 13},
		{"a", "ab", 16},
		{"\x07\x07", "\x07\x01", 27},
	}

	for i, c := range cases {
		rst := bitree.FirstDiff(c.a, c.b)
		if rst != c.want {
			t.Fatalf("%d-th: input: %#v %#v; want: %#v; actual: %#v",
				i+1, c.a, c.b, c.want, rst)
		}
	}
}
