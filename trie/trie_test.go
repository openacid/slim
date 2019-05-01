package trie

import (
	"reflect"
	"testing"

	"github.com/openacid/errors"
	"github.com/openacid/slim/benchhelper"
)

func TestTrie(t *testing.T) {

	type ExpectType struct {
		key []byte
		rst int
	}

	var cases = []struct {
		key      [][]byte
		value    interface{}
		expected []ExpectType
	}{
		{
			key: [][]byte{
				{'a', 'b', 'c'},
				{'a', 'b', 'd'},
				{'b', 'c', 'd'},
				{'b', 'c', 'e'},
				{'c', 'd', 'e'},
			},
			value: []int{
				0,
				1,
				2,
				3,
				4,
			},
			expected: []ExpectType{
				{[]byte{'a', 'b', 'c'}, 0},
				{[]byte{'a', 'b', 'd'}, 1},
				{[]byte{'b', 'c', 'd'}, 2},
				{[]byte{'b', 'c', 'e'}, 3},
				{[]byte{'c', 'd', 'e'}, 4},
			},
		},
		{
			key: [][]byte{
				{'a', 'b', 'c'},
				{'a', 'b', 'c', 'd'},
				{'b', 'c'},
				{'b', 'c', 'd'},
				{'b', 'c', 'd', 'e'},
			},
			value: []int{
				0,
				1,
				2,
				3,
				4,
			},
			expected: []ExpectType{
				{[]byte{'a', 'b', 'c'}, 0},
				{[]byte{'a', 'b', 'c', 'd'}, 1},
				{[]byte{'b', 'c'}, 2},
				{[]byte{'b', 'c', 'd'}, 3},
				{[]byte{'b', 'c', 'd', 'e'}, 4},
			},
		},
	}

	for _, c := range cases {

		trie, _ := NewTrie(c.key, c.value, false)
		for _, ex := range c.expected {
			_, rst, _ := trie.Search(ex.key)

			if !reflect.DeepEqual(ex.rst, rst) {
				t.Fatal("ks: ", ex.key, "expected value: ", ex.rst, "rst: ", rst)
			}
		}

		trie.Squash()
		for _, ex := range c.expected {
			_, rst, _ := trie.Search(ex.key)

			if !reflect.DeepEqual(ex.rst, rst) {
				t.Error("ks: ", ex.key, "expected value: ", ex.rst, "rst: ", rst)
			}
		}

		trie.Squash()
		trie.Squash()
		for _, ex := range c.expected {
			_, rst, _ := trie.Search(ex.key)

			if !reflect.DeepEqual(ex.rst, rst) {
				t.Error("ks: ", ex.key, "expected value: ", ex.rst, "rst: ", rst)
			}
		}

		trie, _ = NewTrie(c.key, c.value, true)
		for _, ex := range c.expected {
			_, rst, _ := trie.Search(ex.key)

			if !reflect.DeepEqual(ex.rst, rst) {
				t.Error("ks: ", ex.key, "expected value: ", ex.rst, "rst: ", rst)
			}
		}
	}
}

