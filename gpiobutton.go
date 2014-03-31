package main

import (
	"github.com/davecheney/gpio"
	"github.com/davecheney/gpio/rpi"
	"log"
)

type button struct {
	gpio.Pin
	num int		
}

func (b button) cb() {
	log.Println(b.num)
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
	res := button{pin, n}
	err = pin.BeginWatch(gpio.EdgeBoth, res.cb)
	if err != nil {
		pin.Close()
		return nil, err
	}
	return &res, nil
}

