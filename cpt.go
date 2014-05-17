package main

import "fmt"

type node struct {
	value    string
	edgename string
	children []*node
}

type tree struct {
	root *node
}

func (t *tree) print(n *node) {
	if n == nil {
	} else {
		for _, c := range n.children {
			t.print(c)
		}
		fmt.Println(n)

	}
}

func (t *tree) insert(root *node, k, v string) {
	n, m := t.lookup(root, v)
	fmt.Println("lookup found", n, m)
	if n == nil {
		newnode := &node{v, k, nil}
		root.children = append(root.children, newnode)
	}

}

// returns the matching prefix between the two
func matchprefix(a, b string) string {
	res := ""
	if len(a) < len(b) {
		for i, c := range a {
			if a[i] == b[i] {
				res += string(c)
			}
		}
	} else {
		for i, c := range b {
			if b[i] == a[i] {
				res += string(c)
			}
		}
	}
	return res
}

func (t *tree) lookup(n *node, s string) (nres *node, match string) {
	if s == "" {
		return n, ""
	}
	for _, c := range n.children {
		match = matchprefix(c.edgename, s)
		if match == "" {
			continue
		} else {
			var m string
			nres, m = t.lookup(c, s[:len(m)])
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
