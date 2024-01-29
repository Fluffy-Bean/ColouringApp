package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	gui "github.com/gen2brain/raylib-go/raygui"
	raylib "github.com/gen2brain/raylib-go/raylib"
)

const (
	applicationTitle           = "Colouring App"
	applicationMinWindowWidth  = int32(800)
	applicationMinWindowHeight = int32(600)

	defaultProjectName   = "NewProject"
	defaultProjectWidth  = 700
	defaultProjectHeight = 530
)

const (
	DirAssets   = "./assets/"
	DirUserData = "./userData/"
)

const (
	StateNormal = iota
	StateDrawing
	StateFileMenu
	StateNewCanvas
)

var (
	applicationState           = StateNormal
	applicationShouldQuit      = false
	applicationShowDebugValues = false
	applicationWindowWidth     = applicationMinWindowWidth
	applicationWindowHeight    = applicationMinWindowHeight
)

var (
	newPenStroke         = penTool{}
	newPenStrokeSafeZone = 1

	toolPanelWidth  = float32(350)
	toolPanelOffset = applicationWindowWidth - int32(toolPanelWidth)

	toolPanelColourPicker       = raylib.Orange
	toolPanelColourPickerHeight = float32(250)

	toolPanelBrushSize = float32(10)

	isEditingCanvasName = false
)

var (
	canvas *Canvas

	shouldCreateNewCanvas = true

	newCanvasName          = defaultProjectName
	isEditingNewCanvasName = false

	newCanvasWidth          = defaultProjectWidth
	isEditingNewCanvasWidth = false

	newCanvasHeight          = defaultProjectHeight
	isEditingNewCanvasHeight = false

	newCanvasColor     = raylib.White
	newCanvasImagePath = ""
)

