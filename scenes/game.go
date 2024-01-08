package scenes

import (
	"ColouringApp/application"

	gui "github.com/gen2brain/raylib-go/raygui"
	raylib "github.com/gen2brain/raylib-go/raylib"
)

func Game() {
	var (
		scenePaused = false
	)

	// load resources here

	for !application.ShouldQuit {
		application.ShouldQuit = raylib.WindowShouldClose()
		if application.CurrentScene != application.SceneGame {
			break
		}

		if raylib.IsKeyPressed(raylib.KeyEscape) {
			scenePaused = !scenePaused
		}

		raylib.BeginDrawing()
		raylib.ClearBackground(raylib.Black)

		raylib.DrawText("Game", 100, 100, 20, raylib.White)

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
