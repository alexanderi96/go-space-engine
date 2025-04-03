// Package constants provides physical constants for the physics engine
package constants

// Universal physical constants
const (
	// G is the universal gravitational constant (m³/kg⋅s²)
	G = 6.67430e-11

	// SpeedOfLight is the speed of light in vacuum (m/s)
	SpeedOfLight = 299792458.0

	// PlanckConstant is the Planck constant (J⋅s)
	PlanckConstant = 6.62607015e-34

	// ReducedPlanckConstant is the reduced Planck constant (J⋅s)
	ReducedPlanckConstant = PlanckConstant / (2.0 * Pi)

	// BoltzmannConstant is the Boltzmann constant (J/K)
	BoltzmannConstant = 1.380649e-23

	// StefanBoltzmannConstant is the Stefan-Boltzmann constant (W/m²⋅K⁴)
	StefanBoltzmannConstant = 5.670374419e-8

	// ElectronCharge is the elementary charge of the electron (C)
	ElectronCharge = 1.602176634e-19

	// VacuumPermittivity is the permittivity of vacuum (F/m)
	VacuumPermittivity = 8.8541878128e-12

	// VacuumPermeability is the permeability of vacuum (H/m)
	VacuumPermeability = 1.25663706212e-6

	// AvogadroNumber is Avogadro's number (mol⁻¹)
	AvogadroNumber = 6.02214076e23

	// GasConstant is the universal gas constant (J/mol⋅K)
	GasConstant = 8.31446261815324

	// Pi is the ratio of a circle's circumference to its diameter
	Pi = 3.14159265358979323846
)

// Astronomical constants
const (
	// SolarMass is the mass of the Sun (kg)
	SolarMass = 1.989e30

	// EarthMass is the mass of the Earth (kg)
	EarthMass = 5.972e24

	// EarthRadius is the average radius of the Earth (m)
	EarthRadius = 6.371e6

	// AstronomicalUnit is the astronomical unit (m)
	AstronomicalUnit = 1.495978707e11

	// LightYear is the light year (m)
	LightYear = 9.4607304725808e15

	// Parsec is the parsec (m)
	Parsec = 3.0856775814671916e16
)

// Simulation constants
const (
	// DefaultTimeStep is the default time step for simulation (s)
	DefaultTimeStep = 0.01

	// DefaultGravity is the default gravity acceleration on Earth (m/s²)
	DefaultGravity = 9.80665

	// DefaultTemperature is the default ambient temperature (K)
	DefaultTemperature = 293.15 // 20°C

	// DefaultPressure is the default atmospheric pressure at sea level (Pa)
	DefaultPressure = 101325.0

	// Epsilon is a small value used for float equality comparisons
	Epsilon = 1e-10
)
