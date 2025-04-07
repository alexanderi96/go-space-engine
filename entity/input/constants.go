// Package input provides interfaces and implementations for input handling
package input

// Tipi di eventi
const (
	EventKeyDown    = "key_down"
	EventKeyUp      = "key_up"
	EventMouseMove  = "mouse_move"
	EventMouseDown  = "mouse_down"
	EventMouseUp    = "mouse_up"
	EventMouseWheel = "mouse_wheel"
)

// Codici dei tasti (compatibili con quelli di g3n/window)
const (
	// Lettere
	KeyA = 65
	KeyB = 66
	KeyC = 67
	KeyD = 68
	KeyE = 69
	KeyF = 70
	KeyG = 71
	KeyH = 72
	KeyI = 73
	KeyJ = 74
	KeyK = 75
	KeyL = 76
	KeyM = 77
	KeyN = 78
	KeyO = 79
	KeyP = 80
	KeyQ = 81
	KeyR = 82
	KeyS = 83
	KeyT = 84
	KeyU = 85
	KeyV = 86
	KeyW = 87
	KeyX = 88
	KeyY = 89
	KeyZ = 90

	// Numeri
	Key0 = 48
	Key1 = 49
	Key2 = 50
	Key3 = 51
	Key4 = 52
	Key5 = 53
	Key6 = 54
	Key7 = 55
	Key8 = 56
	Key9 = 57

	// Tasti funzione
	KeyF1  = 290
	KeyF2  = 291
	KeyF3  = 292
	KeyF4  = 293
	KeyF5  = 294
	KeyF6  = 295
	KeyF7  = 296
	KeyF8  = 297
	KeyF9  = 298
	KeyF10 = 299
	KeyF11 = 300
	KeyF12 = 301

	// Tasti speciali
	KeySpace        = 32
	KeyApostrophe   = 39
	KeyComma        = 44
	KeyMinus        = 45
	KeyPeriod       = 46
	KeySlash        = 47
	KeySemicolon    = 59
	KeyEqual        = 61
	KeyLeftBracket  = 91
	KeyBackslash    = 92
	KeyRightBracket = 93
	KeyGraveAccent  = 96
	KeyEscape       = 256
	KeyEnter        = 257
	KeyTab          = 258
	KeyBackspace    = 259
	KeyInsert       = 260
	KeyDelete       = 261
	KeyRight        = 262
	KeyLeft         = 263
	KeyDownArrow    = 264
	KeyUpArrow      = 265
	KeyPageUp       = 266
	KeyPageDown     = 267
	KeyHome         = 268
	KeyEnd          = 269
	KeyCapsLock     = 280
	KeyScrollLock   = 281
	KeyNumLock      = 282
	KeyPrintScreen  = 283
	KeyPause        = 284
	KeyLeftShift    = 340
	KeyLeftControl  = 341
	KeyLeftAlt      = 342
	KeyLeftSuper    = 343
	KeyRightShift   = 344
	KeyRightControl = 345
	KeyRightAlt     = 346
	KeyRightSuper   = 347
	KeyMenu         = 348
)

// Azioni dei tasti
const (
	Press   = 1
	Release = 0
	Repeat  = 2
)

// Modificatori dei tasti
const (
	ModShift   = 0x0001
	ModControl = 0x0002
	ModAlt     = 0x0004
	ModSuper   = 0x0008
)

// Pulsanti del mouse
const (
	MouseButton1 = 0
	MouseButton2 = 1
	MouseButton3 = 2
	MouseButton4 = 3
	MouseButton5 = 4
	MouseButton6 = 5
	MouseButton7 = 6
	MouseButton8 = 7

	// Alias comuni
	MouseButtonLeft   = MouseButton1
	MouseButtonRight  = MouseButton2
	MouseButtonMiddle = MouseButton3
)
