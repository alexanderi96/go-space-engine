// Package material provides implementations for physical materials
package material

import (
	"github.com/alexanderi96/go-space-engine/core/units"
)

// Material represents the physical properties of a material
type Material interface {
	// Name returns the name of the material
	Name() string

	// Density returns the density of the material
	Density() units.Quantity

	// SpecificHeat returns the specific heat capacity of the material
	SpecificHeat() units.Quantity

	// ThermalConductivity returns the thermal conductivity of the material
	ThermalConductivity() units.Quantity

	// Emissivity returns the emissivity of the material
	Emissivity() float64

	// Elasticity returns the elasticity of the material
	Elasticity() float64

	// Color returns the color of the material as RGB
	Color() [3]float64
}

// BasicMaterial implements a basic material
type BasicMaterial struct {
	name                string
	density             units.Quantity
	specificHeat        units.Quantity
	thermalConductivity units.Quantity
	emissivity          float64
	elasticity          float64
	color               [3]float64
}

// NewBasicMaterial creates a new basic material
func NewBasicMaterial(
	name string,
	density units.Quantity,
	specificHeat units.Quantity,
	thermalConductivity units.Quantity,
	emissivity float64,
	elasticity float64,
	color [3]float64,
) *BasicMaterial {
	return &BasicMaterial{
		name:                name,
		density:             density,
		specificHeat:        specificHeat,
		thermalConductivity: thermalConductivity,
		emissivity:          emissivity,
		elasticity:          elasticity,
		color:               color,
	}
}

// Name returns the name of the material
func (m *BasicMaterial) Name() string {
	return m.name
}

// Density returns the density of the material
func (m *BasicMaterial) Density() units.Quantity {
	return m.density
}

// SpecificHeat returns the specific heat capacity of the material
func (m *BasicMaterial) SpecificHeat() units.Quantity {
	return m.specificHeat
}

// ThermalConductivity returns the thermal conductivity of the material
func (m *BasicMaterial) ThermalConductivity() units.Quantity {
	return m.thermalConductivity
}

// Emissivity returns the emissivity of the material
func (m *BasicMaterial) Emissivity() float64 {
	return m.emissivity
}

// Elasticity returns the elasticity of the material
func (m *BasicMaterial) Elasticity() float64 {
	return m.elasticity
}

// Color returns the color of the material
func (m *BasicMaterial) Color() [3]float64 {
	return m.color
}

// Predefined materials
var (
	// Iron represents iron
	Iron = NewBasicMaterial(
		"Iron",
		units.NewQuantity(7874.0, units.Kilogram), // kg/m³ (density)
		units.NewQuantity(450.0, units.Joule),     // J/(kg·K) (specific heat capacity)
		units.NewQuantity(80.2, units.Watt),       // W/(m·K) (thermal conductivity)
		0.3,                                       // Emissivity
		0.7,                                       // Elasticity
		[3]float64{0.6, 0.6, 0.6},                 // Gray color
	)

	// Copper represents copper
	Copper = NewBasicMaterial(
		"Copper",
		units.NewQuantity(8960.0, units.Kilogram), // kg/m³
		units.NewQuantity(386.0, units.Joule),     // J/(kg·K)
		units.NewQuantity(401.0, units.Watt),      // W/(m·K)
		0.03,                                      // Emissivity
		0.75,                                      // Elasticity
		[3]float64{0.85, 0.45, 0.2},               // Copper color
	)

	// Ice represents ice
	Ice = NewBasicMaterial(
		"Ice",
		units.NewQuantity(917.0, units.Kilogram), // kg/m³
		units.NewQuantity(2108.0, units.Joule),   // J/(kg·K)
		units.NewQuantity(2.18, units.Watt),      // W/(m·K)
		0.97,                                     // Emissivity
		0.3,                                      // Elasticity
		[3]float64{0.8, 0.9, 0.95},               // Light blue color
	)

	// Water represents water
	Water = NewBasicMaterial(
		"Water",
		units.NewQuantity(997.0, units.Kilogram), // kg/m³
		units.NewQuantity(4186.0, units.Joule),   // J/(kg·K)
		units.NewQuantity(0.6, units.Watt),       // W/(m·K)
		0.95,                                     // Emissivity
		0.0,                                      // Elasticity (fluid)
		[3]float64{0.0, 0.3, 0.8},                // Blue color
	)

	// Rock represents rock
	Rock = NewBasicMaterial(
		"Rock",
		units.NewQuantity(2700.0, units.Kilogram), // kg/m³
		units.NewQuantity(840.0, units.Joule),     // J/(kg·K)
		units.NewQuantity(2.0, units.Watt),        // W/(m·K)
		0.8,                                       // Emissivity
		0.4,                                       // Elasticity
		[3]float64{0.5, 0.5, 0.5},                 // Gray color
	)
)

// Composition represents a composition of materials
type Composition struct {
	materials map[Material]float64 // Material -> Volume fraction (0-1)
}

// NewComposition creates a new composition
func NewComposition(fractions map[Material]float64) *Composition {
	// Normalize fractions to ensure they sum to 1
	total := 0.0
	for _, fraction := range fractions {
		total += fraction
	}

	normalized := make(map[Material]float64)
	for material, fraction := range fractions {
		normalized[material] = fraction / total
	}

	return &Composition{
		materials: normalized,
	}
}

// GetEffectiveProperties calculates the effective properties of the composition
func (c *Composition) GetEffectiveProperties() (Material, error) {
	// Calculate weighted average properties
	name := "Composite"
	densityValue := 0.0
	specificHeatValue := 0.0
	thermalConductivityValue := 0.0
	emissivityValue := 0.0
	elasticityValue := 0.0
	colorR, colorG, colorB := 0.0, 0.0, 0.0

	for material, fraction := range c.materials {
		densityValue += material.Density().Value() * fraction
		specificHeatValue += material.SpecificHeat().Value() * fraction
		thermalConductivityValue += material.ThermalConductivity().Value() * fraction
		emissivityValue += material.Emissivity() * fraction
		elasticityValue += material.Elasticity() * fraction

		color := material.Color()
		colorR += color[0] * fraction
		colorG += color[1] * fraction
		colorB += color[2] * fraction
	}

	// Create a new material with the calculated properties
	return NewBasicMaterial(
		name,
		units.NewQuantity(densityValue, units.Kilogram),
		units.NewQuantity(specificHeatValue, units.Joule),
		units.NewQuantity(thermalConductivityValue, units.Watt),
		emissivityValue,
		elasticityValue,
		[3]float64{colorR, colorG, colorB},
	), nil
}
