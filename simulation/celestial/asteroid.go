package celestial

import (
	"math"
	"math/rand"

	"github.com/alexanderi96/go-space-engine/core/units"
	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/physics/body"
	"github.com/alexanderi96/go-space-engine/physics/force"
	"github.com/alexanderi96/go-space-engine/physics/material"
	"github.com/alexanderi96/go-space-engine/simulation/world"
)

// AsteroidParams defines parameters for creating asteroid fields
type AsteroidParams struct {
	// Base parameters
	Count       int     // Number of asteroids to create
	MinDistance float64 // Minimum distance between asteroids

	// Belt parameters
	InnerRadius float64 // Inner radius of the belt
	OuterRadius float64 // Outer radius of the belt
	Height      float64 // Vertical height of the belt

	// Size and mass parameters
	MinSize float64 // Minimum asteroid size
	MaxSize float64 // Maximum asteroid size
	MinMass float64 // Minimum asteroid mass
	MaxMass float64 // Maximum asteroid mass

	// Orbital parameters
	Eccentricity float64 // Maximum orbital eccentricity
	Inclination  float64 // Maximum orbital inclination (radians)

	// Central body
	CentralBody body.Body // Central body for orbital calculations
}

// DefaultAsteroidParams returns default parameters for an asteroid belt
func DefaultAsteroidParams() AsteroidParams {
	return AsteroidParams{
		Count:        100,
		MinDistance:  1.0,
		InnerRadius:  60.0,
		OuterRadius:  80.0,
		Height:       2.0,
		MinSize:      0.2,
		MaxSize:      0.8,
		MinMass:      1e15,
		MaxMass:      1e17,
		Eccentricity: 0.1,
		Inclination:  0.2, // About 11.5 degrees
		CentralBody:  nil,
	}
}

// CreateAsteroidField creates an asteroid field with the specified parameters
func CreateAsteroidField(w world.World, params AsteroidParams) []body.Body {
	if params.CentralBody == nil {
		return nil
	}

	centralMass := params.CentralBody.Mass().Value()
	asteroids := make([]body.Body, 0, params.Count)
	positions := make([]vector.Vector3, 0, params.Count)

	for i := 0; i < params.Count; i++ {
		// Generate a random position in the belt
		// Use square root distribution for uniform area density
		radius := params.InnerRadius +
			math.Sqrt(rand.Float64())*(params.OuterRadius-params.InnerRadius)

		// Random angle and height
		angle := rand.Float64() * 2 * math.Pi
		height := (rand.Float64()*2 - 1) * params.Height

		// Basic position before orbital adjustments
		x := radius * math.Cos(angle)
		y := height
		z := radius * math.Sin(angle)

		// Apply eccentricity and inclination
		// This is a simplified approach - a proper implementation would use Keplerian elements
		eccentricity := rand.Float64() * params.Eccentricity
		inclination := (rand.Float64()*2 - 1) * params.Inclination

		// Apply inclination (rotate around X axis)
		cosInc := math.Cos(inclination)
		sinInc := math.Sin(inclination)
		yRotated := y*cosInc - z*sinInc
		zRotated := y*sinInc + z*cosInc

		position := vector.NewVector3(x, yRotated, zRotated)

		// Check minimum distance from other asteroids
		tooClose := false
		for _, pos := range positions {
			if position.Distance(pos) < params.MinDistance {
				tooClose = true
				break
			}
		}

		// If too close, try again
		if tooClose {
			i--
			continue
		}

		positions = append(positions, position)

		// Calculate asteroid size and mass
		size := params.MinSize + rand.Float64()*(params.MaxSize-params.MinSize)
		mass := params.MinMass + rand.Float64()*(params.MaxMass-params.MinMass)

		// Calculate orbital velocity
		baseSpeed := force.CalculateOrbitalVelocity(centralMass, radius)

		// Adjust for eccentricity
		// At perihelion, velocity is higher; at aphelion, it's lower
		// This is a simplified adjustment
		velocityFactor := 1.0
		if eccentricity > 0 {
			// We're not calculating true anomaly correctly here,
			// but this provides a reasonable approximation for visualization
			velocityFactor = math.Sqrt((1 + eccentricity*math.Cos(angle)) /
				(1 - eccentricity*math.Cos(angle)))
		}

		orbitSpeed := baseSpeed * velocityFactor

		// Calculate velocity perpendicular to radius in orbital plane
		velX := -orbitSpeed * math.Sin(angle)
		velZ := orbitSpeed * math.Cos(angle)

		// Apply inclination to velocity
		velY := 0.0
		if inclination != 0 {
			velY = velZ * math.Sin(inclination)
			velZ = velZ * math.Cos(inclination)
		}

		velocity := vector.NewVector3(velX, velY, velZ)

		// Create the asteroid
		asteroid := body.NewRigidBody(
			units.NewQuantity(mass, units.Kilogram),
			units.NewQuantity(size, units.Meter),
			position,
			velocity,
			material.Rock, // Use default rock material
		)

		asteroids = append(asteroids, asteroid)
		w.AddBody(asteroid)
	}

	return asteroids
}

