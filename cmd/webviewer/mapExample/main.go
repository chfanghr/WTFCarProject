package main

import (
	"fmt"
	"github.com/chfanghr/WTFCarProject/location"
	"github.com/chfanghr/WTFCarProject/map2d"
	"github.com/chfanghr/WTFCarProject/rpcprotocal"
)

func main() {
	m := map2d.Map2D{}
	m.Map.Size.X = 100
	m.Map.Size.Y = 100

	m.Map.Barriers = []map2d.Barrier{
		{
			Required: [3]rpcprotocal.Point2D{
				{
					X: 10,
					Y: 40,
				},
				{
					X: 80,
					Y: 20,
				},
				{
					X: 97,
					Y: 60,
				},
			},
		},
	}

	res := m.ComputePathTo(*location.NewPoint2D(0, 0), *location.NewPoint2D(99, 88))
	fmt.Println(res)
}
