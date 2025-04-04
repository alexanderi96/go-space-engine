// Package main provides an example of using G3N with the physics engine via direct adapter
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
	"github.com/alexanderi96/go-space-engine/render/g3n"
	"github.com/alexanderi96/go-space-engine/simulation/config"
	"github.com/alexanderi96/go-space-engine/simulation/world"
)

func main() {
	log.Println("Initializing G3N Physics example with Direct Adapter")

	// Crea la configurazione della simulazione
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

	// Crea il mondo della simulazione
	w := world.NewPhysicalWorld(cfg.GetWorldBounds())

	// Aggiungi la forza gravitazionale
	gravityForce := force.NewGravitationalForce()
	gravityForce.SetTheta(0.5) // Set the theta value for the Barnes-Hut algorithm
	w.AddForce(gravityForce)

	// Crea alcuni corpi
	createBodies(w)

	// Crea l'adapter G3N diretto
	adapter := g3n.NewG3NAdapter()

	// Configura l'adapter
	adapter.SetBackgroundColor(g3n.NewColor(0.2, 0.2, 0.2, 1.0)) // Dark blue background for space

	// Variabili per il timing
	lastUpdateTime := time.Now()

	// Avvia il loop di rendering
	adapter.Run(func(deltaTime time.Duration) {
		// Calcola il delta time
		currentTime := time.Now()
		dt := currentTime.Sub(lastUpdateTime).Seconds()
		lastUpdateTime = currentTime

		// Limita il delta time per evitare instabilità
		if dt > 0.1 {
			dt = 0.1
		}

		// Esegui un passo della simulazione
		w.Step(0.01)

		// Renderizza il mondo
		adapter.RenderWorld(w)
	})

	log.Println("Example completed")
}

// createBodies creates some bodies in the world
func createBodies(w world.World) {
	log.Println("Creating the solar system")
	createSolarSystem(w)
}

// createSolarSystem creates a realistic solar system
func createSolarSystem(w world.World) {
	log.Println("Creating the sun")

	// Massa fissa del sole - valore elevato per garantire orbite stabili
	// In una simulazione, i rapporti relativi sono più importanti dei valori assoluti
	solarMass := 1e11 // Simplified value

	log.Printf("Massa del sole: %e kg", solarMass)

	sun := body.NewRigidBody(
		units.NewQuantity(solarMass, units.Kilogram),
		units.NewQuantity(5.0, units.Meter),                        // Raggio del sole (scalato)
		vector.NewVector3(0, 0, 0),                                 // Posizione al centro
		vector.NewVector3(0, 0, 0),                                 // Velocità zero
		createMaterial("Sun", 0.9, 0.5, [3]float64{1.0, 0.8, 0.0}), // Yellow color
	)
	sun.SetStatic(false) // The sun is static (does not move)
	w.AddBody(sun)
	log.Printf("Sun created: ID=%v, Position=%v", sun.ID(), sun.Position())

	// Crea i pianeti
	log.Println("Creating the planets")

	// Definiamo le distanze dei pianeti
	// Simply increasing progressively
	distances := []float64{20, 30, 40, 50, 70, 90, 110, 130}

	// Planet names
	names := []string{"Mercury", "Venus", "Earth", "Mars", "Jupiter", "Saturn", "Uranus", "Neptune"}

	// Planet radii (scaled)
	radii := []float64{0.8, 1.2, 1.3, 1.0, 2.5, 2.2, 1.8, 1.8}

	// Planet masses as fractions of the solar mass
	// The exact values are not important, what matters is that they are much smaller than the sun
	massFractions := []float64{1e-7, 2e-7, 2e-7, 1e-7, 1e-6, 9e-7, 4e-7, 5e-7}

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

	// Crea ogni pianeta
	for i := 0; i < len(names); i++ {
		// Crea il pianeta - la velocità orbitale verrà calcolata all'interno della funzione
		createPlanet(
			w,
			names[i],                   // Nome
			solarMass*massFractions[i], // Massa
			radii[i],                   // Raggio
			distances[i],               // Distanza
			solarMass,                  // Massa dell'oggetto centrale (sole)
			vector.NewVector3(0, 1, 0), // Piano dell'orbita
			colors[i],                  // Colore
		)
	}

	// Create an asteroid belt
	log.Println("Creating the asteroid belt")
	createAsteroidBelt(w, 200, solarMass, 60.0, 80.0) // Passiamo la massa solare come massa centrale
}

