// Package g3n fornisce un'implementazione dell'interfaccia RenderAdapter utilizzando G3N
package g3n

import (
	"time"

	"github.com/alexanderi96/go-space-engine/physics/body"
	"github.com/alexanderi96/go-space-engine/render/adapter"
	"github.com/alexanderi96/go-space-engine/simulation/world"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
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

// G3NAdapter è un adattatore per il rendering con G3N
type G3NAdapter struct {
	app        *app.Application
	scene      *core.Node
	camera     *camera.Camera
	cameraCtrl *camera.OrbitControl
	bodyMeshes map[body.ID]*BodyMesh
	bgColor    adapter.Color
	debugMode  bool
}

// NewG3NAdapter crea un nuovo adattatore G3N
func NewG3NAdapter() *G3NAdapter {
	return &G3NAdapter{
		bodyMeshes: make(map[body.ID]*BodyMesh),
		bgColor:    adapter.NewColor(1.0, 1.0, 1.0, 1.0), // Sfondo bianco
		debugMode:  false,
	}
}

// GetRenderer restituisce il renderer (implementazione dell'interfaccia RenderAdapter)
func (ga *G3NAdapter) GetRenderer() adapter.Renderer {
	// Questo adapter non utilizza l'interfaccia Renderer standard
	// Restituisce nil perché implementa direttamente i metodi necessari
	return nil
}

// RenderWorld renderizza il mondo
func (ga *G3NAdapter) RenderWorld(w world.World) {
	// Aggiorna la posizione dei mesh
	for _, b := range w.GetBodies() {
		if bodyMesh, exists := ga.bodyMeshes[b.ID()]; exists {
			pos := b.Position()
			bodyMesh.Mesh.SetPosition(float32(pos.X()), float32(pos.Y()), float32(pos.Z()))
			if bodyMesh.Light != nil {
				bodyMesh.Light.SetPosition(float32(pos.X()), float32(pos.Y()), float32(pos.Z()))
			}
		} else {
			// Se il corpo non ha un mesh associato, crealo
			ga.createMeshForBody(b)
		}
	}
}

// Run avvia il loop di rendering
func (ga *G3NAdapter) Run(updateFunc func(deltaTime time.Duration)) {
	// Inizializza l'applicazione G3N se non è già stata inizializzata
	if ga.app == nil {
		ga.initialize()
	}

	// Avvia il loop di rendering
	ga.app.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		// Operazioni OpenGL esplicite
		gl := ga.app.Gls()
		gl.Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		gl.Enable(gls.DEPTH_TEST)

		// Chiama la funzione di aggiornamento fornita
		if updateFunc != nil {
			updateFunc(deltaTime)
		}

		// Renderizza la scena
		renderer.Render(ga.scene, ga.camera)

		// Disabilita il depth testing dopo il rendering
		gl.Disable(gls.DEPTH_TEST)
	})
}

// initialize inizializza l'adapter
func (ga *G3NAdapter) initialize() {
	// Crea l'applicazione G3N
	ga.app = app.App()

	// Crea la scena
	ga.scene = core.NewNode()

	// Crea la camera
	ga.camera = camera.New(1)
	ga.camera.SetPosition(0, 1, 15)
	ga.camera.LookAt(&math32.Vector3{0, 0, 0}, &math32.Vector3{0, 1, 0})
	ga.scene.Add(ga.camera)

	// Crea il controllo orbitale della camera
	ga.cameraCtrl = camera.NewOrbitControl(ga.camera)

	// Imposta il colore di sfondo
	ga.app.Gls().ClearColor(float32(ga.bgColor.R), float32(ga.bgColor.G), float32(ga.bgColor.B), float32(ga.bgColor.A))

	// Aggiungi luci
	ambLight := light.NewAmbient(&math32.Color{0.8, 0.8, 0.8}, 1.0)
	ga.scene.Add(ambLight)

	pointLight1 := light.NewPoint(&math32.Color{1, 1, 1}, 2.0)
	pointLight1.SetPosition(10, 10, 10)
	ga.scene.Add(pointLight1)

	pointLight2 := light.NewPoint(&math32.Color{1, 1, 1}, 2.0)
	pointLight2.SetPosition(-10, 10, 10)
	ga.scene.Add(pointLight2)

	pointLight3 := light.NewPoint(&math32.Color{1, 1, 1}, 2.0)
	pointLight3.SetPosition(0, 10, -10)
	ga.scene.Add(pointLight3)

	// Crea gli assi
	axes := helper.NewAxes(2)
	ga.scene.Add(axes)

	// Crea una griglia per riferimento
	grid := helper.NewGrid(20, 1, &math32.Color{0.4, 0.4, 0.4})
	ga.scene.Add(grid)
}

// createMeshForBody crea un mesh per un corpo fisico
func (ga *G3NAdapter) createMeshForBody(b body.Body) {
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
	ga.scene.Add(bodyLight)

	// Aggiungi il mesh alla scena
	ga.scene.Add(mesh)

	// Memorizza il BodyMesh nella mappa
	ga.bodyMeshes[b.ID()] = &BodyMesh{
		Mesh:  mesh,
		Light: bodyLight,
	}
}

// SetDebugMode imposta la modalità di debug
func (ga *G3NAdapter) SetDebugMode(debug bool) {
	ga.debugMode = debug
}

// IsDebugMode restituisce true se la modalità di debug è attiva
func (ga *G3NAdapter) IsDebugMode() bool {
	return ga.debugMode
}

// SetRenderOctree imposta se renderizzare l'octree
func (ga *G3NAdapter) SetRenderOctree(render bool) {
	// Non implementato in questo adapter
}

// IsRenderOctree restituisce true se l'octree viene renderizzato
func (ga *G3NAdapter) IsRenderOctree() bool {
	return false
}

// SetRenderBoundingBoxes imposta se renderizzare i bounding box
func (ga *G3NAdapter) SetRenderBoundingBoxes(render bool) {
	// Non implementato in questo adapter
}

// IsRenderBoundingBoxes restituisce true se i bounding box vengono renderizzati
func (ga *G3NAdapter) IsRenderBoundingBoxes() bool {
	return false
}

// SetRenderVelocities imposta se renderizzare i vettori velocità
func (ga *G3NAdapter) SetRenderVelocities(render bool) {
	// Non implementato in questo adapter
}

// IsRenderVelocities restituisce true se i vettori velocità vengono renderizzati
func (ga *G3NAdapter) IsRenderVelocities() bool {
	return false
}

// SetRenderAccelerations imposta se renderizzare i vettori accelerazione
func (ga *G3NAdapter) SetRenderAccelerations(render bool) {
	// Non implementato in questo adapter
}

// IsRenderAccelerations restituisce true se i vettori accelerazione vengono renderizzati
func (ga *G3NAdapter) IsRenderAccelerations() bool {
	return false
}

// SetRenderForces imposta se renderizzare i vettori forza
func (ga *G3NAdapter) SetRenderForces(render bool) {
	// Non implementato in questo adapter
}

// IsRenderForces restituisce true se i vettori forza vengono renderizzati
func (ga *G3NAdapter) IsRenderForces() bool {
	return false
}

// SetBackgroundColor imposta il colore di sfondo
func (ga *G3NAdapter) SetBackgroundColor(color adapter.Color) {
	ga.bgColor = color
	if ga.app != nil {
		ga.app.Gls().ClearColor(float32(color.R), float32(color.G), float32(color.B), float32(color.A))
	}
}

// NewColor crea un nuovo colore
func NewColor(r, g, b, a float64) adapter.Color {
	return adapter.NewColor(r, g, b, a)
}
