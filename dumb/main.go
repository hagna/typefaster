package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var verbose = flag.Bool("v", false, "verbose?")
var iphod = flag.String("iphod", "iphod.txt", "iphod file name")
var mkdir = flag.String("mkdir", "", "make the iphod into a trie on disk")

func main() {
	flag.Parse()
	cb := func(key, value []string) {
		res := []string{}
		// skip first no encoded root
		for _, v := range key[1:] {
			res = append(res, v)
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
	if *iphod != "" {
		if *mkdir != "" {

			fh, err := os.Open(*iphod)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer fh.Close()
			scanner := bufio.NewScanner(fh)
			scanner.Scan() // chomp first line
			for scanner.Scan() {
				l := strings.Fields(scanner.Text())
				v := strings.ToLower(l[1])
				phonemes := l[2]
				p := []string{*mkdir}
				p = append(p, strings.Split(phonemes, ".")...)
				cb(p, []string{v})

			}
		}
	}
}
