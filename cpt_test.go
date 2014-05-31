package typefaster

import (
	"strings"
	"testing"
)

func TestInsertDup(t *testing.T) {
	tree := NewMemTree("root")
	Insert(tree, "T.AH.P", "top")
	Insert(tree, "T.AH.P", "top")
	if len(tree.Root().Children) > 1 {
		t.Fatal("should only be one child")
	}
}

func TestInsert(t *testing.T) {
	tree := NewMemTree("root")
	Insert(tree, "T.AH.P", "top")
	Insert(tree, "T.AH.P.S", "tops")
	tree.Print(tree.Root(), "")
	if len(tree.Root().Children) == 2 {
		t.Fatal("should only two Children")
	}
}

func TestInsertSplit(t *testing.T) {
	tree := NewMemTree("root")
	Insert(tree, "T.AH.P", "top")
	Insert(tree, "T.AH.T", "tot")
	if len(tree.Root().Children) == 2 {
		t.Fatal("should only two Children")
	}
}

func TestInsertMore(t *testing.T) {
	tree := NewMemTree("root")
	Insert(tree, "test", "top")
	Insert(tree, "slow", "top")
	Insert(tree, "water", "top")
	Insert(tree, "slower", "top")
	Insert(tree, "slowest", "top")
	tree.Print(tree.Root(), "")
	if len(tree.Root().Children) == 2 {
		t.Fatal("should only two Children")
	}

}

func TestBug1(t *testing.T) {
	tree := NewMemTree("root")
	Insert(tree, "A", "top")
	Insert(tree, "Alpha", "top")
	Insert(tree, "Anaconda", "top")
	Insert(tree, "Al", "top")
	tree.Print(tree.Root(), "")
	if len(tree.Root().Children) != 1 {
		t.Fatal("should have one child")
	}
	n := tree.Root().Children[0]
	if n.Edgename != "A" {
		t.Fatal("should be A")
	}
	if len(n.Children) != 2 {
		t.Fatal("should have 2 Children of A, but we have", n.Children)
	}
}

func TestBug2(t *testing.T) {
	tree := NewMemTree("root")
	Insert(tree, "A", "top")
	Insert(tree, "Alpha", "top")
	// This next shouldn't be under Alpha
	Insert(tree, "Aughat", "top")
	Insert(tree, "Ao", "top")
	tree.Print(tree.Root(), "")
	if len(tree.Root().Children) != 1 {
		t.Fatal("should have one child")
	}
	n := tree.Root().Children[0]
	if n.Edgename != "A" {
		t.Fatal("should be A")
	}
	if len(n.Children) != 3 {
		t.Fatal("should have 3 Children of A, but we have", n.Children)
	}
}

func isFound(s string, tree MemTree, t *testing.T) {
	_, _, c := tree.Lookup(tree.Root(), s)
	tree.Print(tree.Root(), "")
	if c != s {
		t.Fatal("didn't find word", s, "but found", c)
	}
}

func TestBug3(t *testing.T) {
	l := `Water
Watering
Waterings
Waterink`
	s := strings.Split(l, "\n")
	tree := NewMemTree("root")
	for _, v := range s {
		Insert(tree, v, "")
	}
	for _, v := range s {
		isFound(v, *tree, t)
	}
	if len(tree.Root().Children) != 1 {
		t.Fatal("should have one child", tree.Root().Children)
	}
	n := tree.Root().Children[0]	
	if n.Edgename != "Water" {
		t.Fatal("should not be", n.Edgename)
	}
	if len(n.Children) != 1 {
		t.Fatal("wrong children count for", n.Children)
	}
	if tree.Root().Children[0].Children[0].Children[1].Edgename != "g" {
		t.Fatal("wrong edgename", tree.Root().Children[0].Children[0].Children[1].Edgename)
	}
	if tree.Root().Children[0].Children[0].Children[0].Edgename != "k" {
		t.Fatal("wrong edgename", tree.Root().Children[0].Children[0].Children[0].Edgename)
	}
	if tree.Root().Children[0].Children[0].Children[1].Children[0].Edgename != "s" {
		t.Fatal("wrong edgename", tree.Root().Children[0].Children[0].Children[1].Children[0].Edgename)
	}
}

func TestBug4(t *testing.T) {
	l := `abstain
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
abstruse`
	s := strings.Split(l, "\n")
	tree := NewMemTree("root")
	for _, v := range s {
		Insert(tree, v, "")
	}
	_, _, c := tree.Lookup(tree.Root(), "abstention")
	tree.Print(tree.Root(), "")
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
	tree := NewMemTree("root")
	for _, v := range s {
		Insert(tree, v, "")
	}
	_, _, c := tree.Lookup(tree.Root(), "aardvark")
	tree.Print(tree.Root(), "")
	if c != "aardvark" {
		t.Fatal("didn't find word but found", c)
	}
}

func TestMkdir(t *testing.T) {
	l := `abstain
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
abstruse`
	s := strings.Split(l, "\n")
	tree := NewMemTree("root")
	for _, v := range s {
		Insert(tree, v, "")
	}
	called := 0
	callme := func(k, v []string) {
		t.Log(k)
		called += len(k)
	}
	tree.mkdir(tree.Root(), []string{"root"}, callme)
	x := 41
	if x != called {
		t.Fatal("called", called, "times but should have been", x)
	}
}

func TestEncodeDecode(t *testing.T) {
	s := "AA.AE.AH.AO.AW.AY.B.CH.D.DH.EH.ER.EY.F.G.HH.IH.IY.JH.K.L.M.N.NG.OW.OY.P.R.S.SH.T.TH.UH.UW.V.W.Y.Z.ZH"
	if s != decode(encode(s)) {
		t.Log("encoded == ", encode(s))
		t.Log("decoded == ", decode(encode(s)))
		t.Fatal("failed to decode the encoded string")
		
	}
}
	
