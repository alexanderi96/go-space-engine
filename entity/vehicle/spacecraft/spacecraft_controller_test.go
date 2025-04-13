package spacecraft

import (
	"math"
	"testing"

	"github.com/alexanderi96/go-space-engine/core/units"
	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/physics/body"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMaterial is a mock implementation of the Material interface
type MockMaterial struct {
	mock.Mock
}

func (m *MockMaterial) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockMaterial) Density() units.Quantity {
	args := m.Called()
	return args.Get(0).(units.Quantity)
}

func (m *MockMaterial) SpecificHeat() units.Quantity {
	args := m.Called()
	return args.Get(0).(units.Quantity)
}

func (m *MockMaterial) ThermalConductivity() units.Quantity {
	args := m.Called()
	return args.Get(0).(units.Quantity)
}

func (m *MockMaterial) Emissivity() float64 {
	args := m.Called()
	return args.Get(0).(float64)
}

func (m *MockMaterial) Elasticity() float64 {
	args := m.Called()
	return args.Get(0).(float64)
}

// MockBody is a mock implementation of the Body interface
type MockBody struct {
	mock.Mock
}

func (m *MockBody) ID() uuid.UUID {
	args := m.Called()
	return args.Get(0).(uuid.UUID)
}

func (m *MockBody) Position() vector.Vector3 {
	args := m.Called()
	return args.Get(0).(vector.Vector3)
}

func (m *MockBody) SetPosition(pos vector.Vector3) {
	m.Called(pos)
}

func (m *MockBody) Velocity() vector.Vector3 {
	args := m.Called()
	return args.Get(0).(vector.Vector3)
}

func (m *MockBody) SetVelocity(vel vector.Vector3) {
	m.Called(vel)
}

func (m *MockBody) Acceleration() vector.Vector3 {
	args := m.Called()
	return args.Get(0).(vector.Vector3)
}

func (m *MockBody) SetAcceleration(acc vector.Vector3) {
	m.Called(acc)
}

func (m *MockBody) Rotation() vector.Vector3 {
	args := m.Called()
	return args.Get(0).(vector.Vector3)
}

func (m *MockBody) SetRotation(rot vector.Vector3) {
	m.Called(rot)
}

func (m *MockBody) AngularVelocity() vector.Vector3 {
	args := m.Called()
	return args.Get(0).(vector.Vector3)
}

func (m *MockBody) SetAngularVelocity(angVel vector.Vector3) {
	m.Called(angVel)
}

func (m *MockBody) Mass() units.Quantity {
	args := m.Called()
	return args.Get(0).(units.Quantity)
}

func (m *MockBody) SetMass(mass units.Quantity) {
	m.Called(mass)
}

func (m *MockBody) Radius() units.Quantity {
	args := m.Called()
	return args.Get(0).(units.Quantity)
}

func (m *MockBody) SetRadius(radius units.Quantity) {
	m.Called(radius)
}

func (m *MockBody) Material() body.Material {
	args := m.Called()
	return args.Get(0).(body.Material)
}

func (m *MockBody) SetMaterial(mat body.Material) {
	m.Called(mat)
}

func (m *MockBody) ApplyForce(force vector.Vector3) {
	m.Called(force)
}

func (m *MockBody) ApplyTorque(torque vector.Vector3) {
	m.Called(torque)
}

func (m *MockBody) Update(dt float64) {
	m.Called(dt)
}

func (m *MockBody) Temperature() units.Quantity {
	args := m.Called()
	return args.Get(0).(units.Quantity)
}

func (m *MockBody) SetTemperature(temp units.Quantity) {
	m.Called(temp)
}

func (m *MockBody) AddHeat(heat units.Quantity) {
	m.Called(heat)
}

func (m *MockBody) IsStatic() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockBody) SetStatic(static bool) {
	m.Called(static)
}

// MockEntity is a mock implementation of the Entity interface
type MockEntity struct {
	mock.Mock
}

