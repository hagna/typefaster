package main

import (
    "fmt"
    "strconv"
    "flag"
    "os"
    "bufio"
    "strings"
)

var verbose = flag.Bool("v", false, "verbose?")
var iphod = flag.String("iphod", "iphod.txt", "iphod file name")

type iphodrecord struct {
    nphones int
    phonemes string
}

var IPHOD map[string]iphodrecord

func readiphod() {
    IPHOD = make(map[string]iphodrecord)
    fh, err := os.Open(*iphod)
    if err != nil {
        fmt.Println(err)
    }
    defer fh.Close()
    scanner := bufio.NewScanner(fh)
    for scanner.Scan() {
        l := strings.Fields(scanner.Text())
        v := l[1]
        nphones, err := strconv.Atoi(l[5])
        if err != nil {
            fmt.Println(err)
            nphones = 0
        }
        phonemes := l[2]
        IPHOD[v] = iphodrecord{nphones, phonemes}
    }
}

func main() {
    flag.Parse()
    readiphod()
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
            if res, ok := IPHOD[w]; ok {
                if *verbose {
                    fmt.Println(res.phonemes)
                }
                total += res.nphones
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
