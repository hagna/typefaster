package main

import (
	"github.com/davecheney/gpio"
	"github.com/davecheney/gpio/rpi"
	"log"
	"os"
	"os/signal"
	"time"
)

type Srpi struct {
	chassis gpio.Pin
	serial  gpio.Pin
	clock   gpio.Pin
	curpin  int
	npins   int
}

/*
   cycles clock high then low with
   the right pulse width and setup time
   according to datasheet for SN74HC595
*/
func (s *Srpi) cycle_clock() {
	// setup time T_su
	time.Sleep(125 * time.Nanosecond)
	s.clock.Set()
	// pulse time T_w
	time.Sleep(120 * time.Nanosecond)
	s.clock.Clear()
}

/*
   Clear the shift register(s)
*/
func (s *Srpi) clearit() {
	s.serial.Clear()
	for i := 0; i < s.npins; i++ {
		s.cycle_clock()
	}
}

/*
   Close everything and clear the shift register(s)
*/
func (s *Srpi) Close() {
	s.clearit()
	s.serial.Close()
	s.clock.Close()
	s.chassis.EndWatch()
	s.chassis.Close()
}

/*
   The callback called when chassis gets a signal,
   used for detecting which key.
*/
func (s *Srpi) chassis_cb() {
	log.Println("button pressed", s.curpin)
}

func NewSrpi() *Srpi {
	ser, err := gpio.OpenPin(rpi.GPIO24, gpio.ModeOutput)
	if err != nil {
		log.Fatal("Error opening pin", err)
	}
	clk, err := gpio.OpenPin(rpi.GPIO25, gpio.ModeOutput)
	if err != nil {
		log.Fatal("Error opening pin", err)
	}
	pin, err := gpio.OpenPin(rpi.GPIO23, gpio.ModeInput)
	if err != nil {
		log.Fatal("Error opening pin", err)
	}

	srpi := new(Srpi)
	srpi.serial = ser
	srpi.clock = clk
	srpi.chassis = pin
	srpi.curpin = 0
	srpi.npins = 8
	srpi.clearit()

	srpi.chassis.BeginWatch(gpio.EdgeRising, srpi.chassis_cb)
	// turn the led off on exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			srpi.Close()
			os.Exit(0)
		}
	}()

	for {
		srpi.curpin = 0
		log.Println("Set serial high to feed a bit to shift register")
		srpi.serial.Set()
		srpi.cycle_clock()
		srpi.serial.Clear()
		log.Println("clock", srpi.curpin)
		for i := 0; i < srpi.npins-1; i++ {
			srpi.cycle_clock()
			time.Sleep(1 * time.Second)
			srpi.curpin = i + 1
			log.Println("clock", srpi.curpin)
		}
	}
	return srpi

}