// createPlanet creates a planet
func createPlanet(w world.World, name string, mass, radius, distance, centralMass float64, orbitPlane vector.Vector3, color [3]float64) body.Body {
	// Calcola la velocità orbitale usando la funzione agnostica del package force
	orbitSpeed := force.CalculateOrbitalVelocity(centralMass, distance)

	log.Printf("Creating planet %s: distance=%f, radius=%f, calculated orbit speed=%f", name, distance, radius, orbitSpeed)

	// Random angle for the initial position (to distribute planets around the sun)
	angle := rand.Float64() * 2 * math.Pi

	// Calculate the initial position
	position := vector.NewVector3(
		distance*math.Cos(angle),
		0,
		distance*math.Sin(angle),
	)

	// Calculate the orbital velocity (perpendicular to the position)
	// This is the key for stable orbits: velocity must be perpendicular to the radius
	velocity := vector.NewVector3(
		-orbitSpeed*math.Sin(angle), // Componente x
		0,                           // Componente y (piano xy)
		orbitSpeed*math.Cos(angle),  // Componente z
	)

	// Create the planet
	planet := body.NewRigidBody(
		units.NewQuantity(mass, units.Kilogram),
		units.NewQuantity(radius, units.Meter),
		position,
		velocity,
		createMaterial(name, 0.7, 0.5, color),
	)

	// Add the planet to the world
	w.AddBody(planet)
	log.Printf("Planet %s added: ID=%v, Position=%v, Velocity=%v", name, planet.ID(), planet.Position(), planet.Velocity())

	return planet
}

// createAsteroidBelt creates an asteroid belt
func createAsteroidBelt(w world.World, count int, centralMass, minDistance, maxDistance float64) {
	log.Printf("Creating %d asteroids", count)

	for i := 0; i < count; i++ {
		// Generate a random position in the asteroid belt
		distance := minDistance + rand.Float64()*(maxDistance-minDistance)
		angle := rand.Float64() * 2 * math.Pi

		x := distance * math.Cos(angle)
		z := distance * math.Sin(angle)
		y := (rand.Float64()*2 - 1) * 5 // Wider vertical distribution

		position := vector.NewVector3(x, y, z)

		// Calculate the orbital velocity using the agnostic function from the force package
		baseSpeed := force.CalculateOrbitalVelocity(centralMass, distance)

		// Add a small random variation to the velocity
		speed := baseSpeed //* (0.95 + rand.Float64()*0.1) // 95-105% of the base speed

		// The velocity must be perpendicular to the radius for a circular orbit
		// For an asteroid with y-axis ≠ 0, we need to calculate the perpendicular vector correctly
		radialDirection := position.Normalize()

		// "Up" vector in the y-axis
		up := vector.NewVector3(0, 1, 0)

		// Get the perpendicular vector by doing the cross product
		tangentialDirection := up.Cross(radialDirection).Normalize()

		// If the result is almost zero (asteroid almost on the y-axis), use another direction
		if tangentialDirection.Length() < 0.1 {
			tangentialDirection = vector.NewVector3(1, 0, 0).Cross(radialDirection).Normalize()
		}

		velocity := tangentialDirection.Scale(speed)

		// Create the asteroid with reduced mass
		asteroid := body.NewRigidBody(
			units.NewQuantity(rand.Float64()*10+1, units.Kilogram), // Much smaller mass than planets
			units.NewQuantity(rand.Float64()*0.3+0.1, units.Meter), // Smaller size
			position,
			velocity,
			physMaterial.Rock,
		)

		// Add the asteroid to the world
		w.AddBody(asteroid)
	}

	log.Printf("Asteroid belt created")
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
