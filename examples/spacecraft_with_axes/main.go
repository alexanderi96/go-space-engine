// Package main provides an example of spacecraft control with reference axes
package main

import (
	"log"
	"time"

	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/entity/input"
	"github.com/alexanderi96/go-space-engine/entity/vehicle/spacecraft"
	"github.com/alexanderi96/go-space-engine/render/g3n"
	"github.com/alexanderi96/go-space-engine/simulation/config"
	"github.com/alexanderi96/go-space-engine/simulation/world"

	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
)

const (
	// Simulation parameters
	simulationTimeStep = 1.0 / 60.0 // 60 Hz simulation

	// Spacecraft controller parameters
	maxAngularVelocity = 60.0 // 60 degrees per second
	angularDamping     = 0.3  // 30% damping
	pidProportional    = 3.0  // P gain
	pidIntegral        = 0.1  // I gain
	pidDerivative      = 0.5  // D gain

	// Reference axes parameters
	axisLength = 1.0  // Length of the reference axes
	axisWidth  = 0.05 // Width of the reference axes
)

// InputHandler handles keyboard input for spacecraft control
type InputHandler struct {
	controller *spacecraft.SpacecraftController
	thrustKey  bool
	yawLeft    bool
	yawRight   bool
	pitchUp    bool
	pitchDown  bool
	rollLeft   bool
	rollRight  bool
}

// NewInputHandler creates a new input handler
func NewInputHandler(controller *spacecraft.SpacecraftController) *InputHandler {
	return &InputHandler{
		controller: controller,
	}
}

// HandleEvent handles input events
func (h *InputHandler) HandleEvent(event input.Event) bool {
	// Check if this is a key event
	if keyEvent, ok := event.(*input.KeyEvent); ok {
		// Handle key press/release
		if keyEvent.Type == input.EventKeyDown {
			switch keyEvent.Key {
			case int(window.KeyW): // Forward thrust
				h.thrustKey = true
			case int(window.KeyA): // Yaw left
				h.yawLeft = true
			case int(window.KeyD): // Yaw right
				h.yawRight = true
			case int(window.KeyS): // Pitch down
				h.pitchDown = true
			case int(window.KeyUp): // Pitch up
				h.pitchUp = true
			case int(window.KeyLeft): // Roll left
				h.rollLeft = true
			case int(window.KeyRight): // Roll right
				h.rollRight = true
			}
		} else if keyEvent.Type == input.EventKeyUp {
			switch keyEvent.Key {
			case int(window.KeyW): // Forward thrust
				h.thrustKey = false
			case int(window.KeyA): // Yaw left
				h.yawLeft = false
			case int(window.KeyD): // Yaw right
				h.yawRight = false
			case int(window.KeyS): // Pitch down
				h.pitchDown = false
			case int(window.KeyUp): // Pitch up
				h.pitchUp = false
			case int(window.KeyLeft): // Roll left
				h.rollLeft = false
			case int(window.KeyRight): // Roll right
				h.rollRight = false
			}
		}
		return true
	}
	return false
}

// Update updates the spacecraft controller based on input
func (h *InputHandler) Update() {
	// Apply thrust if the thrust key is pressed
	if h.thrustKey {
		h.controller.SetThrustLevel(1.0)
	} else {
		h.controller.SetThrustLevel(0.0)
	}

	// Calculate rotation based on key presses
	yawAmount := 0.0
	if h.yawLeft {
		yawAmount -= 1.0
	}
	if h.yawRight {
		yawAmount += 1.0
	}

	pitchAmount := 0.0
	if h.pitchUp {
		pitchAmount -= 1.0
	}
	if h.pitchDown {
		pitchAmount += 1.0
	}

	rollAmount := 0.0
	if h.rollLeft {
		rollAmount -= 1.0
	}
	if h.rollRight {
		rollAmount += 1.0
	}

	// Get current rotation
	currentRotation := h.controller.GetEntity().GetRotation()

	// Calculate target rotation based on input
	targetYaw := currentRotation.Y()
	targetPitch := currentRotation.X()
	targetRoll := currentRotation.Z()

	// Update target rotation based on input
	// The rotation speed is controlled by the controller's PID parameters
	if yawAmount != 0 {
		targetYaw += yawAmount * 5.0 // 5 degrees per update
	}
	if pitchAmount != 0 {
		targetPitch += pitchAmount * 5.0
	}
	if rollAmount != 0 {
		targetRoll += rollAmount * 5.0
	}

	// Set the target rotation
	h.controller.SetTargetRotation(vector.NewVector3(targetPitch, targetYaw, targetRoll))
}

