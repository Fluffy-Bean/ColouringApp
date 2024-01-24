package main

import (
	"ColouringApp/application"
	"ColouringApp/scenes"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	raylib.SetConfigFlags(raylib.FlagWindowResizable)

	raylib.InitWindow(application.WindowWidth, application.WindowHeight, application.WindowTitle)
	raylib.InitAudioDevice()

	raylib.SetTargetFPS(application.WindowFPS)
	//raylib.SetExitKey(0) // disable exit key

	// MAIN LOOP
	for !application.ShouldQuit {
		switch application.CurrentScene {
		case application.ScenePlayerData:
			scenes.PlayerData()
		case application.SceneTitle:
			scenes.Title()
		case application.SceneOptions:
			scenes.Options()
		case application.SceneDrawing:
			scenes.Drawing()
		default:
			panic("Unknown scene")
		}
	}

	// QUIT
	raylib.CloseAudioDevice()
	raylib.CloseWindow()

	// GOODBYE
}
