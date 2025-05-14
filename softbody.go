package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type MassPoint struct {
	pos [2]float64
	vel [2]float64

	force  [2]float64
	mass   float64
	radius float64

	index int32
}

func (mp *MassPoint) update(delta_time float64) {
	mp.force = addVec2(mp.force, [2]float64{0, GRAITY * mp.mass})

	mp.vel = addVec2(mp.vel, descaleVec2(scaleVec2(mp.force, delta_time), mp.mass))

	mp.pos = addVec2(mp.pos, scaleVec2(mp.vel, delta_time))

	mp.force = [2]float64{0, 0}
}

func (mp *MassPoint) render() {
	rl.DrawCircle(int32(mp.pos[0]), int32(mp.pos[1]), float32(mp.radius), rl.Red)
}

func (mp *MassPoint) resolveCollision(polygons []Polygon, points []MassPoint) {
	for _, point := range points {
		var dist = lengthVec2(subVec2(point.pos, mp.pos))

		if dist < mp.radius+point.radius {
			var normalized_push_vec = normalizedVec2(subVec2(mp.pos, point.pos))
			var push_vec = scaleVec2(normalized_push_vec, dist/2)

			mp.pos = addVec2(mp.pos, push_vec)

			mp.vel = subVec2(mp.vel, mulVec2(scaleVec2(mulVec2(mp.vel, normalized_push_vec), 2), normalized_push_vec))
		}
	}

	for _, polygon := range polygons {
		if polygon.IsCollidingWithPoint(mp.pos) {
			mp.pos = polygon.last_closest_point

			var normalized_push_vec = polygon.last_norm_push_vec

			mp.vel = subVec2(mp.vel, mulVec2(scaleVec2(mulVec2(mp.vel, normalized_push_vec), 2), normalized_push_vec))
		}
	}
}

func (mp *MassPoint) applyForce(force [2]float64) {
	mp.force = addVec2(mp.force, force)
}

type Spring struct {
	point_a *MassPoint
	point_b *MassPoint

	stiffnes    float64
	rest_length float64
	damping     float64

	index int32
}

func (s *Spring) update(delta_time float64) {
	// Hooke's Law
	spring_vector := subVec2(s.point_a.pos, s.point_b.pos)
	current_length := lengthVec2(spring_vector)
	direction := normalizedVec2(spring_vector)

	displacement := current_length - s.rest_length
	spring_force := -s.stiffnes * displacement

	relative_vel := subVec2(s.point_a.vel, s.point_b.vel)
	damping_force := -s.damping * dotVec2(relative_vel, direction)

	force_total := spring_force + damping_force

	s.point_a.applyForce(scaleVec2(direction, force_total))
	s.point_b.applyForce(scaleVec2(direction, -force_total))
}

func (s *Spring) render() {
	rl.DrawLine(int32(s.point_a.pos[0]), int32(s.point_a.pos[1]), int32(s.point_b.pos[0]), int32(s.point_b.pos[1]), rl.RayWhite)
}

type SoftBody struct {
	points  []MassPoint
	springs []Spring

	stiffnes        float64
	damping         float64
	mass            float64
	point_radius    float64
	physic_steps    int32
	collision_steps int32
	volume_force    float64
}

func (sb *SoftBody) update(delta_time float64, polygons []Polygon) {
	delta_time = delta_time / float64(sb.physic_steps)

	for i := range sb.physic_steps {
		i += 0

		for index := range len(sb.points) {
			var point *MassPoint = &sb.points[index]
			point.update(delta_time)

			point.index = int32(index)
		}

		for index := range len(sb.springs) {
			var spring *Spring = &sb.springs[index]
			spring.update(delta_time)

			spring.index = int32(index)
		}

		var poss [][2]float64 = [][2]float64{}
		for index := range len(sb.points) {
			poss = append(poss, sb.points[index].pos)
		}

		var middle_point = averageVec2(poss)

		for i := range sb.collision_steps {
			i += 0

			for index := range len(sb.points) {
				var point *MassPoint = &sb.points[index]
				point.resolveCollision(polygons, sb.points)

				var force [2]float64 = scaleVec2(normalizedVec2(subVec2(point.pos, middle_point)), sb.volume_force)
				point.applyForce(force)
			}
		}
	}
}

func (sb *SoftBody) render() {
	for _, spring := range sb.springs {
		spring.render()
	}

	for _, point := range sb.points {
		point.render()
	}
}

func (sb *SoftBody) addPoint(pos [2]float64, mass float64, point_radius float64) *int32 {
	var index = int32(len(sb.points))
	sb.points = append(sb.points, MassPoint{pos: pos, mass: mass, index: index, radius: point_radius})

	return &sb.points[index].index
}

func (sb *SoftBody) addSpring(index_p1 int32, index_p2 int32, stiffnes float64, rest_length float64, damping float64) *int32 {
	var index = int32(len(sb.springs))

	var p1, p2 *MassPoint = &sb.points[index_p1], &sb.points[index_p2]
	sb.springs = append(sb.springs, Spring{point_a: p1, point_b: p2, stiffnes: stiffnes, rest_length: rest_length, damping: damping})

	return &sb.springs[index].index
}

func (sb *SoftBody) addPointDefault(pos [2]float64) *int32 {
	return sb.addPoint(pos, sb.mass, sb.point_radius)
}

func (sb *SoftBody) addSpringDefault(index_p1 int32, index_p2 int32, rest_length float64) *int32 {
	return sb.addSpring(index_p1, index_p2, sb.stiffnes, rest_length, sb.damping)
}
