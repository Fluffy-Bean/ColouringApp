package scenes

import (
	"ColouringApp/application"

	gui "github.com/gen2brain/raylib-go/raygui"
	raylib "github.com/gen2brain/raylib-go/raylib"
)

func Title() {
	var (
		titleText = application.WindowTitle
		mapImage  = raylib.LoadTexture(application.DirAssets + "Map.png")

		mapX = 0
		mapY = 0
	)

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

		// Map thing?
		raylib.DrawRectangleLines(120, 39, application.WindowWidth-130, application.WindowHeight-49, raylib.White)
		raylib.BeginScissorMode(121, 40, application.WindowWidth-132, application.WindowHeight-51)

		mapX += 1
		mapX = 0
		mapY = 0
		if mapX > 1920 {
			mapX = 0
		}

		mapY += 1
		if mapY > 1080 {
			mapY = 0
		}
		raylib.DrawTexture(mapImage, int32(-mapX), int32(-mapY), raylib.White)

		raylib.EndScissorMode()
		raylib.EndDrawing()
	}

	raylib.UnloadTexture(mapImage)
}
