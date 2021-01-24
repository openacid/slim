package kv3

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewShard(t *testing.T) {

	ta := require.New(t)

	cases := []struct {
		input     int32
		inputKeys []string
		want      *shard
	}{
		{0, []string{}, &shard{"", []string{}}},
		{0, []string{""}, &shard{"", []string{""}}},
		{0, []string{"", "a"}, &shard{"", []string{"", "a"}}},
		{1, []string{"a", "ac"}, &shard{"a", []string{"", "c"}}},
	}

	for i, c := range cases {
		got := newShard(c.input, c.inputKeys)
		ta.Equal(c.want, got, "%d-th: case: %+v", i+1, c)
	}
}

func TestShard(t *testing.T) {

	ta := require.New(t)

	keysAbe_Abo := []string{
		"Abe",
		"Abel",
		"Abelia",
		"Abelian",
		"Abelicea",
		"Abelite",
		"Abelmoschus",
		"Abelonian",
		"Abencerrages",
		"Aberdeen",
		"Aberdonian",
		"Aberia",
		"Abhorson",
		"Abie",
		"Abies",
		"Abietineae",
		"Abiezer",
		"Abigail",
		"Abipon",
		"Abitibi",
		"Abkhas",
		"Abkhasian",
		"Ablepharus",
		"Abnaki",
		"Abner",
		"Abo",
	}

	_ = keysAbe_Abo

	keysAaron_Abuta := []string{
		"Aaron",
		"Aaronic",
		"Aaronical",
		"Aaronite",
		"Aaronitic",
		"Aaru",
		"Ab",
		"Ababdeh",
		"Ababua",
		"Abadite",
		"Abama",
		"Abanic",
		"Abantes",
		"Abarambo",
		"Abaris",
		"Abasgi",
		"Abassin",
		"Abatua",
		"Abba",
		"Abbadide",
		"Abbasside",
		"Abbie",
		"Abby",
		"Abderian",
		"Abderite",
		"Abdiel",
		"Abdominales",
		"Abe",
		"Abel",
		"Abelia",
		"Abelian",
		"Abelicea",
		"Abelite",
		"Abelmoschus",
		"Abelonian",
		"Abencerrages",
		"Aberdeen",
		"Aberdonian",
		"Aberia",
		"Abhorson",
		"Abie",
		"Abies",
		"Abietineae",
		"Abiezer",
		"Abigail",
		"Abipon",
		"Abitibi",
		"Abkhas",
		"Abkhasian",
		"Ablepharus",
		"Abnaki",
		"Abner",
		"Abo",
		"Abobra",
		"Abongo",
		"Abraham",
		"Abrahamic",
		"Abrahamidae",
		"Abrahamite",
		"Abrahamitic",
		"Abram",
		"Abramis",
		"Abranchiata",
		"Abrocoma",
		"Abroma",
		"Abronia",
		"Abrus",
		"Absalom",
		"Absaroka",
		"Absi",
		"Absyrtus",
		"Abu",
		"Abundantia",
		"Abuta",
	}

	_ = keysAaron_Abuta

	cases := []struct {
		input []string
		thre  int32
		want  []*shard
	}{
		// {
		//     input: []string{"", "a"},
		//     thre:  1,
		// },
		// {
		//     input: []string{"", "a"},
		//     thre:  2,
		// },
		// {
		//     input: []string{"", "a"},
		//     thre:  3,
		// },
		{
			input: keysAbe_Abo,
			thre:  2,
		},
		// {
		//     input: keysAbe_Abo,
		//     thre:  4,
		// },
		// {
		//     input: keysAbe_Abo,
		//     thre:  5,
		// },
		// {
		//     input: keysAaron_Abuta,
		//     thre:  8,
		// },
		// {
		//     // input: testkeys.Load("200kweb2")[:500],
		//     input: testkeys.Load("200kweb2"),
		//     thre:  10,
		// },
	}

	for i, c := range cases {
		got := splitKeys3(c.input, c.thre)
		fmt.Println("got:")
		for _, g := range got {
			fmt.Println(g)
		}

		j := 0
		for gi, g := range got {
			if gi > 0 {
				ta.Greater(g.prefix, got[gi-1].prefix)
			}
			ta.LessOrEqual(len(g.suffixes), int(c.thre))
			for si, suf := range g.suffixes {
				ta.Equal(c.input[j], g.prefix+suf, "%d-th: shard-idx: %d, suffix-idx:%d, case: %+v", i+1, gi, si, c)
				j++
			}
		}

		ta.Equal(len(c.input), j)
		fmt.Println(len(c.input), "n-shard:", len(got), len(c.input)/len(got))
	}
}
