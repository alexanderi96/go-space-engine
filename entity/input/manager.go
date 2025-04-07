// Package input provides interfaces and implementations for input handling
package input

// InputManager gestisce gli handler di input e il dispatching degli eventi
type InputManager struct {
	handlers []InputHandler
}

// NewInputManager crea un nuovo InputManager
func NewInputManager() *InputManager {
	return &InputManager{
		handlers: make([]InputHandler, 0),
	}
}

// RegisterInputHandler registra un handler di input
func (m *InputManager) RegisterInputHandler(handler InputHandler) {
	m.handlers = append(m.handlers, handler)
}

// UnregisterInputHandler rimuove un handler di input
func (m *InputManager) UnregisterInputHandler(handler InputHandler) {
	for i, h := range m.handlers {
		if h == handler {
			m.handlers = append(m.handlers[:i], m.handlers[i+1:]...)
			break
		}
	}
}

// DispatchEvent invia un evento a tutti gli handler registrati
func (m *InputManager) DispatchEvent(event InputEvent) {
	for _, handler := range m.handlers {
		handler.HandleInput(event)
	}
}
