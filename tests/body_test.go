// Package tests provides tests for the physics engine
package tests

import (
	"math"
	"testing"

	"github.com/alexanderi96/go-space-engine/core/units"
	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/physics/body"
	"github.com/alexanderi96/go-space-engine/physics/material"
)

// TestRigidBodyCreation verifies that a rigid body is created correctly
func TestRigidBodyCreation(t *testing.T) {
	// Create a rigid body
	position := vector.NewVector3(1, 2, 3)
	velocity := vector.NewVector3(4, 5, 6)
	mass := units.NewQuantity(10, units.Kilogram)
	radius := units.NewQuantity(1, units.Meter)

	rb := body.NewRigidBody(mass, radius, position, velocity, material.Iron)

	// Verify that the body was created correctly
	if rb.Position().X() != position.X() || rb.Position().Y() != position.Y() || rb.Position().Z() != position.Z() {
		t.Errorf("Position not set correctly: got %v, want %v", rb.Position(), position)
	}

	if rb.Velocity().X() != velocity.X() || rb.Velocity().Y() != velocity.Y() || rb.Velocity().Z() != velocity.Z() {
		t.Errorf("Velocity not set correctly: got %v, want %v", rb.Velocity(), velocity)
	}

	if rb.Mass().Value() != mass.Value() {
		t.Errorf("Mass not set correctly: got %v, want %v", rb.Mass().Value(), mass.Value())
	}

	if rb.Radius().Value() != radius.Value() {
		t.Errorf("Radius not set correctly: got %v, want %v", rb.Radius().Value(), radius.Value())
	}

	if rb.Material() != material.Iron {
		t.Errorf("Material not set correctly: got %v, want %v", rb.Material(), material.Iron)
	}

	if rb.IsStatic() {
		t.Errorf("Body should not be static by default")
	}
}

// TestRigidBodySetters verifies that the rigid body setters work correctly
func TestRigidBodySetters(t *testing.T) {
	// Create a rigid body
	rb := body.NewRigidBody(
		units.NewQuantity(10, units.Kilogram),
		units.NewQuantity(1, units.Meter),
		vector.NewVector3(0, 0, 0),
		vector.NewVector3(0, 0, 0),
		material.Iron,
	)

	// Set the position
	newPosition := vector.NewVector3(1, 2, 3)
	rb.SetPosition(newPosition)
	if rb.Position().X() != newPosition.X() || rb.Position().Y() != newPosition.Y() || rb.Position().Z() != newPosition.Z() {
		t.Errorf("SetPosition failed: got %v, want %v", rb.Position(), newPosition)
	}

	// Set the velocity
	newVelocity := vector.NewVector3(4, 5, 6)
	rb.SetVelocity(newVelocity)
	if rb.Velocity().X() != newVelocity.X() || rb.Velocity().Y() != newVelocity.Y() || rb.Velocity().Z() != newVelocity.Z() {
		t.Errorf("SetVelocity failed: got %v, want %v", rb.Velocity(), newVelocity)
	}

	// Set the acceleration
	newAcceleration := vector.NewVector3(7, 8, 9)
	rb.SetAcceleration(newAcceleration)
	if rb.Acceleration().X() != newAcceleration.X() || rb.Acceleration().Y() != newAcceleration.Y() || rb.Acceleration().Z() != newAcceleration.Z() {
		t.Errorf("SetAcceleration failed: got %v, want %v", rb.Acceleration(), newAcceleration)
	}

	// Set the mass
	newMass := units.NewQuantity(20, units.Kilogram)
	rb.SetMass(newMass)
	if rb.Mass().Value() != newMass.Value() {
		t.Errorf("SetMass failed: got %v, want %v", rb.Mass().Value(), newMass.Value())
	}

	// Set the radius
	newRadius := units.NewQuantity(2, units.Meter)
	rb.SetRadius(newRadius)
	if rb.Radius().Value() != newRadius.Value() {
		t.Errorf("SetRadius failed: got %v, want %v", rb.Radius().Value(), newRadius.Value())
	}

	// Set the material
	rb.SetMaterial(material.Copper)
	if rb.Material() != material.Copper {
		t.Errorf("SetMaterial failed: got %v, want %v", rb.Material(), material.Copper)
	}

	// Set the static state
	rb.SetStatic(true)
	if !rb.IsStatic() {
		t.Errorf("SetStatic failed: body should be static")
	}

	// Verify that velocity and acceleration are zero when the body is static
	if rb.Velocity().Length() != 0 || rb.Acceleration().Length() != 0 {
		t.Errorf("Static body should have zero velocity and acceleration")
	}
}

