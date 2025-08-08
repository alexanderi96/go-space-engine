// Package raylib provides an implementation of the RenderAdapter interface using Raylib
package raylib

import (
	"math"
	"time"

	"github.com/alexanderi96/go-space-engine/core/units"
	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/physics/body"
	"github.com/alexanderi96/go-space-engine/physics/space"
	"github.com/alexanderi96/go-space-engine/render/adapter"
	"github.com/alexanderi96/go-space-engine/simulation/world"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/google/uuid"
)

// RaylibRenderer implements the Renderer interface using Raylib
type RaylibRenderer struct {
	width            int32
	height           int32
	title            string
	isInitialized    bool
	isRunning        bool
	camera           rl.Camera3D
	bgColor          adapter.Color
	mouseCaptured    bool
	cameraSpeed      float32
	mouseSensitivity float32
}

// NewRaylibRenderer creates a new Raylib renderer
func NewRaylibRenderer(width, height int32, title string) *RaylibRenderer {
	return &RaylibRenderer{
		width:            width,
		height:           height,
		title:            title,
		bgColor:          adapter.NewColor(0.1, 0.1, 0.1, 1.0), // Dark background for space
		mouseCaptured:    false,
		cameraSpeed:      50.0,
		mouseSensitivity: 0.003,
		camera: rl.Camera3D{
			Position:   rl.Vector3{X: 0, Y: 50, Z: 150},
			Target:     rl.Vector3{X: 0, Y: 0, Z: 0},
			Up:         rl.Vector3{X: 0, Y: 1, Z: 0},
			Fovy:       45.0,
			Projection: rl.CameraPerspective,
		},
	}
}

// Initialize initializes the renderer
func (r *RaylibRenderer) Initialize() error {
	rl.InitWindow(r.width, r.height, r.title)
	rl.SetTargetFPS(60)

	// Don't disable cursor by default - let user control it
	r.isInitialized = true
	r.isRunning = true
	return nil
}

// Shutdown closes the renderer
func (r *RaylibRenderer) Shutdown() error {
	if r.isInitialized {
		rl.CloseWindow()
		r.isInitialized = false
		r.isRunning = false
	}
	return nil
}

// BeginFrame starts a new frame
func (r *RaylibRenderer) BeginFrame() {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Color{
		R: uint8(r.bgColor.R * 255),
		G: uint8(r.bgColor.G * 255),
		B: uint8(r.bgColor.B * 255),
		A: uint8(r.bgColor.A * 255),
	})
	rl.BeginMode3D(r.camera)
}

// EndFrame ends the current frame
func (r *RaylibRenderer) EndFrame() {
	rl.EndMode3D()

	// Draw UI elements (instructions)
	r.drawUI()

	rl.EndDrawing()

	// Handle input and update camera manually
	r.handleInput()

	// Check if window should close
	if rl.WindowShouldClose() {
		r.isRunning = false
	}
}

// drawUI draws the user interface
func (r *RaylibRenderer) drawUI() {
	// Draw instructions
	instructions := []string{
		"Controls:",
		"Right Click: Toggle mouse capture",
		"WASD: Move camera",
		"Mouse: Look around (when captured)",
		"Scroll: Zoom in/out",
		"ESC: Exit",
	}

	y := int32(10)
	for _, instruction := range instructions {
		rl.DrawText(instruction, 10, y, 20, rl.White)
		y += 25
	}

	// Show mouse capture status
	status := "Mouse: Free"
	if r.mouseCaptured {
		status = "Mouse: Captured"
	}
	rl.DrawText(status, 10, r.height-30, 20, rl.Yellow)
}

