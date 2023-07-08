package main

import (
	"math"
	"sync"
	"syscall/js"

	"github.com/sportshead/tinyraytracer/bitmap"
	. "github.com/sportshead/tinyraytracer/vectors"
)

const fov = math.Pi / 2

var (
	width  int
	height int
)

var ctx js.Value

func main() {
	done := make(chan struct{})

	js.Global().Call("fetch", "envmap.jpg").Call("then", js.FuncOf(func(this js.Value, args []js.Value) any {
		res := args[0]
		return res.Call("blob")
	})).Call("then", js.FuncOf(func(this js.Value, args []js.Value) any {
		blob := args[0]
		return js.Global().Call("createImageBitmap", blob)
	})).Call("then", js.FuncOf(func(this js.Value, args []js.Value) any {
		jsBitmap := args[0]

		offscreenCanvas := js.Global().Get("OffscreenCanvas").New(jsBitmap.Get("width"), jsBitmap.Get("height"))
		context := offscreenCanvas.Call("getContext", "2d")
		context.Call("drawImage", jsBitmap, 0, 0)

		return context.Call("getImageData", 0, 0, jsBitmap.Get("width"), jsBitmap.Get("height"))
	})).Call("then", js.FuncOf(start)).Call("then", js.FuncOf(func(this js.Value, args []js.Value) any {
		imageBitmap := args[0]
		ctx.Call("transferFromImageBitmap", imageBitmap)
		close(done)
		return nil
	}))
	<-done
}

// start(envmap: ImageData): Promise<ImageBitmap>
func start(this js.Value, args []js.Value) any {
	envmapImageData := args[0]
	envmapJsBuffer := envmapImageData.Get("data")
	envmapBuffer := make([]byte, envmapJsBuffer.Length())

	js.CopyBytesToGo(envmapBuffer, envmapJsBuffer)

	envmap := bitmap.Bitmap{
		Data:   envmapBuffer,
		Width:  envmapImageData.Get("width").Int(),
		Height: envmapImageData.Get("height").Int(),
	}

	canvas := js.Global().Get("document").Call("getElementById", "main")
	ctx = canvas.Call("getContext", "bitmaprenderer")

	width = canvas.Get("clientWidth").Int()
	height = canvas.Get("clientHeight").Int()

	canvas.Set("width", width)
	canvas.Set("height", height)

	spheres := []Sphere{
		{Vec3f{-3, 0, -16}, 2, Materials[Ivory]},
		{Vec3f{-1, -1.5, -12}, 2, Materials[Glass]},
		{Vec3f{1.5, -0.5, -18}, 3, Materials[RedRubber]},
		{Vec3f{7, 5, -18}, 4, Materials[Mirror]},
	}

	lights := []Light{
		{Vec3f{-20, 20, 20}, 1.5},
		{Vec3f{30, 50, -25}, 1.8},
		{Vec3f{30, 20, 30}, 1.7},
	}

	frameBuffer := render(spheres, lights, envmap)

	jsBuffer := js.Global().Get("Uint8ClampedArray").New(len(frameBuffer.Data))
	js.CopyBytesToJS(jsBuffer, frameBuffer.Data)
	imageData := js.Global().Get("ImageData").New(
		jsBuffer,
		width,
		height,
	)

	imageBitmap := js.Global().Call("createImageBitmap", imageData)
	return imageBitmap
}

func sceneIntersect(orig, dir Vec3f, spheres []Sphere) (pt Vec3f, N Vec3f, material Material, intersect bool) {
	material = Material{1, Vec4f{1, 0, 0, 0}, Vec3f{0, 0, 0}, 0}

	nearestDist := 1e10
	if math.Abs(dir[1]) > .001 { // intersect the ray with the checkerboard, avoid division by zero
		d := -(orig[1] + 4) / dir[1] // the checkerboard plane has equation y = -4
		p := orig.Add(dir.Mul(d))
		if d > .001 && d < nearestDist && math.Abs(p[0]) < 10 && p[2] < -10 && p[2] > -30 {
			nearestDist = d
			pt = p
			N = Vec3f{0, 1, 0}
			if (int(0.5*pt[0]+1000)+int(0.5*pt[2]))&1 == 1 {
				material.DiffuseColor = Vec3f{0.3, 0.3, 0.3}
			} else {
				material.DiffuseColor = Vec3f{0.3, 0.2, 0.1}
			}
		}
	}

	for _, sphere := range spheres {
		dist, rayIntersect := sphere.RayIntersect(orig, dir)
		if !rayIntersect || dist > nearestDist {
			continue
		}
		nearestDist = dist
		pt = orig.Add(dir.Mul(dist))
		N = pt.Sub(sphere.Center).Normalize()
		material = sphere.Material
	}

	intersect = nearestDist < 1000
	return
}

