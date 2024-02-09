package main

import (
	raylib "github.com/gen2brain/raylib-go/raylib"
	"math"
)

type rectangleTool struct {
	StartPos raylib.Vector2
	EndPos   raylib.Vector2
	Rounded  bool
	Color    raylib.Color
	Size     float32
}

func (r *rectangleTool) Render() raylib.Texture2D {
	offset := raylib.Vector2Scale(canvas.Offset, -1)
	texture := raylib.LoadRenderTexture(int32(canvas.Size.X), int32(canvas.Size.Y))
	startPosOffset := raylib.Vector2Add(r.StartPos, offset)
	endPosOffset := raylib.Vector2Add(r.EndPos, offset)

	raylib.BeginTextureMode(texture)
	raylib.ClearBackground(raylib.Fade(raylib.Black, 0))

	if r.Rounded {
		// Linesss
		raylib.DrawLineEx(startPosOffset, raylib.NewVector2(endPosOffset.X, startPosOffset.Y), r.Size, r.Color)
		raylib.DrawLineEx(startPosOffset, raylib.NewVector2(startPosOffset.X, endPosOffset.Y), r.Size, r.Color)
		raylib.DrawLineEx(endPosOffset, raylib.NewVector2(endPosOffset.X, startPosOffset.Y), r.Size, r.Color)
		raylib.DrawLineEx(endPosOffset, raylib.NewVector2(startPosOffset.X, endPosOffset.Y), r.Size, r.Color)

		// Roundy
		raylib.DrawCircle(int32(startPosOffset.X), int32(startPosOffset.Y), r.Size/2, r.Color)
		raylib.DrawCircle(int32(endPosOffset.X), int32(startPosOffset.Y), r.Size/2, r.Color)
		raylib.DrawCircle(int32(startPosOffset.X), int32(endPosOffset.Y), r.Size/2, r.Color)
		raylib.DrawCircle(int32(endPosOffset.X), int32(endPosOffset.Y), r.Size/2, r.Color)
	} else {
		startPos := raylib.NewVector2(
			float32(math.Min(float64(startPosOffset.X), float64(endPosOffset.X))),
			float32(math.Min(float64(startPosOffset.Y), float64(endPosOffset.Y))),
		)
		endPos := raylib.NewVector2(
			float32(math.Max(float64(startPosOffset.X), float64(endPosOffset.X))),
			float32(math.Max(float64(startPosOffset.Y), float64(endPosOffset.Y))),
		)

		raylib.DrawLineEx(raylib.NewVector2(startPos.X-(r.Size/2), startPos.Y), raylib.NewVector2(endPos.X+(r.Size/2), startPos.Y), r.Size, r.Color)
		raylib.DrawLineEx(raylib.NewVector2(endPos.X, startPos.Y-(r.Size/2)), raylib.NewVector2(endPos.X, endPos.Y+(r.Size/2)), r.Size, r.Color)
		raylib.DrawLineEx(raylib.NewVector2(startPos.X-(r.Size/2), endPos.Y), raylib.NewVector2(endPos.X+(r.Size/2), endPos.Y), r.Size, r.Color)
		raylib.DrawLineEx(raylib.NewVector2(startPos.X, startPos.Y-(r.Size/2)), raylib.NewVector2(startPos.X, endPos.Y+(r.Size/2)), r.Size, r.Color)
	}

	raylib.EndTextureMode()

	return texture.Texture
}

func (r *rectangleTool) Draw() {
	if r.Rounded {
		// Linesss
		raylib.DrawLineEx(r.StartPos, raylib.NewVector2(r.EndPos.X, r.StartPos.Y), r.Size, r.Color)
		raylib.DrawLineEx(r.StartPos, raylib.NewVector2(r.StartPos.X, r.EndPos.Y), r.Size, r.Color)
		raylib.DrawLineEx(r.EndPos, raylib.NewVector2(r.EndPos.X, r.StartPos.Y), r.Size, r.Color)
		raylib.DrawLineEx(r.EndPos, raylib.NewVector2(r.StartPos.X, r.EndPos.Y), r.Size, r.Color)

		// Roundy
		raylib.DrawCircle(int32(r.StartPos.X), int32(r.StartPos.Y), r.Size/2, r.Color)
		raylib.DrawCircle(int32(r.EndPos.X), int32(r.StartPos.Y), r.Size/2, r.Color)
		raylib.DrawCircle(int32(r.StartPos.X), int32(r.EndPos.Y), r.Size/2, r.Color)
		raylib.DrawCircle(int32(r.EndPos.X), int32(r.EndPos.Y), r.Size/2, r.Color)
	} else {
		startPos := raylib.NewVector2(
			float32(math.Min(float64(r.StartPos.X), float64(r.EndPos.X))),
			float32(math.Min(float64(r.StartPos.Y), float64(r.EndPos.Y))),
		)
		endPos := raylib.NewVector2(
			float32(math.Max(float64(r.StartPos.X), float64(r.EndPos.X))),
			float32(math.Max(float64(r.StartPos.Y), float64(r.EndPos.Y))),
		)

		raylib.DrawLineEx(raylib.NewVector2(startPos.X-(r.Size/2), startPos.Y), raylib.NewVector2(endPos.X+(r.Size/2), startPos.Y), r.Size, r.Color)
		raylib.DrawLineEx(raylib.NewVector2(endPos.X, startPos.Y-(r.Size/2)), raylib.NewVector2(endPos.X, endPos.Y+(r.Size/2)), r.Size, r.Color)
		raylib.DrawLineEx(raylib.NewVector2(startPos.X-(r.Size/2), endPos.Y), raylib.NewVector2(endPos.X+(r.Size/2), endPos.Y), r.Size, r.Color)
		raylib.DrawLineEx(raylib.NewVector2(startPos.X, startPos.Y-(r.Size/2)), raylib.NewVector2(startPos.X, endPos.Y+(r.Size/2)), r.Size, r.Color)
	}
}
