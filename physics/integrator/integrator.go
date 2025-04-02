// Package integrator fornisce integratori numerici per le equazioni del moto
package integrator

import (
	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/physics/body"
	"github.com/google/uuid"
)

// Integrator rappresenta un integratore numerico per le equazioni del moto
type Integrator interface {
	// Integrate integra le equazioni del moto per un corpo
	Integrate(b body.Body, dt float64)
	// IntegrateAll integra le equazioni del moto per tutti i corpi
	IntegrateAll(bodies []body.Body, dt float64)
}

// EulerIntegrator implementa l'integratore di Euler
type EulerIntegrator struct{}

// NewEulerIntegrator crea un nuovo integratore di Euler
func NewEulerIntegrator() *EulerIntegrator {
	return &EulerIntegrator{}
}

// Integrate integra le equazioni del moto per un corpo usando il metodo di Euler
func (ei *EulerIntegrator) Integrate(b body.Body, dt float64) {
	// Se il corpo è statico, non fare nulla
	if b.IsStatic() {
		return
	}

	// Aggiorna la posizione: x(t+dt) = x(t) + v(t)*dt
	newPosition := b.Position().Add(b.Velocity().Scale(dt))
	b.SetPosition(newPosition)

	// Aggiorna la velocità: v(t+dt) = v(t) + a(t)*dt
	newVelocity := b.Velocity().Add(b.Acceleration().Scale(dt))
	b.SetVelocity(newVelocity)

	// Resetta l'accelerazione (verrà ricalcolata nel prossimo ciclo)
	b.SetAcceleration(vector.Zero3())
}

// IntegrateAll integra le equazioni del moto per tutti i corpi usando il metodo di Euler
func (ei *EulerIntegrator) IntegrateAll(bodies []body.Body, dt float64) {
	for _, b := range bodies {
		ei.Integrate(b, dt)
	}
}

// VerletIntegrator implementa l'integratore di Verlet
type VerletIntegrator struct {
	// Mappa che memorizza le posizioni precedenti dei corpi
	previousPositions map[uuid.UUID]vector.Vector3
}

// NewVerletIntegrator crea un nuovo integratore di Verlet
func NewVerletIntegrator() *VerletIntegrator {
	return &VerletIntegrator{
		previousPositions: make(map[uuid.UUID]vector.Vector3),
	}
}

// Integrate integra le equazioni del moto per un corpo usando il metodo di Verlet
func (vi *VerletIntegrator) Integrate(b body.Body, dt float64) {
	// Se il corpo è statico, non fare nulla
	if b.IsStatic() {
		return
	}

	// Ottieni la posizione corrente
	currentPosition := b.Position()

	// Verifica se esiste una posizione precedente per questo corpo
	previousPosition, exists := vi.previousPositions[b.ID()]

	if !exists {
		// Se non esiste una posizione precedente, usa l'integratore di Euler per il primo passo
		// x(t-dt) = x(t) - v(t)*dt + 0.5*a(t)*dt^2
		previousPosition = currentPosition.Sub(b.Velocity().Scale(dt)).Add(b.Acceleration().Scale(0.5 * dt * dt))
		vi.previousPositions[b.ID()] = previousPosition
	}

	// Calcola la nuova posizione usando l'algoritmo di Verlet
	// x(t+dt) = 2*x(t) - x(t-dt) + a(t)*dt^2
	newPosition := currentPosition.Scale(2).Sub(previousPosition).Add(b.Acceleration().Scale(dt * dt))

	// Calcola la nuova velocità
	// v(t+dt) = (x(t+dt) - x(t-dt)) / (2*dt)
	newVelocity := newPosition.Sub(previousPosition).Scale(1.0 / (2.0 * dt))

	// Aggiorna la posizione precedente per il prossimo passo
	vi.previousPositions[b.ID()] = currentPosition

	// Aggiorna la posizione e la velocità del corpo
	b.SetPosition(newPosition)
	b.SetVelocity(newVelocity)

	// Resetta l'accelerazione (verrà ricalcolata nel prossimo ciclo)
	b.SetAcceleration(vector.Zero3())
}

// IntegrateAll integra le equazioni del moto per tutti i corpi usando il metodo di Verlet
func (vi *VerletIntegrator) IntegrateAll(bodies []body.Body, dt float64) {
	for _, b := range bodies {
		vi.Integrate(b, dt)
	}
}

// RK4Integrator implementa l'integratore di Runge-Kutta di quarto ordine
type RK4Integrator struct{}

// NewRK4Integrator crea un nuovo integratore di Runge-Kutta di quarto ordine
func NewRK4Integrator() *RK4Integrator {
	return &RK4Integrator{}
}

// Integrate integra le equazioni del moto per un corpo usando il metodo di Runge-Kutta di quarto ordine
func (rk *RK4Integrator) Integrate(b body.Body, dt float64) {
	// Se il corpo è statico, non fare nulla
	if b.IsStatic() {
		return
	}

	// Stato iniziale
	x0 := b.Position()
	v0 := b.Velocity()
	a0 := b.Acceleration()

	// Primo passo (k1)
	k1v := a0.Scale(dt)
	k1x := v0.Scale(dt)

	// Secondo passo (k2)
	// Nota: In un'implementazione completa, dovremmo calcolare nuove posizioni e velocità
	// e ricalcolare le accelerazioni, ma qui usiamo a0 come approssimazione
	v1 := v0.Add(k1v.Scale(0.5))
	a1 := a0
	k2v := a1.Scale(dt)
	k2x := v1.Scale(dt)

	// Terzo passo (k3)
	// Nota: In un'implementazione completa, dovremmo calcolare nuove posizioni e velocità
	// e ricalcolare le accelerazioni, ma qui usiamo a0 come approssimazione
	v2 := v0.Add(k2v.Scale(0.5))
	a2 := a0
	k3v := a2.Scale(dt)
	k3x := v2.Scale(dt)

	// Quarto passo (k4)
	// Nota: In un'implementazione completa, dovremmo calcolare nuove posizioni e velocità
	// e ricalcolare le accelerazioni, ma qui usiamo a0 come approssimazione
	v3 := v0.Add(k3v)
	a3 := a0
	k4v := a3.Scale(dt)
	k4x := v3.Scale(dt)

	// Calcola la nuova posizione e velocità
	// x(t+dt) = x(t) + (k1x + 2*k2x + 2*k3x + k4x) / 6
	// v(t+dt) = v(t) + (k1v + 2*k2v + 2*k3v + k4v) / 6
	newPosition := x0.Add(k1x.Add(k2x.Scale(2)).Add(k3x.Scale(2)).Add(k4x).Scale(1.0 / 6.0))
	newVelocity := v0.Add(k1v.Add(k2v.Scale(2)).Add(k3v.Scale(2)).Add(k4v).Scale(1.0 / 6.0))

	// Aggiorna la posizione e la velocità del corpo
	b.SetPosition(newPosition)
	b.SetVelocity(newVelocity)

	// Resetta l'accelerazione (verrà ricalcolata nel prossimo ciclo)
	b.SetAcceleration(vector.Zero3())
}

// IntegrateAll integra le equazioni del moto per tutti i corpi usando il metodo di Runge-Kutta di quarto ordine
func (rk *RK4Integrator) IntegrateAll(bodies []body.Body, dt float64) {
	for _, b := range bodies {
		rk.Integrate(b, dt)
	}
}
