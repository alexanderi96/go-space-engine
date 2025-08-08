// Package raylib provides an implementation of the RenderAdapter interface using Raylib
package raylib

import (
	"math"
	"time"

	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/physics/body"
	"github.com/alexanderi96/go-space-engine/render/adapter"
	"github.com/alexanderi96/go-space-engine/simulation/world"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// ControllableRaylibAdapter extends RaylibAdapter to support controllable bodies
type ControllableRaylibAdapter struct {
	*RaylibAdapter
	controllableBody body.ControllableBody
	controlMode      bool // true = control body, false = control camera
}

// NewControllableRaylibAdapter creates a new controllable Raylib adapter
func NewControllableRaylibAdapter(width, height int32, title string) *ControllableRaylibAdapter {
	adapter := NewRaylibAdapter(width, height, title)
	return &ControllableRaylibAdapter{
		RaylibAdapter:    adapter,
		controllableBody: nil,
		controlMode:      false, // Default to camera control
	}
}

// SetControllableBody sets the controllable body
func (cra *ControllableRaylibAdapter) SetControllableBody(b body.ControllableBody) {
	cra.controllableBody = b
}

// GetControllableBody returns the controllable body
func (cra *ControllableRaylibAdapter) GetControllableBody() body.ControllableBody {
	return cra.controllableBody
}

// ToggleControlMode toggles between camera control and body control
func (cra *ControllableRaylibAdapter) ToggleControlMode() {
	cra.controlMode = !cra.controlMode

	// Se passiamo alla modalità controllo navicella, posiziona la camera dietro di essa
	if cra.controlMode && cra.controllableBody != nil {
		cra.updateCameraToFollowSpacecraft()
	}
}

// updateCameraToFollowSpacecraft posiziona la camera dietro la navicella
func (cra *ControllableRaylibAdapter) updateCameraToFollowSpacecraft() {
	if cra.controllableBody == nil {
		return
	}

	// Ottieni la posizione e la rotazione della navicella
	spacecraftPos := cra.controllableBody.Position()
	spacecraftRot := cra.controllableBody.Rotation()

	// Calcola la direzione in cui la navicella sta guardando
	// Assumiamo che la direzione iniziale sia (0, 0, -1) e applichiamo la rotazione
	forward := vector.NewVector3(0, 0, -1)

	// Applica la rotazione Y (yaw)
	cosY := math.Cos(spacecraftRot.Y())
	sinY := math.Sin(spacecraftRot.Y())
	rotatedForward := vector.NewVector3(
		forward.X()*cosY-forward.Z()*sinY,
		forward.Y(),
		forward.X()*sinY+forward.Z()*cosY,
	)

	// Calcola la posizione della camera (dietro la navicella)
	// Distanza della camera dalla navicella
	cameraDistance := 5.0

	// Posizione della camera = posizione navicella - direzione * distanza
	cameraPos := spacecraftPos.Add(rotatedForward.Scale(-cameraDistance))

	// Aggiungi un po' di altezza alla camera per una visuale migliore
	cameraPos = vector.NewVector3(cameraPos.X(), cameraPos.Y()+2.0, cameraPos.Z())

	// Imposta la posizione e il target della camera
	cra.SetCameraPosition(cameraPos)
	cra.SetCameraTarget(spacecraftPos)
}

// IsControllingBody returns true if controlling the body
func (cra *ControllableRaylibAdapter) IsControllingBody() bool {
	return cra.controlMode
}

// HandleControllableBodyInput handles input for the controllable body
func (cra *ControllableRaylibAdapter) HandleControllableBodyInput() {
	if cra.controllableBody == nil || !cra.controllableBody.IsControllable() {
		return
	}

	// Get delta time
	deltaTime := rl.GetFrameTime()

	// Detect key presses for movement
	moveForward := rl.IsKeyDown(rl.KeyW)
	moveBackward := rl.IsKeyDown(rl.KeyS)
	moveLeft := rl.IsKeyDown(rl.KeyA)
	moveRight := rl.IsKeyDown(rl.KeyD)
	moveUp := rl.IsKeyDown(rl.KeyE)
	moveDown := rl.IsKeyDown(rl.KeyQ)

	// Detect key presses for rotation
	rotateLeft := rl.IsKeyDown(rl.KeyLeft)
	rotateRight := rl.IsKeyDown(rl.KeyRight)
	rotateUp := rl.IsKeyDown(rl.KeyUp)
	rotateDown := rl.IsKeyDown(rl.KeyDown)

	// Pass input to the controllable body
	cra.controllableBody.HandleInput(float64(deltaTime), moveForward, moveBackward, moveLeft, moveRight, moveUp, moveDown, rotateLeft, rotateRight, rotateUp, rotateDown)
}

// ProcessEvents processes renderer events
func (cra *ControllableRaylibAdapter) ProcessEvents() {
	// Toggle control mode with Tab key
	if rl.IsKeyPressed(rl.KeyTab) {
		cra.ToggleControlMode()
	}

	// Handle input based on control mode
	if cra.IsControllingBody() {
		// Handle input for the controllable body
		cra.HandleControllableBodyInput()

		// Aggiorna la posizione della camera per seguire la navicella
		cra.updateCameraToFollowSpacecraft()
	} else {
		// Handle input for the camera (default behavior)
		cra.renderer.handleInput()
	}
}

// Run starts the rendering loop with a custom update function
func (cra *ControllableRaylibAdapter) Run(updateFunc func(deltaTime time.Duration)) {
	if !cra.renderer.isInitialized {
		cra.Initialize()
	}

	lastTime := time.Now()

	for cra.renderer.IsRunning() {
		currentTime := time.Now()
		deltaTime := currentTime.Sub(lastTime)
		lastTime = currentTime

		// Process events (including input handling)
		cra.ProcessEvents()

		// Call the update function
		if updateFunc != nil {
			updateFunc(deltaTime)
		}
	}

	cra.Shutdown()
}

// RenderWorld renders the world with enhanced features
func (cra *ControllableRaylibAdapter) RenderWorld(w world.World) {
	// Start frame
	cra.renderer.BeginFrame()

	// Render all bodies
	bodies := w.GetBodies()
	cra.renderer.RenderBodies(bodies)

	// Track rendered bodies
	for _, b := range bodies {
		cra.bodyMeshes[b.ID()] = true
	}

	// Highlight the controllable body if in control mode
	if cra.IsControllingBody() && cra.controllableBody != nil {
		position := cra.controllableBody.Position()
		radius := cra.controllableBody.Radius().Value()

		// Draw a yellow outline around the controllable body
		cra.renderer.RenderSphere(position, radius*1.1, adapter.NewColor(1.0, 1.0, 0.0, 0.3))
	}

	// Render debug information if requested
	if cra.IsDebugMode() {
		bounds := w.GetBounds()
		cra.renderer.RenderAABB(bounds, adapter.NewColor(0.5, 0.5, 0.5, 0.5))
	}

	// Render other debug information (same as RaylibAdapter)
	// ...

	// End frame
	cra.renderer.EndFrame()

	// Draw UI elements
	cra.DrawUI()
}

// DrawUI draws the user interface
func (cra *ControllableRaylibAdapter) DrawUI() {
	// Draw control mode indicator
	modeText := "Camera Control (Tab to switch)"
	if cra.IsControllingBody() {
		modeText = "Spacecraft Control (Tab to switch)"
	}
	rl.DrawText(modeText, 10, 10, 20, rl.Yellow)

	// Draw controls
	if cra.IsControllingBody() {
		controls := []string{
			"W/S: Move forward/backward",
			"A/D: Move left/right",
			"Q/E: Move down/up",
			"Arrows: Rotate",
		}

		y := int32(40)
		for _, control := range controls {
			rl.DrawText(control, 10, y, 20, rl.White)
			y += 25
		}
	}
}