func calcOrig(dir, point, N Vec3f) Vec3f {
	if dir.Dot(N) < 0 {
		return point.Sub(N.Mul(1e-3))
	}
	return point.Add(N.Mul(1e-3))
}

func castRay(orig, dir Vec3f, spheres []Sphere, lights []Light, depth int, envmap bitmap.Bitmap) Vec3f {
	point, N, material, intersect := sceneIntersect(orig, dir, spheres)
	if depth > 4 || !intersect {
		normalized := dir.Normalize()
		// https://github.com/tpaschalis/go-tinyraytracer/blob/b95d3b3be7a231154881944dd6953955b9af2ae4/main.go#L107
		pixel, err := envmap.GetPixel(int((normalized[0]/2+0.5)*0.5*float64(envmap.Width)), int((1-(normalized[1]/2+0.5))*float64(envmap.Height)))
		if err != nil {
			panic(err)
		}
		return pixel
	}

	reflectDir := dir.Reflect(N)
	refractDir := dir.Refract(N, material.RefractiveIndex, 1).Normalize()

	reflectOrig := calcOrig(reflectDir, point, N)
	refractOrig := calcOrig(refractDir, point, N)

	reflectColor := castRay(reflectOrig, reflectDir, spheres, lights, depth+1, envmap)
	refractColor := castRay(refractOrig, refractDir, spheres, lights, depth+1, envmap)

	diffuseLightIntensity := 0.0
	specularLightIntensity := 0.0
	for _, light := range lights {
		lightDir := light.Position.Sub(point).Normalize()
		lightDistance := light.Position.Sub(point).Norm()

		shadowOrig := calcOrig(lightDir, point, N)

		shadowPt, _, _, shadowIntersect := sceneIntersect(shadowOrig, lightDir, spheres)
		if shadowIntersect && shadowPt.Sub(shadowOrig).Norm() < lightDistance {
			continue
		}

		diffuseLightIntensity += light.Intensity * math.Max(0, lightDir.Dot(N))
		specularLightIntensity += light.Intensity * math.Pow(math.Max(0, lightDir.Reflect(N).Dot(dir)), material.SpecularExponent)
	}
	return material.DiffuseColor.Mul(diffuseLightIntensity * material.Albedo[0]).Add(Vec3f{1.0, 1.0, 1.0}.Mul(specularLightIntensity * material.Albedo[1])).Add(reflectColor.Mul(material.Albedo[2])).Add(refractColor.Mul(material.Albedo[3]))
}

func render(spheres []Sphere, lights []Light, envmap bitmap.Bitmap) bitmap.Bitmap {
	frameBuffer := bitmap.NewBitmap(width, height)

	wg := sync.WaitGroup{}
	wg.Add(height)
	for j := 0; j < height; j++ {
		go func(j int) {
			for i := 0; i < width; i++ {
				x := (2*(float64(i)+0.5)/float64(width) - 1) * math.Tan(fov/2) * float64(width) / float64(height)
				y := -(2*(float64(j)+0.5)/float64(height) - 1) * math.Tan(fov/2)
				dir := (Vec3f{x, y, -1}).Normalize()
				frameBuffer.SetPixel(i, j, castRay(Vec3f{0, 0, 0}, dir, spheres, lights, 0, envmap))
			}
			wg.Done()
		}(j)
	}

	wg.Wait()

	return *frameBuffer
}
