package scenes

import (
	"ColouringApp/application"
	"fmt"
	gui "github.com/gen2brain/raylib-go/raygui"
	raylib "github.com/gen2brain/raylib-go/raylib"
)

const (
	ModeDrawing = iota
	ModeLine
)

type stroke struct {
	Color  raylib.Color
	Size   float32
	Points []raylib.Vector2
}

func Drawing() {
	var (
		canvasSize = raylib.NewVector2(500, 430)
		canvas     = raylib.LoadRenderTexture(int32(canvasSize.X), int32(canvasSize.Y))

		sidePanelWidth     float32 = 300
		sidePanelRelativeX         = application.WindowWidth - int32(sidePanelWidth)

		drawing = false
		//drawingMode = ModeDrawing

		currentStroke = stroke{}
		strokes       = []stroke{currentStroke}
		undoneStrokes = []stroke{}

		colourPickerVal            = raylib.Orange
		colourPickerHeight float32 = 200

		brushSize float32 = 10
	)

	refreshCanvas := func() {
		raylib.BeginTextureMode(canvas)
		raylib.ClearBackground(raylib.White)

		for i := 1; i < len(strokes); i++ {
			for j := 1; j < len(strokes[i].Points); j++ {
				startPos := raylib.Vector2Subtract(strokes[i].Points[j-1], raylib.NewVector2(10, 10))
				endPos := raylib.Vector2Subtract(strokes[i].Points[j], raylib.NewVector2(10, 10))
				raylib.DrawLineEx(startPos, endPos, strokes[i].Size, strokes[i].Color)
				raylib.DrawCircle(int32(endPos.X), int32(endPos.Y), strokes[i].Size/2, strokes[i].Color)
			}
		}

		raylib.EndTextureMode()
	}

	// Create canvas
	refreshCanvas()

	for !application.ShouldQuit {
		// DEFAULT
		{
			application.ShouldQuit = raylib.WindowShouldClose()
			if application.CurrentScene != application.SceneDrawing {
				break
			}
			if raylib.IsWindowResized() {
				application.WindowWidth = int32(raylib.GetScreenWidth())
				application.WindowHeight = int32(raylib.GetScreenHeight())

				sidePanelRelativeX = application.WindowWidth - int32(sidePanelWidth)
			}
		}

		// INPUT
		{
			if raylib.IsMouseButtonPressed(raylib.MouseLeftButton) {
				drawing = false

				if raylib.CheckCollisionPointRec(raylib.GetMousePosition(), raylib.NewRectangle(10, 10, canvasSize.X, canvasSize.Y)) {
					drawing = true
					currentStroke = stroke{
						Color: colourPickerVal,
						Size:  brushSize,
					}
				}
			}
			if raylib.IsMouseButtonDown(raylib.MouseLeftButton) && drawing {
				var safeZone float32 = 5

				if len(currentStroke.Points) <= 1 {
					currentStroke.Points = append(currentStroke.Points, raylib.GetMousePosition())
				} else if raylib.Vector2Distance(currentStroke.Points[len(currentStroke.Points)-1], raylib.GetMousePosition()) > safeZone {
					currentStroke.Points = append(currentStroke.Points, raylib.GetMousePosition())
				}
			}
			if raylib.IsMouseButtonReleased(raylib.MouseLeftButton) && currentStroke.Points != nil {
				drawing = false

				strokes = append(strokes, currentStroke)
				currentStroke = stroke{}
				undoneStrokes = []stroke{}

				refreshCanvas()
			}

			if raylib.IsKeyDown(raylib.KeyLeftControl) && raylib.IsKeyDown(raylib.KeyLeftShift) && raylib.IsKeyPressed(raylib.KeyZ) {
				if len(undoneStrokes) > 1 {
					strokes = append(strokes, undoneStrokes[len(undoneStrokes)-1])
					undoneStrokes = undoneStrokes[:len(undoneStrokes)-1]
				}

				refreshCanvas()
			} else if raylib.IsKeyDown(raylib.KeyLeftControl) && raylib.IsKeyPressed(raylib.KeyZ) {
				if len(strokes) > 1 {
					undoneStrokes = append(undoneStrokes, strokes[len(strokes)-1])
					strokes = strokes[:len(strokes)-1]
				}

				refreshCanvas()
			}
		}

		// UPDATE
		{
			if drawing {
				gui.SetState(gui.STATE_DISABLED)
			} else {
				gui.SetState(gui.STATE_NORMAL)
			}
		}

		// DRAW
		raylib.BeginDrawing()
		{
			raylib.ClearBackground(raylib.White)
			gui.Grid(raylib.NewRectangle(0, 0, float32(application.WindowWidth), float32(application.WindowHeight)), "", 30, 1, &raylib.Vector2{})

			// Canvas stuff
			raylib.DrawRectangle(20, 20, int32(canvasSize.X), int32(canvasSize.Y), raylib.Fade(raylib.Black, 0.3))
			raylib.BeginScissorMode(10, 10, int32(canvasSize.X), int32(canvasSize.Y))
			{
				raylib.DrawTexturePro(canvas.Texture, raylib.NewRectangle(0, 0, float32(canvas.Texture.Width), float32(-canvas.Texture.Height)), raylib.NewRectangle(10, 10, canvasSize.X, canvasSize.Y), raylib.Vector2{}, 0, raylib.White)
				for i := 1; i < len(currentStroke.Points); i++ {
					raylib.DrawLineEx(currentStroke.Points[i-1], currentStroke.Points[i], currentStroke.Size, currentStroke.Color)
					raylib.DrawCircle(int32(currentStroke.Points[i].X), int32(currentStroke.Points[i].Y), currentStroke.Size/2, currentStroke.Color)
				}
				raylib.DrawCircleLines(int32(raylib.GetMousePosition().X), int32(raylib.GetMousePosition().Y), brushSize/2, raylib.Black)
			}
			raylib.EndScissorMode()
			raylib.DrawRectangleLines(10, 10, int32(canvasSize.X), int32(canvasSize.Y), raylib.Gray)

			// UI stuff
			raylib.BeginScissorMode(sidePanelRelativeX, 0, int32(sidePanelWidth), application.WindowHeight)
			{
				raylib.DrawRectangle(sidePanelRelativeX, 0, int32(sidePanelWidth), application.WindowHeight, raylib.Fade(raylib.White, 0.7))

				if gui.Button(raylib.NewRectangle(float32(sidePanelRelativeX+10), 10, 25, 25), gui.IconText(gui.ICON_CROSS, "")) {
					application.CurrentScene = application.SceneTitle
				}
				if gui.Button(raylib.NewRectangle(float32(sidePanelRelativeX+20+25), 10, 25, 25), gui.IconText(gui.ICON_FOLDER_SAVE, "")) {
				}

				colourPickerVal = gui.ColorPicker(raylib.NewRectangle(float32(sidePanelRelativeX+10), 45, sidePanelWidth-40, colourPickerHeight), "Color", colourPickerVal)
				brushSize = gui.Slider(raylib.NewRectangle(float32(sidePanelRelativeX+80), 65+colourPickerHeight, 200, 20), "Brush Size", "", brushSize, 1, 100)
			}
			raylib.EndScissorMode()
			raylib.DrawRectangleLines(sidePanelRelativeX, 0, int32(sidePanelWidth), application.WindowHeight, raylib.Gray)

			// Info
			{
				var text string

				text = fmt.Sprintf("Strokes: %d | Stroke Points: %d", len(strokes), len(currentStroke.Points))
				gui.StatusBar(raylib.NewRectangle(0, float32(application.WindowHeight-20), 200, 20), text)

				text = fmt.Sprintf("Canvas Size: %dx%d", int32(canvasSize.X), int32(canvasSize.Y))
				gui.StatusBar(raylib.NewRectangle(199, float32(application.WindowHeight-20), 200, 20), text)
			}
		}
		raylib.EndDrawing()
	}

	// unload resources here
}
