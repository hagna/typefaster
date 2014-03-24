package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/hagna/typefaster/rawkb"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
    // consonants
    N = 1 << 4
    T = 1 << 3
    R = 1 << 1
    S = 1 << 2
    D = 1 << 5
    L = 1 << 1 | 1 << 4
    DH = 1 << 2 | 1 << 3
    Z = 1 << 3 | 1 << 4
    M = 1 << 1 | 1 << 2
    K = 1 << 2 | 1 << 3 | 1 << 4
    V = 1 << 1 | 1 << 3
    W = 1 << 1 | 1 << 2 | 1 << 3 | 1 << 4
    P = 1 << 1 | 1 << 2 | 1 << 3
    F = 1 << 1 | 1 << 5
    B = 1 << 4 | 1 << 5
    HH = 1 << 2 | 1 << 4
    NG = 1 << 2 | 1 << 3 | 1 << 4 | 1 << 5
    SH = 1 << 1 | 1 << 3 | 1 << 4
    G = 1 << 3 | 1 << 4 | 1 << 5
    Y = 1 << 1 | 1 << 2 | 1 << 3 | 1 << 4 | 1 << 5
    CH = 1 << 2 | 1 << 5
    JH = 1 << 1 | 1 << 4 | 1 << 5
    TH = 1 << 1 | 1 << 2 | 1 << 4
    ZH = 1 << 1 | 1 << 3 | 1 << 4 | 1 << 5
    
    // vowels
    AA = 1 
    IH2 = 1 | 1 << 4  // maybe get rid of this one? dennis?
    AO = 1 | 1 << 2
    IH = 1 | 1 << 1
    AE = 1 | 1 << 3
    EH = 1 | 1 << 2 | 1 << 3 | 1 << 4
    IY = 1 | 1 << 2 | 1 << 3 
    EY = 1 | 1 << 5
    AH = 1 | 1 << 3 | 1 << 4
    UW = 1 | 1 << 2 | 1 << 3 | 1 << 4 | 1 << 5
    AY = 1 | 1 << 4 | 1 << 5
    OW = 1 | 1 << 2 | 1 << 4
    UH = 1 | 1 << 3 | 1 << 4 | 1 << 5
    ER = 1 | 1 << 2 | 1 << 5
    AW = 1 | 1 << 2 | 1 << 3 | 1 << 5
    YU = 1 | 1 << 3 | 1 << 5  // really Y UW
    OY = 1 | 1 << 2 | 1 << 4 | 1 << 5
    
)

type phone struct {
	cmu    string // from cmupd
	klat   string // klattese
	ipa    string // IPA
	mbrola string // for mbrola
	espeak string // for espeak
	ispeak string // for apple speak
	des    string // deseret alphabet
}

var verbose = flag.Bool("v", false, "verbose?")
var iphod = flag.String("iphod", "iphod.txt", "iphod file name")
var interactive = flag.Bool("i", false, "interactive mode")

type iphodrecord struct {
	nphones  int
	phonemes string
}

var IPHOD map[string]iphodrecord

