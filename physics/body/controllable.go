// Package body provides interfaces and implementations for physical bodies
package body

import (
	"math"

	"github.com/alexanderi96/go-space-engine/core/vector"
)

// ControllableBody rappresenta un corpo che può essere controllato dall'utente
type ControllableBody interface {
	Body

	// IsControllable restituisce true se il corpo è attualmente controllabile
	IsControllable() bool

	// SetControllable imposta se il corpo è controllabile
	SetControllable(controllable bool)

	// HandleInput gestisce l'input dell'utente per controllare il corpo
	// deltaTime è il tempo trascorso dall'ultimo frame
	// moveForward, moveBackward, moveLeft, moveRight, moveUp, moveDown sono i comandi di movimento
	// rotateLeft, rotateRight, rotateUp, rotateDown sono i comandi di rotazione
	HandleInput(deltaTime float64, moveForward, moveBackward, moveLeft, moveRight, moveUp, moveDown, rotateLeft, rotateRight, rotateUp, rotateDown bool)
}

// ControllableRigidBody implementa un corpo rigido controllabile
type ControllableRigidBody struct {
	*RigidBody

	isControllable bool
	moveSpeed      float64 // Velocità di movimento (forza di propulsione)
	rotateSpeed    float64 // Velocità di rotazione (radianti/secondo)
}

// NewControllableRigidBody crea un nuovo corpo rigido controllabile
func NewControllableRigidBody(rb *RigidBody, moveSpeed, rotateSpeed float64) *ControllableRigidBody {
	return &ControllableRigidBody{
		RigidBody:      rb,
		isControllable: true,
		moveSpeed:      moveSpeed,
		rotateSpeed:    rotateSpeed,
	}
}

// IsControllable restituisce true se il corpo è attualmente controllabile
func (crb *ControllableRigidBody) IsControllable() bool {
	return crb.isControllable
}

// SetControllable imposta se il corpo è controllabile
func (crb *ControllableRigidBody) SetControllable(controllable bool) {
	crb.isControllable = controllable
}

// HandleInput gestisce l'input dell'utente per controllare il corpo
func (crb *ControllableRigidBody) HandleInput(deltaTime float64, moveForward, moveBackward, moveLeft, moveRight, moveUp, moveDown, rotateLeft, rotateRight, rotateUp, rotateDown bool) {
	if !crb.isControllable {
		return
	}

	// Calcola la direzione di movimento basata sulla rotazione del corpo
	forward := vector.NewVector3(0, 0, -1) // Direzione iniziale (verso lo schermo)

	// Applica la rotazione del corpo alla direzione forward
	rotation := crb.Rotation()

	// Calcola la matrice di rotazione (semplificata per questo esempio)
	// In una implementazione completa, dovresti usare quaternioni o matrici di rotazione
	cosY := float64(0.0)
	sinY := float64(0.0)

	// Evita calcoli NaN
	if !math.IsNaN(rotation.Y()) {
		cosY = math.Cos(rotation.Y())
		sinY = math.Sin(rotation.Y())
	}

	// Ruota il vettore forward attorno all'asse Y
	rotatedForward := vector.NewVector3(
		forward.X()*cosY-forward.Z()*sinY,
		forward.Y(),
		forward.X()*sinY+forward.Z()*cosY,
	)

	// Calcola il vettore right (perpendicolare a forward)
	right := vector.NewVector3(
		rotatedForward.Z(),
		0,
		-rotatedForward.X(),
	)

	// Calcola il vettore up (perpendicolare a forward e right)
	up := vector.NewVector3(0, 1, 0)

	// Calcola la forza di propulsione basata sull'input
	thrustForce := vector.Zero3()

	// Applica i comandi di movimento
	if moveForward {
		thrustForce = thrustForce.Add(rotatedForward)
	}
	if moveBackward {
		thrustForce = thrustForce.Add(rotatedForward.Scale(-1))
	}
	if moveRight {
		thrustForce = thrustForce.Add(right)
	}
	if moveLeft {
		thrustForce = thrustForce.Add(right.Scale(-1))
	}
	if moveUp {
		thrustForce = thrustForce.Add(up)
	}
	if moveDown {
		thrustForce = thrustForce.Add(up.Scale(-1))
	}

	// Normalizza la forza di propulsione se non è zero
	if thrustForce.Length() > 0 {
		thrustForce = thrustForce.Normalize()

		// Applica la forza di propulsione
		thrustForce = thrustForce.Scale(crb.moveSpeed)

		// Applica la forza al corpo
		crb.ApplyForce(thrustForce)
	}

	// TODO: Implementare la stabilizzazione
	// Quando non ci sono input di movimento, applicare una forza contraria alla velocità
	// per rallentare gradualmente il corpo

	// TODO: Implementare limiti di velocità
	// Controllare la velocità corrente e limitarla se supera un valore massimo

	// Gestisci la rotazione
	torque := vector.Zero3()

	if rotateLeft {
		torque = torque.Add(vector.NewVector3(0, 1, 0))
	}
	if rotateRight {
		torque = torque.Add(vector.NewVector3(0, -1, 0))
	}
	if rotateUp {
		torque = torque.Add(vector.NewVector3(1, 0, 0))
	}
	if rotateDown {
		torque = torque.Add(vector.NewVector3(-1, 0, 0))
	}

	// Applica il torque se non è zero
	if torque.Length() > 0 {
		torque = torque.Scale(crb.rotateSpeed)

		// Applica il torque al corpo
		crb.ApplyTorque(torque)
	}
}
