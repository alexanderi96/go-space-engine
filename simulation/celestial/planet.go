package celestial

import (
	"math"
	"math/rand"

	"github.com/alexanderi96/go-space-engine/core/units"
	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/physics/body"
	"github.com/alexanderi96/go-space-engine/physics/force"
	physMaterial "github.com/alexanderi96/go-space-engine/physics/material"
	"github.com/alexanderi96/go-space-engine/simulation/world"
)

// BodyType represents the type of celestial body
type BodyType int

const (
	// Star represents a star
	Star BodyType = iota
	// TerrestrialPlanet represents a small rocky planet
	TerrestrialPlanet
	// GasGiant represents a large gas giant
	GasGiant
	// IceGiant represents an ice giant
	IceGiant
	// DwarfPlanet represents a dwarf planet
	DwarfPlanet
	// Moon represents a moon
	Moon
	// Asteroid represents an asteroid
	Asteroid
	// Comet represents a comet
	Comet
)

// PlanetParams defines parameters for creating a planet
type PlanetParams struct {
	// Name of the planet
	Name string

	// Physical parameters
	Mass   float64  // Mass in kg
	Radius float64  // Radius in m
	Type   BodyType // Type of celestial body

	// Orbital parameters
	Distance     float64        // Distance from central body
	Eccentricity float64        // Orbital eccentricity (0-1)
	Inclination  float64        // Orbital inclination in radians
	InitialAngle float64        // Initial angle in orbit (radians)
	OrbitPlane   vector.Vector3 // Orbital plane normal vector

	// Position and velocity (overrides orbital parameters if set)
	Position vector.Vector3 // Direct position setting
	Velocity vector.Vector3 // Direct velocity setting

	// Material
	Material physMaterial.Material // Material of the body

	// Reference
	CentralBody body.Body // Central body for orbital calculations
}

// DefaultPlanetParams returns default parameters for a planet
func DefaultPlanetParams() PlanetParams {
	return PlanetParams{
		Name:         "Planet",
		Mass:         5.972e24, // Earth mass
		Radius:       6371.0,   // Earth radius
		Type:         TerrestrialPlanet,
		Distance:     149.6e9, // 1 AU
		Eccentricity: 0.0,     // Circular orbit
		Inclination:  0.0,     // No inclination
		InitialAngle: 0.0,     // Start at angle 0
		OrbitPlane:   vector.NewVector3(0, 1, 0),
		Position:     vector.Zero3(), // Will be calculated from orbital params if not set
		Velocity:     vector.Zero3(), // Will be calculated from orbital params if not set
		Material:     physMaterial.Rock,
		CentralBody:  nil,
	}
}

