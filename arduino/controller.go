package arduino

import (
	"encoding/json"
	"github.com/chfanghr/WTFCarProject/hardware"
	"sync"
	"time"
)

func NewArduinoControllerViaI2C(b Board, i hardware.I2C) *ArduinoControllerViaI2c {
	return &ArduinoControllerViaI2c{
		i: i,
		m: new(sync.Mutex),
		b: b,
	}
}

type ArduinoControllerViaI2c struct {
	i hardware.I2C
	m sync.Locker
	b Board
}

func (a *ArduinoControllerViaI2c) withMutex(job func() (interface{}, error)) (res interface{}, err error) {
	a.m.Lock()
	defer a.m.Unlock()
	res, err = job()
	return
}
func (a *ArduinoControllerViaI2c) Command(req hardware.CommandRequest) (hardware.CommandResponse, error) {
	res, err := a.withMutex(func() (interface{}, error) {
		e := json.NewEncoder(a.i)
		err := e.Encode(req)
		if err != nil {
			return nil, err
		}
		<-time.After(30 * time.Millisecond)
		d := json.NewDecoder(a.i)
		v := &hardware.CommandResponse{}
		err = d.Decode(v)
		if err != nil {
			return nil, err
		}
		return *v, nil
	})
	if err != nil {
		return hardware.CommandResponse{}, err
	}
	return res.(hardware.CommandResponse), nil
}
func (a *ArduinoControllerViaI2c) IsValidPin(p hardware.PinNumber) error {
	_, err := a.withMutex(func() (interface{}, error) {
		return nil, a.b.IsValidPin(p)
	})
	return err
}
