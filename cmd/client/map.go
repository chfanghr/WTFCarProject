package client

import (
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
