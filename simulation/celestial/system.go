package celestial

import (
	"math"
	"math/rand"

	"github.com/alexanderi96/go-space-engine/core/constants"
	"github.com/alexanderi96/go-space-engine/core/units"
	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/physics/body"
	physMaterial "github.com/alexanderi96/go-space-engine/physics/material"
	"github.com/alexanderi96/go-space-engine/simulation/world"
)

// SystemParams defines parameters for creating star systems
type SystemParams struct {
	// System parameters
	StarMass   float64 // Mass of the central star (kg)
	StarRadius float64 // Radius of the central star (m)

	// Planet parameters
	PlanetCount   int     // Number of planets
	MoonFrequency float64 // Frequency of moons (0-1)

	// Asteroid parameters
	AsteroidBelt  bool    // Whether to include an asteroid belt
	AsteroidCount int     // Number of asteroids
	AsteroidMinR  float64 // Inner radius of asteroid belt
	AsteroidMaxR  float64 // Outer radius of asteroid belt

	// Comet parameters
	CometCloud bool // Whether to include a comet cloud
	CometCount int  // Number of comets
}

// DefaultSystemParams returns default parameters for a star system
func DefaultSystemParams() SystemParams {
	return SystemParams{
		StarMass:      1.989e30, // Solar mass
		StarRadius:    695700.0, // Solar radius (km)
		PlanetCount:   8,        // Number of planets
		MoonFrequency: 0.5,      // 50% chance of moons per planet
		AsteroidBelt:  true,     // Include asteroid belt
		AsteroidCount: 200,      // Number of asteroids
		AsteroidMinR:  60.0,     // Inner radius (scaled)
		AsteroidMaxR:  80.0,     // Outer radius (scaled)
		CometCloud:    true,     // Include comet cloud
		CometCount:    50,       // Number of comets
	}
}

// CreateStarSystem creates a complete star system with planets, moons, asteroids and comets
func CreateStarSystem(w world.World, params SystemParams) (star body.Body, planets []body.Body) {
	// Create the star
	starParams := DefaultPlanetParams()
	starParams.Name = "Star"
	starParams.Type = Star
	starParams.Mass = params.StarMass
	starParams.Radius = params.StarRadius

	star = GeneratePlanet(w, starParams)

	planets = make([]body.Body, 0, params.PlanetCount)

	// Create planets with logarithmic spacing
	for i := 0; i < params.PlanetCount; i++ {
		// Calculate distance with logarithmic spacing
		distanceFactor := 1.0 + 0.3*float64(i) // Gives a reasonable progression
		distance := 20.0 * math.Pow(1.5, distanceFactor)

		// Decide planet type based on distance
		planetType := TerrestrialPlanet
		if i > 3 { // Gas giants in outer system
			if rand.Float64() > 0.5 {
				planetType = GasGiant
			} else {
				planetType = IceGiant
			}
		}

		// Generate random orbital parameters
		inclination := (rand.Float64() * 5.0) * (math.Pi / 180.0) // 0-5 degrees
		eccentricity := rand.Float64() * 0.1                      // 0-0.1

		// Create planet parameters
		planetParams := DefaultPlanetParams()
		planetParams.Name = "Planet-" + string(rune(i+65)) // A, B, C, ...
		planetParams.Type = planetType
		planetParams.Distance = distance
		planetParams.Eccentricity = eccentricity
		planetParams.Inclination = inclination
		planetParams.InitialAngle = rand.Float64() * 2.0 * math.Pi
		planetParams.CentralBody = star

		// Calculate mass based on type
		switch planetType {
		case TerrestrialPlanet:
			planetParams.Mass = 1e24 * (0.1 + rand.Float64()*1.9) // 0.1 to 2 Earth masses
		case GasGiant:
			planetParams.Mass = 1e27 * (0.1 + rand.Float64()*1.9) // 0.1 to 2 Jupiter masses
		case IceGiant:
			planetParams.Mass = 1e26 * (0.5 + rand.Float64()) // 0.5 to 1.5 Neptune masses
		}

		// Create the planet
		planet := GeneratePlanet(w, planetParams)
		planets = append(planets, planet)

		// Add moons with probability based on moon frequency
		if rand.Float64() < params.MoonFrequency {
			// More moons for larger planets
			moonCount := 0
			switch planetType {
			case TerrestrialPlanet:
				moonCount = int(rand.Float64()*2) + 1 // 1-2 moons
			case GasGiant:
				moonCount = int(rand.Float64()*10) + 3 // 3-12 moons
			case IceGiant:
				moonCount = int(rand.Float64()*5) + 1 // 1-5 moons
			}

			// Create moons
			for m := 0; m < moonCount; m++ {
				GenerateMoon(w, planet)
			}
		}
	}

	// Create asteroid belt if requested
	if params.AsteroidBelt {
		asteroidParams := DefaultAsteroidParams()
		asteroidParams.Count = params.AsteroidCount
		asteroidParams.InnerRadius = params.AsteroidMinR
		asteroidParams.OuterRadius = params.AsteroidMaxR
		asteroidParams.CentralBody = star

		CreateAsteroidField(w, asteroidParams)
	}

	// Create comet cloud if requested
	if params.CometCloud {
		CreateCometCloud(w, star, params.CometCount)
	}

	return star, planets
}

