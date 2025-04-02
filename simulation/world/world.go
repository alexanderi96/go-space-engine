// Package world fornisce l'implementazione del mondo della simulazione
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

// World rappresenta il mondo della simulazione
type World interface {
	// AddBody aggiunge un corpo al mondo
	AddBody(b body.Body)
	// RemoveBody rimuove un corpo dal mondo
	RemoveBody(id uuid.UUID)
	// GetBody restituisce un corpo dal mondo
	GetBody(id uuid.UUID) body.Body
	// GetBodies restituisce tutti i corpi nel mondo
	GetBodies() []body.Body
	// GetBodyCount restituisce il numero di corpi nel mondo
	GetBodyCount() int

	// AddForce aggiunge una forza al mondo
	AddForce(f force.Force)
	// RemoveForce rimuove una forza dal mondo
	RemoveForce(f force.Force)
	// GetForces restituisce tutte le forze nel mondo
	GetForces() []force.Force

	// SetIntegrator imposta l'integratore numerico
	SetIntegrator(i integrator.Integrator)
	// GetIntegrator restituisce l'integratore numerico
	GetIntegrator() integrator.Integrator

	// SetCollider imposta il rilevatore di collisioni
	SetCollider(c collision.Collider)
	// GetCollider restituisce il rilevatore di collisioni
	GetCollider() collision.Collider

	// SetCollisionResolver imposta il risolutore di collisioni
	SetCollisionResolver(r collision.CollisionResolver)
	// GetCollisionResolver restituisce il risolutore di collisioni
	GetCollisionResolver() collision.CollisionResolver

	// SetSpatialStructure imposta la struttura spaziale
	SetSpatialStructure(s space.SpatialStructure)
	// GetSpatialStructure restituisce la struttura spaziale
	GetSpatialStructure() space.SpatialStructure

	// SetBounds imposta i limiti del mondo
	SetBounds(bounds *space.AABB)
	// GetBounds restituisce i limiti del mondo
	GetBounds() *space.AABB

	// Step avanza la simulazione di un passo temporale
	Step(dt float64)

	// Clear rimuove tutti i corpi e le forze dal mondo
	Clear()
}

// WorkerPool rappresenta un pool di worker per il calcolo parallelo
type WorkerPool struct {
	numWorkers int
	tasks      chan func()
	wg         sync.WaitGroup
}

// NewWorkerPool crea un nuovo pool di worker
func NewWorkerPool(numWorkers int) *WorkerPool {
	pool := &WorkerPool{
		numWorkers: numWorkers,
		tasks:      make(chan func(), numWorkers*10), // Buffer per le task
	}

	// Avvia i worker
	for i := 0; i < numWorkers; i++ {
		go pool.worker()
	}

	return pool
}

// worker esegue le task dal canale
func (wp *WorkerPool) worker() {
	for task := range wp.tasks {
		task()
		wp.wg.Done()
	}
}

// Submit invia una task al pool
func (wp *WorkerPool) Submit(task func()) {
	wp.wg.Add(1)
	wp.tasks <- task
}

// Wait attende che tutte le task siano completate
func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}

// PhysicalWorld implementa l'interfaccia World
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

