package main

import "math"

type Sphere struct {
	center Vec3f
	radius float64
}

// math i dont understand
func (s *Sphere) RayIntersect(orig, dir Vec3f) (float64, bool) {
	L := s.center.Sub(orig)
	tca := L.Dot(dir)
	d2 := L.Dot(L) - tca*tca
	rSq := s.radius * s.radius
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
