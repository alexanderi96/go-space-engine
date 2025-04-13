// Package celestial provides generators for celestial bodies and formations
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

// FormationParams defines parameters for celestial formation
type FormationParams struct {
	// Core parameters
	Count       int     // Number of bodies to create
	MinDistance float64 // Minimum distance between bodies

	// Size/Position parameters
	MinRadius float64 // Inner radius of the formation
	MaxRadius float64 // Outer radius of the formation

	// Mass parameters
	MinMass float64 // Minimum mass
	MaxMass float64 // Maximum mass

	// Size parameters
	MinSize float64 // Minimum physical size
	MaxSize float64 // Maximum physical size

	// Formation shape parameters
	Arms   int     // Number of arms (for spiral)
	Turns  float64 // Number of turns (for spiral)
	Height float64 // Maximum height (for ring/disk)
	Orbits bool    // Generate orbital velocities

	// Reference
	CentralBody body.Body // Central body for orbital calculations
}

// DefaultFormationParams returns default formation parameters
func DefaultFormationParams() FormationParams {
	return FormationParams{
		Count:       100,
		MinDistance: 2.0,
		MinRadius:   20.0,
		MaxRadius:   40.0,
		MinMass:     1e5,
		MaxMass:     1e6,
		MinSize:     0.5,
		MaxSize:     1.5,
		Arms:        2,
		Turns:       1.5,
		Height:      1.0,
		Orbits:      true,
		CentralBody: nil,
	}
}

// CreateRingFormation creates bodies distributed in a ring around a central body
func CreateRingFormation(w world.World, params FormationParams) []body.Body {
	// If no central body is provided, we can't calculate orbital velocities
	centralMass := 0.0
	if params.CentralBody != nil {
		centralMass = params.CentralBody.Mass().Value()
	}

	bodies := make([]body.Body, 0, params.Count)
	positions := make([]vector.Vector3, 0, params.Count)

	for i := 0; i < params.Count; i++ {
		// Generate a point in a ring
		angle := rand.Float64() * 2 * math.Pi
		radius := params.MinRadius + rand.Float64()*(params.MaxRadius-params.MinRadius)

		// Generate a small vertical variation
		height := (rand.Float64()*2 - 1) * params.Height

		x := radius * math.Cos(angle)
		y := height
		z := radius * math.Sin(angle)

		position := vector.NewVector3(x, y, z)

		// Check minimum distance from other bodies
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

		// Calculate mass
		bodyMass := params.MinMass
		if params.MaxMass > params.MinMass {
			bodyMass = params.MinMass + rand.Float64()*(params.MaxMass-params.MinMass)
		}

		// Calculate size
		bodySize := params.MinSize
		if params.MaxSize > params.MinSize {
			bodySize = params.MinSize + rand.Float64()*(params.MaxSize-params.MinSize)
		}

		// Calculate orbital velocity if needed
		velocity := vector.Zero3()
		if params.Orbits && centralMass > 0 {
			orbitSpeed := force.CalculateOrbitalVelocity(centralMass, radius)

			// The velocity must be perpendicular to the radius on the ring plane
			velocity = vector.NewVector3(
				-orbitSpeed*math.Sin(angle),         // X component
				(rand.Float64()*2-1)*0.1*orbitSpeed, // Small random vertical component
				orbitSpeed*math.Cos(angle),          // Z component
			)
		}

		// Create the body
		b := body.NewRigidBody(
			units.NewQuantity(bodyMass, units.Kilogram),
			units.NewQuantity(bodySize, units.Meter),
			position,
			velocity,
			material.Rock, // Default material
		)

		bodies = append(bodies, b)
		w.AddBody(b)
	}

	return bodies
}

// CreateDiskFormation creates bodies distributed in a disk around a central body
// This is similar to ring but with bodies distributed throughout the disk, not just the edge
func CreateDiskFormation(w world.World, params FormationParams) []body.Body {
	// If no central body is provided, we can't calculate orbital velocities
	centralMass := 0.0
	if params.CentralBody != nil {
		centralMass = params.CentralBody.Mass().Value()
	}

	bodies := make([]body.Body, 0, params.Count)
	positions := make([]vector.Vector3, 0, params.Count)

	for i := 0; i < params.Count; i++ {
		// Generate a point in a disk
		// Square root distribution for uniform area density
		radius := params.MinRadius + math.Sqrt(rand.Float64())*(params.MaxRadius-params.MinRadius)
		angle := rand.Float64() * 2 * math.Pi

		// Generate a small vertical variation
		height := (rand.Float64()*2 - 1) * params.Height

		x := radius * math.Cos(angle)
		y := height
		z := radius * math.Sin(angle)

		position := vector.NewVector3(x, y, z)

		// Check minimum distance from other bodies
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

		// Calculate mass
		bodyMass := params.MinMass
		if params.MaxMass > params.MinMass {
			bodyMass = params.MinMass + rand.Float64()*(params.MaxMass-params.MinMass)
		}

		// Calculate size
		bodySize := params.MinSize
		if params.MaxSize > params.MinSize {
			bodySize = params.MinSize + rand.Float64()*(params.MaxSize-params.MinSize)
		}

		// Calculate orbital velocity if needed
		velocity := vector.Zero3()
		if params.Orbits && centralMass > 0 {
			orbitSpeed := force.CalculateOrbitalVelocity(centralMass, radius)

			// The velocity must be perpendicular to the radius on the disk plane
			velocity = vector.NewVector3(
				-orbitSpeed*math.Sin(angle),         // X component
				(rand.Float64()*2-1)*0.1*orbitSpeed, // Small random vertical component
				orbitSpeed*math.Cos(angle),          // Z component
			)
		}

		// Create the body
		b := body.NewRigidBody(
			units.NewQuantity(bodyMass, units.Kilogram),
			units.NewQuantity(bodySize, units.Meter),
			position,
			velocity,
			material.Rock, // Default material
		)

		bodies = append(bodies, b)
		w.AddBody(b)
	}

	return bodies
}

