package application

var (
	WindowTitle        = "Colouring App"
	WindowWidth  int32 = 800
	WindowHeight int32 = 600
	WindowFPS    int32 = 144
)

const (
	ScenePlayerData = iota
	SceneTitle
	SceneOptions
	SceneDrawing
)

const (
	DirAssets   = "./assets/"
	DirUserData = "./userData/"
)

var (
	ShouldQuit   = false
	CurrentScene = ScenePlayerData
)
