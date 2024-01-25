package main

import (
	"fmt"
	"os"

	gui "github.com/gen2brain/raylib-go/raygui"
	raylib "github.com/gen2brain/raylib-go/raylib"
)

const (
	WindowTitle     = "Colouring App"
	WindowMinWidth  = int32(800)
	WindowMinHeight = int32(600)
	WindowFPS       = int32(144)
)

var (
	WindowWidth  = WindowMinWidth
	WindowHeight = WindowMinHeight
)

const (
	DirAssets   = "./assets/"
	DirUserData = "./userData/"
)

const (
	StateNone = iota
	StateDrawing
	StateFileMenu
)

func checkDirs() {
	if _, err := os.Stat(DirUserData); os.IsNotExist(err) {
		err := os.Mkdir(DirUserData, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	raylib.SetConfigFlags(raylib.FlagWindowResizable)
	raylib.SetTraceLogLevel(raylib.LogTrace)
	raylib.SetConfigFlags(raylib.FlagMsaa4xHint)

	raylib.InitWindow(WindowWidth, WindowHeight, WindowTitle)
	raylib.SetWindowMinSize(int(WindowMinWidth), int(WindowMinHeight))

	raylib.InitAudioDevice()

	raylib.SetTargetFPS(WindowFPS)
	//raylib.SetExitKey(0) // disable exit key

	var (
		camera = raylib.NewCamera2D(raylib.NewVector2(0, 0), raylib.NewVector2(0, 0), 0, 1)

		canvas        = NewCanvas("NewProject", raylib.NewVector2(600, 530), raylib.NewVector2(15, 15))
		currentStroke = penTool{}
		drawing       = false

		sidePanelWidth     = float32(350)
		sidePanelRelativeX = WindowWidth - int32(sidePanelWidth)

		colourPickerVal    = raylib.Orange
		colourPickerHeight = float32(250)

		brushSize = float32(10)

		fileNameEditing = false

		menu = StateNone

		appShouldQuit = false
	)

	// check if userData exists
	checkDirs()

	for !appShouldQuit {
		// LOOP
		appShouldQuit = raylib.WindowShouldClose()

		if raylib.IsWindowResized() {
			WindowWidth = int32(raylib.GetScreenWidth())
			WindowHeight = int32(raylib.GetScreenHeight())
			sidePanelRelativeX = WindowWidth - int32(sidePanelWidth)
		}

		// INPUT
		{
			if raylib.IsKeyPressed(raylib.KeyF8) {
				AddToast("This is a toast message!")
			}

			if raylib.IsMouseButtonPressed(raylib.MouseLeftButton) {
				if raylib.CheckCollisionPointRec(raylib.GetMousePosition(), raylib.NewRectangle(float32(WindowWidth-int32(sidePanelWidth)), 0, sidePanelWidth, float32(WindowHeight))) {
					drawing = false
				} else if raylib.CheckCollisionPointRec(raylib.GetMousePosition(), raylib.NewRectangle(10, 10, canvas.Size.X, canvas.Size.Y)) {
					drawing = true
					currentStroke = penTool{
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
				canvas.Strokes = append(canvas.Strokes, currentStroke)
				canvas.UndoneStrokes = []penTool{}
				canvas.Refresh = true

				currentStroke = penTool{}
				drawing = false
			}

			if raylib.IsKeyDown(raylib.KeyLeftControl) && raylib.IsKeyDown(raylib.KeyLeftShift) && raylib.IsKeyPressed(raylib.KeyZ) {
				canvas.Redo()
			} else if raylib.IsKeyDown(raylib.KeyLeftControl) && raylib.IsKeyPressed(raylib.KeyZ) {
				canvas.Undo()
			} else if raylib.IsKeyDown(raylib.KeyLeftControl) && raylib.IsKeyPressed(raylib.KeyS) {
				canvas.Save()
			}
		}

		// UPDATE
		{
			if drawing {
				gui.SetState(gui.STATE_DISABLED)
			} else {
				gui.SetState(gui.STATE_NORMAL)
			}

			canvas.Update()

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
				raylib.DrawRectangle(int32(canvas.Offset.X)+10, int32(canvas.Offset.Y)+10, int32(canvas.Size.X), int32(canvas.Size.Y), raylib.Fade(raylib.Black, 0.3))
				canvas.Draw()

				raylib.BeginScissorMode(int32(canvas.Offset.X), int32(canvas.Offset.Y), int32(canvas.Size.X), int32(canvas.Size.Y))
				currentStroke.Draw(raylib.NewVector2(0, 0))
				raylib.EndScissorMode()

				if drawing {
					raylib.DrawRectangleLines(int32(canvas.Offset.X), int32(canvas.Offset.Y), int32(canvas.Size.X), int32(canvas.Size.Y), raylib.DarkGray)
					raylib.DrawCircleLines(int32(raylib.GetMousePosition().X), int32(raylib.GetMousePosition().Y), brushSize/2, raylib.Black)
				} else {
					raylib.DrawRectangleLines(int32(canvas.Offset.X), int32(canvas.Offset.Y), int32(canvas.Size.X), int32(canvas.Size.Y), raylib.Gray)
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
					canvas.Save()
				}

				if gui.Button(raylib.NewRectangle(float32(WindowWidth-70), 10, 25, 25), gui.IconText(gui.ICON_UNDO, "")) {
					canvas.Undo()
				}
				if gui.Button(raylib.NewRectangle(float32(WindowWidth-35), 10, 25, 25), gui.IconText(gui.ICON_REDO, "")) {
					canvas.Redo()
				}

				colourPickerVal = gui.ColorPicker(raylib.NewRectangle(float32(sidePanelRelativeX+10), 45, sidePanelWidth-45, colourPickerHeight), "Color", colourPickerVal)

				gui.Label(raylib.NewRectangle(float32(sidePanelRelativeX+10), 55+colourPickerHeight, 60, 20), "Brush Size")
				brushSize = gui.Slider(raylib.NewRectangle(float32(sidePanelRelativeX+80), 55+colourPickerHeight, sidePanelWidth-90, 20), "", "", brushSize, 1, 100)

				gui.Label(raylib.NewRectangle(float32(sidePanelRelativeX+10), 115+colourPickerHeight, 60, 20), "File Name")
				if gui.TextBox(raylib.NewRectangle(float32(sidePanelRelativeX+80), 115+colourPickerHeight, sidePanelWidth-90, 20), &canvas.Name, 40, fileNameEditing) {
					fileNameEditing = !fileNameEditing
				}
			}
			raylib.EndScissorMode()
			raylib.DrawRectangleLines(sidePanelRelativeX, 0, int32(sidePanelWidth), WindowHeight, raylib.Gray)

			// Info
			{
				var text string

				text = fmt.Sprintf("Strokes: %d | Points: %d", len(canvas.Strokes), len(currentStroke.Points))
				gui.StatusBar(raylib.NewRectangle(0, float32(WindowHeight-20), 200, 20), text)

				text = fmt.Sprintf("Canvas Size: %dx%d", int32(canvas.Size.X), int32(canvas.Size.Y))
				gui.StatusBar(raylib.NewRectangle(199, float32(WindowHeight-20), 200, 20), text)
			}

			switch menu {
			case StateFileMenu:
				raylib.DrawRectangle(0, 0, WindowWidth, WindowHeight, raylib.Fade(raylib.Black, 0.5))
				choice := gui.MessageBox(raylib.NewRectangle(float32(WindowWidth/2-200), float32(WindowHeight/2-100), 400, 200), "File", "This is a message box", "Ok")
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
