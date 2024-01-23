package scenes

import (
	"ColouringApp/application"

	gui "github.com/gen2brain/raylib-go/raygui"
	raylib "github.com/gen2brain/raylib-go/raylib"
)

func Title() {
	var (
		titleText = application.WindowTitle
	)

	for !application.ShouldQuit {
		application.ShouldQuit = raylib.WindowShouldClose()
		if application.CurrentScene != application.SceneTitle {
			break
		}

		// ToDo: Remove this
		application.CurrentScene = application.SceneDrawing

		if raylib.IsWindowResized() {
			application.WindowWidth = int32(raylib.GetScreenWidth())
			application.WindowHeight = int32(raylib.GetScreenHeight())
		}

		raylib.BeginDrawing()
		raylib.ClearBackground(raylib.White)

		if gui.Button(raylib.NewRectangle(10, 10, 40, 40), gui.IconText(gui.ICON_CROSS, "")) {
			application.ShouldQuit = true
		}
		raylib.DrawText(titleText, (application.WindowWidth-raylib.MeasureText(titleText, 20))/2, 20, 20, raylib.Black)
		if gui.Button(raylib.NewRectangle(float32(application.WindowWidth-50), 10, 40, 40), gui.IconText(gui.ICON_GEAR, "")) {
			application.CurrentScene = application.SceneOptions
		}

		raylib.EndScissorMode()

		if gui.Button(raylib.NewRectangle(float32((application.WindowWidth-100)/2), float32(application.WindowHeight-70), 100, 40), "Start") {
			application.CurrentScene = application.SceneDrawing
		}

		raylib.EndDrawing()
	}
}
