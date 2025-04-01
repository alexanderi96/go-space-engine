// Package vector fornisce implementazioni e interfacce per vettori 3D e 4D
package vector

import (
	"math"
)

// Vector3 rappresenta un vettore tridimensionale
type Vector3 interface {
	// Componenti
	X() float64
	Y() float64
	Z() float64

	// Operazioni vettoriali
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

	// Conversione
	ToArray() [3]float64
	ToVector4(t float64) Vector4
}

// Vector4 rappresenta un vettore quadridimensionale (spazio-tempo)
type Vector4 interface {
	// Componenti
	X() float64
	Y() float64
	Z() float64
	T() float64

	// Operazioni vettoriali
	Add(v Vector4) Vector4
	Sub(v Vector4) Vector4
	Scale(s float64) Vector4
	Dot(v Vector4) float64
	Length() float64
	LengthSquared() float64
	Normalize() Vector4

	// Operazioni specifiche per lo spazio-tempo
	SpaceLength() float64
	IsTimelike() bool
	IsSpacelike() bool
	IsLightlike() bool
	ProperTime() float64

	// Conversione
	ToArray() [4]float64
	ToVector3() Vector3
}

// Vec3 implementa l'interfaccia Vector3
type Vec3 struct {
	x, y, z float64
}

// NewVector3 crea un nuovo vettore tridimensionale
func NewVector3(x, y, z float64) Vector3 {
	return &Vec3{x, y, z}
}

// Zero3 restituisce un vettore tridimensionale nullo
func Zero3() Vector3 {
	return &Vec3{0, 0, 0}
}

// X restituisce la componente x del vettore
func (v *Vec3) X() float64 {
	return v.x
}

// Y restituisce la componente y del vettore
func (v *Vec3) Y() float64 {
	return v.y
}

// Z restituisce la componente z del vettore
func (v *Vec3) Z() float64 {
	return v.z
}

// Add somma due vettori
func (v *Vec3) Add(other Vector3) Vector3 {
	return &Vec3{
		v.x + other.X(),
		v.y + other.Y(),
		v.z + other.Z(),
	}
}

// Sub sottrae due vettori
func (v *Vec3) Sub(other Vector3) Vector3 {
	return &Vec3{
		v.x - other.X(),
		v.y - other.Y(),
		v.z - other.Z(),
	}
}

// Scale moltiplica un vettore per uno scalare
func (v *Vec3) Scale(s float64) Vector3 {
	return &Vec3{
		v.x * s,
		v.y * s,
		v.z * s,
	}
}

// Dot calcola il prodotto scalare tra due vettori
func (v *Vec3) Dot(other Vector3) float64 {
	return v.x*other.X() + v.y*other.Y() + v.z*other.Z()
}

// Cross calcola il prodotto vettoriale tra due vettori
func (v *Vec3) Cross(other Vector3) Vector3 {
	return &Vec3{
		v.y*other.Z() - v.z*other.Y(),
		v.z*other.X() - v.x*other.Z(),
		v.x*other.Y() - v.y*other.X(),
	}
}

// LengthSquared calcola il quadrato della lunghezza del vettore
func (v *Vec3) LengthSquared() float64 {
	return v.x*v.x + v.y*v.y + v.z*v.z
}

// Length calcola la lunghezza del vettore
func (v *Vec3) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}

// Normalize normalizza il vettore
func (v *Vec3) Normalize() Vector3 {
	length := v.Length()
	if length < 1e-10 {
		return &Vec3{0, 0, 0}
	}
	return v.Scale(1.0 / length)
}

// DistanceSquared calcola il quadrato della distanza tra due vettori
func (v *Vec3) DistanceSquared(other Vector3) float64 {
	dx := v.x - other.X()
	dy := v.y - other.Y()
	dz := v.z - other.Z()
	return dx*dx + dy*dy + dz*dz
}

// Distance calcola la distanza tra due vettori
func (v *Vec3) Distance(other Vector3) float64 {
	return math.Sqrt(v.DistanceSquared(other))
}

