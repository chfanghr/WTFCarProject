package location

import "math/rand"

const (
	LocationReqFailed = iota
)

type Engine interface {
	GetLocation() (Point2D, error)
}

type FakeLocationEngine struct{}

func NewFakeLocationEngine() *FakeLocationEngine {
	return &FakeLocationEngine{}
}

func (FakeLocationEngine) GetLocation() (Point2D, error) {
	return *NewPoint2D(rand.Float64(), rand.Float64()), nil
}
