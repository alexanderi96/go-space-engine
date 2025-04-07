// Package main provides an example of spacecraft control using the physics engine
package main

import (
	"log"
	"time"

	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/entity"
	"github.com/alexanderi96/go-space-engine/entity/vehicle/spacecraft"
	"github.com/alexanderi96/go-space-engine/simulation/config"
	"github.com/alexanderi96/go-space-engine/simulation/world"
)

const (
	// Simulation parameters
	simulationTimeStep = 1.0 / 60.0 // 60 Hz simulation
	simulationDuration = 10.0       // 10 seconds total

	// Spacecraft controller parameters
	maxAngularVelocity = 60.0 // 60 degrees per second
	angularDamping     = 0.3  // 30% damping
	pidProportional    = 3.0  // P gain
	pidIntegral        = 0.1  // I gain
	pidDerivative      = 0.5  // D gain

	// Mission parameters
	firstRotationTime  = 2.0 // Time to complete first rotation
	firstThrustTime    = 4.0 // Time to complete first thrust phase
	secondRotationTime = 6.0 // Time to complete second rotation
)

func main() {
	log.Println("Starting spacecraft control simulation...")

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
	entity, controller := createSpacecraft()

	// Add the entity to the world
	w.AddBody(entity.GetBody())

	log.Printf("Initial position: %v\n", entity.GetPosition())
	log.Printf("Initial rotation: %v\n", entity.GetRotation())

	// Main simulation loop
	elapsed := 0.0
	lastUpdateTime := time.Now()
	lastPrintTime := -1 // Initialize to -1 to ensure first print at 0 seconds

	for elapsed < simulationDuration {
		// Calculate real delta time
		currentTime := time.Now()
		dt := currentTime.Sub(lastUpdateTime).Seconds()
		lastUpdateTime = currentTime

		// Limit delta time to avoid instability
		if dt > 0.1 {
			dt = 0.1
		}

		// Apply mission control logic
		updateMissionControl(controller, elapsed)

		// Update the controller and world
		controller.Update(dt)
		w.Step(dt)

		// Print state exactly once per second
		currentSecond := int(elapsed)
		if currentSecond != lastPrintTime {
			printSpacecraftState(elapsed, entity)
			lastPrintTime = currentSecond
		}

		elapsed += dt

		// Sleep to control simulation speed (run at 10x speed for testing)
		time.Sleep(time.Duration(dt * float64(time.Second) / 10))
	}

	log.Println("Final position:", entity.GetPosition())
	log.Println("Simulation completed!")
}

// createSpacecraft creates and configures a spacecraft with its controller
func createSpacecraft() (entity.Entity, *spacecraft.SpacecraftController) {
	// Create a spacecraft with default configuration
	config := spacecraft.DefaultSpacecraftConfig()
	config.Position = vector.NewVector3(0, 0, 0)
	config.Velocity = vector.NewVector3(0, 0, 0)

	// Create the spacecraft and controller
	entity, controller := spacecraft.CreateSpacecraft(config)

	// Configure the spacecraft controller
	controller.SetMaxAngularVelocity(maxAngularVelocity)
	controller.SetAngularDamping(angularDamping)
	controller.SetRotationPIDGains(pidProportional, pidIntegral, pidDerivative)

	// Set initial target rotation to current rotation
	controller.SetTargetRotation(entity.GetRotation())

	return entity, controller
}

// updateMissionControl updates the spacecraft controller based on mission phase
func updateMissionControl(controller *spacecraft.SpacecraftController, elapsed float64) {
	if elapsed < firstRotationTime {
		// First phase: rotate to point in a direction (90 degrees)
		targetRotation := vector.NewVector3(0, 90, 0)
		controller.SetTargetRotation(targetRotation)
		controller.SetThrustLevel(0.0) // No thrust while rotating
	} else if elapsed < firstThrustTime {
		// Second phase: apply thrust to move forward
		controller.SetThrustLevel(0.5) // 50% thrust
	} else if elapsed < secondRotationTime {
		// Third phase: rotate in the opposite direction (270 degrees)
		targetRotation := vector.NewVector3(0, 270, 0)
		controller.SetTargetRotation(targetRotation)
		controller.SetThrustLevel(0.0) // No thrust while rotating
	} else {
		// Final phase: apply thrust again
		controller.SetThrustLevel(0.7) // 70% thrust
	}
}

// printSpacecraftState prints the current state of the spacecraft
func printSpacecraftState(elapsed float64, entity entity.Entity) {
	log.Printf("Time: %.1fs\n", elapsed)
	log.Printf("Position: %v\n", entity.GetPosition())
	log.Printf("Rotation: %v\n", entity.GetRotation())
	log.Printf("Velocity: %v\n", entity.GetVelocity())
	log.Printf("Angular Velocity: %v\n", entity.GetAngularVelocity())
	log.Println("---")
}
