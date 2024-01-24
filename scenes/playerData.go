package scenes

import (
	"ColouringApp/application"
	"os"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

func PlayerData() {
	// Load player data here
	for !application.ShouldQuit {
		// DEFAULT
		{
			application.ShouldQuit = raylib.WindowShouldClose()
			if application.CurrentScene != application.ScenePlayerData {
				break
			}
		}

		// check if userData exists
		if _, err := os.Stat(application.DirUserData); os.IsNotExist(err) {
			err := os.Mkdir(application.DirUserData, 0755)
			if err != nil {
				panic(err)
			}
		}

		// DRAW
		{
			raylib.BeginDrawing()
			raylib.ClearBackground(raylib.Black)
			raylib.DrawText("Loading...", 10, application.WindowHeight-30, 20, raylib.White)
			raylib.EndDrawing()
		}

		application.CurrentScene = application.SceneTitle
	}
}
