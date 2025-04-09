package spacecraft

import (
	"math"

	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/entity"
	"github.com/alexanderi96/go-space-engine/entity/controller"
)

// SpacecraftController implements a controller for spacecraft movement
type SpacecraftController struct {
	*controller.BaseController

	// Maximum thrust force that can be applied
	maxThrust float64

	// Maximum torque that can be applied for rotation
	maxTorque float64

	// Current thrust level (0.0 to 1.0)
	thrustLevel float64

	// Maximum angular velocity (degrees per second)
	maxAngularVelocity float64

	// PID controller parameters for rotation
	rotationPID struct {
		kP float64 // Proportional gain
		kI float64 // Integral gain
		kD float64 // Derivative gain

		// Error tracking
		lastError  vector.Vector3
		errorSum   vector.Vector3
		errorDelta vector.Vector3

		// Target rotation
		targetRotation vector.Vector3
	}

	// Damping factor for angular velocity (0.0 to 1.0)
	angularDamping float64
}

// NewSpacecraftController creates a new spacecraft controller
func NewSpacecraftController(entity entity.Entity, maxThrust, maxTorque float64) *SpacecraftController {
	controller := &SpacecraftController{
		BaseController:     controller.NewBaseController(entity),
		maxThrust:          maxThrust,
		maxTorque:          maxTorque,
		thrustLevel:        0.0,
		maxAngularVelocity: 45.0, // 45 degrees per second default
		angularDamping:     0.2,  // 20% damping by default
	}

	// Initialize PID controller with default values
	controller.rotationPID.kP = 2.0
	controller.rotationPID.kI = 0.1
	controller.rotationPID.kD = 0.5

	// Initialize vector fields to avoid nil pointer dereference
	controller.rotationPID.lastError = vector.Zero3()
	controller.rotationPID.errorSum = vector.Zero3()
	controller.rotationPID.errorDelta = vector.Zero3()
	controller.rotationPID.targetRotation = vector.Zero3()

	return controller
}

// SetMaxAngularVelocity sets the maximum angular velocity in degrees per second
func (c *SpacecraftController) SetMaxAngularVelocity(maxVelocity float64) {
	if maxVelocity < 0 {
		maxVelocity = 0
	}
	c.maxAngularVelocity = maxVelocity
}

// SetAngularDamping sets the damping factor for angular velocity
func (c *SpacecraftController) SetAngularDamping(damping float64) {
	if damping < 0.0 {
		damping = 0.0
	} else if damping > 1.0 {
		damping = 1.0
	}
	c.angularDamping = damping
}

// SetRotationPIDGains sets the PID controller gains for rotation
func (c *SpacecraftController) SetRotationPIDGains(kP, kI, kD float64) {
	c.rotationPID.kP = kP
	c.rotationPID.kI = kI
	c.rotationPID.kD = kD
}

// SetTargetRotation sets the target rotation for the PID controller
func (c *SpacecraftController) SetTargetRotation(targetRotation vector.Vector3) {
	c.rotationPID.targetRotation = targetRotation
}

// SetThrustLevel sets the current thrust level (0.0 to 1.0)
func (c *SpacecraftController) SetThrustLevel(level float64) {
	// Clamp thrust level between 0 and 1
	if level < 0.0 {
		level = 0.0
	} else if level > 1.0 {
		level = 1.0
	}
	c.thrustLevel = level
}

// NormalizeAngle normalizes an angle to the range [-180, 180]
func NormalizeAngle(angle float64) float64 {
	// Normalize to [0, 360)
	angle = math.Mod(angle, 360.0)
	if angle < 0 {
		angle += 360.0
	}

	// Convert to [-180, 180)
	if angle >= 180.0 {
		angle -= 360.0
	}

	return angle
}

// ApplyMainThrust applies the main thrust in the forward direction
func (c *SpacecraftController) ApplyMainThrust() {
	if c.thrustLevel <= 0.0 {
		return
	}

	// Get the entity's rotation
	rotation := c.GetEntity().GetRotation()

	// Normalize rotation angles
	normalizedRotation := vector.NewVector3(
		NormalizeAngle(rotation.X()),
		NormalizeAngle(rotation.Y()),
		NormalizeAngle(rotation.Z()),
	)

	// Convert rotation angles to radians
	rotX := normalizedRotation.X() * (math.Pi / 180.0)
	rotY := normalizedRotation.Y() * (math.Pi / 180.0)
	rotZ := normalizedRotation.Z() * (math.Pi / 180.0)

	// Start with forward vector (0,0,1)
	forwardVector := vector.NewVector3(0, 0, 1)

	// Apply rotations in order: Z (roll), X (pitch), Y (yaw)
	// This is a simplified rotation calculation and might not be 100% accurate for all cases

	// Roll (Z-axis rotation)
	cosZ := math.Cos(rotZ)
	sinZ := math.Sin(rotZ)
	tempX := forwardVector.X()*cosZ - forwardVector.Y()*sinZ
	tempY := forwardVector.X()*sinZ + forwardVector.Y()*cosZ
	forwardVector = vector.NewVector3(tempX, tempY, forwardVector.Z())

	// Pitch (X-axis rotation)
	cosX := math.Cos(rotX)
	sinX := math.Sin(rotX)
	tempY = forwardVector.Y()*cosX - forwardVector.Z()*sinX
	tempZ := forwardVector.Y()*sinX + forwardVector.Z()*cosX
	forwardVector = vector.NewVector3(forwardVector.X(), tempY, tempZ)

	// Yaw (Y-axis rotation)
	cosY := math.Cos(rotY)
	sinY := math.Sin(rotY)
	tempX = forwardVector.X()*cosY + forwardVector.Z()*sinY
	tempZ = -forwardVector.X()*sinY + forwardVector.Z()*cosY
	forwardVector = vector.NewVector3(tempX, forwardVector.Y(), tempZ)

	// Normalize the direction vector
	thrustDir := forwardVector.Normalize()

	// Apply the thrust force
	thrustForce := thrustDir.Scale(c.maxThrust * c.thrustLevel)
	c.ApplyForce(thrustForce)
}

