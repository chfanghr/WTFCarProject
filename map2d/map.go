package map2d

import (
	"encoding/json"
	"github.com/chfanghr/WTFCarProject/grid"
	"github.com/chfanghr/WTFCarProject/location"
	"math"
	"sort"
)

type Map2D struct {
	Map struct {
		Size struct {
			X float64 `json:"x"`
			Y float64 `json:"y"`
		} `json:"size"`
		Barriers []struct {
			Required [3]location.Point2D `json:"required"`
			Optional []location.Point2D  `json:"optional"`
		} `json:"barriers"`
	} `json:"map"`
}

func NewMap2d(raw []byte) (*Map2D, error) {
	res := &Map2D{}
	err := json.Unmarshal(raw, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (m *Map2D) isValid() bool {
	return m.Map.Size.X > 0 && m.Map.Size.Y > 0 &&
		func() bool {
			for _, b := range m.Map.Barriers {
				isOutOfMap := func(p location.Point2D) bool {
					return !(p.GetX() > m.Map.Size.X || p.GetY() > m.Map.Size.Y || p.GetX() < 0 || p.GetY() < 0)
				}
				for _, p := range b.Required {
					if !isOutOfMap(p) {
						return false
					}
				}
				for _, p := range b.Optional {
					if !isOutOfMap(p) {
						return false
					}
				}
			}
			return true
		}() &&
		func() bool {
			for _, b := range m.Map.Barriers {
				tmp := append(b.Optional, b.Required[1:]...)
				if b.Required[0].IsOnSameLine(tmp...) {
					return false
				}
			}
			return true
		}()
}

const perGridSize = 0.01

func init() {
	if perGridSize <= 0 {
		panic(nil)
	}
}

type gridPosition struct{ x, y int }
type vector2D struct{ x, y int }

func angleBetween(v1, v2 vector2D) float64 {
	return math.Acos(float64(v1.x*v2.x+v1.y*v2.y) / (vecLen(v1) * vecLen(v2)))
}
func vecLen(v vector2D) float64 {
	return math.Sqrt(float64(v.x*v.x + v.y*v.y))
}
func isGridPositionInTriangle(p gridPosition, tri [3]gridPosition) bool {
	var xs, ys []int
	for _, trip := range tri {
		xs = append(xs, trip.x)
		ys = append(ys, trip.y)
	}
	sort.Ints(xs)
	sort.Ints(ys)
	if p.x > xs[len(xs)-1] || p.y > ys[len(xs)-1] {
		return false
	}
	var dists = make(map[float64]gridPosition)
	for _, trip := range tri {
		x2 := float64((p.x - trip.x) * (p.x - trip.x))
		y2 := float64((p.y - trip.y) * (p.y - trip.y))
		dists[math.Sqrt(x2+y2)] = trip
	}
	var p1, p2, p3 gridPosition
	var keys []float64
	for k := range dists {
		keys = append(keys, k)
	}
	sort.Float64s(keys)
	if len(keys) == 2 {
		p1, p2 = dists[keys[1]], dists[keys[1]]
	} else {
		p1, p2 = dists[keys[2]], dists[keys[1]]
	}
	p3 = dists[keys[0]]
	p1p := vector2D{p.x - p1.x, p.y - p1.y}
	p2p := vector2D{p.x - p2.x, p.y - p2.y}
	p1p3 := vector2D{p3.x - p1.x, p3.y - p1.y}
	p2p3 := vector2D{p3.x - p2.x, p3.y - p2.y}
	p1p2 := vector2D{p2.x - p1.x, p2.y - p1.y}
	p2p1 := vector2D{0 - p1p2.x, 0 - p1p2.y}

	if ap1p2p1p3, ap1pp1p3 := angleBetween(p1p2, p1p3), angleBetween(p1p, p1p3); ap1pp1p3 > ap1p2p1p3 {
		return false
	}
	if ap2p1p2p3, ap2pp2p3 := angleBetween(p2p1, p2p3), angleBetween(p2p, p2p3); ap2pp2p3 > ap2p1p2p3 {
		return false
	}
	return true
}
func pointToGridPosition(ps ...location.Point2D) (res []gridPosition) {
	for _, p := range ps {
		res = append(res, gridPosition{x: int(p.GetY() / perGridSize), y: int(p.GetY() / perGridSize)})
	}
	return
}
func (m *Map2D) toGrid() *grid.Grid {
	if !m.isValid() {
		return nil
	}
	gridXSize, gridYSize := int(m.Map.Size.X/perGridSize), int(m.Map.Size.Y/perGridSize)
	g := grid.NewGrid(gridXSize, gridYSize)

	//for _, b := range m.Map.Barriers {
	//	//TODO
	//	//ps := pointToGridPosition(append(b.Optional, b.Required[:]...)...)
	//}
	return g
}
