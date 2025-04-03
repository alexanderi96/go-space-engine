// Package collision provides interfaces and implementations for collision detection and resolution
package collision

import (
	"math"

	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/physics/body"
)

// CollisionInfo contains information about a collision
type CollisionInfo struct {
	BodyA       body.Body      // First body involved in the collision
	BodyB       body.Body      // Second body involved in the collision
	Point       vector.Vector3 // Contact point
	Normal      vector.Vector3 // Collision normal (from A to B)
	Depth       float64        // Penetration depth
	HasCollided bool           // Indicates if a collision occurred
}

// Collider detects collisions between bodies
type Collider interface {
	// CheckCollision checks if two bodies collide
	CheckCollision(a, b body.Body) CollisionInfo
}

// CollisionResolver resolves collisions
type CollisionResolver interface {
	// ResolveCollision resolves a collision
	ResolveCollision(info CollisionInfo)
}

// SphereCollider implements a collision detector for spherical bodies
type SphereCollider struct{}

// NewSphereCollider creates a new collision detector for spherical bodies
func NewSphereCollider() *SphereCollider {
	return &SphereCollider{}
}

// CheckCollision checks if two spherical bodies collide
func (sc *SphereCollider) CheckCollision(a, b body.Body) CollisionInfo {
	// Calculate the direction vector from a to b
	direction := b.Position().Sub(a.Position())

	// Calculate the squared distance
	distanceSquared := direction.LengthSquared()

	// Calculate the sum of radii
	radiusA := a.Radius().Value()
	radiusB := b.Radius().Value()
	sumRadii := radiusA + radiusB

	// Check if there is a collision
	hasCollided := distanceSquared < sumRadii*sumRadii

	// If there is no collision, return an empty collision info
	if !hasCollided {
		return CollisionInfo{
			BodyA:       a,
			BodyB:       b,
			HasCollided: false,
		}
	}

	// Calculate the distance
	distance := math.Sqrt(distanceSquared)

	// Calculate the collision normal
	var normal vector.Vector3
	if distance > 1e-10 {
		normal = direction.Scale(1.0 / distance)
	} else {
		// If the bodies completely overlap, use a default normal
		normal = vector.NewVector3(1, 0, 0)
	}

	// Calculate the penetration depth
	depth := sumRadii - distance

	// Calculate the contact point
	point := a.Position().Add(normal.Scale(radiusA))

	return CollisionInfo{
		BodyA:       a,
		BodyB:       b,
		Point:       point,
		Normal:      normal,
		Depth:       depth,
		HasCollided: true,
	}
}

// ImpulseResolver implements an impulse-based collision resolver
type ImpulseResolver struct {
	restitution float64 // Coefficient of restitution (elasticity)
}

// NewImpulseResolver creates a new impulse-based collision resolver
func NewImpulseResolver(restitution float64) *ImpulseResolver {
	return &ImpulseResolver{
		restitution: restitution,
	}
}

// ResolveCollision resolves a collision using the impulse method
func (ir *ImpulseResolver) ResolveCollision(info CollisionInfo) {
	// If there was no collision, do nothing
	if !info.HasCollided {
		return
	}

	a := info.BodyA
	b := info.BodyB

	// If both bodies are static, do nothing
	if a.IsStatic() && b.IsStatic() {
		return
	}

	// Calculate the relative velocity
	relativeVelocity := b.Velocity().Sub(a.Velocity())

	// Calculate the relative velocity along the normal
	velocityAlongNormal := relativeVelocity.Dot(info.Normal)

	// If the bodies are moving away from each other, do nothing
	if velocityAlongNormal > 0 {
		return
	}

	// Calculate the coefficient of restitution (elasticity)
	// Use the minimum between the resolver's restitution and the materials' elasticity
	elasticityA := a.Material().Elasticity()
	elasticityB := b.Material().Elasticity()
	restitution := math.Min(ir.restitution, math.Min(elasticityA, elasticityB))

	// Calculate the scalar impulse
	// j = -(1 + e) * velocityAlongNormal / (1/massA + 1/massB)
	massA := a.Mass().Value()
	massB := b.Mass().Value()

	// Handle static bodies (infinite mass)
	var inverseMassA, inverseMassB float64
	if a.IsStatic() {
		inverseMassA = 0
	} else {
		inverseMassA = 1.0 / massA
	}

	if b.IsStatic() {
		inverseMassB = 0
	} else {
		inverseMassB = 1.0 / massB
	}

	// Calculate the scalar impulse
	j := -(1.0 + restitution) * velocityAlongNormal
	j /= inverseMassA + inverseMassB

	// Apply the impulse
	impulse := info.Normal.Scale(j)

	if !a.IsStatic() {
		a.SetVelocity(a.Velocity().Sub(impulse.Scale(inverseMassA)))
	}

	if !b.IsStatic() {
		b.SetVelocity(b.Velocity().Add(impulse.Scale(inverseMassB)))
	}

	// Correct the penetration (position resolution)
	ir.resolvePosition(info)
}

// resolvePosition corrects the penetration between bodies
func (ir *ImpulseResolver) resolvePosition(info CollisionInfo) {
	a := info.BodyA
	b := info.BodyB

	// If both bodies are static, do nothing
	if a.IsStatic() && b.IsStatic() {
		return
	}

	// Position correction constant (0.2 - 0.8)
	const percent = 0.2
	// Minimum penetration threshold
	const slop = 0.01

	// Calculate the inverse masses
	massA := a.Mass().Value()
	massB := b.Mass().Value()
	var inverseMassA, inverseMassB float64
	if a.IsStatic() {
		inverseMassA = 0
	} else {
		inverseMassA = 1.0 / massA
	}

	if b.IsStatic() {
		inverseMassB = 0
	} else {
		inverseMassB = 1.0 / massB
	}

	// Calculate the sum of inverse masses
	totalInverseMass := inverseMassA + inverseMassB
	if totalInverseMass <= 0 {
		return // Both bodies have infinite mass
	}

	// Calculate the correction
	correction := math.Max(info.Depth-slop, 0.0) * percent

	// Apply the correction
	if !a.IsStatic() {
		a.SetPosition(a.Position().Sub(info.Normal.Scale(correction * inverseMassA / totalInverseMass)))
	}

	if !b.IsStatic() {
		b.SetPosition(b.Position().Add(info.Normal.Scale(correction * inverseMassB / totalInverseMass)))
	}
}
