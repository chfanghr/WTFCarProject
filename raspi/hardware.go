package raspi

import (
	"errors"
	"github.com/chfanghr/Backend/hardware"
)

type RaspiController struct{}

func (c RaspiController) IsValidPin(p hardware.PinNumber) error {
	if _, ok := pins[uint8(p)]; !ok {
		return ErrNorAValidPin
	}
	return nil
}
func (c RaspiController) Command(req hardware.CommandRequest) (hardware.CommandResponse, error) {
	switch req.CommType {
	case hardware.Command_GPIO:
		dp, err := ExportPin(req.CommParameter[2].(uint8))
		if err != nil {
			return hardware.CommandResponse{}, err
		}
		switch req.CommMethod {
		case hardware.GPIO_PinMode:
			switch req.CommParameter[3].(hardware.PinValue) {
			case hardware.GPIO_INPUT:
				err = dp.SetPinMode(INPUT)
			case hardware.GPIO_OUTPUT:
				err = dp.SetPinMode(OUTPUT)
			default:
				err = errors.New("unsupported mode")
			}

			if err != nil {
				return *hardware.NewCommandResponse(hardware.Command_GPIO, hardware.OperationFailed), err
			}
			return *hardware.NewCommandResponse(hardware.Command_GPIO, hardware.OperationSucceeded), nil
		case hardware.GPIO_DigitalWrite:
			switch req.CommParameter[3].(hardware.PinValue) {
			case hardware.GPIO_HIGH:
				err = dp.DigitalWrite(HIGH)
			case hardware.GPIO_LOW:
				err = dp.DigitalWrite(LOW)
			default:
				err = errors.New("unsupported value")
			}
			if err != nil {
				if err != nil {
					return *hardware.NewCommandResponse(hardware.Command_GPIO, hardware.OperationFailed), err
				}
				return *hardware.NewCommandResponse(hardware.Command_GPIO, hardware.OperationSucceeded), nil
			}
		case hardware.GPIO_DigitalRead:
			v, err := dp.DigitalRead()
			if err != nil {
				if err != nil {
					return *hardware.NewCommandResponse(hardware.Command_GPIO, hardware.OperationFailed, 0), err
				}
				return *hardware.NewCommandResponse(hardware.Command_GPIO, hardware.OperationSucceeded, v), nil
			}
		default:
			return hardware.CommandResponse{}, errors.New("unsupported method")
		}
	default:
		return hardware.CommandResponse{}, errors.New("unsupported type")
	}
	return hardware.CommandResponse{}, errors.New("unknown error")
}
