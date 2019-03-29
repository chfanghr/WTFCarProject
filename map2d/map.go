package map2d

import (
	"encoding/json"
	"github.com/chfanghr/WTFCarProject/grid"
	"github.com/chfanghr/WTFCarProject/location"
	"github.com/chfanghr/WTFCarProject/rpcprotocal"
)

const perGridSize = 0.01

type Barrier struct {
	Required [3]rpcprotocal.Point2D `json:"required"`
	Optional []rpcprotocal.Point2D  `json:"optional"`
}
type Map2D struct {
	Map struct {
		Size struct {
			X float64 `json:"x"`
			Y float64 `json:"y"`
		} `json:"size"`
		Barriers []Barrier `json:"barriers"`
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

//func (m *Map2D) isValid() bool {
//	return m.Map.Size.X > 0 && m.Map.Size.Y > 0 &&
//		func() bool {
//			for _, b := range m.Map.Barriers {
//				isOutOfMap := func(p location.Point2D) bool {
//					return !(p.GetX() > m.Map.Size.X || p.GetY() > m.Map.Size.Y || p.GetX() < 0 || p.GetY() < 0)
//				}
//				for _, p := range b.Required {
//					if !isOutOfMap(*rpcprotocal.Point2DToLocationPoint2D(p)) {
//						return false
//					}
//				}
//				for _, p := range b.Optional {
//					if !isOutOfMap(*rpcprotocal.Point2DToLocationPoint2D(p)) {
//						return false
//					}
//				}
//			}
//			return true
//		}() &&
//		func() bool {
//			for _, b := range m.Map.Barriers {
//				var ps []location.Point2D
//				for _, p := range b.Required {
//					ps = append(ps, *rpcprotocal.Point2DToLocationPoint2D(p))
//				}
//				for _, p := range b.Optional {
//					ps = append(ps, *rpcprotocal.Point2DToLocationPoint2D(p))
//				}
//				if ps[0].IsOnSameLine(ps[1:]...) {
//					return false
//				}
//			}
//			return true
//		}()
//}

func (m *Map2D) toGrid() *grid.Grid {
	//if !m.isValid() {
	//	return nil
	//}
	gridXSize, gridYSize := int(m.Map.Size.X/perGridSize), int(m.Map.Size.Y/perGridSize)
	g := grid.NewGrid(gridXSize, gridYSize)
	var barriers [][][]float64
	for _, b := range m.Map.Barriers {
		var ps []location.Point2D
		for _, p := range b.Required {
			ps = append(ps, *rpcprotocal.Point2DToLocationPoint2D(p))
		}
		for _, p := range b.Optional {
			ps = append(ps, *rpcprotocal.Point2DToLocationPoint2D(p))
		}
		barrier := grid2DsToFloat2DArray(pointsToGrid2DArray(ps...)...)
		barriers = append(barriers, barrier)
	}
	for x := 0; x < gridXSize; x++ {
		for y := 0; y < gridYSize; y++ {
			for _, b := range barriers {
				if InPolygon([]float64{float64(x), float64(y)}, b) {
					g.SetBarrier(x, y, true)
				}
			}
		}
	}
	return g
}

func (m *Map2D) ComputePathTo(f, t location.Point2D) (ress []location.Point2D) {
	//if !m.isValid() {
	//	return nil
	//}
	fg, tg := pointToGrid2D(f), pointToGrid2D(t)
	res := m.toGrid().GetShortestPath(fg.x, fg.y, tg.x, tg.y)
	for _, v := range res {
		ress = append(ress, *location.NewPoint2D(float64(v.X)*perGridSize, float64(v.Y)*perGridSize))
	}
	return
}
