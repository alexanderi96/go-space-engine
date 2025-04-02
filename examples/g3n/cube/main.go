// Package main fornisce un esempio di simulazione con corpi disposti in forme geometriche prestabilite
package main

import (
	"log"
	"math"
	"math/rand"
	"os"
	"runtime/pprof"
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

const (
	shouldBeProfiled = true
)

func main() {

	log.Println("Inizializzazione della simulazione con corpi in forme geometriche con interazione gravitazionale")

	if shouldBeProfiled {
		f, err := os.Create("cpu.pprof")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// Crea la configurazione della simulazione
	cfg := config.NewSimulationBuilder().
		WithTimeStep(0.01).
		WithMaxBodies(5000).
		WithGravity(true).
		WithCollisions(true).
		WithBoundaryCollisions(true).
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

	// Scegli una delle configurazioni seguenti
	// Puoi commentare/decommentare quella desiderata

	// Creazione di un corpo centrale massiccio (opzionale)
	// createCentralBody(w)

	// Crea un cubo di corpi (non statici, influenzati dalla gravità)
	createCuboidFormation(w, 128, 50.0, 100.0, 5.0)

	// Altre formazioni disponibili:
	// createSphereFormation(w, 300, 40.0, 80.0, 5.0)
	// createRingFormation(w, 200, 60.0, 80.0, 5.0)
	// createSpiralFormation(w, 200, 20.0, 80.0, 5.0, 3)
	// createBinarySystem(w, 200, 5.0)

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
		w.Step(0.01)

		// Renderizza il mondo
		adapter.RenderWorld(w)
	})

	log.Println("Simulazione completata")
}

// createCentralBody crea un corpo centrale massiccio
func createCentralBody(w world.World) {
	log.Println("Creazione del corpo centrale")

	// Massa elevata per il corpo centrale
	centralMass := 1.5e11

	centralBody := body.NewRigidBody(
		units.NewQuantity(centralMass, units.Kilogram),
		units.NewQuantity(8.0, units.Meter),
		vector.NewVector3(0, 0, 0),
		vector.NewVector3(0, 0, 0),
		createMaterial("CentralBody", 0.9, 0.5, [3]float64{1.0, 0.6, 0.0}),
	)

	// Il corpo centrale può essere statico o dinamico
	// Se vuoi che TUTTO sia influenzato dalla gravità, commenta la linea seguente
	centralBody.SetStatic(true)

	w.AddBody(centralBody)
	log.Printf("Corpo centrale creato: ID=%v, Posizione=%v", centralBody.ID(), centralBody.Position())
}

// createCuboidFormation crea corpi distribuiti in un cubo
func createCuboidFormation(w world.World, count int, minSize, maxSize, minDistance float64) {
	log.Printf("Creazione di %d corpi in formazione cubica", count)

	// Massa del corpo centrale (se presente)
	// centralMass := 1.5e11

	// Determina il numero di corpi per lato per ottenere un cubo perfetto
	// Calcola la radice cubica arrotondata all'intero più vicino
	bodiesPerSide := int(math.Ceil(math.Pow(float64(count), 1.0/3.0)))
	actualCount := bodiesPerSide * bodiesPerSide * bodiesPerSide

	log.Printf("Creazione di un cubo %dx%dx%d con %d corpi totali",
		bodiesPerSide, bodiesPerSide, bodiesPerSide, actualCount)

	// Calcola la spaziatura tra i corpi
	spacing := minDistance + 1.0 // Assicura una distanza minima tra i corpi

	// Determina la dimensione totale del cubo
	cubeSize := float64(bodiesPerSide-1) * spacing
	halfSize := cubeSize / 2.0

	// Posizioni dei corpi già creati
	positions := make([]vector.Vector3, 0, actualCount)

	// Crea il reticolo cubico
	for x := 0; x < bodiesPerSide; x++ {
		for y := 0; y < bodiesPerSide; y++ {
			for z := 0; z < bodiesPerSide; z++ {
				// Calcola la posizione nel reticolo
				posX := float64(x)*spacing - halfSize
				posY := float64(y)*spacing - halfSize
				posZ := float64(z)*spacing - halfSize

				position := vector.NewVector3(posX, posY, posZ)
				positions = append(positions, position)

				// Crea un corpo con massa casuale ma non troppo grande
				bodyMass := (rand.Float64()*20 + 5.0) * 1e9 // Massa tra 5 e 25

				// Calcola una piccola velocità iniziale casuale per evitare una configurazione perfettamente stabile
				// velMagnitude := rand.Float64() * 0.5 // Velocità iniziale ridotta
				// velX := (rand.Float64()*2.0 - 1.0) * velMagnitude
				// velY := (rand.Float64()*2.0 - 1.0) * velMagnitude
				// velZ := (rand.Float64()*2.0 - 1.0) * velMagnitude

				// velocity := vector.NewVector3(velX, velY, velZ)

				// Crea il corpo
				newBody := body.NewRigidBody(
					units.NewQuantity(bodyMass, units.Kilogram),
					units.NewQuantity(rand.Float64()*0.5+0.5, units.Meter), // Raggio casuale
					position,
					// velocity,
					vector.NewVector3(0, 0, 0),
					createRandomMaterial(),
				)

				// Importante: NON impostare il corpo come statico
				// newBody.SetStatic(false) - questo è il comportamento predefinito

				w.AddBody(newBody)
			}
		}
	}

	log.Printf("Formazione cubica creata con %d corpi", len(positions))
}

