// Package tests fornisce test per il motore fisico
package tests

import (
	"math"
	"testing"

	"github.com/alexanderi96/go-space-engine/core/units"
	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/physics/body"
	"github.com/alexanderi96/go-space-engine/physics/material"
)

// TestRigidBodyCreation verifica che un corpo rigido venga creato correttamente
func TestRigidBodyCreation(t *testing.T) {
	// Crea un corpo rigido
	position := vector.NewVector3(1, 2, 3)
	velocity := vector.NewVector3(4, 5, 6)
	mass := units.NewQuantity(10, units.Kilogram)
	radius := units.NewQuantity(1, units.Meter)

	rb := body.NewRigidBody(mass, radius, position, velocity, material.Iron)

	// Verifica che il corpo sia stato creato correttamente
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

// TestRigidBodySetters verifica che i setter del corpo rigido funzionino correttamente
func TestRigidBodySetters(t *testing.T) {
	// Crea un corpo rigido
	rb := body.NewRigidBody(
		units.NewQuantity(10, units.Kilogram),
		units.NewQuantity(1, units.Meter),
		vector.NewVector3(0, 0, 0),
		vector.NewVector3(0, 0, 0),
		material.Iron,
	)

	// Imposta la posizione
	newPosition := vector.NewVector3(1, 2, 3)
	rb.SetPosition(newPosition)
	if rb.Position().X() != newPosition.X() || rb.Position().Y() != newPosition.Y() || rb.Position().Z() != newPosition.Z() {
		t.Errorf("SetPosition failed: got %v, want %v", rb.Position(), newPosition)
	}

	// Imposta la velocità
	newVelocity := vector.NewVector3(4, 5, 6)
	rb.SetVelocity(newVelocity)
	if rb.Velocity().X() != newVelocity.X() || rb.Velocity().Y() != newVelocity.Y() || rb.Velocity().Z() != newVelocity.Z() {
		t.Errorf("SetVelocity failed: got %v, want %v", rb.Velocity(), newVelocity)
	}

	// Imposta l'accelerazione
	newAcceleration := vector.NewVector3(7, 8, 9)
	rb.SetAcceleration(newAcceleration)
	if rb.Acceleration().X() != newAcceleration.X() || rb.Acceleration().Y() != newAcceleration.Y() || rb.Acceleration().Z() != newAcceleration.Z() {
		t.Errorf("SetAcceleration failed: got %v, want %v", rb.Acceleration(), newAcceleration)
	}

	// Imposta la massa
	newMass := units.NewQuantity(20, units.Kilogram)
	rb.SetMass(newMass)
	if rb.Mass().Value() != newMass.Value() {
		t.Errorf("SetMass failed: got %v, want %v", rb.Mass().Value(), newMass.Value())
	}

	// Imposta il raggio
	newRadius := units.NewQuantity(2, units.Meter)
	rb.SetRadius(newRadius)
	if rb.Radius().Value() != newRadius.Value() {
		t.Errorf("SetRadius failed: got %v, want %v", rb.Radius().Value(), newRadius.Value())
	}

	// Imposta il materiale
	rb.SetMaterial(material.Copper)
	if rb.Material() != material.Copper {
		t.Errorf("SetMaterial failed: got %v, want %v", rb.Material(), material.Copper)
	}

	// Imposta lo stato statico
	rb.SetStatic(true)
	if !rb.IsStatic() {
		t.Errorf("SetStatic failed: body should be static")
	}

	// Verifica che la velocità e l'accelerazione siano zero quando il corpo è statico
	if rb.Velocity().Length() != 0 || rb.Acceleration().Length() != 0 {
		t.Errorf("Static body should have zero velocity and acceleration")
	}
}

// TestRigidBodyForce verifica che l'applicazione di forze funzioni correttamente
func TestRigidBodyForce(t *testing.T) {
	// Crea un corpo rigido
	mass := units.NewQuantity(10, units.Kilogram)
	rb := body.NewRigidBody(
		mass,
		units.NewQuantity(1, units.Meter),
		vector.NewVector3(0, 0, 0),
		vector.NewVector3(0, 0, 0),
		material.Iron,
	)

	// Applica una forza
	force := vector.NewVector3(10, 20, 30)
	rb.ApplyForce(force)

	// Verifica che l'accelerazione sia stata calcolata correttamente (a = F/m)
	expectedAcceleration := force.Scale(1.0 / mass.Value())
	if rb.Acceleration().X() != expectedAcceleration.X() || rb.Acceleration().Y() != expectedAcceleration.Y() || rb.Acceleration().Z() != expectedAcceleration.Z() {
		t.Errorf("ApplyForce failed: got acceleration %v, want %v", rb.Acceleration(), expectedAcceleration)
	}

	// Verifica che l'applicazione di una forza a un corpo statico non abbia effetto
	rb.SetStatic(true)
	rb.SetAcceleration(vector.Zero3())
	rb.ApplyForce(force)
	if rb.Acceleration().Length() != 0 {
		t.Errorf("ApplyForce should have no effect on static body")
	}
}

// TestRigidBodyUpdate verifica che l'aggiornamento del corpo funzioni correttamente
func TestRigidBodyUpdate(t *testing.T) {
	// Crea un corpo rigido
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

	// Aggiorna il corpo con un passo temporale
	dt := 0.1
	rb.Update(dt)

	// Verifica che la posizione e la velocità siano state aggiornate correttamente
	// Usando l'integrazione di Verlet:
	// x(t+dt) = x(t) + v(t)*dt + 0.5*a(t)*dt^2
	// v(t+dt) = v(t) + 0.5*(a(t) + a(t+dt))*dt
	// Ma poiché a(t+dt) = 0 (l'accelerazione viene resettata), abbiamo:
	// v(t+dt) = v(t) + 0.5*a(t)*dt

	expectedPosition := position.Add(velocity.Scale(dt)).Add(acceleration.Scale(0.5 * dt * dt))
	if !vectorsAlmostEqual(rb.Position(), expectedPosition, 1e-10) {
		t.Errorf("Update failed: got position %v, want %v", rb.Position(), expectedPosition)
	}

	expectedVelocity := velocity.Add(acceleration.Scale(0.5 * dt))
	if !vectorsAlmostEqual(rb.Velocity(), expectedVelocity, 1e-10) {
		t.Errorf("Update failed: got velocity %v, want %v", rb.Velocity(), expectedVelocity)
	}

	// Verifica che l'accelerazione sia stata resettata
	if rb.Acceleration().Length() != 0 {
		t.Errorf("Update failed: acceleration should be reset to zero")
	}

	// Verifica che l'aggiornamento di un corpo statico non abbia effetto
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

// TestRigidBodyHeat verifica che la gestione del calore funzioni correttamente
func TestRigidBodyHeat(t *testing.T) {
	// Crea un corpo rigido
	rb := body.NewRigidBody(
		units.NewQuantity(10, units.Kilogram),
		units.NewQuantity(1, units.Meter),
		vector.NewVector3(0, 0, 0),
		vector.NewVector3(0, 0, 0),
		material.Iron,
	)

	// Verifica che la temperatura iniziale sia quella ambiente
	initialTemp := rb.Temperature().Value()
	if initialTemp < 273.15 || initialTemp > 303.15 { // Tra 0°C e 30°C
		t.Errorf("Initial temperature should be around room temperature: got %v K", initialTemp)
	}

	// Imposta una nuova temperatura
	newTemp := units.NewQuantity(373.15, units.Kelvin) // 100°C
	rb.SetTemperature(newTemp)
	if rb.Temperature().Value() != newTemp.Value() {
		t.Errorf("SetTemperature failed: got %v, want %v", rb.Temperature().Value(), newTemp.Value())
	}

	// Aggiungi calore
	heat := units.NewQuantity(1000, units.Joule)
	rb.AddHeat(heat)

	// Verifica che la temperatura sia aumentata
	if rb.Temperature().Value() <= newTemp.Value() {
		t.Errorf("AddHeat failed: temperature did not increase")
	}
}

// vectorsAlmostEqual verifica se due vettori sono quasi uguali entro una tolleranza
func vectorsAlmostEqual(a, b vector.Vector3, tolerance float64) bool {
	return math.Abs(a.X()-b.X()) < tolerance &&
		math.Abs(a.Y()-b.Y()) < tolerance &&
		math.Abs(a.Z()-b.Z()) < tolerance
}
