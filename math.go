package main

import "math"

func addVec2(vec1 [2]float64, vec2 [2]float64) [2]float64 {
	vec1[0] += vec2[0]
	vec1[1] += vec2[1]

	return vec1
}

func subVec2(vec1 [2]float64, vec2 [2]float64) [2]float64 {
	vec1[0] -= vec2[0]
	vec1[1] -= vec2[1]

	return vec1
}

func mulVec2(vec1 [2]float64, vec2 [2]float64) [2]float64 {
	vec1[0] *= vec2[0]
	vec1[1] *= vec2[1]

	return vec1
}

func scaleVec2(vec [2]float64, scale float64) [2]float64 {
	vec[0] *= scale
	vec[1] *= scale

	return vec
}

func divVec2(vec1 [2]float64, vec2 [2]float64) [2]float64 {
	if vec2[0] == 0 && vec2[1] == 0 {
		return vec1
	}

	vec1[0] /= vec2[0]
	vec1[1] /= vec2[1]

	return vec1
}

func descaleVec2(vec [2]float64, divisor float64) [2]float64 {
	if divisor == 0 {
		return vec
	}

	vec[0] /= divisor
	vec[1] /= divisor

	return vec
}

func lengthVec2(vec [2]float64) float64 {
	return math.Sqrt(vec[0]*vec[0] + vec[1]*vec[1])
}

func normalizedVec2(vec [2]float64) [2]float64 {
	if vec[0] == 0 && vec[1] == 0 {
		return vec
	}

	var length = lengthVec2(vec)

	return [2]float64{vec[0] / length, vec[1] / length}
}

func dotVec2(vec1 [2]float64, vec2 [2]float64) float64 {
	return vec1[0]*vec2[0] + vec1[1]*vec2[1]
}

func averageVec2(vecs [][2]float64) [2]float64 {
	var x float64 = 0
	var y float64 = 0

	for _, vec := range vecs {
		x += vec[0]
		y += vec[1]
	}

	var count float64 = float64(len(vecs))
	x /= count
	y /= count

	return [2]float64{x, y}
}

func absVec2(vec [2]float64) [2]float64 {
	vec[0] = math.Abs(vec[0])
	vec[1] = math.Abs(vec[1])

	return vec
}

// credits chatGPT
func getClosestPointOnLine(p1, p2, point [2]float64) [2]float64 {
	// Richtungsvektor der Linie
	dx := p2[0] - p1[0]
	dy := p2[1] - p1[1]

	// Vektor von p1 zum externen Punkt
	px := point[0] - p1[0]
	py := point[1] - p1[1]

	// Skalarprodukt und Betragsquadrat des Richtungsvektors
	dot := px*dx + py*dy
	lenSq := dx*dx + dy*dy

	// Parameter t bestimmen (wie weit entlang der Linie)
	t := dot / lenSq

	// Den Projektionspunkt berechnen: p1 + t * (Richtung)
	closest := [2]float64{
		p1[0] + t*dx,
		p1[1] + t*dy,
	}

	return closest
}
