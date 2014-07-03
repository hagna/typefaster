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