package typefaster

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
	n, m := t.lookup(root, k)
	fmt.Println("lookup found", n, m)
	if n == nil {
		newnode := &node{v, k, nil}
		root.children = append(root.children, newnode)
		return
	}
	if m == n.edgename {
		// simple case just add the rest
		fmt.Println("simple case add a node")
		nk := k[len(m):]
		if len(nk) > 0 {
			newnode := &node{v, nk, nil}
			n.children = append(n.children, newnode)
		} else {
			fmt.Println("node exists already")
		}
	} else {
		fmt.Println("split the node case")
		mp := matchprefix(m, n.edgename)
		nk := k[len(mp):]
		newnodeA := &node{v, nk, nil}
		rnk := n.edgename[len(mp):]
		newnodeB := &node{n.value, rnk, n.children}
		n.edgename = mp
		n.children = append(n.children, newnodeA)
		n.children = append(n.children, newnodeB)
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
	fmt.Println("lookup", *n, s)
	if s == "" {
		return n, ""
	}
	for _, c := range n.children {
		fmt.Println("checking child", *c)
		match = matchprefix(c.edgename, s)
		fmt.Printf("got match ('%s', '%s') -> '%s'\n", c.edgename, s, match)
		if match == "" {
			continue
		} else {
			var m string
			nres, m = t.lookup(c, s[len(m):])
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
