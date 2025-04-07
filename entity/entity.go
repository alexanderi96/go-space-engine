package entity

import (
	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/physics/body"
)

// Entity represents a controllable entity in the simulation
type Entity interface {
	// GetID returns the unique identifier of the entity
	GetID() string

	// GetBody returns the underlying physical body
	GetBody() body.Body

	// GetPosition returns the current position of the entity
	GetPosition() vector.Vector3

	// GetRotation returns the current rotation of the entity
	GetRotation() vector.Vector3

	// GetVelocity returns the current velocity of the entity
	GetVelocity() vector.Vector3

	// GetAngularVelocity returns the current angular velocity of the entity
	GetAngularVelocity() vector.Vector3

	// SetAngularVelocity sets the angular velocity of the entity
	SetAngularVelocity(angVel vector.Vector3)

	// Update is called every simulation step to update the entity's state
	Update(deltaTime float64)
}

// BaseEntity provides a basic implementation of the Entity interface
type BaseEntity struct {
	id         string
	body       body.Body
	rotation   vector.Vector3
	angularVel vector.Vector3
}

// NewBaseEntity creates a new base entity with the given parameters
func NewBaseEntity(id string, body body.Body) *BaseEntity {
	return &BaseEntity{
		id:         id,
		body:       body,
		rotation:   vector.Zero3(),
		angularVel: vector.Zero3(),
	}
}

// GetID returns the unique identifier of the entity
func (e *BaseEntity) GetID() string {
	return e.id
}

// GetBody returns the underlying physical body
func (e *BaseEntity) GetBody() body.Body {
	return e.body
}

// GetPosition returns the current position of the entity
func (e *BaseEntity) GetPosition() vector.Vector3 {
	return e.body.Position()
}

// GetRotation returns the current rotation of the entity
func (e *BaseEntity) GetRotation() vector.Vector3 {
	return e.rotation
}

// GetVelocity returns the current velocity of the entity
func (e *BaseEntity) GetVelocity() vector.Vector3 {
	return e.body.Velocity()
}

// GetAngularVelocity returns the current angular velocity of the entity
func (e *BaseEntity) GetAngularVelocity() vector.Vector3 {
	return e.angularVel
}

// SetAngularVelocity sets the angular velocity of the entity
func (e *BaseEntity) SetAngularVelocity(angVel vector.Vector3) {
	e.angularVel = angVel
}

// Update updates the entity's state based on its physical body
func (e *BaseEntity) Update(deltaTime float64) {
	// Update rotation based on angular velocity
	e.rotation = e.rotation.Add(e.angularVel.Scale(deltaTime))

	// Update the physical body
	e.body.Update(deltaTime)
}
