// Package vector provides implementations and interfaces for 3D and 4D vectors
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
	ToVector4(t float64) Vector4
}

// Vector4 represents a four-dimensional vector (space-time)
type Vector4 interface {
	// Components
	X() float64
	Y() float64
	Z() float64
	T() float64

	// Vector operations
	Add(v Vector4) Vector4
	Sub(v Vector4) Vector4
	Scale(s float64) Vector4
	Dot(v Vector4) float64
	Length() float64
	LengthSquared() float64
	Normalize() Vector4

	// Space-time specific operations
	SpaceLength() float64
	IsTimelike() bool
	IsSpacelike() bool
	IsLightlike() bool
	ProperTime() float64

	// Conversion
	ToArray() [4]float64
	ToVector3() Vector3
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

// ToVector4 converts the 3D vector to a 4D vector with the specified time component
func (v *Vec3) ToVector4(t float64) Vector4 {
	return &Vec4{v.x, v.y, v.z, t}
}

// Vec4 implements the Vector4 interface
type Vec4 struct {
	x, y, z, t float64
}

// NewVector4 creates a new four-dimensional vector
func NewVector4(x, y, z, t float64) Vector4 {
	return &Vec4{x, y, z, t}
}

// Zero4 returns a zero four-dimensional vector
func Zero4() Vector4 {
	return &Vec4{0, 0, 0, 0}
}

// X returns the x component of the vector
func (v *Vec4) X() float64 {
	return v.x
}

// Y returns the y component of the vector
func (v *Vec4) Y() float64 {
	return v.y
}

// Z returns the z component of the vector
func (v *Vec4) Z() float64 {
	return v.z
}

// T returns the time component of the vector
func (v *Vec4) T() float64 {
	return v.t
}

// Add sums two vectors
func (v *Vec4) Add(other Vector4) Vector4 {
	return &Vec4{
		v.x + other.X(),
		v.y + other.Y(),
		v.z + other.Z(),
		v.t + other.T(),
	}
}

// Sub subtracts two vectors
func (v *Vec4) Sub(other Vector4) Vector4 {
	return &Vec4{
		v.x - other.X(),
		v.y - other.Y(),
		v.z - other.Z(),
		v.t - other.T(),
	}
}

// Scale multiplies a vector by a scalar
func (v *Vec4) Scale(s float64) Vector4 {
	return &Vec4{
		v.x * s,
		v.y * s,
		v.z * s,
		v.t * s,
	}
}

// Dot calculates the dot product between two vectors (with Minkowski metric)
func (v *Vec4) Dot(other Vector4) float64 {
	// We use the Minkowski metric (-,+,+,+)
	return -v.t*other.T() + v.x*other.X() + v.y*other.Y() + v.z*other.Z()
}

// LengthSquared calculates the squared length of the vector (with Minkowski metric)
func (v *Vec4) LengthSquared() float64 {
	return v.Dot(v)
}

// Length calculates the length of the vector (with Minkowski metric)
func (v *Vec4) Length() float64 {
	l2 := v.LengthSquared()
	if l2 < 0 {
		return -math.Sqrt(-l2)
	}
	return math.Sqrt(l2)
}

// Normalize normalizes the vector
func (v *Vec4) Normalize() Vector4 {
	length := v.Length()
	if math.Abs(length) < 1e-10 {
		return &Vec4{0, 0, 0, 0}
	}
	return v.Scale(1.0 / length)
}

// SpaceLength calculates the spatial length of the vector
func (v *Vec4) SpaceLength() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z)
}

// IsTimelike checks if the vector is timelike
func (v *Vec4) IsTimelike() bool {
	return v.LengthSquared() < 0
}

// IsSpacelike checks if the vector is spacelike
func (v *Vec4) IsSpacelike() bool {
	return v.LengthSquared() > 0
}

// IsLightlike checks if the vector is lightlike
func (v *Vec4) IsLightlike() bool {
	return math.Abs(v.LengthSquared()) < 1e-10
}

// ProperTime calculates the proper time associated with the vector
func (v *Vec4) ProperTime() float64 {
	if !v.IsTimelike() {
		return 0
	}
	return math.Abs(v.Length())
}

// ToArray converts the vector to an array
func (v *Vec4) ToArray() [4]float64 {
	return [4]float64{v.x, v.y, v.z, v.t}
}

// ToVector3 converts the 4D vector to a 3D vector
func (v *Vec4) ToVector3() Vector3 {
	return &Vec3{v.x, v.y, v.z}
}
