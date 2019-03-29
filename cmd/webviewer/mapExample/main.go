package main

import (
	"encoding/json"
	"fmt"
	"github.com/chfanghr/WTFCarProject/map2d"
	"github.com/chfanghr/WTFCarProject/rpcprotocal"
	"os"
)

func main() {
	m := map2d.Map2D{}
	m.Map.Size.X = 1000
	m.Map.Size.Y = 1000

	m.Map.Barriers = []map2d.Barrier{
		{
			Required: [3]rpcprotocal.Point2D{
				{
					X: 100,
					Y: 200,
				},
				{
					X: 80,
					Y: 290,
				},
				{
					X: 97,
					Y: 200,
				},
			},
		},
	}

	json.NewEncoder(os.Stdout).Encode(m)
	b, _ := json.Marshal(m)

	m1, err := map2d.NewMap2d(b)
	fmt.Println(m1, err)
}