func TestTrieSearch(t *testing.T) {

	var key = [][]byte{
		{'a', 'b', 'c'},
		{'a', 'b', 'c', 'd'},
		{'a', 'b', 'd'},
		{'a', 'b', 'd', 'e'},
		{'b', 'c'},
		{'b', 'c', 'd'},
		{'b', 'c', 'd', 'e'},
		{'c', 'd', 'e'},
	}
	var value = []int{
		0,
		1,
		2,
		3,
		4,
		5,
		6,
		7,
	}

	var trie, _ = NewTrie(key, value, false)

	type ExpectType struct {
		ltVal interface{}
		eqVal interface{}
		gtVal interface{}
	}

	var cases = []struct {
		key      []byte
		expected ExpectType
	}{
		{
			[]byte{'a', 'b', 'c'},
			ExpectType{nil, 0, 1},
		},
		{
			[]byte{'a', 'b', 'd'},
			ExpectType{1, 2, 3},
		},
		{
			[]byte{'b', 'c', 'd'},
			ExpectType{4, 5, 6},
		},
		{
			[]byte{'b', 'c', 'e'},
			ExpectType{6, nil, 7},
		},
		{
			[]byte{'c', 'd', 'e'},
			ExpectType{6, 7, nil},
		},
		{
			[]byte{'a', 'c', 'b'},
			ExpectType{3, nil, 4},
		},
		{
			[]byte{90, 'a', 'v'},
			ExpectType{nil, nil, 0},
		},
		{
			[]byte{'a', 'b'},
			ExpectType{nil, nil, 0},
		},
		{
			[]byte{'a', 'c'},
			ExpectType{3, nil, 4},
		},
		{
			[]byte{90},
			ExpectType{nil, nil, 0},
		},
		{
			[]byte{'a', 'b', 'c', 'd', 'e'},
			ExpectType{1, nil, 2},
		},
	}

	for _, c := range cases {

		ltVal, eqVal, gtVal := trie.Search(c.key)

		if !reflect.DeepEqual(c.expected.ltVal, ltVal) {
			t.Error("key: ", c.key, "expected lt value: ", c.expected.ltVal, "rst: ", ltVal)
		}
		if !reflect.DeepEqual(c.expected.eqVal, eqVal) {
			t.Error("key: ", c.key, "expected eq value: ", c.expected.eqVal, "rst: ", eqVal)
		}
		if !reflect.DeepEqual(c.expected.gtVal, gtVal) {
			t.Error("key: ", c.key, "expected gt value: ", c.expected.gtVal, "rst: ", gtVal)
		}
	}

	var squashedCases = []struct {
		key      []byte
		expected ExpectType
	}{
		{
			[]byte{'a', 'b'},
			ExpectType{nil, nil, 0},
		},
		{
			[]byte{'a', 'b', 'c'},
			ExpectType{nil, 0, 1},
		},
		{
			[]byte{'a', 'd', 'c'},
			ExpectType{nil, 0, 1},
		},
		{
			[]byte{'a', 'b', 'd'},
			ExpectType{1, 2, 3},
		},
		{
			[]byte{'a', 'c', 'd'},
			ExpectType{1, 2, 3},
		},
		{
			[]byte{'c', 'd', 'e'},
			ExpectType{6, 7, nil},
		},
		{
			[]byte{'c', 'f', 'e'},
			ExpectType{6, 7, nil},
		},
		{
			[]byte{'c', 'f', 'f'},
			ExpectType{6, 7, nil},
		},
		{
			[]byte{'c'},
			ExpectType{6, nil, 7},
		},
		{
			[]byte{'a', 'c'},
			ExpectType{nil, nil, 0},
		},
		{
			[]byte{'a', 'b', 'c', 'd', 'e'},
			ExpectType{1, nil, 2},
		},
	}

	trie.Squash()
	for _, c := range squashedCases {

		ltVal, eqVal, gtVal := trie.Search(c.key)

		if !reflect.DeepEqual(c.expected.ltVal, ltVal) {
			t.Error("key: ", c.key, "expected lt value: ", c.expected.ltVal, "rst: ", ltVal)
		}
		if !reflect.DeepEqual(c.expected.eqVal, eqVal) {
			t.Error("key: ", c.key, "expected eq value: ", c.expected.eqVal, "rst: ", eqVal)
		}
		if !reflect.DeepEqual(c.expected.gtVal, gtVal) {
			t.Error("key: ", c.key, "expected gt value: ", c.expected.gtVal, "rst: ", gtVal)
		}
	}
}

