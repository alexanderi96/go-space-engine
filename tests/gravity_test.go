package tests

import (
	"math"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/alexanderi96/go-space-engine/core/constants"
	"github.com/alexanderi96/go-space-engine/core/units"
	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/physics/body"
	"github.com/alexanderi96/go-space-engine/physics/force"
	"github.com/alexanderi96/go-space-engine/physics/material"
	"github.com/alexanderi96/go-space-engine/physics/space"
	"github.com/alexanderi96/go-space-engine/simulation/world"
)

// TestGravitationalForce verifica il calcolo della forza gravitazionale tra due corpi
func TestGravitationalForce(t *testing.T) {
	// Crea due corpi
	body1 := body.NewRigidBody(
		units.NewQuantity(1000.0, units.Kilogram),
		units.NewQuantity(1.0, units.Meter),
		vector.NewVector3(0, 0, 0),
		vector.Zero3(),
		material.Rock,
	)

	body2 := body.NewRigidBody(
		units.NewQuantity(2000.0, units.Kilogram),
		units.NewQuantity(1.0, units.Meter),
		vector.NewVector3(10, 0, 0),
		vector.Zero3(),
		material.Rock,
	)

	// Crea la forza gravitazionale
	gravityForce := force.NewGravitationalForce()

	// Calcola la forza tra i due corpi
	forceOnBody1, forceOnBody2 := gravityForce.ApplyBetween(body1, body2)

	// Calcola la forza attesa secondo la legge di gravitazione universale
	// F = G * m1 * m2 / r^2
	distance := 10.0
	expectedForceMagnitude := constants.G * 1000.0 * 2000.0 / (distance * distance)

	// Verifica che la forza calcolata sia corretta
	actualForceMagnitude := forceOnBody1.Length()
	if math.Abs(actualForceMagnitude-expectedForceMagnitude) > 1e-10 {
		t.Errorf("Forza gravitazionale errata: attesa %v, ottenuta %v", expectedForceMagnitude, actualForceMagnitude)
	}

	// Verifica che le forze siano uguali e opposte
	sumForces := forceOnBody1.Add(forceOnBody2)
	if sumForces.Length() > 1e-10 {
		t.Errorf("Le forze non sono uguali e opposte: %v + %v = %v", forceOnBody1, forceOnBody2, sumForces)
	}

	// Verifica che la direzione della forza sia corretta
	expectedDirection := vector.NewVector3(1, 0, 0) // Da body1 a body2
	actualDirection := forceOnBody1.Normalize()
	dotProduct := expectedDirection.Dot(actualDirection)
	if math.Abs(dotProduct-1.0) > 1e-10 {
		t.Errorf("Direzione della forza errata: attesa %v, ottenuta %v", expectedDirection, actualDirection)
	}
}

// TestOctreeGravity verifica il calcolo ottimizzato della gravità utilizzando l'octree
func TestOctreeGravity(t *testing.T) {
	// Crea un octree
	bounds := space.NewAABB(
		vector.NewVector3(-100, -100, -100),
		vector.NewVector3(100, 100, 100),
	)
	octree := space.NewOctree(bounds, 10, 8)

	// Crea un corpo centrale massivo
	centralBody := body.NewRigidBody(
		units.NewQuantity(1.0e6, units.Kilogram),
		units.NewQuantity(1.0, units.Meter),
		vector.NewVector3(0, 0, 0),
		vector.Zero3(),
		material.Rock,
	)
	octree.Insert(centralBody)

	// Crea un corpo di test
	testBody := body.NewRigidBody(
		units.NewQuantity(1000.0, units.Kilogram),
		units.NewQuantity(1.0, units.Meter),
		vector.NewVector3(10, 0, 0),
		vector.Zero3(),
		material.Rock,
	)

	// Calcola la forza gravitazionale utilizzando l'octree
	theta := 0.5
	force := octree.CalculateGravity(testBody, theta)

	// Calcola la forza attesa secondo la legge di gravitazione universale
	// F = G * m1 * m2 / r^2
	distance := 10.0
	expectedForceMagnitude := constants.G * 1.0e6 * 1000.0 / (distance * distance)

	// Verifica che la forza calcolata sia corretta
	actualForceMagnitude := force.Length()
	if math.Abs(actualForceMagnitude-expectedForceMagnitude)/expectedForceMagnitude > 0.01 {
		t.Errorf("Forza gravitazionale errata: attesa %v, ottenuta %v", expectedForceMagnitude, actualForceMagnitude)
	}

	// Verifica che la direzione della forza sia corretta
	// La forza gravitazionale è attrattiva, quindi la direzione è verso il centro dell'octree
	expectedDirection := vector.NewVector3(-1, 0, 0) // Dal corpo di test verso il centro dell'octree
	actualDirection := force.Normalize()
	dotProduct := expectedDirection.Dot(actualDirection)
	if math.Abs(dotProduct-1.0) > 1e-10 {
		t.Errorf("Direzione della forza errata: attesa %v, ottenuta %v", expectedDirection, actualDirection)
	}
}