func (m *MockEntity) GetID() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockEntity) GetBody() body.Body {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(body.Body)
}

func (m *MockEntity) GetPosition() vector.Vector3 {
	args := m.Called()
	return args.Get(0).(vector.Vector3)
}

func (m *MockEntity) GetRotation() vector.Vector3 {
	args := m.Called()
	return args.Get(0).(vector.Vector3)
}

func (m *MockEntity) GetVelocity() vector.Vector3 {
	args := m.Called()
	return args.Get(0).(vector.Vector3)
}

func (m *MockEntity) GetAngularVelocity() vector.Vector3 {
	args := m.Called()
	return args.Get(0).(vector.Vector3)
}

func (m *MockEntity) SetAngularVelocity(angVel vector.Vector3) {
	m.Called(angVel)
}

func (m *MockEntity) Update(deltaTime float64) {
	m.Called(deltaTime)
}

// TestNewSpacecraftController tests the creation of a new spacecraft controller
func TestNewSpacecraftController(t *testing.T) {
	// Create a mock entity
	mockEntity := new(MockEntity)

	// Create a new spacecraft controller
	maxThrust := 1000.0
	maxTorque := 100.0
	controller := NewSpacecraftController(mockEntity, maxThrust, maxTorque)

	// Assert that the controller was created correctly
	assert.Equal(t, mockEntity, controller.GetEntity())
	assert.Equal(t, maxThrust, controller.maxThrust)
	assert.Equal(t, maxTorque, controller.maxTorque)
	assert.Equal(t, 0.0, controller.thrustLevel)
	assert.Equal(t, 45.0, controller.maxAngularVelocity)
	assert.Equal(t, 0.2, controller.angularDamping)

	// Check PID controller initialization
	assert.Equal(t, 2.0, controller.rotationPID.kP)
	assert.Equal(t, 0.1, controller.rotationPID.kI)
	assert.Equal(t, 0.5, controller.rotationPID.kD)

	// Check vector initializations
	assert.Equal(t, vector.Zero3(), controller.rotationPID.lastError)
	assert.Equal(t, vector.Zero3(), controller.rotationPID.errorSum)
	assert.Equal(t, vector.Zero3(), controller.rotationPID.errorDelta)
	assert.Equal(t, vector.Zero3(), controller.rotationPID.targetRotation)
}

// TestSetMaxAngularVelocity tests the SetMaxAngularVelocity method
func TestSetMaxAngularVelocity(t *testing.T) {
	// Create a mock entity
	mockEntity := new(MockEntity)

	// Create a new spacecraft controller
	controller := NewSpacecraftController(mockEntity, 1000.0, 100.0)

	// Test setting a positive value
	controller.SetMaxAngularVelocity(60.0)
	assert.Equal(t, 60.0, controller.maxAngularVelocity)

	// Test setting a negative value (should be clamped to 0)
	controller.SetMaxAngularVelocity(-10.0)
	assert.Equal(t, 0.0, controller.maxAngularVelocity)
}

// TestSetAngularDamping tests the SetAngularDamping method
func TestSetAngularDamping(t *testing.T) {
	// Create a mock entity
	mockEntity := new(MockEntity)

	// Create a new spacecraft controller
	controller := NewSpacecraftController(mockEntity, 1000.0, 100.0)

	// Test setting a value in range
	controller.SetAngularDamping(0.5)
	assert.Equal(t, 0.5, controller.angularDamping)

	// Test setting a value below 0 (should be clamped to 0)
	controller.SetAngularDamping(-0.1)
	assert.Equal(t, 0.0, controller.angularDamping)

	// Test setting a value above 1 (should be clamped to 1)
	controller.SetAngularDamping(1.5)
	assert.Equal(t, 1.0, controller.angularDamping)
}

