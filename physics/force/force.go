// Package force fornisce interfacce e implementazioni per le forze fisiche
package force

import (
	"math"

	"github.com/alexanderi96/go-space-engine/core/constants"
	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/physics/body"
)

// Force rappresenta una forza fisica
type Force interface {
	// Apply applica la forza a un corpo e restituisce il vettore forza
	Apply(b body.Body) vector.Vector3

	// ApplyBetween applica la forza tra due corpi e restituisce i vettori forza per ciascun corpo
	ApplyBetween(a, b body.Body) (vector.Vector3, vector.Vector3)

	// IsGlobal restituisce true se la forza è globale (applicata a tutti i corpi)
	IsGlobal() bool
}

// GravitationalForce implementa la forza gravitazionale
type GravitationalForce struct {
	G     float64 // Costante gravitazionale
	Theta float64 // Parametro di approssimazione per l'algoritmo Barnes-Hut
}

// NewGravitationalForce crea una nuova forza gravitazionale
func NewGravitationalForce() *GravitationalForce {
	return &GravitationalForce{
		G:     constants.G,
		Theta: 0.5, // Valore predefinito che bilancia precisione ed efficienza
	}
}

// Apply applica la forza gravitazionale a un corpo (non fa nulla per un singolo corpo)
func (gf *GravitationalForce) Apply(b body.Body) vector.Vector3 {
	// La gravità richiede due corpi per essere applicata
	return vector.Zero3()
}

// ApplyBetween applica la forza gravitazionale tra due corpi
func (gf *GravitationalForce) ApplyBetween(a, b body.Body) (vector.Vector3, vector.Vector3) {
	// Calcola il vettore direzione da a a b
	direction := b.Position().Sub(a.Position())

	// Calcola la distanza al quadrato
	distanceSquared := direction.LengthSquared()

	// Evita divisione per zero o forze troppo grandi
	if distanceSquared < 1e-10 {
		return vector.Zero3(), vector.Zero3()
	}

	// Normalizza la direzione
	distance := math.Sqrt(distanceSquared)
	normalizedDirection := direction.Scale(1.0 / distance)

	// Calcola la forza secondo la legge di gravitazione universale
	// F = G * m1 * m2 / r^2
	massA := a.Mass().Value()
	massB := b.Mass().Value()
	forceMagnitude := gf.G * massA * massB / distanceSquared

	// Calcola i vettori forza (direzioni opposte)
	forceOnA := normalizedDirection.Scale(forceMagnitude)
	forceOnB := normalizedDirection.Scale(-forceMagnitude)

	return forceOnA, forceOnB
}

// IsGlobal restituisce true perché la gravità è una forza globale
func (gf *GravitationalForce) IsGlobal() bool {
	return true
}

// SetTheta imposta il parametro di approssimazione per l'algoritmo Barnes-Hut
func (gf *GravitationalForce) SetTheta(theta float64) {
	gf.Theta = theta
}

// GetTheta restituisce il parametro di approssimazione per l'algoritmo Barnes-Hut
func (gf *GravitationalForce) GetTheta() float64 {
	return gf.Theta
}

// ConstantForce implementa una forza costante
type ConstantForce struct {
	force vector.Vector3
}

// NewConstantForce crea una nuova forza costante
func NewConstantForce(force vector.Vector3) *ConstantForce {
	return &ConstantForce{
		force: force,
	}
}

// Apply applica la forza costante a un corpo
func (cf *ConstantForce) Apply(b body.Body) vector.Vector3 {
	return cf.force
}

// ApplyBetween applica la forza costante tra due corpi (applica solo al primo)
func (cf *ConstantForce) ApplyBetween(a, b body.Body) (vector.Vector3, vector.Vector3) {
	return cf.force, vector.Zero3()
}

// IsGlobal restituisce false perché la forza costante non è globale
func (cf *ConstantForce) IsGlobal() bool {
	return false
}

// SpringForce implementa una forza elastica
type SpringForce struct {
	stiffness  float64 // Costante elastica
	restLength float64 // Lunghezza a riposo
	damping    float64 // Coefficiente di smorzamento
}

// NewSpringForce crea una nuova forza elastica
func NewSpringForce(stiffness, restLength, damping float64) *SpringForce {
	return &SpringForce{
		stiffness:  stiffness,
		restLength: restLength,
		damping:    damping,
	}
}

// Apply applica la forza elastica a un corpo (non fa nulla per un singolo corpo)
func (sf *SpringForce) Apply(b body.Body) vector.Vector3 {
	// La forza elastica richiede due corpi per essere applicata
	return vector.Zero3()
}

// ApplyBetween applica la forza elastica tra due corpi
func (sf *SpringForce) ApplyBetween(a, b body.Body) (vector.Vector3, vector.Vector3) {
	// Calcola il vettore direzione da a a b
	direction := b.Position().Sub(a.Position())

	// Calcola la distanza
	distance := direction.Length()

	// Evita divisione per zero
	if distance < 1e-10 {
		return vector.Zero3(), vector.Zero3()
	}

	// Normalizza la direzione
	normalizedDirection := direction.Scale(1.0 / distance)

	// Calcola l'estensione della molla (distanza - lunghezza a riposo)
	extension := distance - sf.restLength

	// Calcola la forza elastica secondo la legge di Hooke
	// F = -k * x
	springForceMagnitude := sf.stiffness * extension

	// Calcola la velocità relativa lungo la direzione della molla
	relativeVelocity := b.Velocity().Sub(a.Velocity())
	relativeVelocityAlongSpring := relativeVelocity.Dot(normalizedDirection)

	// Calcola la forza di smorzamento
	// F_damping = -c * v
	dampingForceMagnitude := sf.damping * relativeVelocityAlongSpring

	// Calcola la forza totale
	totalForceMagnitude := springForceMagnitude + dampingForceMagnitude

	// Calcola i vettori forza (direzioni opposte)
	forceOnA := normalizedDirection.Scale(totalForceMagnitude)
	forceOnB := normalizedDirection.Scale(-totalForceMagnitude)

	return forceOnA, forceOnB
}

// IsGlobal restituisce false perché la forza elastica non è globale
func (sf *SpringForce) IsGlobal() bool {
	return false
}

// DragForce implementa una forza di resistenza
type DragForce struct {
	coefficient float64 // Coefficiente di resistenza
}

// NewDragForce crea una nuova forza di resistenza
func NewDragForce(coefficient float64) *DragForce {
	return &DragForce{
		coefficient: coefficient,
	}
}

// Apply applica la forza di resistenza a un corpo
func (df *DragForce) Apply(b body.Body) vector.Vector3 {
	// La forza di resistenza è proporzionale al quadrato della velocità
	// F = -c * v^2 * v_hat
	velocity := b.Velocity()
	speed := velocity.Length()

	// Se la velocità è quasi zero, non c'è resistenza
	if speed < 1e-10 {
		return vector.Zero3()
	}

	// Calcola la direzione della velocità
	direction := velocity.Scale(1.0 / speed)

	// Calcola la forza di resistenza (in direzione opposta alla velocità)
	forceMagnitude := df.coefficient * speed * speed
	force := direction.Scale(-forceMagnitude)

	return force
}

// ApplyBetween applica la forza di resistenza tra due corpi (applica a entrambi separatamente)
func (df *DragForce) ApplyBetween(a, b body.Body) (vector.Vector3, vector.Vector3) {
	return df.Apply(a), df.Apply(b)
}

// IsGlobal restituisce true perché la forza di resistenza è globale
func (df *DragForce) IsGlobal() bool {
	return true
}
