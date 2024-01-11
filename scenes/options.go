package scenes

import (
	"ColouringApp/application"
	"fmt"

	gui "github.com/gen2brain/raylib-go/raygui"
	raylib "github.com/gen2brain/raylib-go/raylib"
)

func Options() {
	var (
		titlePos        float32 = 10
		titleForwardPos float32 = 45

		centerPos  float32 = 10
		backPos    float32 = -application.WindowWidth + 10
		forwardPos float32 = application.WindowWidth + 10

		rootPanel     = true
		controlsPanel = false
		graphicPanel  = false

		rootPos     = centerPos
		controlsPos = forwardPos
		graphicPos  = forwardPos

		titleControls = titleForwardPos
		titleGraphics = titleForwardPos
	)
	// load resources here

	fmt.Println("Options")

	for !application.ShouldQuit {
		application.ShouldQuit = raylib.WindowShouldClose()
		if application.CurrentScene != application.SceneOptions {
			break
		}

		if rootPanel {
			rootPos = raylib.Lerp(rootPos, centerPos, 0.1)
		} else {
			rootPos = raylib.Lerp(rootPos, backPos, 0.1)
		}
		if controlsPanel {
			controlsPos = raylib.Lerp(controlsPos, centerPos, 0.1)
			titleControls = raylib.Lerp(titleControls, titlePos, 0.1)
		} else {
			controlsPos = raylib.Lerp(controlsPos, forwardPos, 0.1)
			titleControls = raylib.Lerp(titleControls, titleForwardPos, 0.1)
		}
		if graphicPanel {
			graphicPos = raylib.Lerp(graphicPos, centerPos, 0.1)
			titleGraphics = raylib.Lerp(titleGraphics, titlePos, 0.1)
		} else {
			graphicPos = raylib.Lerp(graphicPos, forwardPos, 0.1)
			titleGraphics = raylib.Lerp(titleGraphics, titleForwardPos, 0.1)
		}

		raylib.BeginDrawing()
		raylib.ClearBackground(raylib.Black)

		raylib.DrawText("Options", 10, 10, 20, raylib.White)
		raylib.BeginScissorMode(0, 0, application.WindowWidth, 40)
		raylib.DrawText("| Controls", 95, int32(titleControls), 20, raylib.White)
		raylib.DrawText("| Graphics", 95, int32(titleGraphics), 20, raylib.White)
		raylib.EndScissorMode()

		raylib.DrawLine(10, 40, 790, 40, raylib.White)
		if gui.Button(raylib.NewRectangle(application.WindowWidth-110, 10, 100, 20), "Main Menu") {
			application.CurrentScene = application.SceneTitle
		}

		// ROOT PANEL FOR SETTINGS
		if gui.Button(raylib.NewRectangle(rootPos, 50, 100, 20), "Controls") {
			rootPanel = false
			controlsPanel = true
		}
		if gui.Button(raylib.NewRectangle(rootPos, 80, 100, 20), "Graphics") {
			rootPanel = false
			graphicPanel = true
		}

		// CONTROLS PANEL
		raylib.DrawText("Controls", int32(controlsPos), 50, 20, raylib.White)
		if gui.Button(raylib.NewRectangle(controlsPos, 80, 100, 20), "Back") {
			rootPanel = true
			controlsPanel = false
		}

		// GRAPHICS PANEL
		raylib.DrawText("Graphics", int32(graphicPos), 50, 20, raylib.White)
		if gui.Button(raylib.NewRectangle(graphicPos, 80, 100, 20), "Back") {
			rootPanel = true
			graphicPanel = false
		}

		raylib.EndDrawing()
	}

	// unload resources here
}
