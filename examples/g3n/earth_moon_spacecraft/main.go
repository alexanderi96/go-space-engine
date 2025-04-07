package main

import (
	"log"
	"time"

	"github.com/alexanderi96/go-space-engine/core/units"
	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/entity"
	"github.com/alexanderi96/go-space-engine/entity/vehicle/spacecraft"
	"github.com/alexanderi96/go-space-engine/physics/body"
	"github.com/alexanderi96/go-space-engine/physics/force"
	physMaterial "github.com/alexanderi96/go-space-engine/physics/material"
	"github.com/alexanderi96/go-space-engine/render/g3n"
	"github.com/alexanderi96/go-space-engine/simulation/config"
	"github.com/alexanderi96/go-space-engine/simulation/world"
	"github.com/g3n/engine/math32"
)

const (
	// Parametri del controller della navicella
	maxAngularVelocity = 60.0 // 60 gradi al secondo
	angularDamping     = 0.3  // 30% di smorzamento
	pidProportional    = 3.0  // Guadagno P
	pidIntegral        = 0.1  // Guadagno I
	pidDerivative      = 0.5  // Guadagno D
	maxThrust          = 5000.0
	maxTorque          = 1000.0
)

func main() {
	log.Println("Initializing Earth-Moon-Spacecraft Example")

	// Crea la configurazione della simulazione
	cfg := config.NewSimulationBuilder().
		WithTimeStep(0.01).
		WithMaxBodies(100).
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
	gravityForce.SetTheta(0.5) // Imposta il valore theta per l'algoritmo Barnes-Hut
	w.AddForce(gravityForce)

	// Crea i corpi celesti (Terra e Luna) e la navicella
	_, _, spacecraftEntity, spacecraftController := createBodies(w)

	// Crea l'adapter G3N
	adapter := g3n.NewG3NAdapter()

	// Configura l'adapter
	adapter.SetBackgroundColor(g3n.NewColor(0.9, 0.9, 0.9, 1.0)) // Sfondo blu scuro per lo spazio

	// Abilita la modalità di debug per la diagnostica
	adapter.SetDebugMode(true)

	// Crea l'handler di input per la navicella
	inputHandler := NewSpacecraftInputHandler(spacecraftController)

	// Registra l'handler di input
	adapter.RegisterInputHandler(inputHandler)

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

		// Aggiorna l'handler di input
		inputHandler.Update(dt)

		// Esegui un passo della simulazione
		w.Step(dt)

		// Ottieni la posizione aggiornata della navicella
		if spacecraftEntity != nil {
			spacecraftPos := spacecraftEntity.GetBody().Position()

			// Aggiorna la posizione della camera per seguire la navicella
			// Posiziona la camera dietro e leggermente sopra la navicella
			adapter.GetCamera().SetPosition(
				float32(spacecraftPos.X()-10),
				float32(spacecraftPos.Y()+5),
				float32(spacecraftPos.Z()-10),
			)
			adapter.GetCamera().LookAt(
				&math32.Vector3{
					float32(spacecraftPos.X()),
					float32(spacecraftPos.Y()),
					float32(spacecraftPos.Z()),
				},
				&math32.Vector3{0, 1, 0},
			)
		}

		log.Printf("Spacecraft: ID=%v, Position=%v", spacecraftEntity.GetID(), spacecraftEntity.GetBody().Position())

		// Renderizza il mondo
		adapter.RenderWorld(w)
	})

	log.Println("Example completed")
}

// createBodies crea i corpi celesti e la navicella
func createBodies(w world.World) (earth, moon body.Body, spacecraft entity.Entity, controller *spacecraft.SpacecraftController) {
	// Crea la Terra
	earth = body.NewRigidBody(
		units.NewQuantity(5.97e4, units.Kilogram),                    // Massa della Terra
		units.NewQuantity(1.3, units.Meter),                          // Raggio della Terra (come in solar_system)
		vector.NewVector3(0, 0, 0),                                   // Posizione al centro
		vector.NewVector3(0, 0, 0),                                   // Velocità zero
		createMaterial("Earth", 0.7, 0.5, [3]float64{0.0, 0.3, 0.8}), // Colore blu
	)
	w.AddBody(earth)

	// Crea la Luna
	moon = body.NewRigidBody(
		units.NewQuantity(7.34e2, units.Kilogram),                   // Massa della Luna
		units.NewQuantity(0.4, units.Meter),                         // Raggio della Luna (proporzionato)
		vector.NewVector3(20, 0, 0),                                 // Posizione a 20 unità dalla Terra
		vector.NewVector3(0, 0, 1.5),                                // Velocità orbitale
		createMaterial("Moon", 0.7, 0.5, [3]float64{0.8, 0.8, 0.8}), // Colore grigio
	)
	w.AddBody(moon)

	// Crea la navicella
	spacecraft, controller = createSpacecraft(w, earth)

	return earth, moon, spacecraft, controller
}

// createSpacecraft crea una navicella controllabile vicino alla Terra
func createSpacecraft(w world.World, planet body.Body) (entity.Entity, *spacecraft.SpacecraftController) {
	log.Println("Creating controllable spacecraft near Earth")

	// Ottieni la posizione e la velocità del pianeta
	planetPos := planet.Position()
	planetVel := planet.Velocity()
	planetRadius := planet.Radius().Value()

	// Calcola la posizione della navicella (1.5 volte il raggio del pianeta)
	offset := vector.NewVector3(0, 0, 1).Scale(planetRadius * 1.5)
	spacecraftPos := planetPos.Add(offset)

	// La navicella inizia con la stessa velocità del pianeta
	spacecraftVel := vector.NewVector3(planetVel.X(), planetVel.Y(), planetVel.Z())

	// Crea la configurazione della navicella
	config := spacecraft.SpacecraftConfig{
		Mass:      1000.0,
		Radius:    0.5,
		MaxThrust: maxThrust,
		MaxTorque: maxTorque,
		Position:  spacecraftPos,
		Velocity:  spacecraftVel,
		Rotation:  vector.Zero3(),
		IsCube:    true,                      // Indica che la navicella deve essere un cubo
		Color:     [3]float64{1.0, 1.0, 1.0}, // Colore bianco
	}

	// Crea la navicella e il controller
	spacecraftEntity, controller := spacecraft.CreateSpacecraft(config)

	// Configura il controller
	controller.SetMaxAngularVelocity(maxAngularVelocity)
	controller.SetAngularDamping(angularDamping)
	controller.SetRotationPIDGains(pidProportional, pidIntegral, pidDerivative)

	// Aggiungi il corpo della navicella al mondo (se non è già stato aggiunto dal factory)
	w.AddBody(spacecraftEntity.GetBody())

	log.Printf("Spacecraft created: ID=%v, Position=%v", spacecraftEntity.GetID(), spacecraftPos)

	return spacecraftEntity, controller
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