// createSphereFormation crea corpi distribuiti in una sfera
func createSphereFormation(w world.World, count int, minRadius, maxRadius, minDistance float64) {
	log.Printf("Creazione di %d corpi in formazione sferica", count)

	// Massa del corpo centrale
	centralMass := 1.5e11

	positions := make([]vector.Vector3, 0, count)

	for i := 0; i < count; i++ {
		// Genera un punto sulla sfera con distribuzione uniforme
		phi := rand.Float64() * 2 * math.Pi
		costheta := rand.Float64()*2 - 1
		theta := math.Acos(costheta)

		// Raggio random tra minRadius e maxRadius
		radius := minRadius + rand.Float64()*(maxRadius-minRadius)

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

		// Calcola la velocità orbitale (ma con una componente casuale)
		distance := position.Length()
		baseSpeed := math.Sqrt(constants.G*centralMass/distance) * 0.8 // 80% della velocità orbitale teorica

		// Calcola direzione della velocità con un po' di casualità
		radialDirection := position.Normalize()

		// Scegliamo un vettore di riferimento
		reference := vector.NewVector3(0, 1, 0)
		if math.Abs(radialDirection.Dot(reference)) > 0.9 {
			reference = vector.NewVector3(1, 0, 0)
		}

		// Calcola il vettore perpendicolare
		tangent := reference.Cross(radialDirection).Normalize()

		// Aggiungi componente casuale alla velocità
		randomFactor := 0.2 // 20% di casualità
		velX := tangent.X() * baseSpeed * (1.0 + (rand.Float64()*2-1)*randomFactor)
		velY := tangent.Y() * baseSpeed * (1.0 + (rand.Float64()*2-1)*randomFactor)
		velZ := tangent.Z() * baseSpeed * (1.0 + (rand.Float64()*2-1)*randomFactor)

		velocity := vector.NewVector3(velX, velY, velZ)

		// Crea il corpo
		newBody := body.NewRigidBody(
			units.NewQuantity(rand.Float64()*50+10, units.Kilogram),
			units.NewQuantity(rand.Float64()*0.5+0.5, units.Meter),
			position,
			velocity,
			createRandomMaterial(),
		)

		w.AddBody(newBody)
	}

	log.Printf("Formazione sferica creata con %d corpi", len(positions))
}

