package main

import (
	"flag"
	"fmt"
	"github.com/hagna/typefaster"
	"github.com/hagna/pt"
	"io"
	"log"
	"strings"
	"os"
	"os/signal"
	"github.com/huin/goserial"
)


var verbose = flag.Bool("v", false, "verbose?")
var treename = flag.String("treename", "root", "name of tree directory")


func decode(a uint8) string {
	if res, ok := typefaster.Phones[a]; ok {
		return res.Cmu
	}
	return fmt.Sprintf("%x", a)
}

type Mcs struct {
	buf uint8
	Tree *pt.Tree
	Cnode *pt.Node
	cword []string
	iword int
	serial io.ReadWriteCloser
}

func NewMcs() *Mcs {
	m := new(Mcs)
	m.Tree = pt.NewTree(*treename)
	c := new(goserial.Config)
	c.Name = "/dev/ttyAMA0"
	c.Baud = 9600
        s, err := goserial.OpenPort(c)
        if err != nil {
                log.Fatal(err)
        }
        
        n, err := s.Write([]byte{97})
        if err != nil {
                log.Fatal(err)
        } else {
		fmt.Println("serial write got back", n)
	}
	m.serial = s
 
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

/* Turns keystate ([]bool) into a useful value for decode. 8 keys only :) */
func decodestate(keys []bool) uint8 {
	var res uint8 = 0
	for i, v := range keys {
		if v {
			res |= 1 << uint8(i)
		}
	}
	return res
}

func isLast(m uint8) bool {
	return (m & 0xc0) != 0
}

/* Decode strokes this ought to run at some high rate in hz */
func (m *Mcs) keystates(keys []bool) bool {
	if keysup(keys) {
		// got a stroke
		if m.buf == 0xff {
			return false //quit
		}
		if m.buf != 0 {
			fmt.Printf("%x\n", m.buf)
			
			phon := decode(0x3f & m.buf)
			ebuf := typefaster.Encode(phon)
			m.cword = append(m.cword, ebuf)
			fmt.Println(ebuf)
			if isLast(m.buf) {
				we := strings.Join(m.cword, "")
				we = we[:len(we)-1] 
				a, i := m.Tree.Lookup(m.Tree.Root, we, 0)
				if a.Name != we {
					fmt.Printf("closest match to \"%s\" was \"%s\"\n", typefaster.Decode(we), typefaster.Decode(a.Name[:i]))
				} 
				if len(a.Value) == 0 {
					fmt.Println("Here are all the spellings with a common prefix.")
					fmt.Println(a)
					//m.Tree.Print(os.Stdout, a, "")
				} else {
					m.serial.Write(append([]byte(a.Value[0]), 0x20))
					fmt.Println(a.Value)
				}

				m.cword = []string{}
			}

			m.buf = 0
			fmt.Println(phon)

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
/*
	printem := func (keys []bool) bool {
		fmt.Println(keys)
		return true
	}
	fmt.Println(mcs)
	done := NewSR(nkeys, printem)
*/
 	done := NewSR(nkeys, mcs.keystates)
	<-done

}

/*
want curses? here they are -->
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
	pi_shiftreg_interact()
	return

}
