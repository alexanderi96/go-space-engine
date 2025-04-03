// Package space provides spatial structures to optimize spatial queries
package space

import (
	"math"
	"sync"

	"github.com/alexanderi96/go-space-engine/core/constants"
	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/physics/body"
)

// TaskSubmitter represents an interface for submitting tasks to be executed in parallel
type TaskSubmitter interface {
	// Submit submits a task to be executed
	Submit(task func())
	// Wait waits for all tasks to be completed
	Wait()
}

// Region represents a region of space
type Region interface {
	// Contains checks if a point is contained in the region
	Contains(point vector.Vector3) bool
	// ContainsSphere checks if a sphere is contained in the region
	ContainsSphere(center vector.Vector3, radius float64) bool
	// Intersects checks if the region intersects another region
	Intersects(other Region) bool
}

// AABB represents an Axis-Aligned Bounding Box
type AABB struct {
	Min vector.Vector3 // Minimum point (bottom left back corner)
	Max vector.Vector3 // Maximum point (top right front corner)
}

// NewAABB creates a new AABB
func NewAABB(min, max vector.Vector3) *AABB {
	return &AABB{
		Min: min,
		Max: max,
	}
}

// Contains checks if a point is contained in the AABB
func (aabb *AABB) Contains(point vector.Vector3) bool {
	return point.X() >= aabb.Min.X() && point.X() <= aabb.Max.X() &&
		point.Y() >= aabb.Min.Y() && point.Y() <= aabb.Max.Y() &&
		point.Z() >= aabb.Min.Z() && point.Z() <= aabb.Max.Z()
}

// ContainsSphere checks if a sphere is contained in the AABB
func (aabb *AABB) ContainsSphere(center vector.Vector3, radius float64) bool {
	// Calculate the squared distance between the sphere center and the closest point of the AABB
	closestX := math.Max(aabb.Min.X(), math.Min(center.X(), aabb.Max.X()))
	closestY := math.Max(aabb.Min.Y(), math.Min(center.Y(), aabb.Max.Y()))
	closestZ := math.Max(aabb.Min.Z(), math.Min(center.Z(), aabb.Max.Z()))

	distanceSquared := (closestX-center.X())*(closestX-center.X()) +
		(closestY-center.Y())*(closestY-center.Y()) +
		(closestZ-center.Z())*(closestZ-center.Z())

	// The sphere is contained if the squared distance is less than or equal to the squared radius
	return distanceSquared <= radius*radius
}

// Intersects checks if the AABB intersects another AABB
func (aabb *AABB) Intersects(other Region) bool {
	otherAABB, ok := other.(*AABB)
	if !ok {
		// If the other region is not an AABB, use a generic implementation
		return false
	}

	// Two AABBs intersect if they overlap in all three dimensions
	return aabb.Min.X() <= otherAABB.Max.X() && aabb.Max.X() >= otherAABB.Min.X() &&
		aabb.Min.Y() <= otherAABB.Max.Y() && aabb.Max.Y() >= otherAABB.Min.Y() &&
		aabb.Min.Z() <= otherAABB.Max.Z() && aabb.Max.Z() >= otherAABB.Min.Z()
}

// Center returns the center of the AABB
func (aabb *AABB) Center() vector.Vector3 {
	return aabb.Min.Add(aabb.Max).Scale(0.5)
}

// Size returns the dimensions of the AABB
func (aabb *AABB) Size() vector.Vector3 {
	return aabb.Max.Sub(aabb.Min)
}

// SpatialStructure represents a spatial structure to optimize spatial queries
type SpatialStructure interface {
	// Insert inserts a body into the structure
	Insert(b body.Body)
	// Remove removes a body from the structure
	Remove(b body.Body)
	// Update updates the position of a body in the structure
	Update(b body.Body)
	// UpdateAll updates the position of multiple bodies in the structure
	UpdateAll(bodies []body.Body, taskSubmitter TaskSubmitter)
	// Query returns all bodies that might interact with the specified region
	Query(region Region) []body.Body
	// QuerySphere returns all bodies that might interact with the specified sphere
	QuerySphere(center vector.Vector3, radius float64) []body.Body
	// Clear removes all bodies from the structure
	Clear()
}

// Octree implements an optimized spatial structure based on an octree
type Octree struct {
	bounds     *AABB       // Octree bounds
	maxObjects int         // Maximum number of objects per node
	maxLevels  int         // Maximum number of levels
	level      int         // Current level
	objects    []body.Body // Objects in this node
	children   [8]*Octree  // Octree children
	divided    bool        // Indicates if the octree has been divided

	// Fields for gravity calculation
	totalMass    float64        // Total mass of all bodies in this node and its children
	centerOfMass vector.Vector3 // Center of mass of all bodies in this node and its children

	// Mutex to protect concurrent access
	mutex sync.RWMutex
}

