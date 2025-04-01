// Package main fornisce un esempio di utilizzo del motore fisico
package main

import (
	"fmt"
	"math"
	"time"

	"github.com/alexanderi96/go-space-engine/core/units"
	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/physics/body"
	"github.com/alexanderi96/go-space-engine/physics/force"
	"github.com/alexanderi96/go-space-engine/physics/material"
	"github.com/alexanderi96/go-space-engine/simulation/config"
	"github.com/alexanderi96/go-space-engine/simulation/world"
)

func main() {
	// Crea una configurazione per la simulazione
	cfg := config.NewSimulationBuilder().
		WithTimeStep(0.01).
		WithMaxBodies(100).
		WithGravity(true).
		WithCollisions(true).
		WithBoundaryCollisions(true).
		WithWorldBounds(
			vector.NewVector3(-10, -10, -10),
			vector.NewVector3(10, 10, 10),
		).
		WithRestitution(0.7).
		WithIntegratorType("verlet").
		Build()

	// Crea il mondo della simulazione
	w := world.NewPhysicalWorld(cfg.GetWorldBounds())

	// Aggiungi la forza gravitazionale
	gravityForce := force.NewGravitationalForce()
	w.AddForce(gravityForce)

	// Crea alcuni corpi
	createBodies(w)

	// Esegui la simulazione
	runSimulation(w, cfg)
}

// createBodies crea alcuni corpi nel mondo
func createBodies(w world.World) {
	// Crea un corpo centrale massivo (come un sole)
	sun := body.NewRigidBody(
		units.NewQuantity(1.0e6, units.Kilogram),
		units.NewQuantity(1.0, units.Meter),
		vector.NewVector3(0, 0, 0),
		vector.NewVector3(0, 0, 0),
		material.Iron,
	)
	sun.SetStatic(true) // Il sole è statico (non si muove)
	w.AddBody(sun)

	// Crea alcuni pianeti in orbita
	createPlanet(w, 3.0, 0.3, 0.5, vector.NewVector3(0, 1, 0), material.Rock)
	createPlanet(w, 5.0, 0.4, 0.3, vector.NewVector3(0, 1, 0), material.Ice)
	createPlanet(w, 7.0, 0.5, 0.2, vector.NewVector3(0, 1, 0), material.Copper)

	// Crea alcune lune
	createMoon(w, 3.0, 0.3, 0.7, 0.1, vector.NewVector3(0, 0, 1), material.Ice)
	createMoon(w, 5.0, 0.4, 0.9, 0.15, vector.NewVector3(0, 0, 1), material.Rock)

	// Crea alcuni asteroidi casuali
	for i := 0; i < 10; i++ {
		angle := float64(i) * 0.628 // 2*pi/10
		distance := 9.0

		x := distance * math.Cos(angle)
		z := distance * math.Sin(angle)

		asteroid := body.NewRigidBody(
			units.NewQuantity(100.0, units.Kilogram),
			units.NewQuantity(0.1, units.Meter),
			vector.NewVector3(x, 0, z),
			vector.NewVector3(-z*0.3, 0, x*0.3), // Velocità tangenziale
			material.Rock,
		)
		w.AddBody(asteroid)
	}
}

// createPlanet crea un pianeta in orbita
func createPlanet(w world.World, distance, radius, speed float64, orbitPlane vector.Vector3, mat material.Material) body.Body {
	// Calcola la posizione iniziale
	position := vector.NewVector3(distance, 0, 0)

	// Calcola la velocità orbitale (perpendicolare alla posizione)
	velocity := orbitPlane.Cross(position).Normalize().Scale(speed)

	// Crea il pianeta
	planet := body.NewRigidBody(
		units.NewQuantity(1000.0, units.Kilogram),
		units.NewQuantity(radius, units.Meter),
		position,
		velocity,
		mat,
	)

	// Aggiungi il pianeta al mondo
	w.AddBody(planet)

	return planet
}

// createMoon crea una luna in orbita attorno a un pianeta
func createMoon(w world.World, planetDistance, planetRadius, moonDistance, moonRadius float64, orbitPlane vector.Vector3, mat material.Material) body.Body {
	// Calcola la posizione del pianeta
	planetPosition := vector.NewVector3(planetDistance, 0, 0)

	// Calcola la posizione della luna rispetto al pianeta
	moonRelativePosition := vector.NewVector3(moonDistance, 0, 0)

	// Calcola la posizione assoluta della luna
	moonPosition := planetPosition.Add(moonRelativePosition)

	// Calcola la velocità orbitale del pianeta
	planetVelocity := orbitPlane.Cross(planetPosition).Normalize().Scale(math.Sqrt(1.0 / planetDistance))

	// Calcola la velocità orbitale della luna rispetto al pianeta
	moonRelativeVelocity := orbitPlane.Cross(moonRelativePosition).Normalize().Scale(math.Sqrt(10.0 / moonDistance))

	// Calcola la velocità assoluta della luna
	moonVelocity := planetVelocity.Add(moonRelativeVelocity)

	// Crea la luna
	moon := body.NewRigidBody(
		units.NewQuantity(100.0, units.Kilogram),
		units.NewQuantity(moonRadius, units.Meter),
		moonPosition,
		moonVelocity,
		mat,
	)

	// Aggiungi la luna al mondo
	w.AddBody(moon)

	return moon
}

// runSimulation esegue la simulazione
func runSimulation(w world.World, cfg *config.Config) {
	// Parametri della simulazione
	timeStep := cfg.TimeStep
	totalTime := 100.0   // Tempo totale della simulazione (secondi)
	printInterval := 1.0 // Intervallo di stampa (secondi)

	// Variabili per il timing
	lastPrintTime := 0.0
	startTime := time.Now()

	// Loop di simulazione
	for t := 0.0; t < totalTime; t += timeStep {
		// Esegui un passo della simulazione
		w.Step(timeStep)

		// Stampa lo stato della simulazione a intervalli regolari
		if t-lastPrintTime >= printInterval {
			// Calcola il tempo reale trascorso
			elapsedTime := time.Since(startTime).Seconds()

			// Stampa lo stato della simulazione
			fmt.Printf("Simulation time: %.2f s, Real time: %.2f s, Bodies: %d\n",
				t, elapsedTime, w.GetBodyCount())

			// Aggiorna il tempo dell'ultima stampa
			lastPrintTime = t
		}
	}

	// Stampa il tempo totale di esecuzione
	fmt.Printf("Simulation completed in %.2f seconds\n", time.Since(startTime).Seconds())
}
