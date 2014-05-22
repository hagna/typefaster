package typefaster

import (
	"strings"
	"testing"
)

func TestInsertDup(t *testing.T) {
	tree := Tree{&node{"Root", "", nil}}
	tree.Insert(tree.Root, "T.AH.P", "top")
	tree.Insert(tree.Root, "T.AH.P", "top")
	if len(tree.Root.Children) > 1 {
		t.Fatal("should only be one child")
	}
}

func TestInsert(t *testing.T) {
	tree := Tree{&node{"Root", "", nil}}
	tree.Insert(tree.Root, "T.AH.P", "top")
	tree.Insert(tree.Root, "T.AH.P.S", "tops")
	tree.Print(tree.Root, "")
	if len(tree.Root.Children) == 2 {
		t.Fatal("should only two Children")
	}
}

func TestInsertSplit(t *testing.T) {
	tree := Tree{&node{"Root", "", nil}}
	tree.Insert(tree.Root, "T.AH.P", "top")
	tree.Insert(tree.Root, "T.AH.T", "tot")
	if len(tree.Root.Children) == 2 {
		t.Fatal("should only two Children")
	}
}

func TestInsertMore(t *testing.T) {
	tree := Tree{&node{"Root", "", nil}}
	tree.Insert(tree.Root, "test", "top")
	tree.Insert(tree.Root, "slow", "top")
	tree.Insert(tree.Root, "water", "top")
	tree.Insert(tree.Root, "slower", "top")
	tree.Insert(tree.Root, "slowest", "top")
	tree.Print(tree.Root, "")
	if len(tree.Root.Children) == 2 {
		t.Fatal("should only two Children")
	}

}

func TestBug1(t *testing.T) {
	tree := Tree{&node{"Root", "", nil}}
	tree.Insert(tree.Root, "A", "top")
	tree.Insert(tree.Root, "Alpha", "top")
	tree.Insert(tree.Root, "Anaconda", "top")
	tree.Insert(tree.Root, "Al", "top")
	tree.Print(tree.Root, "")
	if len(tree.Root.Children) != 1 {
		t.Fatal("should have one child")
	}
	if tree.Root.Children[0].Edgename != "A" {
		t.Fatal("should be A")
	}
	if len(tree.Root.Children[0].Children) != 2 {
		t.Fatal("should have 2 Children of A, but we have", tree.Root.Children[0].Children)
	}
}

func TestBug2(t *testing.T) {
	tree := Tree{&node{"Root", "", nil}}
	tree.Insert(tree.Root, "A", "top")
	tree.Insert(tree.Root, "Alpha", "top")
	// This next shouldn't be under Alpha
	tree.Insert(tree.Root, "Aughat", "top")
	tree.Insert(tree.Root, "Ao", "top")
	tree.Print(tree.Root, "")
	if len(tree.Root.Children) != 1 {
		t.Fatal("should have one child")
	}
	if tree.Root.Children[0].Edgename != "A" {
		t.Fatal("should be A")
	}
	if len(tree.Root.Children[0].Children) != 3 {
		t.Fatal("should have 3 Children of A, but we have", tree.Root.Children[0].Children)
	}
}

func showChildren(t *Tree, test *testing.T) {
	for i, v := range t.Root.Children {
		test.Log(i, v)
	}
}

func TestBug3(t *testing.T) {
	tree := Tree{&node{"Root", "", nil}}
	tree.Insert(tree.Root, "Water", "")
	tree.Insert(tree.Root, "Watering", "")
	tree.Insert(tree.Root, "Waterings", "")
	tree.Insert(tree.Root, "Waterink", "")
	tree.Print(tree.Root, "")
	if len(tree.Root.Children) != 1 {
		t.Fatal("should have one child")
	}
	if tree.Root.Children[0].Edgename != "Water" {
		t.Fatal("should not be", tree.Root.Children[0].Edgename)
	}
	if len(tree.Root.Children[0].Children) != 1 {
		t.Fatal("wrong children count for", tree.Root.Children[0].Children)
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

func TestBug4(t *testing.T) {
	l := `
abstain
abstained
abstaining
abstention
abstentions
abstinence
abstinent
abstract
abstracted
abstraction
abstractions
abstracts
abstruse
	`
	s := strings.Split(l, "\n")
	tree := Tree{&node{"Root", "", nil}}
	for _, v := range s {
		tree.Insert(tree.Root, v, "")
	}
	_, _, c := tree.Lookup(tree.Root, "abstention")
	tree.Print(tree.Root, "")
	if c != "abstention" {
		t.Fatal("didn't find word but found", c)
	}
}

func TestBugAAA(t *testing.T) {
	l := `
a
aaa
aardvark
aaron
`
	s := strings.Split(l, "\n")
	tree := Tree{&node{"Root", "", nil}}
	for _, v := range s {
		tree.Insert(tree.Root, v, "")
	}
	_, _, c := tree.Lookup(tree.Root, "aardvark")
	tree.Print(tree.Root, "")
	if c != "aardvark" {
		t.Fatal("didn't find word but found", c)
	}
}
