package trie_test

import (
	"fmt"
	"testing"

	"github.com/openacid/slim/trie"
)

type T struct {
	branches map[int][]int
}

func nodeID(node interface{}) int {
	if node == nil {
		return 0
	}
	return node.(int)
}

func (t *T) Child(node, branch interface{}) interface{} {
	if children, ok := t.branches[nodeID(node)]; ok {
		return children[branch.(int)]
	}
	return nil
}
func (t *T) Branches(node interface{}) []interface{} {
	rst := []interface{}{}
	if children, ok := t.branches[nodeID(node)]; ok {
		for i := 0; i < len(children); i++ {
			rst = append(rst, i)
		}
	}
	return rst
}

func (t *T) NodeID(node interface{}) string {
	return fmt.Sprintf("%02d", nodeID(node))
}

func (t *T) NodeInfo(node interface{}) string {
	return "(foo)"
}

func (t *T) LeafVal(node interface{}) (interface{}, bool) {
	if _, ok := t.branches[nodeID(node)]; ok {
		return nil, false
	}
	return "leaf", true
}

func TestToString(t *testing.T) {
	// 0 --> 1 --> 3 -->5
	//               `->6
	//   `-> 2 --> 4
	tree := &T{
		branches: map[int][]int{
			0: {1, 2},
			1: {3},
			2: {4},
			3: {5, 6},
		},
	}

	rst := trie.ToString(tree)
	want := `
#00(foo)*2
   -000->#01(foo)
            -000->#03(foo)*2
                     -000->#05(foo)=leaf
                     -001->#06(foo)=leaf
   -001->#02(foo)
            -000->#04(foo)=leaf`[1:]
	if want != rst {
		t.Fatalf("expect: \n%v\n; but: \n%v\n", want, rst)
	}
}