// CreateSphereFormation creates bodies distributed in a sphere around a central body
func CreateSphereFormation(w world.World, params FormationParams) []body.Body {
	// If no central body is provided, we can't calculate orbital velocities
	centralMass := 0.0
	if params.CentralBody != nil {
		centralMass = params.CentralBody.Mass().Value()
	}

	bodies := make([]body.Body, 0, params.Count)
	positions := make([]vector.Vector3, 0, params.Count)

	for i := 0; i < params.Count; i++ {
		// Generate a point on the sphere with uniform distribution
		phi := rand.Float64() * 2 * math.Pi
		costheta := rand.Float64()*2 - 1
		theta := math.Acos(costheta)

		// Random radius between minRadius and maxRadius
		radius := params.MinRadius + rand.Float64()*(params.MaxRadius-params.MinRadius)

		x := radius * math.Sin(theta) * math.Cos(phi)
		y := radius * math.Sin(theta) * math.Sin(phi)
		z := radius * math.Cos(theta)

		position := vector.NewVector3(x, y, z)

		// Check minimum distance from other bodies
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

		// Calculate mass
		bodyMass := params.MinMass
		if params.MaxMass > params.MinMass {
			bodyMass = params.MinMass + rand.Float64()*(params.MaxMass-params.MinMass)
		}

		// Calculate size
		bodySize := params.MinSize
		if params.MaxSize > params.MinSize {
			bodySize = params.MinSize + rand.Float64()*(params.MaxSize-params.MinSize)
		}

		// Calculate orbital velocity if needed
		velocity := vector.Zero3()
		if params.Orbits && centralMass > 0 {
			// Calculate the orbital velocity using the agnostic function
			baseSpeed := force.CalculateOrbitalVelocity(centralMass, radius) * 0.8 // 80% of theoretical for stability

			// Calculate velocity direction with some randomness
			radialDirection := position.Normalize()

			// Choose a reference vector
			reference := vector.NewVector3(0, 1, 0)
			if math.Abs(radialDirection.Dot(reference)) > 0.9 {
				reference = vector.NewVector3(1, 0, 0)
			}

			// Calculate the perpendicular vector
			tangent := reference.Cross(radialDirection).Normalize()

			// Add random component to the velocity
			randomFactor := 0.2 // 20% randomness
			velX := tangent.X() * baseSpeed * (1.0 + (rand.Float64()*2-1)*randomFactor)
			velY := tangent.Y() * baseSpeed * (1.0 + (rand.Float64()*2-1)*randomFactor)
			velZ := tangent.Z() * baseSpeed * (1.0 + (rand.Float64()*2-1)*randomFactor)

			velocity = vector.NewVector3(velX, velY, velZ)
		}

		// Create the body
		b := body.NewRigidBody(
			units.NewQuantity(bodyMass, units.Kilogram),
			units.NewQuantity(bodySize, units.Meter),
			position,
			velocity,
			material.Rock, // Default material
		)

		bodies = append(bodies, b)
		w.AddBody(b)
	}

	return bodies
}

