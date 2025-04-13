// Package body provides interfaces and implementations for physical bodies
package body

import (
	"github.com/alexanderi96/go-space-engine/core/units"
	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/google/uuid"
)

// ID represents a unique identifier for a body
// type ID uuid.UUID

// NewID generates a new unique ID
// func NewID() ID {

// 	// Generates a random ID of 16 characters
// 	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
// 	result := make([]byte, 16)
// 	for i := range result {
// 		result[i] = chars[rand.Intn(len(chars))]
// 	}

// 	return ID(fmt.Sprintf("%s-%d", string(result), time.Now().UnixNano()))
// }

// Material represents the physical properties of a material
type Material interface {
	// Name returns the name of the material
	Name() string

	// Density returns the density of the material
	Density() units.Quantity

	// SpecificHeat returns the specific heat capacity of the material
	SpecificHeat() units.Quantity

	// ThermalConductivity returns the thermal conductivity of the material
	ThermalConductivity() units.Quantity

	// Emissivity returns the emissivity of the material
	Emissivity() float64

	// Elasticity returns the elasticity of the material
	Elasticity() float64
}

// Body represents a physical body in the engine
type Body interface {
	// ID returns the unique identifier of the body
	ID() uuid.UUID

	// Position returns the position of the body
	Position() vector.Vector3
	// SetPosition sets the position of the body
	SetPosition(pos vector.Vector3)

	// Velocity returns the velocity of the body
	Velocity() vector.Vector3
	// SetVelocity sets the velocity of the body
	SetVelocity(vel vector.Vector3)

	// Acceleration returns the acceleration of the body
	Acceleration() vector.Vector3
	// SetAcceleration sets the acceleration of the body
	SetAcceleration(acc vector.Vector3)

	// Rotation returns the rotation of the body
	Rotation() vector.Vector3
	// SetRotation sets the rotation of the body
	SetRotation(rot vector.Vector3)

	// AngularVelocity returns the angular velocity of the body
	AngularVelocity() vector.Vector3
	// SetAngularVelocity sets the angular velocity of the body
	SetAngularVelocity(angVel vector.Vector3)

	// Mass returns the mass of the body
	Mass() units.Quantity
	// SetMass sets the mass of the body
	SetMass(mass units.Quantity)

	// Radius returns the radius of the body
	Radius() units.Quantity
	// SetRadius sets the radius of the body
	SetRadius(radius units.Quantity)

	// Material returns the material of the body
	Material() Material
	// SetMaterial sets the material of the body
	SetMaterial(mat Material)

	// ApplyForce applies a force to the body
	ApplyForce(force vector.Vector3)

	// ApplyTorque applies a torque to the body
	ApplyTorque(torque vector.Vector3)

	// Update updates the state of the body
	Update(dt float64)

	// Temperature returns the temperature of the body
	Temperature() units.Quantity
	// SetTemperature sets the temperature of the body
	SetTemperature(temp units.Quantity)

	// AddHeat adds heat to the body
	AddHeat(heat units.Quantity)

	// IsStatic returns true if the body is static (does not move)
	IsStatic() bool
	// SetStatic sets whether the body is static
	SetStatic(static bool)
}

// RigidBody implements a rigid body
type RigidBody struct {
	id           uuid.UUID
	position     vector.Vector3
	velocity     vector.Vector3
	acceleration vector.Vector3
	rotation     vector.Vector3
	angularVel   vector.Vector3
	angularAcc   vector.Vector3
	mass         units.Quantity
	radius       units.Quantity
	material     Material
	temperature  units.Quantity
	isStatic     bool
}

// NewRigidBody creates a new rigid body
func NewRigidBody(
	mass units.Quantity,
	radius units.Quantity,
	position vector.Vector3,
	velocity vector.Vector3,
	mat Material,
) *RigidBody {
	return &RigidBody{
		id:           uuid.New(),
		position:     position,
		velocity:     velocity,
		acceleration: vector.Zero3(),
		rotation:     vector.Zero3(),
		angularVel:   vector.Zero3(),
		angularAcc:   vector.Zero3(),
		mass:         mass,
		radius:       radius,
		material:     mat,
		temperature:  units.NewQuantity(293.15, units.Kelvin), // Room temperature (20°C)
		isStatic:     false,
	}
}

// ID returns the unique identifier of the body
func (rb *RigidBody) ID() uuid.UUID {
	return rb.id
}

// Position returns the position of the body
func (rb *RigidBody) Position() vector.Vector3 {
	return rb.position
}

// SetPosition sets the position of the body
func (rb *RigidBody) SetPosition(pos vector.Vector3) {
	rb.position = pos
}

// Velocity returns the velocity of the body
func (rb *RigidBody) Velocity() vector.Vector3 {
	return rb.velocity
}

// SetVelocity sets the velocity of the body
func (rb *RigidBody) SetVelocity(vel vector.Vector3) {
	rb.velocity = vel

	// If the body is static, velocity must be zero
	if rb.isStatic {
		rb.velocity = vector.Zero3()
	}
}

// Acceleration returns the acceleration of the body
func (rb *RigidBody) Acceleration() vector.Vector3 {
	return rb.acceleration
}

// SetAcceleration sets the acceleration of the body
func (rb *RigidBody) SetAcceleration(acc vector.Vector3) {
	rb.acceleration = acc

	// If the body is static, acceleration must be zero
	if rb.isStatic {
		rb.acceleration = vector.Zero3()
	}
}

// Mass returns the mass of the body
func (rb *RigidBody) Mass() units.Quantity {
	return rb.mass
}

// SetMass sets the mass of the body
func (rb *RigidBody) SetMass(mass units.Quantity) {
	if mass.Unit().Type() != units.Mass {
		panic("Mass must be a mass quantity")
	}
	rb.mass = mass
}