// handleInput handles keyboard and mouse input
func (r *RaylibRenderer) handleInput() {
	// Toggle mouse capture with right click
	if rl.IsMouseButtonPressed(rl.MouseRightButton) {
		r.mouseCaptured = !r.mouseCaptured
		if r.mouseCaptured {
			rl.DisableCursor()
		} else {
			rl.EnableCursor()
		}
	}

	// Camera movement with WASD
	deltaTime := rl.GetFrameTime()
	moveSpeed := r.cameraSpeed * deltaTime

	// Calculate camera direction vectors
	forward := rl.Vector3Subtract(r.camera.Target, r.camera.Position)
	forward = rl.Vector3Normalize(forward)

	right := rl.Vector3CrossProduct(forward, r.camera.Up)
	right = rl.Vector3Normalize(right)

	up := r.camera.Up

	// Movement
	if rl.IsKeyDown(rl.KeyW) {
		r.camera.Position = rl.Vector3Add(r.camera.Position, rl.Vector3Scale(forward, moveSpeed))
		r.camera.Target = rl.Vector3Add(r.camera.Target, rl.Vector3Scale(forward, moveSpeed))
	}
	if rl.IsKeyDown(rl.KeyS) {
		r.camera.Position = rl.Vector3Subtract(r.camera.Position, rl.Vector3Scale(forward, moveSpeed))
		r.camera.Target = rl.Vector3Subtract(r.camera.Target, rl.Vector3Scale(forward, moveSpeed))
	}
	if rl.IsKeyDown(rl.KeyA) {
		r.camera.Position = rl.Vector3Subtract(r.camera.Position, rl.Vector3Scale(right, moveSpeed))
		r.camera.Target = rl.Vector3Subtract(r.camera.Target, rl.Vector3Scale(right, moveSpeed))
	}
	if rl.IsKeyDown(rl.KeyD) {
		r.camera.Position = rl.Vector3Add(r.camera.Position, rl.Vector3Scale(right, moveSpeed))
		r.camera.Target = rl.Vector3Add(r.camera.Target, rl.Vector3Scale(right, moveSpeed))
	}
	if rl.IsKeyDown(rl.KeyQ) {
		r.camera.Position = rl.Vector3Subtract(r.camera.Position, rl.Vector3Scale(up, moveSpeed))
		r.camera.Target = rl.Vector3Subtract(r.camera.Target, rl.Vector3Scale(up, moveSpeed))
	}
	if rl.IsKeyDown(rl.KeyE) {
		r.camera.Position = rl.Vector3Add(r.camera.Position, rl.Vector3Scale(up, moveSpeed))
		r.camera.Target = rl.Vector3Add(r.camera.Target, rl.Vector3Scale(up, moveSpeed))
	}

	// Mouse look (only when captured)
	if r.mouseCaptured {
		mouseDelta := rl.GetMouseDelta()

		if mouseDelta.X != 0 || mouseDelta.Y != 0 {
			// Horizontal rotation (yaw)
			yaw := mouseDelta.X * r.mouseSensitivity

			// Vertical rotation (pitch)
			pitch := -mouseDelta.Y * r.mouseSensitivity

			// Apply rotations
			r.rotateCamera(yaw, pitch)
		}
	}

	// Zoom with mouse wheel
	wheelMove := rl.GetMouseWheelMove()
	if wheelMove != 0 {
		zoomSpeed := moveSpeed * 10
		zoom := rl.Vector3Scale(forward, wheelMove*zoomSpeed)
		r.camera.Position = rl.Vector3Add(r.camera.Position, zoom)
		r.camera.Target = rl.Vector3Add(r.camera.Target, zoom)
	}
}

// rotateCamera rotates the camera around its position
func (r *RaylibRenderer) rotateCamera(yaw, pitch float32) {
	// Simplified camera rotation using direct vector manipulation
	// Get the direction vector from camera to target
	direction := rl.Vector3Subtract(r.camera.Target, r.camera.Position)
	distance := rl.Vector3Length(direction)

	if distance < 0.001 {
		return // Avoid division by zero
	}

	// Normalize direction
	direction = rl.Vector3Normalize(direction)

	// Create rotation around Y axis (yaw)
	cosYaw := float32(math.Cos(float64(yaw)))
	sinYaw := float32(math.Sin(float64(yaw)))

	// Rotate direction vector around Y axis
	newX := direction.X*cosYaw - direction.Z*sinYaw
	newZ := direction.X*sinYaw + direction.Z*cosYaw
	direction.X = newX
	direction.Z = newZ

	// Apply pitch rotation (around the right vector)
	// Calculate right vector
	right := rl.Vector3CrossProduct(direction, r.camera.Up)
	right = rl.Vector3Normalize(right)

	// Rotate around right vector for pitch
	cosPitch := float32(math.Cos(float64(pitch)))
	sinPitch := float32(math.Sin(float64(pitch)))

	// Apply pitch rotation
	newDirection := rl.Vector3{
		X: direction.X*cosPitch + r.camera.Up.X*sinPitch,
		Y: direction.Y*cosPitch + r.camera.Up.Y*sinPitch,
		Z: direction.Z*cosPitch + r.camera.Up.Z*sinPitch,
	}

	// Clamp pitch to avoid flipping
	if newDirection.Y > 0.95 {
		newDirection.Y = 0.95
	}
	if newDirection.Y < -0.95 {
		newDirection.Y = -0.95
	}

	// Normalize and scale back to original distance
	newDirection = rl.Vector3Normalize(newDirection)
	r.camera.Target = rl.Vector3Add(r.camera.Position, rl.Vector3Scale(newDirection, distance))
}

