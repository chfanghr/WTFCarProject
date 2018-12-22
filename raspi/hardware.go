package raspi

import (
	"errors"
	"github.com/chfanghr/Backend/hardware"
	"sync"
)

func NewRaspiController() *RaspiController {
	return &RaspiController{
		m: new(sync.Mutex),
		p: make(map[hardware.SerialBusIdentifier]*SerialPort),
	}
}

type RaspiController struct {
	m sync.Locker
	p map[hardware.SerialBusIdentifier]*SerialPort
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
					if len(req.CommParameter) < 2 {
						return *hardware.NewCommandResponse(hardware.CommandGpio, hardware.OperationFailed), errors.New("len of parameters too short")
					}
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
				case hardware.GpioDigitalWrite:
					if len(req.CommParameter) < 2 {
						return *hardware.NewCommandResponse(hardware.CommandGpio, hardware.OperationFailed), errors.New("len of parameters too short")
					}
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
				case hardware.GpioDigitalRead:
					if len(req.CommParameter) < 1 {
						return *hardware.NewCommandResponse(hardware.CommandGpio, hardware.OperationFailed), errors.New("len of parameters too short")
					}
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
			/*
				case hardware.CommandSerialHost:
					switch req.CommMethod {
					case hardware.SerialHostBegin:
						if len(req.CommParameter)<2{
							return *hardware.NewCommandResponse(hardware.CommandSerialHost, hardware.OperationFailed), errors.New("len of parameters too short")
						}
						if _, ok := c.p[req.CommParameter[0].(hardware.SerialBusIdentifier)]; ok {
							return *hardware.NewCommandResponse(hardware.CommandSerial, hardware.OperationFailed), errors.New("port already begin")
						}
						cf := SerialConfig{
							Name: fmt.Sprint("/dev/ttyUSB", req.CommParameter[0].(hardware.SerialBusIdentifier)),
							Baud: int(req.CommParameter[1].(hardware.SerialBaudRate)),
						}
						p, err := OpenSerialPort(&cf)
						if err != nil {
							return *hardware.NewCommandResponse(hardware.CommandSerial, hardware.OperationFailed), err
						}
						c.p[req.CommParameter[0].(hardware.SerialBusIdentifier)] = p
						return *hardware.NewCommandResponse(hardware.CommandSerial, hardware.OperationSucceeded), nil
					case hardware.SerialHostEnd:
						if len(req.CommParameter)<1{
							return *hardware.NewCommandResponse(hardware.CommandSerialHost, hardware.OperationFailed), errors.New("len of parameters too short")
						}
						if p, ok := c.p[req.CommParameter[0].(hardware.SerialBusIdentifier)]; ok {
							err := p.Close()
							delete(c.p, req.CommParameter[0].(hardware.SerialBusIdentifier))
							if err != nil {
								return *hardware.NewCommandResponse(hardware.CommandSerial, hardware.OperationFailed), err
							}
							return *hardware.NewCommandResponse(hardware.CommandSerial, hardware.OperationSucceeded), nil
						} else {
							return *hardware.NewCommandResponse(hardware.CommandSerial, hardware.OperationFailed), errors.New("port not exist")
						}
					case hardware.SerialHostSetTimeout:
						if len(req.CommParameter)<2{
							return *hardware.NewCommandResponse(hardware.CommandSerialHost, hardware.OperationFailed), errors.New("len of parameters too short")
						}
						if p, ok := c.p[req.CommParameter[0].(hardware.SerialBusIdentifier)]; ok {
							err := p.Close()
							delete(c.p, req.CommParameter[0].(hardware.SerialBusIdentifier))
							if err != nil {
								return *hardware.NewCommandResponse(hardware.CommandSerial, hardware.OperationFailed), err
							}
							cf := SerialConfig{
								Name: fmt.Sprint("/dev/ttyUSB", req.CommParameter[0].(hardware.SerialBusIdentifier)),
								Baud: int(req.CommParameter[1].(hardware.SerialBaudRate)),
								ReadTimeout:req.CommParameter[0]
							}
							p, err := OpenSerialPort(&cf)
							if err != nil {
								return *hardware.NewCommandResponse(hardware.CommandSerial, hardware.OperationFailed), err
							}
							return *hardware.NewCommandResponse(hardware.CommandSerial, hardware.OperationSucceeded), nil
						} else {
							return *hardware.NewCommandResponse(hardware.CommandSerial, hardware.OperationFailed), errors.New("port not exist")
						}
					default:
						return *hardware.NewCommandResponse(hardware.CommandSerial, hardware.OperationFailed), errors.New("unsupported method")
					}
			*/
		default:
			return *hardware.NewCommandResponse(hardware.CommandGpio, hardware.OperationFailed), errors.New("unsupported type")
		}
	})
	return res.(hardware.CommandResponse), err
}
