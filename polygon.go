package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func angle(point [2]float64, middle_point [2]float64) float64 {
	return math.Atan2(point[1]-middle_point[1], point[0]-middle_point[0])
}

func quicksort_points(points [][2]float64, middle_point [2]float64) [][2]float64 {
	if len(points) <= 1 {
		return points
	}

	pivot := points[0]
	pivotAngle := angle(pivot, middle_point)

	var less, greater [][2]float64

	for _, p := range points[1:] {
		if angle(p, middle_point) <= pivotAngle {
			less = append(less, p)
		} else {
			greater = append(greater, p)
		}
	}

	less = quicksort_points(less, middle_point)
	greater = quicksort_points(greater, middle_point)

	return append(append(less, pivot), greater...)
}

type Polygon struct {
	points [][2]float64
	debug  bool

	last_closest_point [2]float64
	last_closest_dist  float64
	last_norm_push_vec [2]float64
}

func (p *Polygon) addPoint(point [2]float64) {
	p.points = append(p.points, point)
}

func (p *Polygon) addPoints(point [][2]float64) {
	p.points = append(p.points, point...)
}

func (p *Polygon) render() {
	var point_count = len(p.points)

	var middle_point [2]float64 = averageVec2(p.points)

	var points_sorted [][2]float64 = quicksort_points(p.points, middle_point)

	// fmt.Printf("%v\n", points_sorted)

	for index, point := range points_sorted {
		index += 1

		if index >= point_count {
			index = 0
		}

		var next_point [2]float64 = points_sorted[index]

		rl.DrawLine(int32(point[0]), int32(point[1]), int32(next_point[0]), int32(next_point[1]), rl.White)

		if p.debug {
			rl.DrawCircle(int32(point[0]), int32(point[1]), 3, rl.Gray)
		}
	}

	if p.debug {
		rl.DrawCircle(int32(middle_point[0]), int32(middle_point[1]), 5, rl.Red)
	}
}

// credits to ChanGPT
func (p *Polygon) IsCollidingWithPoint(point [2]float64) bool {
	pts := p.points
	n := len(pts)
	count := 0

	// Setze last_closest_dist auf einen großen Wert
	p.last_closest_dist = math.Inf(1)

	for i := 0; i < n; i++ {
		p1 := pts[i]
		p2 := pts[(i+1)%n]

		closest := closestPointOnLine(p1, p2, point)
		dist := lengthVec2(subVec2(point, closest))

		if dist < p.last_closest_dist {
			p.last_closest_dist = dist
			p.last_closest_point = closest

			var push_vec = subVec2(p2, p1)
			push_vec = [2]float64{-push_vec[1], push_vec[0]}

			p.last_norm_push_vec = normalizedVec2(push_vec)
		}

		// Prüfe, ob die horizontale Linie bei point[1] die Kante schneidet
		if (p1[1] > point[1]) != (p2[1] > point[1]) {
			// x-Koordinate des Schnittpunkts
			xInter := p1[0] + (point[1]-p1[1])*(p2[0]-p1[0])/(p2[1]-p1[1])
			if xInter > point[0] {
				count++
			}
		}
	}

	return count%2 == 1
}

func closestPointOnLine(a, b, P [2]float64) [2]float64 {
	dx := b[0] - a[0]
	dy := b[1] - a[1]
	len2 := dx*dx + dy*dy
	if len2 == 0 {
		return a
	}
	t := ((P[0]-a[0])*dx + (P[1]-a[1])*dy) / len2
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}
	return [2]float64{a[0] + t*dx, a[1] + t*dy}
}
