package oid

import (
	"fmt"
	"testing"
)

func createTestTree() (t *Tree) {
	t = NewTree()
	t.AddNode([]uint32{1, 3, 6, 1, 4, 1, 9999, 1})
	t.AddNode([]uint32{1, 3, 6, 1, 4, 1, 9999, 2})
	t.AddNode([]uint32{1, 3, 6, 1, 4, 1, 9999, 3})
	return
}

func TestForRangeLeaves(t *testing.T) {
	tree := createTestTree()
	tree.ForRangeLeaves(func(oid []uint32, node *Node) {
		fmt.Println(String(oid))
	})
}
