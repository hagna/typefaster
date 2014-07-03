package typefaster

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
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
type Node struct {
	Key      string
	Value    []string
	Children []*Node
	Parent   *Node
	Edgename string
}

/* a node in the compact prefix tree stored on disk */
type disknode struct {
	Key      string            `json:"key"`
	Value    []string          `json:"value"`
	Children map[string]string `json:"children"`
	Parent   string            `json:"parent"`
	Edgename string            `json:"edgename"`
	Hash     string            `json:"hash"`
}

type Tree interface {
	Insert(key, value string)
	Lookup(parent *Node, searchfor string) (c *Node, p, m string)
	Root() *Node
	String() string
}

func (t *DiskTree) Root() *Node {
	return t.root.toMem()
}

type MemTree struct {
	root *Node
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
		debug(err)
	}
	res := new(DiskTree)
	res.root = new(disknode)
	res.root.Children = make(map[string]string)
	res.root.Key = dirname + "asdfasdfjkl;ajsdl;fkjaskl;djasdf"
	res.root.Hash = smash(res.root.Key)
	res.path = dirname
	if _, err := os.Stat(dirname + "/" + res.root.Hash); os.IsNotExist(err) {
		debug("no such file or directory: %s CREATING", dirname)

		res.write(res.root)
	} else {
		debug("dir exists already")
		res.root = res.dnodeFromHash(res.root.Hash)
	}
	return res
}

// allusion to Schneier's description of a one-way hash function
func smash(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

/*
creates new child and writes it to disk along with the parent
*/
func (t *DiskTree) addChild(root *disknode, edgename, key string, value []string) {
	newnode := new(disknode)
	newnode.Value = value
	newnode.Edgename = edgename
	newnode.Key = key
	newnode.Hash = smash(key)
	newnode.Parent = root.Hash
	if root.Children == nil {
		root.Children = make(map[string]string)
	}
	root.Children[string(edgename[0])] = newnode.Hash
	debugf("add disk child %+v\n", newnode)
	t.write(root)
	t.write(newnode)
}

func debug(i ...interface{}) {
	return
	var msg string
	if _, fname, lineno, ok := runtime.Caller(1); !ok {
		debug("couldn't get line number")
	} else {
		j := strings.LastIndex(fname, "/")
		fname = fname[j+1:]
		msg = fmt.Sprintf("%s:%d ", fname, lineno)
	}

	fmt.Printf(msg)
	fmt.Println(i...)
}

func debugf(format string, i ...interface{}) {
	return
	var msg string
	if _, fname, lineno, ok := runtime.Caller(1); !ok {
		debug("couldn't get line number")
	} else {
		j := strings.LastIndex(fname, "/")
		fname = fname[j+1:]
		msg = fmt.Sprintf("%s:%d ", fname, lineno)
	}
	fmt.Printf(msg+format, i...)
}

func (t *DiskTree) Insert(k, v string) {
	debug("insert", k, v)
	root := t.root.toMem()
	n, i := t.Lookup(root, k, 0)
	commonprefix := k[:i]
	debug("Insert", k, "and commonprefix is", commonprefix)

	debugf("Lookup returns node '%+v' mathced chars = '%v' match '%v'\n", n, i, k[:i])

	debug("is it the root?")
	if n == root {
		debug("addChild")
		latestroot := t.dnodeFromHash(t.root.Hash)
		t.addChild(latestroot, k, k, []string{v})
		debug("yes")
		return
	}
	debug("no")

	debug("is it a complete match?")
	if k == n.Key {
		dn := t.dnodeFromNode(n)
		dn.Value = append(dn.Value, v)
		t.write(dn)

		debug("node", n, "already found append value here TODO")
		debug("yes")
		return
	}
	debug("no")

	debug("does commonprefix consume the whole tree so far?")
	// the best match matches the whole key (including n.Edgename)
	if commonprefix == n.Key {
		// but if it is longer than the key it's a simple add
		if len(k) > len(n.Key) {
			e := k[len(commonprefix):]
			dn := t.dnodeFromNode(n)
			t.addChild(dn, e, k, []string{v})
			debug("yes")
			return
		}
	}
	debug("no")

	// otherwise it's a split because it matches part of n.Edgename

	debug("split them then")

	/* say we have the string "key" and we add "ketones"
	   then the left node will be "y" the right node will be "tones"
	   and the middle will be "ke"
	*/

	mid := t.dnodeFromNode(n)
	children := make(map[string]string)

	// whatever is left in n.Key after taking out the length of common prefix
	lname := n.Key[len(commonprefix):]
	rname := k[len(commonprefix):]

	// index of edgename
	ie := strings.LastIndex(n.Key, n.Edgename)
	midname := n.Key[ie:len(commonprefix)]

	debug("into", lname, rname)

	// left node (preserve the old string)
	leftnode := new(disknode)
	leftnode.Value = n.Value
	leftnode.Edgename = lname
	leftnode.Key = mid.Key
	leftnode.Hash = mid.Hash
	leftnode.Children = mid.Children
	children[string(leftnode.Edgename[0])] = leftnode.Hash

	// update the middle node
	mid.Edgename = midname
	mid.Value = []string{}
	mid.Key = commonprefix
	mid.Hash = smash(commonprefix)
	leftnode.Parent = mid.Hash

	// if you have 'cats' and try to add 'cat'
	// you'll have this empty right node case
	if rname != "" {
		// right node (add the new string)
		rightnode := new(disknode)
		rightnode.Value = append(rightnode.Value, v)
		rightnode.Edgename = rname
		rightnode.Key = k
		rightnode.Hash = smash(k)
		children[string(rightnode.Edgename[0])] = rightnode.Hash
		rightnode.Parent = mid.Hash
		t.write(rightnode)
	}

	mid.Children = children

	// also update mid's parent hash
	midparent := t.dnodeFromHash(mid.Parent)
	midparent.Children[string(midname[0])] = mid.Hash

	t.write(midparent)
	t.write(mid)
	t.write(leftnode)

}

func (n *disknode) toMem() *Node {
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
		debug(err)
	}
	fname := t.path + "/" + a.Hash
	debug("writing", string(dat), "to", fname)
	err = ioutil.WriteFile(fname, dat, 0666)
	if err != nil {
		debug(err)
	}
}

