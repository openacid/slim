package bitree

import (
	"fmt"
	"sort"
	"strings"
)

func ToScaled(a byte) string {
	rst := []string{}
	for i := 7; i >= 0; i-- {
		if a&(1<<uint(i)) != 0 {
			rst = append(rst, ".1")
		} else {
			rst = append(rst, ".0")
		}
	}

	return strings.Join(rst, "")

}

func ToBin(s string) string {
	rst := []string{}
	for i := 0; i < len(s); i++ {
		rst = append(rst, ToScaled(s[i]))
	}
	return strings.Join(rst, " ")
}

func MakeSteps(is []int) []int {
	rst := make([]int, len(is))
	p0 := 0
	for i, p := range is {
		rst[i] = p - p0
		p0 = p
	}

	return rst
}

type SigBit struct {
	Row int
	Pos int
}

func FindSignificantBits(s []string, nbit int) ([]string, []int) {
	prev := s[0]
	strs := []string{s[0]}
	ibits := map[int]bool{}
	for _, str := range s[1:] {
		i := FirstDiff(prev, str)
		// fmt.Printf("found diff bit of %s %s at: %d\n", ToBin(prev), ToBin(str), i)
		ibits[i] = true
		strs = append(strs, str)
		if len(ibits) == nbit {
			break
		}
		prev = str
	}

	is := make([]int, 0, len(ibits))
	for k := range ibits {
		is = append(is, k)
	}

	sort.Ints(is)

	return strs, is
}

func MinSigBit(sb []SigBit) int {

	min := 0
	for i := 1; i < len(sb); i++ {
		if sb[i].Pos < sb[min].Pos {
			min = i
		}
	}
	return min
}

type TNode struct {
	SigBits   []SigBit
	Leaf      bool
	LeafValue int
	L, R      int
	Level     int
}

type SigTrie []TNode

func (t SigTrie) Strings(i int) []string {

	node := t[i]
	fmt.Printf("tostring: %d, %#v\n", i, node)
	nodestr := strings.Repeat(" ", node.Level*4)+
	fmt.Sprintf("%02d)", i)
	if node.Leaf {
		return []string{nodestr + fmt.Sprintf(":%d", node.LeafValue)}
	}

	rst := t.Strings(node.L)
	rst = append(rst, nodestr)
	rst = append(rst, t.Strings(node.R)...)

	return rst
}

func (t SigTrie) String() string {
	return strings.Join(t.Strings(0), "\n")
}

func NewSigTrie(sb []SigBit) SigTrie {

	if len(sb) == 0 {
		return nil
	}

	nodes := SigBitsBulidTrie(sb)
	return SigTrie(nodes)
}

// SigBitsBulidTrie form a trie and return the root.
func SigBitsBulidTrie(sigbits []SigBit) []TNode {

	children := make([]TNode, 0, len(sigbits)*2+1)
	children = append(children, TNode{
		SigBits: sigbits,
		Level:   0,
	})

	for i := 0; i < len(children); i++ {
		n := &children[i]

		if n.Leaf {
			continue
		}

		min := MinSigBit(n.SigBits)

		fmt.Println("sigbits:", n.SigBits)
		fmt.Println("min:", min)

		minRow := n.SigBits[min].Row

		var t TNode
		if min == 0 {
			t = TNode{
				Leaf:      true,
				LeafValue: minRow,
				Level:     n.Level + 1,
			}
		} else {
			t = TNode{
				Leaf:    false,
				SigBits: n.SigBits[0:min],
				Level:   n.Level + 1,
			}
		}
		n.L = len(children)
		children = append(children, t)
		fmt.Printf("append left child: %#v\n", t)
		// fmt.Printf("children: %#v\n", children)

		if min == len(n.SigBits)-1 {
			t = TNode{
				Leaf:      true,
				LeafValue: minRow + 1,
				Level:     n.Level + 1,
			}
		} else {
			t = TNode{
				Leaf:    false,
				SigBits: n.SigBits[min+1:],
				Level:   n.Level + 1,
			}
		}
		n.R = len(children)
		children = append(children, t)
		fmt.Printf("append right child: %#v\n", t)

		fmt.Printf("processed: %#v\n", n)
		// fmt.Printf("children: %#v\n", children)
		// fmt.Println("current trie:", children)
	}

	return children
}

func FindSignificantBits222(s []string) []SigBit {
	prev := s[0]
	sigbits := make([]SigBit, 0)
	for row, str := range s[1:] {
		i := FirstDiff(prev, str)
		sigbits = append(sigbits, SigBit{Row: row, Pos: i})
		// fmt.Printf("found diff bit of %s %s at: %d\n", ToBin(prev), ToBin(str), i)
		prev = str
	}

	return sigbits
}

// string is converted to infinit bits:
// every char is converted to binary form,
// then add a "1" before every bit:
//
//    "a" -> 0x61 -> 01100001 -> 1011111010101011
//
// Then append infinite number of 0 after end of string.
func FirstDiff(a, b string) int {

	la := len(a)
	lb := len(b)

	l := la
	if lb < la {
		l = lb
	}

	for i := 0; i < l; i++ {
		if a[i] == b[i] {
			continue
		}

		v := a[i] ^ b[i]

		// find the most significant 1
		j := 0
		for j = 0; j < 8; j++ {
			if v >= 128 {
				break
			}
			v <<= 1
		}

		return (i*8+j)*2 + 1
	}

	// one is shorter and is a prefix of the other.
	return l * 8 * 2

}
