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


func TestBug1(t *testing.T) {
	tree := Tree{&node{"Root", "", nil}}
	tree.Insert(tree.Root, "A", "top")
	tree.Insert(tree.Root, "Alpha", "top")
	tree.Insert(tree.Root, "Anaconda", "top")
	tree.Insert(tree.Root, "Al", "top")
	tree.Print(tree.Root)
	if len(tree.Root.Children) != 1 {
		t.Fatalf("should have one child")
	}
	if tree.Root.Children[0].Edgename != "A" {
		t.Fatalf("should be A")
	}
	if len(tree.Root.Children[0].Children) != 2 {
		t.Fatalf("should have 2 Children of A, but we have", tree.Root.Children[0].Children)
	}
}

func TestBug2(t *testing.T) {
	tree := Tree{&node{"Root", "", nil}}
	tree.Insert(tree.Root, "A", "top")
	tree.Insert(tree.Root, "Alpha", "top")
	// This next shouldn't be under Alpha
	tree.Insert(tree.Root, "Aughat", "top")
	tree.Insert(tree.Root, "Ao", "top")
	tree.Print(tree.Root)
	if len(tree.Root.Children) != 1 {
		t.Fatalf("should have one child")
	}
	if tree.Root.Children[0].Edgename != "A" {
		t.Fatalf("should be A")
	}
	if len(tree.Root.Children[0].Children) != 3 {
		t.Fatalf("should have 3 Children of A, but we have", tree.Root.Children[0].Children)
	}
}

func TestBug3(t *testing.T) {
	tree := Tree{&node{"Root", "", nil}}
	tree.Insert(tree.Root, "Water", "")
	tree.Insert(tree.Root, "Watering", "")
	tree.Insert(tree.Root, "Waterings", "")
	tree.Insert(tree.Root, "Waterink", "")
	tree.Print(tree.Root)
	if len(tree.Root.Children) != 1 {
		t.Fatalf("should have one child")
	}
	if tree.Root.Children[0].Edgename != "Water" {
		t.Fatalf("should not be", tree.Root.Children[0].Edgename)
	}
	if len(tree.Root.Children[0].Children) != 1 {
		t.Fatalf("wrong children count for", tree.Root.Children[0].Children)
	}
	if tree.Root.Children[0].Children[0].Children[1].Edgename != "g" {
		t.Fatal("wrong edgename", tree.Root.Children[0].Children[0].Children[1].Edgename)
	}
	if tree.Root.Children[0].Children[0].Children[0].Edgename != "k" {
		t.Fatal("wrong edgename", tree.Root.Children[0].Children[0].Children[0].Edgename)
	}
	if tree.Root.Children[0].Children[0].Children[1].Children[0].Edgename != "s" {
		t.Fatal("wrong edgename", tree.Root.Children[0].Children[0].Children[1].Children[0].Edgename)
	}
}


