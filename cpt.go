package typefaster

import "fmt"

type node struct {
	Value    string
	Edgename string
	Children []*node
}

type Tree struct {
	Root *node
}

func NewTree(rootval string) *Tree {
	return &Tree{&node{"root", "", nil}}
}

func (t *Tree) Print(n *node) {
	if n == nil {
	} else {
		for _, c := range n.Children {
			t.Print(c)
		}
		fmt.Println(n)

	}
}

func (t *Tree) Insert(root *node, k, v string) {
	n, m := t.Lookup(root, k)
	if n == nil {
		newnode := &node{v, k, nil}
		fmt.Println("add child", newnode)
		root.Children = append(root.Children, newnode)
		return
	}
	if m == n.Edgename {
		// simple case just add the rest
		nk := k[len(m):]
		if len(nk) > 0 {
			newnode := &node{v, nk, nil}
			n.Children = append(n.Children, newnode)
			fmt.Println("add child (simple)", newnode)
		} else {
			fmt.Println("node exists already")
		}
	} else {
		mp := matchprefix(m, n.Edgename)
		nk := k[len(mp):]
		newnodeA := &node{v, nk, nil}
		rnk := n.Edgename[len(mp):]
		newnodeB := &node{n.Value, rnk, n.Children}
		n.Edgename = mp
		n.Children = append(n.Children, newnodeA)
		n.Children = append(n.Children, newnodeB)
		fmt.Println("add child (split a)", newnodeA)
		fmt.Println("add child (split b)", newnodeB)
	}

}

// returns the matching prefix between the two
func matchprefix(a, b string) string {
	res := ""
	if len(a) < len(b) {
		for i, c := range a {
			if a[i] == b[i] {
				res += string(c)
			} else {
				break
			}
		}
	} else {
		for i, c := range b {
			if b[i] == a[i] {
				res += string(c)
			} else {
				break
			}
		}
	}
	return res
}

func (t *Tree) Lookup(n *node, s string) (nres *node, match string) {
	fmt.Printf("Lookup: NODE<%v> for '%s'\n", *n, s)
	if s == "" {
		return n, ""
	}
	for _, c := range n.Children {
		fmt.Printf("\tchild %s", c.Edgename)
		match = matchprefix(c.Edgename, s)
		if match == "" {
			fmt.Println(" does not match")
			continue
		} else {
			fmt.Println(" matches", len(match), "characters ->", match)
			var m string
			nres, m = t.Lookup(c, s[len(match):])
			match += m
			// for a partial match
			if nres == nil {
				return c, match
			}
			return nres, match
		}
	}
	return nil, ""
}
