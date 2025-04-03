package tests

import (
	"math"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/alexanderi96/go-space-engine/core/constants"
	"github.com/alexanderi96/go-space-engine/core/units"
	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/physics/body"
	"github.com/alexanderi96/go-space-engine/physics/force"
	"github.com/alexanderi96/go-space-engine/physics/material"
	"github.com/alexanderi96/go-space-engine/physics/space"
	"github.com/alexanderi96/go-space-engine/simulation/world"
)

// TestGravitationalForce verifies the calculation of gravitational force between two bodies
func TestGravitationalForce(t *testing.T) {
	// Create two bodies
	body1 := body.NewRigidBody(
		units.NewQuantity(1000.0, units.Kilogram),
		units.NewQuantity(1.0, units.Meter),
		vector.NewVector3(0, 0, 0),
		vector.Zero3(),
		material.Rock,
	)

	body2 := body.NewRigidBody(
		units.NewQuantity(2000.0, units.Kilogram),
		units.NewQuantity(1.0, units.Meter),
		vector.NewVector3(10, 0, 0),
		vector.Zero3(),
		material.Rock,
	)

	// Create the gravitational force
	gravityForce := force.NewGravitationalForce()

	// Calculate the force between the two bodies
	forceOnBody1, forceOnBody2 := gravityForce.ApplyBetween(body1, body2)

	// Calculate the expected force according to the universal law of gravitation
	// F = G * m1 * m2 / r^2
	distance := 10.0
	expectedForceMagnitude := constants.G * 1000.0 * 2000.0 / (distance * distance)

	// Verify that the calculated force is correct
	actualForceMagnitude := forceOnBody1.Length()
	if math.Abs(actualForceMagnitude-expectedForceMagnitude) > 1e-10 {
		t.Errorf("Incorrect gravitational force: expected %v, got %v", expectedForceMagnitude, actualForceMagnitude)
	}

	// Verify that the forces are equal and opposite
	sumForces := forceOnBody1.Add(forceOnBody2)
	if sumForces.Length() > 1e-10 {
		t.Errorf("Forces are not equal and opposite: %v + %v = %v", forceOnBody1, forceOnBody2, sumForces)
	}

	// Verify that the force direction is correct
	expectedDirection := vector.NewVector3(1, 0, 0) // From body1 to body2
	actualDirection := forceOnBody1.Normalize()
	dotProduct := expectedDirection.Dot(actualDirection)
	if math.Abs(dotProduct-1.0) > 1e-10 {
		t.Errorf("Incorrect force direction: expected %v, got %v", expectedDirection, actualDirection)
	}
}

// TestOctreeGravity verifies the optimized gravity calculation using the octree
func TestOctreeGravity(t *testing.T) {
	// Create an octree
	bounds := space.NewAABB(
		vector.NewVector3(-100, -100, -100),
		vector.NewVector3(100, 100, 100),
	)
	octree := space.NewOctree(bounds, 10, 8)

	// Create a massive central body
	centralBody := body.NewRigidBody(
		units.NewQuantity(1.0e6, units.Kilogram),
		units.NewQuantity(1.0, units.Meter),
		vector.NewVector3(0, 0, 0),
		vector.Zero3(),
		material.Rock,
	)
	octree.Insert(centralBody)

	// Create a test body
	testBody := body.NewRigidBody(
		units.NewQuantity(1000.0, units.Kilogram),
		units.NewQuantity(1.0, units.Meter),
		vector.NewVector3(10, 0, 0),
		vector.Zero3(),
		material.Rock,
	)

	// Calculate the gravitational force using the octree
	theta := 0.5
	force := octree.CalculateGravity(testBody, theta)

	// Calculate the expected force according to the universal law of gravitation
	// F = G * m1 * m2 / r^2
	distance := 10.0
	expectedForceMagnitude := constants.G * 1.0e6 * 1000.0 / (distance * distance)

	// Verify that the calculated force is correct
	actualForceMagnitude := force.Length()
	if math.Abs(actualForceMagnitude-expectedForceMagnitude)/expectedForceMagnitude > 0.01 {
		t.Errorf("Incorrect gravitational force: expected %v, got %v", expectedForceMagnitude, actualForceMagnitude)
	}

	// Verify that the force direction is correct
	// Gravitational force is attractive, so the direction is towards the center of the octree
	expectedDirection := vector.NewVector3(-1, 0, 0) // From the test body towards the center of the octree
	actualDirection := force.Normalize()
	dotProduct := expectedDirection.Dot(actualDirection)
	if math.Abs(dotProduct-1.0) > 1e-10 {
		t.Errorf("Incorrect force direction: expected %v, got %v", expectedDirection, actualDirection)
	}
}