// CreateBinarySystem creates a binary star system
func CreateBinarySystem(w world.World, params SystemParams) (stars []body.Body, planets []body.Body) {
	// Create two stars
	// Split the mass between the two stars (60/40 split)
	mass1 := params.StarMass * 0.6
	mass2 := params.StarMass * 0.4

	// Calculate radii based on mass (roughly R ∝ M^0.8 for main sequence stars)
	radius1 := params.StarRadius * math.Pow(mass1/params.StarMass, 0.8)
	radius2 := params.StarRadius * math.Pow(mass2/params.StarMass, 0.8)

	// Orbital separation (scaled based on radii)
	separation := (radius1 + radius2) * 5 // Keep stars reasonably separated

	// Calculate orbital velocities
	totalMass := mass1 + mass2

	// Calculate orbital speed based on Kepler's laws
	speed := math.Sqrt(constants.G * totalMass / separation)

	// Create the stars at opposite points in their orbit
	star1Position := vector.NewVector3(-separation/2, 0, 0)
	star2Position := vector.NewVector3(separation/2, 0, 0)

	// Velocities perpendicular to the line connecting the stars
	// For a simplified circular orbit
	star1Velocity := vector.NewVector3(0, 0, speed*math.Sqrt(mass2/totalMass))
	star2Velocity := vector.NewVector3(0, 0, -speed*math.Sqrt(mass1/totalMass))

	// Create first star
	star1 := body.NewRigidBody(
		units.NewQuantity(mass1, units.Kilogram),
		units.NewQuantity(scaleRadius(radius1, Star), units.Meter),
		star1Position,
		star1Velocity,
		createStarMaterial(),
	)

	// Create second star
	star2 := body.NewRigidBody(
		units.NewQuantity(mass2, units.Kilogram),
		units.NewQuantity(scaleRadius(radius2, Star), units.Meter),
		star2Position,
		star2Velocity,
		createStarMaterial(),
	)

	// Add stars to world
	w.AddBody(star1)
	w.AddBody(star2)

	stars = []body.Body{star1, star2}
	planets = make([]body.Body, 0)

	// Create planets orbiting the center of mass
	if params.PlanetCount > 0 {
		// Create a simplified FormationParams for circumbinary planets
		formParams := DefaultFormationParams()
		formParams.Count = params.PlanetCount
		formParams.MinRadius = separation * 2 // Planets must be outside the binary orbit
		formParams.MaxRadius = separation * 5 // Extend to several times the binary separation
		formParams.MinMass = 1e24             // Earth-like masses
		formParams.MaxMass = 1e27             // Jupiter-like masses
		formParams.MinSize = 1.0
		formParams.MaxSize = 3.0
		formParams.Orbits = true

		// Set combined mass for orbital calculations
		// We'll position these bodies at the center of mass
		combinedBody := body.NewRigidBody(
			units.NewQuantity(totalMass, units.Kilogram),
			units.NewQuantity(1.0, units.Meter),
			vector.NewVector3(0, 0, 0), // Center of mass
			vector.Zero3(),
			createStarMaterial(),
		)

		// Don't add this body to the world, it's just for orbital calculations
		formParams.CentralBody = combinedBody

		// Create planets in a disk around the binary
		planets = CreateDiskFormation(w, formParams)
	}

	// Create asteroid belt if requested
	if params.AsteroidBelt {
		asteroidParams := DefaultAsteroidParams()
		asteroidParams.Count = params.AsteroidCount
		asteroidParams.InnerRadius = separation * 3
		asteroidParams.OuterRadius = separation * 4

		// Create a proxy body for the combined mass
		combinedBody := body.NewRigidBody(
			units.NewQuantity(totalMass, units.Kilogram),
			units.NewQuantity(1.0, units.Meter),
			vector.NewVector3(0, 0, 0), // Center of mass
			vector.Zero3(),
			createStarMaterial(),
		)

		asteroidParams.CentralBody = combinedBody

		CreateAsteroidField(w, asteroidParams)
	}

	return stars, planets
}

// CreateGalaxy creates a simplified galaxy model
func CreateGalaxy(w world.World, starCount int, size float64) []body.Body {
	stars := make([]body.Body, 0, starCount)

	// Create a spiral formation for stars
	formParams := DefaultFormationParams()
	formParams.Count = starCount
	formParams.MinRadius = 0.1 * size
	formParams.MaxRadius = size
	formParams.MinDistance = size / 100
	formParams.MinMass = 1e30 // Star masses
	formParams.MaxMass = 5e30
	formParams.MinSize = 2.0 // Visual sizes
	formParams.MaxSize = 5.0
	formParams.Height = size * 0.05 // Thin disk
	formParams.Arms = 4
	formParams.Turns = 1.5

	// Create a central black hole
	blackHoleMass := 1e36 // Supermassive
	blackHole := body.NewRigidBody(
		units.NewQuantity(blackHoleMass, units.Kilogram),
		units.NewQuantity(10.0, units.Meter), // Visual size
		vector.NewVector3(0, 0, 0),
		vector.Zero3(),
		createBlackHoleMaterial(),
	)

	w.AddBody(blackHole)
	formParams.CentralBody = blackHole
	formParams.Orbits = true

	// Create spiral arms of stars
	stars = CreateSpiralFormation(w, formParams)

	// Add the central black hole to the returned stars
	return append([]body.Body{blackHole}, stars...)
}

// createBlackHoleMaterial creates a material for black holes
func createBlackHoleMaterial() physMaterial.Material {
	return physMaterial.NewBasicMaterial(
		"BlackHole",
		units.NewQuantity(5000, units.Kilogram),
		units.NewQuantity(800, units.Joule),
		units.NewQuantity(0.0, units.Watt),
		0.0,                       // Zero emissivity
		0.5,                       // Medium elasticity
		[3]float64{0.1, 0.0, 0.1}, // Very dark purple
	)
}
