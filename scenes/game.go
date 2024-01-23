package scenes

import (
	"ColouringApp/application"

	gui "github.com/gen2brain/raylib-go/raygui"
	raylib "github.com/gen2brain/raylib-go/raylib"
)

func Game() {
	var (
		scenePaused = false

		gridPos             = raylib.NewVector2(0, 0)
		gridSize            = raylib.NewVector2(100, (application.WindowHeight-20)/5)
		gridDensity float32 = 5

		brushSize float32 = 1
		color             = raylib.Orange
		canvas    [][]raylib.Color
	)

	// Create canvas
	for x := 0; x < int(gridSize.X); x++ {
		canvas = append(canvas, []raylib.Color{})
		for y := 0; y < int(gridSize.Y); y++ {
			canvas[x] = append(canvas[x], raylib.White)
		}
	}
	if gridPos.X >= 0 && gridPos.X < gridSize.X && gridPos.Y >= 0 && gridPos.Y < gridSize.Y {
		for x := int(gridPos.X - brushSize/2); x < int(gridPos.X+brushSize/2); x++ {
			for y := int(gridPos.Y - brushSize/2); y < int(gridPos.Y+brushSize/2); y++ {
				if x >= 0 && x < int(gridSize.X) && y >= 0 && y < int(gridSize.Y) {
					canvas[x][y] = color
				}
			}
		}
	}

	// load resources here

	for !application.ShouldQuit {
		application.ShouldQuit = raylib.WindowShouldClose()
		if application.CurrentScene != application.SceneGame {
			break
		}
		if raylib.IsKeyPressed(raylib.KeyEscape) {
			scenePaused = !scenePaused
		}

		// INPUT
		if raylib.IsMouseButtonDown(raylib.MouseLeftButton) {
			if gridPos.X >= 0 && gridPos.X < gridSize.X && gridPos.Y >= 0 && gridPos.Y < gridSize.Y {
				for x := int(gridPos.X - brushSize/2); x < int(gridPos.X+brushSize/2); x++ {
					for y := int(gridPos.Y - brushSize/2); y < int(gridPos.Y+brushSize/2); y++ {
						if x >= 0 && x < int(gridSize.X) && y >= 0 && y < int(gridSize.Y) {
							canvas[x][y] = color
						}
					}
				}
			}
		}

		// UPDATE

		// DRAW
		raylib.BeginDrawing()
		raylib.ClearBackground(raylib.Black)

		gui.Grid(raylib.NewRectangle(10, 10, gridDensity*gridSize.X, gridDensity*gridSize.Y), "", gridDensity, 1, &gridPos)
		// Default grid doesnt show up
		raylib.DrawRectangle(0, 0, application.WindowWidth, application.WindowHeight, raylib.LightGray)

		for x := 0; x < int(gridSize.X); x++ {
			for y := 0; y < int(gridSize.Y); y++ {
				pos := raylib.NewVector2(float32(x)*gridDensity, float32(y)*gridDensity)
				raylib.DrawRectangle(int32(pos.X)+10, int32(pos.Y)+10, int32(gridDensity), int32(gridDensity), canvas[x][y])
			}
		}
		raylib.DrawRectangleLines(10, 10, int32(gridDensity*gridSize.X), int32(gridDensity*gridSize.Y), raylib.Gray)

		if gridPos.X >= 0 && gridPos.X < gridSize.X && gridPos.Y >= 0 && gridPos.Y < gridSize.Y {
			for x := int(gridPos.X - brushSize/2); x < int(gridPos.X+brushSize/2); x++ {
				for y := int(gridPos.Y - brushSize/2); y < int(gridPos.Y+brushSize/2); y++ {
					if x >= 0 && x < int(gridSize.X) && y >= 0 && y < int(gridSize.Y) {
						pos := raylib.NewVector2(float32(x)*gridDensity, float32(y)*gridDensity)
						raylib.DrawRectangle(int32(pos.X)+10, int32(pos.Y)+10, int32(gridDensity), int32(gridDensity), raylib.Fade(color, 0.5))
					}
				}
			}
		}

		color = gui.ColorPicker(raylib.NewRectangle(float32(30+int32(gridDensity*gridSize.X)), 10, application.WindowWidth-float32(65+int32(gridDensity*gridSize.X)), 200), "Color", color)

		brushSize = gui.Slider(raylib.NewRectangle(float32(90+int32(gridDensity*gridSize.X)), 30+200, 200, 20), "Brush Size", "", brushSize, 1, 10)

		raylib.DrawLine(20+int32(gridDensity*gridSize.X), 10, 20+int32(gridDensity*gridSize.X), 10+int32(gridDensity*gridSize.Y), raylib.Gray)
		raylib.DrawLine(30+int32(gridDensity*gridSize.X), 20+200, application.WindowWidth-10, 20+200, raylib.Gray)

		if scenePaused {
			raylib.DrawRectangle(0, 0, application.WindowWidth, application.WindowHeight, raylib.Fade(raylib.Black, 0.5))
			raylib.DrawText("Paused", 10, 10, 20, raylib.White)
			raylib.DrawLine(10, 40, 790, 40, raylib.White)
			if gui.Button(raylib.NewRectangle(application.WindowWidth-110, 10, 100, 20), "Unpause") {
				scenePaused = false
			}

			if gui.Button(raylib.NewRectangle(10, 50, 100, 20), "Main Menu") {
				application.CurrentScene = application.SceneTitle
			}
		}

		raylib.EndDrawing()
	}

	// unload resources here
}
