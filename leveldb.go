package typefaster

import (
	"log"
	"github.com/jmhodges/levigo"
)

type LDB struct {
	*levigo.DB
}

func NewLDB(dbname string) *LDB {
	n := new(LDB)
	opts := levigo.NewOptions()
	opts.SetCache(levigo.NewLRUCache(3<<30))
	opts.SetCreateIfMissing(true)
	db, err := levigo.Open(dbname, opts)
	if err != nil {
		log.Fatal(err)
	}
	n.DB = db
	return n
}

func (l *LDB) Lookup(search string, i int) (*Node, int) {
	ro := levigo.NewReadOptions()
	n, err := l.Get(ro, []byte(search))
	ro.Close()
	if err != nil {
		log.Println(err)
		return nil, 0
	}
	v := string(n)
	return NewNode(v, v, nil), len(v)
}