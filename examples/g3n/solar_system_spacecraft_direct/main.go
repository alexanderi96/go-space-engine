// Package main provides an example of using G3N directly with the physics engine
// with a controllable spacecraft near Earth
package main

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/alexanderi96/go-space-engine/core/units"
	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/entity"
	"github.com/alexanderi96/go-space-engine/entity/vehicle/spacecraft"
	"github.com/alexanderi96/go-space-engine/physics/body"
	"github.com/alexanderi96/go-space-engine/physics/force"
	physMaterial "github.com/alexanderi96/go-space-engine/physics/material"
	"github.com/alexanderi96/go-space-engine/simulation/config"
	"github.com/alexanderi96/go-space-engine/simulation/world"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/window"
)

const (
	// Spacecraft controller parameters
	maxAngularVelocity = 60.0 // 60 degrees per second
	angularDamping     = 0.3  // 30% damping
	pidProportional    = 3.0  // P gain
	pidIntegral        = 0.1  // I gain
	pidDerivative      = 0.5  // D gain
	maxThrust          = 5000.0
	maxTorque          = 1000.0
)

// Global variables for spacecraft control
var (
	spacecraftEntity     entity.Entity
	spacecraftController *spacecraft.SpacecraftController
	thrustLevel          = 0.0
	earthPosition        vector.Vector3

	// G3N related variables
	application    *app.Application
	scene          *core.Node
	cam            *camera.Camera
	cameraCtrl     *camera.OrbitControl
	bodyMeshes     = make(map[string]*graphic.Mesh)
	spacecraftMesh *graphic.Mesh
)

func main() {
	log.Println("Initializing Solar System with Controllable Spacecraft (Direct G3N)")

	// Create the simulation configuration
	cfg := config.NewSimulationBuilder().
		WithTimeStep(0.01).
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

	// Add the gravitational force
	gravityForce := force.NewGravitationalForce()
	gravityForce.SetTheta(0.5) // Set the theta value for the Barnes-Hut algorithm
	w.AddForce(gravityForce)

	// Create the celestial bodies
	createBodies(w)

	// Initialize G3N
	initializeG3N()

	// Create meshes for all bodies in the world
	createMeshesForBodies(w)

	// Configure input handling
	setupInputHandling()

	// Variables for timing
	lastUpdateTime := time.Now()

	// Start the rendering loop
	application.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		// Calculate the delta time
		currentTime := time.Now()
		dt := currentTime.Sub(lastUpdateTime).Seconds()
		lastUpdateTime = currentTime

		// Limit the delta time to avoid instability
		if dt > 0.1 {
			dt = 0.1
		}

		// Update the spacecraft controller
		if spacecraftController != nil {
			// Apply the current thrust level
			spacecraftController.SetThrustLevel(thrustLevel)

			// Update the controller
			spacecraftController.Update(dt)
		}

		// Execute a simulation step
		w.Step(dt)

		// Update mesh positions based on body positions
		updateMeshPositions(w)

		// Update camera to follow spacecraft if needed
		updateCamera()

		// Render the scene
		renderer.Render(scene, cam)
	})

	log.Println("Example completed")
}

// initializeG3N initializes the G3N engine
func initializeG3N() {
	// Create the G3N application
	application = app.App()

	// Set the background color (dark blue for space)
	application.Gls().ClearColor(0.0, 0.0, 0.2, 1.0)

	// Create the scene
	scene = core.NewNode()

	// Create the camera
	cam = camera.New(1)
	cam.SetPosition(0, 50, 150)
	cam.LookAt(&math32.Vector3{0, 0, 0}, &math32.Vector3{0, 1, 0})
	scene.Add(cam)

	// Create the orbital camera control
	cameraCtrl = camera.NewOrbitControl(cam)

	// Add a handler for window resizing
	application.Subscribe(window.OnWindowSize, onWindowResize)

	// Set the initial aspect ratio of the camera
	width, height := application.GetSize()
	aspect := float32(width) / float32(height)
	cam.SetAspect(aspect)

	// Add lights
	// Softer ambient light for a space effect
	ambLight := light.NewAmbient(&math32.Color{0.3, 0.3, 0.4}, 0.5)
	scene.Add(ambLight)

	// More intense and distant point lights to illuminate the entire solar system
	pointLight1 := light.NewPoint(&math32.Color{1, 1, 1}, 5.0)
	pointLight1.SetPosition(50, 50, 50)
	pointLight1.SetLinearDecay(0.1)
	pointLight1.SetQuadraticDecay(0.01)
	scene.Add(pointLight1)

	pointLight2 := light.NewPoint(&math32.Color{1, 1, 1}, 5.0)
	pointLight2.SetPosition(-50, 50, 50)
	pointLight2.SetLinearDecay(0.1)
	pointLight2.SetQuadraticDecay(0.01)
	scene.Add(pointLight2)

	pointLight3 := light.NewPoint(&math32.Color{1, 1, 1}, 5.0)
	pointLight3.SetPosition(0, 50, -50)
	pointLight3.SetLinearDecay(0.1)
	pointLight3.SetQuadraticDecay(0.01)
	scene.Add(pointLight3)
}