// TestSetRotationPIDGains tests the SetRotationPIDGains method
func TestSetRotationPIDGains(t *testing.T) {
	// Create a mock entity
	mockEntity := new(MockEntity)

	// Create a new spacecraft controller
	controller := NewSpacecraftController(mockEntity, 1000.0, 100.0)

	// Test setting new PID gains
	kP := 3.0
	kI := 0.2
	kD := 0.7
	controller.SetRotationPIDGains(kP, kI, kD)

	// Assert that the gains were set correctly
	assert.Equal(t, kP, controller.rotationPID.kP)
	assert.Equal(t, kI, controller.rotationPID.kI)
	assert.Equal(t, kD, controller.rotationPID.kD)
}

// TestSetTargetRotation tests the SetTargetRotation method
func TestSetTargetRotation(t *testing.T) {
	// Create a mock entity
	mockEntity := new(MockEntity)

	// Create a new spacecraft controller
	controller := NewSpacecraftController(mockEntity, 1000.0, 100.0)

	// Test setting a new target rotation
	targetRotation := vector.NewVector3(45, 90, 180)
	controller.SetTargetRotation(targetRotation)

	// Assert that the target rotation was set correctly
	assert.Equal(t, targetRotation, controller.rotationPID.targetRotation)
}

// TestSetThrustLevel tests the SetThrustLevel method
func TestSetThrustLevel(t *testing.T) {
	// Create a mock entity
	mockEntity := new(MockEntity)

	// Create a new spacecraft controller
	controller := NewSpacecraftController(mockEntity, 1000.0, 100.0)

	// Test setting a value in range
	controller.SetThrustLevel(0.5)
	assert.Equal(t, 0.5, controller.thrustLevel)

	// Test setting a value below 0 (should be clamped to 0)
	controller.SetThrustLevel(-0.1)
	assert.Equal(t, 0.0, controller.thrustLevel)

	// Test setting a value above 1 (should be clamped to 1)
	controller.SetThrustLevel(1.5)
	assert.Equal(t, 1.0, controller.thrustLevel)
}

// TestNormalizeAngle tests the NormalizeAngle function
func TestNormalizeAngle(t *testing.T) {
	// Test cases
	testCases := []struct {
		input    float64
		expected float64
	}{
		{0.0, 0.0},
		{45.0, 45.0},
		{180.0, -180.0},
		{181.0, -179.0},
		{-45.0, -45.0},
		{-180.0, -180.0},
		{-181.0, 179.0},
		{360.0, 0.0},
		{720.0, 0.0},
		{-360.0, 0.0},
		{-720.0, 0.0},
	}

	// Run test cases
	for _, tc := range testCases {
		result := NormalizeAngle(tc.input)
		assert.InDelta(t, tc.expected, result, 0.001, "NormalizeAngle(%f) = %f, expected %f", tc.input, result, tc.expected)
	}
}

// TestApplyMainThrust tests the ApplyMainThrust method
func TestApplyMainThrust(t *testing.T) {
	// Create a mock body
	mockBody := new(MockBody)

	// Create a mock entity
	mockEntity := new(MockEntity)
	mockEntity.On("GetBody").Return(mockBody)
	mockEntity.On("GetRotation").Return(vector.NewVector3(0, 0, 0))

	// Create a new spacecraft controller
	maxThrust := 1000.0
	controller := NewSpacecraftController(mockEntity, maxThrust, 100.0)

	// Set thrust level to 0.5 (50%)
	controller.SetThrustLevel(0.5)

	// Expected thrust force (forward direction with 50% thrust)
	expectedForce := vector.NewVector3(0, 0, 500.0)

	// Set up expectations
	mockBody.On("ApplyForce", mock.MatchedBy(func(force vector.Vector3) bool {
		// Check if the force is approximately equal to the expected force
		return math.Abs(force.X()-expectedForce.X()) < 0.001 &&
			math.Abs(force.Y()-expectedForce.Y()) < 0.001 &&
			math.Abs(force.Z()-expectedForce.Z()) < 0.001
	})).Return()

	// Call ApplyMainThrust
	controller.ApplyMainThrust()

	// Assert that the mock methods were called
	mockEntity.AssertExpectations(t)
	mockBody.AssertExpectations(t)

	// Test with zero thrust level
	mockEntity = new(MockEntity)
	mockBody = new(MockBody)
	// No expectations needed since we return early when thrustLevel <= 0
	controller = NewSpacecraftController(mockEntity, maxThrust, 100.0)
	controller.SetThrustLevel(0.0)

	// Call ApplyMainThrust (should return early without calling any methods)
	controller.ApplyMainThrust()

	// No methods should be called
	mockEntity.AssertNotCalled(t, "GetBody")
	mockEntity.AssertNotCalled(t, "GetRotation")
}

