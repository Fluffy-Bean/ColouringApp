package application

import (
	"fmt"
	"time"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

var (
	toasts      = []Toast{}
	toastHeight = float32(0)
	toastMaxAge = 5 * time.Second
)

type Toast struct {
	Text string
	Age  time.Time
}

func AddToast(text string) {
	toast := Toast{Text: text, Age: time.Now()}
	toasts = append(toasts, toast)

	fmt.Printf("Added Toast: '%s'\n", text)
}

func UpdateToasts() {
	if len(toasts) != 0 {
		toastHeight = raylib.Lerp(toastHeight, float32(20*len(toasts))+10, 0.1)
	} else {
		toastHeight = raylib.Lerp(toastHeight, 0, 0.1)
	}

	for i := 0; i < len(toasts); i++ {
		if time.Since(toasts[i].Age) > toastMaxAge {
			toasts = append(toasts[:i], toasts[i+1:]...)
			i--
		}
	}
}

func DrawToasts() {
	raylib.BeginScissorMode(0, 0, WindowWidth, int32(toastHeight))
	raylib.DrawRectangle(0, 0, WindowWidth, WindowHeight, raylib.Fade(raylib.Black, 0.5))
	for i := 0; i < len(toasts); i++ {
		text := fmt.Sprintf("%s (%s)", toasts[i].Text, time.Since(toasts[i].Age).Round(time.Second))
		raylib.DrawText(text, 10, int32(20*i)+10, 10, raylib.White)
	}
	raylib.EndScissorMode()
}
