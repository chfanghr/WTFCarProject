package generalserialhost

import (
	"errors"
	"github.com/chfanghr/Backend/hardware"
	"sync"
)

type GeneralSerialHost struct {
	c   hardware.Controller
	m   sync.Locker
	num hardware.SerialBusIdentifier
}

func (e *GeneralSerialHost) withMutex(job func() (interface{}, error)) (res interface{}, err error) {
	e.m.Lock()
	defer e.m.Unlock()
	res, err = job()
	return
}
func (g *GeneralSerialHost) Available() (int, error) {
	res, err := g.withMutex(func() (i interface{}, e error) {
		req := hardware.NewCommandRequest(hardware.CommandSerialHost, hardware.SerialHostAvailable, g.num)
		res, err := g.c.Command(*req)
		if err != nil {
			return nil, err
		}
		err = res.Check(hardware.CommandSerialHost)
		if err != nil {
			return nil, err
		}
		return res.GetParameter(), nil
	})
	if err != nil {
		return 0, err
	}
	return res.([]hardware.CommandParameter)[0].(int), nil
}
func (g *GeneralSerialHost) AvailableToWrite() (int, error) {
	res, err := g.withMutex(func() (i interface{}, e error) {
		req := hardware.NewCommandRequest(hardware.CommandSerialHost, hardware.SerialHostAvailableToWrite, g.num)
		res, err := g.c.Command(*req)
		if err != nil {
			return nil, err
		}
		err = res.Check(hardware.CommandSerialHost)
		if err != nil {
			return nil, err
		}
		return res.GetParameter(), nil
	})
	if err != nil {
		return 0, err
	}
	return res.([]hardware.CommandParameter)[0].(int), nil
}
func (g *GeneralSerialHost) SetTimeout(t int) error {
	_, err := g.withMutex(func() (i interface{}, e error) {
		req := hardware.NewCommandRequest(hardware.CommandSerialHost, hardware.SerialHostSetTimeout, g.num, t)
		res, err := g.c.Command(*req)
		if err != nil {
			return nil, err
		}
		err = res.Check(hardware.CommandSerialHost)
		if err != nil {
			return nil, err
		}
		return res.GetParameter(), nil
	})
	if err != nil {
		return err
	}
	return nil
}
func (g *GeneralSerialHost) WriteAndRead(w, r []byte) (wr int, rr int, err error) {
	res, err := g.withMutex(func() (i interface{}, e error) {
		res, err := func() (i interface{}, e error) {
			req := hardware.NewCommandRequest(hardware.CommandSerialHost, hardware.SerialHostWrite, g.num, w)
			res, err := g.c.Command(*req)
			if err != nil {
				return nil, err
			}
			err = res.Check(hardware.CommandSerialHost)
			if err != nil {
				return nil, err
			}
			return res.GetParameter(), nil
		}()
		if err != nil {
			return struct {
				wr, rr int
			}{
				wr: 0,
				rr: 0,
			}, err
		}
		wr := res.([]hardware.CommandParameter)[0].(int)
		res, err = func() (i interface{}, e error) {
			req := hardware.NewCommandRequest(hardware.CommandSerialHost, hardware.SerialHostReadBytes, g.num, len(r))
			res, err := g.c.Command(*req)
			if err != nil {
				return nil, err
			}
			err = res.Check(hardware.CommandSerialHost)
			if err != nil {
				return nil, err
			}
			return res.GetParameter(), nil
		}()
		if err != nil {
			return struct {
				wr, rr int
			}{
				wr: wr,
				rr: 0,
			}, err
		}
		rr := res.([]hardware.CommandParameter)[0].(int)
		return struct {
			wr, rr int
		}{
			wr: wr,
			rr: rr,
		}, nil
	})
	return res.(struct{ wr, rr int }).wr, res.(struct{ wr, rr int }).rr, err
}
func (g *GeneralSerialHost) Begin(baud int) error {
	_, err := g.withMutex(func() (i interface{}, e error) {
		req := hardware.NewCommandRequest(hardware.CommandSerial, hardware.SerialHostBegin, g.num, baud)
		res, err := g.c.Command(*req)
		if err != nil {
			return nil, err
		}
		err = res.Check(hardware.CommandSerial)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}
func (g *GeneralSerialHost) End() error {
	_, err := g.withMutex(func() (i interface{}, e error) {
		req := hardware.NewCommandRequest(hardware.CommandSerial, hardware.SerialHostEnd, g.num)
		res, err := g.c.Command(*req)
		if err != nil {
			return nil, err
		}
		err = res.Check(hardware.CommandSerial)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	return err
}
func (g *GeneralSerialHost) Read(p []byte) (n int, err error) {
	res, err := g.withMutex(func() (i interface{}, e error) {
		req := hardware.NewCommandRequest(hardware.CommandSerialHost, hardware.SerialHostReadBytes, g.num, len(p))
		res, err := g.c.Command(*req)
		if err != nil {
			return nil, err
		}
		err = res.Check(hardware.CommandSerialHost)
		if err != nil {
			return nil, err
		}
		return res.GetParameter(), nil
	})
	if err != nil {
		return 0, err
	}
	//unused res[0]:len of data (n bytes) that controller read from serial port
	//instead,use len()
	l := len(res.([]hardware.CommandParameter)[1].(hardware.SerialData))
	i := copy(p, res.([]hardware.CommandParameter)[1].(hardware.SerialData))
	if l > len(p) {
		return i, errors.New("host error")
	}
	return i, nil
}
func (g *GeneralSerialHost) Write(p []byte) (n int, err error) {
	res, err := g.withMutex(func() (i interface{}, e error) {
		req := hardware.NewCommandRequest(hardware.CommandSerialHost, hardware.SerialHostWrite, g.num, p)
		res, err := g.c.Command(*req)
		if err != nil {
			return nil, err
		}
		err = res.Check(hardware.CommandSerialHost)
		if err != nil {
			return nil, err
		}
		return res.GetParameter(), nil
	})
	if err != nil {
		return 0, err
	}
	return res.([]hardware.CommandParameter)[0].(int), nil
}

func NewGeneralSerialHost(c hardware.Controller, num hardware.SerialBusIdentifier) *GeneralSerialHost {
	return &GeneralSerialHost{c: c, num: num, m: new(sync.Mutex)}
}
