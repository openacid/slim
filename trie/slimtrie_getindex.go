package trie

// Index oriented query is different from the basic one.
// It returns the original position of a leaf in the key-value array for
// creating a slimtrie.
// This is useful when a user need to separate storage for keys and values.

import "github.com/openacid/low/bitmap"

// GetIndex looks up for key and return the its position in the sorted kv list that is used to build slim.
// If no such key is found, it returns -1.
//
// Since 0.5.12
func (st *SlimTrie) GetIndex(key string) int32 {

	ns := st.inner

	if ns.NodeTypeBM == nil {
		return -1
	}

	l := int32(8 * len(key))
	qr := &querySession{
		keyBitLen: l,
		key:       key,
	}

	cur := &walkingCursor{id: 0, lvl: 1}

	keyIdx := int32(0)

	for {
		st.getNode(cur.id, qr)
		if qr.isInner == 0 {
			// leaf
			break
		}

		if qr.hasInnerPrefix {
			r := prefixCompare(key[keyIdx>>3:], qr.innerPrefix)
			if r != 0 {
				return -1
			}
			keyIdx = keyIdx&(^7) + qr.innerPrefixLen
		} else {
			keyIdx += qr.innerPrefixLen
		}

		if keyIdx > l {
			return -1
		}

		lchID, has := st.getLeftChildID(qr, keyIdx)
		if has == 0 {
			// no such branch of label
			return -1
		}
		cur.nextLevel(qr.ithInner, st, lchID+1)

		if keyIdx == l {
			// must be a leaf
			break
		}

		keyIdx += qr.wordSize
	}

	// currId must not be -1

	// if keyIdx == l the leaf does not have leaf prefix
	if keyIdx <= l {
		tail := key[keyIdx>>3:]
		// the quick path: break from `if keyIdx == l`, qr is old.
		r := st.cmpLeafPrefix(tail, qr)
		if r != 0 {
			return -1
		}
	}

	return st.cursorLeafIndex(cur, true)
}

// GetLRIndex looks up for two indexes l and r so that keys[l] <= key <= keys[r]
// If a exact match is found, it returns (l,l);
// If no exact match is found, it returns (l, l+1); l could be -1 and l+1 could be `len(keys)`
//
// Since 0.5.12
func (st *SlimTrie) GetLRIndex(key string) (int32, int32) {
	ns := st.inner

	if ns.NodeTypeBM == nil {
		return -1, 0
	}

	l := int32(8 * len(key))
	qr := &querySession{
		keyBitLen: l,
		key:       key,
	}

	leftCur := &walkingCursor{id: -1, lvl: -1}
	eqCur := &walkingCursor{id: 0, lvl: 1}

	keyIdx := int32(0)

	for {
		st.getNode(eqCur.id, qr)
		if qr.isInner == 0 {
			// leaf
			break
		}

		if qr.hasInnerPrefix {
			r := prefixCompare(key[keyIdx>>3:], qr.innerPrefix)
			if r == 0 {
				keyIdx = keyIdx&(^7) + qr.innerPrefixLen
			} else if r < 0 {
				// key < prefix
				eqCur.id = -1
				break
			} else {
				// key > prefix
				*leftCur = *eqCur
				eqCur.id = -1
				break
			}
		} else {
			keyIdx += qr.innerPrefixLen
		}

		if keyIdx > l {
			// same as key < prefix
			eqCur.id = -1
			break
		}

		leftChild, has := st.getLeftChildID(qr, keyIdx)

		// left most and right most child from this node
		leftMostChild, _ := bitmap.Rank128(ns.Inners.Words, ns.Inners.RankIndex, qr.from)
		leftMostChild++

		if leftChild >= leftMostChild {
			*leftCur = *eqCur
			leftCur.nextLevel(qr.ithInner, st, leftChild)
		}

		if has == 0 {
			eqCur.id = -1
			break
		}
		eqCur.nextLevel(qr.ithInner, st, leftChild+1)

		if keyIdx == l {
			// must be a leaf
			break
		}

		keyIdx += qr.wordSize
	}

	// currId must not be -1

	if eqCur.id != -1 {
		// if keyIdx == l the leaf does not have leaf prefix
		if keyIdx <= l {
			tail := key[keyIdx>>3:]
			// the quick path: break from `if keyIdx == l`, qr is old.
			r := st.cmpLeafPrefix(tail, qr)
			if r == -1 {
				// key < pref
				eqCur.id = -1
			} else if r == 1 {
				// key > pref
				*leftCur = *eqCur
				eqCur.id = -1
			}
		}
	}

	if eqCur.id != -1 {
		i := st.cursorLeafIndex(eqCur, true)
		return i, i
	}

	if leftCur.id != -1 {
		st.rightMostCursor(leftCur)
		i := st.cursorLeafIndex(leftCur, true)
		return i, i + 1
	}

	// key < all record
	return -1, 0
}
