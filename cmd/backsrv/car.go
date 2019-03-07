package main

import (
	"github.com/chfanghr/Backend/a4950"
	"github.com/chfanghr/Backend/arduino"
	"github.com/chfanghr/Backend/car"
	"github.com/chfanghr/Backend/generalgpio"
	"github.com/chfanghr/Backend/generalir"
	"github.com/chfanghr/Backend/hardware"
	"github.com/chfanghr/Backend/location"
	"github.com/chfanghr/Backend/raspi"
	"sync"
	"time"
)

type Car struct {
	a            hardware.Controller
	le           location.Engine
	hub          hardware.GPIOHub
	MA, MB       hardware.Motor
	ir           hardware.IR
	m            sync.Locker
	lastLocation location.Point2D
}

func NewCar(I2CAddr uint8, I2CBus int, MotorAIN1 hardware.PinNumber, MotorAIN2 hardware.PinNumber,
	MotorBIN1 hardware.PinNumber, MotorBIN2 hardware.PinNumber, IRPin hardware.PinNumber, BluetoothHost string, IBeaconName string, IBeaconUUID string) (*Car, error) {
	c := &Car{m: new(sync.Mutex)}
	dev, err := raspi.NewI2C(I2CAddr, I2CBus)
	if err != nil {
		return nil, err
	}

	MA1 := generalgpio.NewGeneralGPIO(c.hub, MotorAIN1)
	MA2 := generalgpio.NewGeneralGPIO(c.hub, MotorAIN2)
	MB1 := generalgpio.NewGeneralGPIO(c.hub, MotorBIN1)
	MB2 := generalgpio.NewGeneralGPIO(c.hub, MotorBIN2)

	c.a = arduino.NewArduinoControllerViaI2C(arduino.ArduinoProMini{}, dev)
	c.hub = generalgpio.NewGeneralGPIOHub(c.a)
	c.MA = a4950.NewA4950(MA1, MA2)
	c.MB = a4950.NewA4950(MB1, MB2)
	c.ir = generalir.NewGeneralIR(c.a, IRPin)
	c.le = location.NewFakeLocationEngine()
	return c, nil
}
func (c *Car) withMutex(job func() (interface{}, error)) (interface{}, error) {
	c.m.Lock()
	defer c.m.Unlock()
	res, err := job()
	return res, err
}
func (c *Car) GetLocation() (location.Point2D, error) {
	res, err := c.withMutex(func() (i interface{}, e error) {
		return c.lastLocation, nil //c.le.GetLocation()
	})
	return res.(location.Point2D), err
}
func (c *Car) MoveTo(l location.Point2D) error {
	_, err := c.withMutex(func() (i interface{}, e error) {
		// TODO:calculate path (last location,new location)
		x := l.GetX() - c.lastLocation.GetX()
		y := l.GetY() - c.lastLocation.GetY()

		c.lastLocation = l

		var goStraight = func(distance float64) {
			const oneMeterA = time.Duration(time.Millisecond * 1000)
			const oneMeterB = time.Duration(time.Millisecond * 1000)
			err := c.MA.Forward(time.Duration(distance) * oneMeterA)
			if err != nil {
				return
			}
			c.MB.Forward(time.Duration(distance) * oneMeterB)
		}
		return nil, nil
	})
	return err
}
func (c *Car) LastMovementStatus() int {
	res, _ := c.withMutex(func() (i interface{}, e error) {
		return car.Succeeded, nil
	})
	return res.(int)
}
func (c *Car) StopMovement() error {
	_, err := c.withMutex(func() (i interface{}, e error) {
		if e = c.MA.Brake(); e != nil {
			return
		}
		if e = c.MB.Brake(); e != nil {
			return
		}
		if e = c.MA.Coast(); e != nil {
			return
		}
		if e = c.MB.Coast(); e != nil {
			return
		}
		return
	})
	return err
}
func (c *Car) IRSend(data hardware.IRData) error {
	_, err := c.withMutex(func() (i interface{}, e error) {
		return nil, c.ir.Send(data)
	})
	return err
}
