package main

import (
	"math"
	"sync"
	"syscall/js"
)

const fov = math.Pi / 2

var (
	width  int
	height int
)

func main() {
	canvas := js.Global().Get("document").Call("getElementById", "main")
	ctx := canvas.Call("getContext", "bitmaprenderer")

	width = canvas.Get("clientWidth").Int()
	height = canvas.Get("clientHeight").Int()

	canvas.Set("width", width)
	canvas.Set("height", height)

	ivory := Material{
		RefractiveIndex:  1.0,
		Albedo:           Vec4f{0.6, 0.3, 0.1, 0.0},
		DiffuseColor:     Vec3f{0.4, 0.4, 0.3},
		SpecularExponent: 50.0,
	}

	glass := Material{
		RefractiveIndex:  1.5,
		Albedo:           Vec4f{0.0, 0.5, 0.1, 0.8},
		DiffuseColor:     Vec3f{0.6, 0.7, 0.8},
		SpecularExponent: 125.0,
	}

	redRubber := Material{
		RefractiveIndex:  1.0,
		Albedo:           Vec4f{0.9, 0.1, 0.0, 0.0},
		DiffuseColor:     Vec3f{0.3, 0.1, 0.1},
		SpecularExponent: 10.0,
	}

	mirror := Material{
		RefractiveIndex:  1.0,
		Albedo:           Vec4f{0.0, 10.0, 0.8, 0.0},
		DiffuseColor:     Vec3f{1.0, 1.0, 1.0},
		SpecularExponent: 1425.0,
	}

	spheres := []Sphere{
		{Vec3f{-3, 0, -16}, 2, ivory},
		{Vec3f{-1, -1.5, -12}, 2, glass},
		{Vec3f{1.5, -0.5, -18}, 3, redRubber},
		{Vec3f{7, 5, -18}, 4, mirror},
	}

	lights := []Light{
		{Vec3f{-20, 20, 20}, 1.5},
		{Vec3f{30, 50, -25}, 1.8},
		{Vec3f{30, 20, 30}, 1.7},
	}

	framebuffer := render(spheres, lights)

	jsbuffer := js.Global().Get("Uint8ClampedArray").New(len(framebuffer.Data()))
	js.CopyBytesToJS(jsbuffer, framebuffer.Data())
	imageData := js.Global().Get("ImageData").New(
		jsbuffer,
		width,
		height,
	)

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
			hit = orig.Add(dir.Mul(dist))
			N = hit.Sub(sphere.Center).Normalize()
			material = sphere.Material
		}
	}
	intersect = spheresDist < 1000
	return
}

func calcOrig(dir, point, N Vec3f) Vec3f {
	if dir.Dot(N) < 0 {
		return point.Sub(N.Mul(1e-3))
	}
	return point.Add(N.Mul(1e-3))
}

func castRay(orig, dir Vec3f, spheres []Sphere, lights []Light, depth int) Vec3f {
	point, N, material, intersect := sceneIntersect(orig, dir, spheres)
	if depth > 4 || !intersect {
		return Vec3f{0.2, 0.7, 0.8} // background color
	}

	reflectDir := dir.Reflect(N)
	refractDir := dir.Refract(N, material.RefractiveIndex, 1).Normalize()

	reflectOrig := calcOrig(reflectDir, point, N)
	refractOrig := calcOrig(refractDir, point, N)

	reflectColor := castRay(reflectOrig, reflectDir, spheres, lights, depth+1)
	refractColor := castRay(refractOrig, refractDir, spheres, lights, depth+1)

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

func render(spheres []Sphere, lights []Light) Bitmap {
	framebuffer := NewBitmap(width, height)

	wg := sync.WaitGroup{}
	wg.Add(height)
	for j := 0; j < height; j++ {
		go func(j int) {
			for i := 0; i < width; i++ {
				x := (2*(float64(i)+0.5)/float64(width) - 1) * math.Tan(fov/2) * float64(width) / float64(height)
				y := -(2*(float64(j)+0.5)/float64(height) - 1) * math.Tan(fov/2)
				dir := (Vec3f{x, y, -1}).Normalize()
				framebuffer.SetPixel(i, j, castRay(Vec3f{0, 0, 0}, dir, spheres, lights, 0))

			}
			wg.Done()
		}(j)
	}

	wg.Wait()

	return *framebuffer
}