// createRingFormation crea corpi distribuiti in un anello
func createRingFormation(w world.World, count int, minRadius, maxRadius, minDistance float64) {
	log.Printf("Creazione di %d corpi in formazione ad anello", count)

	// Massa del corpo centrale
	centralMass := 1.5e11

	positions := make([]vector.Vector3, 0, count)

	for i := 0; i < count; i++ {
		// Genera un punto in un anello
		angle := rand.Float64() * 2 * math.Pi
		radius := minRadius + rand.Float64()*(maxRadius-minRadius)

		// Genera una piccola variazione verticale
		height := (rand.Float64()*2 - 1) * (maxRadius - minRadius) * 0.1

		x := radius * math.Cos(angle)
		y := height
		z := radius * math.Sin(angle)

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

		// Calcola la velocità orbitale
		orbitSpeed := math.Sqrt(constants.G * centralMass / radius)

		// La velocità deve essere perpendicolare al raggio sul piano dell'anello
		velocity := vector.NewVector3(
			-orbitSpeed*math.Sin(angle),         // Componente x
			(rand.Float64()*2-1)*0.1*orbitSpeed, // Piccola componente verticale casuale
			orbitSpeed*math.Cos(angle),          // Componente z
		)

		// Crea il corpo
		newBody := body.NewRigidBody(
			units.NewQuantity(rand.Float64()*20+5, units.Kilogram),
			units.NewQuantity(rand.Float64()*0.5+0.3, units.Meter),
			position,
			velocity,
			createRandomMaterial(),
		)

		w.AddBody(newBody)
	}

	log.Printf("Formazione ad anello creata con %d corpi", len(positions))
}

// createSpiralFormation crea corpi distribuiti in una spirale
func createSpiralFormation(w world.World, count int, minRadius, maxRadius, minDistance float64, arms int) {
	log.Printf("Creazione di %d corpi in formazione a spirale con %d bracci", count, arms)

	// Massa del corpo centrale
	centralMass := 1.5e11

	positions := make([]vector.Vector3, 0, count)

	// Parametri della spirale
	turns := 2.0 // Numero di giri completi della spirale

	for i := 0; i < count; i++ {
		// Scegli un braccio casuale
		arm := rand.Intn(arms)

		// Parametro t varia da 0 a 1 lungo la spirale
		t := rand.Float64()

		// Angolo base per questo braccio della spirale
		baseAngle := 2.0 * math.Pi * float64(arm) / float64(arms)

		// Angolo che aumenta con t
		angle := baseAngle + turns*2.0*math.Pi*t

		// Il raggio aumenta con t
		radius := minRadius + t*(maxRadius-minRadius)

		// Aggiungi un po' di variazione al raggio
		radius += (rand.Float64()*2 - 1) * (maxRadius - minRadius) * 0.05

		// Coordiante x, y, z
		x := radius * math.Cos(angle)
		y := (rand.Float64()*2 - 1) * maxRadius * 0.05 // Piccola variazione sull'asse y
		z := radius * math.Sin(angle)

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

		// Calcola la velocità orbitale
		orbitSpeed := math.Sqrt(constants.G*centralMass/radius) * 0.9

		// La velocità deve essere perpendicolare al raggio
		velocity := vector.NewVector3(
			-orbitSpeed*math.Sin(angle),
			(rand.Float64()*2-1)*0.05*orbitSpeed, // Piccola componente verticale
			orbitSpeed*math.Cos(angle),
		)

		// Crea il corpo
		newBody := body.NewRigidBody(
			units.NewQuantity(rand.Float64()*15+5, units.Kilogram),
			units.NewQuantity(rand.Float64()*0.4+0.3, units.Meter),
			position,
			velocity,
			createRandomMaterial(),
		)

		w.AddBody(newBody)
	}

	log.Printf("Formazione a spirale creata con %d corpi", len(positions))
}

