package main

import (
	"syscall/js"
)

var (
	width  int
	height int

	console js.Value
)

func main() {
	console = js.Global().Get("console")

	canvas := js.Global().Get("document").Call("getElementById", "main")
	ctx := canvas.Call("getContext", "bitmaprenderer")

	width = canvas.Get("clientWidth").Int()
	height = canvas.Get("clientHeight").Int()

	canvas.Set("width", width)
	canvas.Set("height", height)

	framebuffer := NewBitmap(width, height)

	render(framebuffer)

	jsbuffer := js.Global().Get("Uint8ClampedArray").New(len(framebuffer.Data()))
	js.CopyBytesToJS(jsbuffer, framebuffer.Data())
	imageData := js.Global().Get("ImageData").New(
		jsbuffer,
		width,
		height,
	)
	//console.Call("log", imageData)
	imageBitmap := js.Global().Call("createImageBitmap", imageData)

	done := make(chan js.Value)

	imageBitmap.Call("then", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		result := args[0]

		done <- result
		return nil
	}))

	ctx.Call("transferFromImageBitmap", <-done)
}

func ftob(f float64) byte {
	return byte(f * 255)
}

func render(framebuffer *Bitmap) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			framebuffer.Set(x, y, 0, ftob(float64(y)/float64(height)))
			framebuffer.Set(x, y, 1, ftob(float64(x)/float64(width)))
			framebuffer.Set(x, y, 3, 0xFF)
		}
	}
}
