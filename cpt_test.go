package typefaster

import (
	"strings"
	"testing"
)

func TestInsertDup(t *testing.T) {
	tree := Tree{NewNode("root", "", nil)}
	tree.Insert(tree.Root, "T.AH.P", "top")
	tree.Insert(tree.Root, "T.AH.P", "top")
	if len(tree.Root.Children) > 1 {
		t.Fatal("should only be one child")
	}
}

func TestInsert(t *testing.T) {
	tree := Tree{NewNode("root", "", nil)}
	tree.Insert(tree.Root, "T.AH.P", "top")
	tree.Insert(tree.Root, "T.AH.P.S", "tops")
	tree.Print(tree.Root, "")
	if len(tree.Root.Children) == 2 {
		t.Fatal("should only two Children")
	}
}

func TestInsertSplit(t *testing.T) {
	tree := Tree{NewNode("root", "", nil)}
	tree.Insert(tree.Root, "T.AH.P", "top")
	tree.Insert(tree.Root, "T.AH.T", "tot")
	if len(tree.Root.Children) == 2 {
		t.Fatal("should only two Children")
	}
}

func TestInsertMore(t *testing.T) {
	tree := Tree{NewNode("root", "", nil)}
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
	tree := Tree{NewNode("root", "", nil)}
	tree.Insert(tree.Root, "A", "top")
	tree.Insert(tree.Root, "Alpha", "top")
	tree.Insert(tree.Root, "Anaconda", "top")
	tree.Insert(tree.Root, "Al", "top")
	tree.Print(tree.Root, "")
	if len(tree.Root.Children) != 1 {
		t.Fatal("should have one child")
	}
	n := tree.Root.Children[0].GetNode()
	if n.Edgename != "A" {
		t.Fatal("should be A")
	}
	if len(n.Children) != 2 {
		t.Fatal("should have 2 Children of A, but we have", n.Children)
	}
}

func TestBug2(t *testing.T) {
	tree := Tree{NewNode("root", "", nil)}
	tree.Insert(tree.Root, "A", "top")
	tree.Insert(tree.Root, "Alpha", "top")
	// This next shouldn't be under Alpha
	tree.Insert(tree.Root, "Aughat", "top")
	tree.Insert(tree.Root, "Ao", "top")
	tree.Print(tree.Root, "")
	if len(tree.Root.Children) != 1 {
		t.Fatal("should have one child")
	}
	n := tree.Root.Children[0].GetNode()
	if n.Edgename != "A" {
		t.Fatal("should be A")
	}
	if len(n.Children) != 3 {
		t.Fatal("should have 3 Children of A, but we have", n.Children)
	}
}

func isFound(s string, tree Tree, t *testing.T) {
	_, _, c := tree.Lookup(tree.Root, s)
	tree.Print(tree.Root, "")
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
	tree := Tree{NewNode("root", "", nil)}
	for _, v := range s {
		tree.Insert(tree.Root, v, "")
	}
	for _, v := range s {
		isFound(v, tree, t)
	}
	if len(tree.Root.Children) != 1 {
		t.Fatal("should have one child", tree.Root.Children)
	}
	n := tree.Root.Children[0].GetNode()	
	if n.Edgename != "Water" {
		t.Fatal("should not be", n.Edgename)
	}
	if len(n.Children) != 1 {
		t.Fatal("wrong children count for", n.Children)
	}
	if tree.Root.Children[0].GetNode().Children[0].GetNode().Children[1].GetNode().Edgename != "g" {
		t.Fatal("wrong edgename", tree.Root.Children[0].GetNode().Children[0].GetNode().Children[1].GetNode().Edgename)
	}
	if tree.Root.Children[0].GetNode().Children[0].GetNode().Children[0].GetNode().Edgename != "k" {
		t.Fatal("wrong edgename", tree.Root.Children[0].GetNode().Children[0].GetNode().Children[0].GetNode().Edgename)
	}
	if tree.Root.Children[0].GetNode().Children[0].GetNode().Children[1].GetNode().Children[0].GetNode().Edgename != "s" {
		t.Fatal("wrong edgename", tree.Root.Children[0].GetNode().Children[0].GetNode().Children[1].GetNode().Children[0].GetNode().Edgename)
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
	tree := Tree{NewNode("root", "", nil)}
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
	tree := Tree{NewNode("root", "", nil)}
	for _, v := range s {
		tree.Insert(tree.Root, v, "")
	}
	_, _, c := tree.Lookup(tree.Root, "aardvark")
	tree.Print(tree.Root, "")
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
	tree := Tree{NewNode("root", "", nil)}
	for _, v := range s {
		tree.Insert(tree.Root, v, "")
	}
	called := 0
	callme := func(k, v []string) {
		t.Log(k)
		called += len(k)
	}
	tree.mkdir(tree.Root, []string{"root"}, callme)
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
	
