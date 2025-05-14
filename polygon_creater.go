package main

import rl "github.com/gen2brain/raylib-go/raylib"

type PolygonCreator struct {
	current_polygon Polygon

	dragg_polygon_index int32
	last_dragged_index  int32
	initialized         bool
}

func (pc *PolygonCreator) init() {
	pc.dragg_polygon_index = -1
	pc.current_polygon.debug = true
}

func (pc *PolygonCreator) update(delta_time float64, simulation *Simulation) {
	if !pc.initialized {
		pc.init()
		pc.initialized = true
	}

	var mouse_pos [2]float64 = [2]float64{float64(rl.GetMouseX()), float64(rl.GetMouseY())}

	if rl.IsKeyDown(rl.KeyC) {
		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			pc.current_polygon.addPoint(mouse_pos)
			pc.last_dragged_index = -1
		}
	} else if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
		var mouse_delta [2]float64 = [2]float64{float64(rl.GetMouseDelta().X), float64(rl.GetMouseDelta().Y)}

		if pc.dragg_polygon_index != -1 {
			for index := range len(simulation.polygons[pc.dragg_polygon_index].points) {
				simulation.polygons[pc.dragg_polygon_index].points[index] = addVec2(simulation.polygons[pc.dragg_polygon_index].points[index], mouse_delta)
			}
		} else if pc.current_polygon.IsCollidingWithPoint(mouse_pos) {
			for index := range len(pc.current_polygon.points) {
				pc.current_polygon.points[index] = addVec2(pc.current_polygon.points[index], mouse_delta)
			}

			pc.last_dragged_index = -1
		} else {
			for index := range len(simulation.polygons) {
				var polygon *Polygon = &simulation.polygons[index]

				if polygon.IsCollidingWithPoint(mouse_pos) {
					pc.dragg_polygon_index = int32(index)
					pc.last_dragged_index = int32(index)

					for index := range len(polygon.points) {
						polygon.points[index] = addVec2(polygon.points[index], mouse_delta)
					}

					break
				}
			}
		}
	} else {
		pc.dragg_polygon_index = -1
	}

	if rl.IsKeyDown(rl.KeyEnter) {
		pc.current_polygon.debug = false
		simulation.addPolygon(pc.current_polygon)

		pc.current_polygon = Polygon{debug: true}
	}

	if rl.IsKeyPressed(rl.KeyDelete) {
		if pc.last_dragged_index > -1 {
			if len(simulation.polygons) > 1 {
				simulation.polygons = append(simulation.polygons[:pc.last_dragged_index], simulation.polygons[pc.last_dragged_index+1:]...)
				pc.last_dragged_index = -2
			} else if pc.last_dragged_index == -1 {
				simulation.polygons = []Polygon{}
			}
		} else {
			pc.current_polygon = Polygon{debug: true}
		}
	}
}

func (pc *PolygonCreator) render() {
	pc.current_polygon.render()
}