func readiphod() error {
	IPHOD = make(map[string]iphodrecord)
	fh, err := os.Open(*iphod)
	if err != nil {
		fmt.Println(err)
		return err
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
	return nil
}

type Pad struct {
	id          string
	boardstate  uint8
	keyspressed uint8
}


var Phones = map[uint8]phone{
    AA: phone{cmu: "AA"},
    AE: phone{cmu: "AE"},
    AH: phone{cmu: "AH"},
    AO: phone{cmu: "AO"},
    AW: phone{cmu: "AW"},
    AY: phone{cmu: "AY"},
    B: phone{cmu: "B"},
    CH: phone{cmu: "CH"},
    D: phone{cmu: "D"},
    DH: phone{cmu: "DH"},
    EH: phone{cmu: "EH"},
    ER: phone{cmu: "ER"},
    EY: phone{cmu: "EY"},
    F: phone{cmu: "F"},
    G: phone{cmu: "G"},
    HH: phone{cmu: "HH"},
    IH: phone{cmu: "IH"},
    IY: phone{cmu: "IY"},
    JH: phone{cmu: "JH"},
    K: phone{cmu: "K"},
    L: phone{cmu: "L"},
    M: phone{cmu: "M"},
    N: phone{cmu: "N"},
    NG: phone{cmu: "NG"},
    OW: phone{cmu: "OW"},
    OY: phone{cmu: "OY"},
    P: phone{cmu: "P"},
    R: phone{cmu: "R"},
    S: phone{cmu: "S"},
    SH: phone{cmu: "SH"},
    T: phone{cmu: "T"},
    TH: phone{cmu: "TH"},
    UH: phone{cmu: "UH"},
    UW: phone{cmu: "UW"},
    V: phone{cmu: "V"},
    W: phone{cmu: "W"},
    Y: phone{cmu: "Y"},
    Z: phone{cmu: "Z"},
    ZH: phone{cmu: "ZH"},
}

func decode(a uint8) string {
    if res, ok := Phones[a]; ok {
        return res.cmu
    }
    return fmt.Sprintf("%x", a)
}

func (p *Pad) sendphone() {
	log.Println("sendphone: ", decode(p.keyspressed))
}

func (p *Pad) keyup(keyno uint8) {
	b := p.boardstate
	b = b & ^(1 << keyno)
	p.boardstate = b
	if b == 0 {
		p.sendphone()
		p.keyspressed = 0
	}
}

func (p *Pad) keydown(keyno uint8) {
	log.Println(p.id, keyno)
	b := p.boardstate
	b = b | (1 << keyno)
	p.boardstate = b
	p.keyspressed = p.keyspressed | (1 << keyno)
}

type Mcs struct {
	keys     map[uint16]uint8
	watched  map[uint16]func()
	states   map[string]func(k uint16) string
	kbevents chan uint16
	pad      [2]Pad
}

func (m *Mcs) state_log(k uint16) string {
	if callable, ok := m.watched[k]; ok {
		callable()
	}
	return "loggy"
}

func iskeyup(k uint16) bool {
	if k > 0x80 {
		return true
	}
	return false
}

func (mcs *Mcs) state_learn(k uint16) string {
	b := k
	log.Println("learning the keys starting with", k)
	for iskeyup(b) {
		log.Println("learn: discarding keyup", b)
		b = <-mcs.kbevents
	}
	for j := 0; j < 2; j++ {
		var p Pad = mcs.pad[j]
		for i := 0; i < 6; i++ {
			var j int = i
			mcs.watched[b] = func() {
				p.keydown(uint8(j))
			}
			b = <-mcs.kbevents
			mcs.watched[b] = func() {
				p.keyup(uint8(j))
			}
			log.Println("learned key", i)
			b = <-mcs.kbevents
		}
	}
	return "loggy"
}

func NewMcs() *Mcs {
	m := new(Mcs)
	m.states = make(map[string]func(k uint16) string)
	m.keys = make(map[uint16]uint8, 12)
	m.kbevents = make(chan uint16, 2)
	m.watched = make(map[uint16]func())
	m.states["loggy"] = m.state_log
	m.states["init"] = m.state_learn
	m.pad = [2]Pad{{"left", 0, 0}, {"right", 0, 0}}
	log.Println(m.pad)
	return m
}

func interact() {
	/*	logfile, err := os.Create("log")
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(logfile)

	*/
	if rawkb.SetupKeyboard() == 0 {
		log.Println("must be on the console for raw keyboard access")
		return
	} else {
		log.Println("setting up rawkb")
		defer rawkb.RestoreKeyboard()
	}
	mcs := NewMcs()
	log.Println(mcs)
	log.Println(mcs.kbevents)
	go func() {
		mcs.kbevents <- 0x81
		for {
			b, ok := rawkb.ReadOnce()
			if ok != 255 {
				mcs.kbevents <- b

			}
			time.Sleep(1 * time.Microsecond)
		}
	}()

	m := mcs.states["init"]
inf:
	for {
		select {
		case b := <-mcs.kbevents:
			next_state := m(b)
			next_method, ok := mcs.states[next_state]
			if !ok {
				log.Println("no such state", next_state)
			} else {
				m = next_method
			}
			if b == 1 {
				break inf
			}
		}
	}

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
	/*
		if err := readiphod(); err != nil {
			log.Println("problem reading iphod")
			return
		}
	*/
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
