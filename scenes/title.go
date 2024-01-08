package scenes

import (
	"ColouringApp/application"

	gui "github.com/gen2brain/raylib-go/raygui"
	raylib "github.com/gen2brain/raylib-go/raylib"
)

func Title() {
	var (
		titleText = "Example Game"
	)

	// load resources here

	for !application.ShouldQuit {
		application.ShouldQuit = raylib.WindowShouldClose()
		if application.CurrentScene != application.SceneTitle {
			break
		}

		raylib.BeginDrawing()
		raylib.ClearBackground(raylib.Black)

		raylib.DrawText(titleText, 10, 10, 20, raylib.White)
		raylib.DrawLine(10, 40, 790, 40, raylib.White)

		if gui.Button(raylib.NewRectangle(10, 50, 100, 20), "Start") {
			application.CurrentScene = application.SceneGame
		}
		if gui.Button(raylib.NewRectangle(10, 80, 100, 20), "Options") {
			application.CurrentScene = application.SceneOptions
		}
		if gui.Button(raylib.NewRectangle(10, 110, 100, 20), "Quit") {
			application.ShouldQuit = true
		}

		raylib.EndDrawing()
	}

	// unload resources here
}
