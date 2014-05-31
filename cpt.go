package typefaster

import (
	"fmt"
	"log"
	"strings"
	"os"
)

type node struct {
	Value    []string
	Edgename string
	Children []*node
}

type Tree interface {
	NewNode(value, edgename string, children []*node) *node 
	Root() *node 
	Splitnode(n, l, r *node)  // new, left, right are the args
	Addnode(parent, newchild *node)
	// c is the closest node
	// p is the partial match on that node
	// m is the match in the tree so far
	Lookup(parent *node, searchfor string) (c *node, p, m string) 
	String() string
}

type MemTree struct {
	root *node
}

func (m MemTree) Root() *node {
	return m.root
}

func (m MemTree) Splitnode(n, l, r *node) {
}

func (m MemTree) Addnode(root, n *node) {
}

func (m MemTree) String() string {
	return fmt.Sprintf("%+v", m)
}

func NewMemTree(rootval string) *MemTree {
	i := new(MemTree)
	i.root = i.NewNode("root", "", nil)
	return i
}

// depth first search 
func (t *MemTree) Print(n *node, prefix string) {
	if len(n.Children) == 0 {
		fmt.Println(prefix, n.Value)
	} else {
		for _, c := range n.Children {
			t.Print(c, prefix+c.Edgename)
		}
		if len(n.Value) != 0 {
			fmt.Println(prefix, n.Value)
		}

	}
}

func (t *MemTree) Mkdir(n *node, prefix []string) {
	cb :=  func(key, value []string) {
		res := []string{}
		// skip first no encoded root
		for _, v := range key[1:] {
			res = append(res, decode(v))
		}
		dir := key[0] + "/" + strings.Join(res, "/")
		fmt.Println("DIR", dir)
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			fmt.Println(err)
		}
		vname := dir + "/value"
		if _, err := os.Stat(vname); os.IsNotExist(err) {
			f, err := os.Create(vname)
			if err != nil {
				fmt.Println(err)
			}
			for _, v := range value {
				fmt.Fprintf(f, "%s\n", v)
			}
			f.Close()
		} else {
			f, err := os.OpenFile(vname, os.O_RDWR, 0666)
			if err != nil {
				fmt.Println(err)
			}
			for _, v := range value {
				fmt.Fprintf(f, "%s\n", v)
			}
			f.Close()
		}
	}
	t.mkdir(n, prefix, cb)
}

func (t *MemTree) mkdir(n *node, prefix []string, cb func(s, v []string)) {
	if len(n.Children) == 0 {
		fmt.Println(prefix, n.Value)
		cb(prefix, n.Value)
	} else {
		for _, c := range n.Children {
			t.mkdir(c, append(prefix, c.Edgename), cb)
		}
		if len(n.Value) != 0 {
			cb(prefix, n.Value)
			fmt.Println(prefix, n.Value)
		}

	}
}

// tells us if we have a duplicate edge
// as in aaardvark
func wellFormed(part, match, edgename, k string) bool {
	res := true
	if len(match) <= len(edgename) {
		partmatch := match + part + edgename[len(match):]
		if strings.HasPrefix(k, partmatch) {
		} else {
			log.Printf("wellFormed: (no) '%s' != '%s'\n", partmatch, k)
			return false
		}
	}
	return res
}


func (t MemTree) NewNode(value, edgename string, children []*node) *node {
	v := []string{}
	if value != "" {
		v = append(v, value)
	}
	log.Println("NewNode: value is", v, len(v))
	log.Println("NewNode: edgename", edgename)
	res := &node{v, edgename, children}
	return res
}

func Insert(t Tree, k, v string) {
	log.Println("insert", k, v)
	root := t.Root()
	n, part, m := t.Lookup(root, k)
	
	log.Printf("Lookup returns node '%+v' part '%v' match '%v'\n", n, part, m)
	if n == nil {
		newnode := t.NewNode(v, k, nil)
		log.Println("add child", newnode)
		root.Children = append(root.Children, newnode)
		return
	}
	if n.Edgename == part || n.Edgename == m {
		// simple case just add the rest
		if !wellFormed(part, m, n.Edgename, k) {
			log.Println("would not be well formed")
		} else {
			nk := k[len(m):]
			if len(nk) > 0 {
				newnode := t.NewNode(v, nk, nil)
				t.Addnode(n, newnode)
				/*
				n.Children = append(n.Children, newnode)
*/
				log.Println("add child (simple)", newnode)
			} else {
				n.Value = append(n.Value, v)
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
	newnodeA := t.NewNode(v, nk, nil)
	rnk := n.Edgename[len(mp):]
	newnodeB := t.NewNode("", rnk, n.Children) 
	newnodeB.Value = n.Value
	
	n.Edgename = mp
	n.Value = []string{}
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
func (t MemTree) Lookup(n *node, s string) (nres *node, part, match string) {
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
