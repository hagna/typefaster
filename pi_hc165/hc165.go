package main

import (
	"github.com/hagna/gpio"
	"github.com/hagna/gpio/rpi"
	rpi2 "github.com/hagna/rpi"
	"log"
	"time"
)

var keymap []uint = []uint{3, 4, 2, 0, 5, 6, 7, 1}

type srpi struct {
	shld    gpio.Pin
	clkinh  gpio.Pin
	clk     gpio.Pin
	chassis gpio.Pin
	npins   int
}

/*
	cycle the clock return value of chassis after rising edge
*/
func (s *srpi) clock() (res bool) {
	s.clk.Set()
	time.Sleep(120 * time.Nanosecond)
	res = rpi2.GPIOGet(rpi.GPIO22) // TODO use gofix to make const
	s.clk.Clear()
	time.Sleep(150 * time.Nanosecond)
	return res
}

/*
	shift
*/
func (s *srpi) Shift() {
	time.Sleep(120 * time.Nanosecond)
	s.clkinh.Clear()
	time.Sleep(120 * time.Nanosecond)
}

/*
	Load
*/
func (s *srpi) Load() {
	s.shld.Clear()
	time.Sleep(120 * time.Nanosecond)
}

/*
   Close everything and clear the shift register(s)
*/
func (s *srpi) Close() {
	s.clkinh.Close()
	s.clk.Close()
	s.shld.Close()
	s.chassis.Close()
}

/*
Takes a number of keys and a callback cb which will be called with a
slice of pressed keys If cb returns false the goroutine in here will exit and
true will be sent on done.
*/
func NewSR(nkeys int, cb func(b []bool) bool) (done chan bool) {
	chassis, err := gpio.OpenPin(rpi.GPIO22, gpio.ModeInput)
	if err != nil {
		log.Fatal("Error opening pin", err)
	}
	ser, err := gpio.OpenPin(rpi.GPIO24, gpio.ModeOutput)
	if err != nil {
		log.Fatal("Error opening pin", err)
	}
	clk, err := gpio.OpenPin(rpi.GPIO25, gpio.ModeOutput)
	if err != nil {
		log.Fatal("Error opening pin", err)
	}
	pin, err := gpio.OpenPin(rpi.GPIO23, gpio.ModeOutput)
	if err != nil {
		log.Fatal("Error opening pin", err)
	}

	srpi := new(srpi)
	srpi.clkinh = ser
	srpi.clk = pin
	srpi.shld = clk
	srpi.chassis = chassis
	srpi.npins = nkeys

	done = make(chan bool)
	go func() {
		// this loop could run as fast as 5 Mhz according to the datasheet for hc165
		for {
			keys := make([]bool, srpi.npins)
			srpi.clk.Set()
			srpi.clkinh.Set()
			time.Sleep(250 * time.Nanosecond)
			srpi.clk.Clear()
			srpi.Load()
			time.Sleep(10 * time.Millisecond)
			srpi.shld.Set()
			for i := 0; i < srpi.npins; i++ {
				if srpi.clock() {
					k := keymap[i]
					//log.Println("key", i, "keymap", k)
					keys[k] = true
				}
				srpi.Shift()
			}
			if cb(keys) == false {
				break
			}
		}
		done <- true
		srpi.Close()
	}()
	return done
}
