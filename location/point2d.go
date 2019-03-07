package location

import "math"

const (
	Yaxis int8 = iota
	antiYaxis
	Xaxis
	antiXaxis
)

type Point2D struct {
	x, y   float64
	Facing int8
}

func NewPoint2D(x float64, y float64) *Point2D {
	return &Point2D{x: x, y: y}
}

func (p *Point2D) SetX(v float64) {
	p.x = v
}

func (p *Point2D) SetY(v float64) {
	p.y = v
}

func (p Point2D) GetX() float64 {
	return p.x
}

func (p Point2D) GetY() float64 {
	return p.y
}

func (p Point2D) DistanceTo(d Point2D) float64 {
	return math.Sqrt((p.x-d.x)*(p.x-d.x) + (p.y-d.y)*(p.y-d.y))
}
