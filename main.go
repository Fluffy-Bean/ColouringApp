package main

import (
	"fmt"
	gui "github.com/gen2brain/raylib-go/raygui"
	raylib "github.com/gen2brain/raylib-go/raylib"
	"os"
)

func main() {
	raylib.SetConfigFlags(raylib.FlagWindowResizable)
	raylib.SetTraceLogLevel(raylib.LogTrace)
	raylib.SetConfigFlags(raylib.FlagMsaa4xHint)

	raylib.InitWindow(WindowWidth, WindowHeight, WindowTitle)
	raylib.InitAudioDevice()

	raylib.SetTargetFPS(WindowFPS)
	//raylib.SetExitKey(0) // disable exit key

	var (
		camera = raylib.NewCamera2D(raylib.NewVector2(0, 0), raylib.NewVector2(0, 0), 0, 1)
		//cameraMoveOffset = raylib.NewVector2(10, 10)

		canvasSize    = raylib.NewVector2(500, 430)
		canvasScale   = float32(1)
		canvas        = raylib.LoadRenderTexture(int32(canvasSize.X), int32(canvasSize.Y))
		canvasRefresh = true

		sidePanelWidth     = float32(350)
		sidePanelRelativeX = WindowWidth - int32(sidePanelWidth)

		drawing = false
		//drawingMode = ModeDrawing

		currentStroke = penTool{}
		strokes       = []penTool{currentStroke}
		undoneStrokes = []penTool{}

		colourPickerVal    = raylib.Orange
		colourPickerHeight = float32(200)

		brushSize    = float32(10)
		brushOpacity = float32(1)

		fileName        = "NewProject"
		fileNameEditing = false

		menu = StateNone
	)

	// check if userData exists
	if _, err := os.Stat(DirUserData); os.IsNotExist(err) {
		err := os.Mkdir(DirUserData, 0755)
		if err != nil {
			panic(err)
		}
	}

	WindowWidth = int32(raylib.GetScreenWidth())
	WindowHeight = int32(raylib.GetScreenHeight())

	refreshCanvas := func() {
		raylib.BeginTextureMode(canvas)
		raylib.ClearBackground(raylib.White)

		for i := 1; i < len(strokes); i++ {
			strokes[i].Draw(raylib.NewVector2(-10, -10))
		}

		raylib.EndTextureMode()
	}

	redoStroke := func() {
		if len(undoneStrokes) > 0 {
			strokes = append(strokes, undoneStrokes[len(undoneStrokes)-1])
			undoneStrokes = undoneStrokes[:len(undoneStrokes)-1]

			AddToast("Redo")
			canvasRefresh = true
		}
	}
	undoStroke := func() {
		// 1 because I dont know why
		if len(strokes) > 1 {
			undoneStrokes = append(undoneStrokes, strokes[len(strokes)-1])
			strokes = strokes[:len(strokes)-1]

			canvasRefresh = true
			AddToast("Undo")
		}
	}

	saveImage := func() {
		if fileName == "" {
			AddToast("Please enter a file name")
		} else {
			image := raylib.LoadImageFromTexture(canvas.Texture)

			raylib.ImageRotate(image, 180)
			raylib.ImageFlipHorizontal(image)

			raylib.ExportImage(*image, DirUserData+fileName+".png")

			AddToast("Drawing saved at " + DirUserData + fileName + ".png")
		}
	}

	for !ShouldQuit {
		// LOOP
		ShouldQuit = raylib.WindowShouldClose()

		if raylib.IsWindowResized() {
			WindowWidth = int32(raylib.GetScreenWidth())
			WindowHeight = int32(raylib.GetScreenHeight())

			sidePanelRelativeX = WindowWidth - int32(sidePanelWidth)
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
				AddToast("This is a toast message!")
			}

			if raylib.IsMouseButtonPressed(raylib.MouseLeftButton) {
				if raylib.CheckCollisionPointRec(raylib.GetMousePosition(), raylib.NewRectangle(float32(WindowWidth-int32(sidePanelWidth)), 0, sidePanelWidth, float32(WindowHeight))) {
					drawing = false
				} else if raylib.CheckCollisionPointRec(raylib.GetMousePosition(), raylib.NewRectangle(10, 10, canvasSize.X, canvasSize.Y)) {
					drawing = true
					currentStroke = penTool{
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
				currentStroke = penTool{}
				undoneStrokes = []penTool{}

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

			UpdateToasts()
		}

		// DRAW
		raylib.BeginDrawing()
		{
			raylib.ClearBackground(raylib.White)
			gui.Grid(raylib.NewRectangle(0, 0, float32(WindowWidth), float32(WindowHeight)), "", 30, 1, &raylib.Vector2{})

			// Canvas stuff
			raylib.BeginMode2D(camera)
			{
				raylib.DrawRectangle(20, 20, int32(canvasSize.X), int32(canvasSize.Y), raylib.Fade(raylib.Black, 0.3))
				raylib.DrawTexturePro(canvas.Texture, raylib.NewRectangle(0, 0, float32(canvas.Texture.Width), float32(-canvas.Texture.Height)), raylib.NewRectangle(10, 10, canvasSize.X, canvasSize.Y), raylib.Vector2{}, 0, raylib.White)

				//raylib.BeginScissorMode(10, 10, int32(canvasSize.X), int32(canvasSize.Y))
				currentStroke.Draw(raylib.NewVector2(0, 0))
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
			raylib.BeginScissorMode(sidePanelRelativeX, 0, int32(sidePanelWidth), WindowHeight)
			{
				raylib.DrawRectangle(sidePanelRelativeX, 0, int32(sidePanelWidth), WindowHeight, raylib.Fade(raylib.White, 0.7))

				if gui.Button(raylib.NewRectangle(float32(sidePanelRelativeX+10), 10, 25, 25), gui.IconText(gui.ICON_CROSS, "")) {
					menu = StateFileMenu
				}
				if gui.Button(raylib.NewRectangle(float32(sidePanelRelativeX+20+25), 10, 25, 25), gui.IconText(gui.ICON_FOLDER_SAVE, "")) {
					saveImage()
				}

				if gui.Button(raylib.NewRectangle(float32(WindowWidth-70), 10, 25, 25), gui.IconText(gui.ICON_UNDO, "")) {
					undoStroke()
				}
				if gui.Button(raylib.NewRectangle(float32(WindowWidth-35), 10, 25, 25), gui.IconText(gui.ICON_REDO, "")) {
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
			raylib.DrawRectangleLines(sidePanelRelativeX, 0, int32(sidePanelWidth), WindowHeight, raylib.Gray)

			// Info
			{
				var text string

				text = fmt.Sprintf("Strokes: %d | Points: %d", len(strokes), len(currentStroke.Points))
				gui.StatusBar(raylib.NewRectangle(0, float32(WindowHeight-20), 200, 20), text)

				text = fmt.Sprintf("Canvas Size: %dx%d | Scale: %v", int32(canvasSize.X), int32(canvasSize.Y), canvasScale)
				gui.StatusBar(raylib.NewRectangle(199, float32(WindowHeight-20), 200, 20), text)
			}

			switch menu {
			case StateFileMenu:
				raylib.DrawRectangle(0, 0, WindowWidth, WindowHeight, raylib.Fade(raylib.Black, 0.5))
				choice := gui.MessageBox(raylib.NewRectangle(float32(WindowWidth/2-200), float32(WindowHeight/2-100), 400, 200), "File", "This is a message box", "OK")
				if choice == 0 || choice == 1 {
					menu = StateNone
				}
			default:
				menu = StateNone
			}

			// Draw toasts
			DrawToasts()
		}
		raylib.EndDrawing()
	}

	// QUIT
	raylib.CloseAudioDevice()
	raylib.CloseWindow()

	// GOODBYE
}
