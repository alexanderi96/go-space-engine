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
	entity     entity.Entity
	angularAcc vector.Vector3
}

// NewBaseController creates a new base controller for the given entity
func NewBaseController(entity entity.Entity) *BaseController {
	return &BaseController{
		entity:     entity,
		angularAcc: vector.Zero3(),
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
	// Store the angular acceleration to be applied during update
	c.angularAcc = c.angularAcc.Add(torque)
}

// Update updates the controller's state
func (c *BaseController) Update(deltaTime float64) {
	// Apply angular acceleration to update angular velocity
	if c.angularAcc.LengthSquared() > 1e-10 {
		// Get current angular velocity
		currentAngVel := c.entity.GetAngularVelocity()

		// Calculate new angular velocity based on angular acceleration
		newAngVel := currentAngVel.Add(c.angularAcc.Scale(deltaTime))

		// Set the new angular velocity on the entity
		c.entity.SetAngularVelocity(newAngVel)

		// Reset angular acceleration
		c.angularAcc = vector.Zero3()
	}

	// Update the entity
	c.entity.Update(deltaTime)
}
