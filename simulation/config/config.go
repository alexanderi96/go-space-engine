// Package config provides configuration for the simulation
package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/alexanderi96/go-space-engine/core/constants"
	"github.com/alexanderi96/go-space-engine/core/units"
	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/physics/space"
)

// Config represents the simulation configuration
type Config struct {
	// General configuration
	TimeStep           float64 `json:"timeStep"`           // Simulation time step (s)
	MaxBodies          int     `json:"maxBodies"`          // Maximum number of bodies
	GravityEnabled     bool    `json:"gravityEnabled"`     // Indicates if gravity is enabled
	GravityConstant    float64 `json:"gravityConstant"`    // Gravitational constant (m³/kg⋅s²)
	CollisionsEnabled  bool    `json:"collisionsEnabled"`  // Indicates if collisions are enabled
	BoundaryCollisions bool    `json:"boundaryCollisions"` // Indicates if boundary collisions are enabled

	// Octree configuration
	OctreeMaxObjects int `json:"octreeMaxObjects"` // Maximum number of objects per octree node
	OctreeMaxLevels  int `json:"octreeMaxLevels"`  // Maximum number of octree levels

	// World boundaries configuration
	WorldMin vector.Vector3 `json:"worldMin"` // Minimum point of world boundaries
	WorldMax vector.Vector3 `json:"worldMax"` // Maximum point of world boundaries

	// Physics configuration
	Restitution float64 `json:"restitution"` // Coefficient of restitution (elasticity)

	// Integrator configuration
	IntegratorType string `json:"integratorType"` // Integrator type ("euler", "verlet", "rk4")
}

// NewDefaultConfig creates a new configuration with default values
func NewDefaultConfig() *Config {
	return &Config{
		TimeStep:           0.01,
		MaxBodies:          1000,
		GravityEnabled:     true,
		GravityConstant:    constants.G,
		CollisionsEnabled:  true,
		BoundaryCollisions: true,

		OctreeMaxObjects: 10,
		OctreeMaxLevels:  8,

		WorldMin: vector.NewVector3(-100, -100, -100),
		WorldMax: vector.NewVector3(100, 100, 100),

		Restitution: 0.5,

		IntegratorType: "rk4",
	}
}

// GetWorldBounds returns the world boundaries as AABB
func (c *Config) GetWorldBounds() *space.AABB {
	return space.NewAABB(c.WorldMin, c.WorldMax)
}

// GetTimeStepQuantity returns the time step as a Quantity
func (c *Config) GetTimeStepQuantity() units.Quantity {
	return units.NewQuantity(c.TimeStep, units.Second)
}

// GetGravityConstantQuantity returns the gravitational constant as a Quantity
func (c *Config) GetGravityConstantQuantity() units.Quantity {
	// G has units of m³/(kg⋅s²)
	return units.NewQuantity(c.GravityConstant, units.Newton)
}

// SaveToFile saves the configuration to a file
func (c *Config) SaveToFile(filename string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0644)
}

// LoadFromFile loads the configuration from a file
func LoadFromFile(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// SimulationBuilder is a builder for the simulation
type SimulationBuilder struct {
	config *Config
}

// NewSimulationBuilder creates a new builder for the simulation
func NewSimulationBuilder() *SimulationBuilder {
	return &SimulationBuilder{
		config: NewDefaultConfig(),
	}
}

// WithTimeStep sets the time step
func (b *SimulationBuilder) WithTimeStep(timeStep float64) *SimulationBuilder {
	b.config.TimeStep = timeStep
	return b
}

// WithMaxBodies sets the maximum number of bodies
func (b *SimulationBuilder) WithMaxBodies(maxBodies int) *SimulationBuilder {
	b.config.MaxBodies = maxBodies
	return b
}

// WithGravity sets whether gravity is enabled
func (b *SimulationBuilder) WithGravity(enabled bool) *SimulationBuilder {
	b.config.GravityEnabled = enabled
	return b
}

// WithGravityConstant sets the gravitational constant
func (b *SimulationBuilder) WithGravityConstant(g float64) *SimulationBuilder {
	b.config.GravityConstant = g
	return b
}

// WithCollisions sets whether collisions are enabled
func (b *SimulationBuilder) WithCollisions(enabled bool) *SimulationBuilder {
	b.config.CollisionsEnabled = enabled
	return b
}

// WithBoundaryCollisions sets whether boundary collisions are enabled
func (b *SimulationBuilder) WithBoundaryCollisions(enabled bool) *SimulationBuilder {
	b.config.BoundaryCollisions = enabled
	return b
}

// WithOctreeConfig sets the octree configuration
func (b *SimulationBuilder) WithOctreeConfig(maxObjects, maxLevels int) *SimulationBuilder {
	b.config.OctreeMaxObjects = maxObjects
	b.config.OctreeMaxLevels = maxLevels
	return b
}

// WithWorldBounds sets the world boundaries
func (b *SimulationBuilder) WithWorldBounds(min, max vector.Vector3) *SimulationBuilder {
	b.config.WorldMin = min
	b.config.WorldMax = max
	return b
}

// WithRestitution sets the coefficient of restitution
func (b *SimulationBuilder) WithRestitution(restitution float64) *SimulationBuilder {
	b.config.Restitution = restitution
	return b
}

// WithIntegratorType sets the integrator type
func (b *SimulationBuilder) WithIntegratorType(integratorType string) *SimulationBuilder {
	b.config.IntegratorType = integratorType
	return b
}

// Build returns the configuration
func (b *SimulationBuilder) Build() *Config {
	return b.config
}