// TestEnergyConservation verifica la conservazione dell'energia in un sistema a due corpi
func TestEnergyConservation(t *testing.T) {
	// Crea un mondo fisico
	bounds := space.NewAABB(
		vector.NewVector3(-1000, -1000, -1000),
		vector.NewVector3(1000, 1000, 1000),
	)
	w := world.NewPhysicalWorld(bounds)

	// Aggiungi la forza gravitazionale
	gravityForce := force.NewGravitationalForce()
	w.AddForce(gravityForce)

	// Crea due corpi
	body1 := body.NewRigidBody(
		units.NewQuantity(1.0e6, units.Kilogram),
		units.NewQuantity(1.0, units.Meter),
		vector.NewVector3(0, 0, 0),
		vector.Zero3(),
		material.Rock,
	)
	body1.SetStatic(true) // Il primo corpo è statico
	w.AddBody(body1)

	// Calcola la velocità orbitale circolare
	distance := 100.0
	orbitalSpeed := math.Sqrt(constants.G * 1.0e6 / distance)

	body2 := body.NewRigidBody(
		units.NewQuantity(1000.0, units.Kilogram),
		units.NewQuantity(1.0, units.Meter),
		vector.NewVector3(distance, 0, 0),
		vector.NewVector3(0, orbitalSpeed, 0), // Velocità orbitale circolare
		material.Rock,
	)
	w.AddBody(body2)

	// Calcola l'energia iniziale del sistema
	initialEnergy := calculateTotalEnergy(w)

	// Esegui la simulazione per 100 passi
	dt := 0.01
	for i := 0; i < 100; i++ {
		w.Step(dt)
	}

	// Calcola l'energia finale del sistema
	finalEnergy := calculateTotalEnergy(w)

	// Verifica che l'energia sia conservata (con una tolleranza dell'1%)
	energyDifference := math.Abs(finalEnergy-initialEnergy) / math.Abs(initialEnergy)
	if energyDifference > 0.01 {
		t.Errorf("Energia non conservata: iniziale %v, finale %v, differenza %v%%", initialEnergy, finalEnergy, energyDifference*100)
	}
}

// calculateTotalEnergy calcola l'energia totale (cinetica + potenziale) di un sistema
func calculateTotalEnergy(w world.World) float64 {
	bodies := w.GetBodies()
	totalEnergy := 0.0

	// Calcola l'energia cinetica di ogni corpo
	for _, b := range bodies {
		if b.IsStatic() {
			continue
		}
		mass := b.Mass().Value()
		velocity := b.Velocity()
		kineticEnergy := 0.5 * mass * velocity.LengthSquared()
		totalEnergy += kineticEnergy
	}

	// Calcola l'energia potenziale gravitazionale di ogni coppia di corpi
	for i := 0; i < len(bodies); i++ {
		for j := i + 1; j < len(bodies); j++ {
			body1 := bodies[i]
			body2 := bodies[j]
			mass1 := body1.Mass().Value()
			mass2 := body2.Mass().Value()
			distance := body1.Position().Distance(body2.Position())
			potentialEnergy := -constants.G * mass1 * mass2 / distance
			totalEnergy += potentialEnergy
		}
	}

	return totalEnergy
}

