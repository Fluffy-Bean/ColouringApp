package application

const (
	// WindowTitle used as the Windows Title but also used for the name of the game
	WindowTitle = "Colouring App"

	// Defines the window proportions
	WindowWidth  = 800
	WindowHeight = 600

	// WindowFPS Max FPS the game should run at, used to calucalte delta time
	WindowFPS = 60
)

// Scene IDs used to determine which scene to run
const (
	ScenePlayerData = iota
	SceneTitle
	SceneOptions
	SceneGallery
	SceneDraw
	SceneGame
)

// Directories used to store assets
const (
	DirAssets     = "./assets/"
	DirPlayerData = "./playerData/"
)

var (
	// ShouldQuit is used to determine if the game should quit
	ShouldQuit = false

	// CurrentScene is the scene that is currently running
	CurrentScene = ScenePlayerData
)
