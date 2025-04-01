// Package units fornisce un sistema di unità di misura per il motore fisico
package units

import (
	"fmt"
	"math"
)

// UnitType rappresenta il tipo di unità di misura
type UnitType int

const (
	// Length rappresenta un'unità di lunghezza
	Length UnitType = iota
	// Mass rappresenta un'unità di massa
	Mass
	// Time rappresenta un'unità di tempo
	Time
	// Temperature rappresenta un'unità di temperatura
	Temperature
	// Angle rappresenta un'unità di angolo
	Angle
	// Force rappresenta un'unità di forza
	Force
	// Energy rappresenta un'unità di energia
	Energy
	// Power rappresenta un'unità di potenza
	Power
	// Velocity rappresenta un'unità di velocità
	Velocity
	// Acceleration rappresenta un'unità di accelerazione
	Acceleration
	// Pressure rappresenta un'unità di pressione
	Pressure
)

// Unit rappresenta un'unità di misura
type Unit interface {
	// Type restituisce il tipo di unità
	Type() UnitType
	// Name restituisce il nome dell'unità
	Name() string
	// Symbol restituisce il simbolo dell'unità
	Symbol() string
	// ConvertTo converte un valore da questa unità a un'altra
	ConvertTo(value float64, target Unit) float64
	// ConvertFrom converte un valore da un'altra unità a questa
	ConvertFrom(value float64, source Unit) float64
}

// BaseUnit implementa un'unità di misura di base
type BaseUnit struct {
	unitType UnitType
	name     string
	symbol   string
	factor   float64 // Fattore di conversione rispetto all'unità SI
	offset   float64 // Offset per unità con punto zero diverso (es. temperature)
}

// NewBaseUnit crea una nuova unità di base
func NewBaseUnit(unitType UnitType, name, symbol string, factor, offset float64) *BaseUnit {
	return &BaseUnit{
		unitType: unitType,
		name:     name,
		symbol:   symbol,
		factor:   factor,
		offset:   offset,
	}
}

// Type restituisce il tipo di unità
func (u *BaseUnit) Type() UnitType {
	return u.unitType
}

// Name restituisce il nome dell'unità
func (u *BaseUnit) Name() string {
	return u.name
}

// Symbol restituisce il simbolo dell'unità
func (u *BaseUnit) Symbol() string {
	return u.symbol
}

// ConvertTo converte un valore da questa unità a un'altra
func (u *BaseUnit) ConvertTo(value float64, target Unit) float64 {
	if u.Type() != target.Type() {
		panic(fmt.Sprintf("Cannot convert between different unit types: %v and %v", u.Type(), target.Type()))
	}

	// Converti prima in unità SI
	siValue := (value + u.offset) * u.factor

	// Poi converti da SI all'unità target
	targetUnit, ok := target.(*BaseUnit)
	if !ok {
		panic("Target unit is not a BaseUnit")
	}

	return (siValue / targetUnit.factor) - targetUnit.offset
}

// ConvertFrom converte un valore da un'altra unità a questa
func (u *BaseUnit) ConvertFrom(value float64, source Unit) float64 {
	return source.(*BaseUnit).ConvertTo(value, u)
}

// DerivedUnit implementa un'unità di misura derivata
type DerivedUnit struct {
	BaseUnit
	components map[Unit]int // Mappa di unità base e loro esponenti
}

// NewDerivedUnit crea una nuova unità derivata
func NewDerivedUnit(unitType UnitType, name, symbol string, components map[Unit]int) *DerivedUnit {
	// Calcola il fattore di conversione basato sulle componenti
	factor := 1.0
	for unit, exp := range components {
		baseUnit, ok := unit.(*BaseUnit)
		if !ok {
			panic("Component unit is not a BaseUnit")
		}
		factor *= math.Pow(baseUnit.factor, float64(exp))
	}

	return &DerivedUnit{
		BaseUnit: BaseUnit{
			unitType: unitType,
			name:     name,
			symbol:   symbol,
			factor:   factor,
			offset:   0, // Le unità derivate non hanno offset
		},
		components: components,
	}
}