// TestApplyRotation tests the ApplyRotation method
func TestApplyRotation(t *testing.T) {
	// Create a mock body
	mockBody := new(MockBody)

	// Create a mock entity
	mockEntity := new(MockEntity)
	mockEntity.On("GetBody").Return(mockBody)
	mockEntity.On("GetAngularVelocity").Return(vector.Zero3())

	// Create a new spacecraft controller
	maxTorque := 100.0
	controller := NewSpacecraftController(mockEntity, 1000.0, maxTorque)

	// Test applying rotation around the Y axis with 50% power
	axis := vector.NewVector3(0, 1, 0)
	amount := 0.5

	// Expected torque
	expectedTorque := vector.NewVector3(0, 50.0, 0)

	// Set up expectations
	mockBody.On("ApplyTorque", mock.MatchedBy(func(torque vector.Vector3) bool {
		// Check if the torque is approximately equal to the expected torque
		return math.Abs(torque.X()-expectedTorque.X()) < 0.001 &&
			math.Abs(torque.Y()-expectedTorque.Y()) < 0.001 &&
			math.Abs(torque.Z()-expectedTorque.Z()) < 0.001
	})).Return()

	// Call ApplyRotation
	controller.ApplyRotation(axis, amount)

	// Assert that the mock methods were called
	mockEntity.AssertExpectations(t)
	mockBody.AssertExpectations(t)

	// Test with amount outside range (should be clamped)
	mockEntity = new(MockEntity)
	mockBody = new(MockBody)
	mockEntity.On("GetBody").Return(mockBody)
	mockEntity.On("GetAngularVelocity").Return(vector.Zero3())
	controller = NewSpacecraftController(mockEntity, 1000.0, maxTorque)

	// Expected torque for amount = 1.0 (clamped from 2.0)
	expectedTorque = vector.NewVector3(0, 100.0, 0)

	// Set up expectations
	mockBody.On("ApplyTorque", mock.MatchedBy(func(torque vector.Vector3) bool {
		// Check if the torque is approximately equal to the expected torque
		return math.Abs(torque.X()-expectedTorque.X()) < 0.001 &&
			math.Abs(torque.Y()-expectedTorque.Y()) < 0.001 &&
			math.Abs(torque.Z()-expectedTorque.Z()) < 0.001
	})).Return()

	// Call ApplyRotation with amount outside range
	controller.ApplyRotation(axis, 2.0)

	// Assert that the mock methods were called
	mockEntity.AssertExpectations(t)
	mockBody.AssertExpectations(t)
}

