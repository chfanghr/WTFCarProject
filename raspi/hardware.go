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
		dp, err := ExportPin(req.CommParameter[0].(uint8))
		if err != nil {
			return hardware.CommandResponse{}, err
		} else {
			switch req.CommMethod {
			case hardware.GPIO_PinMode:
				switch req.CommParameter[1].(hardware.PinValue) {
				case hardware.GPIO_INPUT:
					err = dp.SetPinMode(INPUT)
					break
				case hardware.GPIO_OUTPUT:
					err = dp.SetPinMode(OUTPUT)
					break
				default:
					err = errors.New("unsupported mode")
				}
				if err != nil {
					return *hardware.NewCommandResponse(hardware.Command_GPIO, hardware.Operation_Failed), err
				} else {
					return *hardware.NewCommandResponse(hardware.Command_GPIO, hardware.Operation_Succeeded), nil
				}
			case hardware.GPIO_DigitalWrite:
				switch req.CommParameter[1].(hardware.PinValue) {
				case hardware.GPIO_HIGH:
					err = dp.DigitalWrite(HIGH)
					break
				case hardware.GPIO_LOW:
					err = dp.DigitalWrite(LOW)
					break
				default:
					err = errors.New("unsupported value")
				}
				if err != nil {
					return *hardware.NewCommandResponse(hardware.Command_GPIO, hardware.Operation_Failed), err
				} else {
					return *hardware.NewCommandResponse(hardware.Command_GPIO, hardware.Operation_Succeeded), nil
				}
			case hardware.GPIO_DigitalRead:
				v, err := dp.DigitalRead()
				if err != nil {
					return *hardware.NewCommandResponse(hardware.Command_GPIO, hardware.Operation_Failed, 0), err
				} else {
					switch v {
					case HIGH:
						v = uint8(hardware.GPIO_HIGH)
						break
					case LOW:
						v = uint8(hardware.GPIO_LOW)
						break
					}
					return *hardware.NewCommandResponse(hardware.Command_GPIO, hardware.Operation_Succeeded, v), nil
				}
			default:
				return hardware.CommandResponse{}, errors.New("unsupported method")
			}
		}
	default:
		return hardware.CommandResponse{}, errors.New("unsupported type")
	}
}
