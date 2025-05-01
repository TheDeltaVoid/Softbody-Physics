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
	var point_inicies [][]*int32 = make([][]*int32, width)

	for y := range width {
		point_inicies[y] = make([]*int32, height)
	}

	for y := range height {
		for x := range width {
			// fmt.Printf("X: %v, Y: %v\n", x, y)
			point_inicies[x][y] = soft_body.addPointDefault([2]float64{float64(x)*spacing + pos[0], float64(y)*spacing + pos[1]})
		}
	}

	for y := range height {
		for x := range width {
			if x < width-1 {
				soft_body.addSpringDefault(*point_inicies[x][y], *point_inicies[x+1][y], spacing)
			}
			if y < height-1 {
				soft_body.addSpringDefault(*point_inicies[x][y], *point_inicies[x][y+1], spacing)
			}

			if x < width-1 && y < height-1 {
				soft_body.addSpringDefault(*point_inicies[x][y], *point_inicies[x+1][y+1], diagonal_spacing)
			}
			if x < width-1 && y > 0 {
				soft_body.addSpringDefault(*point_inicies[x][y], *point_inicies[x+1][y-1], diagonal_spacing)
			}
		}
	}

	soft_body.collision_steps = 2
	soft_body.physic_steps = 10

	s.addSoftBody(soft_body)
}
