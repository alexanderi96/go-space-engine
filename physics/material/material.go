// Package material fornisce implementazioni per i materiali fisici
package material

import (
	"github.com/alexanderi96/go-space-engine/core/units"
)

// Material rappresenta le proprietà fisiche di un materiale
type Material interface {
	// Name restituisce il nome del materiale
	Name() string

	// Density restituisce la densità del materiale
	Density() units.Quantity

	// SpecificHeat restituisce la capacità termica specifica del materiale
	SpecificHeat() units.Quantity

	// ThermalConductivity restituisce la conducibilità termica del materiale
	ThermalConductivity() units.Quantity

	// Emissivity restituisce l'emissività del materiale
	Emissivity() float64

	// Elasticity restituisce l'elasticità del materiale
	Elasticity() float64

	// Color restituisce il colore del materiale come RGB
	Color() [3]float64
}

// BasicMaterial implementa un materiale base
type BasicMaterial struct {
	name                string
	density             units.Quantity
	specificHeat        units.Quantity
	thermalConductivity units.Quantity
	emissivity          float64
	elasticity          float64
	color               [3]float64
}

// NewBasicMaterial crea un nuovo materiale base
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

// Name restituisce il nome del materiale
func (m *BasicMaterial) Name() string {
	return m.name
}

// Density restituisce la densità del materiale
func (m *BasicMaterial) Density() units.Quantity {
	return m.density
}

// SpecificHeat restituisce la capacità termica specifica del materiale
func (m *BasicMaterial) SpecificHeat() units.Quantity {
	return m.specificHeat
}

// ThermalConductivity restituisce la conducibilità termica del materiale
func (m *BasicMaterial) ThermalConductivity() units.Quantity {
	return m.thermalConductivity
}

// Emissivity restituisce l'emissività del materiale
func (m *BasicMaterial) Emissivity() float64 {
	return m.emissivity
}

// Elasticity restituisce l'elasticità del materiale
func (m *BasicMaterial) Elasticity() float64 {
	return m.elasticity
}

// Color restituisce il colore del materiale
func (m *BasicMaterial) Color() [3]float64 {
	return m.color
}

// Materiali predefiniti
var (
	// Iron rappresenta il ferro
	Iron = NewBasicMaterial(
		"Iron",
		units.NewQuantity(7874.0, units.Kilogram), // kg/m³ (densità)
		units.NewQuantity(450.0, units.Joule),     // J/(kg·K) (capacità termica specifica)
		units.NewQuantity(80.2, units.Watt),       // W/(m·K) (conducibilità termica)
		0.3,                                       // Emissività
		0.7,                                       // Elasticità
		[3]float64{0.6, 0.6, 0.6},                 // Colore grigio
	)

	// Copper rappresenta il rame
	Copper = NewBasicMaterial(
		"Copper",
		units.NewQuantity(8960.0, units.Kilogram), // kg/m³
		units.NewQuantity(386.0, units.Joule),     // J/(kg·K)
		units.NewQuantity(401.0, units.Watt),      // W/(m·K)
		0.03,                                      // Emissività
		0.75,                                      // Elasticità
		[3]float64{0.85, 0.45, 0.2},               // Colore rame
	)

	// Ice rappresenta il ghiaccio
	Ice = NewBasicMaterial(
		"Ice",
		units.NewQuantity(917.0, units.Kilogram), // kg/m³
		units.NewQuantity(2108.0, units.Joule),   // J/(kg·K)
		units.NewQuantity(2.18, units.Watt),      // W/(m·K)
		0.97,                                     // Emissività
		0.3,                                      // Elasticità
		[3]float64{0.8, 0.9, 0.95},               // Colore azzurro chiaro
	)

	// Water rappresenta l'acqua
	Water = NewBasicMaterial(
		"Water",
		units.NewQuantity(997.0, units.Kilogram), // kg/m³
		units.NewQuantity(4186.0, units.Joule),   // J/(kg·K)
		units.NewQuantity(0.6, units.Watt),       // W/(m·K)
		0.95,                                     // Emissività
		0.0,                                      // Elasticità (fluido)
		[3]float64{0.0, 0.3, 0.8},                // Colore blu
	)

	// Rock rappresenta la roccia
	Rock = NewBasicMaterial(
		"Rock",
		units.NewQuantity(2700.0, units.Kilogram), // kg/m³
		units.NewQuantity(840.0, units.Joule),     // J/(kg·K)
		units.NewQuantity(2.0, units.Watt),        // W/(m·K)
		0.8,                                       // Emissività
		0.4,                                       // Elasticità
		[3]float64{0.5, 0.5, 0.5},                 // Colore grigio
	)
)

// Composition rappresenta una composizione di materiali
type Composition struct {
	materials map[Material]float64 // Materiale -> Frazione volumetrica (0-1)
}

// NewComposition crea una nuova composizione
func NewComposition(fractions map[Material]float64) *Composition {
	// Normalizza le frazioni per assicurarsi che sommino a 1
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

// GetEffectiveProperties calcola le proprietà effettive della composizione
func (c *Composition) GetEffectiveProperties() (Material, error) {
	// Calcola le proprietà medie pesate
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

	// Crea un nuovo materiale con le proprietà calcolate
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
