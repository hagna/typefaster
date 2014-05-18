package typefaster

import "testing"

func TestInsertDup(t *testing.T) {
	tree := Tree{&node{"Root", "", nil}}
	tree.Insert(tree.Root, "T.AH.P", "top")
	tree.Insert(tree.Root, "T.AH.P", "top")
	if len(tree.Root.Children) > 1 {
		t.Fatalf("should only be one child")
	}
}

func TestInsert(t *testing.T) {
	tree := Tree{&node{"Root", "", nil}}
	tree.Insert(tree.Root, "T.AH.P", "top")
	tree.Insert(tree.Root, "T.AH.P.S", "tops")
	tree.Print(tree.Root)
	if len(tree.Root.Children) == 2 {
		t.Fatalf("should only two Children")
	}
}

func TestInsertSplit(t *testing.T) {
	tree := Tree{&node{"Root", "", nil}}
	tree.Insert(tree.Root, "T.AH.P", "top")
	tree.Insert(tree.Root, "T.AH.T", "tot")
	if len(tree.Root.Children) == 2 {
		t.Fatalf("should only two Children")
	}
}

func TestInsertMore(t *testing.T) {
	tree := Tree{&node{"Root", "", nil}}
	tree.Insert(tree.Root, "test", "top")
	tree.Insert(tree.Root, "slow", "top")
	tree.Insert(tree.Root, "water", "top")
	tree.Insert(tree.Root, "slower", "top")
	tree.Insert(tree.Root, "slowest", "top")
	tree.Print(tree.Root)
	if len(tree.Root.Children) == 2 {
		t.Fatalf("should only two Children")
	}
	
}
