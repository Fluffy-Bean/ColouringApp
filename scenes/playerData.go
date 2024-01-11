package scenes

import (
	"ColouringApp/application"
	"time"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

func PlayerData() {
	// Load player data here
	for !application.ShouldQuit {
		application.ShouldQuit = raylib.WindowShouldClose()
		if application.CurrentScene != application.ScenePlayerData {
			break
		}

		raylib.BeginDrawing()
		raylib.ClearBackground(raylib.Black)

		raylib.DrawText("Loading...", 10, application.WindowHeight-30, 20, raylib.White)

		raylib.EndDrawing()

		// Just a placeholder
		time.Sleep(1 * time.Second)
		application.CurrentScene = application.SceneTitle
	}
}
