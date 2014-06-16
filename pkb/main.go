package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/hagna/typefaster"
	"log"
	"os"
)

var genfromiphod = flag.String("genfromiphod", "", "generate the tree from the iphod.txt")
var treename = flag.String("treename", "root", "name of tree directory")
var verbose = flag.Bool("v", false, "verbose?")
var print = flag.Bool("print", false, "print?")
var lookup = flag.Bool("lookup", false, "lookup")

func main() {
	flag.Parse()
	if *genfromiphod != "" {
		var tree *typefaster.DiskTree
		var err error
		if tree, err = typefaster.Maketree(*genfromiphod, *treename); err != nil {
				log.Println("problem loading iphod")
		}
		fmt.Println(tree)
	
		return 
	}
	if *print {
		t := typefaster.NewDiskTree(*treename)
		t.Print(os.Stdout, t.Root(), "")
		return
	}
	if *lookup {
		tree := typefaster.NewDiskTree(*treename)
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
