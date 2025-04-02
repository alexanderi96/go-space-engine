// Package main fornisce un esempio di posizionamento di corpi in formazioni geometriche prestabilite
package main

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/alexanderi96/go-space-engine/core/units"
	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/physics/body"
	physMaterial "github.com/alexanderi96/go-space-engine/physics/material"
	"github.com/alexanderi96/go-space-engine/render/g3n"
	"github.com/alexanderi96/go-space-engine/simulation/config"
	"github.com/alexanderi96/go-space-engine/simulation/world"
)

func main() {
	log.Println("Inizializzazione del posizionamento di corpi in formazioni geometriche")

	// Inizializza il generatore di numeri casuali
	rand.Seed(time.Now().UnixNano())

	// Crea la configurazione della simulazione
	cfg := config.NewSimulationBuilder().
		WithTimeStep(0.01).
		WithMaxBodies(5000).
		WithGravity(true).    // Disabilita la gravità per mantenere i corpi fermi nelle loro posizioni
		WithCollisions(true). // Disabilita le collisioni
		WithBoundaryCollisions(true).
		WithWorldBounds(
			vector.NewVector3(-500, -500, -500),
			vector.NewVector3(500, 500, 500),
		).
		WithOctreeConfig(10, 8).
		Build()

	// Crea il mondo della simulazione
	w := world.NewPhysicalWorld(cfg.GetWorldBounds())

	// Scegli una delle configurazioni seguenti:

	// Sfera di corpi
	createSphereFormation(w, 500, 100.0, 5.0)

	// Cubo di corpi
	//createCuboidFormation(w, 500, 80.0, 5.0)

	// Griglia 3D
	//createGrid3D(w, 10, 10, 10, 10.0)

	// Spirale
	//createSpiral(w, 300, 100.0, 5.0, 3)

	// Anello
	//createRingFormation(w, 200, 80.0, 5.0)

	// Crea l'adapter G3N diretto
	adapter := g3n.NewG3NAdapter()

	// Configura l'adapter
	adapter.SetBackgroundColor(g3n.NewColor(0.0, 0.0, 0.05, 1.0)) // Sfondo blu molto scuro per lo spazio

	// Avvia il loop di rendering
	adapter.Run(func(deltaTime time.Duration) {
		// Renderizza il mondo (senza step fisici per mantenere le posizioni)
		adapter.RenderWorld(w)
	})

	log.Println("Rendering completato")
}

// createSphereFormation crea corpi distribuiti su una sfera
func createSphereFormation(w world.World, count int, radius, minDistance float64) {
	log.Printf("Creazione di %d corpi in formazione sferica", count)

	positions := make([]vector.Vector3, 0, count)

	for i := 0; i < count; i++ {
		// Genera un punto sulla sfera con distribuzione uniforme
		phi := rand.Float64() * 2 * math.Pi
		costheta := rand.Float64()*2 - 1
		theta := math.Acos(costheta)

		x := radius * math.Sin(theta) * math.Cos(phi)
		y := radius * math.Sin(theta) * math.Sin(phi)
		z := radius * math.Cos(theta)

		position := vector.NewVector3(x, y, z)

		// Verifica distanza minima da altri corpi
		tooClose := false
		for _, pos := range positions {
			if position.Distance(pos) < minDistance {
				tooClose = true
				break
			}
		}

		// Se troppo vicino, riprova
		if tooClose {
			i--
			continue
		}

		positions = append(positions, position)

		// Crea il corpo con velocità zero per mantenerlo fermo
		newBody := body.NewRigidBody(
			units.NewQuantity(10.0, units.Kilogram),
			units.NewQuantity(1.0, units.Meter),
			position,
			vector.NewVector3(0, 0, 0), // Velocità zero
			createRandomMaterial(),
		)

		w.AddBody(newBody)
	}

	log.Printf("Formazione sferica creata con %d corpi", len(positions))
}

// createCuboidFormation crea corpi distribuiti sulle facce di un cubo
func createCuboidFormation(w world.World, count int, size, minDistance float64) {
	log.Printf("Creazione di %d corpi in formazione cuboidale", count)

	positions := make([]vector.Vector3, 0, count)
	bodiesCreated := 0

	// Continua finché non abbiamo creato abbastanza corpi
	for bodiesCreated < count {
		var position vector.Vector3

		// Scegli casualmente una delle sei facce del cubo
		face := rand.Intn(6)
		switch face {
		case 0: // Faccia +X
			position = vector.NewVector3(
				size,
				(rand.Float64()*2-1)*size,
				(rand.Float64()*2-1)*size,
			)
		case 1: // Faccia -X
			position = vector.NewVector3(
				-size,
				(rand.Float64()*2-1)*size,
				(rand.Float64()*2-1)*size,
			)
		case 2: // Faccia +Y
			position = vector.NewVector3(
				(rand.Float64()*2-1)*size,
				size,
				(rand.Float64()*2-1)*size,
			)
		case 3: // Faccia -Y
			position = vector.NewVector3(
				(rand.Float64()*2-1)*size,
				-size,
				(rand.Float64()*2-1)*size,
			)
		case 4: // Faccia +Z
			position = vector.NewVector3(
				(rand.Float64()*2-1)*size,
				(rand.Float64()*2-1)*size,
				size,
			)
		case 5: // Faccia -Z
			position = vector.NewVector3(
				(rand.Float64()*2-1)*size,
				(rand.Float64()*2-1)*size,
				-size,
			)
		}

		// Verifica distanza minima
		tooClose := false
		for _, pos := range positions {
			if position.Distance(pos) < minDistance {
				tooClose = true
				break
			}
		}

		if tooClose {
			continue
		}

		positions = append(positions, position)

		// Crea il corpo
		newBody := body.NewRigidBody(
			units.NewQuantity(10.0, units.Kilogram),
			units.NewQuantity(1.0, units.Meter),
			position,
			vector.NewVector3(0, 0, 0), // Velocità zero
			createRandomMaterial(),
		)

		w.AddBody(newBody)
		bodiesCreated++
	}

	log.Printf("Formazione cuboidale creata con %d corpi", bodiesCreated)
}