func TestTrieNew(t *testing.T) {

	var cases = []struct {
		keys        [][]byte
		values      interface{}
		expectedErr error
	}{
		{
			keys: [][]byte{
				{'a', 'b', 'c'},
				{'b', 'c', 'd'},
				{'c', 'd', 'e'},
			},
			values:      []int{0, 1, 2},
			expectedErr: nil,
		},
		{
			keys: [][]byte{
				{'a', 'b', 'c'},
				{'a', 'b', 'c'},
				{'c', 'd', 'e'},
			},
			values:      []int{0, 1, 2},
			expectedErr: errors.Wrapf(ErrDuplicateKeys, "key: abc"),
		},
	}
	benchhelper.WantPanic(t, func() {
		NewTrie([][]byte{{'a', 'b', 'c'}}, map[string]int{"abc": 0}, false)
	}, "values is map")

	for _, c := range cases {
		_, err := NewTrie(c.keys, c.values, false)
		if errors.Cause(err) != errors.Cause(c.expectedErr) {
			t.Fatalf("new trie: expectedErr: %v, got: %v", c.expectedErr, err)
		}
	}
}

func TestNewError(t *testing.T) {

	benchhelper.WantPanic(t, func() { NewTrie([][]byte{}, nil, false) }, "values is nil")
	benchhelper.WantPanic(t, func() { NewTrie([][]byte{{1}}, 1, false) }, "values is int")

	cases := []struct {
		keys    [][]byte
		vals    interface{}
		wanterr error
	}{
		{nil, nil, nil},
		{nil, []int{}, nil},
		{nil, 1, nil},
		{[][]byte{}, []int{1}, ErrKVLenNotMatch},
		{[][]byte{{1}}, []int{}, ErrKVLenNotMatch},
		{[][]byte{{1, 2}, {1}}, []int{1, 2}, ErrKeyOutOfOrder},
		{[][]byte{{1, 2}, {1, 1}}, []int{1, 2}, ErrKeyOutOfOrder},
		{[][]byte{{1, 2}, {1, 2}}, []int{1, 2}, ErrDuplicateKeys},
	}

	for i, c := range cases {
		_, err := NewTrie(c.keys, c.vals, false)
		if c.wanterr != errors.Cause(err) {
			t.Fatalf("%d-th: input: keys:%v; vals: %v; wanterr: %v; actual: %v",
				i+1, c.keys, c.vals, c.wanterr, err)
		}
	}
}

func TestAppend(t *testing.T) {

	tr, err := NewTrie([][]byte{{2, 3}, {2, 5}}, []int{1, 2}, false)
	if err != nil {
		t.Fatalf("expect no error but: %v", err)
	}

	cases := []struct {
		keys    []byte
		wanterr error
	}{
		{[]byte{1}, ErrKeyOutOfOrder},
		{[]byte{1, 2}, ErrKeyOutOfOrder},
		{[]byte{2}, ErrKeyOutOfOrder},
		{[]byte{2, 2}, ErrKeyOutOfOrder},
		{[]byte{2, 3}, ErrDuplicateKeys},
		{[]byte{2, 4}, ErrKeyOutOfOrder},
		{[]byte{2, 5}, ErrDuplicateKeys},
		{[]byte{2, 6}, nil},
	}

	for i, c := range cases {
		_, err := tr.Append(c.keys, 1)
		if c.wanterr != errors.Cause(err) {
			t.Fatalf("%d-th: input: %v; want: %v; actual: %v",
				i+1, c.keys, c.wanterr, err)
		}
	}
}

func TestAppendOneKeySquash(t *testing.T) {

	keys := [][]byte{
		{1, 2, 3},
	}
	values := []int{
		3,
	}

	type ExpectType struct {
		key   []byte
		value interface{}
	}

	rt, err := NewTrie(keys, values, true)
	if err != nil {
		t.Error("NewTrie failed: ", err)
	}

	cases := []ExpectType{
		{[]byte{1, 2, 2}, 3},
		{[]byte{1, 2, 3}, 3},
		{[]byte{1, 3, 3}, 3},
		{[]byte{3, 3, 3}, 3},
		{[]byte{4, 4, 4}, 3},
		{[]byte{1, 2}, nil},
		{[]byte{1, 2, 3, 4}, nil},
	}

	for _, c := range cases {
		_, rst, _ := rt.Search(c.key)
		if rst != c.value {
			t.Error("ks: ", c.key, "expected value: ", c.value, "rst: ", rst)
		}
	}

}