// RenderBody renders a single body
func (r *RaylibRenderer) RenderBody(b body.Body) {
	pos := b.Position()

	// Convert radius to standard unit (meters) for consistent rendering
	radius := float32(units.ConvertToStandardUnit(b.Radius()))

	position := rl.Vector3{
		X: float32(pos.X()),
		Y: float32(pos.Y()),
		Z: float32(pos.Z()),
	}

	// Get color based on material
	color := r.getBodyColor(b)

	// Render based on body type
	if b.Material() != nil && b.Material().Name() == "Spacecraft" {
		// Render spacecraft as a cube
		rl.DrawCube(position, radius*2, radius*2, radius*2, color)
		rl.DrawCubeWires(position, radius*2, radius*2, radius*2, rl.White)
	} else {
		// Render other bodies as spheres
		rl.DrawSphere(position, radius, color)
		rl.DrawSphereWires(position, radius, 16, 16, rl.White)
	}
}

// RenderBodies renders all bodies
func (r *RaylibRenderer) RenderBodies(bodies []body.Body) {
	for _, b := range bodies {
		r.RenderBody(b)
	}
}

// RenderAABB renders an AABB (Axis-Aligned Bounding Box)
func (r *RaylibRenderer) RenderAABB(aabb *space.AABB, color adapter.Color) {
	min := aabb.Min
	max := aabb.Max

	rlColor := rl.Color{
		R: uint8(color.R * 255),
		G: uint8(color.G * 255),
		B: uint8(color.B * 255),
		A: uint8(color.A * 255),
	}

	// Calculate center and size
	center := rl.Vector3{
		X: float32((min.X() + max.X()) / 2),
		Y: float32((min.Y() + max.Y()) / 2),
		Z: float32((min.Z() + max.Z()) / 2),
	}

	size := rl.Vector3{
		X: float32(max.X() - min.X()),
		Y: float32(max.Y() - min.Y()),
		Z: float32(max.Z() - min.Z()),
	}

	rl.DrawCubeWires(center, size.X, size.Y, size.Z, rlColor)
}

// RenderOctree renders an octree
func (r *RaylibRenderer) RenderOctree(octree *space.Octree, maxDepth int) {
	// For now, we'll just render a placeholder
	// A full implementation would require access to octree internals
	// which are not exposed in the current interface
	_ = octree
	_ = maxDepth
}

// RenderLine renders a line
func (r *RaylibRenderer) RenderLine(start, end vector.Vector3, color adapter.Color) {
	startPos := rl.Vector3{
		X: float32(start.X()),
		Y: float32(start.Y()),
		Z: float32(start.Z()),
	}

	endPos := rl.Vector3{
		X: float32(end.X()),
		Y: float32(end.Y()),
		Z: float32(end.Z()),
	}

	rlColor := rl.Color{
		R: uint8(color.R * 255),
		G: uint8(color.G * 255),
		B: uint8(color.B * 255),
		A: uint8(color.A * 255),
	}

	rl.DrawLine3D(startPos, endPos, rlColor)
}

// RenderSphere renders a sphere
func (r *RaylibRenderer) RenderSphere(center vector.Vector3, radius float64, color adapter.Color) {
	position := rl.Vector3{
		X: float32(center.X()),
		Y: float32(center.Y()),
		Z: float32(center.Z()),
	}

	rlColor := rl.Color{
		R: uint8(color.R * 255),
		G: uint8(color.G * 255),
		B: uint8(color.B * 255),
		A: uint8(color.A * 255),
	}

	rl.DrawSphere(position, float32(radius), rlColor)
}

// RenderSphereWithUnit renders a sphere with a radius specified as a Quantity
// This allows for proper unit conversion
func (r *RaylibRenderer) RenderSphereWithUnit(center vector.Vector3, radius units.Quantity, color adapter.Color) {
	position := rl.Vector3{
		X: float32(center.X()),
		Y: float32(center.Y()),
		Z: float32(center.Z()),
	}

	rlColor := rl.Color{
		R: uint8(color.R * 255),
		G: uint8(color.G * 255),
		B: uint8(color.B * 255),
		A: uint8(color.A * 255),
	}

	// Convert radius to standard unit (meters)
	standardRadius := units.ConvertToStandardUnit(radius)

	rl.DrawSphere(position, float32(standardRadius), rlColor)
}

