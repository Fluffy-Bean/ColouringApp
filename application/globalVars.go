package application

const (
	WindowTitle  = "Colouring App"
	WindowWidth  = 800
	WindowHeight = 600
	WindowFPS    = 60
)

const (
	ScenePlayerData = iota
	SceneTitle
	SceneOptions
	SceneGame
)

const (
	DirAssets     = "./assets/"
	DirPlayerData = "./playerData/"
)

var CurrentScene = ScenePlayerData

var ShouldQuit = false