// TestRigidBodyForce verifies that force application works correctly
func TestRigidBodyForce(t *testing.T) {
	// Create a rigid body
	mass := units.NewQuantity(10, units.Kilogram)
	rb := body.NewRigidBody(
		mass,
		units.NewQuantity(1, units.Meter),
		vector.NewVector3(0, 0, 0),
		vector.NewVector3(0, 0, 0),
		material.Iron,
	)

	// Apply a force
	force := vector.NewVector3(10, 20, 30)
	rb.ApplyForce(force)

	// Verify that the acceleration was calculated correctly (a = F/m)
	expectedAcceleration := force.Scale(1.0 / mass.Value())
	if rb.Acceleration().X() != expectedAcceleration.X() || rb.Acceleration().Y() != expectedAcceleration.Y() || rb.Acceleration().Z() != expectedAcceleration.Z() {
		t.Errorf("ApplyForce failed: got acceleration %v, want %v", rb.Acceleration(), expectedAcceleration)
	}

	// Verify that applying a force to a static body has no effect
	rb.SetStatic(true)
	rb.SetAcceleration(vector.Zero3())
	rb.ApplyForce(force)
	if rb.Acceleration().Length() != 0 {
		t.Errorf("ApplyForce should have no effect on static body")
	}
}

// TestRigidBodyUpdate verifies that body updating works correctly
func TestRigidBodyUpdate(t *testing.T) {
	// Create a rigid body
	position := vector.NewVector3(0, 0, 0)
	velocity := vector.NewVector3(1, 2, 3)
	acceleration := vector.NewVector3(4, 5, 6)

	rb := body.NewRigidBody(
		units.NewQuantity(10, units.Kilogram),
		units.NewQuantity(1, units.Meter),
		position,
		velocity,
		material.Iron,
	)
	rb.SetAcceleration(acceleration)

	// Update the body with a time step
	dt := 0.1
	rb.Update(dt)

	// Verify that position and velocity were updated correctly
	// Using Verlet integration:
	// x(t+dt) = x(t) + v(t)*dt + 0.5*a(t)*dt^2
	// v(t+dt) = v(t) + 0.5*(a(t) + a(t+dt))*dt
	// But since a(t+dt) = 0 (acceleration is reset), we have:
	// v(t+dt) = v(t) + 0.5*a(t)*dt

	expectedPosition := position.Add(velocity.Scale(dt)).Add(acceleration.Scale(0.5 * dt * dt))
	if !vectorsAlmostEqual(rb.Position(), expectedPosition, 1e-10) {
		t.Errorf("Update failed: got position %v, want %v", rb.Position(), expectedPosition)
	}

	expectedVelocity := velocity.Add(acceleration.Scale(0.5 * dt))
	if !vectorsAlmostEqual(rb.Velocity(), expectedVelocity, 1e-10) {
		t.Errorf("Update failed: got velocity %v, want %v", rb.Velocity(), expectedVelocity)
	}

	// Verify that the acceleration was reset
	if rb.Acceleration().Length() != 0 {
		t.Errorf("Update failed: acceleration should be reset to zero")
	}

	// Verify that updating a static body has no effect
	rb.SetStatic(true)
	rb.SetPosition(position)
	rb.SetVelocity(velocity)
	rb.SetAcceleration(acceleration)

	rb.Update(dt)

	if rb.Position().X() != position.X() || rb.Position().Y() != position.Y() || rb.Position().Z() != position.Z() {
		t.Errorf("Update should have no effect on static body position")
	}

	if rb.Velocity().Length() != 0 || rb.Acceleration().Length() != 0 {
		t.Errorf("Static body should have zero velocity and acceleration")
	}
}

// TestRigidBodyHeat verifies that heat management works correctly
func TestRigidBodyHeat(t *testing.T) {
	// Create a rigid body
	rb := body.NewRigidBody(
		units.NewQuantity(10, units.Kilogram),
		units.NewQuantity(1, units.Meter),
		vector.NewVector3(0, 0, 0),
		vector.NewVector3(0, 0, 0),
		material.Iron,
	)

	// Verify that the initial temperature is room temperature
	initialTemp := rb.Temperature().Value()
	if initialTemp < 273.15 || initialTemp > 303.15 { // Between 0°C and 30°C
		t.Errorf("Initial temperature should be around room temperature: got %v K", initialTemp)
	}

	// Set a new temperature
	newTemp := units.NewQuantity(373.15, units.Kelvin) // 100°C
	rb.SetTemperature(newTemp)
	if rb.Temperature().Value() != newTemp.Value() {
		t.Errorf("SetTemperature failed: got %v, want %v", rb.Temperature().Value(), newTemp.Value())
	}

	// Add heat
	heat := units.NewQuantity(1000, units.Joule)
	rb.AddHeat(heat)

	// Verify that the temperature has increased
	if rb.Temperature().Value() <= newTemp.Value() {
		t.Errorf("AddHeat failed: temperature did not increase")
	}
}

// vectorsAlmostEqual verifies if two vectors are almost equal within a tolerance
func vectorsAlmostEqual(a, b vector.Vector3, tolerance float64) bool {
	return math.Abs(a.X()-b.X()) < tolerance &&
		math.Abs(a.Y()-b.Y()) < tolerance &&
		math.Abs(a.Z()-b.Z()) < tolerance
}
