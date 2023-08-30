package graphics

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

type Vector struct {
	X, Y float64
}

func LinePoints(p1, p2 sdl.Point, distanceFromP1, distanceFromP2 float64) (*sdl.Point, *sdl.Point, bool) {
	// TODO: Test if this convrsion is necessary
	if p1 == p2 {
		return nil, nil, false
	}

	vec1 := FromSdlPoint(p1)
	vec2 := FromSdlPoint(p2)

	// u = v / ||v|| (normalization)
	u := &Vector{X: vec2.X - vec1.X, Y: vec2.Y - vec1.Y}
	u.Normalize()

	// v + d*u
	p1U := MultiplyByScalar(u, distanceFromP1)
	p2U := MultiplyByScalar(u, distanceFromP2)
	return Add(vec1, *p1U).ToSdlPoint(), Sub(vec2, *p2U).ToSdlPoint(), true
}

// For the simplicity of the code, we will use sqrt()
func (v *Vector) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func MultiplyByScalar(v *Vector, scalar float64) *Vector {
	return &Vector{X: v.X * scalar, Y: v.Y * scalar}
}

func Sub(vec1, vec2 Vector) *Vector {
	return &Vector{X: vec1.X - vec2.X, Y: vec1.Y - vec2.Y}
}

func Add(vec1, vec2 Vector) *Vector {
	return &Vector{X: vec1.X + vec2.X, Y: vec1.Y + vec2.Y}
}

func (v *Vector) ToSdlPoint() *sdl.Point {
	return &sdl.Point{X: int32(v.X), Y: int32(v.Y)}
}

func FromSdlPoint(point sdl.Point) Vector {
	return Vector{X: float64(point.X), Y: float64(point.Y)}
}

func (v *Vector) Normalize() {
	length := v.Length()
	v.X = v.X / length
	v.Y = v.Y / length
}
