package main

import (
	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/entity/input"
	"github.com/alexanderi96/go-space-engine/entity/vehicle/spacecraft"
)

// SpacecraftInputHandler gestisce gli input per il controllo della navicella
type SpacecraftInputHandler struct {
	spacecraft  *spacecraft.SpacecraftController
	thrustLevel float64
}

// NewSpacecraftInputHandler crea un nuovo handler di input per la navicella
func NewSpacecraftInputHandler(spacecraft *spacecraft.SpacecraftController) *SpacecraftInputHandler {
	return &SpacecraftInputHandler{
		spacecraft:  spacecraft,
		thrustLevel: 0.0,
	}
}

// HandleInput gestisce un evento di input
func (h *SpacecraftInputHandler) HandleInput(event input.InputEvent) {
	// Gestisce solo gli eventi di tastiera
	if event.GetType() != input.EventKeyDown && event.GetType() != input.EventKeyUp {
		return
	}

	// Type assertion per ottenere l'evento di tastiera
	keyEvent, ok := event.(*input.KeyEvent)
	if !ok {
		return
	}

	// Gestione degli eventi di pressione dei tasti
	if event.GetType() == input.EventKeyDown {
		switch keyEvent.Key {
		case input.KeyW: // Pitch up
			h.spacecraft.ApplyRotation(vector.NewVector3(1, 0, 0), 1.0)
		case input.KeyS: // Pitch down
			h.spacecraft.ApplyRotation(vector.NewVector3(1, 0, 0), -1.0)
		case input.KeyA: // Yaw left
			h.spacecraft.ApplyRotation(vector.NewVector3(0, 1, 0), 1.0)
		case input.KeyD: // Yaw right
			h.spacecraft.ApplyRotation(vector.NewVector3(0, 1, 0), -1.0)
		case input.KeyQ: // Roll left
			h.spacecraft.ApplyRotation(vector.NewVector3(0, 0, 1), 1.0)
		case input.KeyE: // Roll right
			h.spacecraft.ApplyRotation(vector.NewVector3(0, 0, 1), -1.0)
		case input.KeySpace: // Toggle thrust
			if h.thrustLevel <= 0.0 {
				h.thrustLevel = 0.5 // 50% thrust
			} else {
				h.thrustLevel = 0.0 // No thrust
			}
			h.spacecraft.SetThrustLevel(h.thrustLevel)
		}
	}
}

// Update aggiorna lo stato dell'handler
func (h *SpacecraftInputHandler) Update(deltaTime float64) {
	// Aggiorna il controller della navicella
	h.spacecraft.Update(deltaTime)
}
