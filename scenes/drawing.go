package scenes

import (
	"ColouringApp/application"
	gui "github.com/gen2brain/raylib-go/raygui"
	raylib "github.com/gen2brain/raylib-go/raylib"
)

type stroke struct {
	Color  raylib.Color
	Size   float32
	Points []raylib.Vector2
}

func Drawing() {
	var (
		canvasSize    = raylib.NewVector2(500, 430)
		canvas        = raylib.LoadRenderTexture(int32(canvasSize.X), int32(canvasSize.Y))
		currentStroke = stroke{}
		strokes       = []stroke{currentStroke}
		undoneStrokes = []stroke{}

		brushSize float32 = 1
		color             = raylib.Orange
	)

	// Create canvas
	raylib.BeginTextureMode(canvas)
	raylib.ClearBackground(raylib.White)
	raylib.EndTextureMode()

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

	for !application.ShouldQuit {
		application.ShouldQuit = raylib.WindowShouldClose()
		if application.CurrentScene != application.SceneDrawing {
			break
		}

		// INPUT
		if raylib.IsMouseButtonPressed(raylib.MouseLeftButton) {
			currentStroke = stroke{
				Color: color,
				Size:  brushSize,
			}
		}
		if raylib.IsMouseButtonDown(raylib.MouseLeftButton) {
			// if mouse is further than 5 pixels from last point, add it
			var safeZone float32 = 5
			if len(currentStroke.Points) == 0 {
				currentStroke.Points = append(currentStroke.Points, raylib.GetMousePosition())
			} else if raylib.Vector2Distance(currentStroke.Points[len(currentStroke.Points)-1], raylib.GetMousePosition()) > safeZone {
				currentStroke.Points = append(currentStroke.Points, raylib.GetMousePosition())
			}
			//currentStroke.Points = append(currentStroke.Points, raylib.GetMousePosition())
		}
		if raylib.IsMouseButtonReleased(raylib.MouseLeftButton) && currentStroke.Points != nil {
			strokes = append(strokes, currentStroke)
			currentStroke = stroke{}
			undoneStrokes = []stroke{}
			refreshCanvas()
		}

		if raylib.IsKeyDown(raylib.KeyLeftControl) && raylib.IsKeyDown(raylib.KeyLeftShift) && raylib.IsKeyPressed(raylib.KeyZ) {
			if len(undoneStrokes) > 0 {
				strokes = append(strokes, undoneStrokes[len(undoneStrokes)-1])
				undoneStrokes = undoneStrokes[:len(undoneStrokes)-1]
			}
			refreshCanvas()
		} else if raylib.IsKeyDown(raylib.KeyLeftControl) && raylib.IsKeyPressed(raylib.KeyZ) {
			if len(strokes) > 0 {
				undoneStrokes = append(undoneStrokes, strokes[len(strokes)-1])
				strokes = strokes[:len(strokes)-1]
			}
			refreshCanvas()
		}

		// UPDATE

		// DRAW
		raylib.BeginDrawing()
		raylib.ClearBackground(raylib.LightGray)

		raylib.BeginScissorMode(10, 10, int32(canvasSize.X), int32(canvasSize.Y))
		raylib.DrawTexturePro(canvas.Texture, raylib.NewRectangle(0, 0, float32(canvas.Texture.Width), float32(-canvas.Texture.Height)), raylib.NewRectangle(10, 10, canvasSize.X, canvasSize.Y), raylib.Vector2{}, 0, raylib.White)
		for i := 1; i < len(currentStroke.Points); i++ {
			raylib.DrawLineEx(currentStroke.Points[i-1], currentStroke.Points[i], currentStroke.Size, currentStroke.Color)
			raylib.DrawCircle(int32(currentStroke.Points[i].X), int32(currentStroke.Points[i].Y), currentStroke.Size/2, currentStroke.Color)
		}
		raylib.EndScissorMode()
		raylib.DrawRectangleLines(10, 10, int32(canvasSize.X), int32(canvasSize.Y), raylib.Gray)

		color = gui.ColorPicker(raylib.NewRectangle(float32(30+int32(canvasSize.X)), 10, application.WindowWidth-float32(65+int32(canvasSize.X)), 200), "Color", color)

		brushSize = gui.Slider(raylib.NewRectangle(float32(90+int32(canvasSize.X)), 30+200, 200, 20), "Brush Size", "", brushSize, 1, 100)

		raylib.DrawLine(20+int32(canvasSize.X), 10, 20+int32(canvasSize.X), 10+int32(canvasSize.Y), raylib.Gray)
		raylib.DrawLine(30+int32(canvasSize.X), 20+200, application.WindowWidth-10, 20+200, raylib.Gray)

		if gui.Button(raylib.NewRectangle(float32(30+int32(canvasSize.X)), application.WindowHeight-35, application.WindowWidth-float32(40+int32(canvasSize.X)), 25), "Main Menu") {
			application.CurrentScene = application.SceneTitle
		}

		raylib.EndDrawing()
	}

	// unload resources here
}
