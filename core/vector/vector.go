// Package vector provides implementations and interfaces for 3D vectors
package vector

import (
	"math"
)

// Vector3 represents a three-dimensional vector
type Vector3 interface {
	// Components
	X() float64
	Y() float64
	Z() float64

	// Vector operations
	Add(v Vector3) Vector3
	Sub(v Vector3) Vector3
	Scale(s float64) Vector3
	Dot(v Vector3) float64
	Cross(v Vector3) Vector3
	Length() float64
	LengthSquared() float64
	Normalize() Vector3
	Distance(v Vector3) float64
	DistanceSquared(v Vector3) float64

	// Conversion
	ToArray() [3]float64
}

// Vec3 implements the Vector3 interface
type Vec3 struct {
	x, y, z float64
}

// NewVector3 creates a new three-dimensional vector
func NewVector3(x, y, z float64) Vector3 {
	return &Vec3{x, y, z}
}

// Zero3 returns a zero three-dimensional vector
func Zero3() Vector3 {
	return &Vec3{0, 0, 0}
}

// X returns the x component of the vector
func (v *Vec3) X() float64 {
	return v.x
}

// Y returns the y component of the vector
func (v *Vec3) Y() float64 {
	return v.y
}

// Z returns the z component of the vector
func (v *Vec3) Z() float64 {
	return v.z
}

// Add sums two vectors
func (v *Vec3) Add(other Vector3) Vector3 {
	return &Vec3{
		v.x + other.X(),
		v.y + other.Y(),
		v.z + other.Z(),
	}
}

// Sub subtracts two vectors
func (v *Vec3) Sub(other Vector3) Vector3 {
	return &Vec3{
		v.x - other.X(),
		v.y - other.Y(),
		v.z - other.Z(),
	}
}

// Scale multiplies a vector by a scalar
func (v *Vec3) Scale(s float64) Vector3 {
	return &Vec3{
		v.x * s,
		v.y * s,
		v.z * s,
	}
}

// Dot calculates the dot product between two vectors
func (v *Vec3) Dot(other Vector3) float64 {
	return v.x*other.X() + v.y*other.Y() + v.z*other.Z()
}

// Cross calculates the cross product between two vectors
func (v *Vec3) Cross(other Vector3) Vector3 {
	return &Vec3{
		v.y*other.Z() - v.z*other.Y(),
		v.z*other.X() - v.x*other.Z(),
		v.x*other.Y() - v.y*other.X(),
	}
}

// LengthSquared calculates the squared length of the vector
func (v *Vec3) LengthSquared() float64 {
	return v.x*v.x + v.y*v.y + v.z*v.z
}

// Length calculates the length of the vector
func (v *Vec3) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}

// Normalize normalizes the vector
func (v *Vec3) Normalize() Vector3 {
	length := v.Length()
	if length < 1e-10 {
		return &Vec3{0, 0, 0}
	}
	return v.Scale(1.0 / length)
}

// DistanceSquared calculates the squared distance between two vectors
func (v *Vec3) DistanceSquared(other Vector3) float64 {
	dx := v.x - other.X()
	dy := v.y - other.Y()
	dz := v.z - other.Z()
	return dx*dx + dy*dy + dz*dz
}

// Distance calculates the distance between two vectors
func (v *Vec3) Distance(other Vector3) float64 {
	return math.Sqrt(v.DistanceSquared(other))
}

// ToArray converts the vector to an array
func (v *Vec3) ToArray() [3]float64 {
	return [3]float64{v.x, v.y, v.z}
}
