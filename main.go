package main

import (
	"fmt"
	"os"
	"strconv"

	gui "github.com/gen2brain/raylib-go/raygui"
	raylib "github.com/gen2brain/raylib-go/raylib"
)

const (
	WindowTitle     = "Colouring App"
	WindowMinWidth  = int32(800)
	WindowMinHeight = int32(600)
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
	StateNormal = iota
	StateDrawing
	StateFileMenu
)

var canvas *Canvas

func checkDirs() {
	if _, err := os.Stat(DirAssets); os.IsNotExist(err) {
		panic("Assets directory not found")
	}
	if _, err := os.Stat(DirUserData); os.IsNotExist(err) {
		if err := os.Mkdir(DirUserData, 0755); err != nil {
			panic("Could not create userData directory")
		}
	}
}

func main() {
	checkDirs() // Make sure all the directories exist

	raylib.SetConfigFlags(raylib.FlagWindowResizable)
	// raylib.SetConfigFlags(raylib.FlagWindowHighdpi)
	// raylib.SetConfigFlags(raylib.FlagMsaa4xHint)

	raylib.InitWindow(WindowWidth, WindowHeight, WindowTitle)
	raylib.SetWindowMinSize(int(WindowMinWidth), int(WindowMinHeight))
	raylib.SetTargetFPS(int32(raylib.GetMonitorRefreshRate(raylib.GetCurrentMonitor())))
	// raylib.SetExitKey(0) // disable exit key

	// Augh
	var (
		camera = raylib.NewCamera2D(raylib.NewVector2(0, 0), raylib.NewVector2(0, 0), 0, 1)

		currentStroke = penTool{}

		sidePanelWidth     = float32(350)
		sidePanelRelativeX = WindowWidth - int32(sidePanelWidth)

		colourPickerVal    = raylib.Orange
		colourPickerHeight = float32(250)

		brushSize = float32(10)

		fileNameEditing = false

		state         = StateNormal
		appShouldQuit = false

		showCursor     = true
		showDebugStats = false
	)

	// New Canvas stuff
	var (
		createNewCanvas = true

		newProjectName        = "NewProject"
		newProjectNameEditing = false

		newCanvasWidth        = 700
		newCanvasWidthEditing = false

		newCanvasHeight        = 530
		newCanvasHeightEditing = false

		newCanvasBackgroundColor = raylib.White
	)

	var userDataProjects []string
	{
		f, err := os.Open(DirUserData)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		files, err := f.Readdir(-1)
		if err != nil {
			panic(err)
		}

		for _, file := range files {
			if file.Mode().IsRegular() {
				userDataProjects = append(userDataProjects, file.Name())
			}
		}
	}

	// LOOP
	for !appShouldQuit {
		appShouldQuit = raylib.WindowShouldClose()
		if raylib.IsWindowResized() {
			WindowWidth = int32(raylib.GetScreenWidth())
			WindowHeight = int32(raylib.GetScreenHeight())
			sidePanelRelativeX = WindowWidth - int32(sidePanelWidth)
		}

		// CREATE CANVAS
		if createNewCanvas {
			canvasBackground := NewBackground(raylib.NewVector2(float32(newCanvasWidth), float32(newCanvasHeight)), newCanvasBackgroundColor)
			canvas = NewCanvas("NewProject", raylib.NewVector2(float32(newCanvasWidth), float32(newCanvasHeight)), raylib.NewVector2(15, 15), canvasBackground)
			createNewCanvas = false
		}

		// INPUT
		{
			if raylib.IsKeyPressed(raylib.KeyF7) {
				AddToast("This is a test toast")
			}
			if raylib.IsKeyPressed(raylib.KeyF8) {
				showDebugStats = !showDebugStats
			}

			if raylib.IsMouseButtonPressed(raylib.MouseLeftButton) && state == StateNormal {
				if !raylib.CheckCollisionPointRec(raylib.GetMousePosition(), raylib.NewRectangle(float32(WindowWidth-int32(sidePanelWidth)), 0, sidePanelWidth, float32(WindowHeight))) &&
					raylib.CheckCollisionPointRec(raylib.GetMousePosition(), raylib.NewRectangle(10, 10, canvas.Size.X, canvas.Size.Y)) {
					state = StateDrawing
					currentStroke = penTool{
						Size:  brushSize,
						Color: colourPickerVal,
					}
				}
			}

			if raylib.IsMouseButtonDown(raylib.MouseLeftButton) && state == StateDrawing {
				var safeZone float32 = 1

				if len(currentStroke.Points) <= 1 {
					currentStroke.Points = append(currentStroke.Points, raylib.GetMousePosition())
				} else if raylib.Vector2Distance(currentStroke.Points[len(currentStroke.Points)-1], raylib.GetMousePosition()) > safeZone {
					currentStroke.Points = append(currentStroke.Points, raylib.GetMousePosition())
				}

				state = StateDrawing
			}

			if raylib.IsMouseButtonReleased(raylib.MouseLeftButton) && currentStroke.Points != nil {
				canvas.AddStroke(currentStroke.Render())
				currentStroke = penTool{}
				state = StateNormal
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
			UpdateToasts()
			canvas.Update()
			if state != StateNormal {
				gui.SetState(gui.STATE_DISABLED)
			} else {
				gui.SetState(gui.STATE_NORMAL)
			}

			if raylib.CheckCollisionPointRec(raylib.GetMousePosition(), raylib.NewRectangle(float32(WindowWidth-int32(sidePanelWidth)), 0, sidePanelWidth, float32(WindowHeight))) {
				showCursor = false
			} else {
				showCursor = true
			}
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
				currentStroke.Draw()
				raylib.EndScissorMode()

				raylib.DrawRectangleLines(int32(canvas.Offset.X), int32(canvas.Offset.Y), int32(canvas.Size.X), int32(canvas.Size.Y), raylib.DarkGray)
			}
			raylib.EndMode2D()

			// UI stuff
			raylib.BeginScissorMode(sidePanelRelativeX, 0, int32(sidePanelWidth), WindowHeight)
			{
				raylib.DrawRectangle(sidePanelRelativeX, 0, int32(sidePanelWidth), WindowHeight, raylib.Fade(raylib.White, 0.9))

				if gui.Button(raylib.NewRectangle(float32(sidePanelRelativeX+10), 10, 25, 25), gui.IconText(gui.ICON_FOLDER_OPEN, "")) {
					state = StateFileMenu
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
				raylib.DrawRectangleLines(sidePanelRelativeX, 0, int32(sidePanelWidth), WindowHeight, raylib.Gray)
			}
			raylib.EndScissorMode()

			// Info
			if showDebugStats {
				var text string

				text = fmt.Sprintf("Strokes: %d | Points: %d", len(canvas.Strokes), len(currentStroke.Points))
				gui.StatusBar(raylib.NewRectangle(0, float32(WindowHeight-20), 150, 20), text)

				text = fmt.Sprintf("Canvas Size: %dx%d", int32(canvas.Size.X), int32(canvas.Size.Y))
				gui.StatusBar(raylib.NewRectangle(150, float32(WindowHeight-20), 150, 20), text)

				text = fmt.Sprintf("FPS: %d | DT: %f", raylib.GetFPS(), raylib.GetFrameTime())
				gui.StatusBar(raylib.NewRectangle(300, float32(WindowHeight-20), 170, 20), text)
			}

			// Cursor
			if showCursor {
			}
			raylib.DrawCircleLines(int32(raylib.GetMousePosition().X), int32(raylib.GetMousePosition().Y), brushSize/2, raylib.Black)

			switch state {
			case StateFileMenu:
				gui.SetState(gui.STATE_NORMAL)
				raylib.DrawRectangle(0, 0, WindowWidth, WindowHeight, raylib.Fade(raylib.Black, 0.5))

				windowPos := raylib.NewRectangle(float32((WindowWidth/2)-200), float32((WindowHeight/2)-200), 400, 400)
				if gui.WindowBox(windowPos, "Open or New File") {
					state = StateNormal
				}

				// Magic numbers
				raylib.BeginScissorMode(int32(windowPos.X)+1, int32(windowPos.Y)+24, int32(windowPos.Width)-2, int32(windowPos.Height)-25)

				gui.GroupBox(raylib.NewRectangle(windowPos.X+11, windowPos.Y+34, windowPos.Width-22, 200), "Create New")
				{
					var err error

					gui.Label(raylib.NewRectangle(windowPos.X+21, windowPos.Y+44, 100, 20), "File Name")
					if gui.TextBox(raylib.NewRectangle(windowPos.X+131, windowPos.Y+44, windowPos.Width-152, 20), &newProjectName, 40, newProjectNameEditing) {
						newProjectNameEditing = !newProjectNameEditing
					}

					gui.Label(raylib.NewRectangle(windowPos.X+21, windowPos.Y+74, 100, 20), "Canvas Width")
					lastWidth := newCanvasWidth
					width := fmt.Sprintf("%d", newCanvasWidth)
					if gui.TextBox(raylib.NewRectangle(windowPos.X+131, windowPos.Y+74, windowPos.Width-152, 20), &width, 6, newCanvasWidthEditing) {
						newCanvasWidthEditing = !newCanvasWidthEditing
					}
					if newCanvasWidth, err = strconv.Atoi(width); err != nil {
						newCanvasWidth = lastWidth
					}

					gui.Label(raylib.NewRectangle(windowPos.X+21, windowPos.Y+104, 100, 20), "Canvas Height")
					lastHeight := newCanvasHeight
					height := fmt.Sprintf("%d", newCanvasHeight)
					if gui.TextBox(raylib.NewRectangle(windowPos.X+131, windowPos.Y+104, windowPos.Width-152, 20), &height, 6, newCanvasHeightEditing) {
						newCanvasHeightEditing = !newCanvasHeightEditing
					}
					if newCanvasHeight, err = strconv.Atoi(height); err != nil {
						newCanvasHeight = lastHeight
					}

					colors := []raylib.Color{raylib.Red, raylib.Orange, raylib.Yellow, raylib.Green, raylib.Blue, raylib.Purple, raylib.Pink, raylib.Brown, raylib.Black, raylib.White}
					for i := 0; i < len(colors); i++ {
						posX := windowPos.X + 21 + ((float32(i) * 26) + float32(i)*11)
						posY := windowPos.Y + 149

						if gui.Button(raylib.NewRectangle(posX, posY, 26, 26), "") {
							newCanvasBackgroundColor = colors[i]
						}

						if newCanvasBackgroundColor == colors[i] {
							raylib.DrawRectangle(int32(posX)-4, int32(posY)-4, 26+8, 26+8, raylib.Fade(raylib.Black, 0.2))
							raylib.DrawRectangle(int32(posX), int32(posY), 26, 26, colors[i])
						} else {
							raylib.DrawRectangle(int32(posX), int32(posY), 26, 26, colors[i])
							raylib.DrawRectangleLines(int32(posX), int32(posY), 26, 26, raylib.Black)
						}
					}

					if gui.Button(raylib.NewRectangle(windowPos.X+windowPos.Width-140, windowPos.Y+204, 120, 20), "Create") {
						state = StateNormal
						createNewCanvas = true
						AddToast("Created New Canvas")
					}
				}

				gui.GroupBox(raylib.NewRectangle(windowPos.X+11, windowPos.Y+244, windowPos.Width-22, float32(len(userDataProjects)*20)+10), "Open Existing")
				{
					if gui.Button(raylib.NewRectangle(windowPos.X+21, windowPos.Y+254, windowPos.Width-42, 20), "Maned Wolf") {
						loadedImage := raylib.LoadImage(DirAssets + "manedWolf.jpg")
						{
							raylib.ImageFlipHorizontal(loadedImage)
							raylib.ImageRotate(loadedImage, 180)
						}
						canvas = NewCanvas("NewProject", raylib.NewVector2(float32(loadedImage.Width), float32(loadedImage.Height)), raylib.NewVector2(15, 15), raylib.LoadTextureFromImage(loadedImage))

						raylib.UnloadImage(loadedImage)
						state = StateNormal

						AddToast("Loaded Maned Wolf")
					}

					for i := 0; i < len(userDataProjects); i++ {
						if gui.Button(raylib.NewRectangle(windowPos.X+21, windowPos.Y+274+float32(i*20), windowPos.Width-42, 20), userDataProjects[i]) {
							loadedImage := raylib.LoadImage(DirUserData + userDataProjects[i])
							{
								raylib.ImageFlipHorizontal(loadedImage)
								raylib.ImageRotate(loadedImage, 180)
							}
							canvas = NewCanvas("NewProject", raylib.NewVector2(float32(loadedImage.Width), float32(loadedImage.Height)), raylib.NewVector2(15, 15), raylib.LoadTextureFromImage(loadedImage))

							raylib.UnloadImage(loadedImage)
							state = StateNormal

							AddToast("Loaded " + userDataProjects[i])
						}
					}
				}

				raylib.EndScissorMode()

			default:
			}

			DrawToasts()
		}
		raylib.EndDrawing()
	}

	// GOODBYE
	raylib.CloseWindow()
}
