package typefaster

import (
	"fmt"
	"log"
	"strings"
)

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

func (t *Tree) Print(n *node, prefix string) {
	if len(n.Children) == 0 {
		fmt.Println(prefix)
	} else {
		for _, c := range n.Children {
			t.Print(c, prefix+c.Edgename)
		}
		if n.Value != "" {
			fmt.Println(prefix)
		}

	}
}

// tells us if we have a duplicate edge
// as in aaardvark
func isDup(part, match, edgename, k string) bool {
	res := false
	if len(match) <= len(edgename) {
		partmatch := match + part + edgename[len(match):]
		if strings.HasPrefix(k, partmatch) {
		} else {
			log.Printf("isDup: (yes) '%s' != '%s'\n", partmatch, k)
			return true
		}
	}
	return res
}

func (t *Tree) Insert(root *node, k, v string) {
	log.Println("insert", k, v)
	n, part, m := t.Lookup(root, k)
	log.Printf("Lookup returns node '%+v' part '%v' match '%v'\n", n, part, m)
	if n == nil {
		newnode := &node{v, k, nil}
		log.Println("add child", newnode)
		root.Children = append(root.Children, newnode)
		return
	}
	if n.Edgename == part || n.Edgename == m {
		// simple case just add the rest
		if isDup(part, m, n.Edgename, k) {
			log.Println("would be a dup")
		} else {
			nk := k[len(m):]
			if len(nk) > 0 {
				newnode := &node{v, nk, nil}
				n.Children = append(n.Children, newnode)
				log.Println("add child (simple)", newnode)
			} else {
				log.Println("node exists already")
			}
			return
		}
	}

	if part == "" {
		part = m
	}

	mp := part
	nk := k[len(m):]
	newnodeA := &node{v, nk, nil}
	rnk := n.Edgename[len(mp):]
	newnodeB := &node{n.Value, rnk, n.Children}
	n.Edgename = mp
	n.Value = ""
	n.Children = nil
	n.Children = append(n.Children, newnodeA)
	n.Children = append(n.Children, newnodeB)
	log.Println("add child (split a)", newnodeA)
	log.Println("add child (split b)", newnodeB)

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

/*
Lookup return the partial match of the current node and the match in the tree so far
*/
func (t *Tree) Lookup(n *node, s string) (nres *node, part, match string) {
	log.Printf("Lookup: NODE<%+v> for '%s'\n", *n, s)
	if s == "" {
		return n, "", ""
	}
	for _, c := range n.Children {
		log.Printf("\tchild %s", c.Edgename)
		match = matchprefix(c.Edgename, s)
		if match == "" {
			log.Println(" does not match")
			continue
		} else {
			log.Println(" matches", len(match), "characters ->", match)
			if len(match) < len(c.Edgename) {
				return c, "", match
			}
			var m string
			nres, part, m = t.Lookup(c, s[len(match):])
			match += m
			if part == "" {
				part = m
			}
			// for a partial match
			if nres == nil {
				return c, part, match
			}
			return nres, part, match
		}
	}
	return nil, "", ""
}
