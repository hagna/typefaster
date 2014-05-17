package typefaster

import "testing"

func TestInsert(t *testing.T) {
	tree := tree{&node{"root", "", nil}}
	tree.insert(tree.root, "top", "T.AH.P")
	tree.print(tree.root)
}
