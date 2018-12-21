package main

import (
	"github.com/chfanghr/Backend/car"
	"github.com/chfanghr/Backend/hardware"
	"github.com/chfanghr/Backend/location"
	"log"
	"math/rand"
)

type FakeCar struct {
	l *log.Logger
}

func NewFakeCar(l *log.Logger) *FakeCar {
	return &FakeCar{
		l: l,
	}
}
func (f *FakeCar) GetLocation() (location.Point2D, error) {
	f.l.Println("FakeCar.GetLocation() called")
	return *location.NewPoint2D(rand.Float64(), rand.Float64()), nil
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
