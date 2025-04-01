// Package events fornisce un sistema di eventi per la simulazione
package events

import (
	"sync"

	"github.com/alexanderi96/go-space-engine/physics/body"
	"github.com/alexanderi96/go-space-engine/physics/collision"
)

// EventType rappresenta il tipo di evento
type EventType int

const (
	// BodyAdded evento generato quando un corpo viene aggiunto al mondo
	BodyAdded EventType = iota
	// BodyRemoved evento generato quando un corpo viene rimosso dal mondo
	BodyRemoved
	// Collision evento generato quando due corpi collidono
	Collision
	// BoundaryCollision evento generato quando un corpo collide con i limiti del mondo
	BoundaryCollision
	// SimulationStarted evento generato quando la simulazione inizia
	SimulationStarted
	// SimulationStopped evento generato quando la simulazione si ferma
	SimulationStopped
	// SimulationStep evento generato ad ogni passo della simulazione
	SimulationStep
)

// Event rappresenta un evento nella simulazione
type Event struct {
	// Type è il tipo di evento
	Type EventType
	// Data contiene i dati dell'evento
	Data interface{}
}

// BodyEvent rappresenta un evento relativo a un corpo
type BodyEvent struct {
	// Body è il corpo coinvolto nell'evento
	Body body.Body
}

// CollisionEvent rappresenta un evento di collisione
type CollisionEvent struct {
	// Info contiene le informazioni sulla collisione
	Info collision.CollisionInfo
}

// BoundaryCollisionEvent rappresenta un evento di collisione con i limiti del mondo
type BoundaryCollisionEvent struct {
	// Body è il corpo coinvolto nella collisione
	Body body.Body
	// Boundary è il limite con cui il corpo ha colliso (0=min_x, 1=max_x, 2=min_y, 3=max_y, 4=min_z, 5=max_z)
	Boundary int
}

// SimulationStepEvent rappresenta un evento di passo della simulazione
type SimulationStepEvent struct {
	// DeltaTime è il passo temporale
	DeltaTime float64
	// Time è il tempo totale della simulazione
	Time float64
}

// EventListener rappresenta un ascoltatore di eventi
type EventListener interface {
	// OnEvent viene chiamato quando si verifica un evento
	OnEvent(event Event)
}

// EventSystem rappresenta un sistema di eventi
type EventSystem interface {
	// AddListener aggiunge un ascoltatore per un tipo di evento
	AddListener(listener EventListener, eventType EventType)
	// RemoveListener rimuove un ascoltatore per un tipo di evento
	RemoveListener(listener EventListener, eventType EventType)
	// DispatchEvent invia un evento a tutti gli ascoltatori registrati
	DispatchEvent(event Event)
}

// SimpleEventSystem implementa un sistema di eventi semplice
type SimpleEventSystem struct {
	listeners map[EventType][]EventListener
	mutex     sync.RWMutex
}

// NewSimpleEventSystem crea un nuovo sistema di eventi semplice
func NewSimpleEventSystem() *SimpleEventSystem {
	return &SimpleEventSystem{
		listeners: make(map[EventType][]EventListener),
	}
}

// AddListener aggiunge un ascoltatore per un tipo di evento
func (es *SimpleEventSystem) AddListener(listener EventListener, eventType EventType) {
	es.mutex.Lock()
	defer es.mutex.Unlock()

	if _, exists := es.listeners[eventType]; !exists {
		es.listeners[eventType] = make([]EventListener, 0)
	}

	es.listeners[eventType] = append(es.listeners[eventType], listener)
}

// RemoveListener rimuove un ascoltatore per un tipo di evento
func (es *SimpleEventSystem) RemoveListener(listener EventListener, eventType EventType) {
	es.mutex.Lock()
	defer es.mutex.Unlock()

	if listeners, exists := es.listeners[eventType]; exists {
		for i, l := range listeners {
			if l == listener {
				// Rimuovi l'ascoltatore scambiandolo con l'ultimo e troncando la slice
				lastIndex := len(listeners) - 1
				listeners[i] = listeners[lastIndex]
				es.listeners[eventType] = listeners[:lastIndex]
				break
			}
		}
	}
}

// DispatchEvent invia un evento a tutti gli ascoltatori registrati
func (es *SimpleEventSystem) DispatchEvent(event Event) {
	es.mutex.RLock()
	defer es.mutex.RUnlock()

	if listeners, exists := es.listeners[event.Type]; exists {
		for _, listener := range listeners {
			listener.OnEvent(event)
		}
	}
}

// EventLogger è un ascoltatore di eventi che registra gli eventi
type EventLogger struct {
	// LogFunc è la funzione di logging
	LogFunc func(format string, args ...interface{})
}

// NewEventLogger crea un nuovo logger di eventi
func NewEventLogger(logFunc func(format string, args ...interface{})) *EventLogger {
	return &EventLogger{
		LogFunc: logFunc,
	}
}

// OnEvent viene chiamato quando si verifica un evento
func (el *EventLogger) OnEvent(event Event) {
	switch event.Type {
	case BodyAdded:
		if bodyEvent, ok := event.Data.(BodyEvent); ok {
			el.LogFunc("Body added: %v", bodyEvent.Body.ID())
		}
	case BodyRemoved:
		if bodyEvent, ok := event.Data.(BodyEvent); ok {
			el.LogFunc("Body removed: %v", bodyEvent.Body.ID())
		}
	case Collision:
		if collisionEvent, ok := event.Data.(CollisionEvent); ok {
			el.LogFunc("Collision between %v and %v", collisionEvent.Info.BodyA.ID(), collisionEvent.Info.BodyB.ID())
		}
	case BoundaryCollision:
		if boundaryEvent, ok := event.Data.(BoundaryCollisionEvent); ok {
			boundaries := []string{"min_x", "max_x", "min_y", "max_y", "min_z", "max_z"}
			boundary := "unknown"
			if boundaryEvent.Boundary >= 0 && boundaryEvent.Boundary < len(boundaries) {
				boundary = boundaries[boundaryEvent.Boundary]
			}
			el.LogFunc("Boundary collision: %v with %s", boundaryEvent.Body.ID(), boundary)
		}
	case SimulationStarted:
		el.LogFunc("Simulation started")
	case SimulationStopped:
		el.LogFunc("Simulation stopped")
	case SimulationStep:
		if stepEvent, ok := event.Data.(SimulationStepEvent); ok {
			el.LogFunc("Simulation step: dt=%f, t=%f", stepEvent.DeltaTime, stepEvent.Time)
		}
	default:
		el.LogFunc("Unknown event: %v", event.Type)
	}
}
