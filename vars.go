package main

import raylib "github.com/gen2brain/raylib-go/raylib"

const (
	applicationTitle           = "Colouring App"
	applicationMinWindowWidth  = int32(800)
	applicationMinWindowHeight = int32(600)

	defaultProjectName   = "NewProject"
	defaultProjectWidth  = 700
	defaultProjectHeight = 530

	dirAssets   = "assets"
	dirUserData = "userData"
)

const (
	StateNormal = iota
	StateDrawing
	StateFileMenu
	StateNewCanvas
	StateWindowWantsToDie
)

const (
	toolNone = iota
	toolPen
)

var (
	applicationState           = StateNormal
	applicationShouldQuit      = false
	applicationShowDebugValues = false
	applicationWindowWidth     = applicationMinWindowWidth
	applicationWindowHeight    = applicationMinWindowHeight
	applicationRuntime         = float32(0)
)

var (
	cursorColor = raylib.Black

	newStrokeType     = toolNone
	newPenStroke      = penTool{}
	newStrokeSafeZone = 1

	toolPanelWidth  = float32(350)
	toolPanelOffset = applicationWindowWidth - int32(toolPanelWidth)

	toolPanelColourPicker       = raylib.Orange
	toolPanelColourPickerHeight = float32(250)

	toolPanelBrushSize = float32(10)

	isEditingCanvasName = false
)

var (
	canvas *Canvas

	shouldCreateNewCanvas = true

	newCanvasName          = defaultProjectName
	isEditingNewCanvasName = false

	newCanvasWidth          = defaultProjectWidth
	isEditingNewCanvasWidth = false

	newCanvasHeight          = defaultProjectHeight
	isEditingNewCanvasHeight = false

	newCanvasColor     = raylib.White
	newCanvasImagePath = ""
)