// GeneratePlanet creates a celestial body with the given parameters
func GeneratePlanet(w world.World, params PlanetParams) body.Body {
	// If no material provided, select a default based on type
	mat := params.Material
	if mat == nil {
		mat = getDefaultMaterial(params.Type)
	}

	// If no central body provided but distance is set, create an object with no orbital velocity
	centralMass := 0.0
	if params.CentralBody != nil {
		centralMass = params.CentralBody.Mass().Value()
	}

	// Get position - either use provided position or calculate from orbital parameters
	position := params.Position
	velocity := params.Velocity

	// If position is zero, calculate it from orbital parameters
	if position.X() == 0 && position.Y() == 0 && position.Z() == 0 && params.Distance > 0 {
		// Apply eccentricity to the distance to get actual position
		// For an elliptical orbit, distance is treated as the semi-major axis
		radius := params.Distance
		if params.Eccentricity > 0 {
			// Adjust for eccentricity (approximate - this doesn't handle true anomaly correctly)
			// For a proper implementation, would need a more complex orbital mechanics calculation
			radius = params.Distance * (1 - params.Eccentricity*math.Cos(params.InitialAngle))
		}

		// Calculate position based on distance and angle
		// First calculate position in orbital plane
		x := radius * math.Cos(params.InitialAngle)
		z := radius * math.Sin(params.InitialAngle)
		position = vector.NewVector3(x, 0, z)

		// If orbital plane is not the XZ plane, rotate the position
		if params.OrbitPlane.X() != 0 || params.OrbitPlane.Z() != 0 || params.OrbitPlane.Y() != 1 {
			// This would need a proper quaternion rotation implementation
			// For simplicity we'll use the current implementation, but note
			// that a full 3D rotation would be better
		}

		// Apply inclination
		if params.Inclination != 0 {
			// Rotate position around the X axis by inclination angle
			// This is simplified - a proper implementation would use quaternions
			cosInc := math.Cos(params.Inclination)
			sinInc := math.Sin(params.Inclination)
			y := z * sinInc
			z = z * cosInc
			position = vector.NewVector3(x, y, z)
		}
	}

	// Calculate orbital velocity if not provided
	if velocity.X() == 0 && velocity.Y() == 0 && velocity.Z() == 0 && centralMass > 0 && params.Distance > 0 {
		// Calculate distance for velocity calculation
		radius := position.Length()

		// Calculate base velocity for circular orbit
		baseSpeed := force.CalculateOrbitalVelocity(centralMass, radius)

		// Adjust for eccentricity
		if params.Eccentricity > 0 {
			// In elliptical orbit, velocity is higher at perihelion and lower at aphelion
			// This is a simplified adjustment
			velocityFactor := math.Sqrt((1 + params.Eccentricity*math.Cos(params.InitialAngle)) /
				(1 - params.Eccentricity*math.Cos(params.InitialAngle)))
			baseSpeed *= velocityFactor
		}

		// Direction of velocity is perpendicular to position vector in orbital plane
		// For a circular orbit in the XZ plane:
		velX := -baseSpeed * math.Sin(params.InitialAngle)
		velZ := baseSpeed * math.Cos(params.InitialAngle)
		velocity = vector.NewVector3(velX, 0, velZ)

		// Apply inclination to velocity
		if params.Inclination != 0 {
			// Rotate velocity around the X axis (simplified)
			cosInc := math.Cos(params.Inclination)
			sinInc := math.Sin(params.Inclination)
			velY := velZ * sinInc
			velZ = velZ * cosInc
			velocity = vector.NewVector3(velX, velY, velZ)
		}
	}

	// Calculate scaled radius for visualization
	// This scales down the actual radius for better visualization
	// For a realistic simulation these would be actual values
	visualRadius := scaleRadius(params.Radius, params.Type)

	// Create the body
	b := body.NewRigidBody(
		units.NewQuantity(params.Mass, units.Kilogram),
		units.NewQuantity(visualRadius, units.Meter),
		position,
		velocity,
		mat,
	)

	// Add the body to the world
	w.AddBody(b)
	return b
}

// GenerateStarSystem creates a star with planets
func GenerateStarSystem(w world.World, planets int) body.Body {
	// Create the star
	starParams := DefaultPlanetParams()
	starParams.Name = "Star"
	starParams.Type = Star
	starParams.Mass = 1.989e30   // Solar mass
	starParams.Radius = 695700.0 // Solar radius

	star := GeneratePlanet(w, starParams)

	// Create planets
	for i := 0; i < planets; i++ {
		// Calculate distance with logarithmic spacing
		distanceFactor := 0.4 + 0.3*float64(i) // Gives a reasonable progression
		distance := 20.0 * math.Pow(10, distanceFactor)

		// Decide planet type based on distance
		planetType := TerrestrialPlanet
		if i > planets/2 {
			if rand.Float64() > 0.5 {
				planetType = GasGiant
			} else {
				planetType = IceGiant
			}
		}

		// Calculate mass based on type
		mass := 0.0
		switch planetType {
		case TerrestrialPlanet:
			mass = 1e24 * (0.1 + rand.Float64()*1.9) // 0.1 to 2 Earth masses
		case GasGiant:
			mass = 1e27 * (0.1 + rand.Float64()*1.9) // 0.1 to 2 Jupiter masses
		case IceGiant:
			mass = 1e26 * (0.5 + rand.Float64()) // 0.5 to 1.5 Neptune masses
		}

		// Calculate radius based on type and mass
		radius := 0.0
		switch planetType {
		case TerrestrialPlanet:
			radius = 6371.0 * math.Pow(mass/5.972e24, 0.3) // Scaled to Earth
		case GasGiant:
			radius = 69911.0 * math.Pow(mass/1.898e27, 0.3) // Scaled to Jupiter
		case IceGiant:
			radius = 24622.0 * math.Pow(mass/1.024e26, 0.3) // Scaled to Neptune
		}

		// Add some orbital inclination and eccentricity
		inclination := (rand.Float64() * 10.0) * (math.Pi / 180.0) // 0-10 degrees
		eccentricity := rand.Float64() * 0.1                       // 0-0.1

		// Create planet
		planetParams := DefaultPlanetParams()
		planetParams.Name = "Planet-" + string(rune(i+65)) // A, B, C, ...
		planetParams.Type = planetType
		planetParams.Mass = mass
		planetParams.Radius = radius
		planetParams.Distance = distance
		planetParams.Eccentricity = eccentricity
		planetParams.Inclination = inclination
		planetParams.InitialAngle = rand.Float64() * 2.0 * math.Pi
		planetParams.CentralBody = star

		planet := GeneratePlanet(w, planetParams)

		// Potentially add moons
		moons := 0
		switch planetType {
		case TerrestrialPlanet:
			moons = int(rand.Float64() * 2) // 0-1 moons
		case GasGiant:
			moons = int(rand.Float64() * 12) // 0-11 moons
		case IceGiant:
			moons = int(rand.Float64() * 6) // 0-5 moons
		}

		// Create moons
		for m := 0; m < moons; m++ {
			GenerateMoon(w, planet)
		}
	}

	return star
}