// NewPhysicalWorld crea un nuovo mondo fisico
func NewPhysicalWorld(bounds *space.AABB) *PhysicalWorld {
	// Crea una struttura spaziale (octree) con limiti predefiniti
	spatialStructure := space.NewOctree(bounds, 10, 8)

	// Crea un pool di worker che si adatta al numero di core disponibili
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

// AddBody aggiunge un corpo al mondo
func (w *PhysicalWorld) AddBody(b body.Body) {
	w.bodies[b.ID()] = b
	w.spatialStructure.Insert(b)
}

// RemoveBody rimuove un corpo dal mondo
func (w *PhysicalWorld) RemoveBody(id uuid.UUID) {
	if b, exists := w.bodies[id]; exists {
		w.spatialStructure.Remove(b)
		delete(w.bodies, id)
	}
}

// GetBody restituisce un corpo dal mondo
func (w *PhysicalWorld) GetBody(id uuid.UUID) body.Body {
	return w.bodies[id]
}

// GetBodies restituisce tutti i corpi nel mondo
func (w *PhysicalWorld) GetBodies() []body.Body {
	bodies := make([]body.Body, 0, len(w.bodies))
	for _, b := range w.bodies {
		bodies = append(bodies, b)
	}
	return bodies
}

// GetBodyCount restituisce il numero di corpi nel mondo
func (w *PhysicalWorld) GetBodyCount() int {
	return len(w.bodies)
}

// AddForce aggiunge una forza al mondo
func (w *PhysicalWorld) AddForce(f force.Force) {
	w.forces = append(w.forces, f)
}

// RemoveForce rimuove una forza dal mondo
func (w *PhysicalWorld) RemoveForce(f force.Force) {
	for i, force := range w.forces {
		if force == f {
			// Rimuovi la forza scambiandola con l'ultima e troncando la slice
			lastIndex := len(w.forces) - 1
			w.forces[i] = w.forces[lastIndex]
			w.forces = w.forces[:lastIndex]
			break
		}
	}
}

// GetForces restituisce tutte le forze nel mondo
func (w *PhysicalWorld) GetForces() []force.Force {
	return w.forces
}

// SetIntegrator imposta l'integratore numerico
func (w *PhysicalWorld) SetIntegrator(i integrator.Integrator) {
	w.integrator = i
}

// GetIntegrator restituisce l'integratore numerico
func (w *PhysicalWorld) GetIntegrator() integrator.Integrator {
	return w.integrator
}

// SetCollider imposta il rilevatore di collisioni
func (w *PhysicalWorld) SetCollider(c collision.Collider) {
	w.collider = c
}

// GetCollider restituisce il rilevatore di collisioni
func (w *PhysicalWorld) GetCollider() collision.Collider {
	return w.collider
}

// SetCollisionResolver imposta il risolutore di collisioni
func (w *PhysicalWorld) SetCollisionResolver(r collision.CollisionResolver) {
	w.collisionResolver = r
}

// GetCollisionResolver restituisce il risolutore di collisioni
func (w *PhysicalWorld) GetCollisionResolver() collision.CollisionResolver {
	return w.collisionResolver
}

// SetSpatialStructure imposta la struttura spaziale
func (w *PhysicalWorld) SetSpatialStructure(s space.SpatialStructure) {
	// Trasferisci tutti i corpi dalla vecchia struttura alla nuova
	bodies := w.GetBodies()
	w.spatialStructure = s
	for _, b := range bodies {
		s.Insert(b)
	}
}

// GetSpatialStructure restituisce la struttura spaziale
func (w *PhysicalWorld) GetSpatialStructure() space.SpatialStructure {
	return w.spatialStructure
}

// SetBounds imposta i limiti del mondo
func (w *PhysicalWorld) SetBounds(bounds *space.AABB) {
	w.bounds = bounds
}

// GetBounds restituisce i limiti del mondo
func (w *PhysicalWorld) GetBounds() *space.AABB {
	return w.bounds
}

// Step avanza la simulazione di un passo temporale
func (w *PhysicalWorld) Step(dt float64) {
	// Applica le forze
	w.applyForces()

	// Rileva e risolvi le collisioni
	w.handleCollisions()

	// Integra le equazioni del moto
	w.integrator.IntegrateAll(w.GetBodies(), dt)

	// Aggiorna la struttura spaziale
	w.updateSpatialStructure()
}

// Clear rimuove tutti i corpi e le forze dal mondo
func (w *PhysicalWorld) Clear() {
	w.bodies = make(map[uuid.UUID]body.Body)
	w.forces = make([]force.Force, 0)
	w.spatialStructure.Clear()
}

// applyForces applica tutte le forze a tutti i corpi
func (w *PhysicalWorld) applyForces() {
	bodies := w.GetBodies()

	// Verifica se ci sono forze gravitazionali che possono utilizzare l'octree
	var gravityForce *force.GravitationalForce
	for _, f := range w.forces {
		if gf, ok := f.(*force.GravitationalForce); ok {
			gravityForce = gf
			break
		}
	}

	// Applica le forze globali a tutti i corpi in parallelo
	for _, f := range w.forces {
		if f.IsGlobal() {
			// Se è una forza gravitazionale e abbiamo un octree, usa l'algoritmo Barnes-Hut
			if gravityForce != nil && f == gravityForce {
				octree, ok := w.spatialStructure.(*space.Octree)
				if ok {
					// Usa l'octree per calcolare la gravità in parallelo
					for _, b := range bodies {
						b := b // Cattura la variabile per la goroutine
						w.workerPool.Submit(func() {
							// Calcola la forza gravitazionale usando l'algoritmo Barnes-Hut
							force := octree.CalculateGravity(b, gravityForce.GetTheta())
							b.ApplyForce(force)
						})
					}
					w.workerPool.Wait()
					continue
				}
			}

			// Per altre forze globali, applica normalmente in parallelo
			for _, b := range bodies {
				b := b // Cattura la variabile per la goroutine
				f := f // Cattura la variabile per la goroutine
				w.workerPool.Submit(func() {
					force := f.Apply(b)
					b.ApplyForce(force)
				})
			}
			w.workerPool.Wait()
		}
	}

	// Applica le forze tra coppie di corpi in parallelo
	for i := 0; i < len(bodies); i++ {
		for j := i + 1; j < len(bodies); j++ {
			i, j := i, j // Cattura le variabili per la goroutine
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

// handleCollisions rileva e risolve le collisioni
func (w *PhysicalWorld) handleCollisions() {
	bodies := w.GetBodies()

	// Rileva e risolvi le collisioni tra coppie di corpi in parallelo
	for i := 0; i < len(bodies); i++ {
		i := i // Cattura la variabile per la goroutine
		w.workerPool.Submit(func() {
			// Usa la struttura spaziale per trovare potenziali collisioni
			radius := bodies[i].Radius().Value()
			nearbyBodies := w.spatialStructure.QuerySphere(bodies[i].Position(), radius*2)

			for _, b := range nearbyBodies {
				// Evita di controllare la collisione con se stesso
				if b.ID() == bodies[i].ID() {
					continue
				}

				// Rileva la collisione
				info := w.collider.CheckCollision(bodies[i], b)

				// Risolvi la collisione
				if info.HasCollided {
					w.collisionResolver.ResolveCollision(info)
				}
			}

			// Controlla anche le collisioni con i limiti del mondo
			w.handleBoundaryCollisions(bodies[i])
		})
	}
	w.workerPool.Wait()
}

// handleBoundaryCollisions gestisce le collisioni con i limiti del mondo
func (w *PhysicalWorld) handleBoundaryCollisions(b body.Body) {
	// Se il corpo è statico, non fare nulla
	if b.IsStatic() {
		return
	}

	position := b.Position()
	velocity := b.Velocity()
	radius := b.Radius().Value()

	// Controlla la collisione con i limiti del mondo
	bounds := w.bounds

	// Collisione con il limite inferiore X
	if position.X()-radius < bounds.Min.X() {
		// Correggi la posizione
		b.SetPosition(vector.NewVector3(bounds.Min.X()+radius, position.Y(), position.Z()))

		// Inverti la velocità X con smorzamento
		elasticity := b.Material().Elasticity()
		b.SetVelocity(vector.NewVector3(-velocity.X()*elasticity, velocity.Y(), velocity.Z()))
	}

	// Collisione con il limite superiore X
	if position.X()+radius > bounds.Max.X() {
		// Correggi la posizione
		b.SetPosition(vector.NewVector3(bounds.Max.X()-radius, position.Y(), position.Z()))

		// Inverti la velocità X con smorzamento
		elasticity := b.Material().Elasticity()
		b.SetVelocity(vector.NewVector3(-velocity.X()*elasticity, velocity.Y(), velocity.Z()))
	}

	// Collisione con il limite inferiore Y
	if position.Y()-radius < bounds.Min.Y() {
		// Correggi la posizione
		b.SetPosition(vector.NewVector3(position.X(), bounds.Min.Y()+radius, position.Z()))

		// Inverti la velocità Y con smorzamento
		elasticity := b.Material().Elasticity()
		b.SetVelocity(vector.NewVector3(velocity.X(), -velocity.Y()*elasticity, velocity.Z()))
	}

	// Collisione con il limite superiore Y
	if position.Y()+radius > bounds.Max.Y() {
		// Correggi la posizione
		b.SetPosition(vector.NewVector3(position.X(), bounds.Max.Y()-radius, position.Z()))

		// Inverti la velocità Y con smorzamento
		elasticity := b.Material().Elasticity()
		b.SetVelocity(vector.NewVector3(velocity.X(), -velocity.Y()*elasticity, velocity.Z()))
	}

	// Collisione con il limite inferiore Z
	if position.Z()-radius < bounds.Min.Z() {
		// Correggi la posizione
		b.SetPosition(vector.NewVector3(position.X(), position.Y(), bounds.Min.Z()+radius))

		// Inverti la velocità Z con smorzamento
		elasticity := b.Material().Elasticity()
		b.SetVelocity(vector.NewVector3(velocity.X(), velocity.Y(), -velocity.Z()*elasticity))
	}

	// Collisione con il limite superiore Z
	if position.Z()+radius > bounds.Max.Z() {
		// Correggi la posizione
		b.SetPosition(vector.NewVector3(position.X(), position.Y(), bounds.Max.Z()-radius))

		// Inverti la velocità Z con smorzamento
		elasticity := b.Material().Elasticity()
		b.SetVelocity(vector.NewVector3(velocity.X(), velocity.Y(), -velocity.Z()*elasticity))
	}
}

// updateSpatialStructure aggiorna la struttura spaziale
func (w *PhysicalWorld) updateSpatialStructure() {
	for _, b := range w.bodies {
		w.spatialStructure.Update(b)
	}
}
