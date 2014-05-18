package typefaster

import "testing"

func TestInsertDup(t *testing.T) {
	tree := tree{&node{"root", "", nil}}
	tree.insert(tree.root, "T.AH.P", "top")
	tree.insert(tree.root, "T.AH.P", "top")
	if len(tree.root.children) > 1 {
		t.Fatalf("should only be one child")
	}
}

func TestInsert(t *testing.T) {
	tree := tree{&node{"root", "", nil}}
	tree.insert(tree.root, "T.AH.P", "top")
	tree.insert(tree.root, "T.AH.P.S", "tops")
	tree.print(tree.root)
	if len(tree.root.children) == 2 {
		t.Fatalf("should only two children")
	}
}

func TestInsertSplit(t *testing.T) {
	tree := tree{&node{"root", "", nil}}
	tree.insert(tree.root, "T.AH.P", "top")
	tree.insert(tree.root, "T.AH.T", "tot")
	if len(tree.root.children) == 2 {
		t.Fatalf("should only two children")
	}
}
