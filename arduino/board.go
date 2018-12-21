package arduino

import "github.com/chfanghr/Backend/hardware"

type Board interface {
	IsValidPin(hardware.PinNumber) error
}

type ArduinoProMini struct{}

func (ArduinoProMini) IsValidPin(hardware.PinNumber) error {
	return nil
}
