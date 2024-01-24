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

type freeHand struct {
	Color   raylib.Color
	Size    float32
	Opacity float32
	Points  []raylib.Vector2
}

//type line struct {
//	Color raylib.Color
//	Size  float32
//	Start raylib.Vector2
//	End   raylib.Vector2
//}

func Drawing() {
	var (
		camera = raylib.NewCamera2D(raylib.NewVector2(0, 0), raylib.NewVector2(0, 0), 0, 1)
		//cameraMoveOffset = raylib.NewVector2(10, 10)

		canvasSize    = raylib.NewVector2(500, 430)
		canvasScale   = float32(1)
		canvas        = raylib.LoadRenderTexture(int32(canvasSize.X), int32(canvasSize.Y))
		canvasRefresh = true

		sidePanelWidth     = float32(300)
		sidePanelRelativeX = application.WindowWidth - int32(sidePanelWidth)

		drawing = false
		//drawingMode = ModeDrawing

		currentStroke = freeHand{}
		strokes       = []freeHand{currentStroke}
		undoneStrokes = []freeHand{}

		colourPickerVal    = raylib.Orange
		colourPickerHeight = float32(200)

		brushSize    = float32(10)
		brushOpacity = float32(1)

		fileName        = "NewProject"
		fileNameEditing = false
	)

	application.WindowWidth = int32(raylib.GetScreenWidth())
	application.WindowHeight = int32(raylib.GetScreenHeight())

	refreshCanvas := func() {
		raylib.BeginTextureMode(canvas)
		raylib.ClearBackground(raylib.White)

		for i := 1; i < len(strokes); i++ {
			for j := 1; j < len(strokes[i].Points); j++ {
				startPos := raylib.Vector2Subtract(strokes[i].Points[j-1], raylib.NewVector2(10, 10))
				endPos := raylib.Vector2Subtract(strokes[i].Points[j], raylib.NewVector2(10, 10))
				raylib.DrawLineEx(startPos, endPos, strokes[i].Size, raylib.Fade(strokes[i].Color, strokes[i].Opacity))
				raylib.DrawCircle(int32(endPos.X), int32(endPos.Y), strokes[i].Size/2, raylib.Fade(strokes[i].Color, strokes[i].Opacity))
			}
		}

		raylib.EndTextureMode()
	}

	redoStroke := func() {
		if len(undoneStrokes) > 0 {
			strokes = append(strokes, undoneStrokes[len(undoneStrokes)-1])
			undoneStrokes = undoneStrokes[:len(undoneStrokes)-1]

			application.AddToast("Redo")
			canvasRefresh = true
		}
	}
	undoStroke := func() {
		// 1 because I dont know why
		if len(strokes) > 1 {
			undoneStrokes = append(undoneStrokes, strokes[len(strokes)-1])
			strokes = strokes[:len(strokes)-1]

			canvasRefresh = true
			application.AddToast("Undo")
		}
	}

	saveImage := func() {
		if fileName == "" {
			application.AddToast("Please enter a file name")
		} else {
			image := raylib.LoadImageFromTexture(canvas.Texture)

			raylib.ImageRotate(image, 180)
			raylib.ImageFlipHorizontal(image)

			raylib.ExportImage(*image, application.DirUserData+fileName+".png")

			application.AddToast("Drawing saved at " + application.DirUserData + fileName + ".png")
		}
	}

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
			if raylib.GetMouseWheelMove() != 0 && !drawing {
				canvasScale += float32(raylib.GetMouseWheelMove()) * 0.05
			}
			//if raylib.IsMouseButtonPressed(raylib.MouseMiddleButton) {
			//	cameraMoveOffset = raylib.Vector2Subtract(camera.Target, raylib.GetMousePosition())
			//}
			//if raylib.IsMouseButtonDown(raylib.MouseMiddleButton) {
			//	camera.Target = raylib.Vector2Subtract(raylib.GetMousePosition(), raylib.Vector2Scale(cameraMoveOffset, -1))
			//}

			if raylib.IsKeyPressed(raylib.KeyF8) {
				application.AddToast("This is a toast message!")
			}

			if raylib.IsMouseButtonPressed(raylib.MouseLeftButton) {
				if raylib.CheckCollisionPointRec(raylib.GetMousePosition(), raylib.NewRectangle(float32(application.WindowWidth-int32(sidePanelWidth)), 0, sidePanelWidth, float32(application.WindowHeight))) {
					drawing = false
				} else if raylib.CheckCollisionPointRec(raylib.GetMousePosition(), raylib.NewRectangle(10, 10, canvasSize.X, canvasSize.Y)) {
					drawing = true
					currentStroke = freeHand{
						Color:   colourPickerVal,
						Size:    brushSize,
						Opacity: 1,
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
				strokes = append(strokes, currentStroke)
				currentStroke = freeHand{}
				undoneStrokes = []freeHand{}

				drawing = false
				canvasRefresh = true
			}

			if raylib.IsKeyDown(raylib.KeyLeftControl) && raylib.IsKeyDown(raylib.KeyLeftShift) && raylib.IsKeyPressed(raylib.KeyZ) {
				redoStroke()
			} else if raylib.IsKeyDown(raylib.KeyLeftControl) && raylib.IsKeyPressed(raylib.KeyZ) {
				undoStroke()
			} else if raylib.IsKeyDown(raylib.KeyLeftControl) && raylib.IsKeyPressed(raylib.KeyS) {
				saveImage()
			}
		}

		// UPDATE
		{
			if drawing {
				gui.SetState(gui.STATE_DISABLED)
			} else {
				gui.SetState(gui.STATE_NORMAL)
			}

			if canvasRefresh {
				refreshCanvas()
				canvasRefresh = false
			}

			application.UpdateToasts()
		}

		// DRAW
		raylib.BeginDrawing()
		{
			raylib.ClearBackground(raylib.White)
			gui.Grid(raylib.NewRectangle(0, 0, float32(application.WindowWidth), float32(application.WindowHeight)), "", 30, 1, &raylib.Vector2{})

			// Canvas stuff
			raylib.BeginMode2D(camera)
			{
				raylib.DrawRectangle(20, 20, int32(canvasSize.X), int32(canvasSize.Y), raylib.Fade(raylib.Black, 0.3))
				raylib.DrawTexturePro(canvas.Texture, raylib.NewRectangle(0, 0, float32(canvas.Texture.Width), float32(-canvas.Texture.Height)), raylib.NewRectangle(10, 10, canvasSize.X, canvasSize.Y), raylib.Vector2{}, 0, raylib.White)

				//raylib.BeginScissorMode(10, 10, int32(canvasSize.X), int32(canvasSize.Y))
				for i := 1; i < len(currentStroke.Points); i++ {
					raylib.DrawLineEx(currentStroke.Points[i-1], currentStroke.Points[i], currentStroke.Size, raylib.Fade(currentStroke.Color, currentStroke.Opacity))
					raylib.DrawCircle(int32(currentStroke.Points[i].X), int32(currentStroke.Points[i].Y), currentStroke.Size/2, raylib.Fade(currentStroke.Color, currentStroke.Opacity))
				}
				//raylib.EndScissorMode()

				if drawing {
					raylib.DrawRectangleLines(10, 10, int32(canvasSize.X), int32(canvasSize.Y), raylib.DarkGray)
					raylib.DrawCircleLines(int32(raylib.GetMousePosition().X), int32(raylib.GetMousePosition().Y), brushSize/2, raylib.Black)
				} else {
					raylib.DrawRectangleLines(10, 10, int32(canvasSize.X), int32(canvasSize.Y), raylib.Gray)
					raylib.DrawCircleLines(int32(raylib.GetMousePosition().X), int32(raylib.GetMousePosition().Y), brushSize/2, raylib.Black)
				}
			}
			raylib.EndMode2D()

			// UI stuff
			raylib.BeginScissorMode(sidePanelRelativeX, 0, int32(sidePanelWidth), application.WindowHeight)
			{
				raylib.DrawRectangle(sidePanelRelativeX, 0, int32(sidePanelWidth), application.WindowHeight, raylib.Fade(raylib.White, 0.7))

				if gui.Button(raylib.NewRectangle(float32(sidePanelRelativeX+10), 10, 25, 25), gui.IconText(gui.ICON_CROSS, "")) {
					application.CurrentScene = application.SceneTitle
				}
				if gui.Button(raylib.NewRectangle(float32(sidePanelRelativeX+20+25), 10, 25, 25), gui.IconText(gui.ICON_FOLDER_SAVE, "")) {
					saveImage()
				}

				if gui.Button(raylib.NewRectangle(float32(application.WindowWidth-70), 10, 25, 25), gui.IconText(gui.ICON_UNDO, "")) {
					undoStroke()
				}
				if gui.Button(raylib.NewRectangle(float32(application.WindowWidth-35), 10, 25, 25), gui.IconText(gui.ICON_REDO, "")) {
					redoStroke()
				}

				colourPickerVal = gui.ColorPicker(raylib.NewRectangle(float32(sidePanelRelativeX+10), 45, sidePanelWidth-45, colourPickerHeight), "Color", colourPickerVal)

				gui.Label(raylib.NewRectangle(float32(sidePanelRelativeX+10), 55+colourPickerHeight, 60, 20), "Brush Size")
				brushSize = gui.Slider(raylib.NewRectangle(float32(sidePanelRelativeX+80), 55+colourPickerHeight, sidePanelWidth-90, 20), "", "", brushSize, 1, 100)

				gui.Label(raylib.NewRectangle(float32(sidePanelRelativeX+10), 85+colourPickerHeight, 60, 20), "Brush Opacity")
				brushOpacity = gui.Slider(raylib.NewRectangle(float32(sidePanelRelativeX+80), 85+colourPickerHeight, sidePanelWidth-90, 20), "", "", brushOpacity, 0, 1)

				gui.Label(raylib.NewRectangle(float32(sidePanelRelativeX+10), 115+colourPickerHeight, 60, 20), "File Name")
				if gui.TextBox(raylib.NewRectangle(float32(sidePanelRelativeX+80), 115+colourPickerHeight, sidePanelWidth-90, 20), &fileName, 40, fileNameEditing) {
					fileNameEditing = !fileNameEditing
				}
			}
			raylib.EndScissorMode()
			raylib.DrawRectangleLines(sidePanelRelativeX, 0, int32(sidePanelWidth), application.WindowHeight, raylib.Gray)

			// Info
			{
				var text string

				text = fmt.Sprintf("Strokes: %d | Points: %d", len(strokes), len(currentStroke.Points))
				gui.StatusBar(raylib.NewRectangle(0, float32(application.WindowHeight-20), 200, 20), text)

				text = fmt.Sprintf("Canvas Size: %dx%d | Scale: %v", int32(canvasSize.X), int32(canvasSize.Y), canvasScale)
				gui.StatusBar(raylib.NewRectangle(199, float32(application.WindowHeight-20), 200, 20), text)
			}

			// Draw toasts
			application.DrawToasts()
		}
		raylib.EndDrawing()
	}

	// unload resources here
}
