package main

import (
	"github.com/chfanghr/backend/a4950"
	"github.com/chfanghr/backend/arduino"
	"github.com/chfanghr/backend/generalgpio"
	"github.com/chfanghr/backend/generalir"
	"github.com/chfanghr/backend/hardware"
	"github.com/chfanghr/backend/location"
	"github.com/chfanghr/backend/raspi"
	"sync"
)

type Car struct {
	a      hardware.Controller
	le     location.Engine
	hub    hardware.GPIOHub
	MA, MB hardware.Motor
	ir     hardware.IR
	m      sync.Locker
}

func NewCar(I2CAddr uint8, I2CBus int, MotorAIN1 hardware.PinNumber, MotorAIN2 hardware.PinNumber,
	MotorBIN1 hardware.PinNumber, MotorBIN2 hardware.PinNumber, IRPin hardware.PinNumber, SerialIBeacon string) (*Car, error) {
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

	return c, nil
}
func (c *Car) withMutex(job func() (interface{}, error)) (interface{}, error) {
	c.m.Lock()
	defer c.m.Unlock()
	res, err := job()
	return res, err
}
func (*Car) GetLocation() (location.Point2D, error) {
	panic("implement me")
}
func (*Car) MoveTo(location.Point2D) error {
	panic("implement me")
}
func (*Car) LastMovementStatus() int {
	panic("implement me")
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
