// Package g3n fornisce un'implementazione dell'interfaccia RenderAdapter utilizzando G3N
package g3n

import (
	"time"

	"github.com/alexanderi96/go-space-engine/physics/body"
	"github.com/alexanderi96/go-space-engine/render/adapter"
	"github.com/alexanderi96/go-space-engine/simulation/events"
	"github.com/alexanderi96/go-space-engine/simulation/world"
)

// G3NAdapter è un adattatore per il rendering con G3N
type G3NAdapter struct {
	*adapter.BaseRenderAdapter
	eventListener *G3NEventListener
}

// NewG3NAdapter crea un nuovo adattatore G3N
func NewG3NAdapter() *G3NAdapter {
	renderer := NewG3NRenderer()
	baseAdapter := adapter.NewBaseRenderAdapter(renderer)

	adapter := &G3NAdapter{
		BaseRenderAdapter: baseAdapter,
	}

	// Crea il listener di eventi
	adapter.eventListener = NewG3NEventListener(adapter)

	return adapter
}

// GetG3NRenderer restituisce il renderer G3N sottostante
func (ga *G3NAdapter) GetG3NRenderer() *G3NRenderer {
	return ga.GetRenderer().(*G3NRenderer)
}

// GetEventListener restituisce il listener di eventi
func (ga *G3NAdapter) GetEventListener() *G3NEventListener {
	return ga.eventListener
}

// RenderWorld renderizza il mondo
func (ga *G3NAdapter) RenderWorld(w world.World) {
	// Chiama il metodo base
	ga.BaseRenderAdapter.RenderWorld(w)
}

// Run avvia il loop di rendering
func (ga *G3NAdapter) Run(updateFunc func(deltaTime time.Duration)) {
	// Delega al renderer G3N
	ga.GetG3NRenderer().Run(updateFunc)
}

// SetDebugMode imposta la modalità di debug
func (ga *G3NAdapter) SetDebugMode(debug bool) {
	ga.BaseRenderAdapter.SetDebugMode(debug)
}

// SetRenderVelocities imposta se renderizzare i vettori velocità
func (ga *G3NAdapter) SetRenderVelocities(render bool) {
	ga.BaseRenderAdapter.SetRenderVelocities(render)
}

// SetRenderAccelerations imposta se renderizzare i vettori accelerazione
func (ga *G3NAdapter) SetRenderAccelerations(render bool) {
	ga.BaseRenderAdapter.SetRenderAccelerations(render)
}

// SetRenderBoundingBoxes imposta se renderizzare i bounding box
func (ga *G3NAdapter) SetRenderBoundingBoxes(render bool) {
	ga.BaseRenderAdapter.SetRenderBoundingBoxes(render)
}

// SetRenderOctree imposta se renderizzare l'octree
func (ga *G3NAdapter) SetRenderOctree(render bool) {
	ga.BaseRenderAdapter.SetRenderOctree(render)
}

// SetRenderForces imposta se renderizzare i vettori forza
func (ga *G3NAdapter) SetRenderForces(render bool) {
	ga.BaseRenderAdapter.SetRenderForces(render)
}

// G3NEventListener ascolta gli eventi della simulazione e aggiorna il renderer G3N
type G3NEventListener struct {
	adapter *G3NAdapter
}

// NewG3NEventListener crea un nuovo listener di eventi per G3N
func NewG3NEventListener(adapter *G3NAdapter) *G3NEventListener {
	return &G3NEventListener{
		adapter: adapter,
	}
}

// OnEvent gestisce gli eventi della simulazione
func (l *G3NEventListener) OnEvent(event events.Event) {
	renderer := l.adapter.GetG3NRenderer()

	switch event.Type {
	case events.BodyAdded:
		// Quando un corpo viene aggiunto, lo renderizziamo
		if bodyEvent, ok := event.Data.(events.BodyEvent); ok {
			renderer.RenderBody(bodyEvent.Body)
		}

	case events.BodyRemoved:
		// Quando un corpo viene rimosso, rimuoviamo il suo nodo grafico
		if bodyEvent, ok := event.Data.(events.BodyEvent); ok {
			renderer.nodeMutex.Lock()
			if node, exists := renderer.bodyNodes[bodyEvent.Body.ID()]; exists {
				renderer.scene.Remove(node)
				delete(renderer.bodyNodes, bodyEvent.Body.ID())
			}
			renderer.nodeMutex.Unlock()
		}

	case events.Collision:
		// Quando avviene una collisione, possiamo visualizzare un effetto
		if collisionEvent, ok := event.Data.(events.CollisionEvent); ok {
			// Visualizza un effetto di collisione (ad esempio, una sfera rossa temporanea)
			point := collisionEvent.Info.Point
			renderer.RenderSphere(point, 0.2, adapter.NewColor(1, 0, 0, 0.7))
		}

	case events.BoundaryCollision:
		// Quando un corpo collide con i limiti del mondo, possiamo visualizzare un effetto
		if boundaryEvent, ok := event.Data.(events.BoundaryCollisionEvent); ok {
			// Visualizza un effetto di collisione con i limiti
			pos := boundaryEvent.Body.Position()
			renderer.RenderSphere(pos, 0.2, adapter.NewColor(1, 1, 0, 0.7))
		}

	case events.SimulationStep:
		// Ad ogni passo della simulazione, aggiorniamo la visualizzazione
		// Questo è già gestito dal metodo RenderWorld
	}
}

// G3NWorldObserver osserva il mondo fisico e genera eventi
type G3NWorldObserver struct {
	world       world.World
	eventSystem events.EventSystem
	previousIDs map[body.ID]bool
}

// NewG3NWorldObserver crea un nuovo observer per il mondo fisico
func NewG3NWorldObserver(w world.World, es events.EventSystem) *G3NWorldObserver {
	return &G3NWorldObserver{
		world:       w,
		eventSystem: es,
		previousIDs: make(map[body.ID]bool),
	}
}

// Update aggiorna l'observer
func (o *G3NWorldObserver) Update() {
	// Ottieni tutti i corpi nel mondo
	bodies := o.world.GetBodies()

	// Crea un set di ID dei corpi correnti
	currentIDs := make(map[body.ID]bool)
	for _, b := range bodies {
		currentIDs[b.ID()] = true

		// Se il corpo non era presente nel passo precedente, genera un evento BodyAdded
		if !o.previousIDs[b.ID()] {
			o.eventSystem.DispatchEvent(events.Event{
				Type: events.BodyAdded,
				Data: events.BodyEvent{Body: b},
			})
		}
	}

	// Controlla se ci sono corpi che sono stati rimossi
	for id := range o.previousIDs {
		if !currentIDs[id] {
			// Il corpo è stato rimosso, genera un evento BodyRemoved
			// Nota: non abbiamo il corpo, solo il suo ID
			// Per ora, non generiamo l'evento
		}
	}

	// Aggiorna il set di ID precedenti
	o.previousIDs = currentIDs
}
