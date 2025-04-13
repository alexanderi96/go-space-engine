package spacecraft

import (
	"github.com/alexanderi96/go-space-engine/core/units"
	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/entity"
	"github.com/alexanderi96/go-space-engine/physics/body"
	"github.com/alexanderi96/go-space-engine/physics/material"
	"github.com/google/uuid"
)

// SpacecraftConfig contains configuration parameters for creating a spacecraft
type SpacecraftConfig struct {
	// Physical properties
	Mass   float64 // Mass in kilograms
	Radius float64 // Radius in meters

	// Performance characteristics
	MaxThrust float64 // Maximum thrust force in Newtons
	MaxTorque float64 // Maximum torque in Newton-meters

	// Initial conditions
	Position vector.Vector3 // Initial position
	Velocity vector.Vector3 // Initial velocity
	Rotation vector.Vector3 // Initial rotation

	// Visual properties
	IsCube bool       // Whether the spacecraft should be rendered as a cube
	Color  [3]float64 // Color of the spacecraft
}

// DefaultSpacecraftConfig returns a configuration with sensible defaults
func DefaultSpacecraftConfig() SpacecraftConfig {
	return SpacecraftConfig{
		Mass:      1000.0,  // 1 tonne
		Radius:    5.0,     // 5 meters
		MaxThrust: 10000.0, // 10 kN
		MaxTorque: 1000.0,  // 1 kNm
		Position:  vector.Zero3(),
		Velocity:  vector.Zero3(),
		Rotation:  vector.Zero3(),
		IsCube:    false,
		Color:     [3]float64{0.8, 0.8, 0.9}, // Light metallic blue color
	}
}

// CreateSpacecraft creates a new spacecraft entity with the specified configuration
func CreateSpacecraft(config SpacecraftConfig) (entity.Entity, *SpacecraftController) {
	// Create a default material for the spacecraft
	color := config.Color
	if color == [3]float64{0, 0, 0} {
		color = [3]float64{0.8, 0.8, 0.9} // Default light metallic blue color
	}

	mat := material.NewBasicMaterial(
		"Spacecraft", // Cambiato da "spacecraft" a "Spacecraft" per facilitare l'identificazione
		units.NewQuantity(2700.0, units.Kilogram), // Density of aluminum
		units.NewQuantity(900.0, units.Joule),     // Specific heat of aluminum
		units.NewQuantity(237.0, units.Watt),      // Thermal conductivity of aluminum
		0.7,                                       // Emissivity
		0.8,                                       // Elasticity (slightly bouncy)
		color,                                     // Use the color from config
	)

	// Create the physical body
	physBody := body.NewRigidBody(
		units.NewQuantity(config.Mass, units.Kilogram),
		units.NewQuantity(config.Radius, units.Meter),
		config.Position,
		config.Velocity,
		mat,
	)

	// Set initial rotation on the physical body
	physBody.SetRotation(config.Rotation)

	// Create the entity
	id := uuid.New().String()
	spacecraft := entity.NewBaseEntity(id, physBody)

	// Create the controller
	controller := NewSpacecraftController(spacecraft, config.MaxThrust, config.MaxTorque)

	return spacecraft, controller
}