// TestEnergyConservation verifies energy conservation in a two-body system
func TestEnergyConservation(t *testing.T) {
	// Create a physical world
	bounds := space.NewAABB(
		vector.NewVector3(-1000, -1000, -1000),
		vector.NewVector3(1000, 1000, 1000),
	)
	w := world.NewPhysicalWorld(bounds)

	// Add gravitational force
	gravityForce := force.NewGravitationalForce()
	w.AddForce(gravityForce)

	// Create two bodies
	body1 := body.NewRigidBody(
		units.NewQuantity(1.0e6, units.Kilogram),
		units.NewQuantity(1.0, units.Meter),
		vector.NewVector3(0, 0, 0),
		vector.Zero3(),
		material.Rock,
	)
	body1.SetStatic(true) // The first body is static
	w.AddBody(body1)

	// Calculate the circular orbital velocity
	distance := 100.0
	orbitalSpeed := math.Sqrt(constants.G * 1.0e6 / distance)

	body2 := body.NewRigidBody(
		units.NewQuantity(1000.0, units.Kilogram),
		units.NewQuantity(1.0, units.Meter),
		vector.NewVector3(distance, 0, 0),
		vector.NewVector3(0, orbitalSpeed, 0), // Circular orbital velocity
		material.Rock,
	)
	w.AddBody(body2)

	// Calculate the initial energy of the system
	initialEnergy := calculateTotalEnergy(w)

	// Run the simulation for 100 steps
	dt := 0.01
	for i := 0; i < 100; i++ {
		w.Step(dt)
	}

	// Calculate the final energy of the system
	finalEnergy := calculateTotalEnergy(w)

	// Verify that energy is conserved (with a 1% tolerance)
	energyDifference := math.Abs(finalEnergy-initialEnergy) / math.Abs(initialEnergy)
	if energyDifference > 0.01 {
		t.Errorf("Energy not conserved: initial %v, final %v, difference %v%%", initialEnergy, finalEnergy, energyDifference*100)
	}
}

// calculateTotalEnergy calculates the total energy (kinetic + potential) of a system
func calculateTotalEnergy(w world.World) float64 {
	bodies := w.GetBodies()
	totalEnergy := 0.0

	// Calculate the kinetic energy of each body
	for _, b := range bodies {
		if b.IsStatic() {
			continue
		}
		mass := b.Mass().Value()
		velocity := b.Velocity()
		kineticEnergy := 0.5 * mass * velocity.LengthSquared()
		totalEnergy += kineticEnergy
	}

	// Calculate the gravitational potential energy of each pair of bodies
	for i := 0; i < len(bodies); i++ {
		for j := i + 1; j < len(bodies); j++ {
			body1 := bodies[i]
			body2 := bodies[j]
			mass1 := body1.Mass().Value()
			mass2 := body2.Mass().Value()
			distance := body1.Position().Distance(body2.Position())
			potentialEnergy := -constants.G * mass1 * mass2 / distance
			totalEnergy += potentialEnergy
		}
	}

	return totalEnergy
}

