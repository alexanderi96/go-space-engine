// Package main fornisce un esempio di utilizzo di G3N con il motore fisico
package main

import (
	"log"
	"time"

	"github.com/alexanderi96/go-space-engine/core/units"
	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/physics/body"
	"github.com/alexanderi96/go-space-engine/physics/force"
	physMaterial "github.com/alexanderi96/go-space-engine/physics/material"
	"github.com/alexanderi96/go-space-engine/simulation/config"
	"github.com/alexanderi96/go-space-engine/simulation/world"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/util/helper"
)

// BodyMesh rappresenta un mesh con una luce puntuale associata
type BodyMesh struct {
	Mesh  *graphic.Mesh
	Light *light.Point
}

func main() {
	log.Println("Inizializzazione dell'esempio G3N Physics")

	// Crea la configurazione della simulazione
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
		Build()

	// Crea il mondo della simulazione
	w := world.NewPhysicalWorld(cfg.GetWorldBounds())

	// Aggiungi la forza gravitazionale
	gravityForce := force.NewGravitationalForce()
	w.AddForce(gravityForce)

	// Crea alcuni corpi
	createBodies(w)

	// Crea l'applicazione G3N
	a := app.App()

	// Crea la scena
	scene := core.NewNode()

	// Crea la camera
	cam := camera.New(1)
	cam.SetPosition(0, 15, 150)
	cam.LookAt(&math32.Vector3{0, 0, 0}, &math32.Vector3{0, 1, 0})
	scene.Add(cam)

	// Crea il controllo orbitale della camera
	camera.NewOrbitControl(cam)

	// Imposta il colore di sfondo a bianco
	a.Gls().ClearColor(1.0, 1.0, 1.0, 1.0)

	// Aggiungi luci
	ambLight := light.NewAmbient(&math32.Color{0.8, 0.8, 0.8}, 1.0)
	scene.Add(ambLight)

	pointLight1 := light.NewPoint(&math32.Color{1, 1, 1}, 2.0)
	pointLight1.SetPosition(10, 10, 10)
	scene.Add(pointLight1)

	pointLight2 := light.NewPoint(&math32.Color{1, 1, 1}, 2.0)
	pointLight2.SetPosition(-10, 10, 10)
	scene.Add(pointLight2)

	pointLight3 := light.NewPoint(&math32.Color{1, 1, 1}, 2.0)
	pointLight3.SetPosition(0, 10, -10)
	scene.Add(pointLight3)

	// Crea gli assi
	axes := helper.NewAxes(2)
	scene.Add(axes)

	// Crea una griglia per riferimento
	grid := helper.NewGrid(20, 1, &math32.Color{0.4, 0.4, 0.4})
	scene.Add(grid)

	// Crea l'interfaccia utente
	gui.Manager().Set(scene)

	// Mappa per tenere traccia dei mesh associati ai corpi
	bodyMeshes := make(map[body.ID]*BodyMesh)

	// Crea i mesh per i corpi esistenti
	for _, b := range w.GetBodies() {
		createMeshForBody(b, scene, bodyMeshes)
	}

	// Variabili per il timing
	lastUpdateTime := time.Now()
	simulationTime := 0.0

	// Avvia il loop di rendering
	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
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

		// Aggiorna la posizione dei mesh
		for _, b := range w.GetBodies() {
			if bodyMesh, exists := bodyMeshes[b.ID()]; exists {
				pos := b.Position()
				bodyMesh.Mesh.SetPosition(float32(pos.X()), float32(pos.Y()), float32(pos.Z()))
				if bodyMesh.Light != nil {
					bodyMesh.Light.SetPosition(float32(pos.X()), float32(pos.Y()), float32(pos.Z()))
				}
			}
		}

		// Operazioni OpenGL esplicite
		gl := a.Gls()
		gl.Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		gl.Enable(gls.DEPTH_TEST)

		// Renderizza la scena
		renderer.Render(scene, cam)

		// Disabilita il depth testing dopo il rendering
		gl.Disable(gls.DEPTH_TEST)
	})

	log.Println("Esempio completato")
}