// createReferenceAxes creates reference axes at the specified position
func createReferenceAxes(adapter *g3n.G3NAdapter, position vector.Vector3) {
	// Create X axis (red)
	xGeom := geometry.NewCylinder(float64(axisWidth), float64(axisWidth), float64(axisLength), 8, 1, false)
	xMat := material.NewStandard(&math32.Color{1, 0, 0})
	xMesh := graphic.NewMesh(xGeom, xMat)
	xMesh.SetPosition(float32(position.X()), float32(position.Y()), float32(position.Z()))
	xMesh.SetRotationZ(math32.Pi / 2) // Rotate to align with X axis
	adapter.GetScene().Add(xMesh)

	// Create Y axis (green)
	yGeom := geometry.NewCylinder(float64(axisWidth), float64(axisWidth), float64(axisLength), 8, 1, false)
	yMat := material.NewStandard(&math32.Color{0, 1, 0})
	yMesh := graphic.NewMesh(yGeom, yMat)
	yMesh.SetPosition(float32(position.X()), float32(position.Y()), float32(position.Z()))
	// Y axis is already aligned correctly
	adapter.GetScene().Add(yMesh)

	// Create Z axis (blue)
	zGeom := geometry.NewCylinder(float64(axisWidth), float64(axisWidth), float64(axisLength), 8, 1, false)
	zMat := material.NewStandard(&math32.Color{0, 0, 1})
	zMesh := graphic.NewMesh(zGeom, zMat)
	zMesh.SetPosition(float32(position.X()), float32(position.Y()), float32(position.Z()))
	zMesh.SetRotationX(math32.Pi / 2) // Rotate to align with Z axis
	adapter.GetScene().Add(zMesh)

	// Add small spheres at the end of each axis for better visibility
	sphereRadius := float64(axisWidth * 2)

	// X axis endpoint (red)
	xSphereGeom := geometry.NewSphere(sphereRadius, 16, 16)
	xSphereMat := material.NewStandard(&math32.Color{1, 0, 0})
	xSphereMesh := graphic.NewMesh(xSphereGeom, xSphereMat)
	xSphereMesh.SetPosition(float32(position.X()+axisLength), float32(position.Y()), float32(position.Z()))
	adapter.GetScene().Add(xSphereMesh)

	// Y axis endpoint (green)
	ySphereGeom := geometry.NewSphere(sphereRadius, 16, 16)
	ySphereMat := material.NewStandard(&math32.Color{0, 1, 0})
	ySphereMesh := graphic.NewMesh(ySphereGeom, ySphereMat)
	ySphereMesh.SetPosition(float32(position.X()), float32(position.Y()+axisLength), float32(position.Z()))
	adapter.GetScene().Add(ySphereMesh)

	// Z axis endpoint (blue)
	zSphereGeom := geometry.NewSphere(sphereRadius, 16, 16)
	zSphereMat := material.NewStandard(&math32.Color{0, 0, 1})
	zSphereMesh := graphic.NewMesh(zSphereGeom, zSphereMat)
	zSphereMesh.SetPosition(float32(position.X()), float32(position.Y()), float32(position.Z()+axisLength))
	adapter.GetScene().Add(zSphereMesh)
}

func main() {
	log.Println("Starting spacecraft control simulation with reference axes...")

	// Create the simulation configuration
	cfg := config.NewSimulationBuilder().
		WithTimeStep(simulationTimeStep).
		WithMaxBodies(10).
		WithGravity(false).    // No gravity for this example
		WithCollisions(false). // No collisions for this example
		WithBoundaryCollisions(true).
		WithWorldBounds(
			vector.NewVector3(-500, -500, -500), // Min point
			vector.NewVector3(500, 500, 500),    // Max point
		).
		Build()

	// Create the simulation world
	w := world.NewPhysicalWorld(cfg.GetWorldBounds())

	// Create spacecraft and controller
	spacecraftConfig := spacecraft.DefaultSpacecraftConfig()
	spacecraftConfig.Position = vector.NewVector3(0, 0, 0)
	spacecraftConfig.Velocity = vector.NewVector3(0, 0, 0)
	spacecraftConfig.IsCube = true // Render as a cube for better orientation visibility

	entity, controller := spacecraft.CreateSpacecraft(spacecraftConfig)

	// Configure the spacecraft controller
	controller.SetMaxAngularVelocity(maxAngularVelocity)
	controller.SetAngularDamping(angularDamping)
	controller.SetRotationPIDGains(pidProportional, pidIntegral, pidDerivative)

	// Set initial target rotation to current rotation
	controller.SetTargetRotation(entity.GetRotation())

	// Add the entity to the world
	w.AddBody(entity.GetBody())

	// Create the G3N adapter
	adapter := g3n.NewG3NAdapter()

	// Configure the adapter
	adapter.SetBackgroundColor(g3n.NewColor(0.05, 0.05, 0.1, 1.0)) // Dark blue background for space

	// Create input handler
	inputHandler := NewInputHandler(controller)
	adapter.RegisterInputHandler(inputHandler)

	// Create reference axes at the origin
	createReferenceAxes(adapter, vector.NewVector3(0, 0, 0))

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

		// Update input handler
		inputHandler.Update()

		// Update the controller and world
		controller.Update(dt)
		w.Step(dt)

		// Render the world
		adapter.RenderWorld(w)
	})

	log.Println("Simulation completed!")
}
