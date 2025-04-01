// Package config fornisce la configurazione per la simulazione
package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/alexanderi96/go-space-engine/core/constants"
	"github.com/alexanderi96/go-space-engine/core/units"
	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/physics/space"
)

// Config rappresenta la configurazione della simulazione
type Config struct {
	// Configurazione generale
	TimeStep           float64 `json:"timeStep"`           // Passo temporale della simulazione (s)
	MaxBodies          int     `json:"maxBodies"`          // Numero massimo di corpi
	GravityEnabled     bool    `json:"gravityEnabled"`     // Indica se la gravità è abilitata
	GravityConstant    float64 `json:"gravityConstant"`    // Costante gravitazionale (m³/kg⋅s²)
	CollisionsEnabled  bool    `json:"collisionsEnabled"`  // Indica se le collisioni sono abilitate
	BoundaryCollisions bool    `json:"boundaryCollisions"` // Indica se le collisioni con i limiti sono abilitate

	// Configurazione dell'octree
	OctreeMaxObjects int `json:"octreeMaxObjects"` // Numero massimo di oggetti per nodo dell'octree
	OctreeMaxLevels  int `json:"octreeMaxLevels"`  // Numero massimo di livelli dell'octree

	// Configurazione dei limiti del mondo
	WorldMin vector.Vector3 `json:"worldMin"` // Punto minimo dei limiti del mondo
	WorldMax vector.Vector3 `json:"worldMax"` // Punto massimo dei limiti del mondo

	// Configurazione della fisica
	Restitution float64 `json:"restitution"` // Coefficiente di restituzione (elasticità)

	// Configurazione dell'integratore
	IntegratorType string `json:"integratorType"` // Tipo di integratore ("euler", "verlet", "rk4")
}

// NewDefaultConfig crea una nuova configurazione con valori predefiniti
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

		IntegratorType: "verlet",
	}
}

// GetWorldBounds restituisce i limiti del mondo come AABB
func (c *Config) GetWorldBounds() *space.AABB {
	return space.NewAABB(c.WorldMin, c.WorldMax)
}

// GetTimeStepQuantity restituisce il passo temporale come Quantity
func (c *Config) GetTimeStepQuantity() units.Quantity {
	return units.NewQuantity(c.TimeStep, units.Second)
}

// GetGravityConstantQuantity restituisce la costante gravitazionale come Quantity
func (c *Config) GetGravityConstantQuantity() units.Quantity {
	// G ha unità di misura m³/(kg⋅s²)
	return units.NewQuantity(c.GravityConstant, units.Newton)
}

// SaveToFile salva la configurazione su file
func (c *Config) SaveToFile(filename string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0644)
}

// LoadFromFile carica la configurazione da file
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

// SimulationBuilder è un builder per la simulazione
type SimulationBuilder struct {
	config *Config
}

// NewSimulationBuilder crea un nuovo builder per la simulazione
func NewSimulationBuilder() *SimulationBuilder {
	return &SimulationBuilder{
		config: NewDefaultConfig(),
	}
}

// WithTimeStep imposta il passo temporale
func (b *SimulationBuilder) WithTimeStep(timeStep float64) *SimulationBuilder {
	b.config.TimeStep = timeStep
	return b
}

// WithMaxBodies imposta il numero massimo di corpi
func (b *SimulationBuilder) WithMaxBodies(maxBodies int) *SimulationBuilder {
	b.config.MaxBodies = maxBodies
	return b
}

// WithGravity imposta se la gravità è abilitata
func (b *SimulationBuilder) WithGravity(enabled bool) *SimulationBuilder {
	b.config.GravityEnabled = enabled
	return b
}

// WithGravityConstant imposta la costante gravitazionale
func (b *SimulationBuilder) WithGravityConstant(g float64) *SimulationBuilder {
	b.config.GravityConstant = g
	return b
}

// WithCollisions imposta se le collisioni sono abilitate
func (b *SimulationBuilder) WithCollisions(enabled bool) *SimulationBuilder {
	b.config.CollisionsEnabled = enabled
	return b
}

// WithBoundaryCollisions imposta se le collisioni con i limiti sono abilitate
func (b *SimulationBuilder) WithBoundaryCollisions(enabled bool) *SimulationBuilder {
	b.config.BoundaryCollisions = enabled
	return b
}

// WithOctreeConfig imposta la configurazione dell'octree
func (b *SimulationBuilder) WithOctreeConfig(maxObjects, maxLevels int) *SimulationBuilder {
	b.config.OctreeMaxObjects = maxObjects
	b.config.OctreeMaxLevels = maxLevels
	return b
}

// WithWorldBounds imposta i limiti del mondo
func (b *SimulationBuilder) WithWorldBounds(min, max vector.Vector3) *SimulationBuilder {
	b.config.WorldMin = min
	b.config.WorldMax = max
	return b
}

// WithRestitution imposta il coefficiente di restituzione
func (b *SimulationBuilder) WithRestitution(restitution float64) *SimulationBuilder {
	b.config.Restitution = restitution
	return b
}

// WithIntegratorType imposta il tipo di integratore
func (b *SimulationBuilder) WithIntegratorType(integratorType string) *SimulationBuilder {
	b.config.IntegratorType = integratorType
	return b
}

// Build restituisce la configurazione
func (b *SimulationBuilder) Build() *Config {
	return b.config
}
