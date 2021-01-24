package trie

import (
	"bytes"
	"encoding/binary"
	fmt "fmt"
	"strings"

	"github.com/openacid/errors"
	"github.com/openacid/low/bitmap"
	"github.com/openacid/low/pbcmpl"
	"github.com/openacid/low/vers"
	"github.com/openacid/must"
	"github.com/openacid/slim/array"
)

// Marshal serializes it to byte stream.
//
// Since 0.4.3
func (st *SlimTrie) Marshal() ([]byte, error) {
	var buf []byte
	writer := bytes.NewBuffer(buf)

	_, err := pbcmpl.Marshal(writer, st.inner)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to marshal st.inner")
	}

	return writer.Bytes(), nil
}

// Unmarshal a SlimTrie from a byte stream.
//
// Since 0.4.3
func (st *SlimTrie) Unmarshal(buf []byte) error {

	st.inner = &Slim{}

	reader := bytes.NewReader(buf)

	_, h, err := pbcmpl.ReadHeader(reader)
	if err != nil {
		return errors.WithMessage(err, "failed to unmarshal header")
	}

	ver := h.GetVersion()
	compatible := st.compatibleVersions()

	if !vers.IsCompatible(ver, compatible) {
		return errors.Wrapf(ErrIncompatible,
			fmt.Sprintf(`version: "%s", compatible versions:"%s"`,
				ver,
				strings.Join(compatible, " || ")))
	}

	reader = bytes.NewReader(buf)

	// 0.5.10 and 0.5.11 share the same protobuf format:

	if vers.Check(ver, slimtrieVersion, "0.5.10") {
		_, _, err := pbcmpl.Unmarshal(reader, st.inner)
		if err != nil {
			return errors.WithMessage(err, "failed to unmarshal inner")
		}
		st.init()
		return nil
	}

	// ver: "==1.0.0 || <0.5.10"

	children := &array.Array32{}
	steps := &array.U16{}
	leaves := &array.Array{}
	leaves.EltEncoder = st.encoder

	_, _, err = pbcmpl.Unmarshal(reader, children)
	if err != nil {
		return errors.WithMessage(err, "failed to unmarshal children")
	}

	_, _, err = pbcmpl.Unmarshal(reader, steps)
	if err != nil {
		return errors.WithMessage(err, "failed to unmarshal steps")
	}

	_, _, err = pbcmpl.Unmarshal(reader, leaves)
	if err != nil {
		return errors.WithMessage(err, "failed to unmarshal leaves")
	}

	// backward compatible:

	before000510(st, ver, children, steps, leaves)

	return nil
}

// ProtoMessage implements proto.Message
//
// Since 0.4.3
func (st *SlimTrie) ProtoMessage() {}

// Reset implements proto.Message
//
// Since 0.4.3
func (st *SlimTrie) Reset() {
	st.inner = &Slim{}
	st.vars = nil
	st.levels = []levelInfo{{0, 0, 0, nil}}
}

func before000510(st *SlimTrie, ver string, ch *array.Array32, steps *array.U16, lvs *array.Array) {
	if !vers.Check(ver, "==1.0.0", "<0.5.10") {
		return
	}
	before000510ToNewChildrenArray(st, ver, ch, steps, lvs)
}

func before000510ToNewChildrenArray(st *SlimTrie, ver string, ch *array.Array32, steps *array.U16, lvs *array.Array) {

	// 1.0.0 is the initial version.
	// From 0.5.8 it starts writing version to marshaled data.
	// In 0.5.4 it starts using Bitmap to store Children elements.
	// But 0.5.4 marshals data with version == 1.0.0

	if vers.Check(ver, "==1.0.0", "<0.5.10") {

		// rebuild inner

		type eltType struct {
			oldid    int32
			step     int32
			leafOnly bool
		}

		// before 0.5.10 it stores steps only, no prefix
		c := newCreator(64, true, normalizeOpt(&Opt{}))

		// before 0.5.10 there is no big inner
		c.isBig = false

		queue := make([]*eltType, 0)
		elt := &eltType{
			oldid:    0,
			step:     getStepBefore000510(steps, 0),
			leafOnly: false,
		}
		queue = append(queue, elt)

		nextOldID := int32(1)

		for newid := int32(0); newid < int32(len(queue)); newid++ {
			qelt := queue[newid]
			oldid := qelt.oldid

			hasInner := bmhas(ch.Bitmaps, oldid)
			hasLeaf := bmhas(lvs.Bitmaps, oldid)

			// it could be an empty slimtrie.
			must.Be.True(hasInner || hasLeaf || (len(ch.Bitmaps) == 0 && len(lvs.Bitmaps) == 0))

			if qelt.leafOnly || (!hasInner && hasLeaf) {
				must.Be.OK(func() { bmhas(lvs.Bitmaps, oldid) })

				lv, found := lvs.GetBytes(oldid, st.encoder.GetEncodedSize(nil))
				must.Be.True(found)

				c.addLeaf(newid, lv)
				continue
			}

			if !hasInner {
				continue
			}

			// 16-bit bitmap is same with bmtree bitmap of size =16
			bm := getBM16Child(ch, oldid)

			if hasLeaf {
				// "" is a explicit branch in 0.5.10
				// add leaf node: empty string "" path

				// before 0.5.10, an inner node could also play as a leaf.
				// In 0.5.10, we need to separate them, by adding a
				// leaf-only sub node.
				queue = append(queue, &eltType{
					oldid:    oldid,
					leafOnly: true,
				})
			}

			for range bitmap.ToArray([]uint64{bm}) {
				ee := &eltType{
					oldid:    nextOldID,
					step:     getStepBefore000510(steps, nextOldID),
					leafOnly: false,
				}
				queue = append(queue, ee)
				nextOldID++
			}

			if hasLeaf {
				bm |= 1
			}

			bmidx := bitmap.ToArray([]uint64{bm})

			c.addInner(newid, bmidx, innerSize, qelt.step, "", -1)
		}

		ns := c.build()
		ns.Leaves = c.buildLeaves(nil)

		st.inner = ns
		st.init()
	}
}

func getStepBefore000510(steps *array.U16, nid int32) int32 {
	if bmhas(steps.Bitmaps, nid) {
		stp, found := steps.Get(nid)
		must.Be.True(found)

		// From 0.5.10 step does not include the count of the label word.
		stp--

		// From 0.5.10 step is in bit instead of in 4-bit.
		return int32(stp) * 4
	}
	return 0
}

// check existence for old un-exntended bitmap
func bmhas(bm []uint64, i int32) bool {
	return bitmap.SafeGet1(bm, i) == 1
}

func getBM16Child(ch *array.Array32, idx int32) uint64 {

	// There are two format with version 1.0.0:
	// Before 0.5.4 Child elements are in Elts, every elt is uint32:
	// 16 bit bitmap in lower 16 bit. and the rank in upper 16 bit.
	//
	// Since 0.5.4 Child elements are in BMElts, every child is a 16-bit bitmap

	endian := binary.LittleEndian

	eltIdx, bitset := bitmap.Rank64(ch.Bitmaps, ch.Offsets, idx)
	must.Be.True(bitset == 1, "node must be in children array")

	var bm uint64

	if ch.Flags&array.ArrayFlagIsBitmap == 0 {

		// load from Base.Elts

		v := endian.Uint32(ch.Elts[eltIdx*4:])
		bm = uint64(v & 0xffff)

	} else {

		// load from Base.BMElts

		bm = bitmap.Getw(ch.BMElts.Words, eltIdx, 16)
	}

	// add leaf bit, thus the size is 17
	return bm << 1
}
