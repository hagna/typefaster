package typefaster

import (
	"fmt"
	"crypto/md5"
	"io"
	"log"
	"strings"
	"os"
	"io/ioutil"
	"encoding/json"
)

/*
I want a data structure that is good for storing lots of strings of phonemes
with their corresponding english word(s), for example meat and meet are two
spellings of one phoneme string (M.IY.T in this case).  The structure ought to
not only store a lot of phoneme strings (more than you'd care to leave in
memory), but also make it fast to translate M.IY.T into meet or meat.

But why?  Because I've got plenty of pronunciation data, and I have a phonetic
system for typing that is easy to learn.  I just need to translate the sequence
of phonemes to words.

And how?

Using a compact prefix tree stored on the disk.  Each node is a filename in the
directory structure that will look like this:
root/
	e0/
		e018f5434:
			{
			"key": "M.IY.T"
			"value":["meat","meet"],
			"children":["d8774efef"]
			"parent":["fe88bdbc7"]
			}
		e09898342:
			...
	d8/
		d8342ffee:
			...
		d8eeff332:
			...
		d8774efef:
			{
			"key": "M.IY.T.IH.NG"
			"value":["meeting"],
			"children":["38528ef5b"]
			"parent":["d8774efef"]
			}


*/

/* a node in the compact prefix tree */
type node struct {
	Key	string
	Value    []string
	Children []*node
	Parent	*node
	Edgename string
}

/* a node in the compact prefix tree stored on disk */
type disknode struct {
	Key string	`json:"key"`
	Value []string  `json:"value"`
	Children map[string]string `json:"children"`
	Parent string	`json:"parent"`
	Edgename string `json:"edgename"`
	Hash string	`json:"hash"`
}

type Tree interface {
	Insert(key, value string)
	Lookup(parent *node, searchfor string) (c *node, p, m string) 
	String() string
}

type MemTree struct {
	root *node
}

type DiskTree struct {
	root *disknode
	path string
}

func (m MemTree) String() string {
	return fmt.Sprintf("%+v", m)
}

func NewMemTree(rootval string) *MemTree {
	i := new(MemTree)
	i.root = NewNode("root", "", nil)
	return i
}

func NewDiskTree(dirname string) *DiskTree {
	err := os.MkdirAll(dirname, 0777)
	if err != nil {
		log.Println(err)
	}
	res := new(DiskTree)
	res.root = new(disknode)
	res.root.Children = make(map[string]string)
	res.root.Key = dirname
	res.root.Hash = smash(dirname)
	res.path = dirname
	return res
}

// allusion to Schneier's description of a one-way hash function
func smash(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (t *DiskTree) Insert(k, v string) {
	log.Println("insert", k, v)
	root := t.root.toMem()
	n, part, m := t.Lookup(root, k)

	log.Printf("Lookup returns node '%+v' part '%v' match '%v'\n", n, part, m)
	if n == nil {
		root := t.root
		newnode := new(disknode) 
		newnode.Value = []string{v}
		newnode.Edgename = k
		newnode.Key = k
		newnode.Hash = smash(k)
		newnode.Parent = root.Hash
		root.Children[string(k[0])] = newnode.Hash
		log.Printf("add disk child %+v\n", newnode)
		t.write(root)
		t.write(newnode)
		return
	}
}

func (n *disknode) toMem() *node {
	if n == nil {
		return nil
	}
	res := NewNode("", n.Edgename, nil)
	res.Value = n.Value
	res.Key = n.Key
	return res
}

func (t *DiskTree) write(a *disknode) {
	dat, err := json.Marshal(a)
	if err != nil {
		log.Println(err)
	}
	fname := t.path + "/" + a.Hash
	log.Println("writing", string(dat), "to", fname)
	err = ioutil.WriteFile(fname, dat, 0666)
	if err != nil {
		log.Println(err)
	}
}

func (t *DiskTree) dnodeFromNode(n *node) *disknode {
	dn := t.dnodeFromHash(smash(n.Key))
	return dn
}

func (t *DiskTree) dnodeFromHash(s string) *disknode {
	dn := new(disknode)
	dat, err := ioutil.ReadFile(t.path + "/" + s)
	if err != nil {
		log.Println(err)
		return nil
	}
	err = json.Unmarshal(dat, dn)
	if err != nil {
		log.Println(err)
		return nil
	}
	return dn
}

func (t *DiskTree) Lookup(n *node, search string) (*node, string, string) {
	m := ""
	log.Println("Lookup node", n, "search for", search)
	// it is only nil if we're searching from root
	if n != nil {
		if search == n.Key {
			return n, "", n.Key
		}
		log.Println("matchprefix(",n.Edgename, search, ")")
		m = matchprefix(n.Edgename, search)
		if m == "" {
			log.Println("node", n, "has no prefix in common with", search)
			return nil, "", ""
		}
		// partial match
		if len(m) < len(n.Edgename) {
			log.Println("partial match", n, m)
			return n, "", m
		}
	}
 	if len(m) < len(search) {
		rest := search[len(m):]
		if n == nil {
			n = new(node)
			n.Key = t.root.Key
		}
		dn := t.dnodeFromNode(n)
		if dn != nil {
			// maybe later you'll have the courage to make this map uint8
			// keys instead of string
			if chash, ok := dn.Children[string(rest[0])]; ok {
				d := t.dnodeFromHash(chash)
				s := d.toMem()
				log.Printf("recurse on \"%s\" with child %+v", rest, s)
				nm, p, m2 := t.Lookup(s, rest)
				m += m2
				if p == "" {
					p = m2
				}
				if nm == nil {
					log.Println("didn't find it in the child", s)
					return s, p, m
				}
				return nm, p, m
			}
		}
	}	
	log.Println("returning nil because no case matched")
	return nil, "", ""
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


func NewNode(value, edgename string, children []*node) *node {
	v := []string{}
	if value != "" {
		v = append(v, value)
	}
	log.Println("NewNode: value is", v, len(v))
	log.Println("NewNode: edgename", edgename)
	res := &node{"", v, children, nil, edgename}
	return res
}

func (t MemTree) Insert(k, v string) {
	log.Println("insert", k, v)
	root := t.root
	n, part, m := t.Lookup(root, k)
	
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
	smallest := len(a)
	if smallest > len(b) {
		smallest = len(b)
	}
	for i:=0; i<smallest; i++ {
		if a[i] == b[i] {
			res += string(a[i])
		} else {
			break
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
