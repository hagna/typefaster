package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/hagna/typefaster"
	"log"
	"os"
)

var verbose = flag.Bool("v", false, "verbose?")
var iphod = flag.String("iphod", "iphod.txt", "iphod file name")
var mkdir = flag.String("mkdir", "", "make the iphod into a trie on disk")

func main() {
	flag.Parse()
	if *iphod != "" {
		if *mkdir != "" {
			if tree, err := typefaster.Maketree(*iphod); err != nil {
				log.Println("problem loading iphod")
			} else {
				tree.Mkdir(tree.Root, []string{*mkdir})
			}
			return
		}
		if err := typefaster.Readiphod(*iphod); err != nil {
			log.Println("problem reading iphod")
			return
		}
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