// Radius returns the radius of the body
func (rb *RigidBody) Radius() units.Quantity {
	return rb.radius
}

// SetRadius sets the radius of the body
func (rb *RigidBody) SetRadius(radius units.Quantity) {
	if radius.Unit().Type() != units.Length {
		panic("Radius must be a length quantity")
	}
	rb.radius = radius
}

// Material returns the material of the body
func (rb *RigidBody) Material() Material {
	return rb.material
}

// SetMaterial sets the material of the body
func (rb *RigidBody) SetMaterial(mat Material) {
	rb.material = mat
}

// ApplyForce applies a force to the body
func (rb *RigidBody) ApplyForce(force vector.Vector3) {
	if rb.isStatic {
		return
	}

	// F = m*a => a = F/m
	massValue := rb.mass.Value()
	if massValue <= 0 {
		return
	}

	// Calculate the acceleration (F/m) and add it to the current acceleration
	acceleration := force.Scale(1.0 / massValue)
	rb.acceleration = rb.acceleration.Add(acceleration)
}

// Update updates the state of the body
func (rb *RigidBody) Update(dt float64) {
	// If the body is static, ensure velocity and acceleration are zero
	if rb.isStatic {
		rb.velocity = vector.Zero3()
		rb.acceleration = vector.Zero3()
		rb.angularVel = vector.Zero3()
		rb.angularAcc = vector.Zero3()
		return
	}

	// Verlet integration for linear motion
	// x(t+dt) = x(t) + v(t)*dt + 0.5*a(t)*dt^2
	// v(t+dt) = v(t) + 0.5*(a(t) + a(t+dt))*dt

	// Save the current acceleration
	oldAcceleration := rb.acceleration

	// Update the position
	halfDtSquared := 0.5 * dt * dt
	dtVelocity := rb.velocity.Scale(dt)
	dtAcceleration := rb.acceleration.Scale(halfDtSquared)
	rb.position = rb.position.Add(dtVelocity).Add(dtAcceleration)

	// Reset the acceleration (will be recalculated in the next cycle)
	rb.acceleration = vector.Zero3()

	// Update the velocity (using the average of accelerations)
	avgAcceleration := oldAcceleration.Add(rb.acceleration).Scale(0.5)
	rb.velocity = rb.velocity.Add(avgAcceleration.Scale(dt))

	// Update rotation based on angular velocity
	oldAngularAcc := rb.angularAcc

	// Update rotation
	rb.rotation = rb.rotation.Add(rb.angularVel.Scale(dt))

	// Reset angular acceleration (will be recalculated in the next cycle)
	rb.angularAcc = vector.Zero3()

	// Update angular velocity (using the average of angular accelerations)
	avgAngularAcc := oldAngularAcc.Add(rb.angularAcc).Scale(0.5)
	rb.angularVel = rb.angularVel.Add(avgAngularAcc.Scale(dt))
}

// Temperature returns the temperature of the body
func (rb *RigidBody) Temperature() units.Quantity {
	return rb.temperature
}

// SetTemperature sets the temperature of the body
func (rb *RigidBody) SetTemperature(temp units.Quantity) {
	if temp.Unit().Type() != units.Temperature {
		panic("Temperature must be a temperature quantity")
	}
	rb.temperature = temp
}

// AddHeat adds heat to the body
func (rb *RigidBody) AddHeat(heat units.Quantity) {
	if heat.Unit().Type() != units.Energy {
		panic("Heat must be an energy quantity")
	}

	// Q = m*c*ΔT => ΔT = Q/(m*c)
	massValue := rb.mass.Value()
	if massValue <= 0 {
		return
	}

	// Get the specific heat capacity from the material
	specificHeat := rb.material.SpecificHeat()
	specificHeatValue := specificHeat.Value()

	// Calculate the temperature change
	heatValue := heat.Value()
	deltaTemp := heatValue / (massValue * specificHeatValue)

	// Update the temperature
	currentTemp := rb.temperature.Value()
	newTemp := currentTemp + deltaTemp
	rb.temperature = units.NewQuantity(newTemp, units.Kelvin)
}

// IsStatic returns true if the body is static (does not move)
func (rb *RigidBody) IsStatic() bool {
	return rb.isStatic
}

// SetStatic sets whether the body is static
func (rb *RigidBody) SetStatic(static bool) {
	rb.isStatic = static
	if static {
		rb.velocity = vector.Zero3()
		rb.acceleration = vector.Zero3()
		rb.angularVel = vector.Zero3()
		rb.angularAcc = vector.Zero3()
	}
}

// Rotation returns the rotation of the body
func (rb *RigidBody) Rotation() vector.Vector3 {
	return rb.rotation
}

// SetRotation sets the rotation of the body
func (rb *RigidBody) SetRotation(rot vector.Vector3) {
	rb.rotation = rot
}

// AngularVelocity returns the angular velocity of the body
func (rb *RigidBody) AngularVelocity() vector.Vector3 {
	return rb.angularVel
}

// SetAngularVelocity sets the angular velocity of the body
func (rb *RigidBody) SetAngularVelocity(angVel vector.Vector3) {
	rb.angularVel = angVel

	// If the body is static, angular velocity must be zero
	if rb.isStatic {
		rb.angularVel = vector.Zero3()
	}
}

// ApplyTorque applies a torque to the body
func (rb *RigidBody) ApplyTorque(torque vector.Vector3) {
	if rb.isStatic {
		return
	}

	// Add the torque to the current angular acceleration
	// In a more complex implementation, this would consider the moment of inertia
	rb.angularAcc = rb.angularAcc.Add(torque)
}