// setupInputHandling configures the input handling for the application
func setupInputHandling() {
	// Subscribe to key events
	application.Subscribe(window.OnKeyDown, onKeyDown)
	application.Subscribe(window.OnKeyUp, onKeyUp)
}

// onKeyDown handles key press events
func onKeyDown(evname string, ev interface{}) {
	if spacecraftController == nil {
		return
	}

	kev := ev.(*window.KeyEvent)

	// Handle rotation controls
	switch kev.Key {
	case window.KeyW: // Pitch up
		spacecraftController.ApplyRotation(vector.NewVector3(1, 0, 0), 1.0)
	case window.KeyS: // Pitch down
		spacecraftController.ApplyRotation(vector.NewVector3(1, 0, 0), -1.0)
	case window.KeyA: // Yaw left
		spacecraftController.ApplyRotation(vector.NewVector3(0, 1, 0), 1.0)
	case window.KeyD: // Yaw right
		spacecraftController.ApplyRotation(vector.NewVector3(0, 1, 0), -1.0)
	case window.KeyQ: // Roll left
		spacecraftController.ApplyRotation(vector.NewVector3(0, 0, 1), 1.0)
	case window.KeyE: // Roll right
		spacecraftController.ApplyRotation(vector.NewVector3(0, 0, 1), -1.0)
	}

	// Handle thrust controls
	switch kev.Key {
	case window.KeySpace: // Toggle thrust
		if thrustLevel <= 0.0 {
			thrustLevel = 0.5 // 50% thrust
		} else {
			thrustLevel = 0.0 // No thrust
		}
	case window.KeyLeftShift: // Increase thrust
		thrustLevel += 0.1
		if thrustLevel > 1.0 {
			thrustLevel = 1.0
		}
	case window.KeyLeftControl: // Decrease thrust
		thrustLevel -= 0.1
		if thrustLevel < 0.0 {
			thrustLevel = 0.0
		}
	}

	// Camera controls
	switch kev.Key {
	case window.Key1: // First person view
		if spacecraftMesh != nil {
			// Position camera slightly above the spacecraft
			pos := spacecraftMesh.Position()
			cam.SetPosition(pos.X, pos.Y+2, pos.Z)
			// Look in the direction of spacecraft's forward vector
			// This is simplified and would need to be adjusted based on spacecraft orientation
			cam.LookAt(&math32.Vector3{pos.X, pos.Y, pos.Z - 10}, &math32.Vector3{0, 1, 0})
		}
	case window.Key2: // Third person view
		if spacecraftMesh != nil {
			// Position camera behind and above the spacecraft
			pos := spacecraftMesh.Position()
			cam.SetPosition(pos.X-10, pos.Y+5, pos.Z-10)
			cam.LookAt(&math32.Vector3{pos.X, pos.Y, pos.Z}, &math32.Vector3{0, 1, 0})
		}
	case window.Key3: // Free view (orbit control)
		// Reset to default view
		cam.SetPosition(0, 50, 150)
		cam.LookAt(&math32.Vector3{0, 0, 0}, &math32.Vector3{0, 1, 0})
	}
}

// onKeyUp handles key release events
func onKeyUp(evname string, ev interface{}) {
	// Currently not used, but could be used to implement continuous control
}

// onWindowResize handles window resizing
func onWindowResize(evname string, ev interface{}) {
	// Get the new window dimensions
	width, height := application.GetSize()

	// Update the camera's aspect ratio
	aspect := float32(width) / float32(height)
	cam.SetAspect(aspect)
}

