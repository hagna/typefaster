package main

import (
	"bufio"
	"flag"
	"fmt"
    "github.com/hagna/pkb/rawkb"
	"log"
	"os"
	"strconv"
	"strings"
)

type phon struct {
	cmu    string // from cmupd
	klat   string // klattese
	ipa    string // IPA
	mbrola string // for mbrola
	espeak string // for espeak
	ispeak string // for apple speak
	des    string // deseret alphabet
}

/*


*/
var verbose = flag.Bool("v", false, "verbose?")
var iphod = flag.String("iphod", "iphod.txt", "iphod file name")
var interactive = flag.Bool("i", false, "interactive mode")

type iphodrecord struct {
	nphones  int
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

func interact() {
	logfile, err := os.Create("log")
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logfile)
    if rawkb.SetupKeyboard() == 0 {
        log.Println("must be on the console for raw keyboard access")
        return
    }
    rawkb.RestoreKeyboard()
}


/*
type keys struct {
	termbox_event chan termbox.Event
}

func (k *keys) handle_event(ev *termbox.Event) bool {
	switch ev.Type {
	case termbox.EventKey:
		switch ev.Key {
		case termbox.KeyCtrlQ:
			return false
		default:
			log.Printf("%#v\n", ev)
		}
	}
	return true
}

func (k *keys) draw() {
	return
}

func (m *keys) mainloop() {
	m.termbox_event = make(chan termbox.Event, 20)
	go func() {
		for {
			m.termbox_event <- termbox.PollEvent()
		}
	}()
	for {
		select {
		case ev := <-m.termbox_event:
			ok := m.handle_event(&ev)
			if !ok {
				return
			}
			m.consume_more_events()
			m.draw()
			termbox.Flush()
		}
	}
}

func (m *keys) consume_more_events() {
loop:
	for {
		select {
		case ev := <-m.termbox_event:
			ok := m.handle_event(&ev)
			if !ok {
				break loop
			}
		default:
			break loop
		}
	}
}
*/

func main() {
	flag.Parse()
	readiphod()
	if *interactive {
		interact()
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
