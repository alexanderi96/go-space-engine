// Package main provides an example of using the physics engine
package main

import (
	"fmt"
	"math"
	"time"

	"github.com/alexanderi96/go-space-engine/core/units"
	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/physics/body"
	"github.com/alexanderi96/go-space-engine/physics/force"
	"github.com/alexanderi96/go-space-engine/physics/material"
	"github.com/alexanderi96/go-space-engine/simulation/config"
	"github.com/alexanderi96/go-space-engine/simulation/world"
)

func main() {
	// Create a configuration for the simulation
	cfg := config.NewSimulationBuilder().
		WithTimeStep(0.01).
		WithMaxBodies(100).
		WithGravity(true).
		WithCollisions(true).
		WithBoundaryCollisions(true).
		WithWorldBounds(
			vector.NewVector3(-10, -10, -10),
			vector.NewVector3(10, 10, 10),
		).
		WithRestitution(0.7).
		WithIntegratorType("verlet").
		Build()

	// Create the simulation world
	w := world.NewPhysicalWorld(cfg.GetWorldBounds())

	// Add gravitational force
	gravityForce := force.NewGravitationalForce()
	w.AddForce(gravityForce)

	// Create some bodies
	createBodies(w)

	// Run the simulation
	runSimulation(w, cfg)
}

// createBodies creates some bodies in the world
func createBodies(w world.World) {
	// Create a massive central body (like a sun)
	sun := body.NewRigidBody(
		units.NewQuantity(1.0e6, units.Kilogram),
		units.NewQuantity(1.0, units.Meter),
		vector.NewVector3(0, 0, 0),
		vector.NewVector3(0, 0, 0),
		material.Iron,
	)
	sun.SetStatic(true) // The sun is static (doesn't move)
	w.AddBody(sun)

	// Create some planets in orbit
	createPlanet(w, 3.0, 0.3, 0.5, vector.NewVector3(0, 1, 0), material.Rock)
	createPlanet(w, 5.0, 0.4, 0.3, vector.NewVector3(0, 1, 0), material.Ice)
	createPlanet(w, 7.0, 0.5, 0.2, vector.NewVector3(0, 1, 0), material.Copper)

	// Create some moons
	createMoon(w, 3.0, 0.3, 0.7, 0.1, vector.NewVector3(0, 0, 1), material.Ice)
	createMoon(w, 5.0, 0.4, 0.9, 0.15, vector.NewVector3(0, 0, 1), material.Rock)

	// Create some random asteroids
	for i := 0; i < 10; i++ {
		angle := float64(i) * 0.628 // 2*pi/10
		distance := 9.0

		x := distance * math.Cos(angle)
		z := distance * math.Sin(angle)

		asteroid := body.NewRigidBody(
			units.NewQuantity(100.0, units.Kilogram),
			units.NewQuantity(0.1, units.Meter),
			vector.NewVector3(x, 0, z),
			vector.NewVector3(-z*0.3, 0, x*0.3), // Tangential velocity
			material.Rock,
		)
		w.AddBody(asteroid)
	}
}

// createPlanet creates a planet in orbit
func createPlanet(w world.World, distance, radius, speed float64, orbitPlane vector.Vector3, mat material.Material) body.Body {
	// Calculate the initial position
	position := vector.NewVector3(distance, 0, 0)

	// Calculate the orbital velocity (perpendicular to position)
	velocity := orbitPlane.Cross(position).Normalize().Scale(speed)

	// Create the planet
	planet := body.NewRigidBody(
		units.NewQuantity(1000.0, units.Kilogram),
		units.NewQuantity(radius, units.Meter),
		position,
		velocity,
		mat,
	)

	// Add the planet to the world
	w.AddBody(planet)

	return planet
}

// createMoon creates a moon orbiting around a planet
func createMoon(w world.World, planetDistance, planetRadius, moonDistance, moonRadius float64, orbitPlane vector.Vector3, mat material.Material) body.Body {
	// Calculate the planet position
	planetPosition := vector.NewVector3(planetDistance, 0, 0)

	// Calculate the moon position relative to the planet
	moonRelativePosition := vector.NewVector3(moonDistance, 0, 0)

	// Calculate the absolute position of the moon
	moonPosition := planetPosition.Add(moonRelativePosition)

	// Calculate the orbital velocity of the planet
	planetVelocity := orbitPlane.Cross(planetPosition).Normalize().Scale(math.Sqrt(1.0 / planetDistance))

	// Calculate the orbital velocity of the moon relative to the planet
	moonRelativeVelocity := orbitPlane.Cross(moonRelativePosition).Normalize().Scale(math.Sqrt(10.0 / moonDistance))

	// Calculate the absolute velocity of the moon
	moonVelocity := planetVelocity.Add(moonRelativeVelocity)

	// Create the moon
	moon := body.NewRigidBody(
		units.NewQuantity(100.0, units.Kilogram),
		units.NewQuantity(moonRadius, units.Meter),
		moonPosition,
		moonVelocity,
		mat,
	)

	// Add the moon to the world
	w.AddBody(moon)

	return moon
}

// runSimulation runs the simulation
func runSimulation(w world.World, cfg *config.Config) {
	// Simulation parameters
	timeStep := cfg.TimeStep
	totalTime := 100.0   // Total simulation time (seconds)
	printInterval := 1.0 // Print interval (seconds)

	// Timing variables
	lastPrintTime := 0.0
	startTime := time.Now()

	// Simulation loop
	for t := 0.0; t < totalTime; t += timeStep {
		// Execute a simulation step
		w.Step(timeStep)

		// Print the simulation state at regular intervals
		if t-lastPrintTime >= printInterval {
			// Calculate the elapsed real time
			elapsedTime := time.Since(startTime).Seconds()

			// Print the simulation state
			fmt.Printf("Simulation time: %.2f s, Real time: %.2f s, Bodies: %d\n",
				t, elapsedTime, w.GetBodyCount())

			// Update the last print time
			lastPrintTime = t
		}
	}

	// Print the total execution time
	fmt.Printf("Simulation completed in %.2f seconds\n", time.Since(startTime).Seconds())
}
