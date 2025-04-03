// Package integrator provides numerical integrators for the equations of motion
package integrator

import (
	"sync"

	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/physics/body"
	"github.com/google/uuid"
)

// TaskSubmitter represents an interface for submitting tasks to be executed in parallel
type TaskSubmitter interface {
	// Submit submits a task to be executed
	Submit(task func())
	// Wait waits for all tasks to be completed
	Wait()
}

// Integrator represents a numerical integrator for the equations of motion
type Integrator interface {
	// Integrate integrates the equations of motion for a body
	Integrate(b body.Body, dt float64)
	// IntegrateAll integrates the equations of motion for all bodies
	IntegrateAll(bodies []body.Body, dt float64, taskSubmitter TaskSubmitter)
}

// EulerIntegrator implements the Euler integrator
type EulerIntegrator struct{}

// NewEulerIntegrator creates a new Euler integrator
func NewEulerIntegrator() *EulerIntegrator {
	return &EulerIntegrator{}
}

// Integrate integrates the equations of motion for a body using the Euler method
func (ei *EulerIntegrator) Integrate(b body.Body, dt float64) {
	// If the body is static, do nothing
	if b.IsStatic() {
		return
	}

	// Update position: x(t+dt) = x(t) + v(t)*dt
	newPosition := b.Position().Add(b.Velocity().Scale(dt))
	b.SetPosition(newPosition)

	// Update velocity: v(t+dt) = v(t) + a(t)*dt
	newVelocity := b.Velocity().Add(b.Acceleration().Scale(dt))
	b.SetVelocity(newVelocity)

	// Reset acceleration (will be recalculated in the next cycle)
	b.SetAcceleration(vector.Zero3())
}

// IntegrateAll integrates the equations of motion for all bodies using the Euler method
func (ei *EulerIntegrator) IntegrateAll(bodies []body.Body, dt float64, taskSubmitter TaskSubmitter) {
	for _, b := range bodies {
		b := b // Capture the variable for the goroutine
		taskSubmitter.Submit(func() {
			ei.Integrate(b, dt)
		})
	}
	taskSubmitter.Wait()
}

// VerletIntegrator implements the Verlet integrator
type VerletIntegrator struct {
	// Map that stores the previous positions of bodies
	previousPositions map[uuid.UUID]vector.Vector3
	// Mutex to protect access to the map
	mutex sync.RWMutex
}

// NewVerletIntegrator creates a new Verlet integrator
func NewVerletIntegrator() *VerletIntegrator {
	return &VerletIntegrator{
		previousPositions: make(map[uuid.UUID]vector.Vector3),
		mutex:             sync.RWMutex{},
	}
}

// Integrate integrates the equations of motion for a body using the Verlet method
func (vi *VerletIntegrator) Integrate(b body.Body, dt float64) {
	// If the body is static, do nothing
	if b.IsStatic() {
		return
	}

	// Get the current position
	currentPosition := b.Position()

	// Check if there is a previous position for this body
	vi.mutex.RLock()
	previousPosition, exists := vi.previousPositions[b.ID()]
	vi.mutex.RUnlock()

	if !exists {
		// If there is no previous position, use the Euler integrator for the first step
		// x(t-dt) = x(t) - v(t)*dt + 0.5*a(t)*dt^2
		previousPosition = currentPosition.Sub(b.Velocity().Scale(dt)).Add(b.Acceleration().Scale(0.5 * dt * dt))
		vi.mutex.Lock()
		vi.previousPositions[b.ID()] = previousPosition
		vi.mutex.Unlock()
	}

	// Calculate the new position using the Verlet algorithm
	// x(t+dt) = 2*x(t) - x(t-dt) + a(t)*dt^2
	newPosition := currentPosition.Scale(2).Sub(previousPosition).Add(b.Acceleration().Scale(dt * dt))

	// Calculate the new velocity
	// v(t+dt) = (x(t+dt) - x(t-dt)) / (2*dt)
	newVelocity := newPosition.Sub(previousPosition).Scale(1.0 / (2.0 * dt))

	// Update the previous position for the next step
	vi.mutex.Lock()
	vi.previousPositions[b.ID()] = currentPosition
	vi.mutex.Unlock()

	// Update the position and velocity of the body
	b.SetPosition(newPosition)
	b.SetVelocity(newVelocity)

	// Reset acceleration (will be recalculated in the next cycle)
	b.SetAcceleration(vector.Zero3())
}

