package main

import (
	"time"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

type toast struct {
	Text   string
	Age    time.Time
	MaxAge time.Duration
}

var (
	toasts            []toast
	toastShouldUpdate = true
	toastDimHeight    = float32(0)
)

func addToast(text string) {
	t := toast{Text: text, Age: time.Now(), MaxAge: 1 * time.Second}
	toasts = append(toasts, t)
}

func updateToasts() {
	if len(toasts) != 0 {
		toastDimHeight = raylib.Lerp(toastDimHeight, float32(20*len(toasts))+10, 0.1)
	} else {
		toastDimHeight = raylib.Lerp(toastDimHeight, 0, 0.1)
	}

	var t []toast
	for i := range toasts {
		if time.Since(toasts[i].Age) < toasts[i].MaxAge {
			t = append(t, toasts[i])
		}
	}
	toasts = t

	toastShouldUpdate = int(toastDimHeight) != 0
}

func drawToasts() {
	raylib.BeginScissorMode(0, 0, applicationWindowWidth, int32(toastDimHeight))
	{
		raylib.DrawRectangle(0, 0, applicationWindowWidth, applicationWindowHeight, raylib.Fade(raylib.Black, 0.5))
		for i, t := range toasts {
			raylib.DrawText(t.Text, 10, int32(20*i)+10, 10, raylib.White)
		}
	}
	raylib.EndScissorMode()
}
