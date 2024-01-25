package main

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
)

type penTool struct {
	Color  raylib.Color
	Size   float32
	Points []raylib.Vector2
}

func (p *penTool) Draw(offset raylib.Vector2) {
	for i := 0; i < len(p.Points)-1; i++ {
		startPoint := raylib.Vector2Add(p.Points[i], offset)
		endPoint := raylib.Vector2Add(p.Points[i+1], offset)

		raylib.DrawLineEx(startPoint, endPoint, p.Size, p.Color)
		raylib.DrawCircle(int32(startPoint.X), int32(startPoint.Y), p.Size/2, p.Color)
	}
}