func TestNewTrieSquash(t *testing.T) {
	keys := [][]byte{
		{1, 2, 3, 4, 0},
		{1, 2, 3, 4, 1},
		{1, 2, 3, 4, 2},
		{1, 2, 3, 4, 3},
		{1, 3, 3, 5, 4},
	}

	values := []int{
		0,
		1,
		2,
		3,
		4,
	}

	type ExpectType struct {
		key   []byte
		value interface{}
	}

	// NewTrie with squash==true
	rt, err := NewTrie(keys, values, true)
	if err != nil {
		t.Error("NewTrie failed: ", err)
	}

	cases := []ExpectType{
		{[]byte{1, 2, 3, 4, 0}, 0},
		{[]byte{1, 2, 3, 5, 0}, 0},
		{[]byte{1, 2, 3, 4, 1}, 1},
		{[]byte{1, 2, 3, 5, 1}, 1},
		{[]byte{1, 2, 3, 4, 2}, 2},
		{[]byte{1, 2, 3, 5, 2}, 2},
		{[]byte{1, 2, 3, 4, 3}, 3},
		{[]byte{1, 2, 3, 5, 3}, 3},
		{[]byte{1, 3, 3, 5, 4}, 4},
	}

	for _, c := range cases {
		_, rst, _ := rt.Search(c.key)
		if rst != c.value {
			t.Error("ks: ", c.key, "expected value: ", c.value, "rst: ", rst)
		}
	}

	// NewTrie with squash==false
	rt, err = NewTrie(keys, values, false)
	if err != nil {
		t.Error("NewTri failed: ", err)
	}

	cases = []ExpectType{
		{[]byte{1, 2, 3, 4, 0}, 0},
		{[]byte{1, 2, 3, 5, 0}, nil},
		{[]byte{1, 2, 3, 4, 1}, 1},
		{[]byte{1, 2, 3, 5, 1}, nil},
		{[]byte{1, 2, 3, 4, 2}, 2},
		{[]byte{1, 2, 3, 5, 2}, nil},
		{[]byte{1, 2, 3, 4, 3}, 3},
		{[]byte{1, 2, 3, 5, 3}, nil},
		{[]byte{1, 3, 3, 5, 4}, 4},
	}

	for _, c := range cases {
		_, rst, _ := rt.Search(c.key)
		if rst != c.value {
			t.Error("ks: ", c.key, "expected value: ", c.value, "rst: ", rst)
		}
	}
}

func TestToStrings(t *testing.T) {
	var keys = [][]byte{
		{'a', 'b', 'c'},
		{'a', 'b', 'c', 'd'},
		{'a', 'b', 'd'},
		{'a', 'b', 'd', 'e'},
		{'b', 'c'},
		{'b', 'c', 'd'},
		{'b', 'c', 'd', 'e'},
		{'c', 'd', 'e'},
	}
	var values = []int{0, 1, 2, 3, 4, 5, 6, 7}

	expect := `
*3
-097->
      -098->*2
            -099->*2
                  -00$->=0
                  -100->
                        -00$->=1
            -100->*2
                  -00$->=2
                  -101->
                        -00$->=3
-098->
      -099->*2
            -00$->=4
            -100->*2
                  -00$->=5
                  -101->
                        -00$->=6
-099->
      -100->
            -101->
                  -00$->=7`[1:]

	var trie, _ = NewTrie(keys, values, false)
	trie, err := NewTrie(keys, values, false)
	if err != nil {
		t.Fatalf("expect no err: %s", err)
	}

	if expect != trie.String() {
		t.Fatalf("expect: \n%v\n; but: \n%v\n", expect, trie.String())
	}
}
