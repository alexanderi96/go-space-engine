// Package input provides interfaces and implementations for input handling
package input

// InputEvent rappresenta un evento di input generico
type InputEvent interface {
	// GetType restituisce il tipo di evento (es. "key_down", "mouse_move")
	GetType() string

	// GetSource restituisce la sorgente dell'evento (es. l'oggetto window di g3n)
	GetSource() interface{}
}

// KeyEvent rappresenta un evento di tastiera
type KeyEvent struct {
	Key    int    // Codice del tasto
	Action int    // Azione (premuto, rilasciato)
	Mods   int    // Modificatori (shift, ctrl, alt)
	Type   string // Tipo di evento ("key_down", "key_up")
	Source interface{}
}

// Implementazione dell'interfaccia InputEvent
func (e *KeyEvent) GetType() string {
	return e.Type
}

func (e *KeyEvent) GetSource() interface{} {
	return e.Source
}

// MouseEvent rappresenta un evento del mouse
type MouseEvent struct {
	X      float64
	Y      float64
	Button int
	Action int
	Type   string
	Source interface{}
}

// Implementazione dell'interfaccia InputEvent
func (e *MouseEvent) GetType() string {
	return e.Type
}

func (e *MouseEvent) GetSource() interface{} {
	return e.Source
}

// InputHandler rappresenta un gestore di eventi di input
type InputHandler interface {
	// HandleInput gestisce un evento di input
	HandleInput(event InputEvent)
}
