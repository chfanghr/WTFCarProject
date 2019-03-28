package map2d

import (
	"github.com/chfanghr/WTFCarProject/location"
	"math"
)

func euqal(x float64, y float64) bool {
	v := x - y
	const delta float64 = 1e-6
	if v < delta && v > -delta {
		return true
	}
	return false

}
func little(x float64, y float64) bool {
	if euqal(x, y) {
		return false
	}
	return x < y
}
func littleEqual(x float64, y float64) bool {
	if euqal(x, y) {
		return true
	}
	return x < y
}
func InPolygon(point []float64, vertices [][]float64) bool {
	x := point[0]
	y := point[1]
	sz := len(vertices)
	is_in := false

	for i := 0; i < sz; i++ {
		j := i - 1
		if i == 0 {
			j = sz - 1
		}
		vi := vertices[i]
		vj := vertices[j]

		xmin := vi[0]
		xmax := vj[0]
		if xmin > xmax {
			t := xmin
			xmin = xmax
			xmax = t
		}
		ymin := vi[1]
		ymax := vj[1]
		if ymin > ymax {
			t := ymin
			ymin = ymax
			ymax = t
		}
		// i//j//aixs_x
		if euqal(vj[1], vi[1]) {
			if euqal(y, vi[1]) && littleEqual(xmin, x) && littleEqual(x, xmax) {
				return true
			}
			continue
		}

		xt := (vj[0]-vi[0])*(y-vi[1])/(vj[1]-vi[1]) + vi[0]
		if euqal(xt, x) && littleEqual(ymin, y) && littleEqual(y, ymax) {
			// on edge [vj,vi]
			return true
		}
		if little(x, xt) && littleEqual(ymin, y) && little(y, ymax) {
			is_in = !is_in
		}

	}
	return is_in
}

type grid2D struct{ x, y int }

func pointToGrid2D(p location.Point2D) (res grid2D) {
	res.x = int(math.Round(p.GetX()))
	res.y = int(math.Round(p.GetY()))
	return
}
func pointsToGrid2DArray(ps ...location.Point2D) (res []grid2D) {
	for _, p := range ps {
		res = append(res, pointToGrid2D(p))
	}
	return
}
func grid2DToFloatArray(g grid2D) (res []float64) {
	res = append(res, float64(g.x), float64(g.y))
	return
}
func grid2DsToFloat2DArray(gs ...grid2D) (res [][]float64) {
	for _, g := range gs {
		res = append(res, grid2DToFloatArray(g))
	}
	return
}
