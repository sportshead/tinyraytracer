package main

import (
	"math"
	"syscall/js"
)

const fov = math.Pi / 2

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

	ivory := Material{Vec3f{0.4, 0.4, 0.3}}
	redRubber := Material{Vec3f{0.3, 0.1, 0.1}}

	spheres := []Sphere{
		{Vec3f{-3, 0, -16}, 2, ivory},
		{Vec3f{-1, -1.5, -12}, 2, redRubber},
		{Vec3f{1.5, -0.5, -18}, 3, redRubber},
		{Vec3f{7, 5, -18}, 4, ivory},
	}

	framebuffer := render(spheres)

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
	return byte(math.Max(0, math.Min(255, f*255+0.5)))
}

func sceneIntersect(orig, dir Vec3f, spheres []Sphere) (hit Vec3f, N Vec3f, material Material, intersect bool) {
	spheresDist := math.MaxFloat64
	for _, sphere := range spheres {
		dist, rayIntersect := sphere.RayIntersect(orig, dir)
		if rayIntersect && dist < spheresDist {
			spheresDist = dist
			dd := dir.Dot(Vec3f{dist, dist, dist})
			hit = orig.Add(Vec3f{dd, dd, dd})
			N = hit.Sub(sphere.Center).Normalize()
			material = sphere.Material
		}
	}
	intersect = spheresDist < 1000
	return
}

func castRay(orig, dir Vec3f, spheres []Sphere) Vec3f {
	_, _, material, intersect := sceneIntersect(orig, dir, spheres)
	if !intersect {
		return Vec3f{0.2, 0.7, 0.8} // background color
	}
	return material.DiffuseColor
}

func render(spheres []Sphere) Bitmap {
	framebuffer := NewBitmap(width, height)

	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			x := (2*(float64(i)+0.5)/float64(width) - 1) * math.Tan(fov/2) * float64(width) / float64(height)
			y := -(2*(float64(j)+0.5)/float64(height) - 1) * math.Tan(fov/2)
			dir := (Vec3f{x, y, -1}).Normalize()
			framebuffer.SetPixel(i, j, castRay(Vec3f{0, 0, 0}, dir, spheres))
		}
	}

	return *framebuffer
}
