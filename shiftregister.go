package main

import (
	"github.com/davecheney/gpio"
	"github.com/davecheney/gpio/rpi"
	"log"
	"os"
	"os/signal"
	"time"
	"sync"
)

type Srpi struct {
	shld gpio.Pin
	clkinh  gpio.Pin
	clk   gpio.Pin
	chassis gpio.Pin
	curpin  int
	npins   int
	sync.RWMutex
}

/*
	cycle the clock
*/
func (s *Srpi) clock() {
	s.clk.Set()
	time.Sleep(120 * time.Nanosecond)
	s.clk.Clear()
	time.Sleep(150 * time.Nanosecond)
}

/*
	shift 
*/
func (s *Srpi) Shift() {
	time.Sleep(100 * time.Nanosecond)
	s.shld.Set()
	time.Sleep(120 * time.Nanosecond)
	s.clkinh.Clear()
	time.Sleep(120 * time.Nanosecond)
	s.clock()
}


/*
	Load
*/
func (s *Srpi) Load() {
	time.Sleep(120 * time.Nanosecond)
	s.clkinh.Set()
	time.Sleep(120 * time.Nanosecond)
	s.shld.Clear()
	time.Sleep(120 * time.Nanosecond)
	s.clock()
}


/*
   Close everything and clear the shift register(s)
*/
func (s *Srpi) Close() {
	s.clkinh.Close()
	s.clk.Close()
	s.shld.Close()
	s.chassis.Close()
}

/*
   The callback called when chassis gets a signal,
   used for detecting which key.
*/
func (s *Srpi) chassis_cb() {
	s.RLock()
	log.Println("button pressed", s.curpin)
	s.RUnlock()
}

func NewSrpi() *Srpi {
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

	srpi := new(Srpi)
	srpi.clkinh = ser
	srpi.clk = pin
	srpi.shld = clk
	srpi.chassis = chassis
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			srpi.Close()
			os.Exit(0)
		}
	}()

	for j := 0; j < 5; j++ {
	log.Println("press some keys")
	srpi.Load()
	time.Sleep(3 * time.Second)
	log.Println("Now we shift")
	for i := 0; i < 8; i++ {
		srpi.Shift()
		log.Println("shift", i)
		time.Sleep(1 * time.Second)
	}
	}
	srpi.Close()
	return srpi

}
