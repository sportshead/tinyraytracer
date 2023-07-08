package main

import (
	"math"

	. "github.com/sportshead/tinyraytracer/vectors"
)

type Sphere struct {
	Center   Vec3f
	Radius   float64
	Material Material
}

func (s *Sphere) RayIntersect(orig, dir Vec3f) (float64, bool) {
	L := s.Center.Sub(orig)
	tca := L.Dot(dir)
	d2 := L.Dot(L) - tca*tca
	if d2 > s.Radius*s.Radius {
		return 0, false
	}
	thc := math.Sqrt(s.Radius*s.Radius - d2)
	t0 := tca - thc
	t1 := tca + thc
	if t0 > 0.001 {
		return t0, true
	}
	if t1 > 0.001 {
		return t1, true
	}
	return 0, false
}
