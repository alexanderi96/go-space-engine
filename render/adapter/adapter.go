// Package adapter provides interfaces for rendering
package adapter

import (
	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/physics/body"
	"github.com/alexanderi96/go-space-engine/physics/space"
	"github.com/alexanderi96/go-space-engine/simulation/world"
)

// Color represents an RGBA color
type Color struct {
	R, G, B, A float64 // Color components (0-1)
}

// NewColor creates a new color
func NewColor(r, g, b, a float64) Color {
	return Color{
		R: r,
		G: g,
		B: b,
		A: a,
	}
}

// Renderer represents an interface for rendering
type Renderer interface {
	// Initialize initializes the renderer
	Initialize() error
	// Shutdown closes the renderer
	Shutdown() error

	// BeginFrame starts a new frame
	BeginFrame()
	// EndFrame ends the current frame
	EndFrame()

	// RenderBody renders a body
	RenderBody(b body.Body)
	// RenderBodies renders all bodies
	RenderBodies(bodies []body.Body)

	// RenderAABB renders an AABB
	RenderAABB(aabb *space.AABB, color Color)
	// RenderOctree renders an octree
	RenderOctree(octree *space.Octree, maxDepth int)

	// RenderLine renders a line
	RenderLine(start, end vector.Vector3, color Color)
	// RenderSphere renders a sphere
	RenderSphere(center vector.Vector3, radius float64, color Color)

	// SetCamera sets the camera position and orientation
	SetCamera(position, target, up vector.Vector3)
	// SetCameraFOV sets the camera field of view
	SetCameraFOV(fov float64)

	// SetBackgroundColor sets the background color
	SetBackgroundColor(color Color)

	// GetWidth returns the width of the rendering window
	GetWidth() int
	// GetHeight returns the height of the rendering window
	GetHeight() int

	// IsRunning returns true if the renderer is running
	IsRunning() bool

	// ProcessEvents processes renderer events
	ProcessEvents()
}

// RenderAdapter represents an adapter for rendering
type RenderAdapter interface {
	// GetRenderer returns the renderer
	GetRenderer() Renderer

	// RenderWorld renders the world
	RenderWorld(w world.World)

	// SetDebugMode sets the debug mode
	SetDebugMode(debug bool)
	// IsDebugMode returns true if debug mode is active
	IsDebugMode() bool

	// SetRenderOctree sets whether to render the octree
	SetRenderOctree(render bool)
	// IsRenderOctree returns true if the octree is being rendered
	IsRenderOctree() bool

	// SetRenderBoundingBoxes sets whether to render bounding boxes
	SetRenderBoundingBoxes(render bool)
	// IsRenderBoundingBoxes returns true if bounding boxes are being rendered
	IsRenderBoundingBoxes() bool

	// SetRenderVelocities sets whether to render velocity vectors
	SetRenderVelocities(render bool)
	// IsRenderVelocities returns true if velocity vectors are being rendered
	IsRenderVelocities() bool

	// SetRenderAccelerations sets whether to render acceleration vectors
	SetRenderAccelerations(render bool)
	// IsRenderAccelerations returns true if acceleration vectors are being rendered
	IsRenderAccelerations() bool

	// SetRenderForces sets whether to render force vectors
	SetRenderForces(render bool)
	// IsRenderForces returns true if force vectors are being rendered
	IsRenderForces() bool
}

// BaseRenderAdapter implements a base adapter for rendering
type BaseRenderAdapter struct {
	renderer            Renderer
	debugMode           bool
	renderOctree        bool
	renderBoundingBoxes bool
	renderVelocities    bool
	renderAccelerations bool
	renderForces        bool
}

// NewBaseRenderAdapter creates a new base adapter for rendering
func NewBaseRenderAdapter(renderer Renderer) *BaseRenderAdapter {
	return &BaseRenderAdapter{
		renderer:            renderer,
		debugMode:           false,
		renderOctree:        false,
		renderBoundingBoxes: false,
		renderVelocities:    false,
		renderAccelerations: false,
		renderForces:        false,
	}
}

// GetRenderer returns the renderer
func (ra *BaseRenderAdapter) GetRenderer() Renderer {
	return ra.renderer
}

