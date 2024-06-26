package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

type Canvas struct {
	Name string

	Size   raylib.Vector2
	Offset raylib.Vector2

	Target     raylib.RenderTexture2D
	Background raylib.Texture2D

	Strokes       []raylib.Texture2D
	UndoneStrokes []raylib.Texture2D

	UnsavedChanges bool
	EditingFile    bool
	Refresh        bool
}

func (c *Canvas) Update() {
	if c.Refresh {
		c.Target = raylib.LoadRenderTexture(int32(c.Size.X), int32(c.Size.Y))

		raylib.BeginTextureMode(c.Target)

		raylib.DrawTexturePro(
			c.Background,
			raylib.NewRectangle(0, 0, c.Size.X, -c.Size.Y),
			raylib.NewRectangle(0, 0, c.Size.X, c.Size.Y),
			raylib.Vector2Zero(),
			0,
			raylib.White,
		)

		for _, stroke := range c.Strokes {
			raylib.DrawTexturePro(
				stroke,
				raylib.NewRectangle(0, 0, c.Size.X, -c.Size.Y),
				raylib.NewRectangle(0, 0, c.Size.X, c.Size.Y),
				raylib.Vector2Zero(),
				0,
				raylib.White,
			)
		}
		raylib.EndTextureMode()

		c.UnsavedChanges = true
		c.Refresh = false
	}
}

func (c *Canvas) AddStroke(stroke raylib.Texture2D) {
	c.Strokes = append(c.Strokes, stroke)

	if len(c.UndoneStrokes) > 0 {
		for i := range c.UndoneStrokes {
			raylib.UnloadTexture(c.UndoneStrokes[i])
		}
		c.UndoneStrokes = []raylib.Texture2D{}
	}

	c.Refresh = true
}

func (c *Canvas) Undo() {
	if len(c.Strokes) > 0 {
		c.UndoneStrokes = append(c.UndoneStrokes, c.Strokes[len(c.Strokes)-1])
		c.Strokes = c.Strokes[:len(c.Strokes)-1]
		c.Refresh = true

		addToast("Undo")
	}
}

func (c *Canvas) Redo() {
	if len(c.UndoneStrokes) > 0 {
		c.Strokes = append(c.Strokes, c.UndoneStrokes[len(c.UndoneStrokes)-1])
		c.UndoneStrokes = c.UndoneStrokes[:len(c.UndoneStrokes)-1]
		c.Refresh = true

		addToast("Redo")
	}
}

func (c *Canvas) Draw() {
	raylib.DrawTexturePro(
		c.Target.Texture,
		raylib.NewRectangle(0, 0, c.Size.X, -c.Size.Y),
		raylib.NewRectangle(c.Offset.X, c.Offset.Y, c.Size.X, c.Size.Y),
		raylib.Vector2Zero(),
		0,
		raylib.White,
	)
}

func (c *Canvas) Save(force bool) {
	c.Name = strings.Trim(c.Name, " ")

	// check if the name is empty
	if c.Name == "" {
		addToast("Please enter a file name!")
		return
	}
	// check if file already exists
	if !c.EditingFile {
		_, err := os.Stat(filepath.Join(dirUserData, c.Name+".png"))
		if !errors.Is(err, os.ErrNotExist) {
			if !force {
				applicationState = StateFileExists
				return
			}
		}
	}

	image := raylib.LoadImageFromTexture(c.Target.Texture)

	raylib.ImageRotate(image, 180)
	raylib.ImageFlipHorizontal(image)

	raylib.ExportImage(*image, filepath.Join(dirUserData, c.Name+".png"))

	addToast("Drawing saved as " + c.Name + ".png")

	c.UnsavedChanges = false
	c.EditingFile = true
}

func NewCanvas(name string, size, offset raylib.Vector2, background raylib.Texture2D) *Canvas {
	return &Canvas{
		Name:           name,
		Size:           size,
		Offset:         offset,
		Target:         raylib.LoadRenderTexture(int32(size.X), int32(size.Y)),
		Background:     background,
		Strokes:        []raylib.Texture2D{},
		UndoneStrokes:  []raylib.Texture2D{},
		UnsavedChanges: false,
		EditingFile:    false,
		Refresh:        true,
	}
}

func NewBackgroundColour(size raylib.Vector2, color raylib.Color) raylib.Texture2D {
	texture := raylib.LoadRenderTexture(int32(size.X), int32(size.Y))

	fmt.Println(size)

	raylib.BeginTextureMode(texture)
	raylib.DrawRectangle(0, 0, int32(size.X), int32(size.Y), color)
	raylib.EndTextureMode()

	return texture.Texture
}

func NewBackgroundImage(pathToImage string) raylib.Texture2D {
	loadedImage := raylib.LoadImage(pathToImage)

	// For some reason Images are flipped horizontally and rotated 180 degrees, so we need to undo that...
	raylib.ImageFlipHorizontal(loadedImage)
	raylib.ImageRotate(loadedImage, 180)

	return raylib.LoadTextureFromImage(loadedImage)
}
