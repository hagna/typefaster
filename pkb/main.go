package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/hagna/typefaster"
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
		var tree *typefaster.LDB
		var err error
		if tree, err = typefaster.MakeLDB(*genfromiphod, *leveldbname); err != nil {
			log.Println("problem loading iphod")
		}
		fmt.Println(tree)
		tree.Close()
		return
	}

	
	if *print {
		db := typefaster.NewLDB(*leveldbname)
		ro := levigo.NewReadOptions()
		ro.SetFillCache(false)
		it := db.NewIterator(ro)
		defer it.Close()
		it.SeekToFirst()
		for it = it; it.Valid(); it.Next() {
			fmt.Println(string(it.Key()), string(it.Value()))
		}
		if err := it.GetError(); err != nil {
			log.Fatal(err)
		}
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
