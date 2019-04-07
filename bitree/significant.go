package bitree

import (
	"encoding/binary"
	"fmt"
	gobits "math/bits"
	"sort"
	"strconv"
	"strings"

	"github.com/openacid/slim/bits"
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

func bitmap2Bin(bs []byte) string {
	rst := []string{}
	for _, b := range bs {
		s := strconv.FormatInt(int64(gobits.Reverse8(b)), 2)
		s = fmt.Sprintf("%08s", s)
		rst = append(rst, s)
	}
	return strings.Join(rst, " ")
}

func bitmap64ToBin(bm uint64) string {
	rst := []string{}
	for i := 0; i < 8; i++ {
		b := byte(bm >> uint(i*8))
		s := strconv.FormatInt(int64(gobits.Reverse8(b)), 2)
		s = fmt.Sprintf("%08s", s)
		rst = append(rst, s)
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

func MinSigBit(sb []int, frm, to int) int {

	min := frm
	for i := frm + 1; i < to; i++ {
		if sb[i] < sb[min] {
			min = i
		}
	}
	return min
}

type TNode struct {
	frm, to   int
	Leaf      bool
	LeafValue int
	L, R      int
	Pos       int
}

type SigTrie []TNode

func (t SigTrie) Strings(i int, lvl int) []string {

	node := t[i]
	nodestr := strings.Repeat(" ", lvl*4) +
		fmt.Sprintf("%02d)", i)
	if node.Leaf {
		return []string{nodestr + fmt.Sprintf(":%d", node.LeafValue)}
	}

	rst := t.Strings(node.L, lvl+1)
	rst = append(rst, nodestr+fmt.Sprintf("@%d", node.Pos))
	rst = append(rst, t.Strings(node.R, lvl+1)...)

	return rst
}

func (t SigTrie) String() string {
	return strings.Join(t.Strings(0, 0), "\n")
}

type Bitrie []byte

func (bt Bitrie) String() string {
	bs := []byte(bt)

	rst := []string{}
	for _, b := range bs[1:] {

		// Reverse to place lower bit on the left
		s := strconv.FormatInt(int64(gobits.Reverse8(b)), 2)
		s = fmt.Sprintf("%08s", s)
		rst = append(rst, s)
	}
	bm := strings.Join(rst, " ")

	return bm
}

func GetIthbit(s string, i int) int {
	// TODO test it
	ibyte := i / 16
	if ibyte >= len(s) {
		return 0
	}

	// get non-ending mark bit
	if i&1 == 0 {
		return 1
	}

	b := s[ibyte]
	ibit := (i % 16) / 2
	fmt.Println("get key bit from:", ToScaled(b), ibit)
	fmt.Println(ToScaled(b >> uint(7-ibit)))
	return int((b >> uint(7-ibit)) & 1)
}

func (bt Bitrie) Get(k string) int {

	bm := (uint64(bt[1]) |
		(uint64(bt[2]) << 8) |
		(uint64(bt[3]) << 16) |
		(uint64(bt[4]) << 24) |
		(uint64(bt[5]) << 32) |
		(uint64(bt[6]) << 40) |
		(uint64(bt[7]) << 48) |
		(uint64(bt[8]) << 54))

	leftLeavesMask := uint64(0)
	leftMost := 1

	// trie starts at the second bit in bitmap
	i := 1
	// skip flag byte and bitmap
	posOffset := 1 + 8
	ikey := 0

	// to read the first inner node pos
	posIndex := 0

	for {

		if bm&(1<<uint(i)) != 0 {
			// leaf
			break
		}

		// i is a inner node

		preceding1 := bits.OnesCount64Before(bm, uint(i))

		nthInner := i - preceding1 - 1

		for ; posIndex <= nthInner; posIndex++ {
			pos, n := binary.Varint(bt[posOffset:])
			posOffset += n
			ikey += int(pos)
		}

		bit := GetIthbit(k, ikey)

		i = (i-preceding1)*2 + bit
		leftMost = (leftMost - bits.OnesCount64Before(bm, uint(leftMost))) * 2
		levelMask := ((uint64(1) << uint(i)) - 1) ^ ((uint64(1) << uint(leftMost)) - 1)
		leftLeavesMask |= levelMask
	}

	for i != leftMost {
		// start to filter all left leaves
		i = (i - bits.OnesCount64Before(bm, uint(i))) * 2
		leftMost = (leftMost - bits.OnesCount64Before(bm, uint(leftMost))) * 2
		levelMask := ((uint64(1) << uint(i)) - 1) ^ ((uint64(1) << uint(leftMost)) - 1)
		leftLeavesMask |= levelMask
	}

	return bits.OnesCount64Before(bm&leftLeavesMask, 64)
}

func NewBitrie(s []string) Bitrie {
	sigbits := FindSignificantBits222(s)

	nodes := SigBitsBulidTrie(sigbits)
	// 1 flags
	// 8 byte bitmap
	bt := make([]byte, 1+8)

	p0 := 0
	// n-1 inner nodes, n leaf nodes
	positions := make([]int, 0, len(s)-1)

	// place root at 2nd bit to simplify the math
	for i, n := range nodes {
		j := i + 1
		// mark leaf as 1 thus right most 1 is the size of bitmap
		if n.Leaf {
			bt[1+j/8] |= 1 << uint(j%8)
		} else {
			positions = append(positions, n.Pos-p0)
			// if p0 == 0 {
			p0 = n.Pos
			// }
		}
	}

	fmt.Println("flag and bitmap:", bt)
	fmt.Println(positions)

	// max size for 64 var-int64
	pos := make([]byte, 64*9)
	n := 0
	for _, p := range positions {
		n += binary.PutVarint(pos[n:], int64(p))
	}

	pos = pos[:n]

	bbb := make([]byte, len(bt)+n)
	copy(bbb, bt)
	copy(bbb[len(bt):], pos)

	fmt.Println("result:", bbb)

	return Bitrie(bbb)
}

func NewSigTrie(sb []int) SigTrie {

	if len(sb) == 0 {
		return nil
	}

	nodes := SigBitsBulidTrie(sb)
	return SigTrie(nodes)
}

// SigBitsBulidTrie form a trie and return the root.
func SigBitsBulidTrie(sigbits []int) []TNode {

	children := make([]TNode, 0, len(sigbits)*2+1)
	children = append(children, TNode{
		frm: 0,
		to:  len(sigbits),
	})

	for i := 0; i < len(children); i++ {
		n := &children[i]

		if n.Leaf {
			continue
		}

		min := MinSigBit(sigbits, n.frm, n.to)
		n.Pos = sigbits[min]

		// fmt.Println("sigbits:", sigbits[n.frm:n.to])
		// fmt.Println("min:", min)

		var t TNode
		if min == n.frm {
			t = TNode{
				Leaf:      true,
				LeafValue: min,
			}
		} else {
			t = TNode{
				Leaf: false,
				frm:  n.frm,
				to:   min,
			}
		}
		n.L = len(children)
		children = append(children, t)
		// fmt.Printf("append left child: %#v\n", t)
		// fmt.Printf("children: %#v\n", children)

		if min == n.to-1 {
			t = TNode{
				Leaf:      true,
				LeafValue: min + 1,
			}
		} else {
			t = TNode{
				Leaf: false,
				frm:  min + 1,
				to:   n.to,
			}
		}
		n.R = len(children)
		children = append(children, t)
		// fmt.Printf("append right child: %#v\n", t)

		// fmt.Printf("processed: %#v\n", n)
		// fmt.Printf("children: %#v\n", children)
		// fmt.Println("current trie:", children)
	}

	return children
}

func FindSignificantBits222(s []string) []int {
	prev := s[0]
	sigbits := make([]int, 0)
	for _, str := range s[1:] {
		i := FirstDiff(prev, str)
		sigbits = append(sigbits, i)
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
