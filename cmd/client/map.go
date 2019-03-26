package client

import (
	"github.com/chfanghr/WTFCarProject/grid"
	"github.com/chfanghr/WTFCarProject/location"
)

//example of a map json message
/*
{
	"map":{
		"size":{
			"x": 100,
			"y": 100
		},
		"barriers": [
			{
			"first":{"x": 1,"y": 1},
			"second":{"x": 50,"y": 32},
			"third":{"x": 77,"y": 12},
			"extra":[{"x": 55,"y": 88}]
			}
		]
	}
}
*/

type Map struct {
	Map struct {
		Size struct {
			//size of map in mile
			X float64 `json:"x"`
			Y float64 `json:"y"`
		} `json:"size"`
		Barriers []struct {
			//barriers represents the barriers in map
			//there should be at least three points to locate an area
			First  location.Point2D   `json:"first"`
			Second location.Point2D   `json:"second"`
			Third  location.Point2D   `json:"third"`
			Extra  []location.Point2D `json:"extra"`
		} `json:"barriers"`
	} `json:"map"`
}

func (m *Map) toGrid() *grid.BlockMap { panic(nil /*TODO:implement this*/) }
