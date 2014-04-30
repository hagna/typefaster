package main

import (
	"github.com/davecheney/gpio"
	"github.com/davecheney/gpio/rpi"
	rpi2 "github.com/davecheney/rpi"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

type Srpi struct {
	shld    gpio.Pin
	clkinh  gpio.Pin
	clk     gpio.Pin
	chassis gpio.Pin
	curpin  int
	npins   int
	M       chan func()
	sync.RWMutex
}

/*
	cycle the clock return value of chassis after rising edge
*/
func (s *Srpi) clock() (res bool) {
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
func (s *Srpi) Shift() {
	time.Sleep(100 * time.Nanosecond)
	s.shld.Set()
	time.Sleep(120 * time.Nanosecond)
	s.clkinh.Clear()
	time.Sleep(120 * time.Nanosecond)
}

/*
	Load
*/
func (s *Srpi) Load() {
	s.shld.Clear()
	time.Sleep(120 * time.Nanosecond)
	s.shld.Set()
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
	log.Println("you would think true since we detect a rising edge", s.chassis.Get())
	log.Println("chassis")
}

func lchassis(msg string) {
	f := rpi2.GPIOGet(rpi.GPIO22)
	if f {
		log.Println(msg)
	}
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
	srpi.M = make(chan func())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			srpi.Close()
			os.Exit(0)
		}
	}()

/*	err = srpi.chassis.BeginWatch(gpio.EdgeRising, srpi.chassis_cb)
	if err != nil {
		log.Fatal("could not watch", err)
	}
*/
	go func() {
		for {
			M := <-srpi.M
			M()
			
			srpi.clock()
			time.Sleep(10 * time.Millisecond)
		}
	}()

	go func() {
		for {
			srpi.Load()
			time.Sleep(10 * time.Millisecond)
			srpi.clock()
			for i := 0; i < 8; i++ {
				srpi.Shift()
				if srpi.clock() {
					log.Println("button", i)
				}
			}
		}
	}()
	select {}
	return srpi

}
