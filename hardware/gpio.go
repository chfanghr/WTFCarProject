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
	GPIO_HIGH PinValue = iota
	GPIO_LOW
)
const (
	GPIO_INPUT PinValue = iota
	GPIO_INPUT_PULLUP
	GPIO_INPUT_PULLDOWN
	GPIO_OUTPUT
)
const (
	GPIO_PinMode CommandMethod = iota
	GPIO_DigitalWrite
	GPIO_DigitalRead
	GPIO_AnalogWrite
	GPIO_AnalogRead
)
