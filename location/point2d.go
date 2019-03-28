package location

import (
	"math"
)

type Point2D struct {
	x, y float64
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

func (p Point2D) IsOnSameLine(ps ...Point2D) bool {
	if len(ps) == 1 || len(ps) == 0 {
		return true
	}
	if s, diff := p.IsTheSame(ps...); !s || diff < 3 {
		return true
	}
	for i, p1 := range ps {
		for _, p2 := range ps[i:] {
			if !isPointInSegments(p1, p2, p) {
				return false
			}
		}
	}
	return true
}

func (p Point2D) IsTheSame(ps ...Point2D) (bool, int) {
	if len(ps)+1 == 1 {
		return true, 1
	}
	tmpMap := make(map[Point2D]struct{})
	for _, v := range ps {
		tmpMap[v] = struct{}{}
	}
	tmpMap[p] = struct{}{}
	if len(tmpMap) == 1 {
		return true, 1
	}
	return false, len(tmpMap)
}

func isPointInSegments(Pi Point2D, Pj Point2D, Q Point2D) bool {
	if (Q.x-Pi.x)*(Pj.y-Pi.y) == (Pj.x-Pi.x)*(Q.y-Pi.y) &&
		math.Min(Pi.x, Pj.x) <= Q.x && Q.x <= math.Max(Pi.x, Pj.x) &&
		math.Min(Pi.y, Pj.y) <= Q.x && Q.y <= math.Max(Pi.y, Pj.y) {
		return true
	}
	return false
}