// NewOctree creates a new octree
func NewOctree(bounds *AABB, maxObjects, maxLevels int) *Octree {
	return &Octree{
		bounds:       bounds,
		maxObjects:   maxObjects,
		maxLevels:    maxLevels,
		level:        0,
		objects:      make([]body.Body, 0),
		divided:      false,
		totalMass:    0,
		centerOfMass: vector.Zero3(),
		mutex:        sync.RWMutex{},
	}
}

// Insert inserts a body into the octree
func (ot *Octree) Insert(b body.Body) {
	ot.mutex.Lock()
	defer ot.mutex.Unlock()

	ot.insertUnsafe(b)
}

// insertUnsafe inserts a body into the octree without locking the mutex
func (ot *Octree) insertUnsafe(b body.Body) {
	// If the octree is already divided, insert into the appropriate children
	if ot.divided {
		indices := ot.getIndices(b)
		for _, index := range indices {
			if index != -1 {
				ot.children[index].insertUnsafe(b)
			}
		}

		// Update the center of mass and total mass
		ot.updateMassAndCenterOfMass(b, true)
		return
	}

	// Add the object to this node
	ot.objects = append(ot.objects, b)

	// Update the center of mass and total mass
	ot.updateMassAndCenterOfMass(b, true)

	// Check if it's necessary to divide the octree
	if len(ot.objects) > ot.maxObjects && ot.level < ot.maxLevels {
		// Divide the octree
		ot.split()

		// Redistribute objects to children
		for i := 0; i < len(ot.objects); i++ {
			indices := ot.getIndices(ot.objects[i])
			for _, index := range indices {
				if index != -1 {
					ot.children[index].insertUnsafe(ot.objects[i])
				}
			}
		}

		// Empty the objects of this node
		ot.objects = make([]body.Body, 0)
	}
}

// Remove removes a body from the octree
func (ot *Octree) Remove(b body.Body) {
	ot.mutex.Lock()
	defer ot.mutex.Unlock()

	ot.removeUnsafe(b)
}

// removeUnsafe removes a body from the octree without locking the mutex
func (ot *Octree) removeUnsafe(b body.Body) {
	// If the octree is divided, remove from the appropriate children
	if ot.divided {
		indices := ot.getIndices(b)
		for _, index := range indices {
			if index != -1 {
				ot.children[index].removeUnsafe(b)
			}
		}

		// Update the center of mass and total mass
		ot.updateMassAndCenterOfMass(b, false)
		return
	}

	// Remove the object from this node
	for i, obj := range ot.objects {
		if obj.ID() == b.ID() {
			// Remove the object by swapping it with the last one and truncating the slice
			lastIndex := len(ot.objects) - 1
			ot.objects[i] = ot.objects[lastIndex]
			ot.objects = ot.objects[:lastIndex]

			// Update the center of mass and total mass
			ot.updateMassAndCenterOfMass(b, false)
			break
		}
	}
}

// UpdateAll updates the position of multiple bodies in the octree
func (ot *Octree) UpdateAll(bodies []body.Body, taskSubmitter TaskSubmitter) {
	for _, b := range bodies {
		b := b // Capture the variable for the goroutine
		taskSubmitter.Submit(func() {
			ot.Update(b)
		})
	}
	taskSubmitter.Wait()
}

// Update updates the position of a body in the octree
func (ot *Octree) Update(b body.Body) {
	ot.mutex.Lock()
	defer ot.mutex.Unlock()

	// Remove and reinsert the object
	ot.removeUnsafe(b)
	ot.insertUnsafe(b)
}

// Query returns all bodies that might interact with the specified region
func (ot *Octree) Query(region Region) []body.Body {
	ot.mutex.RLock()
	defer ot.mutex.RUnlock()

	result := make([]body.Body, 0)

	// Check if the region intersects this node
	if !region.Intersects(ot.bounds) {
		return result
	}

	// Add the objects of this node
	result = append(result, ot.objects...)

	// If the octree is divided, query the children
	if ot.divided {
		for i := 0; i < 8; i++ {
			childResult := ot.children[i].Query(region)
			result = append(result, childResult...)
		}
	}

	return result
}

// QuerySphere returns all bodies that might interact with the specified sphere
func (ot *Octree) QuerySphere(center vector.Vector3, radius float64) []body.Body {
	ot.mutex.RLock()
	defer ot.mutex.RUnlock()

	result := make([]body.Body, 0)

	// Check if the sphere intersects this node
	if !ot.bounds.ContainsSphere(center, radius) {
		return result
	}

	// Add the objects of this node
	result = append(result, ot.objects...)

	// If the octree is divided, query the children
	if ot.divided {
		for i := 0; i < 8; i++ {
			childResult := ot.children[i].QuerySphere(center, radius)
			result = append(result, childResult...)
		}
	}

	return result
}

