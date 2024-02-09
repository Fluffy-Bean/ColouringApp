package main

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
)

type penTool struct {
	Size   float32
	Color  raylib.Color
	Points []raylib.Vector2
}

func (p *penTool) Render() raylib.Texture2D {
	offset := raylib.Vector2Scale(canvas.Offset, -1)
	texture := raylib.LoadRenderTexture(int32(canvas.Size.X), int32(canvas.Size.Y))

	raylib.BeginTextureMode(texture)
	raylib.ClearBackground(raylib.Fade(raylib.Black, 0))
	for i := 0; i < len(p.Points)-1; i++ {
		startPointOffset := raylib.Vector2Add(p.Points[i], offset)
		endPointOffset := raylib.Vector2Add(p.Points[i+1], offset)
		raylib.DrawLineEx(startPointOffset, endPointOffset, p.Size, p.Color)
		raylib.DrawCircle(int32(startPointOffset.X), int32(startPointOffset.Y), p.Size/2, p.Color)
	}
	if len(p.Points) > 0 {
		endPointOffset := raylib.Vector2Add(p.Points[len(p.Points)-1], offset)
		raylib.DrawCircle(int32(endPointOffset.X), int32(endPointOffset.Y), p.Size/2, p.Color)
	}
	raylib.EndTextureMode()

	return texture.Texture
}

func (p *penTool) Draw() {
	for i := 0; i < len(p.Points)-1; i++ {
		startPoint := p.Points[i]
		endPoint := p.Points[i+1]
		raylib.DrawLineEx(startPoint, endPoint, p.Size, p.Color)
		raylib.DrawCircle(int32(startPoint.X), int32(startPoint.Y), p.Size/2, p.Color)
	}
	if len(p.Points) > 0 {
		endPoint := p.Points[len(p.Points)-1]
		raylib.DrawCircle(int32(endPoint.X), int32(endPoint.Y), p.Size/2, p.Color)
	}
}
