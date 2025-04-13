// Package main provides an example of simulation with bodies arranged in predefined geometric shapes
package main

import (
	"log"
	"math"
	"math/rand"
	"os"
	"runtime/pprof"
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

const (
	shouldBeProfiled = true
)

func main() {

	log.Println("Initializing simulation with bodies in geometric shapes with gravitational interaction")

	if shouldBeProfiled {
		f, err := os.Create("cpu.pprof")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// Create the simulation configuration
	cfg := config.NewSimulationBuilder().
		WithTimeStep(0.01).
		WithMaxBodies(5000).
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

	// Add the gravitational force
	gravityForce := force.NewGravitationalForce()
	gravityForce.SetTheta(0.5) // Set the theta value for the Barnes-Hut algorithm
	w.AddForce(gravityForce)

	// Choose one of the following configurations
	// You can comment/uncomment the desired one

	// Creation of a massive central body (optional)
	// createCentralBody(w)

	// Create a cube of bodies (non-static, influenced by gravity)
	// createCuboidFormation(w, 512, 50.0, 100.0, 5.0)

	// Other available formations:
	createSphereFormation(w, 300, 40.0, 80.0, 5.0)
	// createRingFormation(w, 200, 60.0, 80.0, 5.0)
	// createSpiralFormation(w, 200, 20.0, 80.0, 5.0, 3)
	// createBinarySystem(w, 200, 5.0)

	// Create the direct G3N adapter
	adapter := g3n.NewG3NAdapter()

	// Configure the adapter
	adapter.SetBackgroundColor(g3n.NewColor(0.0, 0.0, 0.1, 1.0)) // Very dark blue background for space

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

	log.Println("Simulation completed")
}

// createCentralBody creates a massive central body
func createCentralBody(w world.World) {
	log.Println("Creating the central body")

	// High mass for the central body
	centralMass := 1.5e11

	centralBody := body.NewRigidBody(
		units.NewQuantity(centralMass, units.Kilogram),
		units.NewQuantity(8.0, units.Meter),
		vector.NewVector3(0, 0, 0),
		vector.NewVector3(0, 0, 0),
		createMaterial("CentralBody", 0.9, 0.5, [3]float64{1.0, 0.6, 0.0}),
	)

	// The central body can be static or dynamic
	// If you want EVERYTHING to be influenced by gravity, comment out the following line
	centralBody.SetStatic(true)

	w.AddBody(centralBody)
	log.Printf("Central body created: ID=%v, Position=%v", centralBody.ID(), centralBody.Position())
}

// createCuboidFormation creates bodies distributed in a cube
func createCuboidFormation(w world.World, count int, minSize, maxSize, minDistance float64) {
	log.Printf("Creating %d bodies in cubic formation", count)

	// Mass of the central body (if present)
	centralMass := 1.5e11

	// Determine the number of bodies per side to get a perfect cube
	// Calculate the cubic root rounded to the nearest integer
	bodiesPerSide := int(math.Ceil(math.Pow(float64(count), 1.0/3.0)))
	actualCount := bodiesPerSide * bodiesPerSide * bodiesPerSide

	log.Printf("Creating a cube %dx%dx%d with %d total bodies",
		bodiesPerSide, bodiesPerSide, bodiesPerSide, actualCount)

	// Calculate the spacing between bodies
	spacing := minDistance + 1.0 // Ensure a minimum distance between bodies

	// Determine the total size of the cube
	cubeSize := float64(bodiesPerSide-1) * spacing
	halfSize := cubeSize / 2.0

	// Positions of already created bodies
	positions := make([]vector.Vector3, 0, actualCount)

	// Create the cubic lattice
	for x := 0; x < bodiesPerSide; x++ {
		for y := 0; y < bodiesPerSide; y++ {
			for z := 0; z < bodiesPerSide; z++ {
				// Calculate the position in the lattice
				posX := float64(x)*spacing - halfSize
				posY := float64(y)*spacing - halfSize
				posZ := float64(z)*spacing - halfSize

				position := vector.NewVector3(posX, posY, posZ)
				positions = append(positions, position)

				// Create a body with random mass but not too large
				bodyMass := (rand.Float64()*20 + 5.0) * 1e9 // Mass between 5 and 25 * 10^9

				// Calculate the distance from the center (0,0,0) where the central body is
				distanceFromCenter := position.Length()

				// Calculate a proper orbital velocity using the utility function
				// This will give a more realistic orbital motion
				orbitSpeed := 0.0
				if distanceFromCenter > 0 {
					orbitSpeed = force.CalculateOrbitalVelocity(centralMass, distanceFromCenter)

					// Apply a scaling factor to make the simulation visually appealing
					// We reduce the speed to create a more stable initial state
					orbitSpeed *= 0.3
				}

				// Calculate the direction perpendicular to the radius
				radialDirection := position.Normalize()

				// Choose a reference vector that is not parallel to the radial direction
				reference := vector.NewVector3(0, 1, 0)
				if math.Abs(radialDirection.Dot(reference)) > 0.9 {
					reference = vector.NewVector3(1, 0, 0)
				}

				// Calculate the perpendicular direction for orbital motion
				tangentDirection := radialDirection.Cross(reference).Normalize()

				// Calculate the velocity vector
				velocity := tangentDirection.Scale(orbitSpeed)

				// Create the body
				newBody := body.NewRigidBody(
					units.NewQuantity(bodyMass, units.Kilogram),
					units.NewQuantity(rand.Float64()*0.5+0.5, units.Meter), // Random radius
					position,
					velocity,
					createRandomMaterial(),
				)

				// Important: DO NOT set the body as static
				// newBody.SetStatic(false) - this is the default behavior

				w.AddBody(newBody)
			}
		}
	}

	log.Printf("Cubic formation created with %d bodies", len(positions))
}

// createSphereFormation creates bodies distributed in a sphere
func createSphereFormation(w world.World, count int, minRadius, maxRadius, minDistance float64) {
	log.Printf("Creating %d bodies in spherical formation", count)

	// Mass of the central body
	centralMass := 1.5e11

	positions := make([]vector.Vector3, 0, count)

	for i := 0; i < count; i++ {
		// Generate a point on the sphere with uniform distribution
		phi := rand.Float64() * 2 * math.Pi
		costheta := rand.Float64()*2 - 1
		theta := math.Acos(costheta)

		// Random radius between minRadius and maxRadius
		radius := minRadius + rand.Float64()*(maxRadius-minRadius)

		x := radius * math.Sin(theta) * math.Cos(phi)
		y := radius * math.Sin(theta) * math.Sin(phi)
		z := radius * math.Cos(theta)

		position := vector.NewVector3(x, y, z)

		// Check minimum distance from other bodies
		tooClose := false
		for _, pos := range positions {
			if position.Distance(pos) < minDistance {
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

		// Calculate the orbital velocity using the utility function
		orbitSpeed := force.CalculateOrbitalVelocity(centralMass, radius) * 0.8 // 80% of the theoretical orbital velocity

		// Calculate velocity direction with some randomness
		radialDirection := position.Normalize()

		// Choose a reference vector
		reference := vector.NewVector3(0, 1, 0)
		if math.Abs(radialDirection.Dot(reference)) > 0.9 {
			reference = vector.NewVector3(1, 0, 0)
		}

		// Calculate the perpendicular vector
		tangent := radialDirection.Cross(reference).Normalize()

		// Add random component to the velocity
		randomFactor := 0.2 // 20% randomness
		velX := tangent.X() * orbitSpeed * (1.0 + (rand.Float64()*2-1)*randomFactor)
		velY := tangent.Y() * orbitSpeed * (1.0 + (rand.Float64()*2-1)*randomFactor)
		velZ := tangent.Z() * orbitSpeed * (1.0 + (rand.Float64()*2-1)*randomFactor)

		velocity := vector.NewVector3(velX, velY, velZ)

		// Create the body
		newBody := body.NewRigidBody(
			units.NewQuantity(rand.Float64()*50+10, units.Kilogram),
			units.NewQuantity(rand.Float64()*0.5+0.5, units.Meter),
			position,
			velocity,
			createRandomMaterial(),
		)

		w.AddBody(newBody)
	}

	log.Printf("Spherical formation created with %d bodies", len(positions))
}

// createRingFormation creates bodies distributed in a ring
func createRingFormation(w world.World, count int, minRadius, maxRadius, minDistance float64) {
	log.Printf("Creating %d bodies in ring formation", count)

	// Mass of the central body
	centralMass := 1.5e11

	positions := make([]vector.Vector3, 0, count)

	for i := 0; i < count; i++ {
		// Generate a point in a ring
		angle := rand.Float64() * 2 * math.Pi
		radius := minRadius + rand.Float64()*(maxRadius-minRadius)

		// Generate a small vertical variation
		height := (rand.Float64()*2 - 1) * (maxRadius - minRadius) * 0.1

		x := radius * math.Cos(angle)
		y := height
		z := radius * math.Sin(angle)

		position := vector.NewVector3(x, y, z)

		// Check minimum distance from other bodies
		tooClose := false
		for _, pos := range positions {
			if position.Distance(pos) < minDistance {
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

		// Calculate the orbital velocity using the utility function
		orbitSpeed := force.CalculateOrbitalVelocity(centralMass, radius)

		// The velocity must be perpendicular to the radius on the ring plane
		velocity := vector.NewVector3(
			-orbitSpeed*math.Sin(angle),         // X component
			(rand.Float64()*2-1)*0.1*orbitSpeed, // Small random vertical component
			orbitSpeed*math.Cos(angle),          // Z component
		)

		// Create the body
		newBody := body.NewRigidBody(
			units.NewQuantity(rand.Float64()*20+5, units.Kilogram),
			units.NewQuantity(rand.Float64()*0.5+0.3, units.Meter),
			position,
			velocity,
			createRandomMaterial(),
		)

		w.AddBody(newBody)
	}

	log.Printf("Ring formation created with %d bodies", len(positions))
}

// createSpiralFormation creates bodies distributed in a spiral
func createSpiralFormation(w world.World, count int, minRadius, maxRadius, minDistance float64, arms int) {
	log.Printf("Creating %d bodies in spiral formation with %d arms", count, arms)

	// Mass of the central body
	centralMass := 1.5e11

	positions := make([]vector.Vector3, 0, count)

	// Spiral parameters
	turns := 2.0 // Number of complete turns of the spiral

	for i := 0; i < count; i++ {
		// Choose a random arm
		arm := rand.Intn(arms)

		// Parameter t varies from 0 to 1 along the spiral
		t := rand.Float64()

		// Base angle for this arm of the spiral
		baseAngle := 2.0 * math.Pi * float64(arm) / float64(arms)

		// Angle that increases with t
		angle := baseAngle + turns*2.0*math.Pi*t

		// The radius increases with t
		radius := minRadius + t*(maxRadius-minRadius)

		// Add some variation to the radius
		radius += (rand.Float64()*2 - 1) * (maxRadius - minRadius) * 0.05

		// X, Y, Z coordinates
		x := radius * math.Cos(angle)
		y := (rand.Float64()*2 - 1) * maxRadius * 0.05 // Small variation on the y-axis
		z := radius * math.Sin(angle)

		position := vector.NewVector3(x, y, z)

		// Check minimum distance from other bodies
		tooClose := false
		for _, pos := range positions {
			if position.Distance(pos) < minDistance {
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

		// Calculate the orbital velocity using the utility function
		orbitSpeed := force.CalculateOrbitalVelocity(centralMass, radius) * 0.9

		// The velocity must be perpendicular to the radius
		velocity := vector.NewVector3(
			-orbitSpeed*math.Sin(angle),
			(rand.Float64()*2-1)*0.05*orbitSpeed, // Small vertical component
			orbitSpeed*math.Cos(angle),
		)

		// Create the body
		newBody := body.NewRigidBody(
			units.NewQuantity(rand.Float64()*15+5, units.Kilogram),
			units.NewQuantity(rand.Float64()*0.4+0.3, units.Meter),
			position,
			velocity,
			createRandomMaterial(),
		)

		w.AddBody(newBody)
	}

	log.Printf("Spiral formation created with %d bodies", len(positions))
}

// createBinarySystem creates a binary system with two massive central bodies
func createBinarySystem(w world.World, satelliteCount int, minDistance float64) {
	log.Println("Creating a binary system")

	// Create two massive bodies
	mass1 := 7.5e10
	mass2 := 5.0e10
	separation := 40.0

	// Calculation of the orbital velocity for the two central bodies
	// We assume that the bodies orbit around their center of mass
	// totalMass := mass1 + mass2

	// Position of the center of mass
	// centerOfMassX := (mass1*(-separation/2) + mass2*(separation/2)) / totalMass

	// Effective distance of each body from the center of mass
	// dist1 := math.Abs((-separation / 2) - centerOfMassX)
	// dist2 := math.Abs((separation / 2) - centerOfMassX)

	// Using the utility function to calculate orbital velocities
	speed1 := force.CalculateOrbitalVelocity(mass2, separation)
	speed2 := force.CalculateOrbitalVelocity(mass1, separation)

	// Creation of the first central body
	body1 := body.NewRigidBody(
		units.NewQuantity(mass1, units.Kilogram),
		units.NewQuantity(5.0, units.Meter),
		vector.NewVector3(-separation/2, 0, 0),
		vector.NewVector3(0, 0, speed1),
		createMaterial("CentralBody1", 0.9, 0.5, [3]float64{0.9, 0.6, 0.1}),
	)
	w.AddBody(body1)

	// Creation of the second central body
	body2 := body.NewRigidBody(
		units.NewQuantity(mass2, units.Kilogram),
		units.NewQuantity(4.0, units.Meter),
		vector.NewVector3(separation/2, 0, 0),
		vector.NewVector3(0, 0, -speed2),
		createMaterial("CentralBody2", 0.9, 0.5, [3]float64{0.2, 0.6, 0.9}),
	)
	w.AddBody(body2)

	log.Println("Central bodies of the binary system created")

	// Create satellites around the binary system
	createSatellites(w, satelliteCount, separation*1.5, separation*5, minDistance, mass1+mass2)
}

// createSatellites creates satellites around a central point
func createSatellites(w world.World, count int, minRadius, maxRadius, minDistance, centralMass float64) {
	log.Printf("Creating %d satellites", count)

	positions := make([]vector.Vector3, 0, count)

	for i := 0; i < count; i++ {
		// Random position in a sphere
		phi := rand.Float64() * 2 * math.Pi
		costheta := rand.Float64()*2 - 1
		theta := math.Acos(costheta)

		radius := minRadius + rand.Float64()*(maxRadius-minRadius)

		x := radius * math.Sin(theta) * math.Cos(phi)
		y := radius * math.Sin(theta) * math.Sin(phi)
		z := radius * math.Cos(theta)

		position := vector.NewVector3(x, y, z)

		// Check minimum distance from other bodies
		tooClose := false
		for _, pos := range positions {
			if position.Distance(pos) < minDistance {
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

		// Calculation of the orbital velocity using the utility function
		orbitSpeed := force.CalculateOrbitalVelocity(centralMass, radius)

		// Calculate velocity direction
		radialDirection := position.Normalize()

		// Reference vector
		reference := vector.NewVector3(0, 1, 0)
		if math.Abs(radialDirection.Dot(reference)) > 0.9 {
			reference = vector.NewVector3(1, 0, 0)
		}

		// Calculate the perpendicular vector
		tangent := reference.Cross(radialDirection).Normalize()
		velocity := tangent.Scale(orbitSpeed)

		// Create the satellite
		satellite := body.NewRigidBody(
			units.NewQuantity(rand.Float64()*10+1, units.Kilogram),
			units.NewQuantity(rand.Float64()*0.4+0.2, units.Meter),
			position,
			velocity,
			createRandomMaterial(),
		)

		w.AddBody(satellite)
	}

	log.Printf("Satellites created: %d", len(positions))
}

// createRandomMaterial creates a material with random color
func createRandomMaterial() physMaterial.Material {
	// Generate a random color
	r := rand.Float64()*0.7 + 0.3
	g := rand.Float64()*0.7 + 0.3
	b := rand.Float64()*0.7 + 0.3

	return physMaterial.NewBasicMaterial(
		"RandomMaterial",
		units.NewQuantity(5000, units.Kilogram),
		units.NewQuantity(800, units.Joule),
		units.NewQuantity(1.5, units.Watt),
		0.7+rand.Float64()*0.3, // Emissivity between 0.7 and 1.0
		0.3+rand.Float64()*0.6, // Elasticity between 0.3 and 0.9
		[3]float64{r, g, b},    // Random color
	)
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
