package entity

import (
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

// TestNewBaseEntity tests the creation of a new base entity
func TestNewBaseEntity(t *testing.T) {
	// Create a mock body
	mockBody := new(MockBody)

	// Create a new base entity
	id := "test-entity"
	entity := NewBaseEntity(id, mockBody)

	// Assert that the entity was created correctly
	assert.Equal(t, id, entity.GetID())
	assert.Equal(t, mockBody, entity.GetBody())
}

// TestGetID tests the GetID method
func TestGetID(t *testing.T) {
	// Create a mock body
	mockBody := new(MockBody)

	// Create a new base entity
	id := "test-entity"
	entity := NewBaseEntity(id, mockBody)

	// Assert that GetID returns the correct ID
	assert.Equal(t, id, entity.GetID())
}

// TestGetBody tests the GetBody method
func TestGetBody(t *testing.T) {
	// Create a mock body
	mockBody := new(MockBody)

	// Create a new base entity
	entity := NewBaseEntity("test-entity", mockBody)

	// Assert that GetBody returns the correct body
	assert.Equal(t, mockBody, entity.GetBody())
}

// TestGetPosition tests the GetPosition method
func TestGetPosition(t *testing.T) {
	// Create a mock body
	mockBody := new(MockBody)
	position := vector.NewVector3(1, 2, 3)
	mockBody.On("Position").Return(position)

	// Create a new base entity
	entity := NewBaseEntity("test-entity", mockBody)

	// Assert that GetPosition returns the correct position
	assert.Equal(t, position, entity.GetPosition())
	mockBody.AssertExpectations(t)
}

// TestGetRotation tests the GetRotation method
func TestGetRotation(t *testing.T) {
	// Create a mock body
	mockBody := new(MockBody)
	rotation := vector.NewVector3(10, 20, 30)
	mockBody.On("Rotation").Return(rotation)

	// Create a new base entity
	entity := NewBaseEntity("test-entity", mockBody)

	// Assert that GetRotation returns the correct rotation
	assert.Equal(t, rotation, entity.GetRotation())
	mockBody.AssertExpectations(t)
}

// TestGetVelocity tests the GetVelocity method
func TestGetVelocity(t *testing.T) {
	// Create a mock body
	mockBody := new(MockBody)
	velocity := vector.NewVector3(5, 10, 15)
	mockBody.On("Velocity").Return(velocity)

	// Create a new base entity
	entity := NewBaseEntity("test-entity", mockBody)

	// Assert that GetVelocity returns the correct velocity
	assert.Equal(t, velocity, entity.GetVelocity())
	mockBody.AssertExpectations(t)
}

// TestGetAngularVelocity tests the GetAngularVelocity method
func TestGetAngularVelocity(t *testing.T) {
	// Create a mock body
	mockBody := new(MockBody)
	angularVelocity := vector.NewVector3(1, 2, 3)
	mockBody.On("AngularVelocity").Return(angularVelocity)

	// Create a new base entity
	entity := NewBaseEntity("test-entity", mockBody)

	// Assert that GetAngularVelocity returns the correct angular velocity
	assert.Equal(t, angularVelocity, entity.GetAngularVelocity())
	mockBody.AssertExpectations(t)
}

// TestSetAngularVelocity tests the SetAngularVelocity method
func TestSetAngularVelocity(t *testing.T) {
	// Create a mock body
	mockBody := new(MockBody)
	angularVelocity := vector.NewVector3(1, 2, 3)
	mockBody.On("SetAngularVelocity", angularVelocity).Return()

	// Create a new base entity
	entity := NewBaseEntity("test-entity", mockBody)

	// Call SetAngularVelocity
	entity.SetAngularVelocity(angularVelocity)

	// Assert that the mock method was called
	mockBody.AssertExpectations(t)
}

// TestUpdate tests the Update method
func TestUpdate(t *testing.T) {
	// Create a mock body
	mockBody := new(MockBody)
	deltaTime := 0.1
	mockBody.On("Update", deltaTime).Return()

	// Create a new base entity
	entity := NewBaseEntity("test-entity", mockBody)

	// Call Update
	entity.Update(deltaTime)

	// Assert that the mock method was called
	mockBody.AssertExpectations(t)
}
