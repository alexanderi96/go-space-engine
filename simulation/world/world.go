// Package world provides the implementation of the simulation world
package world

import (
	"runtime"
	"sync"

	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/physics/body"
	"github.com/alexanderi96/go-space-engine/physics/collision"
	"github.com/alexanderi96/go-space-engine/physics/force"
	"github.com/alexanderi96/go-space-engine/physics/integrator"
	"github.com/alexanderi96/go-space-engine/physics/space"
	"github.com/google/uuid"
)

// World represents the simulation world
type World interface {
	// AddBody adds a body to the world
	AddBody(b body.Body)
	// RemoveBody removes a body from the world
	RemoveBody(id uuid.UUID)
	// GetBody returns a body from the world
	GetBody(id uuid.UUID) body.Body
	// GetBodies returns all bodies in the world
	GetBodies() []body.Body
	// GetBodyCount returns the number of bodies in the world
	GetBodyCount() int

	// AddForce adds a force to the world
	AddForce(f force.Force)
	// RemoveForce removes a force from the world
	RemoveForce(f force.Force)
	// GetForces returns all forces in the world
	GetForces() []force.Force

	// SetIntegrator sets the numerical integrator
	SetIntegrator(i integrator.Integrator)
	// GetIntegrator returns the numerical integrator
	GetIntegrator() integrator.Integrator

	// SetCollider sets the collision detector
	SetCollider(c collision.Collider)
	// GetCollider returns the collision detector
	GetCollider() collision.Collider

	// SetCollisionResolver sets the collision resolver
	SetCollisionResolver(r collision.CollisionResolver)
	// GetCollisionResolver returns the collision resolver
	GetCollisionResolver() collision.CollisionResolver

	// SetSpatialStructure sets the spatial structure
	SetSpatialStructure(s space.SpatialStructure)
	// GetSpatialStructure returns the spatial structure
	GetSpatialStructure() space.SpatialStructure

	// SetBounds sets the world boundaries
	SetBounds(bounds *space.AABB)
	// GetBounds returns the world boundaries
	GetBounds() *space.AABB

	// Step advances the simulation by one time step
	Step(dt float64)

	// Clear removes all bodies and forces from the world
	Clear()
}

// WorkerPool represents a worker pool for parallel computation
type WorkerPool struct {
	numWorkers int
	tasks      chan func()
	wg         sync.WaitGroup
	mutex      sync.Mutex
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(numWorkers int) *WorkerPool {
	pool := &WorkerPool{
		numWorkers: numWorkers,
		tasks:      make(chan func(), numWorkers*10), // Buffer for tasks
	}

	// Start the workers
	for i := 0; i < numWorkers; i++ {
		go pool.worker()
	}

	return pool
}

// worker executes tasks from the channel
func (wp *WorkerPool) worker() {
	for task := range wp.tasks {
		task()
		wp.wg.Done()
	}
}

// Submit sends a task to the pool
func (wp *WorkerPool) Submit(task func()) {
	wp.mutex.Lock()
	wp.wg.Add(1)
	wp.mutex.Unlock()
	wp.tasks <- task
}

// Wait waits for all tasks to be completed
func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}

// PhysicalWorld implements the World interface
type PhysicalWorld struct {
	bodies            map[uuid.UUID]body.Body
	forces            []force.Force
	integrator        integrator.Integrator
	collider          collision.Collider
	collisionResolver collision.CollisionResolver
	spatialStructure  space.SpatialStructure
	bounds            *space.AABB
	workerPool        *WorkerPool
}

// NewPhysicalWorld creates a new physical world
func NewPhysicalWorld(bounds *space.AABB) *PhysicalWorld {
	// Create a spatial structure (octree) with predefined boundaries
	spatialStructure := space.NewOctree(bounds, 10, 8)

	// Create a worker pool that adapts to the number of available cores
	numCPU := runtime.NumCPU()
	workerPool := NewWorkerPool(numCPU)

	return &PhysicalWorld{
		bodies:            make(map[uuid.UUID]body.Body),
		forces:            make([]force.Force, 0),
		integrator:        integrator.NewVerletIntegrator(),
		collider:          collision.NewSphereCollider(),
		collisionResolver: collision.NewImpulseResolver(0.5),
		spatialStructure:  spatialStructure,
		bounds:            bounds,
		workerPool:        workerPool,
	}
}

// AddBody adds a body to the world
func (w *PhysicalWorld) AddBody(b body.Body) {
	w.bodies[b.ID()] = b
	w.spatialStructure.Insert(b)
}

// RemoveBody removes a body from the world
func (w *PhysicalWorld) RemoveBody(id uuid.UUID) {
	if b, exists := w.bodies[id]; exists {
		w.spatialStructure.Remove(b)
		delete(w.bodies, id)
	}
}

// GetBody returns a body from the world
func (w *PhysicalWorld) GetBody(id uuid.UUID) body.Body {
	return w.bodies[id]
}

