package main

import "math"

type Sphere struct {
	Center   Vec3f
	Radius   float64
	Material Material
}

func (s *Sphere) RayIntersect(orig, dir Vec3f) (float64, bool) {
	L := s.Center.Sub(orig)
	tca := L.Dot(dir)
	d2 := L.Dot(L) - tca*tca
	rSq := s.Radius * s.Radius
	if d2 > rSq {
		return 0, false
	}
	thc := math.Sqrt(rSq - d2)
	t0 := tca - thc
	t1 := tca + thc
	if t0 < 0 {
		t0 = t1
	}
	if t0 < 0 {
		return t0, false
	}
	return t0, true
}