// Clear removes all bodies from the octree
func (ot *Octree) Clear() {
	ot.mutex.Lock()
	defer ot.mutex.Unlock()

	ot.objects = make([]body.Body, 0)

	// Reset the center of mass and total mass
	ot.totalMass = 0
	ot.centerOfMass = vector.Zero3()

	if ot.divided {
		for i := 0; i < 8; i++ {
			ot.children[i].Clear()
			ot.children[i] = nil
		}
		ot.divided = false
	}
}

// split divides the octree into eight children
func (ot *Octree) split() {
	// Calculate the center of the octree
	center := ot.bounds.Center()

	// Create the eight children
	// Order: [0] = Bottom Left Back, [1] = Bottom Right Back, [2] = Bottom Right Front, [3] = Bottom Left Front,
	//        [4] = Top Left Back, [5] = Top Right Back, [6] = Top Right Front, [7] = Top Left Front
	childBounds := [8]*AABB{
		// [0] Bottom Left Back
		NewAABB(
			ot.bounds.Min,
			center,
		),
		// [1] Bottom Right Back
		NewAABB(
			vector.NewVector3(center.X(), ot.bounds.Min.Y(), ot.bounds.Min.Z()),
			vector.NewVector3(ot.bounds.Max.X(), center.Y(), center.Z()),
		),
		// [2] Bottom Right Front
		NewAABB(
			vector.NewVector3(center.X(), ot.bounds.Min.Y(), center.Z()),
			vector.NewVector3(ot.bounds.Max.X(), center.Y(), ot.bounds.Max.Z()),
		),
		// [3] Bottom Left Front
		NewAABB(
			vector.NewVector3(ot.bounds.Min.X(), ot.bounds.Min.Y(), center.Z()),
			vector.NewVector3(center.X(), center.Y(), ot.bounds.Max.Z()),
		),
		// [4] Top Left Back
		NewAABB(
			vector.NewVector3(ot.bounds.Min.X(), center.Y(), ot.bounds.Min.Z()),
			vector.NewVector3(center.X(), ot.bounds.Max.Y(), center.Z()),
		),
		// [5] Top Right Back
		NewAABB(
			vector.NewVector3(center.X(), center.Y(), ot.bounds.Min.Z()),
			vector.NewVector3(ot.bounds.Max.X(), ot.bounds.Max.Y(), center.Z()),
		),
		// [6] Top Right Front
		NewAABB(
			center,
			ot.bounds.Max,
		),
		// [7] Top Left Front
		NewAABB(
			vector.NewVector3(ot.bounds.Min.X(), center.Y(), center.Z()),
			vector.NewVector3(center.X(), ot.bounds.Max.Y(), ot.bounds.Max.Z()),
		),
	}

	// Create the octree children
	for i := 0; i < 8; i++ {
		ot.children[i] = NewOctree(childBounds[i], ot.maxObjects, ot.maxLevels)
		ot.children[i].level = ot.level + 1
	}

	ot.divided = true
}

// getIndices determines which children a body should be inserted into
func (ot *Octree) getIndices(b body.Body) []int {
	result := make([]int, 0, 8)
	center := ot.bounds.Center()
	position := b.Position()
	radius := b.Radius().Value()

	// Determine in which octants the body is located
	top := position.Y()+radius > center.Y()
	bottom := position.Y()-radius < center.Y()
	left := position.X()-radius < center.X()
	right := position.X()+radius > center.X()
	front := position.Z()+radius > center.Z()
	back := position.Z()-radius < center.Z()

	// Bottom Left Back
	if bottom && left && back {
		result = append(result, 0)
	}

	// Bottom Right Back
	if bottom && right && back {
		result = append(result, 1)
	}

	// Bottom Right Front
	if bottom && right && front {
		result = append(result, 2)
	}

	// Bottom Left Front
	if bottom && left && front {
		result = append(result, 3)
	}

	// Top Left Back
	if top && left && back {
		result = append(result, 4)
	}

	// Top Right Back
	if top && right && back {
		result = append(result, 5)
	}

	// Top Right Front
	if top && right && front {
		result = append(result, 6)
	}

	// Top Left Front
	if top && left && front {
		result = append(result, 7)
	}

	if len(result) == 0 {
		result = append(result, -1)
	}

	return result
}

