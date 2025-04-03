// Package g3n provides an implementation of the RenderAdapter interface using G3N
package g3n

import (
	"time"

	"github.com/alexanderi96/go-space-engine/physics/body"
	"github.com/alexanderi96/go-space-engine/render/adapter"
	"github.com/alexanderi96/go-space-engine/simulation/world"
	"github.com/google/uuid"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/window"
)

// BodyMesh represents a mesh with an associated point light
type BodyMesh struct {
	Mesh  *graphic.Mesh
	Light *light.Point
}

// G3NAdapter is an adapter for rendering with G3N
type G3NAdapter struct {
	app        *app.Application
	scene      *core.Node
	camera     *camera.Camera
	cameraCtrl *camera.OrbitControl
	bodyMeshes map[uuid.UUID]*BodyMesh
	bgColor    adapter.Color
	debugMode  bool
}

// NewG3NAdapter creates a new G3N adapter
func NewG3NAdapter() *G3NAdapter {
	return &G3NAdapter{
		bodyMeshes: make(map[uuid.UUID]*BodyMesh),
		bgColor:    adapter.NewColor(1.0, 1.0, 1.0, 1.0), // White background
		debugMode:  false,
	}
}

// GetRenderer returns the renderer (implementation of the RenderAdapter interface)
func (ga *G3NAdapter) GetRenderer() adapter.Renderer {
	// This adapter does not use the standard Renderer interface
	// Returns nil because it directly implements the necessary methods
	return nil
}

// RenderWorld renders the world
func (ga *G3NAdapter) RenderWorld(w world.World) {
	// Update the position of meshes
	for _, b := range w.GetBodies() {
		if bodyMesh, exists := ga.bodyMeshes[b.ID()]; exists {
			pos := b.Position()
			bodyMesh.Mesh.SetPosition(float32(pos.X()), float32(pos.Y()), float32(pos.Z()))
			if bodyMesh.Light != nil {
				bodyMesh.Light.SetPosition(float32(pos.X()), float32(pos.Y()), float32(pos.Z()))
			}
		} else {
			// If the body does not have an associated mesh, create it
			ga.createMeshForBody(b)
		}
	}
}

// Run starts the rendering loop
func (ga *G3NAdapter) Run(updateFunc func(deltaTime time.Duration)) {
	// Initialize the G3N application if it has not already been initialized
	if ga.app == nil {
		ga.initialize()
	}

	// Start the rendering loop
	ga.app.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		// Explicit OpenGL operations
		gl := ga.app.Gls()
		gl.Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
		gl.Enable(gls.DEPTH_TEST)

		// Call the provided update function
		if updateFunc != nil {
			updateFunc(deltaTime)
		}

		// Render the scene
		renderer.Render(ga.scene, ga.camera)

		// Disable depth testing after rendering
		gl.Disable(gls.DEPTH_TEST)
	})
}

// initialize initializes the adapter
func (ga *G3NAdapter) initialize() {
	// Create the G3N application
	ga.app = app.App()

	// Create the scene
	ga.scene = core.NewNode()

	// Create the camera
	ga.camera = camera.New(1)
	ga.camera.SetPosition(0, 50, 150)
	ga.camera.LookAt(&math32.Vector3{0, 0, 0}, &math32.Vector3{0, 1, 0})
	ga.scene.Add(ga.camera)

	// Create the orbital camera control
	ga.cameraCtrl = camera.NewOrbitControl(ga.camera)

	// Set the background color
	ga.app.Gls().ClearColor(float32(ga.bgColor.R), float32(ga.bgColor.G), float32(ga.bgColor.B), float32(ga.bgColor.A))

	// Add a handler for window resizing
	ga.app.Subscribe(window.OnWindowSize, ga.onWindowResize)

	// Set the initial aspect ratio of the camera
	width, height := ga.app.GetSize()
	aspect := float32(width) / float32(height)
	ga.camera.SetAspect(aspect)

	// Add lights
	// Softer ambient light for a space effect
	ambLight := light.NewAmbient(&math32.Color{0.3, 0.3, 0.4}, 0.5)
	ga.scene.Add(ambLight)

	// More intense and distant point lights to illuminate the entire solar system
	pointLight1 := light.NewPoint(&math32.Color{1, 1, 1}, 5.0)
	pointLight1.SetPosition(50, 50, 50)
	pointLight1.SetLinearDecay(0.1)
	pointLight1.SetQuadraticDecay(0.01)
	ga.scene.Add(pointLight1)

	pointLight2 := light.NewPoint(&math32.Color{1, 1, 1}, 5.0)
	pointLight2.SetPosition(-50, 50, 50)
	pointLight2.SetLinearDecay(0.1)
	pointLight2.SetQuadraticDecay(0.01)
	ga.scene.Add(pointLight2)

	pointLight3 := light.NewPoint(&math32.Color{1, 1, 1}, 5.0)
	pointLight3.SetPosition(0, 50, -50)
	pointLight3.SetLinearDecay(0.1)
	pointLight3.SetQuadraticDecay(0.01)
	ga.scene.Add(pointLight3)

	// Create axes
	// We remove the axes and grid for a cleaner visualization of the solar system
}