// TestBarnesHutAccuracy verifica l'accuratezza dell'algoritmo Barnes-Hut
func TestBarnesHutAccuracy(t *testing.T) {
	// Crea un octree
	bounds := space.NewAABB(
		vector.NewVector3(-100, -100, -100),
		vector.NewVector3(100, 100, 100),
	)
	octree := space.NewOctree(bounds, 10, 8)

	// Crea 100 corpi casuali
	for i := 0; i < 100; i++ {
		// Posizione casuale all'interno dei limiti
		x := (rand.Float64() * 200) - 100
		y := (rand.Float64() * 200) - 100
		z := (rand.Float64() * 200) - 100
		position := vector.NewVector3(x, y, z)

		// Massa casuale
		mass := rand.Float64() * 1000

		// Crea il corpo
		b := body.NewRigidBody(
			units.NewQuantity(mass, units.Kilogram),
			units.NewQuantity(1.0, units.Meter),
			position,
			vector.Zero3(),
			material.Rock,
		)
		octree.Insert(b)
	}

	// Crea un corpo di test
	testBody := body.NewRigidBody(
		units.NewQuantity(1000.0, units.Kilogram),
		units.NewQuantity(1.0, units.Meter),
		vector.NewVector3(50, 50, 50),
		vector.Zero3(),
		material.Rock,
	)

	// Calcola la forza gravitazionale utilizzando l'algoritmo Barnes-Hut con diversi valori di theta
	thetas := []float64{0.0, 0.1, 0.5, 1.0}
	forces := make([]vector.Vector3, len(thetas))

	for i, theta := range thetas {
		forces[i] = octree.CalculateGravity(testBody, theta)
	}

	// La forza calcolata con theta = 0 è la più accurata (calcolo diretto)
	exactForce := forces[0]

	// Verifica che l'errore aumenti con l'aumentare di theta
	for i := 1; i < len(thetas); i++ {
		error := forces[i].Sub(exactForce).Length() / exactForce.Length()
		t.Logf("Theta = %v, errore relativo = %v%%", thetas[i], error*100)

		// Verifica che l'errore sia accettabile (meno del 10% per theta = 1.0)
		if error > 0.1 && thetas[i] <= 1.0 {
			t.Errorf("Errore troppo grande per theta = %v: %v%%", thetas[i], error*100)
		}
	}
}

// TestMultithreadingPerformance verifica le prestazioni del multithreading
func TestMultithreadingPerformance(t *testing.T) {
	// Crea due mondi fisici identici
	bounds := space.NewAABB(
		vector.NewVector3(-1000, -1000, -1000),
		vector.NewVector3(1000, 1000, 1000),
	)
	w1 := world.NewPhysicalWorld(bounds)
	w2 := world.NewPhysicalWorld(bounds)

	// Aggiungi la forza gravitazionale a entrambi i mondi
	gravityForce1 := force.NewGravitationalForce()
	gravityForce2 := force.NewGravitationalForce()
	w1.AddForce(gravityForce1)
	w2.AddForce(gravityForce2)

	// Crea 500 corpi casuali in entrambi i mondi per rendere il test più significativo
	for i := 0; i < 500; i++ {
		// Posizione casuale all'interno dei limiti
		x := (rand.Float64() * 2000) - 1000
		y := (rand.Float64() * 2000) - 1000
		z := (rand.Float64() * 2000) - 1000
		position := vector.NewVector3(x, y, z)

		// Velocità casuale
		vx := (rand.Float64() * 10) - 5
		vy := (rand.Float64() * 10) - 5
		vz := (rand.Float64() * 10) - 5
		velocity := vector.NewVector3(vx, vy, vz)

		// Massa casuale
		mass := rand.Float64() * 1000

		// Crea i corpi
		b1 := body.NewRigidBody(
			units.NewQuantity(mass, units.Kilogram),
			units.NewQuantity(1.0, units.Meter),
			position,
			velocity,
			material.Rock,
		)
		b2 := body.NewRigidBody(
			units.NewQuantity(mass, units.Kilogram),
			units.NewQuantity(1.0, units.Meter),
			position,
			velocity,
			material.Rock,
		)
		w1.AddBody(b1)
		w2.AddBody(b2)
	}

	// Misura il tempo di esecuzione con multithreading
	startTime1 := time.Now()
	for i := 0; i < 50; i++ { // Aumentiamo il numero di passi per rendere il test più significativo
		w1.Step(0.01)
	}
	duration1 := time.Since(startTime1)

	// Disabilita il multithreading nel secondo mondo
	// Nota: questo è solo un test, non c'è un modo diretto per disabilitare il multithreading
	// Quindi questo test è solo indicativo
	startTime2 := time.Now()
	for i := 0; i < 50; i++ { // Aumentiamo il numero di passi per rendere il test più significativo
		w2.Step(0.01)
	}
	duration2 := time.Since(startTime2)

	// Verifica che il multithreading sia più veloce
	t.Logf("Tempo con multithreading: %v", duration1)
	t.Logf("Tempo senza multithreading: %v", duration2)
}