// RenderWorld renders the world
func (ra *BaseRenderAdapter) RenderWorld(w world.World) {
	// Start a new frame
	ra.renderer.BeginFrame()

	// Render all bodies
	ra.renderer.RenderBodies(w.GetBodies())

	// Render world boundaries
	bounds := w.GetBounds()
	ra.renderer.RenderAABB(bounds, NewColor(0.5, 0.5, 0.5, 0.5))

	// Render the octree if requested
	if ra.renderOctree {
		if octree, ok := w.GetSpatialStructure().(*space.Octree); ok {
			ra.renderer.RenderOctree(octree, 8)
		}
	}

	// Render bounding boxes if requested
	if ra.renderBoundingBoxes {
		for _, b := range w.GetBodies() {
			position := b.Position()
			radius := b.Radius().Value()
			min := vector.NewVector3(position.X()-radius, position.Y()-radius, position.Z()-radius)
			max := vector.NewVector3(position.X()+radius, position.Y()+radius, position.Z()+radius)
			aabb := space.NewAABB(min, max)
			ra.renderer.RenderAABB(aabb, NewColor(0, 1, 0, 0.5))
		}
	}

	// Render velocity vectors if requested
	if ra.renderVelocities {
		for _, b := range w.GetBodies() {
			position := b.Position()
			velocity := b.Velocity()
			if velocity.Length() > 0.001 {
				end := position.Add(velocity)
				ra.renderer.RenderLine(position, end, NewColor(0, 0, 1, 1))
			}
		}
	}

	// Render acceleration vectors if requested
	if ra.renderAccelerations {
		for _, b := range w.GetBodies() {
			position := b.Position()
			acceleration := b.Acceleration()
			if acceleration.Length() > 0.001 {
				end := position.Add(acceleration)
				ra.renderer.RenderLine(position, end, NewColor(1, 0, 0, 1))
			}
		}
	}

	// Note: when using the Run method, rendering is handled internally
	// by G3N, so it's not necessary to call EndFrame here.
	// However, for compatibility with traditional usage, we call it anyway.
	// The EndFrame method has been modified to avoid conflicts.
	ra.renderer.EndFrame()
}

// SetDebugMode sets the debug mode
func (ra *BaseRenderAdapter) SetDebugMode(debug bool) {
	ra.debugMode = debug
}

// IsDebugMode returns true if debug mode is active
func (ra *BaseRenderAdapter) IsDebugMode() bool {
	return ra.debugMode
}

// SetRenderOctree sets whether to render the octree
func (ra *BaseRenderAdapter) SetRenderOctree(render bool) {
	ra.renderOctree = render
}

// IsRenderOctree returns true if the octree is being rendered
func (ra *BaseRenderAdapter) IsRenderOctree() bool {
	return ra.renderOctree
}

// SetRenderBoundingBoxes sets whether to render bounding boxes
func (ra *BaseRenderAdapter) SetRenderBoundingBoxes(render bool) {
	ra.renderBoundingBoxes = render
}

// IsRenderBoundingBoxes returns true if bounding boxes are being rendered
func (ra *BaseRenderAdapter) IsRenderBoundingBoxes() bool {
	return ra.renderBoundingBoxes
}

// SetRenderVelocities sets whether to render velocity vectors
func (ra *BaseRenderAdapter) SetRenderVelocities(render bool) {
	ra.renderVelocities = render
}

// IsRenderVelocities returns true if velocity vectors are being rendered
func (ra *BaseRenderAdapter) IsRenderVelocities() bool {
	return ra.renderVelocities
}

// SetRenderAccelerations sets whether to render acceleration vectors
func (ra *BaseRenderAdapter) SetRenderAccelerations(render bool) {
	ra.renderAccelerations = render
}

// IsRenderAccelerations returns true if acceleration vectors are being rendered
func (ra *BaseRenderAdapter) IsRenderAccelerations() bool {
	return ra.renderAccelerations
}

// SetRenderForces sets whether to render force vectors
func (ra *BaseRenderAdapter) SetRenderForces(render bool) {
	ra.renderForces = render
}

// IsRenderForces returns true if force vectors are being rendered
func (ra *BaseRenderAdapter) IsRenderForces() bool {
	return ra.renderForces
}
