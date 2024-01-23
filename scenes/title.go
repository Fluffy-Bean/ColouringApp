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

		//mapX = 0
		//mapY = 0
	)

	for !application.ShouldQuit {
		application.ShouldQuit = raylib.WindowShouldClose()
		if application.CurrentScene != application.SceneTitle {
			break
		}

		raylib.BeginDrawing()
		raylib.ClearBackground(raylib.Black)

		if gui.Button(raylib.NewRectangle(10, 10, 40, 40), gui.IconText(gui.ICON_CROSS, "")) {
			application.ShouldQuit = true
		}
		raylib.DrawText(titleText, (application.WindowWidth-raylib.MeasureText(titleText, 20))/2, 20, 20, raylib.White)
		if gui.Button(raylib.NewRectangle(application.WindowWidth-50, 10, 40, 40), gui.IconText(gui.ICON_GEAR, "")) {
			application.CurrentScene = application.SceneOptions
		}

		//raylib.DrawLine(10, 60, 790, 60, raylib.White)

		//// Map thing?
		//raylib.DrawRectangleLines(120, 39, application.WindowWidth-130, application.WindowHeight-49, raylib.White)
		//raylib.BeginScissorMode(121, 40, application.WindowWidth-132, application.WindowHeight-51)
		//raylib.DrawTexture(mapImage, int32(-mapX), int32(-mapY), raylib.White)

		raylib.EndScissorMode()

		if gui.Button(raylib.NewRectangle((application.WindowWidth-100)/2, application.WindowHeight-70, 100, 40), "Start") {
			application.CurrentScene = application.SceneGame
		}

		raylib.EndDrawing()
	}

	raylib.UnloadTexture(mapImage)
}
