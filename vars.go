package main

var (
	WindowTitle        = "Colouring App"
	WindowWidth  int32 = 800
	WindowHeight int32 = 600
	WindowFPS    int32 = 144
)

const (
	StateNone = iota
	StateFileMenu
	StateDrawing
)

const (
	DirAssets   = "./assets/"
	DirUserData = "./userData/"
)

var (
	ShouldQuit = false
)
