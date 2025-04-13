// Package main provides an example of using the celestial simulation framework
package main

import (
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/physics/force"
	"github.com/alexanderi96/go-space-engine/render/g3n"
	"github.com/alexanderi96/go-space-engine/simulation/celestial"
	"github.com/alexanderi96/go-space-engine/simulation/config"
	"github.com/alexanderi96/go-space-engine/simulation/world"
)

// Command line flags for simulation configuration
var (
	simulationType   = flag.String("type", "solar", "Type of simulation to run (solar, binary, asteroid, galaxy)")
	bodyCount        = flag.Int("bodies", 500, "Number of bodies to simulate")
	planetCount      = flag.Int("planets", 8, "Number of planets in star system")
	asteroidCount    = flag.Int("asteroids", 200, "Number of asteroids in belt")
	cometCount       = flag.Int("comets", 50, "Number of comets in cloud")
	moonFrequency    = flag.Float64("moons", 0.5, "Frequency of moons (0-1)")
	includeAsteroids = flag.Bool("belt", true, "Include asteroid belt")
	includeComets    = flag.Bool("cloud", true, "Include comet cloud")
)

func main() {
	// Parse command line flags
	flag.Parse()

	log.Println("Starting Celestial Simulator")
	log.Printf("Simulation type: %s", *simulationType)

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Create the simulation configuration
	cfg := config.NewSimulationBuilder().
		WithTimeStep(0.01).
		WithMaxBodies(10000).
		WithGravity(true).
		WithCollisions(true).
		WithBoundaryCollisions(false). // Space is infinite!
		WithWorldBounds(
			vector.NewVector3(-5000, -5000, -5000),
			vector.NewVector3(5000, 5000, 5000),
		).
		WithOctreeConfig(10, 8).
		Build()

	// Create the simulation world
	w := world.NewPhysicalWorld(cfg.GetWorldBounds())

	// Add the gravitational force
	gravityForce := force.NewGravitationalForce()
	gravityForce.SetTheta(0.5) // Set the theta value for the Barnes-Hut algorithm
	w.AddForce(gravityForce)

	// Create bodies based on simulation type
	createBodies(w, *simulationType)

	// Create the G3N adapter
	adapter := g3n.NewG3NAdapter()

	// Configure the adapter for space visualization
	adapter.SetBackgroundColor(g3n.NewColor(0.1, 0.1, 0.1, 1.0)) // Very dark blue background

	// Variables for timing
	lastUpdateTime := time.Now()
	simulationSpeed := 1.0

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

		// Apply simulation speed factor
		dt *= simulationSpeed

		// Execute a simulation step
		w.Step(dt)

		// Render the world
		adapter.RenderWorld(w)
	})

	log.Println("Simulation completed")
}

// createBodies creates bodies based on the simulation type
func createBodies(w world.World, simType string) {
	switch simType {
	case "solar":
		createSolarSystem(w)
	case "binary":
		createBinarySystem(w)
	case "asteroid":
		createAsteroidBelt(w)
	case "galaxy":
		createGalaxy(w)
	default:
		log.Printf("Unknown simulation type: %s, defaulting to solar system", simType)
		createSolarSystem(w)
	}
}

// createSolarSystem creates a solar system with planets, asteroids and comets
func createSolarSystem(w world.World) {
	log.Println("Creating solar system simulation")

	// Create system parameters
	params := celestial.DefaultSystemParams()
	params.PlanetCount = *planetCount
	params.MoonFrequency = *moonFrequency
	params.AsteroidBelt = *includeAsteroids
	params.AsteroidCount = *asteroidCount
	params.CometCloud = *includeComets
	params.CometCount = *cometCount

	// Create the star system
	star, planets := celestial.CreateStarSystem(w, params)

	log.Printf("Created solar system with 1 star (ID: %v) and %d planets", star.ID(), len(planets))
}

// createBinarySystem creates a binary star system
func createBinarySystem(w world.World) {
	log.Println("Creating binary star system simulation")

	// Create system parameters
	params := celestial.DefaultSystemParams()
	params.PlanetCount = *planetCount
	params.AsteroidBelt = *includeAsteroids
	params.AsteroidCount = *asteroidCount

	// Create the binary system
	stars, planets := celestial.CreateBinarySystem(w, params)

	log.Printf("Created binary system with %d stars and %d planets", len(stars), len(planets))
}

// createAsteroidBelt creates a standalone asteroid belt around a star
func createAsteroidBelt(w world.World) {
	log.Println("Creating asteroid belt simulation")

	// Create a single star
	starParams := celestial.DefaultPlanetParams()
	starParams.Type = celestial.Star

	star := celestial.GeneratePlanet(w, starParams)

	// Create asteroid field
	asteroidParams := celestial.DefaultAsteroidParams()
	asteroidParams.Count = *bodyCount
	asteroidParams.CentralBody = star

	// Create asteroid belt
	asteroids := celestial.CreateAsteroidField(w, asteroidParams)

	log.Printf("Created asteroid belt with %d asteroids", len(asteroids))
}

// createGalaxy creates a spiral galaxy simulation
func createGalaxy(w world.World) {
	log.Println("Creating galaxy simulation")

	// Galaxy size
	galaxySize := 500.0

	// Create galaxy
	stars := celestial.CreateGalaxy(w, *bodyCount, galaxySize)

	log.Printf("Created galaxy with %d stars", len(stars))
}