// createMeshForBody creates a mesh for a physical body
func (ga *G3NAdapter) createMeshForBody(b body.Body) {
	// Create a sphere to represent the body
	radius := float32(b.Radius().Value())

	// Increase the quality of spheres for larger bodies
	var segments, rings int
	if radius > 1.5 {
		segments, rings = 64, 32 // High quality for large planets
	} else if radius > 0.8 {
		segments, rings = 48, 24 // Medium quality for medium planets
	} else {
		segments, rings = 32, 16 // Standard quality for small bodies
	}

	geom := geometry.NewSphere(float64(radius), segments, rings)

	// Create a material based on the physical body's material
	var mat material.IMaterial
	var bodyColor math32.Color

	// Determine the color of the body
	if b.Material() != nil {
		// Here you should map the physical material to a G3N color
		// For simplicity, we use a predefined color for each type of material
		switch b.Material().Name() {
		case "Sun":
			// Special material for the sun with emission
			bodyColor = math32.Color{1.0, 0.8, 0.0}
			// Create a more intense emission color to simulate the brightness of the sun
			emissiveColor := math32.Color{1.0, 0.9, 0.5}
			sunMat := material.NewStandard(&bodyColor)
			sunMat.SetEmissiveColor(&emissiveColor)
			sunMat.SetOpacity(1.0)
			mat = sunMat
		case "Iron":
			bodyColor = math32.Color{0.6, 0.6, 0.6}
		case "Rock":
			bodyColor = math32.Color{0.5, 0.3, 0.2}
		case "Ice":
			bodyColor = math32.Color{0.8, 0.9, 1.0}
		case "Copper":
			bodyColor = math32.Color{0.8, 0.5, 0.2}
		case "Mercury":
			bodyColor = math32.Color{0.7, 0.7, 0.7}
		case "Venus":
			bodyColor = math32.Color{0.9, 0.7, 0.0}
		case "Earth":
			bodyColor = math32.Color{0.0, 0.3, 0.8}
		case "Mars":
			bodyColor = math32.Color{0.8, 0.3, 0.0}
		case "Jupiter":
			bodyColor = math32.Color{0.8, 0.6, 0.4}
		case "Saturn":
			bodyColor = math32.Color{0.9, 0.8, 0.5}
		case "Uranus":
			bodyColor = math32.Color{0.5, 0.8, 0.9}
		case "Neptune":
			bodyColor = math32.Color{0.0, 0.0, 0.8}
		default:
			// Random color based on the body's ID (which is a string)
			id := b.ID()
			hash := 0
			for i := 0; i < len(id); i++ {
				hash = 31*hash + int(id[i])
			}
			if hash < 0 {
				hash = -hash
			}
			r := float32(hash%255) / 255.0
			g := float32((hash/255)%255) / 255.0
			b := float32((hash/(255*255))%255) / 255.0
			bodyColor = math32.Color{r, g, b}
		}

		// If the material has not already been created (as for the sun)
		if mat == nil {
			standardMat := material.NewStandard(&bodyColor)
			standardMat.SetShininess(30)
			mat = standardMat
		}
	} else {
		bodyColor = math32.Color{0.8, 0.8, 0.8}
		mat = material.NewStandard(&bodyColor)
	}

	// Create a mesh with the geometry and material
	mesh := graphic.NewMesh(geom, mat)

	// Set the position of the mesh
	pos := b.Position()
	mesh.SetPosition(float32(pos.X()), float32(pos.Y()), float32(pos.Z()))

	// Create a point light for the body only if it is large enough
	var bodyLight *light.Point
	if radius > 0.5 || b.Material().Name() == "Sun" {
		// Light intensity proportional to the size of the body
		lightIntensity := float32(0.5)
		if b.Material().Name() == "Sun" {
			lightIntensity = 5.0 // The sun is much brighter
		} else if radius > 1.5 {
			lightIntensity = 1.0 // Large planets are brighter
		}

		bodyLight = light.NewPoint(&bodyColor, lightIntensity)
		bodyLight.SetPosition(float32(pos.X()), float32(pos.Y()), float32(pos.Z()))

		// More gradual light decay for the sun
		if b.Material().Name() == "Sun" {
			bodyLight.SetLinearDecay(0.05)
			bodyLight.SetQuadraticDecay(0.005)
		} else {
			bodyLight.SetLinearDecay(0.5)
			bodyLight.SetQuadraticDecay(0.5)
		}

		ga.scene.Add(bodyLight)
	}

	// Add the mesh to the scene
	ga.scene.Add(mesh)

	// Store the BodyMesh in the map
	ga.bodyMeshes[b.ID()] = &BodyMesh{
		Mesh:  mesh,
		Light: bodyLight,
	}
}

