package main

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/alexanderi96/go-space-engine/core/units"
	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/physics/body"
	"github.com/alexanderi96/go-space-engine/physics/force"
	physMaterial "github.com/alexanderi96/go-space-engine/physics/material"
	"github.com/alexanderi96/go-space-engine/render/raylib"
	"github.com/alexanderi96/go-space-engine/simulation/config"
	"github.com/alexanderi96/go-space-engine/simulation/world"
)

func main() {
	log.Println("Initializing Controllable Spacecraft Example")

	// Create simulation configuration
	cfg := config.NewSimulationBuilder().
		WithTimeStep(0.001).
		WithMaxBodies(1000).
		WithGravity(true).
		WithCollisions(true).
		WithBoundaryCollisions(true).
		WithWorldBounds(
			vector.NewVector3(-500, -500, -500),
			vector.NewVector3(500, 500, 500),
		).
		WithOctreeConfig(10, 8).
		Build()

	// Create the simulation world
	w := world.NewPhysicalWorld(cfg.GetWorldBounds())

	// Add gravitational force
	gravityForce := force.NewGravitationalForce()
	gravityForce.SetTheta(0.5) // Set the theta value for the Barnes-Hut algorithm
	w.AddForce(gravityForce)

	// Create some bodies
	createBodies(w)

	// Create a controllable spacecraft
	spacecraft := createSpacecraft(w)

	// Create Controllable Raylib adapter
	adapter := raylib.NewControllableRaylibAdapter(2400, 1600, "Controllable Spacecraft - Raylib")

	// Set the controllable body
	adapter.SetControllableBody(spacecraft)

	// Initialize the adapter
	if err := adapter.Initialize(); err != nil {
		log.Fatalf("Failed to initialize Raylib adapter: %v", err)
	}

	// Set camera position for better view
	adapter.SetCameraPosition(vector.NewVector3(0, 50, 150))
	adapter.SetCameraTarget(vector.Zero3())

	// Enable debug features
	adapter.SetDebugMode(false)
	adapter.SetRenderVelocities(true)
	adapter.SetRenderAccelerations(true)

	// Attiva la modalità controllo navicella all'avvio
	adapter.ToggleControlMode()

	log.Println("Premi TAB per passare dal controllo della navicella al controllo della camera e viceversa")
	log.Println("Usa WASD per muovere la navicella, QE per salire/scendere, e le frecce per ruotare")

	log.Println("Starting simulation with Raylib renderer...")

	// Variables for timing
	lastUpdateTime := time.Now()

	// Run the simulation
	adapter.Run(func(deltaTime time.Duration) {
		// Calculate delta time
		currentTime := time.Now()
		dt := currentTime.Sub(lastUpdateTime).Seconds()
		lastUpdateTime = currentTime

		// Limit delta time to avoid instability
		if dt > 0.1 {
			dt = 0.1
		}

		// Execute a simulation step with fixed timestep for stability
		w.Step(0.001)

		// Render the world
		adapter.RenderWorld(w)
	})

	log.Println("Simulation ended")
}

// createBodies creates some bodies in the world
func createBodies(w world.World) {
	log.Println("Creating the solar system")
	createSolarSystem(w)
}

