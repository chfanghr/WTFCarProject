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
			case GpioInput, GpioOutput, GpioInputPullup, GpioInputPulldown:
				return c.Command(*NewCommandRequest(CommandGpio, g.m, g.p, g.v))
			default:
				return CommandResponse{}, errors.New("invalid mode")
			}
		case GpioDigitalwrite:
			switch g.v {
			case GpioHigh:
				return c.Command(*NewCommandRequest(CommandGpio, g.m, g.p, g.v))
			case GpioLow:
				return c.Command(*NewCommandRequest(CommandGpio, g.m, g.p, g.v))
			default:
				return CommandResponse{}, errors.New("invalid value")
			}
		case GpioAnalogwrite:
			return c.Command(*NewCommandRequest(CommandGpio, g.m, g.p, g.v))
		case GpioAnalogread:
			return c.Command(*NewCommandRequest(CommandGpio, g.m, g.p))
		case GpioDigitalread:
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
	case GpioDigitalread, GpioAnalogread:
		if len(res.GetParameter()) == 0 {
			return errors.New("invalid len of paramerters")
		}
		g.res = res.GetParameter()[0].(PinValue)
		break
	}
	return nil
}
func (g GeneralGPIORequester) GetRes() interface{} {
	return g.res
}