// GetBodies returns all bodies in the world
func (w *PhysicalWorld) GetBodies() []body.Body {
	bodies := make([]body.Body, 0, len(w.bodies))
	for _, b := range w.bodies {
		bodies = append(bodies, b)
	}
	return bodies
}

// GetBodyCount returns the number of bodies in the world
func (w *PhysicalWorld) GetBodyCount() int {
	return len(w.bodies)
}

// AddForce adds a force to the world
func (w *PhysicalWorld) AddForce(f force.Force) {
	w.forces = append(w.forces, f)
}

// RemoveForce removes a force from the world
func (w *PhysicalWorld) RemoveForce(f force.Force) {
	for i, force := range w.forces {
		if force == f {
			// Remove the force by swapping it with the last one and truncating the slice
			lastIndex := len(w.forces) - 1
			w.forces[i] = w.forces[lastIndex]
			w.forces = w.forces[:lastIndex]
			break
		}
	}
}

// GetForces returns all forces in the world
func (w *PhysicalWorld) GetForces() []force.Force {
	return w.forces
}

// SetIntegrator sets the numerical integrator
func (w *PhysicalWorld) SetIntegrator(i integrator.Integrator) {
	w.integrator = i
}

// GetIntegrator returns the numerical integrator
func (w *PhysicalWorld) GetIntegrator() integrator.Integrator {
	return w.integrator
}

// SetCollider sets the collision detector
func (w *PhysicalWorld) SetCollider(c collision.Collider) {
	w.collider = c
}

// GetCollider returns the collision detector
func (w *PhysicalWorld) GetCollider() collision.Collider {
	return w.collider
}

// SetCollisionResolver sets the collision resolver
func (w *PhysicalWorld) SetCollisionResolver(r collision.CollisionResolver) {
	w.collisionResolver = r
}

// GetCollisionResolver returns the collision resolver
func (w *PhysicalWorld) GetCollisionResolver() collision.CollisionResolver {
	return w.collisionResolver
}

// SetSpatialStructure sets the spatial structure
func (w *PhysicalWorld) SetSpatialStructure(s space.SpatialStructure) {
	// Transfer all bodies from the old structure to the new one
	bodies := w.GetBodies()
	w.spatialStructure = s
	for _, b := range bodies {
		s.Insert(b)
	}
}

// GetSpatialStructure returns the spatial structure
func (w *PhysicalWorld) GetSpatialStructure() space.SpatialStructure {
	return w.spatialStructure
}

// SetBounds sets the world boundaries
func (w *PhysicalWorld) SetBounds(bounds *space.AABB) {
	w.bounds = bounds
}

// GetBounds returns the world boundaries
func (w *PhysicalWorld) GetBounds() *space.AABB {
	return w.bounds
}

// Step advances the simulation by one time step
func (w *PhysicalWorld) Step(dt float64) {
	// Apply forces
	w.applyForces()

	// Detect and resolve collisions
	w.handleCollisions()

	// Integrate the equations of motion in parallel
	bodies := w.GetBodies()
	w.integrator.IntegrateAll(bodies, dt, w.workerPool)

	// Update the spatial structure
	w.updateSpatialStructure()
}

// Clear removes all bodies and forces from the world
func (w *PhysicalWorld) Clear() {
	w.bodies = make(map[uuid.UUID]body.Body)
	w.forces = make([]force.Force, 0)
	w.spatialStructure.Clear()
}

// applyForces applies all forces to all bodies
func (w *PhysicalWorld) applyForces() {
	bodies := w.GetBodies()

	// Check if there are gravitational forces that can use the octree
	var gravityForce *force.GravitationalForce
	for _, f := range w.forces {
		if gf, ok := f.(*force.GravitationalForce); ok {
			gravityForce = gf
			break
		}
	}

	// Apply global forces to all bodies in parallel
	for _, f := range w.forces {
		if f.IsGlobal() {
			// If it's a gravitational force and we have an octree, use the Barnes-Hut algorithm
			if gravityForce != nil && f == gravityForce {
				octree, ok := w.spatialStructure.(*space.Octree)
				if ok {
					// Use the octree to calculate gravity in parallel
					for _, b := range bodies {
						b := b // Capture the variable for the goroutine
						w.workerPool.Submit(func() {
							// Calculate the gravitational force using the Barnes-Hut algorithm
							force := octree.CalculateGravity(b, gravityForce.GetTheta())
							b.ApplyForce(force)
						})
					}
					w.workerPool.Wait()
					continue
				}
			}

			// For other global forces, apply normally in parallel
			for _, b := range bodies {
				b := b // Capture the variable for the goroutine
				f := f // Capture the variable for the goroutine
				w.workerPool.Submit(func() {
					force := f.Apply(b)
					b.ApplyForce(force)
				})
			}
			w.workerPool.Wait()
		}
	}

	// Apply forces between pairs of bodies in parallel
	for i := 0; i < len(bodies); i++ {
		for j := i + 1; j < len(bodies); j++ {
			i, j := i, j // Capture the variables for the goroutine
			w.workerPool.Submit(func() {
				for _, f := range w.forces {
					if !f.IsGlobal() {
						forceA, forceB := f.ApplyBetween(bodies[i], bodies[j])
						bodies[i].ApplyForce(forceA)
						bodies[j].ApplyForce(forceB)
					}
				}
			})
		}
	}
	w.workerPool.Wait()
}

