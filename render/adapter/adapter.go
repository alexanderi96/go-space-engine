// Package adapter fornisce interfacce per il rendering
package adapter

import (
	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/physics/body"
	"github.com/alexanderi96/go-space-engine/physics/space"
	"github.com/alexanderi96/go-space-engine/simulation/world"
)

// Color rappresenta un colore RGBA
type Color struct {
	R, G, B, A float64 // Componenti del colore (0-1)
}

// NewColor crea un nuovo colore
func NewColor(r, g, b, a float64) Color {
	return Color{
		R: r,
		G: g,
		B: b,
		A: a,
	}
}

// Renderer rappresenta un'interfaccia per il rendering
type Renderer interface {
	// Initialize inizializza il renderer
	Initialize() error
	// Shutdown chiude il renderer
	Shutdown() error

	// BeginFrame inizia un nuovo frame
	BeginFrame()
	// EndFrame termina il frame corrente
	EndFrame()

	// RenderBody renderizza un corpo
	RenderBody(b body.Body)
	// RenderBodies renderizza tutti i corpi
	RenderBodies(bodies []body.Body)

	// RenderAABB renderizza un AABB
	RenderAABB(aabb *space.AABB, color Color)
	// RenderOctree renderizza un octree
	RenderOctree(octree *space.Octree, maxDepth int)

	// RenderLine renderizza una linea
	RenderLine(start, end vector.Vector3, color Color)
	// RenderSphere renderizza una sfera
	RenderSphere(center vector.Vector3, radius float64, color Color)

	// SetCamera imposta la posizione e l'orientamento della camera
	SetCamera(position, target, up vector.Vector3)
	// SetCameraFOV imposta il campo visivo della camera
	SetCameraFOV(fov float64)

	// SetBackgroundColor imposta il colore di sfondo
	SetBackgroundColor(color Color)

	// GetWidth restituisce la larghezza della finestra di rendering
	GetWidth() int
	// GetHeight restituisce l'altezza della finestra di rendering
	GetHeight() int

	// IsRunning restituisce true se il renderer è in esecuzione
	IsRunning() bool

	// ProcessEvents processa gli eventi del renderer
	ProcessEvents()
}

// RenderAdapter rappresenta un adattatore per il rendering
type RenderAdapter interface {
	// GetRenderer restituisce il renderer
	GetRenderer() Renderer

	// RenderWorld renderizza il mondo
	RenderWorld(w world.World)

	// SetDebugMode imposta la modalità di debug
	SetDebugMode(debug bool)
	// IsDebugMode restituisce true se la modalità di debug è attiva
	IsDebugMode() bool

	// SetRenderOctree imposta se renderizzare l'octree
	SetRenderOctree(render bool)
	// IsRenderOctree restituisce true se l'octree viene renderizzato
	IsRenderOctree() bool

	// SetRenderBoundingBoxes imposta se renderizzare i bounding box
	SetRenderBoundingBoxes(render bool)
	// IsRenderBoundingBoxes restituisce true se i bounding box vengono renderizzati
	IsRenderBoundingBoxes() bool

	// SetRenderVelocities imposta se renderizzare i vettori velocità
	SetRenderVelocities(render bool)
	// IsRenderVelocities restituisce true se i vettori velocità vengono renderizzati
	IsRenderVelocities() bool

	// SetRenderAccelerations imposta se renderizzare i vettori accelerazione
	SetRenderAccelerations(render bool)
	// IsRenderAccelerations restituisce true se i vettori accelerazione vengono renderizzati
	IsRenderAccelerations() bool

	// SetRenderForces imposta se renderizzare i vettori forza
	SetRenderForces(render bool)
	// IsRenderForces restituisce true se i vettori forza vengono renderizzati
	IsRenderForces() bool
}

// BaseRenderAdapter implementa un adattatore di base per il rendering
type BaseRenderAdapter struct {
	renderer            Renderer
	debugMode           bool
	renderOctree        bool
	renderBoundingBoxes bool
	renderVelocities    bool
	renderAccelerations bool
	renderForces        bool
}

// NewBaseRenderAdapter crea un nuovo adattatore di base per il rendering
func NewBaseRenderAdapter(renderer Renderer) *BaseRenderAdapter {
	return &BaseRenderAdapter{
		renderer:            renderer,
		debugMode:           false,
		renderOctree:        false,
		renderBoundingBoxes: false,
		renderVelocities:    false,
		renderAccelerations: false,
		renderForces:        false,
	}
}