// ToArray converte il vettore in un array
func (v *Vec3) ToArray() [3]float64 {
	return [3]float64{v.x, v.y, v.z}
}

// ToVector4 converte il vettore 3D in un vettore 4D con la componente temporale specificata
func (v *Vec3) ToVector4(t float64) Vector4 {
	return &Vec4{v.x, v.y, v.z, t}
}

// Vec4 implementa l'interfaccia Vector4
type Vec4 struct {
	x, y, z, t float64
}

// NewVector4 crea un nuovo vettore quadridimensionale
func NewVector4(x, y, z, t float64) Vector4 {
	return &Vec4{x, y, z, t}
}

// Zero4 restituisce un vettore quadridimensionale nullo
func Zero4() Vector4 {
	return &Vec4{0, 0, 0, 0}
}

// X restituisce la componente x del vettore
func (v *Vec4) X() float64 {
	return v.x
}

// Y restituisce la componente y del vettore
func (v *Vec4) Y() float64 {
	return v.y
}

// Z restituisce la componente z del vettore
func (v *Vec4) Z() float64 {
	return v.z
}

// T restituisce la componente temporale del vettore
func (v *Vec4) T() float64 {
	return v.t
}

// Add somma due vettori
func (v *Vec4) Add(other Vector4) Vector4 {
	return &Vec4{
		v.x + other.X(),
		v.y + other.Y(),
		v.z + other.Z(),
		v.t + other.T(),
	}
}

// Sub sottrae due vettori
func (v *Vec4) Sub(other Vector4) Vector4 {
	return &Vec4{
		v.x - other.X(),
		v.y - other.Y(),
		v.z - other.Z(),
		v.t - other.T(),
	}
}

// Scale moltiplica un vettore per uno scalare
func (v *Vec4) Scale(s float64) Vector4 {
	return &Vec4{
		v.x * s,
		v.y * s,
		v.z * s,
		v.t * s,
	}
}

// Dot calcola il prodotto scalare tra due vettori (con metrica di Minkowski)
func (v *Vec4) Dot(other Vector4) float64 {
	// Utilizziamo la metrica di Minkowski (-,+,+,+)
	return -v.t*other.T() + v.x*other.X() + v.y*other.Y() + v.z*other.Z()
}

// LengthSquared calcola il quadrato della lunghezza del vettore (con metrica di Minkowski)
func (v *Vec4) LengthSquared() float64 {
	return v.Dot(v)
}

// Length calcola la lunghezza del vettore (con metrica di Minkowski)
func (v *Vec4) Length() float64 {
	l2 := v.LengthSquared()
	if l2 < 0 {
		return -math.Sqrt(-l2)
	}
	return math.Sqrt(l2)
}

// Normalize normalizza il vettore
func (v *Vec4) Normalize() Vector4 {
	length := v.Length()
	if math.Abs(length) < 1e-10 {
		return &Vec4{0, 0, 0, 0}
	}
	return v.Scale(1.0 / length)
}

// SpaceLength calcola la lunghezza spaziale del vettore
func (v *Vec4) SpaceLength() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z)
}

// IsTimelike verifica se il vettore è di tipo tempo
func (v *Vec4) IsTimelike() bool {
	return v.LengthSquared() < 0
}

// IsSpacelike verifica se il vettore è di tipo spazio
func (v *Vec4) IsSpacelike() bool {
	return v.LengthSquared() > 0
}

// IsLightlike verifica se il vettore è di tipo luce
func (v *Vec4) IsLightlike() bool {
	return math.Abs(v.LengthSquared()) < 1e-10
}

// ProperTime calcola il tempo proprio associato al vettore
func (v *Vec4) ProperTime() float64 {
	if !v.IsTimelike() {
		return 0
	}
	return math.Abs(v.Length())
}

// ToArray converte il vettore in un array
func (v *Vec4) ToArray() [4]float64 {
	return [4]float64{v.x, v.y, v.z, v.t}
}

// ToVector3 converte il vettore 4D in un vettore 3D
func (v *Vec4) ToVector3() Vector3 {
	return &Vec3{v.x, v.y, v.z}
}