// CreateSpiralFormation creates bodies distributed in a spiral around a central body
func CreateSpiralFormation(w world.World, params FormationParams) []body.Body {
	// If no central body is provided, we can't calculate orbital velocities
	centralMass := 0.0
	if params.CentralBody != nil {
		centralMass = params.CentralBody.Mass().Value()
	}

	bodies := make([]body.Body, 0, params.Count)
	positions := make([]vector.Vector3, 0, params.Count)

	// Default number of arms if not specified
	arms := params.Arms
	if arms <= 0 {
		arms = 2
	}

	turns := params.Turns
	if turns <= 0 {
		turns = 2.0
	}

	for i := 0; i < params.Count; i++ {
		// Choose a random arm
		arm := rand.Intn(arms)

		// Parameter t varies from 0 to 1 along the spiral
		t := rand.Float64()

		// Base angle for this arm of the spiral
		baseAngle := 2.0 * math.Pi * float64(arm) / float64(arms)

		// Angle that increases with t
		angle := baseAngle + turns*2.0*math.Pi*t

		// The radius increases with t
		radius := params.MinRadius + t*(params.MaxRadius-params.MinRadius)

		// Add some variation to the radius
		radius += (rand.Float64()*2 - 1) * (params.MaxRadius - params.MinRadius) * 0.05

		// X, Y, Z coordinates
		x := radius * math.Cos(angle)
		y := (rand.Float64()*2 - 1) * params.Height // Variation on the y-axis
		z := radius * math.Sin(angle)

		position := vector.NewVector3(x, y, z)

		// Check minimum distance from other bodies
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

		// Calculate mass
		bodyMass := params.MinMass
		if params.MaxMass > params.MinMass {
			bodyMass = params.MinMass + rand.Float64()*(params.MaxMass-params.MinMass)
		}

		// Calculate size
		bodySize := params.MinSize
		if params.MaxSize > params.MinSize {
			bodySize = params.MinSize + rand.Float64()*(params.MaxSize-params.MinSize)
		}

		// Calculate orbital velocity if needed
		velocity := vector.Zero3()
		if params.Orbits && centralMass > 0 {
			// Calculate the orbital velocity
			orbitSpeed := force.CalculateOrbitalVelocity(centralMass, radius) * 0.9

			// The velocity must be perpendicular to the radius
			velocity = vector.NewVector3(
				-orbitSpeed*math.Sin(angle),
				(rand.Float64()*2-1)*0.05*orbitSpeed, // Small vertical component
				orbitSpeed*math.Cos(angle),
			)
		}

		// Create the body
		b := body.NewRigidBody(
			units.NewQuantity(bodyMass, units.Kilogram),
			units.NewQuantity(bodySize, units.Meter),
			position,
			velocity,
			material.Rock, // Default material
		)

		bodies = append(bodies, b)
		w.AddBody(b)
	}

	return bodies
}

// CreateCubeFormation creates bodies distributed in a cube
func CreateCubeFormation(w world.World, params FormationParams) []body.Body {
	// If no central body is provided, we can't calculate orbital velocities
	centralMass := 0.0
	if params.CentralBody != nil {
		centralMass = params.CentralBody.Mass().Value()
	}

	bodies := make([]body.Body, 0, params.Count)

	// Calculate the number of bodies per side for a cube
	bodiesPerSide := int(math.Ceil(math.Pow(float64(params.Count), 1.0/3.0)))

	// Calculate the spacing between bodies
	spacing := params.MinDistance + params.MaxSize // Use max size to ensure spacing

	// Calculate the total size of the cube
	cubeSize := float64(bodiesPerSide-1) * spacing
	halfSize := cubeSize / 2.0

	// Create the cubic lattice
	for x := 0; x < bodiesPerSide; x++ {
		for y := 0; y < bodiesPerSide; y++ {
			for z := 0; z < bodiesPerSide; z++ {
				// Skip if we've reached the count
				if len(bodies) >= params.Count {
					break
				}

				// Calculate the position in the lattice
				posX := float64(x)*spacing - halfSize
				posY := float64(y)*spacing - halfSize
				posZ := float64(z)*spacing - halfSize

				position := vector.NewVector3(posX, posY, posZ)

				// Calculate mass
				bodyMass := params.MinMass
				if params.MaxMass > params.MinMass {
					bodyMass = params.MinMass + rand.Float64()*(params.MaxMass-params.MinMass)
				}

				// Calculate size
				bodySize := params.MinSize
				if params.MaxSize > params.MinSize {
					bodySize = params.MinSize + rand.Float64()*(params.MaxSize-params.MinSize)
				}

				// Calculate orbital velocity if needed
				velocity := vector.Zero3()
				if params.Orbits && centralMass > 0 {
					// Calculate the distance from center for orbital calculations
					distanceFromCenter := position.Length()
					if distanceFromCenter > 0 {
						// Calculate the orbital velocity
						orbitSpeed := force.CalculateOrbitalVelocity(centralMass, distanceFromCenter)
						orbitSpeed *= 0.3 // Scale down for stability

						// Calculate direction perpendicular to radius
						radialDirection := position.Normalize()

						// Choose a reference vector
						reference := vector.NewVector3(0, 1, 0)
						if math.Abs(radialDirection.Dot(reference)) > 0.9 {
							reference = vector.NewVector3(1, 0, 0)
						}

						// Calculate the perpendicular vector
						tangent := radialDirection.Cross(reference).Normalize()
						velocity = tangent.Scale(orbitSpeed)
					}
				}

				// Create the body
				b := body.NewRigidBody(
					units.NewQuantity(bodyMass, units.Kilogram),
					units.NewQuantity(bodySize, units.Meter),
					position,
					velocity,
					material.Rock, // Default material
				)

				bodies = append(bodies, b)
				w.AddBody(b)
			}
		}
	}

	return bodies
}
