package main

import (
	"github.com/chfanghr/backend/car"
	"github.com/chfanghr/backend/hardware"
	"github.com/chfanghr/backend/location"
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
	return *location.NewPoint2d(rand.Float64(), rand.Float64()), nil
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
