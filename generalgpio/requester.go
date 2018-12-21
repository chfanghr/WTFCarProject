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
		case GPIO_PinMode:
			switch g.v {
			case GPIO_INPUT:
				return c.Command(*NewCommandRequest(Command_GPIO, g.m, g.p, g.v))
			case GPIO_OUTPUT:
				return c.Command(*NewCommandRequest(Command_GPIO, g.m, g.p, g.v))
			case GPIO_INPUT_PULLUP:
				return c.Command(*NewCommandRequest(Command_GPIO, g.m, g.p, g.v))
			case GPIO_INPUT_PULLDOWN:
				return c.Command(*NewCommandRequest(Command_GPIO, g.m, g.p, g.v))
			default:
				return CommandResponse{}, errors.New("invalid mode")
			}
		case GPIO_DigitalWrite:
			switch g.v {
			case GPIO_HIGH:
				return c.Command(*NewCommandRequest(Command_GPIO, g.m, g.p, g.v))
			case GPIO_LOW:
				return c.Command(*NewCommandRequest(Command_GPIO, g.m, g.p, g.v))
			default:
				return CommandResponse{}, errors.New("invalid value")
			}
		case GPIO_AnalogWrite:
			return c.Command(*NewCommandRequest(Command_GPIO, g.m, g.p, g.v))
		case GPIO_AnalogRead:
			return c.Command(*NewCommandRequest(Command_GPIO, g.m, g.p))
		case GPIO_DigitalRead:
			return c.Command(*NewCommandRequest(Command_GPIO, g.m, g.p))
		default:
			return CommandResponse{}, errors.New("unsupported method")
		}
	}()
	if err == nil {
		return err
	}
	err = res.Check(Command_GPIO)
	if err != nil {
		return err
	}
	switch g.m {
	case GPIO_PinMode:
		g.res = res.GetParameter()[2].(PinValue)
	case GPIO_DigitalWrite:
		g.res = res.GetParameter()[2].(PinValue)
	case GPIO_AnalogWrite:
		g.res = res.GetParameter()[2].(PinValue)
	}
	return nil
}
func (g GeneralGPIORequester) GetRes() interface{} {
	return g.res
}
