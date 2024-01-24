package scenes

import (
	"ColouringApp/application"
	"os"

	gui "github.com/gen2brain/raylib-go/raygui"
	raylib "github.com/gen2brain/raylib-go/raylib"
)

func Title() {
	var (
		titleText = application.WindowTitle
		gallery   []raylib.Texture2D
	)

	application.WindowWidth = int32(raylib.GetScreenWidth())
	application.WindowHeight = int32(raylib.GetScreenHeight())

	// Load gallery here
	files, err := os.ReadDir(application.DirUserData)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		gallery = append(gallery, raylib.LoadTexture(application.DirUserData+file.Name()))
	}

	for !application.ShouldQuit {
		// DEFAULT
		{
			application.ShouldQuit = raylib.WindowShouldClose()
			if application.CurrentScene != application.SceneTitle {
				break
			}
			if raylib.IsWindowResized() {
				application.WindowWidth = int32(raylib.GetScreenWidth())
				application.WindowHeight = int32(raylib.GetScreenHeight())
			}
		}

		// INPUT

		// UPDATE
		{
			application.UpdateToasts()
		}

		// DRAW
		{
			raylib.BeginDrawing()
			raylib.ClearBackground(raylib.White)

			if gui.Button(raylib.NewRectangle(10, 10, 40, 40), gui.IconText(gui.ICON_CROSS, "")) {
				application.ShouldQuit = true
			}
			raylib.DrawText(titleText, (application.WindowWidth-raylib.MeasureText(titleText, 20))/2, 20, 20, raylib.Black)
			if gui.Button(raylib.NewRectangle(float32(application.WindowWidth-50), 10, 40, 40), gui.IconText(gui.ICON_GEAR, "")) {
				application.CurrentScene = application.SceneOptions
			}

			for i := 0; i < len(gallery); i++ {
				raylib.DrawTexturePro(gallery[i], raylib.NewRectangle(0, 0, float32(gallery[i].Width), float32(-gallery[i].Height)), raylib.NewRectangle(float32(10+(i%5)*100), float32(70+(i/5)*100), 100, 100), raylib.Vector2{}, 0, raylib.White)
			}

			if gui.Button(raylib.NewRectangle(float32((application.WindowWidth-100)/2), float32(application.WindowHeight-70), 100, 40), "Start") {
				application.CurrentScene = application.SceneDrawing
			}

			application.DrawToasts()

			raylib.EndDrawing()
		}
	}

	for i := 0; i < len(gallery); i++ {
		raylib.UnloadTexture(gallery[i])
	}
}