// ApplyRotation applies rotation around the specified axis
func (c *SpacecraftController) ApplyRotation(axis vector.Vector3, amount float64) {
	// Clamp rotation amount between -1 and 1
	if amount < -1.0 {
		amount = -1.0
	} else if amount > 1.0 {
		amount = 1.0
	}

	// Calculate and apply the torque with PID control
	torque := axis.Scale(c.maxTorque * amount)

	// Apply angular velocity limits
	angularVel := c.GetEntity().GetAngularVelocity()
	angularSpeed := angularVel.Length()

	// If we're exceeding max angular velocity, reduce the torque
	if angularSpeed > c.maxAngularVelocity {
		// Scale down torque to prevent exceeding max angular velocity
		scaleFactor := c.maxAngularVelocity / angularSpeed
		torque = torque.Scale(scaleFactor)
	}

	c.ApplyTorque(torque)
}

// UpdatePIDController updates the PID controller for rotation
func (c *SpacecraftController) UpdatePIDController(deltaTime float64) {
	// Get current rotation
	currentRotation := c.GetEntity().GetRotation()

	// Normalize current rotation
	normalizedRotation := vector.NewVector3(
		NormalizeAngle(currentRotation.X()),
		NormalizeAngle(currentRotation.Y()),
		NormalizeAngle(currentRotation.Z()),
	)

	// Calculate error (difference between target and current rotation)
	error := c.rotationPID.targetRotation.Sub(normalizedRotation)

	// Normalize error angles to ensure shortest path rotation
	error = vector.NewVector3(
		NormalizeAngle(error.X()),
		NormalizeAngle(error.Y()),
		NormalizeAngle(error.Z()),
	)

	// Calculate error delta (derivative term)
	c.rotationPID.errorDelta = error.Sub(c.rotationPID.lastError).Scale(1.0 / deltaTime)

	// Update error sum (integral term) with anti-windup
	c.rotationPID.errorSum = c.rotationPID.errorSum.Add(error.Scale(deltaTime))

	// Anti-windup: limit the integral term
	maxIntegral := 10.0

	// Create new error sum with clamped values
	x := c.rotationPID.errorSum.X()
	y := c.rotationPID.errorSum.Y()
	z := c.rotationPID.errorSum.Z()

	// Clamp X component
	if x > maxIntegral {
		x = maxIntegral
	} else if x < -maxIntegral {
		x = -maxIntegral
	}

	// Clamp Y component
	if y > maxIntegral {
		y = maxIntegral
	} else if y < -maxIntegral {
		y = -maxIntegral
	}

	// Clamp Z component
	if z > maxIntegral {
		z = maxIntegral
	} else if z < -maxIntegral {
		z = -maxIntegral
	}

	// Create new vector with clamped values
	c.rotationPID.errorSum = vector.NewVector3(x, y, z)

	// Calculate PID output
	pTerm := error.Scale(c.rotationPID.kP)
	iTerm := c.rotationPID.errorSum.Scale(c.rotationPID.kI)
	dTerm := c.rotationPID.errorDelta.Scale(c.rotationPID.kD)

	// Combine terms
	pidOutput := pTerm.Add(iTerm).Add(dTerm)

	// Apply PID output as torque if there's a significant error
	if error.Length() > 0.1 {
		// Normalize the output to get direction
		direction := pidOutput.Normalize()

		// Scale by magnitude and max torque
		magnitude := math.Min(pidOutput.Length(), 1.0)
		torque := direction.Scale(c.maxTorque * magnitude)

		c.ApplyTorque(torque)
	}

	// Store current error for next iteration
	c.rotationPID.lastError = error
}

// ApplyAngularDamping applies damping to angular velocity
func (c *SpacecraftController) ApplyAngularDamping() {
	if c.angularDamping <= 0.0 {
		return
	}

	// Get current angular velocity
	angularVel := c.GetEntity().GetAngularVelocity()

	// Apply damping force (opposite to current angular velocity)
	dampingTorque := angularVel.Scale(-c.angularDamping)
	c.ApplyTorque(dampingTorque)
}

// Update updates the spacecraft controller state
func (c *SpacecraftController) Update(deltaTime float64) {
	// Update PID controller
	c.UpdatePIDController(deltaTime)

	// Apply angular damping
	c.ApplyAngularDamping()

	// Apply main thrust
	c.ApplyMainThrust()

	// Update base controller
	c.BaseController.Update(deltaTime)
}
