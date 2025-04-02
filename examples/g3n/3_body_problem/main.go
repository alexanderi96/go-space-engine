// Package main fornisce un esempio di utilizzo di G3N con il motore fisico tramite adapter diretto
package main

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/alexanderi96/go-space-engine/core/constants"
	"github.com/alexanderi96/go-space-engine/core/units"
	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/physics/body"
	"github.com/alexanderi96/go-space-engine/physics/force"
	physMaterial "github.com/alexanderi96/go-space-engine/physics/material"
	"github.com/alexanderi96/go-space-engine/render/g3n"
	"github.com/alexanderi96/go-space-engine/simulation/config"
	"github.com/alexanderi96/go-space-engine/simulation/world"
)

func main() {
	log.Println("Inizializzazione dell'esempio G3N Physics con Adapter Diretto")

	// Inizializza il generatore di numeri casuali
	rand.Seed(time.Now().UnixNano())

	// Crea la configurazione della simulazione
	cfg := config.NewSimulationBuilder().
		WithTimeStep(0.01).
		WithMaxBodies(1000).
		WithGravity(true).
		WithCollisions(true).
		WithBoundaryCollisions(false). // Disabilitiamo le collisioni con i bordi
		WithWorldBounds(
			vector.NewVector3(-500, -500, -500),
			vector.NewVector3(500, 500, 500),
		).
		WithOctreeConfig(10, 8).
		Build()

	// Crea il mondo della simulazione
	w := world.NewPhysicalWorld(cfg.GetWorldBounds())

	// Aggiungi la forza gravitazionale
	gravityForce := force.NewGravitationalForce()
	gravityForce.SetTheta(0.5) // Imposta il valore di theta per l'algoritmo Barnes-Hut
	w.AddForce(gravityForce)

	// Crea alcuni corpi
	createBodies(w)

	// Crea l'adapter G3N diretto
	adapter := g3n.NewG3NAdapter()

	// Configura l'adapter
	adapter.SetBackgroundColor(g3n.NewColor(0, 0, 0, 1.0)) // Sfondo blu molto scuro per lo spazio

	// Variabili per il timing
	lastUpdateTime := time.Now()

	// Avvia il loop di rendering
	adapter.Run(func(deltaTime time.Duration) {
		// Calcola il delta time
		currentTime := time.Now()
		dt := currentTime.Sub(lastUpdateTime).Seconds()
		lastUpdateTime = currentTime

		// Limita il delta time per evitare instabilità
		if dt > 0.1 {
			dt = 0.1
		}

		// Esegui un passo della simulazione
		w.Step(dt)

		// Renderizza il mondo
		adapter.RenderWorld(w)
	})

	log.Println("Esempio completato")
}

// createBodies crea alcuni corpi nel mondo
func createBodies(w world.World) {
	create3BodySystem(w)
}

// create3BodySystem crea un sistema stabile di tre corpi con masse uguali
func create3BodySystem(w world.World) {
	// Massa uguale per tutti i corpi
	mass := 1.0e14
	log.Printf("Massa di ciascun corpo: %e kg", mass)

	// Raggio dell'orbita
	radius := 10.0

	// Colori per i tre corpi
	colors := [][3]float64{
		{1.0, 0.3, 0.3}, // Rosso
		{0.3, 1.0, 0.3}, // Verde
		{0.3, 0.3, 1.0}, // Blu
	}

	// Nomi dei corpi
	names := []string{"Corpo1", "Corpo2", "Corpo3"}

	// Calcola la velocità orbitale necessaria per un'orbita stabile
	// Per un sistema a tre corpi con masse uguali in configurazione triangolare equilatera
	// La formula corretta è v = sqrt(G*M/r)
	orbitSpeed := math.Sqrt(constants.G * mass / radius)

	// Applica un fattore di scala per rallentare ulteriormente il movimento
	// e rendere la simulazione più visivamente piacevole
	orbitSpeed *= 0.1

	// Crea i tre corpi posizionati ai vertici di un triangolo equilatero
	for i := 0; i < 3; i++ {
		// Calcola l'angolo per questo corpo (120 gradi di distanza l'uno dall'altro)
		angle := float64(i) * (2.0 * math.Pi / 3.0)

		// Calcola la posizione (vertici di un triangolo equilatero)
		position := vector.NewVector3(
			radius*math.Cos(angle),
			0,
			radius*math.Sin(angle),
		)

		// Calcola la velocità (perpendicolare alla posizione per un'orbita circolare)
		velocity := vector.NewVector3(
			-orbitSpeed*math.Sin(angle),
			0,
			orbitSpeed*math.Cos(angle),
		)

		// Crea il corpo
		b := body.NewRigidBody(
			units.NewQuantity(mass, units.Kilogram),
			units.NewQuantity(2.0, units.Meter), // Raggio del corpo
			position,
			velocity,
			createMaterial(names[i], 0.9, 0.5, colors[i]),
		)

		w.AddBody(b)
		log.Printf("%s creato: ID=%v, Posizione=%v, Velocità=%v", names[i], b.ID(), b.Position(), b.Velocity())
	}
}

// createMaterial crea un materiale personalizzato
func createMaterial(name string, emissivity, elasticity float64, color [3]float64) physMaterial.Material {
	return physMaterial.NewBasicMaterial(
		name,
		units.NewQuantity(5000, units.Kilogram),
		units.NewQuantity(800, units.Joule),
		units.NewQuantity(1.5, units.Watt),
		emissivity,
		elasticity,
		color,
	)
}