func main() {
	raylib.SetConfigFlags(raylib.FlagWindowResizable | raylib.FlagWindowHighdpi | raylib.FlagMsaa4xHint)

	raylib.InitWindow(applicationWindowWidth, applicationWindowHeight, applicationTitle)

	raylib.SetWindowMinSize(int(applicationMinWindowWidth), int(applicationMinWindowHeight))
	raylib.SetTargetFPS(int32(raylib.GetMonitorRefreshRate(raylib.GetCurrentMonitor())))
	raylib.SetExitKey(0) // disable exit key

	if _, err := os.Stat(DirAssets); os.IsNotExist(err) {
		panic("Assets directory not found")
	}
	if _, err := os.Stat(DirUserData); os.IsNotExist(err) {
		if err := os.Mkdir(DirUserData, 0755); err != nil {
			panic("Could not create userData directory")
		}
	}

	var userDataProjects []string
	{
		f, err := os.Open(DirUserData)
		if err != nil {
			panic(err)
		}
		defer func() {
			err := f.Close()
			if err != nil {
				panic(err)
			}
		}()

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
	for !applicationShouldQuit {
		applicationShouldQuit = raylib.WindowShouldClose()
		if raylib.IsWindowResized() {
			applicationWindowWidth = int32(raylib.GetScreenWidth())
			applicationWindowHeight = int32(raylib.GetScreenHeight())
			toolPanelOffset = applicationWindowWidth - int32(toolPanelWidth)
		}

		// CREATE CANVAS
		if shouldCreateNewCanvas {
			var canvasBackground raylib.Texture2D

			if newCanvasImagePath != "" {
				// For some reason Images are flipped horizontally and rotated 180 degrees, so we need to undo that...
				loadedImage := raylib.LoadImage(newCanvasImagePath)
				raylib.ImageFlipHorizontal(loadedImage)
				raylib.ImageRotate(loadedImage, 180)

				canvasBackground = raylib.LoadTextureFromImage(loadedImage)

				newCanvasWidth = int(loadedImage.Width)
				newCanvasHeight = int(loadedImage.Height)
			} else {
				canvasBackground = NewBackground(raylib.NewVector2(float32(newCanvasWidth), float32(newCanvasHeight)), newCanvasColor)
			}

			canvas = NewCanvas(newCanvasName, raylib.NewVector2(float32(newCanvasWidth), float32(newCanvasHeight)), raylib.NewVector2(15, 15), canvasBackground)

			shouldCreateNewCanvas = false
			newCanvasName = defaultProjectName
			newCanvasWidth = defaultProjectWidth
			newCanvasHeight = defaultProjectHeight
			newCanvasColor = raylib.White
			newCanvasImagePath = ""
		}

		// INPUT
		{
			if raylib.IsKeyPressed(raylib.KeyF7) {
				AddToast("This is a test toast")
			}
			if raylib.IsKeyPressed(raylib.KeyF8) {
				applicationShowDebugValues = !applicationShowDebugValues
			}
			if raylib.IsKeyPressed(raylib.KeyF12) {
				AddToast("Screenshot saved!")
			}

			if raylib.IsMouseButtonPressed(raylib.MouseLeftButton) && applicationState == StateNormal {
				if !raylib.CheckCollisionPointRec(raylib.GetMousePosition(), raylib.NewRectangle(float32(applicationWindowWidth-int32(toolPanelWidth)), 0, toolPanelWidth, float32(applicationWindowHeight))) &&
					raylib.CheckCollisionPointRec(raylib.GetMousePosition(), raylib.NewRectangle(10, 10, canvas.Size.X, canvas.Size.Y)) {
					applicationState = StateDrawing
					newPenStroke = penTool{
						Size:  toolPanelBrushSize,
						Color: toolPanelColourPicker,
					}
				}
			}

			if raylib.IsMouseButtonDown(raylib.MouseLeftButton) && applicationState == StateDrawing {
				distanceToLastPoint := raylib.Vector2Distance(newPenStroke.Points[len(newPenStroke.Points)-1], raylib.GetMousePosition())

				if distanceToLastPoint > float32(newPenStrokeSafeZone) {
					newPenStroke.Points = append(newPenStroke.Points, raylib.GetMousePosition())
				} else if len(newPenStroke.Points) <= newPenStrokeSafeZone {
					newPenStroke.Points = append(newPenStroke.Points, raylib.GetMousePosition())
				}

				applicationState = StateDrawing
			}

			if raylib.IsMouseButtonReleased(raylib.MouseLeftButton) && newPenStroke.Points != nil {
				canvas.AddStroke(newPenStroke.Render())
				newPenStroke = penTool{}
				applicationState = StateNormal
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
			if applicationState != StateNormal {
				gui.Lock()
			} else {
				gui.Unlock()
			}
		}

		// DRAW
		raylib.BeginDrawing()
		{
			raylib.ClearBackground(raylib.White)
			gui.Grid(raylib.NewRectangle(0, 0, float32(applicationWindowWidth), float32(applicationWindowHeight)), "", 30, 1, &raylib.Vector2{})

			// Canvas stuff
			{
				raylib.DrawRectangle(int32(canvas.Offset.X)+10, int32(canvas.Offset.Y)+10, int32(canvas.Size.X), int32(canvas.Size.Y), raylib.Fade(raylib.Black, 0.3))
				canvas.Draw()

				raylib.BeginScissorMode(int32(canvas.Offset.X), int32(canvas.Offset.Y), int32(canvas.Size.X), int32(canvas.Size.Y))
				newPenStroke.Draw()
				raylib.EndScissorMode()

				raylib.DrawRectangleLines(int32(canvas.Offset.X), int32(canvas.Offset.Y), int32(canvas.Size.X), int32(canvas.Size.Y), raylib.DarkGray)
			}

			// UI stuff
			raylib.BeginScissorMode(toolPanelOffset, 0, int32(toolPanelWidth), applicationWindowHeight)
			{
				raylib.DrawRectangle(toolPanelOffset, 0, int32(toolPanelWidth), applicationWindowHeight, raylib.Fade(raylib.White, 0.9))

				if gui.Button(raylib.NewRectangle(float32(toolPanelOffset+10), 10, 25, 25), gui.IconText(gui.ICON_FOLDER_OPEN, "")) {
					applicationState = StateFileMenu
				}
				if gui.Button(raylib.NewRectangle(float32(toolPanelOffset+20+25), 10, 25, 25), gui.IconText(gui.ICON_FOLDER_SAVE, "")) {
					canvas.Save()
				}

				if gui.Button(raylib.NewRectangle(float32(applicationWindowWidth-70), 10, 25, 25), gui.IconText(gui.ICON_UNDO, "")) {
					canvas.Undo()
				}
				if gui.Button(raylib.NewRectangle(float32(applicationWindowWidth-35), 10, 25, 25), gui.IconText(gui.ICON_REDO, "")) {
					canvas.Redo()
				}

				toolPanelColourPicker = gui.ColorPicker(raylib.NewRectangle(float32(toolPanelOffset+10), 45, toolPanelWidth-45, toolPanelColourPickerHeight), "Color", toolPanelColourPicker)

				gui.Label(raylib.NewRectangle(float32(toolPanelOffset+10), 55+toolPanelColourPickerHeight, 60, 20), "Brush Size")
				toolPanelBrushSize = gui.Slider(raylib.NewRectangle(float32(toolPanelOffset+80), 55+toolPanelColourPickerHeight, toolPanelWidth-90, 20), "", "", toolPanelBrushSize, 1, 100)

				gui.Label(raylib.NewRectangle(float32(toolPanelOffset+10), 115+toolPanelColourPickerHeight, 60, 20), "File Name")
				if gui.TextBox(raylib.NewRectangle(float32(toolPanelOffset+80), 115+toolPanelColourPickerHeight, toolPanelWidth-90, 20), &canvas.Name, 40, isEditingCanvasName) {
					isEditingCanvasName = !isEditingCanvasName
				}
				raylib.DrawRectangleLines(toolPanelOffset, 0, int32(toolPanelWidth), applicationWindowHeight, raylib.Gray)
			}
			raylib.EndScissorMode()

			// Info
			if applicationShowDebugValues {
				var text string

				text = fmt.Sprintf("Strokes: %d | Points: %d", len(canvas.Strokes), len(newPenStroke.Points))
				gui.StatusBar(raylib.NewRectangle(0, float32(applicationWindowHeight-20), 150, 20), text)

				text = fmt.Sprintf("Canvas Size: %dx%d", int32(canvas.Size.X), int32(canvas.Size.Y))
				gui.StatusBar(raylib.NewRectangle(150, float32(applicationWindowHeight-20), 150, 20), text)

				text = fmt.Sprintf("FPS: %d | DT: %f", raylib.GetFPS(), raylib.GetFrameTime())
				gui.StatusBar(raylib.NewRectangle(300, float32(applicationWindowHeight-20), 170, 20), text)
			}

			// Cursor
			raylib.DrawCircleLines(int32(raylib.GetMousePosition().X), int32(raylib.GetMousePosition().Y), toolPanelBrushSize/2, raylib.Black)

			switch applicationState {
			case StateFileMenu:
				windowPos := raylib.NewRectangle(float32((applicationWindowWidth/2)-200), float32((applicationWindowHeight/2)-200), 400, 400)

				raylib.DrawRectangle(0, 0, applicationWindowWidth, applicationWindowHeight, raylib.Fade(raylib.Black, 0.5))

				gui.Unlock()
				if gui.WindowBox(windowPos, "Open or New File") {
					applicationState = StateNormal
				}

				// Magic numbers
				raylib.BeginScissorMode(int32(windowPos.X)+1, int32(windowPos.Y)+24, int32(windowPos.Width)-2, int32(windowPos.Height)-25)

				gui.GroupBox(raylib.NewRectangle(windowPos.X+11, windowPos.Y+34, windowPos.Width-22, 200), "Create New")
				{
					var err error

					gui.Label(raylib.NewRectangle(windowPos.X+21, windowPos.Y+44, 100, 20), "File Name")
					if gui.TextBox(raylib.NewRectangle(windowPos.X+131, windowPos.Y+44, windowPos.Width-152, 20), &newCanvasName, 40, isEditingNewCanvasName) {
						isEditingNewCanvasName = !isEditingNewCanvasName
					}

					gui.Label(raylib.NewRectangle(windowPos.X+21, windowPos.Y+74, 100, 20), "Canvas Width")
					lastWidth := newCanvasWidth
					width := fmt.Sprintf("%d", newCanvasWidth)
					if gui.TextBox(raylib.NewRectangle(windowPos.X+131, windowPos.Y+74, windowPos.Width-152, 20), &width, 6, isEditingNewCanvasWidth) {
						isEditingNewCanvasWidth = !isEditingNewCanvasWidth
					}
					if newCanvasWidth, err = strconv.Atoi(width); err != nil {
						newCanvasWidth = lastWidth
					}

					gui.Label(raylib.NewRectangle(windowPos.X+21, windowPos.Y+104, 100, 20), "Canvas Height")
					lastHeight := newCanvasHeight
					height := fmt.Sprintf("%d", newCanvasHeight)
					if gui.TextBox(raylib.NewRectangle(windowPos.X+131, windowPos.Y+104, windowPos.Width-152, 20), &height, 6, isEditingNewCanvasHeight) {
						isEditingNewCanvasHeight = !isEditingNewCanvasHeight
					}
					if newCanvasHeight, err = strconv.Atoi(height); err != nil {
						newCanvasHeight = lastHeight
					}

					colors := []raylib.Color{raylib.Red, raylib.Orange, raylib.Yellow, raylib.Green, raylib.Blue, raylib.Purple, raylib.Pink, raylib.Brown, raylib.Black, raylib.White}
					for i := 0; i < len(colors); i++ {
						posX := windowPos.X + 21 + ((float32(i) * 26) + float32(i)*11)
						posY := windowPos.Y + 149

						if gui.Button(raylib.NewRectangle(posX, posY, 26, 26), "") {
							newCanvasColor = colors[i]
						}

						if newCanvasColor == colors[i] {
							raylib.DrawRectangle(int32(posX)-4, int32(posY)-4, 26+8, 26+8, raylib.Fade(raylib.Black, 0.2))
							raylib.DrawRectangle(int32(posX), int32(posY), 26, 26, colors[i])
						} else {
							raylib.DrawRectangle(int32(posX), int32(posY), 26, 26, colors[i])
							raylib.DrawRectangleLines(int32(posX), int32(posY), 26, 26, raylib.Black)
						}
					}

					if gui.Button(raylib.NewRectangle(windowPos.X+windowPos.Width-140, windowPos.Y+204, 120, 20), "Create") {
						applicationState = StateNewCanvas
					}
				}

				gui.GroupBox(raylib.NewRectangle(windowPos.X+11, windowPos.Y+244, windowPos.Width-22, float32(len(userDataProjects)*20)+10), "Open Existing")
				{
					if gui.Button(raylib.NewRectangle(windowPos.X+21, windowPos.Y+254, windowPos.Width-42, 20), "Maned Wolf") {
						newCanvasImagePath = DirAssets + "manedWolf.jpg"
						newCanvasName = "ManedWolf"
						applicationState = StateNewCanvas
					}

					for i := 0; i < len(userDataProjects); i++ {
						if gui.Button(raylib.NewRectangle(windowPos.X+21, windowPos.Y+274+float32(i*20), windowPos.Width-42, 20), userDataProjects[i]) {
							newCanvasImagePath = DirUserData + userDataProjects[i]
							splitName := strings.Split(userDataProjects[i], ".")
							newCanvasName = splitName[:len(splitName)-1][0]
							applicationState = StateNewCanvas
						}
					}
				}

				raylib.EndScissorMode()
			case StateNewCanvas:
				windowPos := raylib.NewRectangle(float32((applicationWindowWidth/2)-200), float32((applicationWindowHeight/2)-75), 400, 150)

				raylib.DrawRectangle(0, 0, applicationWindowWidth, applicationWindowHeight, raylib.Fade(raylib.Black, 0.5))

				gui.Unlock()
				choice := gui.MessageBox(windowPos, "New Canvas", "Are you sure you want to create a new canvas?", "No;Yes")

				if choice == 0 || choice == 1 {
					applicationState = StateNormal
					shouldCreateNewCanvas = false
				} else if choice == 2 {
					applicationState = StateNormal
					shouldCreateNewCanvas = true
					AddToast("Created New Canvas: " + canvas.Name)
				}
			default:
			}

			DrawToasts()
		}
		raylib.EndDrawing()
	}

	// GOODBYE
	raylib.CloseWindow()
}