// updateMassAndCenterOfMass updates the center of mass and total mass
func (ot *Octree) updateMassAndCenterOfMass(b body.Body, adding bool) {
	mass := b.Mass().Value()
	position := b.Position()

	if adding {
		// Add the body
		oldTotalMass := ot.totalMass
		ot.totalMass += mass

		if oldTotalMass > 0 {
			// Update the center of mass
			ot.centerOfMass = ot.centerOfMass.Scale(oldTotalMass).Add(position.Scale(mass)).Scale(1.0 / ot.totalMass)
		} else {
			// If it's the first body, the center of mass is its position
			ot.centerOfMass = position
		}
	} else {
		// Remove the body
		if ot.totalMass > mass {
			// Update the center of mass
			oldTotalMass := ot.totalMass
			ot.totalMass -= mass
			ot.centerOfMass = ot.centerOfMass.Scale(oldTotalMass).Sub(position.Scale(mass)).Scale(1.0 / ot.totalMass)
		} else {
			// If it was the last body, reset the center of mass
			ot.totalMass = 0
			ot.centerOfMass = vector.Zero3()
		}
	}

	// If the octree is divided, propagate the update to the children
	if ot.divided {
		// Recalculate the center of mass from the children
		ot.totalMass = 0
		weightedPosition := vector.Zero3()

		for i := 0; i < 8; i++ {
			if ot.children[i] != nil {
				childMass := ot.children[i].totalMass
				ot.totalMass += childMass
				if childMass > 0 {
					weightedPosition = weightedPosition.Add(ot.children[i].centerOfMass.Scale(childMass))
				}
			}
		}

		if ot.totalMass > 0 {
			ot.centerOfMass = weightedPosition.Scale(1.0 / ot.totalMass)
		}
	}
}

// CalculateGravity calculates the gravitational force on a body using the Barnes-Hut algorithm
func (ot *Octree) CalculateGravity(b body.Body, theta float64) vector.Vector3 {
	ot.mutex.RLock()
	defer ot.mutex.RUnlock()

	force := vector.Zero3()
	ot.calculateGravityRecursive(b, theta, &force)
	return force
}

// calculateGravityRecursive recursively calculates the gravitational force
func (ot *Octree) calculateGravityRecursive(b body.Body, theta float64, force *vector.Vector3) {
	// If the octree is not divided or has no bodies, calculate the force directly
	if !ot.divided || ot.totalMass == 0 {
		ot.calculateLeafNodeGravity(b, force)
		return
	}

	// Calculate the node width and the distance from the body to the center of mass
	width := ot.bounds.Max.X() - ot.bounds.Min.X()
	deltaPos := ot.centerOfMass.Sub(b.Position())
	distanceSquared := deltaPos.LengthSquared()

	// Avoid division by zero
	if distanceSquared < 1e-10 {
		return
	}

	// If the width/distance ratio is less than theta, approximate with the center of mass
	if (width * width) < (theta * theta * distanceSquared) {
		ot.approximateGravityWithCenterOfMass(b, force)
		return
	}

	// Otherwise, calculate recursively for each child
	for i := 0; i < 8; i++ {
		if ot.children[i] != nil && ot.children[i].totalMass > 0 {
			ot.children[i].calculateGravityRecursive(b, theta, force)
		}
	}
}

// calculateLeafNodeGravity calculates the gravitational force for each body in the leaf node
func (ot *Octree) calculateLeafNodeGravity(b body.Body, force *vector.Vector3) {

	// Body mass
	bodyMass := b.Mass().Value()
	bodyPos := b.Position()

	// Calculate the force for each body in the node
	for _, obj := range ot.objects {
		// Avoid calculating the force on itself
		if obj.ID() == b.ID() {
			continue
		}

		// Calculate the direction vector
		deltaPos := obj.Position().Sub(bodyPos)
		distanceSquared := deltaPos.LengthSquared()

		// Avoid division by zero
		if distanceSquared <= 1e-10 {
			continue
		}

		// Calculate the gravitational force
		distance := math.Sqrt(distanceSquared)
		direction := deltaPos.Scale(1.0 / distance)

		// F = G * m1 * m2 / r^2
		forceMagnitude := constants.G * bodyMass * obj.Mass().Value() / distanceSquared

		// Add the force to the total force vector
		forceVector := *force
		*force = forceVector.Add(direction.Scale(forceMagnitude))
	}
}

// approximateGravityWithCenterOfMass approximates the gravitational force using the center of mass
func (ot *Octree) approximateGravityWithCenterOfMass(b body.Body, force *vector.Vector3) {

	// Body mass
	bodyMass := b.Mass().Value()
	bodyPos := b.Position()

	// Calculate the direction vector
	deltaPos := ot.centerOfMass.Sub(bodyPos)
	distanceSquared := deltaPos.LengthSquared()

	// Avoid division by zero
	if distanceSquared <= 1e-10 {
		return
	}

	// Calculate the gravitational force
	distance := math.Sqrt(distanceSquared)
	direction := deltaPos.Scale(1.0 / distance)

	// F = G * m1 * m2 / r^2
	forceMagnitude := constants.G * bodyMass * ot.totalMass / distanceSquared

	// Add the force to the total force vector
	forceVector := *force
	*force = forceVector.Add(direction.Scale(forceMagnitude))
}