// handleCollisions detects and resolves collisions
func (w *PhysicalWorld) handleCollisions() {
	bodies := w.GetBodies()

	// Detect and resolve collisions between pairs of bodies in parallel
	for i := 0; i < len(bodies); i++ {
		i := i // Capture the variable for the goroutine
		w.workerPool.Submit(func() {
			// Use the spatial structure to find potential collisions
			radius := bodies[i].Radius().Value()
			nearbyBodies := w.spatialStructure.QuerySphere(bodies[i].Position(), radius*2)

			for _, b := range nearbyBodies {
				// Avoid checking collision with itself
				if b.ID() == bodies[i].ID() {
					continue
				}

				// Detect the collision
				info := w.collider.CheckCollision(bodies[i], b)

				// Resolve the collision
				if info.HasCollided {
					w.collisionResolver.ResolveCollision(info)
				}
			}

			// Also check collisions with world boundaries
			w.handleBoundaryCollisions(bodies[i])
		})
	}
	w.workerPool.Wait()
}

// handleBoundaryCollisions handles collisions with world boundaries
func (w *PhysicalWorld) handleBoundaryCollisions(b body.Body) {
	// If the body is static, do nothing
	if b.IsStatic() {
		return
	}

	// Get the body data only once to reduce method calls
	position := b.Position()
	velocity := b.Velocity()
	radius := b.Radius().Value()
	bounds := w.bounds
	elasticity := b.Material().Elasticity()

	// Flags to track if position or velocity have been modified
	positionChanged := false
	velocityChanged := false
	newPosition := position
	newVelocity := velocity

	// Collision with the lower X boundary
	if position.X()-radius < bounds.Min.X() {
		// Correct the position
		newPosition = vector.NewVector3(bounds.Min.X()+radius, position.Y(), position.Z())
		positionChanged = true

		// Invert the X velocity with damping
		newVelocity = vector.NewVector3(-velocity.X()*elasticity, velocity.Y(), velocity.Z())
		velocityChanged = true
	}

	// Collision with the upper X boundary
	if position.X()+radius > bounds.Max.X() {
		// Correct the position
		newPosition = vector.NewVector3(bounds.Max.X()-radius, position.Y(), position.Z())
		positionChanged = true

		// Invert the X velocity with damping
		newVelocity = vector.NewVector3(-velocity.X()*elasticity, velocity.Y(), velocity.Z())
		velocityChanged = true
	}

	// Collision with the lower Y boundary
	if position.Y()-radius < bounds.Min.Y() {
		// Correct the position
		newPosition = vector.NewVector3(newPosition.X(), bounds.Min.Y()+radius, position.Z())
		positionChanged = true

		// Invert the Y velocity with damping
		newVelocity = vector.NewVector3(newVelocity.X(), -velocity.Y()*elasticity, velocity.Z())
		velocityChanged = true
	}

	// Collision with the upper Y boundary
	if position.Y()+radius > bounds.Max.Y() {
		// Correct the position
		newPosition = vector.NewVector3(newPosition.X(), bounds.Max.Y()-radius, position.Z())
		positionChanged = true

		// Invert the Y velocity with damping
		newVelocity = vector.NewVector3(newVelocity.X(), -velocity.Y()*elasticity, velocity.Z())
		velocityChanged = true
	}

	// Collision with the lower Z boundary
	if position.Z()-radius < bounds.Min.Z() {
		// Correct the position
		newPosition = vector.NewVector3(newPosition.X(), newPosition.Y(), bounds.Min.Z()+radius)
		positionChanged = true

		// Invert the Z velocity with damping
		newVelocity = vector.NewVector3(newVelocity.X(), newVelocity.Y(), -velocity.Z()*elasticity)
		velocityChanged = true
	}

	// Collision with the upper Z boundary
	if position.Z()+radius > bounds.Max.Z() {
		// Correct the position
		newPosition = vector.NewVector3(newPosition.X(), newPosition.Y(), bounds.Max.Z()-radius)
		positionChanged = true

		// Invert the Z velocity with damping
		newVelocity = vector.NewVector3(newVelocity.X(), newVelocity.Y(), -velocity.Z()*elasticity)
		velocityChanged = true
	}

	// Update position and velocity only if necessary
	if positionChanged {
		b.SetPosition(newPosition)
	}
	if velocityChanged {
		b.SetVelocity(newVelocity)
	}
}

// updateSpatialStructure updates the spatial structure
func (w *PhysicalWorld) updateSpatialStructure() {
	// Update the spatial structure in parallel
	w.spatialStructure.UpdateAll(w.GetBodies(), w.workerPool)
}
