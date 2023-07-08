package main

type Material struct {
	RefractiveIndex  float64
	Albedo           Vec4f
	DiffuseColor     Vec3f
	SpecularExponent float64
}

type MaterialType int

const (
	Ivory MaterialType = iota
	Glass
	RedRubber
	Mirror
)

var (
	Materials = map[MaterialType]Material{
		Ivory: {
			RefractiveIndex:  1.0,
			Albedo:           Vec4f{0.6, 0.3, 0.1, 0.0},
			DiffuseColor:     Vec3f{0.4, 0.4, 0.3},
			SpecularExponent: 50.0,
		},
		Glass: {
			RefractiveIndex:  1.5,
			Albedo:           Vec4f{0.0, 0.5, 0.1, 0.8},
			DiffuseColor:     Vec3f{0.6, 0.7, 0.8},
			SpecularExponent: 125.0,
		},
		RedRubber: {
			RefractiveIndex:  1.0,
			Albedo:           Vec4f{0.9, 0.1, 0.0, 0.0},
			DiffuseColor:     Vec3f{0.3, 0.1, 0.1},
			SpecularExponent: 10.0,
		},
		Mirror: {
			RefractiveIndex:  1.0,
			Albedo:           Vec4f{0.0, 10.0, 0.8, 0.0},
			DiffuseColor:     Vec3f{1.0, 1.0, 1.0},
			SpecularExponent: 1425.0,
		},
	}
)
