package trie

import (
	"encoding/binary"
	"xec/sparse"
)

type children struct {
	Bitmap uint16
	Offset uint16
}

type SparseTrie struct {
	Children sparse.SparseArray
	Steps    sparse.SparseArray
	Leaves   sparse.SparseArray
}

const WordMask = 0xf

type ChildConv struct {
}

func (c ChildConv) MarshalElt(d interface{}) []byte {
	child := d.(children)

	b := make([]byte, 4)
	binary.LittleEndian.PutUint16(b[:2], child.Bitmap)
	binary.LittleEndian.PutUint16(b[2:4], child.Offset)

	return b
}

func (c ChildConv) UnmarshalElt(b []byte) (uint32, interface{}) {
	d := children{
		Bitmap: binary.LittleEndian.Uint16(b[:2]),
		Offset: binary.LittleEndian.Uint16(b[2:4]),
	}
	return uint32(4), d
}

func (c ChildConv) GetMarshaledEltSize(b []byte) uint32 {
	return uint32(4)
}

func (st *SparseTrie) FromSparse(root *Node) (err error) {
	if root == nil {
		return
	}

	childIndex, childData := []uint32{}, []children{}
	stepIndex, stepData := []uint32{}, []uint16{}
	leafIndex, leafData := []uint32{}, []interface{}{}

	tq := make([]*Node, 0, 256)
	tq = append(tq, root)

	for nId := uint16(0); ; {
		if len(tq) == 0 {
			break
		}

		node := tq[0]
		tq = tq[1:]

		if len(node.Branches) == 0 {
			continue
		}

		brs := node.Branches

		if brs[0] == leafBranch {
			leafIndex = append(leafIndex, uint32(nId))
			leafData = append(leafData, node.Children[brs[0]].Value)

			brs = brs[1:]
		}

		if node.Step > 1 {
			stepIndex = append(stepIndex, uint32(nId))
			stepData = append(stepData, node.Step)

		}

		if len(brs) > 0 {
			childIndex = append(childIndex, uint32(nId))
			offset := nId + uint16(len(tq)) + uint16(1)

			bitmap := uint16(0)
			for _, b := range brs {
				bitmap |= uint16(1) << (uint16(b) & WordMask)
			}

			ch := children{
				Bitmap: bitmap,
				Offset: offset,
			}

			childData = append(childData, ch)
		}

		for _, b := range brs {
			tq = append(tq, node.Children[b])
		}

		nId++
	}

	err = st.Children.Init(childIndex, childData)
	if err != nil {
		return err
	}

	err = st.Steps.Init(stepIndex, stepData)
	if err != nil {
		return err
	}

	err = st.Leaves.Init(leafIndex, leafData)
	if err != nil {
		return err
	}

	return nil
}

func (st *SparseTrie) Search(key []byte, mode Mode) (value interface{}) {
	return nil
}
