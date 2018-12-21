package rpcprotocal

import "github.com/chfanghr/Backend/location"

type Point2D struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func Point2DFromLocationPoint2D(d location.Point2D) *Point2D {
	return &Point2D{
		X: d.GetX(),
		Y: d.GetY(),
	}
}

func Point2DToLocationPoint2D(d Point2D) *location.Point2D {
	return location.NewPoint2D(d.X, d.Y)
}
