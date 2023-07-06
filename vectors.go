package main

import (
	"fmt"
	"math"
)

type Vec3f [3]float64

func (v Vec3f) Add(other Vec3f) Vec3f {
	return Vec3f{v[0] + other[0], v[1] + other[1], v[2] + other[2]}
}

func (v Vec3f) Sub(other Vec3f) Vec3f {
	return Vec3f{v[0] - other[0], v[1] - other[1], v[2] - other[2]}
}

func (v Vec3f) Mul(f float64) Vec3f {
	return Vec3f{v[0] * f, v[1] * f, v[2] * f}
}

func (v Vec3f) Dot(other Vec3f) float64 {
	return v[0]*other[0] + v[1]*other[1] + v[2]*other[2]
}

func (v Vec3f) Cross(other Vec3f) Vec3f {
	return Vec3f{v[1]*other[2] - v[2]*other[1], v[2]*other[0] - v[0]*other[2], v[0]*other[1] - v[1]*other[0]}
}

func (v Vec3f) Norm() float64 {
	return math.Sqrt(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])
}

func (v Vec3f) Normalize() Vec3f {
	l := v.Norm()
	return Vec3f{v[0] / l, v[1] / l, v[2] / l}
}

func (v Vec3f) String() string {
	return fmt.Sprintf("(%f, %f, %f)", v[0], v[1], v[2])
}