// GenerateMoon creates a moon orbiting a planet
func GenerateMoon(w world.World, planet body.Body) body.Body {
	// Get planet properties
	planetMass := planet.Mass().Value()
	planetRadius := planet.Radius().Value()

	// Generate moon parameters
	moonMass := planetMass * (0.01 + rand.Float64()*0.05)   // 1-6% of planet mass
	moonRadius := planetRadius * (0.1 + rand.Float64()*0.3) // 10-40% of planet radius

	// Calculate orbital distance (3-15 planet radii)
	distance := planetRadius * (3.0 + rand.Float64()*12.0)

	// Calculate moon parameters
	moonParams := DefaultPlanetParams()
	moonParams.Name = "Moon"
	moonParams.Type = Moon
	moonParams.Mass = moonMass
	moonParams.Radius = moonRadius
	moonParams.Distance = distance
	moonParams.Eccentricity = rand.Float64() * 0.1                       // 0-0.1
	moonParams.Inclination = (rand.Float64() * 20.0) * (math.Pi / 180.0) // 0-20 degrees
	moonParams.InitialAngle = rand.Float64() * 2.0 * math.Pi
	moonParams.CentralBody = planet

	return GeneratePlanet(w, moonParams)
}

// getDefaultMaterial returns a default material based on body type
func getDefaultMaterial(bodyType BodyType) physMaterial.Material {
	// Create custom materials for each body type
	switch bodyType {
	case Star:
		return createStarMaterial()
	case TerrestrialPlanet:
		return createTerrestrialMaterial()
	case GasGiant:
		return createGasGiantMaterial()
	case IceGiant:
		return createIceGiantMaterial()
	case DwarfPlanet:
		return createDwarfPlanetMaterial()
	case Moon:
		return createMoonMaterial()
	case Asteroid:
		return physMaterial.Rock
	case Comet:
		return createCometMaterial()
	default:
		return physMaterial.Rock
	}
}

// createStarMaterial creates a material for stars
func createStarMaterial() physMaterial.Material {
	return physMaterial.NewBasicMaterial(
		"Sun",
		units.NewQuantity(5000, units.Kilogram),
		units.NewQuantity(1000, units.Joule),
		units.NewQuantity(10.0, units.Watt),
		0.9,                       // High emissivity
		0.5,                       // Medium elasticity
		[3]float64{1.0, 0.8, 0.0}, // Yellow color
	)
}

// createTerrestrialMaterial creates a material for terrestrial planets
func createTerrestrialMaterial() physMaterial.Material {
	return physMaterial.NewBasicMaterial(
		"Earth",
		units.NewQuantity(5000, units.Kilogram),
		units.NewQuantity(800, units.Joule),
		units.NewQuantity(1.5, units.Watt),
		0.7,                       // Medium emissivity
		0.5,                       // Medium elasticity
		[3]float64{0.0, 0.3, 0.8}, // Blue-green color
	)
}

