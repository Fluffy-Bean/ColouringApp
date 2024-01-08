package application

const (
	WindowTitle  = "Colouring App"
	WindowWidth  = 800
	WindowHeight = 600
	WindowFPS    = 60
)

const (
	SceneTitle = iota
	SceneOptions
	SceneGame
)

var CurrentScene = SceneTitle

var ShouldQuit = false
