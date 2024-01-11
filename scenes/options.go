package scenes

import (
	"ColouringApp/application"
	"fmt"

	gui "github.com/gen2brain/raylib-go/raygui"
	raylib "github.com/gen2brain/raylib-go/raylib"
)

func Options() {
	var (
		centerPos  float32 = 10
		backPos    float32 = -application.WindowWidth + 10
		forwardPos float32 = application.WindowWidth + 10

		rootPanelPos     = true
		controlsPanelPos = false
		graphicPanelPos  = false

		rootPos     = centerPos
		controlsPos = forwardPos
		graphicPos  = forwardPos
	)
	// load resources here

	fmt.Println("Options")

	for !application.ShouldQuit {
		application.ShouldQuit = raylib.WindowShouldClose()
		if application.CurrentScene != application.SceneOptions {
			break
		}

		raylib.BeginDrawing()
		raylib.ClearBackground(raylib.Black)

		raylib.DrawText("Options", 10, 10, 20, raylib.White)
		raylib.DrawLine(10, 40, 790, 40, raylib.White)
		if gui.Button(raylib.NewRectangle(application.WindowWidth-110, 10, 100, 20), "Main Menu") {
			application.CurrentScene = application.SceneTitle
		}

		// ROOT PANEL FOR SETTINGS
		{
			if rootPanelPos {
				rootPos = raylib.Lerp(rootPos, centerPos, 0.1)
			} else {
				rootPos = raylib.Lerp(rootPos, backPos, 0.1)
			}
			if gui.Button(raylib.NewRectangle(rootPos, 50, 100, 20), "Controls") {
				rootPanelPos = false
				controlsPanelPos = true
			}
			if gui.Button(raylib.NewRectangle(rootPos, 80, 100, 20), "Graphics") {
				rootPanelPos = false
				graphicPanelPos = true
			}
		}

		// CONTROLS PANEL
		{
			if controlsPanelPos {
				controlsPos = raylib.Lerp(controlsPos, centerPos, 0.1)
			} else {
				controlsPos = raylib.Lerp(controlsPos, forwardPos, 0.1)
			}

			raylib.DrawText("Controls", int32(controlsPos), 50, 20, raylib.White)
			if gui.Button(raylib.NewRectangle(controlsPos, 80, 100, 20), "Back") {
				rootPanelPos = true
				controlsPanelPos = false
			}
		}

		// GRAPHICS PANEL
		{
			if graphicPanelPos {
				graphicPos = raylib.Lerp(graphicPos, centerPos, 0.1)
			} else {
				graphicPos = raylib.Lerp(graphicPos, forwardPos, 0.1)
			}

			raylib.DrawText("Graphics", int32(graphicPos), 50, 20, raylib.White)
			if gui.Button(raylib.NewRectangle(graphicPos, 80, 100, 20), "Back") {
				rootPanelPos = true
				graphicPanelPos = false
			}
		}

		raylib.EndDrawing()
	}

	// unload resources here
}
