// Package main provides an example of using G3N with the physics engine via direct adapter
package main

import (
	"log"
	"math"
	"time"

	"github.com/alexanderi96/go-space-engine/core/units"
	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/physics/body"
	"github.com/alexanderi96/go-space-engine/physics/force"
	physMaterial "github.com/alexanderi96/go-space-engine/physics/material"
	"github.com/alexanderi96/go-space-engine/render/g3n"
	"github.com/alexanderi96/go-space-engine/simulation/config"
	"github.com/alexanderi96/go-space-engine/simulation/world"
)

func main() {
	log.Println("Initializing G3N Physics example with Direct Adapter")

	// Create the simulation configuration
	cfg := config.NewSimulationBuilder().
		WithTimeStep(0.01).
		WithMaxBodies(1000).
		WithGravity(true).
		WithCollisions(true).
		WithBoundaryCollisions(true). // We disable collisions with the boundaries
		WithWorldBounds(
			vector.NewVector3(-500, -500, -500),
			vector.NewVector3(500, 500, 500),
		).
		WithOctreeConfig(10, 8).
		Build()

	// Create the simulation world
	w := world.NewPhysicalWorld(cfg.GetWorldBounds())

	// Add the gravitational force
	gravityForce := force.NewGravitationalForce()
	gravityForce.SetTheta(0.5) // Set the theta value for the Barnes-Hut algorithm
	w.AddForce(gravityForce)

	// Create some bodies
	createBodies(w)

	// Create the direct G3N adapter
	adapter := g3n.NewG3NAdapter()

	// Configure the adapter
	adapter.SetBackgroundColor(g3n.NewColor(1.0, 1.0, 1.0, 1.0)) // White background

	// Variables for timing
	lastUpdateTime := time.Now()

	// Start the rendering loop
	adapter.Run(func(deltaTime time.Duration) {
		// Calculate the delta time
		currentTime := time.Now()
		dt := currentTime.Sub(lastUpdateTime).Seconds()
		lastUpdateTime = currentTime

		// Limit the delta time to avoid instability
		if dt > 0.1 {
			dt = 0.1
		}

		// Execute a simulation step
		w.Step(dt)

		// Render the world
		adapter.RenderWorld(w)
	})

	log.Println("Example completed")
}

// createBodies creates some bodies in the world
func createBodies(w world.World) {
	create3BodySystem(w)
}

// create3BodySystem creates a stable system of three bodies with equal masses
func create3BodySystem(w world.World) {
	// Equal mass for all bodies
	mass := 1.0e14
	log.Printf("Mass of each body: %e kg", mass)

	// Orbit radius
	radius := 10.0

	// Colors for the three bodies
	colors := [][3]float64{
		{1.0, 0.3, 0.3}, // Red
		{0.3, 1.0, 0.3}, // Green
		{0.3, 0.3, 1.0}, // Blue
	}

	// Names of the bodies
	names := []string{"Body1", "Body2", "Body3"}

	// Calculate the orbital velocity needed for a stable orbit
	// For a three-body system with equal masses in an equilateral triangular configuration
	// Using our agnostic orbital velocity calculation function
	orbitSpeed := force.CalculateOrbitalVelocity(mass, radius)

	// Apply a scale factor to further slow down the movement
	// and make the simulation more visually pleasing
	orbitSpeed *= 0.1

	// Create the three bodies positioned at the vertices of an equilateral triangle
	for i := 0; i < 3; i++ {
		// Calculate the angle for this body (120 degrees apart from each other)
		angle := float64(i) * (2.0 * math.Pi / 3.0)

		// Calculate the position (vertices of an equilateral triangle)
		position := vector.NewVector3(
			radius*math.Cos(angle),
			0,
			radius*math.Sin(angle),
		)

		// Calculate the velocity (perpendicular to the position for a circular orbit)
		velocity := vector.NewVector3(
			-orbitSpeed*math.Sin(angle),
			0,
			orbitSpeed*math.Cos(angle),
		)

		// Create the body
		b := body.NewRigidBody(
			units.NewQuantity(mass, units.Kilogram),
			units.NewQuantity(2.0, units.Meter), // Body radius
			position,
			velocity,
			createMaterial(names[i], 0.9, 0.5, colors[i]),
		)

		w.AddBody(b)
		log.Printf("%s created: ID=%v, Position=%v, Velocity=%v", names[i], b.ID(), b.Position(), b.Velocity())
	}
}

// createMaterial creates a custom material
func createMaterial(name string, emissivity, elasticity float64, color [3]float64) physMaterial.Material {
	return physMaterial.NewBasicMaterial(
		name,
		units.NewQuantity(5000, units.Kilogram),
		units.NewQuantity(800, units.Joule),
		units.NewQuantity(1.5, units.Watt),
		emissivity,
		elasticity,
		color,
	)
}
