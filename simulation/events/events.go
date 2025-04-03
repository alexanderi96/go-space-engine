// Package events provides an event system for the simulation
package events

import (
	"sync"

	"github.com/alexanderi96/go-space-engine/physics/body"
	"github.com/alexanderi96/go-space-engine/physics/collision"
)

// EventType represents the type of event
type EventType int

const (
	// BodyAdded event generated when a body is added to the world
	BodyAdded EventType = iota
	// BodyRemoved event generated when a body is removed from the world
	BodyRemoved
	// Collision event generated when two bodies collide
	Collision
	// BoundaryCollision event generated when a body collides with the world boundaries
	BoundaryCollision
	// SimulationStarted event generated when the simulation starts
	SimulationStarted
	// SimulationStopped event generated when the simulation stops
	SimulationStopped
	// SimulationStep event generated at each simulation step
	SimulationStep
)

// Event represents an event in the simulation
type Event struct {
	// Type is the event type
	Type EventType
	// Data contains the event data
	Data interface{}
}

// BodyEvent represents an event related to a body
type BodyEvent struct {
	// Body is the body involved in the event
	Body body.Body
}

// CollisionEvent represents a collision event
type CollisionEvent struct {
	// Info contains information about the collision
	Info collision.CollisionInfo
}

// BoundaryCollisionEvent represents a collision event with the world boundaries
type BoundaryCollisionEvent struct {
	// Body is the body involved in the collision
	Body body.Body
	// Boundary is the boundary with which the body collided (0=min_x, 1=max_x, 2=min_y, 3=max_y, 4=min_z, 5=max_z)
	Boundary int
}

// SimulationStepEvent represents a simulation step event
type SimulationStepEvent struct {
	// DeltaTime is the time step
	DeltaTime float64
	// Time is the total simulation time
	Time float64
}

// EventListener represents an event listener
type EventListener interface {
	// OnEvent is called when an event occurs
	OnEvent(event Event)
}

// EventSystem represents an event system
type EventSystem interface {
	// AddListener adds a listener for an event type
	AddListener(listener EventListener, eventType EventType)
	// RemoveListener removes a listener for an event type
	RemoveListener(listener EventListener, eventType EventType)
	// DispatchEvent dispatches an event to all registered listeners
	DispatchEvent(event Event)
}

// SimpleEventSystem implements a simple event system
type SimpleEventSystem struct {
	listeners map[EventType][]EventListener
	mutex     sync.RWMutex
}

// NewSimpleEventSystem creates a new simple event system
func NewSimpleEventSystem() *SimpleEventSystem {
	return &SimpleEventSystem{
		listeners: make(map[EventType][]EventListener),
	}
}

// AddListener adds a listener for an event type
func (es *SimpleEventSystem) AddListener(listener EventListener, eventType EventType) {
	es.mutex.Lock()
	defer es.mutex.Unlock()

	if _, exists := es.listeners[eventType]; !exists {
		es.listeners[eventType] = make([]EventListener, 0)
	}

	es.listeners[eventType] = append(es.listeners[eventType], listener)
}

// RemoveListener removes a listener for an event type
func (es *SimpleEventSystem) RemoveListener(listener EventListener, eventType EventType) {
	es.mutex.Lock()
	defer es.mutex.Unlock()

	if listeners, exists := es.listeners[eventType]; exists {
		for i, l := range listeners {
			if l == listener {
				// Remove the listener by swapping it with the last one and truncating the slice
				lastIndex := len(listeners) - 1
				listeners[i] = listeners[lastIndex]
				es.listeners[eventType] = listeners[:lastIndex]
				break
			}
		}
	}
}

// DispatchEvent dispatches an event to all registered listeners
func (es *SimpleEventSystem) DispatchEvent(event Event) {
	es.mutex.RLock()
	defer es.mutex.RUnlock()

	if listeners, exists := es.listeners[event.Type]; exists {
		for _, listener := range listeners {
			listener.OnEvent(event)
		}
	}
}

// EventLogger is an event listener that logs events
type EventLogger struct {
	// LogFunc is the logging function
	LogFunc func(format string, args ...interface{})
}

// NewEventLogger creates a new event logger
func NewEventLogger(logFunc func(format string, args ...interface{})) *EventLogger {
	return &EventLogger{
		LogFunc: logFunc,
	}
}

// OnEvent is called when an event occurs
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
