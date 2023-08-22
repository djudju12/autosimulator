package graphics

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

type VectorInt struct {
	X, Y float64
}

func PointsFromRadius(p1, p2 sdl.Point, r1 float64) (*sdl.Point, *sdl.Point) {
	// u = v / ||v|| (normalization)
	vec1 := FromSdlPoint(p1)
	vec2 := FromSdlPoint(p2)
	u := &VectorInt{X: vec2.X - vec1.X, Y: vec2.Y - vec1.Y}
	u.Normalize()
	// v + d*u
	distanceU := MultiplyByScalar(u, r1)
	return Add(vec1, distanceU).ToSdlPoint(), Sub(vec2, distanceU).ToSdlPoint()
}

// TODO: Mudar!
func (v *VectorInt) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}
func (v *VectorInt) Normalize() {
	length := v.Length()
	v.X = v.X / length
	v.Y = v.Y / length
}

func MultiplyByScalar(v *VectorInt, scalar float64) *VectorInt {
	return &VectorInt{X: v.X * scalar, Y: v.Y * scalar}
}

func Sub(vec1, vec2 *VectorInt) *VectorInt {
	return &VectorInt{X: vec1.X - vec2.X, Y: vec1.Y - vec2.Y}
}

func Add(vec1, vec2 *VectorInt) *VectorInt {
	return &VectorInt{X: vec1.X + vec2.X, Y: vec1.Y + vec2.Y}
}

func (v *VectorInt) ToSdlPoint() *sdl.Point {
	return &sdl.Point{X: int32(v.X), Y: int32(v.Y)}
}

func FromSdlPoint(point sdl.Point) *VectorInt {
	return &VectorInt{X: float64(point.X), Y: float64(point.Y)}
}
