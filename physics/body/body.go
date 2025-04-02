// Package body fornisce interfacce e implementazioni per i corpi fisici
package body

import (
	"github.com/alexanderi96/go-space-engine/core/units"
	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/google/uuid"
)

// ID rappresenta un identificatore univoco per un corpo
// type ID uuid.UUID

// NewID genera un nuovo ID univoco
// func NewID() ID {

// 	// Genera un ID casuale di 16 caratteri
// 	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
// 	result := make([]byte, 16)
// 	for i := range result {
// 		result[i] = chars[rand.Intn(len(chars))]
// 	}

// 	return ID(fmt.Sprintf("%s-%d", string(result), time.Now().UnixNano()))
// }

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
}

// Body rappresenta un corpo fisico nel motore
type Body interface {
	// ID restituisce l'identificatore univoco del corpo
	ID() uuid.UUID

	// Position restituisce la posizione del corpo
	Position() vector.Vector3
	// SetPosition imposta la posizione del corpo
	SetPosition(pos vector.Vector3)

	// Velocity restituisce la velocità del corpo
	Velocity() vector.Vector3
	// SetVelocity imposta la velocità del corpo
	SetVelocity(vel vector.Vector3)

	// Acceleration restituisce l'accelerazione del corpo
	Acceleration() vector.Vector3
	// SetAcceleration imposta l'accelerazione del corpo
	SetAcceleration(acc vector.Vector3)

	// Mass restituisce la massa del corpo
	Mass() units.Quantity
	// SetMass imposta la massa del corpo
	SetMass(mass units.Quantity)

	// Radius restituisce il raggio del corpo
	Radius() units.Quantity
	// SetRadius imposta il raggio del corpo
	SetRadius(radius units.Quantity)

	// Material restituisce il materiale del corpo
	Material() Material
	// SetMaterial imposta il materiale del corpo
	SetMaterial(mat Material)

	// ApplyForce applica una forza al corpo
	ApplyForce(force vector.Vector3)

	// Update aggiorna lo stato del corpo
	Update(dt float64)

	// Temperature restituisce la temperatura del corpo
	Temperature() units.Quantity
	// SetTemperature imposta la temperatura del corpo
	SetTemperature(temp units.Quantity)

	// AddHeat aggiunge calore al corpo
	AddHeat(heat units.Quantity)

	// IsStatic restituisce true se il corpo è statico (non si muove)
	IsStatic() bool
	// SetStatic imposta se il corpo è statico
	SetStatic(static bool)
}

// RigidBody implementa un corpo rigido
type RigidBody struct {
	id           uuid.UUID
	position     vector.Vector3
	velocity     vector.Vector3
	acceleration vector.Vector3
	mass         units.Quantity
	radius       units.Quantity
	material     Material
	temperature  units.Quantity
	isStatic     bool
}

// NewRigidBody crea un nuovo corpo rigido
func NewRigidBody(
	mass units.Quantity,
	radius units.Quantity,
	position vector.Vector3,
	velocity vector.Vector3,
	mat Material,
) *RigidBody {
	return &RigidBody{
		id:           uuid.New(),
		position:     position,
		velocity:     velocity,
		acceleration: vector.Zero3(),
		mass:         mass,
		radius:       radius,
		material:     mat,
		temperature:  units.NewQuantity(293.15, units.Kelvin), // Temperatura ambiente (20°C)
		isStatic:     false,
	}
}

// ID restituisce l'identificatore univoco del corpo
func (rb *RigidBody) ID() uuid.UUID {
	return rb.id
}

// Position restituisce la posizione del corpo
func (rb *RigidBody) Position() vector.Vector3 {
	return rb.position
}

// SetPosition imposta la posizione del corpo
func (rb *RigidBody) SetPosition(pos vector.Vector3) {
	rb.position = pos
}

// Velocity restituisce la velocità del corpo
func (rb *RigidBody) Velocity() vector.Vector3 {
	return rb.velocity
}

// SetVelocity imposta la velocità del corpo
func (rb *RigidBody) SetVelocity(vel vector.Vector3) {
	rb.velocity = vel

	// Se il corpo è statico, la velocità deve essere zero
	if rb.isStatic {
		rb.velocity = vector.Zero3()
	}
}

// Acceleration restituisce l'accelerazione del corpo
func (rb *RigidBody) Acceleration() vector.Vector3 {
	return rb.acceleration
}

// SetAcceleration imposta l'accelerazione del corpo
func (rb *RigidBody) SetAcceleration(acc vector.Vector3) {
	rb.acceleration = acc

	// Se il corpo è statico, l'accelerazione deve essere zero
	if rb.isStatic {
		rb.acceleration = vector.Zero3()
	}
}

