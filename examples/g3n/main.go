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
		WithMaxBodies(1000). // Aumentato il numero massimo di corpi
		WithGravity(true).
		WithCollisions(true).
		WithBoundaryCollisions(true).
		WithWorldBounds(
			vector.NewVector3(-100, -100, -100),
			vector.NewVector3(100, 100, 100),
		).
		WithOctreeConfig(10, 8). // Configurazione ottimale per l'octree
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
	adapter.SetBackgroundColor(g3n.NewColor(0.0, 0.0, 0.1, 1.0)) // Sfondo blu scuro per lo spazio

	// Variabili per il timing
	lastUpdateTime := time.Now()
	simulationTime := 0.0

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
		simulationTime += dt

		// Renderizza il mondo
		adapter.RenderWorld(w)
	})

	log.Println("Esempio completato")
}

// createBodies crea alcuni corpi nel mondo
func createBodies(w world.World) {
	log.Println("Creazione del sistema solare")
	createSolarSystem(w)

	// Crea un campo di asteroidi
	log.Println("Creazione del campo di asteroidi")
	createAsteroidBelt(w, 200, 30.0, 50.0)
}

// createSolarSystem crea un sistema solare realistico
func createSolarSystem(w world.World) {
	log.Println("Creazione del sole")

	// Crea il sole
	sun := body.NewRigidBody(
		units.NewQuantity(1.989e12, units.Kilogram), // Massa del sole aumentata per permettere orbite stabili
		units.NewQuantity(2.0, units.Meter),         // Raggio del sole (scalato)
		vector.NewVector3(0, 0, 0),
		vector.NewVector3(0, 0, 0),
		createMaterial("Sun", 0.9, 0.5, [3]float64{1.0, 0.8, 0.0}), // Colore giallo
	)
	sun.SetStatic(true) // Il sole è statico (non si muove)
	w.AddBody(sun)
	log.Printf("Sole creato: ID=%v, Posizione=%v", sun.ID(), sun.Position())

	// Crea i pianeti
	log.Println("Creazione dei pianeti")

	// Calcola le velocità orbitali corrette per ogni pianeta
	// Utilizziamo la formula v = sqrt(G*M/r) per la velocità orbitale circolare
	solarMass := 1.989e12 // Massa del sole in kg
	G := constants.G      // Costante gravitazionale

	// Mercurio
	mercuryDist := 8.0
	mercurySpeed := math.Sqrt(G*solarMass/mercuryDist) * 0.1
	createPlanet(w, "Mercury", 0.33e3, 0.4, mercuryDist, mercurySpeed, vector.NewVector3(0, 1, 0), [3]float64{0.7, 0.7, 0.7})

	// Venere
	venusDist := 12.0
	venusSpeed := math.Sqrt(G*solarMass/venusDist) * 0.1
	createPlanet(w, "Venus", 4.87e3, 0.6, venusDist, venusSpeed, vector.NewVector3(0, 1, 0), [3]float64{0.9, 0.7, 0.0})

	// Terra
	earthDist := 16.0
	earthSpeed := math.Sqrt(G*solarMass/earthDist) * 0.1
	createPlanet(w, "Earth", 5.97e3, 0.6, earthDist, earthSpeed, vector.NewVector3(0, 1, 0), [3]float64{0.0, 0.3, 0.8})

	// Marte
	marsDist := 20.0
	marsSpeed := math.Sqrt(G*solarMass/marsDist) * 0.1
	createPlanet(w, "Mars", 0.642e3, 0.5, marsDist, marsSpeed, vector.NewVector3(0, 1, 0), [3]float64{0.8, 0.3, 0.0})

	// Giove
	jupiterDist := 28.0
	jupiterSpeed := math.Sqrt(G*solarMass/jupiterDist) * 0.1
	createPlanet(w, "Jupiter", 1898e3, 1.2, jupiterDist, jupiterSpeed, vector.NewVector3(0, 1, 0), [3]float64{0.8, 0.6, 0.4})

	// Saturno
	saturnDist := 36.0
	saturnSpeed := math.Sqrt(G*solarMass/saturnDist) * 0.1
	createPlanet(w, "Saturn", 568e3, 1.0, saturnDist, saturnSpeed, vector.NewVector3(0, 1, 0), [3]float64{0.9, 0.8, 0.5})

	// Urano
	uranusDist := 44.0
	uranusSpeed := math.Sqrt(G*solarMass/uranusDist) * 0.1
	createPlanet(w, "Uranus", 86.8e3, 0.8, uranusDist, uranusSpeed, vector.NewVector3(0, 1, 0), [3]float64{0.5, 0.8, 0.9})

	// Nettuno
	neptuneDist := 52.0
	neptuneSpeed := math.Sqrt(G*solarMass/neptuneDist) * 0.1
	createPlanet(w, "Neptune", 102e3, 0.8, neptuneDist, neptuneSpeed, vector.NewVector3(0, 1, 0), [3]float64{0.0, 0.0, 0.8})
}

// createPlanet crea un pianeta con parametri realistici
func createPlanet(w world.World, name string, mass, radius, distance, speed float64, orbitPlane vector.Vector3, color [3]float64) body.Body {
	log.Printf("Creazione del pianeta %s: distanza=%f, raggio=%f, velocità=%f", name, distance, radius, speed)

	// Calcola la posizione iniziale
	position := vector.NewVector3(distance, 0, 0)

	// Calcola la velocità orbitale (perpendicolare alla posizione)
	velocity := orbitPlane.Cross(position).Normalize().Scale(speed)

	// Crea il pianeta
	planet := body.NewRigidBody(
		units.NewQuantity(mass, units.Kilogram),
		units.NewQuantity(radius, units.Meter),
		position,
		velocity,
		createMaterial(name, 0.7, 0.5, color),
	)

	// Aggiungi il pianeta al mondo
	w.AddBody(planet)
	log.Printf("Pianeta %s aggiunto: ID=%v, Posizione=%v, Velocità=%v", name, planet.ID(), planet.Position(), planet.Velocity())

	return planet
}

// createAsteroidBelt crea un campo di asteroidi
func createAsteroidBelt(w world.World, count int, minDistance, maxDistance float64) {
	log.Printf("Creazione di %d asteroidi", count)

	for i := 0; i < count; i++ {
		// Genera una posizione casuale nel campo di asteroidi
		distance := minDistance + rand.Float64()*(maxDistance-minDistance)
		angle := rand.Float64() * 2 * math.Pi

		x := distance * math.Cos(angle)
		z := distance * math.Sin(angle)
		y := (rand.Float64()*2 - 1) * 2 // Distribuzione verticale limitata

		position := vector.NewVector3(x, y, z)

		// Calcola la velocità orbitale corretta
		speed := math.Sqrt(constants.G*1.989e12/distance) * 0.1
		velocity := vector.NewVector3(-z, 0, x).Normalize().Scale(speed)

		// Crea l'asteroide
		asteroid := body.NewRigidBody(
			units.NewQuantity(rand.Float64()*100+10, units.Kilogram),
			units.NewQuantity(rand.Float64()*0.2+0.1, units.Meter),
			position,
			velocity,
			physMaterial.Rock,
		)

		// Aggiungi l'asteroide al mondo
		w.AddBody(asteroid)
	}

	log.Printf("Campo di asteroidi creato")
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
