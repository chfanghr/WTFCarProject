package arduino

import "github.com/chfanghr/WTFCarProject/hardware"

type Board interface {
	IsValidPin(hardware.PinNumber) error
}

type ArduinoProMini struct{}

func (ArduinoProMini) IsValidPin(hardware.PinNumber) error {
	return nil
}
