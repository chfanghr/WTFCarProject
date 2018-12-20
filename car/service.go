package car

import (
	"errors"
	"github.com/chfanghr/backend/hardware"
	"github.com/chfanghr/backend/location"
	"github.com/chfanghr/backend/rpcprotocal"
	"sync"
)

var (
	operationSucceeded = error(nil)
	operationFailed    = errors.New("operation failed")
)

type GeneralServiceHandler struct {
	m sync.Locker
	c Car
}

func NewGeneralServiceHandler(c Car) *GeneralServiceHandler {
	return &GeneralServiceHandler{
		m: new(sync.Mutex),
		c: c,
	}
}
func (s *GeneralServiceHandler) withMutex(job func() error) error {
	s.m.Lock()
	defer s.m.Unlock()
	err := job()
	return err
}
func (s *GeneralServiceHandler) GetLocation(nouse int, d *rpcprotocal.Point2D) error {
	return s.withMutex(func() error {
		l, err := s.c.GetLocation()
		if err != nil {
			return operationFailed
		}
		d.X = l.GetX()
		d.Y = l.GetY()
		return operationSucceeded
	})
}
func (s *GeneralServiceHandler) MoveTo(d rpcprotocal.Point2D, rep *int) error {
	return s.withMutex(func() error {
		go s.c.MoveTo(*location.NewPoint2D(d.X, d.Y))
		*rep = Moving
		return operationSucceeded
	})
}
func (s *GeneralServiceHandler) LastMovementStatus(nouse int, rep *int) error {
	return s.withMutex(func() error {
		*rep = s.c.LastMovementStatus()
		return operationSucceeded
	})
}
func (s *GeneralServiceHandler) StopMoving(int, *int) error {
	return s.withMutex(func() error {
		if s.c.StopMovement() != nil {
			return operationFailed
		}
		return operationSucceeded
	})
}
func (s *GeneralServiceHandler) IRSend(data hardware.IRData, nouse *int) error {
	if len(data) > hardware.IR_DataMaxLen {
		return operationFailed
	}
	return s.withMutex(func() error {
		if s.c.IRSend(data) != nil {
			return operationFailed
		}
		return operationSucceeded
	})
}