// CreateCometCloud creates a cloud of comets around a star system (simplified Oort cloud)
func CreateCometCloud(w world.World, centralBody body.Body, count int) []body.Body {
	if centralBody == nil {
		return nil
	}

	// Create parameters for the comet cloud
	params := AsteroidParams{
		Count:        count,
		MinDistance:  10.0,  // Comets can be closer to each other
		InnerRadius:  200.0, // Start beyond the planets
		OuterRadius:  500.0, // Extend far out
		Height:       500.0, // Spherical distribution
		MinSize:      0.2,
		MaxSize:      1.0,  // Comets can be larger
		MinMass:      1e14, // Smaller than asteroids
		MaxMass:      1e16,
		Eccentricity: 0.95,        // Highly eccentric orbits
		Inclination:  math.Pi / 2, // All inclinations up to 90 degrees
		CentralBody:  centralBody,
	}

	comets := make([]body.Body, 0, count)
	positions := make([]vector.Vector3, 0, count)

	// Get central body mass
	centralMass := centralBody.Mass().Value()

	for i := 0; i < count; i++ {
		// Generate a point in the spherical shell
		phi := rand.Float64() * 2 * math.Pi
		costheta := rand.Float64()*2 - 1
		theta := math.Acos(costheta)

		// Distance with cubic distribution for more outer comets
		distanceFactor := math.Pow(rand.Float64(), 1.0/3.0)
		radius := params.InnerRadius + distanceFactor*(params.OuterRadius-params.InnerRadius)

		// Position on the sphere
		x := radius * math.Sin(theta) * math.Cos(phi)
		y := radius * math.Sin(theta) * math.Sin(phi)
		z := radius * math.Cos(theta)

		position := vector.NewVector3(x, y, z)

		// Check minimum distance from other comets
		tooClose := false
		for _, pos := range positions {
			if position.Distance(pos) < params.MinDistance {
				tooClose = true
				break
			}
		}

		// If too close, try again
		if tooClose {
			i--
			continue
		}

		positions = append(positions, position)

		// Calculate comet size and mass
		size := params.MinSize + rand.Float64()*(params.MaxSize-params.MinSize)
		mass := params.MinMass + rand.Float64()*(params.MaxMass-params.MinMass)

		// Calculate orbital parameters
		// Eccentricity: most comets have high eccentricity (close to 1)
		eccentricity := 0.7 + rand.Float64()*0.29 // 0.7 - 0.99

		// Calculate orbital velocity
		// For highly eccentric orbits, current velocity depends on where in the orbit we are
		// We'll place comets at random points in their orbits

		// Simplified velocity calculation
		baseSpeed := force.CalculateOrbitalVelocity(centralMass, radius)

		// Reduce speed for eccentric orbit
		// At aphelion (furthest point), velocity is lower
		orbitSpeed := baseSpeed * math.Sqrt((1-eccentricity)/(1+eccentricity))

		// Calculate velocity direction perpendicular to radius vector
		radialDirection := position.Normalize()

		// Choose a random perpendicular direction
		// This isn't physically accurate but provides varied orbits for visualization
		tangent := generatePerpendicularVector(radialDirection)

		velocity := tangent.Scale(orbitSpeed)

		// Create the comet with comet material
		comet := body.NewRigidBody(
			units.NewQuantity(mass, units.Kilogram),
			units.NewQuantity(size, units.Meter),
			position,
			velocity,
			createCometMaterial(),
		)

		comets = append(comets, comet)
		w.AddBody(comet)
	}

	return comets
}

// generatePerpendicularVector returns a unit vector perpendicular to the input vector
func generatePerpendicularVector(v vector.Vector3) vector.Vector3 {
	// Start with a reference vector
	reference := vector.NewVector3(0, 1, 0)

	// If v is too close to the reference, use a different reference
	if math.Abs(v.Dot(reference)) > 0.9 {
		reference = vector.NewVector3(1, 0, 0)
	}

	// Calculate the cross product to get a perpendicular vector
	perpendicular := v.Cross(reference).Normalize()

	// Randomly rotate around the radial direction
	angle := rand.Float64() * 2 * math.Pi
	cosAngle := math.Cos(angle)
	sinAngle := math.Sin(angle)

	// Get a second perpendicular vector
	perpendicular2 := v.Cross(perpendicular).Normalize()

	// Combine the two perpendicular vectors with rotation
	result := perpendicular.Scale(cosAngle).Add(perpendicular2.Scale(sinAngle))

	return result
}
