package controller

import (
	"testing"

	"github.com/alexanderi96/go-space-engine/core/units"
	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/physics/body"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

// TestNewBaseController tests the creation of a new base controller
func TestNewBaseController(t *testing.T) {
	// Create a mock entity
	mockEntity := new(MockEntity)

	// Create a new base controller
	controller := NewBaseController(mockEntity)

	// Assert that the controller was created correctly
	assert.Equal(t, mockEntity, controller.GetEntity())
}

// TestGetEntity tests the GetEntity method
func TestGetEntity(t *testing.T) {
	// Create a mock entity
	mockEntity := new(MockEntity)

	// Create a new base controller
	controller := NewBaseController(mockEntity)

	// Assert that GetEntity returns the correct entity
	assert.Equal(t, mockEntity, controller.GetEntity())
}

// TestApplyForce tests the ApplyForce method
func TestApplyForce(t *testing.T) {
	// Create a mock body
	mockBody := new(MockBody)

	// Create a mock entity
	mockEntity := new(MockEntity)
	mockEntity.On("GetBody").Return(mockBody)

	// Create a force vector
	force := vector.NewVector3(10, 20, 30)

	// Set up expectations
	mockBody.On("ApplyForce", force).Return()

	// Create a new base controller
	controller := NewBaseController(mockEntity)

	// Call ApplyForce
	controller.ApplyForce(force)

	// Assert that the mock methods were called
	mockEntity.AssertExpectations(t)
	mockBody.AssertExpectations(t)
}

// TestApplyTorque tests the ApplyTorque method
func TestApplyTorque(t *testing.T) {
	// Create a mock body
	mockBody := new(MockBody)

	// Create a mock entity
	mockEntity := new(MockEntity)
	mockEntity.On("GetBody").Return(mockBody)

	// Create a torque vector
	torque := vector.NewVector3(1, 2, 3)

	// Set up expectations
	mockBody.On("ApplyTorque", torque).Return()

	// Create a new base controller
	controller := NewBaseController(mockEntity)

	// Call ApplyTorque
	controller.ApplyTorque(torque)

	// Assert that the mock methods were called
	mockEntity.AssertExpectations(t)
	mockBody.AssertExpectations(t)
}

// TestUpdate tests the Update method
func TestUpdate(t *testing.T) {
	// Create a mock entity
	mockEntity := new(MockEntity)

	// Set up expectations
	deltaTime := 0.1
	mockEntity.On("Update", deltaTime).Return()

	// Create a new base controller
	controller := NewBaseController(mockEntity)

	// Call Update
	controller.Update(deltaTime)

	// Assert that the mock method was called
	mockEntity.AssertExpectations(t)
}

// TestNilBody tests handling of nil body
func TestNilBody(t *testing.T) {
	// Create a mock entity that returns nil for GetBody
	mockEntity := new(MockEntity)
	mockEntity.On("GetBody").Return(nil)

	// Create a new base controller
	controller := NewBaseController(mockEntity)

	// These should not panic
	controller.ApplyForce(vector.NewVector3(1, 2, 3))
	controller.ApplyTorque(vector.NewVector3(1, 2, 3))

	// Assert that the mock method was called
	mockEntity.AssertExpectations(t)
}