// createSolarSystem creates a realistic solar system
func createSolarSystem(w world.World) {
	log.Println("Creating the sun")

	// Fixed solar mass - high value to ensure stable orbits
	solarMass := 3e14 // Simplified value

	log.Printf("Solar mass: %e kg", solarMass)

	sun := body.NewRigidBody(
		units.NewQuantity(solarMass, units.Kilogram),
		units.NewQuantity(5.0, units.Kilometer),                    // Solar radius (scaled)
		vector.NewVector3(0, 0, 0),                                 // Position at center
		vector.NewVector3(0, 0, 0),                                 // Zero velocity
		createMaterial("Sun", 0.9, 0.5, [3]float64{1.0, 0.8, 0.0}), // Yellow color
	)
	sun.SetStatic(true) // The sun is static (does not move)
	w.AddBody(sun)
	log.Printf("Sun created: ID=%v, Position=%v", sun.ID(), sun.Position())

	// Create the planets
	log.Println("Creating the planets")

	// Define planet distances
	distances := []float64{20, 30, 40, 50, 70, 90, 110, 130}

	// Planet names
	names := []string{"Mercury", "Venus", "Earth", "Mars", "Jupiter", "Saturn", "Uranus", "Neptune"}

	// Planet radii (scaled)
	radii := []float64{0.8, 1.2, 1.3, 1.0, 2.5, 2.2, 1.8, 1.8}

	// Planet masses as fractions of the solar mass
	massFractions := []float64{1e-6, 2e-6, 2e-6, 1e-6, 1e-6, 9e-6, 4e-6, 5e-6}

	// Planet colors
	colors := [][3]float64{
		{0.7, 0.7, 0.7}, // Mercury
		{0.9, 0.7, 0.0}, // Venus
		{0.0, 0.3, 0.8}, // Earth
		{0.8, 0.3, 0.0}, // Mars
		{0.8, 0.6, 0.4}, // Jupiter
		{0.9, 0.8, 0.5}, // Saturn
		{0.5, 0.8, 0.9}, // Uranus
		{0.0, 0.0, 0.8}, // Neptune
	}

	// Create each planet
	for i := 0; i < len(names); i++ {
		createPlanet(
			w,
			names[i],                   // Name
			solarMass*massFractions[i], // Mass
			radii[i],                   // Radius
			distances[i],               // Distance
			solarMass,                  // Central object mass (sun)
			vector.NewVector3(0, 1, 0), // Orbital plane
			colors[i],                  // Color
		)
	}
}

// createPlanet creates a planet
func createPlanet(w world.World, name string, mass, radius, distance, centralMass float64, orbitPlane vector.Vector3, color [3]float64) body.Body {
	// Calculate orbital velocity
	orbitSpeed := force.CalculateOrbitalVelocity(centralMass, distance)

	log.Printf("Creating planet %s: distance=%f, radius=%f, calculated orbit speed=%f", name, distance, radius, orbitSpeed)

	// Random angle for the initial position
	angle := rand.Float64() * 2 * math.Pi

	// Calculate the initial position
	position := vector.NewVector3(
		distance*math.Cos(angle),
		0,
		distance*math.Sin(angle),
	)

	// Calculate the orbital velocity
	velocity := vector.NewVector3(
		-orbitSpeed*math.Sin(angle), // x component
		0,                           // y component (xy plane)
		orbitSpeed*math.Cos(angle),  // z component
	)

	// Create the planet
	planet := body.NewRigidBody(
		units.NewQuantity(mass, units.Kilogram),
		units.NewQuantity(radius, units.Kilometer),
		position,
		velocity,
		createMaterial(name, 0.7, 0.5, color),
	)

	// Add the planet to the world
	w.AddBody(planet)
	log.Printf("Planet %s added: ID=%v, Position=%v, Velocity=%v", name, planet.ID(), planet.Position(), planet.Velocity())

	return planet
}

// createSpacecraft creates a controllable spacecraft
func createSpacecraft(w world.World) body.ControllableBody {
	// Create a material for the spacecraft
	material := createMaterial("Spacecraft", 0.7, 0.5, [3]float64{1.0, 1.0, 1.0})

	// Create a rigid body for the spacecraft
	rb := body.NewRigidBody(
		units.NewQuantity(100, units.Kilogram), // Ridotta la massa da 1000 a 100 kg per renderla più reattiva
		units.NewQuantity(1.0, units.Meter),
		vector.NewVector3(0, 0, 20), // Position near the sun
		vector.Zero3(),              // Zero initial velocity
		material,
	)

	// Create a controllable rigid body
	// Aumentiamo significativamente la potenza di propulsione (da 1000.0 a 50000.0)
	// per rendere più evidente il movimento della navicella
	spacecraft := body.NewControllableRigidBody(rb, 50000.0, 1.0)

	// Add the spacecraft to the world
	w.AddBody(spacecraft)
	log.Printf("Spacecraft added: ID=%v, Position=%v", spacecraft.ID(), spacecraft.Position())

	return spacecraft
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