// updateCamera updates the camera position to follow the spacecraft if needed
func updateCamera() {
	// This function can be expanded to implement camera following logic
	// For now, we'll rely on the key controls to change the camera view
}

// createMeshesForBodies creates meshes for all bodies in the world
func createMeshesForBodies(w world.World) {
	// Create meshes for all bodies
	for _, b := range w.GetBodies() {
		createMeshForBody(b)
	}
}

// createMeshForBody creates a mesh for a physical body
func createMeshForBody(b body.Body) {
	// Create a sphere to represent the body
	radius := float32(b.Radius().Value())

	// Increase the quality of spheres for larger bodies
	var segments, rings int
	if radius > 1.5 {
		segments, rings = 64, 32 // High quality for large planets
	} else if radius > 0.8 {
		segments, rings = 48, 24 // Medium quality for medium planets
	} else {
		segments, rings = 32, 16 // Standard quality for small bodies
	}

	geom := geometry.NewSphere(float64(radius), segments, rings)

	// Create a material based on the physical body's material
	var mat material.IMaterial
	var bodyColor math32.Color

	// Determine the color of the body
	if b.Material() != nil {
		// Map the physical material to a G3N color
		switch b.Material().Name() {
		case "Sun":
			// Special material for the sun with emission
			bodyColor = math32.Color{1.0, 0.8, 0.0}
			// Create a more intense emission color to simulate the brightness of the sun
			emissiveColor := math32.Color{1.0, 0.9, 0.5}
			sunMat := material.NewStandard(&bodyColor)
			sunMat.SetEmissiveColor(&emissiveColor)
			sunMat.SetOpacity(1.0)
			mat = sunMat
		case "Mercury":
			bodyColor = math32.Color{0.7, 0.7, 0.7}
		case "Venus":
			bodyColor = math32.Color{0.9, 0.7, 0.0}
		case "Earth":
			bodyColor = math32.Color{0.0, 0.3, 0.8}
		case "Mars":
			bodyColor = math32.Color{0.8, 0.3, 0.0}
		case "Jupiter":
			bodyColor = math32.Color{0.8, 0.6, 0.4}
		case "Saturn":
			bodyColor = math32.Color{0.9, 0.8, 0.5}
		case "Uranus":
			bodyColor = math32.Color{0.5, 0.8, 0.9}
		case "Neptune":
			bodyColor = math32.Color{0.0, 0.0, 0.8}
		case "Spacecraft":
			bodyColor = math32.Color{0.9, 0.9, 0.9} // White for spacecraft
		default:
			// Default gray color
			bodyColor = math32.Color{0.7, 0.7, 0.7}
		}

		// If the material has not already been created (as for the sun)
		if mat == nil {
			standardMat := material.NewStandard(&bodyColor)
			standardMat.SetShininess(30)
			mat = standardMat
		}
	} else {
		bodyColor = math32.Color{0.8, 0.8, 0.8}
		mat = material.NewStandard(&bodyColor)
	}

	// Create a mesh with the geometry and material
	mesh := graphic.NewMesh(geom, mat)

	// Set the position of the mesh
	pos := b.Position()
	mesh.SetPosition(float32(pos.X()), float32(pos.Y()), float32(pos.Z()))

	// Add the mesh to the scene
	scene.Add(mesh)

	// Store the mesh in the map
	bodyMeshes[b.ID().String()] = mesh

	// If this is the spacecraft, store a reference to its mesh
	if b.Material() != nil && b.Material().Name() == "Spacecraft" {
		spacecraftMesh = mesh
	}
}

// updateMeshPositions updates the positions of all meshes based on the positions of the bodies
func updateMeshPositions(w world.World) {
	for _, b := range w.GetBodies() {
		if mesh, exists := bodyMeshes[b.ID().String()]; exists {
			pos := b.Position()
			mesh.SetPosition(float32(pos.X()), float32(pos.Y()), float32(pos.Z()))

			// If this is the spacecraft, also update its rotation
			if b.Material() != nil && b.Material().Name() == "Spacecraft" && spacecraftEntity != nil {
				// Get rotation from the entity instead of the body
				rot := spacecraftEntity.GetRotation()
				// Convert from degrees to radians
				rotX := float32(rot.X() * math.Pi / 180.0)
				rotY := float32(rot.Y() * math.Pi / 180.0)
				rotZ := float32(rot.Z() * math.Pi / 180.0)

				// Apply rotation in the correct order (Y-X-Z)
				mesh.SetRotationY(rotY)
				mesh.SetRotationX(rotX)
				mesh.SetRotationZ(rotZ)
			}
		}
	}
}

