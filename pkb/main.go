package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/hagna/typefaster"
	"github.com/hagna/pt"
	"github.com/jmhodges/levigo"
	"log"
	"os"
)

var genfromiphod = flag.String("genfromiphod", "", "generate the tree from the iphod.txt")
//var treename = flag.String("treename", "root", "name of tree directory")
var leveldbname = flag.String("leveldbname", "", "name of level db")
var verbose = flag.Bool("v", false, "verbose?")
var print = flag.Bool("print", false, "print?")
var lookup = flag.Bool("lookup", false, "lookup")

func main() {
	flag.Parse()
	if *genfromiphod != "" {
		db := pt.NewTree(*leveldbname)
		defer db.Close()
		cb := func(word, phonemes string, nphones int) {	
			db.Insert(typefaster.Encode(phonemes), word)
		}
		if err := typefaster.MakeDB(*genfromiphod,  cb); err != nil {
			log.Println("problem loading iphod")
		}
		fmt.Println(db)
		return
	}

	
	if *print {
		db := pt.NewTree(*leveldbname)
		ro := levigo.NewReadOptions()
		cb := func(prefix string, val []string) {
			fmt.Fprintf(os.Stdout, "%s %s\n", typefaster.Decode(prefix), val)
		}
		defer ro.Close()
		db.Dfs(ro,  db.Root, "", cb)
	}
	if *lookup {
		tree := typefaster.NewDiskTree(*leveldbname)
		for _, w := range flag.Args() {
			we := typefaster.Encode(w)
			a, i := tree.Lookup(tree.Root(), we, 0)
			if a.Key != we {
				fmt.Printf("closest match \"%s\"\n", typefaster.Decode(a.Key[:i]))
			} 
			if len(a.Value) == 0 {
				fmt.Println("Here are all the spellings with a common prefix.")
				tree.Print(os.Stdout, a, "")
			} else {
				fmt.Println(a.Value)
			}
		}
		return
	}

	total := 0
	utotal := 0
	ucount := 0

	for _, fname := range flag.Args() {
		fh, err := os.Open(fname)
		if err != nil {
			fmt.Println(err)
		}
		defer fh.Close()
		scanner := bufio.NewScanner(fh)
		scanner.Split(bufio.ScanWords)
		for scanner.Scan() {
			w := scanner.Text()
			if res, ok := typefaster.IPHOD[w]; ok {
				if *verbose {
					fmt.Println(res.Phonemes)
				}
				total += res.Nphones
			} else {
				if *verbose {
					fmt.Printf("Unknown:%s\n", w)
				}
				utotal += len(w)
				ucount += 1
			}
		}
		if !*verbose {
			fmt.Printf("%d phonemes\n%d unknown words of total length %d\n", total, ucount, utotal)
		}
	}
}
