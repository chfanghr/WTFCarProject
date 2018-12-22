package generalgpio

import (
	"errors"
	. "github.com/chfanghr/Backend/hardware"
)

type GeneralGPIORequester struct {
	m   CommandMethod
	p   PinNumber
	v   PinValue
	res PinValue
}

func NewGeneralGpioRequester(m CommandMethod, p PinNumber, v PinValue) *GeneralGPIORequester {
	return &GeneralGPIORequester{
		m:   m,
		p:   p,
		v:   v,
		res: 0,
	}
}
func (g *GeneralGPIORequester) Commit(c Controller) error {
	res, err := func() (CommandResponse, error) {
		switch g.m {
		case GpioPinmode:
			switch g.v {
			case GpioInput:
				return c.Command(*NewCommandRequest(CommandGpio, g.m, g.p, g.v))
			case GpioOutput:
				return c.Command(*NewCommandRequest(CommandGpio, g.m, g.p, g.v))
			case GpioInputPullUp:
				return c.Command(*NewCommandRequest(CommandGpio, g.m, g.p, g.v))
			case GpioInputPullDown:
				return c.Command(*NewCommandRequest(CommandGpio, g.m, g.p, g.v))
			default:
				return CommandResponse{}, errors.New("invalid mode")
			}
		case GpioDigitalWrite:
			switch g.v {
			case GpioHigh:
				return c.Command(*NewCommandRequest(CommandGpio, g.m, g.p, g.v))
			case GpioLow:
				return c.Command(*NewCommandRequest(CommandGpio, g.m, g.p, g.v))
			default:
				return CommandResponse{}, errors.New("invalid value")
			}
		case GpioAnalogWrite:
			return c.Command(*NewCommandRequest(CommandGpio, g.m, g.p, g.v))
		case GpioAnalogRead:
			return c.Command(*NewCommandRequest(CommandGpio, g.m, g.p))
		case GpioDigitalRead:
			return c.Command(*NewCommandRequest(CommandGpio, g.m, g.p))
		default:
			return CommandResponse{}, errors.New("unsupported method")
		}
	}()
	if err != nil {
		return err
	}
	err = res.Check(CommandGpio)
	if err != nil {
		return err
	}
	switch g.m {
	case GpioDigitalRead:
		g.res = res.GetParameter()[0].(PinValue)
		break
	case GpioAnalogRead:
		g.res = res.GetParameter()[0].(PinValue)
		break
	}
	return nil
}
func (g GeneralGPIORequester) GetRes() interface{} {
	return g.res
}
