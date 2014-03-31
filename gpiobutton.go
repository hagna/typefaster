package main

import (
	"github.com/davecheney/gpio"
	"github.com/davecheney/gpio/rpi"
	"time"
	"log"
)

type button struct {
	gpio.Pin
	num int		
	c *time.Timer
	edgeChange chan bool
}

func (b button) cb() {
	b.edgeChange <- b.Get()
}

func (b button) Keyup() {
	log.Println("keyup", b.num)
}

func (b button) Keydown() {
	log.Println("keydown", b.num)
}

func (b button) Close() {
	log.Println("Close:", b)
	b.Pin.Close()
	b.Pin.EndWatch()
}	

func NewButton(n int) (*button, error) {
	pin, err := rpi.OpenPin(n, gpio.ModeInput)
	if err != nil {
		return nil, err
	}
	res := button{Pin:pin, num:n, c:time.NewTimer(1 * time.Second), edgeChange:make(chan bool)}
	err = pin.BeginWatch(gpio.EdgeBoth, res.cb)
	if err != nil {
		pin.Close()
		return nil, err
	}
	go func() {
		currentState := false
		keydown := false
		for {
			select {
			case ec := <- res.edgeChange:
				if keydown == true && ec == true {
					keydown = false
					res.Keyup()
				}
				currentState = ec
				res.c.Reset(100 * time.Millisecond)
			case <-res.c.C:
				if currentState == false {
					keydown = true
					res.Keydown()
				}
			}
		}
	}()
	return &res, nil
}

