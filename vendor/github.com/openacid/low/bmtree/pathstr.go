package bmtree

import (
	"fmt"
)

// PathStr creates a human-readable string of a node searching path.
//
// Since 0.1.9
func PathStr(path uint64) string {

	treeHeight := PathHeight(path)

	l := PathLen(path)
	if l == 0 {
		return ""
	}

	return fmt.Sprintf("%0[1]*[2]b", l, path>>uint(32+treeHeight-l))
}