// Unità di lunghezza
var (
	// Meter è il metro (unità SI di lunghezza)
	Meter = NewBaseUnit(Length, "meter", "m", 1.0, 0.0)
	// Kilometer è il chilometro
	Kilometer = NewBaseUnit(Length, "kilometer", "km", 1000.0, 0.0)
	// Centimeter è il centimetro
	Centimeter = NewBaseUnit(Length, "centimeter", "cm", 0.01, 0.0)
	// Millimeter è il millimetro
	Millimeter = NewBaseUnit(Length, "millimeter", "mm", 0.001, 0.0)
	// Inch è il pollice
	Inch = NewBaseUnit(Length, "inch", "in", 0.0254, 0.0)
	// Foot è il piede
	Foot = NewBaseUnit(Length, "foot", "ft", 0.3048, 0.0)
	// Mile è il miglio
	Mile = NewBaseUnit(Length, "mile", "mi", 1609.344, 0.0)
	// AstronomicalUnit è l'unità astronomica
	AstronomicalUnit = NewBaseUnit(Length, "astronomical unit", "AU", 1.495978707e11, 0.0)
	// LightYear è l'anno luce
	LightYear = NewBaseUnit(Length, "light year", "ly", 9.4607304725808e15, 0.0)
)

// Unità di massa
var (
	// Kilogram è il chilogrammo (unità SI di massa)
	Kilogram = NewBaseUnit(Mass, "kilogram", "kg", 1.0, 0.0)
	// Gram è il grammo
	Gram = NewBaseUnit(Mass, "gram", "g", 0.001, 0.0)
	// Milligram è il milligrammo
	Milligram = NewBaseUnit(Mass, "milligram", "mg", 1e-6, 0.0)
	// Tonne è la tonnellata
	Tonne = NewBaseUnit(Mass, "tonne", "t", 1000.0, 0.0)
	// Pound è la libbra
	Pound = NewBaseUnit(Mass, "pound", "lb", 0.45359237, 0.0)
	// SolarMass è la massa solare
	SolarMass = NewBaseUnit(Mass, "solar mass", "M☉", 1.989e30, 0.0)
)

// Unità di tempo
var (
	// Second è il secondo (unità SI di tempo)
	Second = NewBaseUnit(Time, "second", "s", 1.0, 0.0)
	// Minute è il minuto
	Minute = NewBaseUnit(Time, "minute", "min", 60.0, 0.0)
	// Hour è l'ora
	Hour = NewBaseUnit(Time, "hour", "h", 3600.0, 0.0)
	// Day è il giorno
	Day = NewBaseUnit(Time, "day", "d", 86400.0, 0.0)
	// Year è l'anno
	Year = NewBaseUnit(Time, "year", "yr", 31557600.0, 0.0) // Anno giuliano medio (365.25 giorni)
)

// Unità di temperatura
var (
	// Kelvin è il kelvin (unità SI di temperatura)
	Kelvin = NewBaseUnit(Temperature, "kelvin", "K", 1.0, 0.0)
	// Celsius è il grado Celsius
	Celsius = NewBaseUnit(Temperature, "Celsius", "°C", 1.0, 273.15)
	// Fahrenheit è il grado Fahrenheit
	Fahrenheit = NewBaseUnit(Temperature, "Fahrenheit", "°F", 5.0/9.0, 459.67)
)

// Unità di angolo
var (
	// Radian è il radiante (unità SI di angolo)
	Radian = NewBaseUnit(Angle, "radian", "rad", 1.0, 0.0)
	// Degree è il grado
	Degree = NewBaseUnit(Angle, "degree", "°", math.Pi/180.0, 0.0)
)

// Unità di forza
var (
	// Newton è il newton (unità SI di forza)
	Newton = NewDerivedUnit(Force, "newton", "N", map[Unit]int{
		Kilogram: 1,
		Meter:    1,
		Second:   -2,
	})
)

