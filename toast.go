package main

import (
	"fmt"
	"time"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

const (
	toastMaxAge = 1 * time.Second
)

var (
	toasts      = []toast{}
	toastHeight = float32(0)
)

type toast struct {
	Text string
	Age  time.Time
}

func AddToast(text string) {
	toast := toast{Text: text, Age: time.Now()}
	toasts = append(toasts, toast)

	fmt.Printf("Added toast: '%s'\n", text)
}

func UpdateToasts() {
	if len(toasts) != 0 {
		toastHeight = raylib.Lerp(toastHeight, float32(20*len(toasts))+10, 0.1)
	} else {
		toastHeight = raylib.Lerp(toastHeight, 0, 0.1)
	}

	for i := 0; i < len(toasts); i += 1 {
		if time.Since(toasts[i].Age) > toastMaxAge {
			toasts = append(toasts[:i], toasts[i+1:]...)
			i -= 1
		}
	}
}

func DrawToasts() {
	raylib.BeginScissorMode(0, 0, WindowWidth, int32(toastHeight))
	raylib.DrawRectangle(0, 0, WindowWidth, WindowHeight, raylib.Fade(raylib.Black, 0.5))
	for i := 0; i < len(toasts); i++ {
		//text := fmt.Sprintf("%s (%s)", toasts[i].Text, time.Since(toasts[i].Age).Round(time.Second))
		text := toasts[i].Text
		raylib.DrawText(text, 10, int32(20*i)+10, 10, raylib.White)
	}
	raylib.EndScissorMode()
}
