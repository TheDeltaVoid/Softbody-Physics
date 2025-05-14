package main

type Simulation struct {
	polygons    []Polygon
	soft_bodies []SoftBody

	physic_steps int32
}

func (s *Simulation) update(delta_time float64) {
	delta_time = delta_time / float64(s.physic_steps)

	for i := range s.physic_steps {
		i += 0

		for _, soft_body := range s.soft_bodies {
			soft_body.update(delta_time, s.polygons)
		}
	}
}

func (s *Simulation) addPolygon(polygon Polygon) {
	s.polygons = append(s.polygons, polygon)
}

func (s *Simulation) addSoftBody(soft_body SoftBody) {
	s.soft_bodies = append(s.soft_bodies, soft_body)
}

func (s *Simulation) render() {
	for _, polygon := range s.polygons {
		polygon.render()
	}

	for _, soft_body := range s.soft_bodies {
		soft_body.render()
	}
}

func (s *Simulation) addRectSoftBody(pos [2]float64, width int32, height int32, spacing float64, stiffnes float64, damping float64, mass float64, point_radius float64) {
	var diagonal_spacing = lengthVec2([2]float64{spacing, spacing})

	var soft_body SoftBody = SoftBody{stiffnes: stiffnes, damping: damping, mass: mass, point_radius: point_radius}
	var point_indicies [][]*int32 = make([][]*int32, width)

	for y := range width {
		point_indicies[y] = make([]*int32, height)
	}

	for y := range height {
		for x := range width {
			// fmt.Printf("X: %v, Y: %v\n", x, y)
			point_indicies[x][y] = soft_body.addPointDefault([2]float64{float64(x)*spacing + pos[0], float64(y)*spacing + pos[1]})
		}
	}

	for y := range height {
		for x := range width {
			if x < width-1 {
				soft_body.addSpringDefault(*point_indicies[x][y], *point_indicies[x+1][y], spacing)
			}
			if y < height-1 {
				soft_body.addSpringDefault(*point_indicies[x][y], *point_indicies[x][y+1], spacing)
			}

			if x < width-1 && y < height-1 {
				soft_body.addSpringDefault(*point_indicies[x][y], *point_indicies[x+1][y+1], diagonal_spacing)
			}
			if x < width-1 && y > 0 {
				soft_body.addSpringDefault(*point_indicies[x][y], *point_indicies[x+1][y-1], diagonal_spacing)
			}
		}
	}

	soft_body.collision_steps = 1
	soft_body.physic_steps = 5

	soft_body.volume_force = 1000

	s.addSoftBody(soft_body)
}

func (s *Simulation) addCircleSoftBody(pos [2]float64, radius float64, point_count int32, stiffnes float64, damping float64, mass float64, point_radius float64) {
	var soft_body SoftBody = SoftBody{stiffnes: stiffnes, damping: damping, mass: mass, point_radius: point_radius}
	var point_indicies []*int32 = []*int32{}

	var base_point [2]float64 = [2]float64{0, radius}
	for i := range point_count {
		var rotation float64 = 360 / float64(point_count) * float64(i)

		point_indicies = append(point_indicies, soft_body.addPointDefault(addVec2(rotateVec2(base_point, rotation), pos)))
	}

	var rest_length float64 = lengthVec2(subVec2(soft_body.points[1].pos, soft_body.points[0].pos))

	for index := range point_indicies {
		var next_index = index + 1
		if index+1 > len(point_indicies)-1 {
			next_index = 0
		}

		soft_body.addSpringDefault(int32(index), int32(next_index), rest_length)
	}

	soft_body.collision_steps = 1
	soft_body.physic_steps = 5

	soft_body.volume_force = 50000 / (float64(point_count) * 4)

	s.addSoftBody(soft_body)
}
