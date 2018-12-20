package location

import (
	"errors"
	"math"
)

type Circle struct {
	c Point2D
	r float64
}

func NewCircle(c Point2D, r float64) (*Circle, error) {
	if r < 0 {
		return nil, errors.New("invalid r")
	}
	return &Circle{c: c, r: r}, nil
}

func (c *Circle) SetCore(d Point2D) {
	c.c = d
}

func (c *Circle) SetR(r float64) error {
	if r < 0 {
		return errors.New("invalid r")
	}
	c.r = r
	return nil
}

func (c Circle) GetCore() Point2D {
	return c.c
}

func (c Circle) GetR() float64 {
	return c.r
}

func (c Circle) GetSamplePoints() (res []Point2D) {
	for i := float64(0); i <= 360; i += 20.0 {
		res = append(res, *NewPoint2D(math.Cos(i)*c.r+c.c.x, math.Sin(i)*c.r+c.c.y))
	}
	return
}