// Mass restituisce la massa del corpo
func (rb *RigidBody) Mass() units.Quantity {
	return rb.mass
}

// SetMass imposta la massa del corpo
func (rb *RigidBody) SetMass(mass units.Quantity) {
	if mass.Unit().Type() != units.Mass {
		panic("Mass must be a mass quantity")
	}
	rb.mass = mass
}

// Radius restituisce il raggio del corpo
func (rb *RigidBody) Radius() units.Quantity {
	return rb.radius
}

// SetRadius imposta il raggio del corpo
func (rb *RigidBody) SetRadius(radius units.Quantity) {
	if radius.Unit().Type() != units.Length {
		panic("Radius must be a length quantity")
	}
	rb.radius = radius
}

// Material restituisce il materiale del corpo
func (rb *RigidBody) Material() Material {
	return rb.material
}

// SetMaterial imposta il materiale del corpo
func (rb *RigidBody) SetMaterial(mat Material) {
	rb.material = mat
}

// ApplyForce applica una forza al corpo
func (rb *RigidBody) ApplyForce(force vector.Vector3) {
	if rb.isStatic {
		return
	}

	// F = m*a => a = F/m
	massValue := rb.mass.Value()
	if massValue <= 0 {
		return
	}

	// Calcola l'accelerazione (F/m) e la aggiunge all'accelerazione corrente
	acceleration := force.Scale(1.0 / massValue)
	rb.acceleration = rb.acceleration.Add(acceleration)
}

// Update aggiorna lo stato del corpo
func (rb *RigidBody) Update(dt float64) {
	// Se il corpo è statico, assicurati che velocità e accelerazione siano zero
	if rb.isStatic {
		rb.velocity = vector.Zero3()
		rb.acceleration = vector.Zero3()
		return
	}

	// Integrazione di Verlet
	// x(t+dt) = x(t) + v(t)*dt + 0.5*a(t)*dt^2
	// v(t+dt) = v(t) + 0.5*(a(t) + a(t+dt))*dt

	// Salva l'accelerazione corrente
	oldAcceleration := rb.acceleration

	// Aggiorna la posizione
	halfDtSquared := 0.5 * dt * dt
	dtVelocity := rb.velocity.Scale(dt)
	dtAcceleration := rb.acceleration.Scale(halfDtSquared)
	rb.position = rb.position.Add(dtVelocity).Add(dtAcceleration)

	// Resetta l'accelerazione (verrà ricalcolata nel prossimo ciclo)
	rb.acceleration = vector.Zero3()

	// Aggiorna la velocità (usando la media delle accelerazioni)
	avgAcceleration := oldAcceleration.Add(rb.acceleration).Scale(0.5)
	rb.velocity = rb.velocity.Add(avgAcceleration.Scale(dt))
}

// Temperature restituisce la temperatura del corpo
func (rb *RigidBody) Temperature() units.Quantity {
	return rb.temperature
}

// SetTemperature imposta la temperatura del corpo
func (rb *RigidBody) SetTemperature(temp units.Quantity) {
	if temp.Unit().Type() != units.Temperature {
		panic("Temperature must be a temperature quantity")
	}
	rb.temperature = temp
}

// AddHeat aggiunge calore al corpo
func (rb *RigidBody) AddHeat(heat units.Quantity) {
	if heat.Unit().Type() != units.Energy {
		panic("Heat must be an energy quantity")
	}

	// Q = m*c*ΔT => ΔT = Q/(m*c)
	massValue := rb.mass.Value()
	if massValue <= 0 {
		return
	}

	// Ottieni la capacità termica specifica dal materiale
	specificHeat := rb.material.SpecificHeat()
	specificHeatValue := specificHeat.Value()

	// Calcola la variazione di temperatura
	heatValue := heat.Value()
	deltaTemp := heatValue / (massValue * specificHeatValue)

	// Aggiorna la temperatura
	currentTemp := rb.temperature.Value()
	newTemp := currentTemp + deltaTemp
	rb.temperature = units.NewQuantity(newTemp, units.Kelvin)
}

// IsStatic restituisce true se il corpo è statico (non si muove)
func (rb *RigidBody) IsStatic() bool {
	return rb.isStatic
}

// SetStatic imposta se il corpo è statico
func (rb *RigidBody) SetStatic(static bool) {
	rb.isStatic = static
	if static {
		rb.velocity = vector.Zero3()
		rb.acceleration = vector.Zero3()
	}
}
