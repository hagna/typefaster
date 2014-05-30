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
	Children []Node
}

func (n *node) String() string {
	return fmt.Sprintf("%+v", *n)
} 

func (n *node) GetNode() *nodelike {
	return &nodelike{n.Value, n.Edgename, n.Children}
}	

type nodelike struct {
	Value    []string
	Edgename string
	Children []Node
}

func (n *nodelike) GetNode() *nodelike {
	return n
}	

func (n *nodelike) String() string {
	return fmt.Sprintf("%+v", *n)
} 



/* this permits us to implement another version of node, but without any boilerplate getters and setters, and so far I think it beats boilerplate */
type Node interface {
	GetNode() *nodelike
	String() string
}

type Tree struct {
	Root *node
}

func NewTree(rootval string) *Tree {
	return &Tree{NewNode("root", "", nil)}
}

// depth first search 
func (t *Tree) Print(i Node, prefix string) {
	n := i.GetNode()
	if len(n.Children) == 0 {
		fmt.Println(prefix, n.Value)
	} else {
		for _, i := range n.Children {
			c := i.GetNode()
			t.Print(c, prefix+c.Edgename)
		}
		if len(n.Value) != 0 {
			fmt.Println(prefix, n.Value)
		}

	}
}

func (t *Tree) Mkdir(i Node, prefix []string) {
	n := i.GetNode()
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

func (t *Tree) mkdir(i Node, prefix []string, cb func(s, v []string)) {
	n := i.GetNode()
	if len(n.Children) == 0 {
		fmt.Println(prefix, n.Value)
		cb(prefix, n.Value)
	} else {
		for _, i := range n.Children {
			c := i.GetNode()
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


func NewNode(value, edgename string, children []Node) *node {
	v := []string{}
	if value != "" {
		v = append(v, value)
	}
	log.Println("NewNode: value is", v, len(v))
	log.Println("NewNode: edgename", edgename)
	res := &node{v, edgename, children}
	return res
}

func (t *Tree) Insert(root *node, k, v string) {
	log.Println("insert", k, v)
	i, part, m := t.Lookup(root, k)
	var n *node
	if i == nil {
		n = nil
	} else {
		n = i.GetNode()
	}
	log.Printf("Lookup returns node '%+v' part '%v' match '%v'\n", n, part, m)
	if n == nil {
		newnode := NewNode(v, k, nil)
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
				newnode := NewNode(v, nk, nil)
				n.Children = append(n.Children, newnode)
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
	newnodeA := NewNode(v, nk, nil)
	rnk := n.Edgename[len(mp):]
	newnodeB := NewNode("", rnk, n.Children) 
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
func (t *Tree) Lookup(i Node, s string) (nres Node, part, match string) {
	n := i.GetNode()
	log.Printf("Lookup: NODE<%+v> for '%s'\n", *n, s)
	if s == "" {
		return n, "", ""
	}
	for _, i := range n.Children {
		c := i.GetNode()
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
