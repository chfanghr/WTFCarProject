package raspi

import (
	"errors"
	"github.com/chfanghr/WTFCarProject/hardware"
	"sync"
)

func NewRaspiController() *RaspiController {
	return &RaspiController{
		m: new(sync.Mutex),
	}
}

type RaspiController struct {
	m sync.Locker
	p *SerialPort
}

func (r *RaspiController) withMutex(job func() (interface{}, error)) (interface{}, error) {
	r.m.Lock()
	defer r.m.Unlock()
	res, err := job()
	return res, err
}

func (c RaspiController) IsValidPin(p hardware.PinNumber) error {
	if _, ok := pins[uint8(p)]; !ok {
		return ErrNorAValidPin
	}
	return nil
}

func (c *RaspiController) Command(req hardware.CommandRequest) (hardware.CommandResponse, error) {
	res, err := c.withMutex(func() (interface{}, error) {
		switch req.CommType {
		case hardware.CommandGpio:
			dp, err := ExportPin(req.CommParameter[0].(uint8))
			if err != nil {
				return *hardware.NewCommandResponse(hardware.CommandGpio, hardware.OperationFailed), err
			} else {
				switch req.CommMethod {
				case hardware.GpioPinmode:
					switch req.CommParameter[1].(hardware.PinValue) {
					case hardware.GpioInput:
						err = dp.SetPinMode(INPUT)
						break
					case hardware.GpioOutput:
						err = dp.SetPinMode(OUTPUT)
						break
					default:
						err = errors.New("unsupported mode")
					}
					if err != nil {
						return *hardware.NewCommandResponse(hardware.CommandGpio, hardware.OperationFailed), err
					} else {
						return *hardware.NewCommandResponse(hardware.CommandGpio, hardware.OperationSucceeded), nil
					}
				case hardware.GpioDigitalwrite:
					switch req.CommParameter[1].(hardware.PinValue) {
					case hardware.GpioHigh:
						err = dp.DigitalWrite(HIGH)
						break
					case hardware.GpioLow:
						err = dp.DigitalWrite(LOW)
						break
					default:
						err = errors.New("unsupported value")
					}
					if err != nil {
						return *hardware.NewCommandResponse(hardware.CommandGpio, hardware.OperationFailed), err
					} else {
						return *hardware.NewCommandResponse(hardware.CommandGpio, hardware.OperationSucceeded), nil
					}
				case hardware.GpioDigitalread:
					v, err := dp.DigitalRead()
					if err != nil {
						return *hardware.NewCommandResponse(hardware.CommandGpio, hardware.OperationFailed, 0), err
					} else {
						switch v {
						case HIGH:
							v = uint8(hardware.GpioHigh)
							break
						case LOW:
							v = uint8(hardware.GpioLow)
							break
						}
						return *hardware.NewCommandResponse(hardware.CommandGpio, hardware.OperationSucceeded, v), nil
					}
				default:
					return *hardware.NewCommandResponse(hardware.CommandGpio, hardware.OperationFailed), errors.New("unsupported method")
				}
			}
		/*case hardware.CommandSerialHost:
		switch req.CommMethod {
		case hardware.SerialHostBegin:
			c:= SerialConfig{
				Name:fmt.Sprint("/dev/ttyUSB",)
			}
		}*/
		default:
			return *hardware.NewCommandResponse(hardware.CommandGpio, hardware.OperationFailed), errors.New("unsupported type")
		}
	})
	return res.(hardware.CommandResponse), err
}
