// Package force provides interfaces and implementations for physical forces
package force

import (
	"math"

	"github.com/alexanderi96/go-space-engine/core/constants"
	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/physics/body"
)

// Force represents a physical force
type Force interface {
	// Apply applies the force to a body and returns the force vector
	Apply(b body.Body) vector.Vector3

	// ApplyBetween applies the force between two bodies and returns the force vectors for each body
	ApplyBetween(a, b body.Body) (vector.Vector3, vector.Vector3)

	// IsGlobal returns true if the force is global (applied to all bodies)
	IsGlobal() bool
}

// GravitationalForce implements gravitational force
type GravitationalForce struct {
	G     float64 // Gravitational constant
	Theta float64 // Approximation parameter for the Barnes-Hut algorithm
}

// NewGravitationalForce creates a new gravitational force
func NewGravitationalForce() *GravitationalForce {
	return &GravitationalForce{
		G:     constants.G,
		Theta: 0.5, // Default value that balances precision and efficiency
	}
}

// Apply applies the gravitational force to a body (does nothing for a single body)
func (gf *GravitationalForce) Apply(b body.Body) vector.Vector3 {
	// Gravity requires two bodies to be applied
	return vector.Zero3()
}

// ApplyBetween applies the gravitational force between two bodies
func (gf *GravitationalForce) ApplyBetween(a, b body.Body) (vector.Vector3, vector.Vector3) {
	// Calculate the direction vector from a to b
	direction := b.Position().Sub(a.Position())

	// Calculate the squared distance
	distanceSquared := direction.LengthSquared()

	// Avoid division by zero or too large forces
	if distanceSquared < 1e-10 {
		return vector.Zero3(), vector.Zero3()
	}

	// Normalize the direction
	distance := math.Sqrt(distanceSquared)
	normalizedDirection := direction.Scale(1.0 / distance)

	// Calculate the force according to the universal law of gravitation
	// F = G * m1 * m2 / r^2
	massA := a.Mass().Value()
	massB := b.Mass().Value()
	forceMagnitude := gf.G * massA * massB / distanceSquared

	// Calculate the force vectors (opposite directions)
	forceOnA := normalizedDirection.Scale(forceMagnitude)
	forceOnB := normalizedDirection.Scale(-forceMagnitude)

	return forceOnA, forceOnB
}

// IsGlobal returns true because gravity is a global force
func (gf *GravitationalForce) IsGlobal() bool {
	return true
}

// SetTheta sets the approximation parameter for the Barnes-Hut algorithm
func (gf *GravitationalForce) SetTheta(theta float64) {
	gf.Theta = theta
}

// GetTheta returns the approximation parameter for the Barnes-Hut algorithm
func (gf *GravitationalForce) GetTheta() float64 {
	return gf.Theta
}

// ConstantForce implements a constant force
type ConstantForce struct {
	force vector.Vector3
}

// NewConstantForce creates a new constant force
func NewConstantForce(force vector.Vector3) *ConstantForce {
	return &ConstantForce{
		force: force,
	}
}

// Apply applies the constant force to a body
func (cf *ConstantForce) Apply(b body.Body) vector.Vector3 {
	return cf.force
}

// ApplyBetween applies the constant force between two bodies (applies only to the first one)
func (cf *ConstantForce) ApplyBetween(a, b body.Body) (vector.Vector3, vector.Vector3) {
	return cf.force, vector.Zero3()
}

// IsGlobal returns false because constant force is not global
func (cf *ConstantForce) IsGlobal() bool {
	return false
}

// SpringForce implements a spring force
type SpringForce struct {
	stiffness  float64 // Spring constant
	restLength float64 // Rest length
	damping    float64 // Damping coefficient
}

// NewSpringForce creates a new spring force
func NewSpringForce(stiffness, restLength, damping float64) *SpringForce {
	return &SpringForce{
		stiffness:  stiffness,
		restLength: restLength,
		damping:    damping,
	}
}

// Apply applies the spring force to a body (does nothing for a single body)
func (sf *SpringForce) Apply(b body.Body) vector.Vector3 {
	// Spring force requires two bodies to be applied
	return vector.Zero3()
}

// ApplyBetween applies the spring force between two bodies
func (sf *SpringForce) ApplyBetween(a, b body.Body) (vector.Vector3, vector.Vector3) {
	// Calculate the direction vector from a to b
	direction := b.Position().Sub(a.Position())

	// Calculate the distance
	distance := direction.Length()

	// Avoid division by zero
	if distance < 1e-10 {
		return vector.Zero3(), vector.Zero3()
	}

	// Normalize the direction
	normalizedDirection := direction.Scale(1.0 / distance)

	// Calculate the spring extension (distance - rest length)
	extension := distance - sf.restLength

	// Calculate the spring force according to Hooke's law
	// F = -k * x
	springForceMagnitude := sf.stiffness * extension

	// Calculate the relative velocity along the spring direction
	relativeVelocity := b.Velocity().Sub(a.Velocity())
	relativeVelocityAlongSpring := relativeVelocity.Dot(normalizedDirection)

	// Calculate the damping force
	// F_damping = -c * v
	dampingForceMagnitude := sf.damping * relativeVelocityAlongSpring

	// Calculate the total force
	totalForceMagnitude := springForceMagnitude + dampingForceMagnitude

	// Calculate the force vectors (opposite directions)
	forceOnA := normalizedDirection.Scale(totalForceMagnitude)
	forceOnB := normalizedDirection.Scale(-totalForceMagnitude)

	return forceOnA, forceOnB
}

// IsGlobal returns false because spring force is not global
func (sf *SpringForce) IsGlobal() bool {
	return false
}

// DragForce implements a drag force
type DragForce struct {
	coefficient float64 // Drag coefficient
}

// NewDragForce creates a new drag force
func NewDragForce(coefficient float64) *DragForce {
	return &DragForce{
		coefficient: coefficient,
	}
}

// Apply applies the drag force to a body
func (df *DragForce) Apply(b body.Body) vector.Vector3 {
	// Drag force is proportional to the square of velocity
	// F = -c * v^2 * v_hat
	velocity := b.Velocity()
	speed := velocity.Length()

	// If velocity is almost zero, there is no drag
	if speed < 1e-10 {
		return vector.Zero3()
	}

	// Calculate the velocity direction
	direction := velocity.Scale(1.0 / speed)

	// Calculate the drag force (in the opposite direction of velocity)
	forceMagnitude := df.coefficient * speed * speed
	force := direction.Scale(-forceMagnitude)

	return force
}

// ApplyBetween applies the drag force between two bodies (applies to both separately)
func (df *DragForce) ApplyBetween(a, b body.Body) (vector.Vector3, vector.Vector3) {
	return df.Apply(a), df.Apply(b)
}

// IsGlobal returns true because drag force is global
func (df *DragForce) IsGlobal() bool {
	return true
}