// createGasGiantMaterial creates a material for gas giants
func createGasGiantMaterial() physMaterial.Material {
	return physMaterial.NewBasicMaterial(
		"Jupiter",
		units.NewQuantity(1000, units.Kilogram), // Lower density
		units.NewQuantity(800, units.Joule),
		units.NewQuantity(1.5, units.Watt),
		0.7,                       // Medium emissivity
		0.3,                       // Low elasticity (gas is compressible)
		[3]float64{0.8, 0.6, 0.4}, // Jupiter-like color
	)
}

// createIceGiantMaterial creates a material for ice giants
func createIceGiantMaterial() physMaterial.Material {
	return physMaterial.NewBasicMaterial(
		"Neptune",
		units.NewQuantity(2000, units.Kilogram),
		units.NewQuantity(800, units.Joule),
		units.NewQuantity(1.5, units.Watt),
		0.7,                       // Medium emissivity
		0.4,                       // Medium-low elasticity
		[3]float64{0.0, 0.0, 0.8}, // Blue color
	)
}

// createDwarfPlanetMaterial creates a material for dwarf planets
func createDwarfPlanetMaterial() physMaterial.Material {
	return physMaterial.NewBasicMaterial(
		"DwarfPlanet",
		units.NewQuantity(5000, units.Kilogram),
		units.NewQuantity(800, units.Joule),
		units.NewQuantity(1.5, units.Watt),
		0.7,                       // Medium emissivity
		0.5,                       // Medium elasticity
		[3]float64{0.7, 0.7, 0.7}, // Grey color
	)
}

// createMoonMaterial creates a material for moons
func createMoonMaterial() physMaterial.Material {
	return physMaterial.NewBasicMaterial(
		"Moon",
		units.NewQuantity(5000, units.Kilogram),
		units.NewQuantity(800, units.Joule),
		units.NewQuantity(1.5, units.Watt),
		0.7,                       // Medium emissivity
		0.5,                       // Medium elasticity
		[3]float64{0.7, 0.7, 0.7}, // Grey color
	)
}

// createCometMaterial creates a material for comets
func createCometMaterial() physMaterial.Material {
	return physMaterial.NewBasicMaterial(
		"Comet",
		units.NewQuantity(3000, units.Kilogram),
		units.NewQuantity(800, units.Joule),
		units.NewQuantity(1.5, units.Watt),
		0.7,                       // Medium emissivity
		0.3,                       // Low elasticity (ice/dust mix)
		[3]float64{0.8, 0.8, 0.9}, // Icy blue-white color
	)
}

// scaleRadius scales the actual radius for visualization purposes
func scaleRadius(actualRadius float64, bodyType BodyType) float64 {
	// Calculate a visualization-friendly radius
	// This is necessary because actual astronomical scale differences
	// would make visualization difficult (e.g., a solar system to scale)

	// Base scale factor
	scaleFactor := 0.0

	switch bodyType {
	case Star:
		// Stars are scaled down more to allow better visualization
		scaleFactor = 1.0 / 100000.0
		return 5.0 + actualRadius*scaleFactor // Min 5.0 size
	case GasGiant:
		scaleFactor = 1.0 / 50000.0
		return 2.5 + actualRadius*scaleFactor
	case IceGiant:
		scaleFactor = 1.0 / 20000.0
		return 2.0 + actualRadius*scaleFactor
	case TerrestrialPlanet:
		scaleFactor = 1.0 / 5000.0
		return 1.0 + actualRadius*scaleFactor
	case DwarfPlanet:
		scaleFactor = 1.0 / 2000.0
		return 0.5 + actualRadius*scaleFactor
	case Moon:
		scaleFactor = 1.0 / 1000.0
		return 0.3 + actualRadius*scaleFactor
	case Asteroid:
		scaleFactor = 1.0 / 100.0
		return 0.1 + actualRadius*scaleFactor
	case Comet:
		scaleFactor = 1.0 / 100.0
		return 0.1 + actualRadius*scaleFactor
	default:
		scaleFactor = 1.0 / 10000.0
		return 1.0 + actualRadius*scaleFactor
	}
}
