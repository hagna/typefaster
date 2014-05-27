package main

import (
	"flag"
	"fmt"
	"github.com/hagna/typefaster"
	"log"
	"os"
	"os/signal"
)


var verbose = flag.Bool("v", false, "verbose?")
var iphod = flag.String("iphod", "iphod.txt", "iphod file name")
var pimode = flag.Bool("pi", false, "use shift register connected to raspberry pi")


func decode(a uint8) string {
	if res, ok := typefaster.Phones[a]; ok {
		return res.Cmu
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
want curses here they are -->
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
			fmt.Println("-=-=-=-=-=-=-=-=--=-")
			tree.Print(tree.Root, "")
			fmt.Println("-=-=-=-=-=-=-=-=--=-")
			a, b, c := tree.Lookup(tree.Root, "abstention")
			fmt.Println("is abstention found?", a, b, c)
		}
	}
	pi_shiftreg_interact()
	return

}
