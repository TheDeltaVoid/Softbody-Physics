package main

import rl "github.com/gen2brain/raylib-go/raylib"

type PolygonCreator struct {
	current_polygon Polygon
}

func (pc *PolygonCreator) update(delta_time float64, simulation *Simulation) {
	if !pc.current_polygon.debug {
		pc.current_polygon.debug = true
	}

	var mouse_pos [2]float64 = [2]float64{float64(rl.GetMouseX()), float64(rl.GetMouseY())}

	if rl.IsKeyDown(rl.KeyC) {
		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			pc.current_polygon.addPoint(mouse_pos)
		}
	} else if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
		var mouse_delta [2]float64 = [2]float64{float64(rl.GetMouseDelta().X), float64(rl.GetMouseDelta().Y)}

		if pc.current_polygon.IsCollidingWithPoint(mouse_pos) {
			for index := range len(pc.current_polygon.points) {
				pc.current_polygon.points[index] = addVec2(pc.current_polygon.points[index], mouse_delta)
			}
		} else {
			for index := range len(simulation.polygons) {
				var polygon *Polygon = &simulation.polygons[index]

				if polygon.IsCollidingWithPoint(mouse_pos) {
					for index := range len(polygon.points) {
						polygon.points[index] = addVec2(polygon.points[index], mouse_delta)
					}

					break
				}
			}
		}
	}

	if rl.IsKeyDown(rl.KeyEnter) {
		pc.current_polygon.debug = false
		simulation.addPolygon(pc.current_polygon)

		pc.current_polygon = Polygon{debug: true}
	}
}

func (pc *PolygonCreator) render() {
	pc.current_polygon.render()
}
