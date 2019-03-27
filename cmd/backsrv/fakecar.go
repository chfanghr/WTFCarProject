package main

import (
	"github.com/chfanghr/WTFCarProject/car"
	"github.com/chfanghr/WTFCarProject/hardware"
	"github.com/chfanghr/WTFCarProject/location"
	_ "github.com/gorilla/websocket"
	"log"
)

type FakeCar struct {
	current *location.Point2D
	//TODO Store last movement status here
	l *log.Logger
}

func NewFakeCar(l *log.Logger) *FakeCar {
	//TODO Create a goroutine inside this constructor to provide websocket service and other staff
	return &FakeCar{
		l: l,
	}
}
func (f *FakeCar) GetLocation() (location.Point2D, error) {
	f.l.Println("FakeCar.GetLocation() called")
	return *f.current, nil
}
func (f *FakeCar) MoveTo(l location.Point2D) error {
	f.l.Println("FakeCar.MoveTo() called:", l)
	return nil
}
func (f *FakeCar) LastMovementStatus() int {
	f.l.Println("LastMovementStatus() called")
	return car.Succeeded
}
func (f *FakeCar) StopMovement() error {
	f.l.Println(" StopMovement() called")
	return nil
}
func (f *FakeCar) IRSend(ir hardware.IRData) error {
	f.l.Println("IRSend() called")
	return nil
}
