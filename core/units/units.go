// Package units provides a measurement unit system for the physics engine
package units

import (
	"fmt"
	"math"
)

// UnitType represents the type of measurement unit
type UnitType int

const (
	// Length represents a unit of length
	Length UnitType = iota
	// Mass represents a unit of mass
	Mass
	// Time represents a unit of time
	Time
	// Temperature represents a unit of temperature
	Temperature
	// Angle represents a unit of angle
	Angle
	// Force represents a unit of force
	Force
	// Energy represents a unit of energy
	Energy
	// Power represents a unit of power
	Power
	// Velocity represents a unit of velocity
	Velocity
	// Acceleration represents a unit of acceleration
	Acceleration
	// Pressure represents a unit of pressure
	Pressure
)

// Unit represents a measurement unit
type Unit interface {
	// Type returns the unit type
	Type() UnitType
	// Name returns the unit name
	Name() string
	// Symbol returns the unit symbol
	Symbol() string
	// ConvertTo converts a value from this unit to another
	ConvertTo(value float64, target Unit) float64
	// ConvertFrom converts a value from another unit to this one
	ConvertFrom(value float64, source Unit) float64
}

// BaseUnit implements a base measurement unit
type BaseUnit struct {
	unitType UnitType
	name     string
	symbol   string
	factor   float64 // Conversion factor relative to SI unit
	offset   float64 // Offset for units with different zero point (e.g. temperatures)
}

// NewBaseUnit creates a new base unit
func NewBaseUnit(unitType UnitType, name, symbol string, factor, offset float64) *BaseUnit {
	return &BaseUnit{
		unitType: unitType,
		name:     name,
		symbol:   symbol,
		factor:   factor,
		offset:   offset,
	}
}

// Type returns the unit type
func (u *BaseUnit) Type() UnitType {
	return u.unitType
}

// Name returns the unit name
func (u *BaseUnit) Name() string {
	return u.name
}

// Symbol returns the unit symbol
func (u *BaseUnit) Symbol() string {
	return u.symbol
}

// ConvertTo converts a value from this unit to another
func (u *BaseUnit) ConvertTo(value float64, target Unit) float64 {
	if u.Type() != target.Type() {
		panic(fmt.Sprintf("Cannot convert between different unit types: %v and %v", u.Type(), target.Type()))
	}

	// First convert to SI unit
	siValue := (value + u.offset) * u.factor

	// Then convert from SI to target unit
	targetUnit, ok := target.(*BaseUnit)
	if !ok {
		panic("Target unit is not a BaseUnit")
	}

	return (siValue / targetUnit.factor) - targetUnit.offset
}

// ConvertFrom converts a value from another unit to this one
func (u *BaseUnit) ConvertFrom(value float64, source Unit) float64 {
	return source.(*BaseUnit).ConvertTo(value, u)
}

// DerivedUnit implements a derived measurement unit
type DerivedUnit struct {
	BaseUnit
	components map[Unit]int // Map of base units and their exponents
}

// NewDerivedUnit creates a new derived unit
func NewDerivedUnit(unitType UnitType, name, symbol string, components map[Unit]int) *DerivedUnit {
	// Calculate the conversion factor based on components
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
			offset:   0, // Derived units have no offset
		},
		components: components,
	}
}

// Length units
var (
	// Meter is the meter (SI unit of length)
	Meter = NewBaseUnit(Length, "meter", "m", 1.0, 0.0)
	// Kilometer is the kilometer
	Kilometer = NewBaseUnit(Length, "kilometer", "km", 1000.0, 0.0)
	// Centimeter is the centimeter
	Centimeter = NewBaseUnit(Length, "centimeter", "cm", 0.01, 0.0)
	// Millimeter is the millimeter
	Millimeter = NewBaseUnit(Length, "millimeter", "mm", 0.001, 0.0)
	// Inch is the inch
	Inch = NewBaseUnit(Length, "inch", "in", 0.0254, 0.0)
	// Foot is the foot
	Foot = NewBaseUnit(Length, "foot", "ft", 0.3048, 0.0)
	// Mile is the mile
	Mile = NewBaseUnit(Length, "mile", "mi", 1609.344, 0.0)
	// AstronomicalUnit is the astronomical unit
	AstronomicalUnit = NewBaseUnit(Length, "astronomical unit", "AU", 1.495978707e11, 0.0)
	// LightYear is the light year
	LightYear = NewBaseUnit(Length, "light year", "ly", 9.4607304725808e15, 0.0)
)

// Mass units
var (
	// Kilogram is the kilogram (SI unit of mass)
	Kilogram = NewBaseUnit(Mass, "kilogram", "kg", 1.0, 0.0)
	// Gram is the gram
	Gram = NewBaseUnit(Mass, "gram", "g", 0.001, 0.0)
	// Milligram is the milligram
	Milligram = NewBaseUnit(Mass, "milligram", "mg", 1e-6, 0.0)
	// Tonne is the tonne
	Tonne = NewBaseUnit(Mass, "tonne", "t", 1000.0, 0.0)
	// Pound is the pound
	Pound = NewBaseUnit(Mass, "pound", "lb", 0.45359237, 0.0)
	// SolarMass is the solar mass
	SolarMass = NewBaseUnit(Mass, "solar mass", "M☉", 1.989e30, 0.0)
)

// Time units
var (
	// Second is the second (SI unit of time)
	Second = NewBaseUnit(Time, "second", "s", 1.0, 0.0)
	// Minute is the minute
	Minute = NewBaseUnit(Time, "minute", "min", 60.0, 0.0)
	// Hour is the hour
	Hour = NewBaseUnit(Time, "hour", "h", 3600.0, 0.0)
	// Day is the day
	Day = NewBaseUnit(Time, "day", "d", 86400.0, 0.0)
	// Year is the year
	Year = NewBaseUnit(Time, "year", "yr", 31557600.0, 0.0) // Julian mean year (365.25 days)
)