// createBinarySystem crea un sistema binario con due corpi centrali massicci
func createBinarySystem(w world.World, satelliteCount int, minDistance float64) {
	log.Println("Creazione di un sistema binario")

	// Crea due corpi massicci
	mass1 := 7.5e10
	mass2 := 5.0e10
	separation := 40.0

	// Calcolo della velocità orbitale per i due corpi centrali
	// Assumiamo che i corpi orbitino attorno al loro centro di massa
	totalMass := mass1 + mass2

	// Posizione del centro di massa
	centerOfMassX := (mass1*(-separation/2) + mass2*(separation/2)) / totalMass

	// Distanza effettiva di ciascun corpo dal centro di massa
	dist1 := math.Abs((-separation / 2) - centerOfMassX)
	dist2 := math.Abs((separation / 2) - centerOfMassX)

	// Velocità orbitale
	orbitPeriod := 2 * math.Pi * math.Sqrt(math.Pow(separation, 3)/(constants.G*totalMass))
	speed1 := 2 * math.Pi * dist1 / orbitPeriod
	speed2 := 2 * math.Pi * dist2 / orbitPeriod

	// Creazione del primo corpo centrale
	body1 := body.NewRigidBody(
		units.NewQuantity(mass1, units.Kilogram),
		units.NewQuantity(5.0, units.Meter),
		vector.NewVector3(-separation/2, 0, 0),
		vector.NewVector3(0, 0, speed1),
		createMaterial("CentralBody1", 0.9, 0.5, [3]float64{0.9, 0.6, 0.1}),
	)
	w.AddBody(body1)

	// Creazione del secondo corpo centrale
	body2 := body.NewRigidBody(
		units.NewQuantity(mass2, units.Kilogram),
		units.NewQuantity(4.0, units.Meter),
		vector.NewVector3(separation/2, 0, 0),
		vector.NewVector3(0, 0, -speed2),
		createMaterial("CentralBody2", 0.9, 0.5, [3]float64{0.2, 0.6, 0.9}),
	)
	w.AddBody(body2)

	log.Println("Corpi centrali del sistema binario creati")

	// Crea satelliti attorno al sistema binario
	createSatellites(w, satelliteCount, separation*1.5, separation*5, minDistance, mass1+mass2)
}

// createSatellites crea satelliti attorno a un punto centrale
func createSatellites(w world.World, count int, minRadius, maxRadius, minDistance, centralMass float64) {
	log.Printf("Creazione di %d satelliti", count)

	positions := make([]vector.Vector3, 0, count)

	for i := 0; i < count; i++ {
		// Posizione casuale in una sfera
		phi := rand.Float64() * 2 * math.Pi
		costheta := rand.Float64()*2 - 1
		theta := math.Acos(costheta)

		radius := minRadius + rand.Float64()*(maxRadius-minRadius)

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

		// Calcolo della velocità orbitale
		distance := position.Length()
		orbitSpeed := math.Sqrt(constants.G * centralMass / distance)

		// Calcola direzione della velocità
		radialDirection := position.Normalize()

		// Vettore di riferimento
		reference := vector.NewVector3(0, 1, 0)
		if math.Abs(radialDirection.Dot(reference)) > 0.9 {
			reference = vector.NewVector3(1, 0, 0)
		}

		// Calcola il vettore perpendicolare
		tangent := reference.Cross(radialDirection).Normalize()
		velocity := tangent.Scale(orbitSpeed)

		// Crea il satellite
		satellite := body.NewRigidBody(
			units.NewQuantity(rand.Float64()*10+1, units.Kilogram),
			units.NewQuantity(rand.Float64()*0.4+0.2, units.Meter),
			position,
			velocity,
			createRandomMaterial(),
		)

		w.AddBody(satellite)
	}

	log.Printf("Satelliti creati: %d", len(positions))
}

// createRandomMaterial crea un materiale con colore casuale
func createRandomMaterial() physMaterial.Material {
	// Genera un colore casuale
	r := rand.Float64()*0.7 + 0.3
	g := rand.Float64()*0.7 + 0.3
	b := rand.Float64()*0.7 + 0.3

	return physMaterial.NewBasicMaterial(
		"RandomMaterial",
		units.NewQuantity(5000, units.Kilogram),
		units.NewQuantity(800, units.Joule),
		units.NewQuantity(1.5, units.Watt),
		0.7+rand.Float64()*0.3, // Emissività tra 0.7 e 1.0
		0.3+rand.Float64()*0.6, // Elasticità tra 0.3 e 0.9
		[3]float64{r, g, b},    // Colore casuale
	)
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
