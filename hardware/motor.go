package hardware

import "time"

type Motor interface {
	Forward(time.Duration) error
	Backward(time.Duration) error
	ChopForward(PinValue, time.Duration) error
	ChopReverse(PinValue, time.Duration) error
	Brake() error
	Coast() error
}
