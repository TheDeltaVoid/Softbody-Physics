package main

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	fmt.Printf("%v\n", math.Sin(6.282))

	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.InitWindow(1080, 720, "Softbody Simulation")
	rl.SetTargetFPS(FPS)

	var sim Simulation = Simulation{physic_steps: 1}

	var poly_creator = PolygonCreator{}

	sim.addRectSoftBody([2]float64{100, 100}, 10, 10, 20, 200, 5, 0.1, 5)
	//sim.addCircleSoftBody([2]float64{200, 200}, 100, 200, 200, 5, 0.1, 5)

	var paused bool = false
	var delta_time float64
	for !rl.WindowShouldClose() {
		delta_time = float64(rl.GetFrameTime())

		if rl.IsKeyPressed(rl.KeyP) {
			paused = !paused
		}

		if !paused {
			sim.update(delta_time)
		}

		poly_creator.update(delta_time, &sim)

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		sim.render()
		poly_creator.render()

		rl.DrawFPS(5, 5)
		rl.EndDrawing()
	}

	rl.CloseWindow()
}
