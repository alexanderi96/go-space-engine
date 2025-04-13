package controller

import (
	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/entity"
)

// Controller defines the interface for entity controllers
type Controller interface {
	// GetEntity returns the controlled entity
	GetEntity() entity.Entity

	// ApplyForce applies a force to the entity at the given point
	ApplyForce(force vector.Vector3)

	// ApplyTorque applies a torque to the entity
	ApplyTorque(torque vector.Vector3)

	// Update updates the controller's state
	Update(deltaTime float64)
}

// BaseController provides a basic implementation of the Controller interface
type BaseController struct {
	entity entity.Entity
}

// NewBaseController creates a new base controller for the given entity
func NewBaseController(entity entity.Entity) *BaseController {
	return &BaseController{
		entity: entity,
	}
}

// GetEntity returns the controlled entity
func (c *BaseController) GetEntity() entity.Entity {
	return c.entity
}

// ApplyForce applies a force to the entity
func (c *BaseController) ApplyForce(force vector.Vector3) {
	body := c.entity.GetBody()
	if body != nil {
		body.ApplyForce(force)
	}
}

// ApplyTorque applies a torque to the entity
func (c *BaseController) ApplyTorque(torque vector.Vector3) {
	// Apply torque directly to the physical body
	body := c.entity.GetBody()
	if body != nil {
		body.ApplyTorque(torque)
	}
}

// Update updates the controller's state
func (c *BaseController) Update(deltaTime float64) {
	// Update the entity (which updates the physical body)
	c.entity.Update(deltaTime)
}
