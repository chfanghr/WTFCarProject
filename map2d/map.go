package map2d

import (
	"encoding/json"
	"github.com/chfanghr/WTFCarProject/grid"
	"github.com/chfanghr/WTFCarProject/location"
	"math"
)

type Map2d struct {
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

func NewMap2d(raw []byte) (*Map2d, error) {
	res := &Map2d{}
	err := json.Unmarshal(raw, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (m *Map2d) isValid() bool {
	return m.Map.Size.X > 0 && m.Map.Size.Y > 0 && func() bool {
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
	}() && func() bool {
		for _, b := range m.Map.Barriers {
			tmp := append(b.Optional, b.Required[1:]...)
			if b.Required[0].IsOnSameLine(tmp...) {
				return false
			}
		}
		return true
	}()
}

func (m *Map2d) toGrid() *grid.Grid {
	if m.Map.Size.X+1 > math.MaxInt32 || m.Map.Size.Y+1 > math.MaxInt32 {
		return nil
	}
	gridX, gridY := int(m.Map.Size.X+1), int(m.Map.Size.Y+1)
	g := grid.NewGrid(gridX, gridY)
	if g == nil {
		return nil
	}
	panic(nil) //TODO
}
