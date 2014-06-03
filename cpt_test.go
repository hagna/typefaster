package typefaster

import (
	"strings"
	"testing"
	"encoding/json"
)

func TestInsertDup(t *testing.T) {
	tree := NewMemTree("root")
	tree.Insert("T.AH.P", "top")
	tree.Insert("T.AH.P", "top")
	if len(tree.root.Children) > 1 {
		t.Fatal("should only be one child")
	}
}

func TestInsert(t *testing.T) {
	tree := NewMemTree("root")
	tree.Insert("T.AH.P", "top")
	tree.Insert("T.AH.P.S", "tops")
	tree.Print(tree.root, "")
	if len(tree.root.Children) == 2 {
		t.Fatal("should only two Children")
	}
}

func TestInsertSplit(t *testing.T) {
	tree := NewMemTree("root")
	tree.Insert("T.AH.P", "top")
	tree.Insert("T.AH.T", "tot")
	if len(tree.root.Children) == 2 {
		t.Fatal("should only two Children")
	}
}

func TestInsertMore(t *testing.T) {
	tree := NewMemTree("root")
	tree.Insert("test", "top")
	tree.Insert("slow", "top")
	tree.Insert("water", "top")
	tree.Insert("slower", "top")
	tree.Insert("slowest", "top")
	tree.Print(tree.root, "")
	if len(tree.root.Children) == 2 {
		t.Fatal("should only two Children")
	}

}

func TestBug1(t *testing.T) {
	tree := NewMemTree("root")
	tree.Insert("A", "top")
	tree.Insert("Alpha", "top")
	tree.Insert("Anaconda", "top")
	tree.Insert("Al", "top")
	tree.Print(tree.root, "")
	if len(tree.root.Children) != 1 {
		t.Fatal("should have one child")
	}
	n := tree.root.Children[0]
	if n.Edgename != "A" {
		t.Fatal("should be A")
	}
	if len(n.Children) != 2 {
		t.Fatal("should have 2 Children of A, but we have", n.Children)
	}
}

func TestBug2(t *testing.T) {
	tree := NewMemTree("root")
	tree.Insert("A", "top")
	tree.Insert("Alpha", "top")
	// This next shouldn't be under Alpha
	tree.Insert("Aughat", "top")
	tree.Insert("Ao", "top")
	tree.Print(tree.root, "")
	if len(tree.root.Children) != 1 {
		t.Fatal("should have one child")
	}
	n := tree.root.Children[0]
	if n.Edgename != "A" {
		t.Fatal("should be A")
	}
	if len(n.Children) != 3 {
		t.Fatal("should have 3 Children of A, but we have", n.Children)
	}
}

func isFound(s string, tree MemTree, t *testing.T) {
	_, _, c := tree.Lookup(tree.root, s)
	tree.Print(tree.root, "")
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
		tree.Insert(v, "")
	}
	for _, v := range s {
		isFound(v, *tree, t)
	}
	if len(tree.root.Children) != 1 {
		t.Fatal("should have one child", tree.root.Children)
	}
	n := tree.root.Children[0]	
	if n.Edgename != "Water" {
		t.Fatal("should not be", n.Edgename)
	}
	if len(n.Children) != 1 {
		t.Fatal("wrong children count for", n.Children)
	}
	if tree.root.Children[0].Children[0].Children[1].Edgename != "g" {
		t.Fatal("wrong edgename", tree.root.Children[0].Children[0].Children[1].Edgename)
	}
	if tree.root.Children[0].Children[0].Children[0].Edgename != "k" {
		t.Fatal("wrong edgename", tree.root.Children[0].Children[0].Children[0].Edgename)
	}
	if tree.root.Children[0].Children[0].Children[1].Children[0].Edgename != "s" {
		t.Fatal("wrong edgename", tree.root.Children[0].Children[0].Children[1].Children[0].Edgename)
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
		tree.Insert(v, "")
	}
	_, _, c := tree.Lookup(tree.root, "abstention")
	tree.Print(tree.root, "")
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
		tree.Insert(v, "")
	}
	_, _, c := tree.Lookup(tree.root, "aardvark")
	tree.Print(tree.root, "")
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
		tree.Insert(v, "")
	}
	called := 0
	callme := func(k, v []string) {
		t.Log(k)
		called += len(k)
	}
	tree.mkdir(tree.root, []string{"root"}, callme)
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
	
func TestDiskNodeJson(t *testing.T) {
	m := make(map[string]string)
	m["a"] = "b"
	s := disknode{"key", []string{"val1", "val2"}, m, "parent", "edgename", "hash"}
	b, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	c := new(disknode)
	err = json.Unmarshal(b, c)
	if err != nil {
		t.Fatal(err)
	}
	if c.Key != s.Key {
		t.Fatal("not the same decoded json", c, b)
	}
	if c.Value[0] != s.Value[0] {
		t.Fatal("not the same decoded json", c, b)
	}
	if c.Value[1] != s.Value[1] {
		t.Fatal("not the same decoded json", c, b)
	}
	if c.Children["a"] != s.Children["a"] {
		t.Fatal("not the same decoded json", c, b)
	}
	if c.Parent != s.Parent {
		t.Fatal("not the same decoded json", c, b)
	}
	if c.Edgename != s.Edgename {
		t.Fatal("not the same decoded json", c, b)
	}
	if c.Hash != s.Hash {
		t.Fatal("not the same decoded json", c, b)
	}

}	
	