func (t *DiskTree) dnodeFromNode(n *Node) *disknode {
	dn := t.dnodeFromHash(smash(n.Key))
	return dn
}

func (t *DiskTree) dnodeFromHash(s string) *disknode {
	dn := new(disknode)
	dat, err := ioutil.ReadFile(t.path + "/" + s)
	if err != nil {
		debug(err)
		return nil
	}
	err = json.Unmarshal(dat, dn)
	if err != nil {
		debug(err)
		return nil
	}
	return dn
}

/* fetch the already existing child of node n from the disk that starts with c */
func (t *DiskTree) fetchChild(n *Node, c string) *Node {
	dn := t.dnodeFromNode(n)
	v, ok := dn.Children[c]
	if ok {
		dv := t.dnodeFromHash(v)
		return dv.toMem()
	}
	return nil
}

/*
	Lookup takes the node to start from, the string to search for, and a
	count of how many chars are matched already.

	It returns the node that matches most closely and the number of
	characters (starting from 0) that match.

*/
func (t *DiskTree) Lookup(n *Node, search string, i int) (*Node, int) {

	if n == nil {
		return nil, i
	}
	dn := t.dnodeFromNode(n)
	debugf("Lookup(%+v, \"%s\", %d)\n", dn, search, i)
	match := matchprefix(n.Edgename, search[i:])
	i += len(match)
	if i < len(search) && len(n.Edgename) == len(match) {
		child := t.fetchChild(n, string(search[i]))
		c, i := t.Lookup(child, search, i)
		if c != nil {
			return c, i
		}
	}
	return n, i
}

// depth first print
// too much converting between node and disknode
func (t *DiskTree) Print(w io.Writer, n *Node, prefix string) {
	dn := t.dnodeFromNode(n)
	if len(dn.Children) == 0 {
		fmt.Fprintf(w, "%s %s\n", Decode(prefix), n.Value)
	} else {
		for _, c := range dn.Children {
			cnode := t.dnodeFromHash(c)
			t.Print(w, cnode.toMem(), prefix+cnode.Edgename)
		}
		if len(n.Value) != 0 {
			fmt.Fprintf(w, "%s %s\n", Decode(prefix), n.Value)
		}
	}
}

// depth first search
func (t *MemTree) Print(n *Node, prefix string) {
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

func (t *MemTree) Mkdir(n *Node, prefix []string) {
	cb := func(key, value []string) {
		res := []string{}
		// skip first no encoded root
		for _, v := range key[1:] {
			res = append(res, Decode(v))
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

func (t *MemTree) mkdir(n *Node, prefix []string, cb func(s, v []string)) {
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
			debugf("wellFormed: (no) '%s' != '%s'\n", partmatch, k)
			return false
		}
	}
	return res
}

func NewNode(value, edgename string, children []*Node) *Node {
	v := []string{}
	if value != "" {
		v = append(v, value)
	}
	res := &Node{"", v, children, nil, edgename}
	return res
}

func (t MemTree) Insert(k, v string) {
	debug("insert", k, v)
	root := t.root
	n, part, m := t.Lookup(root, k)

	debugf("Lookup returns node '%+v' part '%v' match '%v'\n", n, part, m)
	if n == nil {
		newnode := NewNode(v, k, nil)
		debug("add child", newnode)
		root.Children = append(root.Children, newnode)
		return
	}
	if n.Edgename == part || n.Edgename == m {
		// simple case just add the rest
		if !wellFormed(part, m, n.Edgename, k) {
			debug("would not be well formed")
		} else {
			nk := k[len(m):]
			if len(nk) > 0 {
				newnode := NewNode(v, nk, nil)
				n.Children = append(n.Children, newnode)
				debug("add child (simple)", newnode)
			} else {
				n.Value = append(n.Value, v)
				debug("node exists already")
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
	debug("add child (split a)", newnodeA)
	debug("add child (split b)", newnodeB)

}

// returns the matching prefix between the two
func matchprefix(a, b string) string {
	res := ""
	smallest := len(a)
	if smallest > len(b) {
		smallest = len(b)
	}
	for i := 0; i < smallest; i++ {
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
func (t MemTree) Lookup(n *Node, s string) (nres *Node, part, match string) {
	debugf("Lookup: NODE<%+v> for '%s'\n", *n, s)
	if s == "" {
		return n, "", ""
	}
	for _, c := range n.Children {
		debugf("\tchild %s", c.Edgename)
		match = matchprefix(c.Edgename, s)
		if match == "" {
			debug(" does not match")
			continue
		} else {
			debug(" matches", len(match), "characters ->", match)
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
