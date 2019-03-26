package a4950

import (
	"github.com/chfanghr/WTFCarProject/hardware"
	"sync"
	"time"
)

const FOREVER = time.Duration(0)

func NewA4950(IN1, IN2 hardware.PWMPin) *A4950 {
	return &A4950{
		m:  new(sync.Mutex),
		i1: IN1,
		i2: IN2,
	}
}

type A4950 struct {
	m      sync.Locker
	i1, i2 hardware.PWMPin
}

func checkErrors(e ...error) error {
	for _, i := range e {
		if i != nil {
			return i
		}
	}
	return nil
}

func (a *A4950) Forward(d time.Duration) error {
	return a.withMutex(func() error {
		err1 := a.i1.DigitalWrite(hardware.GpioHigh)
		err2 := a.i2.DigitalWrite(hardware.GpioLow)
		if err := checkErrors(err1, err2); err != nil {
			a.Brake()
			return err
		}
		if d == 0 {
			return nil
		}
		<-time.After(d)
		return a.brake()
	})
}

func (a *A4950) Backward(d time.Duration) error {
	return a.withMutex(func() error {
		err1 := a.i1.DigitalWrite(hardware.GpioLow)
		err2 := a.i2.DigitalWrite(hardware.GpioHigh)
		if err := checkErrors(err1, err2); err != nil {
			_ = a.brake()
			return err
		}

		if d == 0 {
			return nil
		}
		<-time.After(d)
		return a.brake()
	})
}

func (a *A4950) ChopForward(v hardware.PinValue, d time.Duration) error {
	return a.withMutex(func() error {
		err1 := a.i1.DigitalWrite(hardware.GpioHigh)
		err2 := a.i2.AnalogWrite(v)
		if err := checkErrors(err1, err2); err != nil {
			_ = a.brake()
			return err
		}
		if d == 0 {
			return nil
		}
		<-time.After(d)
		return a.brake()
	})
}

func (a *A4950) ChopReverse(v hardware.PinValue, d time.Duration) error {
	return a.withMutex(func() error {
		err1 := a.i2.DigitalWrite(hardware.GpioHigh)
		err2 := a.i1.AnalogWrite(v)
		if err := checkErrors(err1, err2); err != nil {
			_ = a.brake()
			return err
		}

		if d == 0 {
			return nil
		}
		<-time.After(d)
		return a.brake()
	})
}

func (a *A4950) Brake() error {
	return a.withMutex(func() error {
		return a.brake()
	})
}

func (a *A4950) brake() error {
	err1 := a.i1.DigitalWrite(hardware.GpioHigh)
	err2 := a.i2.DigitalWrite(hardware.GpioHigh)
	return checkErrors(err1, err2)
}

func (a *A4950) Coast() error {
	return a.withMutex(func() error {
		return a.coast()
	})
}

func (a *A4950) coast() error {
	err1 := a.i1.DigitalWrite(hardware.GpioLow)
	err2 := a.i2.DigitalWrite(hardware.GpioLow)
	return checkErrors(err1, err2)
}
func (a *A4950) withMutex(job func() error) error {
	a.m.Lock()
	defer a.m.Unlock()
	err := job()
	return err
}
