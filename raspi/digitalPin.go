package raspi

import (
	"errors"
	"sync"
)

//To satisfy arduino API
const (
	//INPUT pin mode INPUT
	INPUT = 0 //"in"
	//OUTPUT pin mode OUTPUT
	OUTPUT = 1
	//HIGH pin value HIGH
	HIGH = 1
	//LOW pin value LOW
	LOW = 0

	_GPIOClassPath = "/sys/class/gpio"
)

//Errors that are used by the package

//ErrPinNotExported When a given hasn't been imported,this error will be generated.
var ErrPinNotExported = errors.New("Given pin has not been exported")

//ErrInvalidPinValue When a given value is invalid,this error will be generated.
var ErrInvalidPinValue = errors.New("Given value is invalid")

//ErrInvalidPinMode When a given pin mode is invalid,this error will be generated
var ErrInvalidPinMode = errors.New("Given mode is invalid")

//DigitalPin This structure represents a digital pin.
type DigitalPin struct {
	realPin uint8
	useable bool
	lock    *sync.Mutex
}

//DigitalWrite This function sets the value of the pin.
func (_this *DigitalPin) DigitalWrite(value uint8) error {
	_this.lock.Lock()
	defer _this.lock.Unlock()

	if !isPinExported(_this.realPin) || !_this.useable {
		exportPin(_this.realPin)
		if !isPinExported(_this.realPin) {
			return ErrPinNotExported
		}
	}
	return digitalWrite(_this.realPin, value)
}

//DigitalRead This function gets the value of the pin
func (_this *DigitalPin) DigitalRead() (uint8, error) {
	_this.lock.Lock()
	defer _this.lock.Unlock()

	if !isPinExported(_this.realPin) || !_this.useable {
		exportPin(_this.realPin)
		if !isPinExported(_this.realPin) {
			return emptyValue, ErrPinNotExported
		}
	}
	return digitalRead(_this.realPin)
}

//SetPinMode This function sets the mode of the pin.
func (_this *DigitalPin) SetPinMode(mode uint8) error {
	_this.lock.Lock()
	defer _this.lock.Unlock()

	if !isPinExported(_this.realPin) || !_this.useable {
		exportPin(_this.realPin)
		if !isPinExported(_this.realPin) {
			return ErrPinNotExported
		}
	}
	return setPinMode(_this.realPin, mode)
}

//GetPinMode This function gets the mode of the pin.
func (_this *DigitalPin) GetPinMode() (uint8, error) {
	_this.lock.Lock()
	defer _this.lock.Unlock()
	if !isPinExported(_this.realPin) || !_this.useable {
		exportPin(_this.realPin)
		if !isPinExported(_this.realPin) {
			return emptyMode, ErrPinNotExported
		}
	}
	return getPinMode(_this.realPin)
}

//IsUseAble This function checks the pin's status.
func (_this *DigitalPin) IsUseAble() bool {
	_this.lock.Lock()
	defer _this.lock.Unlock()
	if _this.useable {
		if !isPinExported(_this.realPin) {
			exportPin(_this.realPin)
			if !isPinExported(_this.realPin) {
				return false
			}
		}
		return true
	}
	return false
}

func (_this *DigitalPin) PinMode(m uint8) error {
	return _this.SetPinMode(m)
}
