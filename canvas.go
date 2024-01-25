package main

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
)

type Canvas struct {
	Name string

	Size   raylib.Vector2
	Offset raylib.Vector2

	Target raylib.RenderTexture2D

	Strokes       []penTool
	UndoneStrokes []penTool

	Refresh bool
}

func (c *Canvas) Update() {
	if c.Refresh {
		c.Target = raylib.LoadRenderTexture(int32(c.Size.X), int32(c.Size.Y))

		raylib.BeginTextureMode(c.Target)
		raylib.ClearBackground(raylib.White)
		for _, mark := range c.Strokes {
			mark.Draw(raylib.Vector2Scale(c.Offset, -1))
		}
		raylib.EndTextureMode()

		c.Refresh = false
	}
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

func NewCanvas(name string, size, offset raylib.Vector2) *Canvas {
	return &Canvas{
		Name:          name,
		Size:          size,
		Offset:        offset,
		Target:        raylib.LoadRenderTexture(int32(size.X), int32(size.Y)),
		Strokes:       []penTool{},
		UndoneStrokes: []penTool{},
		Refresh:       true,
	}
}
