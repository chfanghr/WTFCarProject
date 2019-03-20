package hardware

type PinNumber uint8
type PinValue uint8
type GPIO interface {
	PinMode(PinValue) error
	DigitalWrite(PinValue) error
	DigitalRead() (PinValue, error)
}
type AnalogPin interface {
	GPIO
	AnalogRead() (PinValue, error)
}
type PWMPin interface {
	GPIO
	AnalogWrite(PinValue) error
}
type GPIOHub interface {
	PinMode(PinNumber, PinValue) error
	DigitalWrite(PinNumber, PinValue) error
	DigitalRead(PinNumber) (PinValue, error)
	AnalogWrite(PinNumber, PinValue) error
	AnalogRead(PinNumber) (PinValue, error)
}

const (
	GpioLow PinValue = iota
	GpioHigh
)
const (
	GpioInput PinValue = iota
	GpioInputPullup
	GpioInputPulldown
	GpioOutput
)
const (
	GpioPinmode CommandMethod = iota
	GpioDigitalwrite
	GpioDigitalread
	GpioAnalogwrite
	GpioAnalogread
)
