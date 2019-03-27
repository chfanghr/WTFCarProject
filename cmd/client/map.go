package client

import (
	"github.com/chfanghr/WTFCarProject/grid"
	"github.com/chfanghr/WTFCarProject/location"
	"math"
)

type mapData struct {
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

func (m *mapData) isValid() bool {
	return m.Map.Size.X > 0 && m.Map.Size.Y > 0 && func() bool {
		for _, b := range m.Map.Barriers {
			tmp := append(b.Optional, b.Required[1:]...)
			if b.Required[0].IsOnSameLine(tmp...) {
				return false
			}
		}
		return true
	}()
}

func (m *mapData) toGrid() *grid.Grid {
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
