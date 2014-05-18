package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"github.com/hagna/typefaster"
)

const (
	// consonants
	N  = 1 << 4
	T  = 1 << 3
	R  = 1 << 1
	S  = 1 << 2
	D  = 1 << 5
	L  = 1<<1 | 1<<4
	DH = 1<<2 | 1<<3
	Z  = 1<<3 | 1<<4
	M  = 1<<1 | 1<<2
	K  = 1<<2 | 1<<3 | 1<<4
	V  = 1<<1 | 1<<3
	W  = 1<<1 | 1<<2 | 1<<3 | 1<<4
	P  = 1<<1 | 1<<2 | 1<<3
	F  = 1<<1 | 1<<5
	B  = 1<<4 | 1<<5
	HH = 1<<2 | 1<<4
	NG = 1<<2 | 1<<3 | 1<<4 | 1<<5
	SH = 1<<1 | 1<<3 | 1<<4
	G  = 1<<3 | 1<<4 | 1<<5
	Y  = 1<<1 | 1<<2 | 1<<3 | 1<<4 | 1<<5
	CH = 1<<2 | 1<<5
	JH = 1<<1 | 1<<4 | 1<<5
	TH = 1<<1 | 1<<2 | 1<<4
	ZH = 1<<1 | 1<<3 | 1<<4 | 1<<5

	// vowels
	AA  = 1
	IH2 = 1 | 1<<4 // maybe get rid of this one? dennis?
	AO  = 1 | 1<<2
	IH  = 1 | 1<<1
	AE  = 1 | 1<<3
	EH  = 1 | 1<<2 | 1<<3 | 1<<4
	IY  = 1 | 1<<2 | 1<<3
	EY  = 1 | 1<<5
	AH  = 1 | 1<<3 | 1<<4
	UW  = 1 | 1<<2 | 1<<3 | 1<<4 | 1<<5
	AY  = 1 | 1<<4 | 1<<5
	OW  = 1 | 1<<2 | 1<<4
	UH  = 1 | 1<<3 | 1<<4 | 1<<5
	ER  = 1 | 1<<2 | 1<<5
	AW  = 1 | 1<<2 | 1<<3 | 1<<5
	YU  = 1 | 1<<3 | 1<<5 // really Y UW
	OY  = 1 | 1<<2 | 1<<4 | 1<<5
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
var pimode = flag.Bool("pi", false, "use shift register connected to raspberry pi")

var Phones = map[uint8]phone{
	AA: phone{cmu: "AA"},
	AE: phone{cmu: "AE"},
	AH: phone{cmu: "AH"},
	AO: phone{cmu: "AO"},
	AW: phone{cmu: "AW"},
	AY: phone{cmu: "AY"},
	B:  phone{cmu: "B"},
	CH: phone{cmu: "CH"},
	D:  phone{cmu: "D"},
	DH: phone{cmu: "DH"},
	EH: phone{cmu: "EH"},
	ER: phone{cmu: "ER"},
	EY: phone{cmu: "EY"},
	F:  phone{cmu: "F"},
	G:  phone{cmu: "G"},
	HH: phone{cmu: "HH"},
	IH: phone{cmu: "IH"},
	IY: phone{cmu: "IY"},
	JH: phone{cmu: "JH"},
	K:  phone{cmu: "K"},
	L:  phone{cmu: "L"},
	M:  phone{cmu: "M"},
	N:  phone{cmu: "N"},
	NG: phone{cmu: "NG"},
	OW: phone{cmu: "OW"},
	OY: phone{cmu: "OY"},
	P:  phone{cmu: "P"},
	R:  phone{cmu: "R"},
	S:  phone{cmu: "S"},
	SH: phone{cmu: "SH"},
	T:  phone{cmu: "T"},
	TH: phone{cmu: "TH"},
	UH: phone{cmu: "UH"},
	UW: phone{cmu: "UW"},
	V:  phone{cmu: "V"},
	W:  phone{cmu: "W"},
	Y:  phone{cmu: "Y"},
	Z:  phone{cmu: "Z"},
	ZH: phone{cmu: "ZH"},
}

func decode(a uint8) string {
	if res, ok := Phones[a]; ok {
		return res.cmu
	}
	return fmt.Sprintf("%x", a)
}

type Mcs struct {
	buf uint8
}

func NewMcs() *Mcs {
	m := new(Mcs)
	return m
}

func keysup(keys []bool) bool {
	res := true
	for _, v := range keys {
		if v {
			return false
		}
	}
	return res
}

/* Turns keystate ([]bool) into a useful value for decode */
func decodestate(keys []bool) uint8 {
	var res uint8 = 0
	for i, v := range keys {
		if v {
			res |= 1 << uint8(i)
		}
	}
	return res
}

/* Decode strokes this ought to run at some high rate in hz */
func (m *Mcs) keystates(keys []bool) bool {
	if keysup(keys) {
		if m.buf == 0xff {
			return false //quit
		}
		if m.buf != 0 {
			fmt.Print(decode(m.buf))
			m.buf = 0
		}
	} else {
		m.buf |= decodestate(keys)
		return true
	}
	return true
}

func pi_shiftreg_interact() {
	mcs := NewMcs()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			os.Exit(0)
		}
	}()
	nkeys := 8
	done := NewSR(nkeys, mcs.keystates)
	<-done

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
	if *iphod != "" {
		if tree, err := typefaster.Maketree(*iphod); err != nil {
			log.Println("problem reading iphod")
			return
		} else {
		tree.Print(tree.Root)
		}
	}
	pi_shiftreg_interact()
	return

}