// SetCamera sets the camera position and orientation
func (r *RaylibRenderer) SetCamera(position, target, up vector.Vector3) {
	r.camera.Position = rl.Vector3{
		X: float32(position.X()),
		Y: float32(position.Y()),
		Z: float32(position.Z()),
	}
	r.camera.Target = rl.Vector3{
		X: float32(target.X()),
		Y: float32(target.Y()),
		Z: float32(target.Z()),
	}
	r.camera.Up = rl.Vector3{
		X: float32(up.X()),
		Y: float32(up.Y()),
		Z: float32(up.Z()),
	}
}

// SetCameraFOV sets the camera field of view
func (r *RaylibRenderer) SetCameraFOV(fov float64) {
	r.camera.Fovy = float32(fov)
}

// SetBackgroundColor sets the background color
func (r *RaylibRenderer) SetBackgroundColor(color adapter.Color) {
	r.bgColor = color
}

// GetWidth returns the width of the rendering window
func (r *RaylibRenderer) GetWidth() int {
	return int(r.width)
}

// GetHeight returns the height of the rendering window
func (r *RaylibRenderer) GetHeight() int {
	return int(r.height)
}

// IsRunning returns true if the renderer is running
func (r *RaylibRenderer) IsRunning() bool {
	return r.isRunning && r.isInitialized && !rl.WindowShouldClose()
}

// ProcessEvents processes renderer events
func (r *RaylibRenderer) ProcessEvents() {
	// Raylib handles events internally, but we can add custom logic here if needed
}

// getBodyColor returns the appropriate color for a body based on its material
func (r *RaylibRenderer) getBodyColor(b body.Body) rl.Color {
	if b.Material() == nil {
		return rl.Gray
	}

	switch b.Material().Name() {
	case "Spacecraft":
		return rl.White
	case "Sun":
		return rl.Yellow
	case "Iron":
		return rl.Color{R: 153, G: 153, B: 153, A: 255}
	case "Rock":
		return rl.Color{R: 128, G: 77, B: 51, A: 255}
	case "Ice":
		return rl.Color{R: 204, G: 230, B: 255, A: 255}
	case "Copper":
		return rl.Color{R: 204, G: 128, B: 51, A: 255}
	case "Mercury":
		return rl.Color{R: 179, G: 179, B: 179, A: 255}
	case "Venus":
		return rl.Color{R: 230, G: 179, B: 0, A: 255}
	case "Earth":
		return rl.Color{R: 0, G: 77, B: 204, A: 255}
	case "Mars":
		return rl.Color{R: 204, G: 77, B: 0, A: 255}
	case "Jupiter":
		return rl.Color{R: 204, G: 153, B: 102, A: 255}
	case "Saturn":
		return rl.Color{R: 230, G: 204, B: 128, A: 255}
	case "Uranus":
		return rl.Color{R: 128, G: 204, B: 230, A: 255}
	case "Neptune":
		return rl.Color{R: 0, G: 0, B: 204, A: 255}
	default:
		// Generate color based on body ID
		id := b.ID()
		hash := 0
		for i := 0; i < len(id); i++ {
			hash = 31*hash + int(id[i])
		}
		if hash < 0 {
			hash = -hash
		}
		return rl.Color{
			R: uint8(hash % 255),
			G: uint8((hash / 255) % 255),
			B: uint8((hash / (255 * 255)) % 255),
			A: 255,
		}
	}
}

// RaylibAdapter is an adapter for rendering with Raylib
type RaylibAdapter struct {
	*adapter.BaseRenderAdapter
	renderer   *RaylibRenderer
	bodyMeshes map[uuid.UUID]bool // Track which bodies have been rendered
}

// NewRaylibAdapter creates a new Raylib adapter
func NewRaylibAdapter(width, height int32, title string) *RaylibAdapter {
	renderer := NewRaylibRenderer(width, height, title)
	baseAdapter := adapter.NewBaseRenderAdapter(renderer)

	return &RaylibAdapter{
		BaseRenderAdapter: baseAdapter,
		renderer:          renderer,
		bodyMeshes:        make(map[uuid.UUID]bool),
	}
}

// GetRenderer returns the renderer
func (ra *RaylibAdapter) GetRenderer() adapter.Renderer {
	return ra.renderer
}