// SetDebugMode sets the debug mode
func (ga *G3NAdapter) SetDebugMode(debug bool) {
	ga.debugMode = debug
}

// IsDebugMode returns true if debug mode is active
func (ga *G3NAdapter) IsDebugMode() bool {
	return ga.debugMode
}

// SetRenderOctree sets whether to render the octree
func (ga *G3NAdapter) SetRenderOctree(render bool) {
	// Not implemented in this adapter
}

// IsRenderOctree returns true if the octree is being rendered
func (ga *G3NAdapter) IsRenderOctree() bool {
	return false
}

// SetRenderBoundingBoxes sets whether to render bounding boxes
func (ga *G3NAdapter) SetRenderBoundingBoxes(render bool) {
	// Not implemented in this adapter
}

// IsRenderBoundingBoxes returns true if bounding boxes are being rendered
func (ga *G3NAdapter) IsRenderBoundingBoxes() bool {
	return false
}

// SetRenderVelocities sets whether to render velocity vectors
func (ga *G3NAdapter) SetRenderVelocities(render bool) {
	// Not implemented in this adapter
}

// IsRenderVelocities returns true if velocity vectors are being rendered
func (ga *G3NAdapter) IsRenderVelocities() bool {
	return false
}

// SetRenderAccelerations sets whether to render acceleration vectors
func (ga *G3NAdapter) SetRenderAccelerations(render bool) {
	// Not implemented in this adapter
}

// IsRenderAccelerations returns true if acceleration vectors are being rendered
func (ga *G3NAdapter) IsRenderAccelerations() bool {
	return false
}

// SetRenderForces sets whether to render force vectors
func (ga *G3NAdapter) SetRenderForces(render bool) {
	// Not implemented in this adapter
}

// IsRenderForces returns true if force vectors are being rendered
func (ga *G3NAdapter) IsRenderForces() bool {
	return false
}

// SetBackgroundColor sets the background color
func (ga *G3NAdapter) SetBackgroundColor(color adapter.Color) {
	ga.bgColor = color
	if ga.app != nil {
		ga.app.Gls().ClearColor(float32(color.R), float32(color.G), float32(color.B), float32(color.A))
	}
}

// NewColor creates a new color
func NewColor(r, g, b, a float64) adapter.Color {
	return adapter.NewColor(r, g, b, a)
}

// onWindowResize handles window resizing
func (ga *G3NAdapter) onWindowResize(evname string, ev interface{}) {
	// Get the new window dimensions
	width, height := ga.app.GetSize()

	// Update the camera's aspect ratio
	aspect := float32(width) / float32(height)
	ga.camera.SetAspect(aspect)

	// Explicitly set the viewport using the OpenGL API
	gl := ga.app.Gls()
	gl.Viewport(0, 0, int32(width), int32(height))
}
