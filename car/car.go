package car

import (
	"github.com/chfanghr/WTFCarProject/hardware"
	"github.com/chfanghr/WTFCarProject/location"
	"github.com/chfanghr/WTFCarProject/rpcprotocal"
)

const (
	Succeeded int = iota
	Failed
	Moving
)

type Car interface {
	location.Engine
	MoveTo(location.Point2D) error
	LastMovementStatus() int
	StopMovement() error
	IRSend(hardware.IRData) error
}

type Service interface {
	GetLocation(int, *rpcprotocal.Point2D) error
	MoveTo(rpcprotocal.Point2D, *int) error
	LastMovementStatus(int, *int) error
	StopMoving(int, *int) error
	IRSend(hardware.IRData, *int) error
}

type Client interface {
	Car
	IsServiceAvailable() error
}