// createMeshForBody crea un mesh per un corpo fisico
func createMeshForBody(b body.Body, scene *core.Node, bodyMeshes map[body.ID]*BodyMesh) {
	// Crea una sfera per rappresentare il corpo
	radius := float32(b.Radius().Value())
	geom := geometry.NewSphere(float64(radius), 32, 16)

	// Crea un materiale in base al materiale del corpo fisico
	mat := material.NewStandard(&math32.Color{0.8, 0.8, 0.8})

	// Se il materiale del corpo ha un colore, usalo
	var bodyColor math32.Color
	if b.Material() != nil {
		// Qui dovresti mappare il materiale fisico a un colore G3N
		// Per semplicità, usiamo un colore predefinito per ogni tipo di materiale
		switch b.Material().Name() {
		case "Iron":
			bodyColor = math32.Color{0.6, 0.6, 0.6}
		case "Rock":
			bodyColor = math32.Color{0.5, 0.3, 0.2}
		case "Ice":
			bodyColor = math32.Color{0.8, 0.9, 1.0}
		case "Copper":
			bodyColor = math32.Color{0.8, 0.5, 0.2}
		default:
			// Colore casuale basato sull'ID del corpo (che è una stringa)
			id := string(b.ID())
			hash := 0
			for i := 0; i < len(id); i++ {
				hash = 31*hash + int(id[i])
			}
			if hash < 0 {
				hash = -hash
			}
			r := float32(hash%255) / 255.0
			g := float32((hash/255)%255) / 255.0
			b := float32((hash/(255*255))%255) / 255.0
			bodyColor = math32.Color{r, g, b}
		}
		mat.SetColor(&bodyColor)
	} else {
		bodyColor = math32.Color{0.8, 0.8, 0.8}
	}

	// Crea un mesh con la geometria e il materiale
	mesh := graphic.NewMesh(geom, mat)

	// Imposta la posizione del mesh
	pos := b.Position()
	mesh.SetPosition(float32(pos.X()), float32(pos.Y()), float32(pos.Z()))

	// Crea una luce puntuale per il corpo
	bodyLight := light.NewPoint(&bodyColor, 0.5)
	bodyLight.SetPosition(float32(pos.X()), float32(pos.Y()), float32(pos.Z()))
	bodyLight.SetLinearDecay(1.0)
	bodyLight.SetQuadraticDecay(1.0)
	scene.Add(bodyLight)

	// Crea un BodyMesh
	bodyMesh := &BodyMesh{
		Mesh:  mesh,
		Light: bodyLight,
	}

	// Aggiungi il mesh alla scena
	scene.Add(mesh)

	// Memorizza il BodyMesh nella mappa
	bodyMeshes[b.ID()] = bodyMesh
}

// createTemporarySphere crea una sfera temporanea nella scena
func createTemporarySphere(scene *core.Node, position vector.Vector3, radius float64, color *math32.Color) {
	// Crea una sfera
	geom := geometry.NewSphere(radius, 16, 8)
	mat := material.NewStandard(color)
	mat.SetOpacity(0.7)
	mat.SetTransparent(true)

	mesh := graphic.NewMesh(geom, mat)
	mesh.SetPosition(float32(position.X()), float32(position.Y()), float32(position.Z()))

	// Aggiungi il mesh alla scena temporaneamente
	scene.Add(mesh)

	// Rimuovi il mesh dopo un po'
	go func() {
		time.Sleep(time.Millisecond * 100)
		scene.Remove(mesh)
	}()
}

// createBodies crea alcuni corpi nel mondo
func createBodies(w world.World) {
	log.Println("Creazione del sole")
	// Crea un corpo centrale massivo (come un sole)
	sun := body.NewRigidBody(
		units.NewQuantity(1.0e6, units.Kilogram),
		units.NewQuantity(1.0, units.Meter),
		vector.NewVector3(0, 0, 0),
		vector.NewVector3(0, 0, 0),
		physMaterial.Iron,
	)
	sun.SetStatic(true) // Il sole è statico (non si muove)
	w.AddBody(sun)
	log.Printf("Sole creato: ID=%v, Posizione=%v", sun.ID(), sun.Position())

	// Crea alcuni pianeti in orbita
	log.Println("Creazione dei pianeti")
	planet1 := createPlanet(w, 3.0, 0.3, 0.5, vector.NewVector3(0, 1, 0), physMaterial.Rock)
	log.Printf("Pianeta 1 creato: ID=%v, Posizione=%v", planet1.ID(), planet1.Position())
	planet2 := createPlanet(w, 5.0, 0.4, 0.3, vector.NewVector3(0, 1, 0), physMaterial.Ice)
	log.Printf("Pianeta 2 creato: ID=%v, Posizione=%v", planet2.ID(), planet2.Position())

	planet3 := createPlanet(w, 7.0, 0.5, 0.2, vector.NewVector3(0, 1, 0), physMaterial.Copper)
	log.Printf("Pianeta 3 creato: ID=%v, Posizione=%v", planet3.ID(), planet3.Position())
}

// createPlanet crea un pianeta in orbita
func createPlanet(w world.World, distance, radius, speed float64, orbitPlane vector.Vector3, mat physMaterial.Material) body.Body {
	log.Printf("Creazione di un pianeta: distanza=%f, raggio=%f, velocità=%f", distance, radius, speed)
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
	log.Printf("Pianeta aggiunto al mondo: ID=%v, Posizione=%v, Velocità=%v", planet.ID(), planet.Position(), planet.Velocity())

	return planet
}