// TestUpdatePIDController tests the UpdatePIDController method
func TestUpdatePIDController(t *testing.T) {
	// Create a mock body
	mockBody := new(MockBody)

	// Create a mock entity
	mockEntity := new(MockEntity)
	mockEntity.On("GetBody").Return(mockBody)

	// Current rotation
	currentRotation := vector.NewVector3(10, 20, 30)
	mockEntity.On("GetRotation").Return(currentRotation)

	// Create a new spacecraft controller
	controller := NewSpacecraftController(mockEntity, 1000.0, 100.0)

	// Set target rotation
	targetRotation := vector.NewVector3(15, 25, 35)
	controller.SetTargetRotation(targetRotation)

	// Expected error
	expectedError := vector.NewVector3(5, 5, 5)

	// Set up expectations for ApplyTorque
	mockBody.On("ApplyTorque", mock.Anything).Return()

	// Call UpdatePIDController
	deltaTime := 0.1
	controller.UpdatePIDController(deltaTime)

	// Assert that the mock methods were called
	mockEntity.AssertExpectations(t)
	mockBody.AssertExpectations(t)

	// Check that the error was calculated correctly
	assert.InDelta(t, expectedError.X(), controller.rotationPID.lastError.X(), 0.001)
	assert.InDelta(t, expectedError.Y(), controller.rotationPID.lastError.Y(), 0.001)
	assert.InDelta(t, expectedError.Z(), controller.rotationPID.lastError.Z(), 0.001)
}

// TestApplyAngularDamping tests the ApplyAngularDamping method
func TestApplyAngularDamping(t *testing.T) {
	// Create a mock body
	mockBody := new(MockBody)

	// Create a mock entity
	mockEntity := new(MockEntity)
	mockEntity.On("GetBody").Return(mockBody)

	// Angular velocity
	angularVelocity := vector.NewVector3(10, 20, 30)
	mockEntity.On("GetAngularVelocity").Return(angularVelocity)

	// Create a new spacecraft controller
	controller := NewSpacecraftController(mockEntity, 1000.0, 100.0)
	controller.SetAngularDamping(0.5)

	// Expected damping torque (opposite to angular velocity, scaled by damping factor)
	expectedTorque := vector.NewVector3(-5, -10, -15)

	// Set up expectations
	mockBody.On("ApplyTorque", mock.MatchedBy(func(torque vector.Vector3) bool {
		// Check if the torque is approximately equal to the expected torque
		return math.Abs(torque.X()-expectedTorque.X()) < 0.001 &&
			math.Abs(torque.Y()-expectedTorque.Y()) < 0.001 &&
			math.Abs(torque.Z()-expectedTorque.Z()) < 0.001
	})).Return()

	// Call ApplyAngularDamping
	controller.ApplyAngularDamping()

	// Assert that the mock methods were called
	mockEntity.AssertExpectations(t)
	mockBody.AssertExpectations(t)

	// Test with zero damping
	mockEntity = new(MockEntity)
	mockBody = new(MockBody)
	// No expectations needed since we return early when angularDamping <= 0
	controller = NewSpacecraftController(mockEntity, 1000.0, 100.0)
	controller.SetAngularDamping(0.0)

	// Call ApplyAngularDamping (should return early without calling any methods)
	controller.ApplyAngularDamping()

	// No methods should be called
	mockEntity.AssertNotCalled(t, "GetBody")
	mockEntity.AssertNotCalled(t, "GetAngularVelocity")
}

// TestUpdate tests the Update method
func TestUpdate(t *testing.T) {
	// Create a mock body
	mockBody := new(MockBody)

	// Create a mock entity with expectations
	mockEntity := new(MockEntity)
	// GetBody is not called because there's no significant error in the PID controller
	mockEntity.On("GetRotation").Return(vector.Zero3()) // Called by UpdatePIDController
	mockEntity.On("Update", 0.1).Return()               // Called by BaseController.Update

	// Create a new spacecraft controller with zero thrust and damping
	controller := NewSpacecraftController(mockEntity, 1000.0, 100.0)
	controller.SetThrustLevel(0.0)
	controller.SetAngularDamping(0.0)
	controller.SetTargetRotation(vector.Zero3())

	// Set up expectations
	mockBody.On("ApplyTorque", mock.Anything).Return().Maybe()

	// Call Update
	controller.Update(0.1)

	// Assert that the mock methods were called
	mockEntity.AssertExpectations(t)
}
