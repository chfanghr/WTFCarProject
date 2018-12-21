package generalir

import (
	"github.com/chfanghr/Backend/hardware"
	"sync"
)

type GeneralIR struct {
	c hardware.Controller
	m sync.Locker
	p hardware.PinNumber
}

func NewGeneralIR(c hardware.Controller, p hardware.PinNumber) *GeneralIR {
	return &GeneralIR{
		c: c,
		p: p,
		m: new(sync.Mutex),
	}
}
func (g *GeneralIR) Send(d hardware.IRData) error {
	req := NewIRRequester(g.p, hardware.IR_SendData, d)
	return req.Commit(g.c)
}