// GetRenderer restituisce il renderer
func (ra *BaseRenderAdapter) GetRenderer() Renderer {
	return ra.renderer
}

// RenderWorld renderizza il mondo
func (ra *BaseRenderAdapter) RenderWorld(w world.World) {
	// Inizia un nuovo frame
	ra.renderer.BeginFrame()

	// Renderizza tutti i corpi
	ra.renderer.RenderBodies(w.GetBodies())

	// Renderizza i limiti del mondo
	bounds := w.GetBounds()
	ra.renderer.RenderAABB(bounds, NewColor(0.5, 0.5, 0.5, 0.5))

	// Renderizza l'octree se richiesto
	if ra.renderOctree {
		if octree, ok := w.GetSpatialStructure().(*space.Octree); ok {
			ra.renderer.RenderOctree(octree, 8)
		}
	}

	// Renderizza i bounding box se richiesto
	if ra.renderBoundingBoxes {
		for _, b := range w.GetBodies() {
			position := b.Position()
			radius := b.Radius().Value()
			min := vector.NewVector3(position.X()-radius, position.Y()-radius, position.Z()-radius)
			max := vector.NewVector3(position.X()+radius, position.Y()+radius, position.Z()+radius)
			aabb := space.NewAABB(min, max)
			ra.renderer.RenderAABB(aabb, NewColor(0, 1, 0, 0.5))
		}
	}

	// Renderizza i vettori velocità se richiesto
	if ra.renderVelocities {
		for _, b := range w.GetBodies() {
			position := b.Position()
			velocity := b.Velocity()
			if velocity.Length() > 0.001 {
				end := position.Add(velocity)
				ra.renderer.RenderLine(position, end, NewColor(0, 0, 1, 1))
			}
		}
	}

	// Renderizza i vettori accelerazione se richiesto
	if ra.renderAccelerations {
		for _, b := range w.GetBodies() {
			position := b.Position()
			acceleration := b.Acceleration()
			if acceleration.Length() > 0.001 {
				end := position.Add(acceleration)
				ra.renderer.RenderLine(position, end, NewColor(1, 0, 0, 1))
			}
		}
	}

	// Termina il frame
	ra.renderer.EndFrame()
}

// SetDebugMode imposta la modalità di debug
func (ra *BaseRenderAdapter) SetDebugMode(debug bool) {
	ra.debugMode = debug
}

// IsDebugMode restituisce true se la modalità di debug è attiva
func (ra *BaseRenderAdapter) IsDebugMode() bool {
	return ra.debugMode
}

// SetRenderOctree imposta se renderizzare l'octree
func (ra *BaseRenderAdapter) SetRenderOctree(render bool) {
	ra.renderOctree = render
}

// IsRenderOctree restituisce true se l'octree viene renderizzato
func (ra *BaseRenderAdapter) IsRenderOctree() bool {
	return ra.renderOctree
}

// SetRenderBoundingBoxes imposta se renderizzare i bounding box
func (ra *BaseRenderAdapter) SetRenderBoundingBoxes(render bool) {
	ra.renderBoundingBoxes = render
}

// IsRenderBoundingBoxes restituisce true se i bounding box vengono renderizzati
func (ra *BaseRenderAdapter) IsRenderBoundingBoxes() bool {
	return ra.renderBoundingBoxes
}

// SetRenderVelocities imposta se renderizzare i vettori velocità
func (ra *BaseRenderAdapter) SetRenderVelocities(render bool) {
	ra.renderVelocities = render
}

// IsRenderVelocities restituisce true se i vettori velocità vengono renderizzati
func (ra *BaseRenderAdapter) IsRenderVelocities() bool {
	return ra.renderVelocities
}

// SetRenderAccelerations imposta se renderizzare i vettori accelerazione
func (ra *BaseRenderAdapter) SetRenderAccelerations(render bool) {
	ra.renderAccelerations = render
}

// IsRenderAccelerations restituisce true se i vettori accelerazione vengono renderizzati
func (ra *BaseRenderAdapter) IsRenderAccelerations() bool {
	return ra.renderAccelerations
}

// SetRenderForces imposta se renderizzare i vettori forza
func (ra *BaseRenderAdapter) SetRenderForces(render bool) {
	ra.renderForces = render
}

// IsRenderForces restituisce true se i vettori forza vengono renderizzati
func (ra *BaseRenderAdapter) IsRenderForces() bool {
	return ra.renderForces
}