// createGrid3D crea una griglia 3D di corpi
func createGrid3D(w world.World, nx, ny, nz int, spacing float64) {
	log.Printf("Creazione di una griglia 3D %dx%dx%d", nx, ny, nz)

	// totalBodies := nx * ny * nz
	bodiesCreated := 0

	// Calcola l'offset per centrare la griglia all'origine
	offsetX := (float64(nx) - 1) * spacing / 2
	offsetY := (float64(ny) - 1) * spacing / 2
	offsetZ := (float64(nz) - 1) * spacing / 2

	// Crea la griglia
	for i := 0; i < nx; i++ {
		for j := 0; j < ny; j++ {
			for k := 0; k < nz; k++ {
				x := float64(i)*spacing - offsetX
				y := float64(j)*spacing - offsetY
				z := float64(k)*spacing - offsetZ

				position := vector.NewVector3(x, y, z)

				// Crea il corpo
				newBody := body.NewRigidBody(
					units.NewQuantity(10.0, units.Kilogram),
					units.NewQuantity(1.0, units.Meter),
					position,
					vector.NewVector3(0, 0, 0), // Velocità zero
					createRandomMaterial(),
				)

				w.AddBody(newBody)
				bodiesCreated++
			}
		}
	}

	log.Printf("Griglia 3D creata con %d corpi", bodiesCreated)
}

// createSpiral crea una spirale di corpi
func createSpiral(w world.World, count int, maxRadius, minDistance float64, turns int) {
	log.Printf("Creazione di %d corpi in una spirale", count)

	positions := make([]vector.Vector3, 0, count)

	// Angolo totale da coprire (in radianti)
	totalAngle := float64(turns) * 2 * math.Pi

	// Incremento dell'angolo per ogni corpo
	angleStep := totalAngle / float64(count)

	// Incremento del raggio per ogni corpo
	radiusStep := maxRadius / float64(count)

	// Incremento dell'altezza per ogni corpo
	heightStep := 50.0 / float64(count)

	for i := 0; i < count; i++ {
		angle := float64(i) * angleStep
		radius := float64(i) * radiusStep
		height := float64(i)*heightStep - 25.0 // centrato verticalmente

		x := radius * math.Cos(angle)
		z := radius * math.Sin(angle)
		y := height

		position := vector.NewVector3(x, y, z)

		// Verifica distanza minima
		tooClose := false
		for _, pos := range positions {
			if position.Distance(pos) < minDistance {
				tooClose = true
				break
			}
		}

		if tooClose {
			i--
			continue
		}

		positions = append(positions, position)

		// Crea il corpo
		newBody := body.NewRigidBody(
			units.NewQuantity(10.0, units.Kilogram),
			units.NewQuantity(1.0, units.Meter),
			position,
			vector.NewVector3(0, 0, 0), // Velocità zero
			createRandomMaterial(),
		)

		w.AddBody(newBody)
	}

	log.Printf("Spirale creata con %d corpi", len(positions))
}

// createRingFormation crea un anello di corpi
func createRingFormation(w world.World, count int, radius, minDistance float64) {
	log.Printf("Creazione di %d corpi in un anello", count)

	positions := make([]vector.Vector3, 0, count)

	// Angolo tra due corpi adiacenti (in radianti)
	angleStep := 2 * math.Pi / float64(count)

	for i := 0; i < count; i++ {
		angle := float64(i) * angleStep

		// Calcola la posizione sul cerchio
		x := radius * math.Cos(angle)
		z := radius * math.Sin(angle)
		y := 0.0 // anello piatto sul piano xz

		position := vector.NewVector3(x, y, z)

		// Verifica che la distanza minima sia rispettata
		if minDistance > 0 {
			tooClose := false
			for _, pos := range positions {
				if position.Distance(pos) < minDistance {
					tooClose = true
					break
				}
			}

			if tooClose {
				// Riduci il numero di corpi da creare
				count--
				continue
			}
		}

		positions = append(positions, position)

		// Crea il corpo
		newBody := body.NewRigidBody(
			units.NewQuantity(10.0, units.Kilogram),
			units.NewQuantity(1.0, units.Meter),
			position,
			vector.NewVector3(0, 0, 0), // Velocità zero
			createRandomMaterial(),
		)

		w.AddBody(newBody)
	}

	log.Printf("Anello creato con %d corpi", len(positions))
}

// createRandomMaterial crea un materiale con colore casuale
func createRandomMaterial() physMaterial.Material {
	// Genera un colore casuale
	color := [3]float64{
		rand.Float64(), // R
		rand.Float64(), // G
		rand.Float64(), // B
	}

	return physMaterial.NewBasicMaterial(
		"RandomMaterial",
		units.NewQuantity(5000, units.Kilogram),
		units.NewQuantity(800, units.Joule),
		units.NewQuantity(1.5, units.Watt),
		0.7, // emissività
		0.5, // elasticità
		color,
	)
}
