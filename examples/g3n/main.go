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
	adapter.SetBackgroundColor(g3n.NewColor(0.0, 0.0, 0.05, 1.0)) // Sfondo blu molto scuro per lo spazio

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
	log.Println("Creazione del sistema solare")
	createSolarSystem(w)

	// Crea un campo di asteroidi
	log.Println("Creazione del campo di asteroidi")
	createAsteroidBelt(w, 200, 60.0, 80.0)
}

// createSolarSystem crea un sistema solare realistico
func createSolarSystem(w world.World) {
	log.Println("Creazione del sole")

	// Massa fissa del sole - valore elevato per garantire orbite stabili
	// In una simulazione, i rapporti relativi sono più importanti dei valori assoluti
	solarMass := 1.0e11 // Valore semplificato

	log.Printf("Massa del sole: %e kg", solarMass)

	sun := body.NewRigidBody(
		units.NewQuantity(solarMass, units.Kilogram),
		units.NewQuantity(5.0, units.Meter),                        // Raggio del sole (scalato)
		vector.NewVector3(0, 0, 0),                                 // Posizione al centro
		vector.NewVector3(0, 0, 0),                                 // Velocità zero
		createMaterial("Sun", 0.9, 0.5, [3]float64{1.0, 0.8, 0.0}), // Colore giallo
	)
	sun.SetStatic(true) // Il sole è statico (non si muove)
	w.AddBody(sun)
	log.Printf("Sole creato: ID=%v, Posizione=%v", sun.ID(), sun.Position())

	// Crea i pianeti
	log.Println("Creazione dei pianeti")

	// Definiamo le distanze dei pianeti
	// Semplicemente aumentando progressivamente
	distances := []float64{20, 30, 40, 50, 70, 90, 110, 130}

	// Nomi dei pianeti
	names := []string{"Mercury", "Venus", "Earth", "Mars", "Jupiter", "Saturn", "Uranus", "Neptune"}

	// Raggi dei pianeti (scalati)
	radii := []float64{0.8, 1.2, 1.3, 1.0, 2.5, 2.2, 1.8, 1.8}

	// Masse dei pianeti come frazioni della massa solare
	// I valori esatti non sono importanti, l'importante è che siano molto più piccoli del sole
	massFractions := []float64{1e-7, 2e-7, 2e-7, 1e-7, 1e-6, 9e-7, 4e-7, 5e-7}

	// Colori dei pianeti
	colors := [][3]float64{
		{0.7, 0.7, 0.7}, // Mercury
		{0.9, 0.7, 0.0}, // Venus
		{0.0, 0.3, 0.8}, // Earth
		{0.8, 0.3, 0.0}, // Mars
		{0.8, 0.6, 0.4}, // Jupiter
		{0.9, 0.8, 0.5}, // Saturn
		{0.5, 0.8, 0.9}, // Uranus
		{0.0, 0.0, 0.8}, // Neptune
	}

	// Crea ogni pianeta
	for i := 0; i < len(names); i++ {
		// Calcola la velocità orbitale usando la formula corretta: v = sqrt(G*M/r)
		// Dove G è la costante gravitazionale, M è la massa del sole, r è la distanza
		orbitSpeed := math.Sqrt(constants.G*solarMass/distances[i]) - 1

		// Crea il pianeta
		createPlanet(
			w,
			names[i],                   // Nome
			solarMass*massFractions[i], // Massa
			radii[i],                   // Raggio
			distances[i],               // Distanza
			orbitSpeed,                 // Velocità orbitale calcolata
			vector.NewVector3(0, 1, 0), // Piano dell'orbita
			colors[i],                  // Colore
		)
	}
}

// createPlanet crea un pianeta
func createPlanet(w world.World, name string, mass, radius, distance, speed float64, orbitPlane vector.Vector3, color [3]float64) body.Body {
	log.Printf("Creazione del pianeta %s: distanza=%f, raggio=%f, velocità=%f", name, distance, radius, speed)

	// Angolo casuale per la posizione iniziale (per distribuire i pianeti intorno al sole)
	angle := rand.Float64() * 2 * math.Pi

	// Calcola la posizione iniziale
	position := vector.NewVector3(
		distance*math.Cos(angle),
		0,
		distance*math.Sin(angle),
	)

	// Calcola la velocità orbitale (perpendicolare alla posizione)
	// Questa è la chiave per orbite stabili: la velocità deve essere perpendicolare al raggio
	velocity := vector.NewVector3(
		-speed*math.Sin(angle), // Componente x
		0,                      // Componente y (piano xy)
		speed*math.Cos(angle),  // Componente z
	)

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

	// Massa del sole (deve corrispondere a quella usata per i pianeti)
	solarMass := 1.0e8

	for i := 0; i < count; i++ {
		// Genera una posizione casuale nel campo di asteroidi
		distance := minDistance + rand.Float64()*(maxDistance-minDistance)
		angle := rand.Float64() * 2 * math.Pi

		x := distance * math.Cos(angle)
		z := distance * math.Sin(angle)
		y := (rand.Float64()*2 - 1) * 5 // Distribuzione verticale più ampia

		position := vector.NewVector3(x, y, z)

		// Calcola la velocità orbitale corretta: v = sqrt(G*M/r)
		baseSpeed := math.Sqrt(constants.G * solarMass / distance)

		// Aggiungi una piccola variazione casuale alla velocità
		speed := baseSpeed * (0.95 + rand.Float64()*0.1) // 95-105% della velocità base

		// La velocità deve essere perpendicolare al raggio per un'orbita circolare
		// Per un asteroide nell'asse y ≠ 0, dobbiamo calcolare il vettore perpendicolare correttamente
		radialDirection := position.Normalize()

		// Vettore "su" nell'asse y
		up := vector.NewVector3(0, 1, 0)

		// Ottieni il vettore perpendicolare facendo il prodotto vettoriale
		tangentialDirection := up.Cross(radialDirection).Normalize()

		// Se il risultato è quasi zero (asteroide quasi sull'asse y), usa un'altra direzione
		if tangentialDirection.Length() < 0.1 {
			tangentialDirection = vector.NewVector3(1, 0, 0).Cross(radialDirection).Normalize()
		}

		velocity := tangentialDirection.Scale(speed)

		// Crea l'asteroide con massa ridotta
		asteroid := body.NewRigidBody(
			units.NewQuantity(rand.Float64()*10+1, units.Kilogram), // Massa molto più piccola dei pianeti
			units.NewQuantity(rand.Float64()*0.3+0.1, units.Meter), // Dimensione più piccola
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
