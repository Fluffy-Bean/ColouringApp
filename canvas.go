package main

import (
	"fmt"
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

	Refresh bool
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

		AddToast("Undo")
	}
}

func (c *Canvas) Redo() {
	if len(c.UndoneStrokes) > 0 {
		c.Strokes = append(c.Strokes, c.UndoneStrokes[len(c.UndoneStrokes)-1])
		c.UndoneStrokes = c.UndoneStrokes[:len(c.UndoneStrokes)-1]
		c.Refresh = true

		AddToast("Redo")
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

func (c *Canvas) Save() {
	c.Name = strings.Trim(c.Name, " ")

	if c.Name == "" {
		AddToast("Please enter a file name!")
	} else {
		image := raylib.LoadImageFromTexture(c.Target.Texture)

		raylib.ImageRotate(image, 180)
		raylib.ImageFlipHorizontal(image)

		raylib.ExportImage(*image, DirUserData+c.Name+".png")

		AddToast("Drawing saved as " + c.Name + ".png")
	}
}

func NewCanvas(name string, size, offset raylib.Vector2, background raylib.Texture2D) *Canvas {
	return &Canvas{
		Name:          name,
		Size:          size,
		Offset:        offset,
		Target:        raylib.LoadRenderTexture(int32(size.X), int32(size.Y)),
		Background:    background,
		Strokes:       []raylib.Texture2D{},
		UndoneStrokes: []raylib.Texture2D{},
		Refresh:       true,
	}
}

func NewBackground(size raylib.Vector2, color raylib.Color) raylib.Texture2D {
	texture := raylib.LoadRenderTexture(int32(size.X), int32(size.Y))

	fmt.Println(size)

	raylib.BeginTextureMode(texture)
	raylib.DrawRectangle(0, 0, int32(size.X), int32(size.Y), color)
	raylib.EndTextureMode()

	return texture.Texture
}