// createBodies creates some bodies in the world
func createBodies(w world.World) {
	log.Println("Creating the solar system")
	createSolarSystem(w)
}

// createSolarSystem creates a realistic solar system
func createSolarSystem(w world.World) {
	log.Println("Creating the sun")

	// Fixed mass for the sun - high value to ensure stable orbits
	// In a simulation, relative ratios are more important than absolute values
	solarMass := 2e12 // Simplified value

	log.Printf("Sun mass: %e kg", solarMass)

	sun := body.NewRigidBody(
		units.NewQuantity(solarMass, units.Kilogram),
		units.NewQuantity(5.0, units.Meter),                        // Sun radius (scaled)
		vector.NewVector3(0, 0, 0),                                 // Position at center
		vector.NewVector3(0, 0, 0),                                 // Zero velocity
		createMaterial("Sun", 0.9, 0.5, [3]float64{1.0, 0.8, 0.0}), // Yellow color
	)
	sun.SetStatic(false) // The sun is not static (it can move)
	w.AddBody(sun)
	log.Printf("Sun created: ID=%v, Position=%v", sun.ID(), sun.Position())

	// Create the planets
	log.Println("Creating the planets")

	// Define planet distances
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

	// Create each planet
	for i := 0; i < len(names); i++ {
		// Create the planet - orbital velocity will be calculated inside the function
		planet := createPlanet(
			w,
			names[i],                   // Name
			solarMass*massFractions[i], // Mass
			radii[i],                   // Radius
			distances[i],               // Distance
			solarMass,                  // Mass of the central object (sun)
			vector.NewVector3(0, 1, 0), // Orbit plane
			colors[i],                  // Color
		)

		// Save Earth's position for spacecraft placement
		if names[i] == "Earth" {
			earthPosition = planet.Position()

			// Create a spacecraft near Earth
			createSpacecraft(w, planet)
		}
	}

	// Create an asteroid belt
	log.Println("Creating the asteroid belt")
	createAsteroidBelt(w, 200, solarMass, 60.0, 80.0) // Pass the solar mass as central mass
}

// createSpacecraft creates a controllable spacecraft near the given planet
func createSpacecraft(w world.World, planet body.Body) {
	log.Println("Creating controllable spacecraft near Earth")

	// Get planet position and velocity
	planetPos := planet.Position()
	planetVel := planet.Velocity()
	planetRadius := planet.Radius().Value()

	// Calculate spacecraft position (1.5 times the planet's radius away from the planet)
	offset := vector.NewVector3(0, planetRadius*1.5, 0)
	spacecraftPos := planetPos.Add(offset)

	// Set spacecraft initial velocity to match the planet's orbital velocity
	// This ensures it starts in a similar orbit
	spacecraftVel := planetVel

	// Create spacecraft configuration
	config := spacecraft.DefaultSpacecraftConfig()
	config.Position = spacecraftPos
	config.Velocity = spacecraftVel
	config.Mass = 100.0 // Much smaller than planets

	// Create the spacecraft and controller
	spacecraftEntity, spacecraftController = spacecraft.CreateSpacecraft(config)

	// Configure the spacecraft controller
	spacecraftController.SetMaxAngularVelocity(maxAngularVelocity)
	spacecraftController.SetAngularDamping(angularDamping)
	spacecraftController.SetRotationPIDGains(pidProportional, pidIntegral, pidDerivative)

	// Set a custom material for the spacecraft
	spacecraftEntity.GetBody().SetMaterial(createMaterial("Spacecraft", 0.8, 0.6, [3]float64{0.9, 0.9, 0.9}))

	// Add the spacecraft to the world
	w.AddBody(spacecraftEntity.GetBody())

	log.Printf("Spacecraft created: Position=%v, Velocity=%v",
		spacecraftEntity.GetPosition(), spacecraftEntity.GetVelocity())
}

// createPlanet creates a planet
func createPlanet(w world.World, name string, mass, radius, distance, centralMass float64, orbitPlane vector.Vector3, color [3]float64) body.Body {
	// Calculate the orbital velocity using the agnostic function from the force package
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
		-orbitSpeed*math.Sin(angle), // X component
		0,                           // Y component (xy plane)
		orbitSpeed*math.Cos(angle),  // Z component
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