// Unità di energia
var (
	// Joule è il joule (unità SI di energia)
	Joule = NewDerivedUnit(Energy, "joule", "J", map[Unit]int{
		Kilogram: 1,
		Meter:    2,
		Second:   -2,
	})
)

// Unità di potenza
var (
	// Watt è il watt (unità SI di potenza)
	Watt = NewDerivedUnit(Power, "watt", "W", map[Unit]int{
		Kilogram: 1,
		Meter:    2,
		Second:   -3,
	})
)

// Unità di velocità
var (
	// MeterPerSecond è il metro al secondo (unità SI di velocità)
	MeterPerSecond = NewDerivedUnit(Velocity, "meter per second", "m/s", map[Unit]int{
		Meter:  1,
		Second: -1,
	})
	// KilometerPerHour è il chilometro all'ora
	KilometerPerHour = NewDerivedUnit(Velocity, "kilometer per hour", "km/h", map[Unit]int{
		Kilometer: 1,
		Hour:      -1,
	})
)

// Unità di accelerazione
var (
	// MeterPerSecondSquared è il metro al secondo quadrato (unità SI di accelerazione)
	MeterPerSecondSquared = NewDerivedUnit(Acceleration, "meter per second squared", "m/s²", map[Unit]int{
		Meter:  1,
		Second: -2,
	})
)

// Unità di pressione
var (
	// Pascal è il pascal (unità SI di pressione)
	Pascal = NewDerivedUnit(Pressure, "pascal", "Pa", map[Unit]int{
		Kilogram: 1,
		Meter:    -1,
		Second:   -2,
	})
)

// Quantity rappresenta una quantità fisica con un valore e un'unità
type Quantity struct {
	value float64
	unit  Unit
}

// NewQuantity crea una nuova quantità
func NewQuantity(value float64, unit Unit) Quantity {
	return Quantity{
		value: value,
		unit:  unit,
	}
}

// Value restituisce il valore della quantità
func (q Quantity) Value() float64 {
	return q.value
}

// Unit restituisce l'unità della quantità
func (q Quantity) Unit() Unit {
	return q.unit
}

// ConvertTo converte la quantità in un'altra unità
func (q Quantity) ConvertTo(unit Unit) Quantity {
	return NewQuantity(q.unit.ConvertTo(q.value, unit), unit)
}

// String restituisce una rappresentazione testuale della quantità
func (q Quantity) String() string {
	return fmt.Sprintf("%g %s", q.value, q.unit.Symbol())
}

// Add somma due quantità (convertendo se necessario)
func (q Quantity) Add(other Quantity) Quantity {
	if q.unit.Type() != other.unit.Type() {
		panic(fmt.Sprintf("Cannot add quantities of different types: %v and %v", q.unit.Type(), other.unit.Type()))
	}

	// Converti l'altra quantità nell'unità di questa
	otherValue := other.unit.ConvertTo(other.value, q.unit)
	return NewQuantity(q.value+otherValue, q.unit)
}

// Sub sottrae due quantità (convertendo se necessario)
func (q Quantity) Sub(other Quantity) Quantity {
	if q.unit.Type() != other.unit.Type() {
		panic(fmt.Sprintf("Cannot subtract quantities of different types: %v and %v", q.unit.Type(), other.unit.Type()))
	}

	// Converti l'altra quantità nell'unità di questa
	otherValue := other.unit.ConvertTo(other.value, q.unit)
	return NewQuantity(q.value-otherValue, q.unit)
}

// Mul moltiplica una quantità per uno scalare
func (q Quantity) Mul(scalar float64) Quantity {
	return NewQuantity(q.value*scalar, q.unit)
}

// Div divide una quantità per uno scalare
func (q Quantity) Div(scalar float64) Quantity {
	if scalar == 0 {
		panic("Division by zero")
	}
	return NewQuantity(q.value/scalar, q.unit)
}