// IntegrateAll integrates the equations of motion for all bodies in parallel using the Verlet method
func (vi *VerletIntegrator) IntegrateAll(bodies []body.Body, dt float64, taskSubmitter TaskSubmitter) {
	for _, b := range bodies {
		b := b // Capture the variable for the goroutine
		taskSubmitter.Submit(func() {
			vi.Integrate(b, dt)
		})
	}
	taskSubmitter.Wait()
}

// RK4Integrator implements the fourth-order Runge-Kutta integrator
type RK4Integrator struct{}

// NewRK4Integrator creates a new fourth-order Runge-Kutta integrator
func NewRK4Integrator() *RK4Integrator {
	return &RK4Integrator{}
}

// Integrate integrates the equations of motion for a body using the fourth-order Runge-Kutta method
func (rk *RK4Integrator) Integrate(b body.Body, dt float64) {
	// If the body is static, do nothing
	if b.IsStatic() {
		return
	}

	// Initial state
	x0 := b.Position()
	v0 := b.Velocity()
	a0 := b.Acceleration()

	// First step (k1)
	k1v := a0.Scale(dt)
	k1x := v0.Scale(dt)

	// Second step (k2)
	// Note: In a complete implementation, we should calculate new positions and velocities
	// and recalculate accelerations, but here we use a0 as an approximation
	v1 := v0.Add(k1v.Scale(0.5))
	a1 := a0
	k2v := a1.Scale(dt)
	k2x := v1.Scale(dt)

	// Third step (k3)
	// Note: In a complete implementation, we should calculate new positions and velocities
	// and recalculate accelerations, but here we use a0 as an approximation
	v2 := v0.Add(k2v.Scale(0.5))
	a2 := a0
	k3v := a2.Scale(dt)
	k3x := v2.Scale(dt)

	// Fourth step (k4)
	// Note: In a complete implementation, we should calculate new positions and velocities
	// and recalculate accelerations, but here we use a0 as an approximation
	v3 := v0.Add(k3v)
	a3 := a0
	k4v := a3.Scale(dt)
	k4x := v3.Scale(dt)

	// Calculate the new position and velocity
	// x(t+dt) = x(t) + (k1x + 2*k2x + 2*k3x + k4x) / 6
	// v(t+dt) = v(t) + (k1v + 2*k2v + 2*k3v + k4v) / 6
	newPosition := x0.Add(k1x.Add(k2x.Scale(2)).Add(k3x.Scale(2)).Add(k4x).Scale(1.0 / 6.0))
	newVelocity := v0.Add(k1v.Add(k2v.Scale(2)).Add(k3v.Scale(2)).Add(k4v).Scale(1.0 / 6.0))

	// Update the position and velocity of the body
	b.SetPosition(newPosition)
	b.SetVelocity(newVelocity)

	// Reset acceleration (will be recalculated in the next cycle)
	b.SetAcceleration(vector.Zero3())
}

// IntegrateAll integrates the equations of motion for all bodies in parallel using the fourth-order Runge-Kutta method
func (rk *RK4Integrator) IntegrateAll(bodies []body.Body, dt float64, taskSubmitter TaskSubmitter) {
	for _, b := range bodies {
		b := b // Capture the variable for the goroutine
		taskSubmitter.Submit(func() {
			rk.Integrate(b, dt)
		})
	}
	taskSubmitter.Wait()
}