// Temperature units
var (
	// Kelvin is the kelvin (SI unit of temperature)
	Kelvin = NewBaseUnit(Temperature, "kelvin", "K", 1.0, 0.0)
	// Celsius is the degree Celsius
	Celsius = NewBaseUnit(Temperature, "Celsius", "°C", 1.0, 273.15)
	// Fahrenheit is the degree Fahrenheit
	Fahrenheit = NewBaseUnit(Temperature, "Fahrenheit", "°F", 5.0/9.0, 459.67)
)

// Angle units
var (
	// Radian is the radian (SI unit of angle)
	Radian = NewBaseUnit(Angle, "radian", "rad", 1.0, 0.0)
	// Degree is the degree
	Degree = NewBaseUnit(Angle, "degree", "°", math.Pi/180.0, 0.0)
)

// Force units
var (
	// Newton is the newton (SI unit of force)
	Newton = NewDerivedUnit(Force, "newton", "N", map[Unit]int{
		Kilogram: 1,
		Meter:    1,
		Second:   -2,
	})
)

// Energy units
var (
	// Joule is the joule (SI unit of energy)
	Joule = NewDerivedUnit(Energy, "joule", "J", map[Unit]int{
		Kilogram: 1,
		Meter:    2,
		Second:   -2,
	})
)

// Power units
var (
	// Watt is the watt (SI unit of power)
	Watt = NewDerivedUnit(Power, "watt", "W", map[Unit]int{
		Kilogram: 1,
		Meter:    2,
		Second:   -3,
	})
)

// Velocity units
var (
	// MeterPerSecond is the meter per second (SI unit of velocity)
	MeterPerSecond = NewDerivedUnit(Velocity, "meter per second", "m/s", map[Unit]int{
		Meter:  1,
		Second: -1,
	})
	// KilometerPerHour is the kilometer per hour
	KilometerPerHour = NewDerivedUnit(Velocity, "kilometer per hour", "km/h", map[Unit]int{
		Kilometer: 1,
		Hour:      -1,
	})
)

// Acceleration units
var (
	// MeterPerSecondSquared is the meter per second squared (SI unit of acceleration)
	MeterPerSecondSquared = NewDerivedUnit(Acceleration, "meter per second squared", "m/s²", map[Unit]int{
		Meter:  1,
		Second: -2,
	})
)

// Pressure units
var (
	// Pascal is the pascal (SI unit of pressure)
	Pascal = NewDerivedUnit(Pressure, "pascal", "Pa", map[Unit]int{
		Kilogram: 1,
		Meter:    -1,
		Second:   -2,
	})
)

// Quantity represents a physical quantity with a value and a unit
type Quantity struct {
	value float64
	unit  Unit
}

// NewQuantity creates a new quantity
func NewQuantity(value float64, unit Unit) Quantity {
	return Quantity{
		value: value,
		unit:  unit,
	}
}

// Value returns the value of the quantity
func (q Quantity) Value() float64 {
	return q.value
}

// Unit returns the unit of the quantity
func (q Quantity) Unit() Unit {
	return q.unit
}

// ConvertTo converts the quantity to another unit
func (q Quantity) ConvertTo(unit Unit) Quantity {
	return NewQuantity(q.unit.ConvertTo(q.value, unit), unit)
}

// String returns a textual representation of the quantity
func (q Quantity) String() string {
	return fmt.Sprintf("%g %s", q.value, q.unit.Symbol())
}

// Add sums two quantities (converting if necessary)
func (q Quantity) Add(other Quantity) Quantity {
	if q.unit.Type() != other.unit.Type() {
		panic(fmt.Sprintf("Cannot add quantities of different types: %v and %v", q.unit.Type(), other.unit.Type()))
	}

	// Convert the other quantity to this unit
	otherValue := other.unit.ConvertTo(other.value, q.unit)
	return NewQuantity(q.value+otherValue, q.unit)
}

// Sub subtracts two quantities (converting if necessary)
func (q Quantity) Sub(other Quantity) Quantity {
	if q.unit.Type() != other.unit.Type() {
		panic(fmt.Sprintf("Cannot subtract quantities of different types: %v and %v", q.unit.Type(), other.unit.Type()))
	}

	// Convert the other quantity to this unit
	otherValue := other.unit.ConvertTo(other.value, q.unit)
	return NewQuantity(q.value-otherValue, q.unit)
}

// Mul multiplies a quantity by a scalar
func (q Quantity) Mul(scalar float64) Quantity {
	return NewQuantity(q.value*scalar, q.unit)
}

// Div divides a quantity by a scalar
func (q Quantity) Div(scalar float64) Quantity {
	if scalar == 0 {
		panic("Division by zero")
	}
	return NewQuantity(q.value/scalar, q.unit)
}

// ConvertToStandardUnit converts a quantity to a standard unit based on its type
// This ensures consistent rendering regardless of the original unit
func ConvertToStandardUnit(quantity Quantity) float64 {
	switch quantity.Unit().Type() {
	case Length:
		// Convert all lengths to meters
		return quantity.ConvertTo(Meter).Value()
	case Mass:
		// Convert all masses to kilograms
		return quantity.ConvertTo(Kilogram).Value()
	case Temperature:
		// Convert all temperatures to kelvin
		return quantity.ConvertTo(Kelvin).Value()
	default:
		// For other types, just return the value
		return quantity.Value()
	}
}