// Initialize initializes the adapter
func (ra *RaylibAdapter) Initialize() error {
	return ra.renderer.Initialize()
}

// Shutdown shuts down the adapter
func (ra *RaylibAdapter) Shutdown() error {
	return ra.renderer.Shutdown()
}

// Run starts the rendering loop with a custom update function
func (ra *RaylibAdapter) Run(updateFunc func(deltaTime time.Duration)) {
	if !ra.renderer.isInitialized {
		ra.Initialize()
	}

	lastTime := time.Now()

	for ra.renderer.IsRunning() {
		currentTime := time.Now()
		deltaTime := currentTime.Sub(lastTime)
		lastTime = currentTime

		// Call the update function
		if updateFunc != nil {
			updateFunc(deltaTime)
		}

		// Process events
		ra.renderer.ProcessEvents()
	}

	ra.Shutdown()
}

// RenderWorld renders the world with enhanced features
func (ra *RaylibAdapter) RenderWorld(w world.World) {
	// Start frame
	ra.renderer.BeginFrame()

	// Render all bodies
	bodies := w.GetBodies()
	ra.renderer.RenderBodies(bodies)

	// Track rendered bodies
	for _, b := range bodies {
		ra.bodyMeshes[b.ID()] = true
	}

	// Render world boundaries if in debug mode
	if ra.IsDebugMode() {
		bounds := w.GetBounds()
		ra.renderer.RenderAABB(bounds, adapter.NewColor(0.5, 0.5, 0.5, 0.5))
	}

	// Render octree if requested
	if ra.IsRenderOctree() {
		if octree, ok := w.GetSpatialStructure().(*space.Octree); ok {
			ra.renderer.RenderOctree(octree, 8)
		}
	}

	// Render bounding boxes if requested
	if ra.IsRenderBoundingBoxes() {
		for _, b := range bodies {
			position := b.Position()
			radius := b.Radius().Value()
			min := vector.NewVector3(position.X()-radius, position.Y()-radius, position.Z()-radius)
			max := vector.NewVector3(position.X()+radius, position.Y()+radius, position.Z()+radius)
			aabb := space.NewAABB(min, max)
			ra.renderer.RenderAABB(aabb, adapter.NewColor(0, 1, 0, 0.5))
		}
	}

	// Render velocity vectors if requested
	if ra.IsRenderVelocities() {
		for _, b := range bodies {
			position := b.Position()
			velocity := b.Velocity()
			if velocity.Length() > 0.001 {
				// Scale velocity for visibility
				scaledVelocity := velocity.Scale(10.0)
				end := position.Add(scaledVelocity)
				ra.renderer.RenderLine(position, end, adapter.NewColor(0, 0, 1, 1))
			}
		}
	}

	// Render acceleration vectors if requested
	if ra.IsRenderAccelerations() {
		for _, b := range bodies {
			position := b.Position()
			acceleration := b.Acceleration()
			if acceleration.Length() > 0.001 {
				// Scale acceleration for visibility
				scaledAcceleration := acceleration.Scale(100.0)
				end := position.Add(scaledAcceleration)
				ra.renderer.RenderLine(position, end, adapter.NewColor(1, 0, 0, 1))
			}
		}
	}

	// Render force vectors if requested
	if ra.IsRenderForces() {
		for _, b := range bodies {
			_ = b // Avoid unused variable error
			// Note: We would need to access forces from the body
			// This might require extending the body interface
			// For now, we'll skip this implementation
		}
	}

	// End frame
	ra.renderer.EndFrame()
}

// SetBackgroundColor sets the background color
func (ra *RaylibAdapter) SetBackgroundColor(color adapter.Color) {
	ra.renderer.SetBackgroundColor(color)
}

// GetCamera returns the camera (Raylib-specific method)
func (ra *RaylibAdapter) GetCamera() *rl.Camera3D {
	return &ra.renderer.camera
}

// SetCameraPosition sets the camera position
func (ra *RaylibAdapter) SetCameraPosition(position vector.Vector3) {
	ra.renderer.camera.Position = rl.Vector3{
		X: float32(position.X()),
		Y: float32(position.Y()),
		Z: float32(position.Z()),
	}
}

// SetCameraTarget sets the camera target
func (ra *RaylibAdapter) SetCameraTarget(target vector.Vector3) {
	ra.renderer.camera.Target = rl.Vector3{
		X: float32(target.X()),
		Y: float32(target.Y()),
		Z: float32(target.Z()),
	}
}