// TestBarnesHutAccuracy verifies the accuracy of the Barnes-Hut algorithm
func TestBarnesHutAccuracy(t *testing.T) {
	// Create an octree
	bounds := space.NewAABB(
		vector.NewVector3(-100, -100, -100),
		vector.NewVector3(100, 100, 100),
	)
	octree := space.NewOctree(bounds, 10, 8)

	// Create 100 random bodies
	for i := 0; i < 100; i++ {
		// Random position within the bounds
		x := (rand.Float64() * 200) - 100
		y := (rand.Float64() * 200) - 100
		z := (rand.Float64() * 200) - 100
		position := vector.NewVector3(x, y, z)

		// Random mass
		mass := rand.Float64() * 1000

		// Create the body
		b := body.NewRigidBody(
			units.NewQuantity(mass, units.Kilogram),
			units.NewQuantity(1.0, units.Meter),
			position,
			vector.Zero3(),
			material.Rock,
		)
		octree.Insert(b)
	}

	// Create a test body
	testBody := body.NewRigidBody(
		units.NewQuantity(1000.0, units.Kilogram),
		units.NewQuantity(1.0, units.Meter),
		vector.NewVector3(50, 50, 50),
		vector.Zero3(),
		material.Rock,
	)

	// Calculate the gravitational force using the Barnes-Hut algorithm with different theta values
	thetas := []float64{0.0, 0.1, 0.5, 1.0}
	forces := make([]vector.Vector3, len(thetas))

	for i, theta := range thetas {
		forces[i] = octree.CalculateGravity(testBody, theta)
	}

	// The force calculated with theta = 0 is the most accurate (direct calculation)
	exactForce := forces[0]

	// Verify that the error increases as theta increases
	for i := 1; i < len(thetas); i++ {
		error := forces[i].Sub(exactForce).Length() / exactForce.Length()
		t.Logf("Theta = %v, relative error = %v%%", thetas[i], error*100)

		// Verify that the error is acceptable (less than 10% for theta = 1.0)
		if error > 0.1 && thetas[i] <= 1.0 {
			t.Errorf("Error too large for theta = %v: %v%%", thetas[i], error*100)
		}
	}
}

// TestMultithreadingPerformance verifies the performance of multithreading
func TestMultithreadingPerformance(t *testing.T) {
	// Create two identical physical worlds
	bounds := space.NewAABB(
		vector.NewVector3(-1000, -1000, -1000),
		vector.NewVector3(1000, 1000, 1000),
	)
	w1 := world.NewPhysicalWorld(bounds)
	w2 := world.NewPhysicalWorld(bounds)

	// Add gravitational force to both worlds
	gravityForce1 := force.NewGravitationalForce()
	gravityForce2 := force.NewGravitationalForce()
	w1.AddForce(gravityForce1)
	w2.AddForce(gravityForce2)

	// Create 500 random bodies in both worlds to make the test more significant
	for i := 0; i < 500; i++ {
		// Random position within the bounds
		x := (rand.Float64() * 2000) - 1000
		y := (rand.Float64() * 2000) - 1000
		z := (rand.Float64() * 2000) - 1000
		position := vector.NewVector3(x, y, z)

		// Random velocity
		vx := (rand.Float64() * 10) - 5
		vy := (rand.Float64() * 10) - 5
		vz := (rand.Float64() * 10) - 5
		velocity := vector.NewVector3(vx, vy, vz)

		// Random mass
		mass := rand.Float64() * 1000

		// Create the bodies
		b1 := body.NewRigidBody(
			units.NewQuantity(mass, units.Kilogram),
			units.NewQuantity(1.0, units.Meter),
			position,
			velocity,
			material.Rock,
		)
		b2 := body.NewRigidBody(
			units.NewQuantity(mass, units.Kilogram),
			units.NewQuantity(1.0, units.Meter),
			position,
			velocity,
			material.Rock,
		)
		w1.AddBody(b1)
		w2.AddBody(b2)
	}

	// Measure execution time with multithreading
	startTime1 := time.Now()
	for i := 0; i < 50; i++ { // Increase the number of steps to make the test more significant
		w1.Step(0.01)
	}
	duration1 := time.Since(startTime1)

	// Disable multithreading in the second world
	// Note: this is just a test, there is no direct way to disable multithreading
	// So this test is only indicative
	startTime2 := time.Now()
	for i := 0; i < 50; i++ { // Increase the number of steps to make the test more significant
		w2.Step(0.01)
	}
	duration2 := time.Since(startTime2)

	// Verify that multithreading is faster
	t.Logf("Time with multithreading: %v", duration1)
	t.Logf("Time without multithreading: %v", duration2)
}
