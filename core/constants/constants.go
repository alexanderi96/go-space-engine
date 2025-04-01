// Package constants fornisce costanti fisiche per il motore fisico
package constants

// Costanti fisiche universali
const (
	// G è la costante gravitazionale universale (m³/kg⋅s²)
	G = 6.67430e-11

	// SpeedOfLight è la velocità della luce nel vuoto (m/s)
	SpeedOfLight = 299792458.0

	// PlanckConstant è la costante di Planck (J⋅s)
	PlanckConstant = 6.62607015e-34

	// ReducedPlanckConstant è la costante di Planck ridotta (J⋅s)
	ReducedPlanckConstant = PlanckConstant / (2.0 * Pi)

	// BoltzmannConstant è la costante di Boltzmann (J/K)
	BoltzmannConstant = 1.380649e-23

	// StefanBoltzmannConstant è la costante di Stefan-Boltzmann (W/m²⋅K⁴)
	StefanBoltzmannConstant = 5.670374419e-8

	// ElectronCharge è la carica elementare dell'elettrone (C)
	ElectronCharge = 1.602176634e-19

	// VacuumPermittivity è la permittività del vuoto (F/m)
	VacuumPermittivity = 8.8541878128e-12

	// VacuumPermeability è la permeabilità del vuoto (H/m)
	VacuumPermeability = 1.25663706212e-6

	// AvogadroNumber è il numero di Avogadro (mol⁻¹)
	AvogadroNumber = 6.02214076e23

	// GasConstant è la costante universale dei gas (J/mol⋅K)
	GasConstant = 8.31446261815324

	// Pi è il rapporto tra la circonferenza e il diametro di un cerchio
	Pi = 3.14159265358979323846
)

// Costanti astronomiche
const (
	// SolarMass è la massa del Sole (kg)
	SolarMass = 1.989e30

	// EarthMass è la massa della Terra (kg)
	EarthMass = 5.972e24

	// EarthRadius è il raggio medio della Terra (m)
	EarthRadius = 6.371e6

	// AstronomicalUnit è l'unità astronomica (m)
	AstronomicalUnit = 1.495978707e11

	// LightYear è l'anno luce (m)
	LightYear = 9.4607304725808e15

	// Parsec è il parsec (m)
	Parsec = 3.0856775814671916e16
)

// Costanti per la simulazione
const (
	// DefaultTimeStep è il passo temporale predefinito per la simulazione (s)
	DefaultTimeStep = 0.01

	// DefaultGravity è l'accelerazione di gravità predefinita sulla Terra (m/s²)
	DefaultGravity = 9.80665

	// DefaultTemperature è la temperatura ambiente predefinita (K)
	DefaultTemperature = 293.15 // 20°C

	// DefaultPressure è la pressione atmosferica predefinita a livello del mare (Pa)
	DefaultPressure = 101325.0

	// Epsilon è un valore piccolo usato per confronti di uguaglianza tra float
	Epsilon = 1e-10
)
