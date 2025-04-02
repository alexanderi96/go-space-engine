// Package collision fornisce interfacce e implementazioni per il rilevamento e la risoluzione delle collisioni
package collision

import (
	"math"

	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/physics/body"
)

// CollisionInfo contiene informazioni sulla collisione
type CollisionInfo struct {
	BodyA       body.Body      // Primo corpo coinvolto nella collisione
	BodyB       body.Body      // Secondo corpo coinvolto nella collisione
	Point       vector.Vector3 // Punto di contatto
	Normal      vector.Vector3 // Normale di collisione (da A a B)
	Depth       float64        // Profondità di penetrazione
	HasCollided bool           // Indica se c'è stata una collisione
}

// Collider rileva le collisioni tra corpi
type Collider interface {
	// CheckCollision verifica se due corpi collidono
	CheckCollision(a, b body.Body) CollisionInfo
}

// CollisionResolver risolve le collisioni
type CollisionResolver interface {
	// ResolveCollision risolve una collisione
	ResolveCollision(info CollisionInfo)
}

// SphereCollider implementa un rilevatore di collisioni per corpi sferici
type SphereCollider struct{}

// NewSphereCollider crea un nuovo rilevatore di collisioni per corpi sferici
func NewSphereCollider() *SphereCollider {
	return &SphereCollider{}
}

// CheckCollision verifica se due corpi sferici collidono
func (sc *SphereCollider) CheckCollision(a, b body.Body) CollisionInfo {
	// Calcola il vettore direzione da a a b
	direction := b.Position().Sub(a.Position())

	// Calcola la distanza al quadrato
	distanceSquared := direction.LengthSquared()

	// Calcola la somma dei raggi
	radiusA := a.Radius().Value()
	radiusB := b.Radius().Value()
	sumRadii := radiusA + radiusB

	// Verifica se c'è una collisione
	hasCollided := distanceSquared < sumRadii*sumRadii

	// Se non c'è collisione, restituisci un'informazione di collisione vuota
	if !hasCollided {
		return CollisionInfo{
			BodyA:       a,
			BodyB:       b,
			HasCollided: false,
		}
	}

	// Calcola la distanza
	distance := math.Sqrt(distanceSquared)

	// Calcola la normale di collisione
	var normal vector.Vector3
	if distance > 1e-10 {
		normal = direction.Scale(1.0 / distance)
	} else {
		// Se i corpi sono sovrapposti completamente, usa una normale predefinita
		normal = vector.NewVector3(1, 0, 0)
	}

	// Calcola la profondità di penetrazione
	depth := sumRadii - distance

	// Calcola il punto di contatto
	point := a.Position().Add(normal.Scale(radiusA))

	return CollisionInfo{
		BodyA:       a,
		BodyB:       b,
		Point:       point,
		Normal:      normal,
		Depth:       depth,
		HasCollided: true,
	}
}

// ImpulseResolver implementa un risolutore di collisioni basato sugli impulsi
type ImpulseResolver struct {
	restitution float64 // Coefficiente di restituzione (elasticità)
}

// NewImpulseResolver crea un nuovo risolutore di collisioni basato sugli impulsi
func NewImpulseResolver(restitution float64) *ImpulseResolver {
	return &ImpulseResolver{
		restitution: restitution,
	}
}

// ResolveCollision risolve una collisione usando il metodo degli impulsi
func (ir *ImpulseResolver) ResolveCollision(info CollisionInfo) {
	// Se non c'è stata una collisione, non fare nulla
	if !info.HasCollided {
		return
	}

	a := info.BodyA
	b := info.BodyB

	// Se entrambi i corpi sono statici, non fare nulla
	if a.IsStatic() && b.IsStatic() {
		return
	}

	// Calcola la velocità relativa
	relativeVelocity := b.Velocity().Sub(a.Velocity())

	// Calcola la velocità relativa lungo la normale
	velocityAlongNormal := relativeVelocity.Dot(info.Normal)

	// Se i corpi si stanno allontanando, non fare nulla
	if velocityAlongNormal > 0 {
		return
	}

	// Calcola il coefficiente di restituzione (elasticità)
	// Usa il minimo tra il coefficiente di restituzione del risolutore e quello dei materiali
	elasticityA := a.Material().Elasticity()
	elasticityB := b.Material().Elasticity()
	restitution := math.Min(ir.restitution, math.Min(elasticityA, elasticityB))

	// Calcola l'impulso scalare
	// j = -(1 + e) * velocityAlongNormal / (1/massA + 1/massB)
	massA := a.Mass().Value()
	massB := b.Mass().Value()

	// Gestisci corpi statici (massa infinita)
	var inverseMassA, inverseMassB float64
	if a.IsStatic() {
		inverseMassA = 0
	} else {
		inverseMassA = 1.0 / massA
	}

	if b.IsStatic() {
		inverseMassB = 0
	} else {
		inverseMassB = 1.0 / massB
	}

	// Calcola l'impulso scalare
	j := -(1.0 + restitution) * velocityAlongNormal
	j /= inverseMassA + inverseMassB

	// Applica l'impulso
	impulse := info.Normal.Scale(j)

	if !a.IsStatic() {
		a.SetVelocity(a.Velocity().Sub(impulse.Scale(inverseMassA)))
	}

	if !b.IsStatic() {
		b.SetVelocity(b.Velocity().Add(impulse.Scale(inverseMassB)))
	}

	// Correggi la penetrazione (risoluzione della posizione)
	ir.resolvePosition(info)
}

// resolvePosition corregge la penetrazione tra i corpi
func (ir *ImpulseResolver) resolvePosition(info CollisionInfo) {
	a := info.BodyA
	b := info.BodyB

	// Se entrambi i corpi sono statici, non fare nulla
	if a.IsStatic() && b.IsStatic() {
		return
	}

	// Costante di correzione della posizione (0.2 - 0.8)
	const percent = 0.2
	// Soglia minima di penetrazione
	const slop = 0.01

	// Calcola le masse inverse
	massA := a.Mass().Value()
	massB := b.Mass().Value()
	var inverseMassA, inverseMassB float64
	if a.IsStatic() {
		inverseMassA = 0
	} else {
		inverseMassA = 1.0 / massA
	}

	if b.IsStatic() {
		inverseMassB = 0
	} else {
		inverseMassB = 1.0 / massB
	}

	// Calcola la somma delle masse inverse
	totalInverseMass := inverseMassA + inverseMassB
	if totalInverseMass <= 0 {
		return // Entrambi i corpi hanno massa infinita
	}

	// Calcola la correzione
	correction := math.Max(info.Depth-slop, 0.0) * percent

	// Applica la correzione
	if !a.IsStatic() {
		a.SetPosition(a.Position().Sub(info.Normal.Scale(correction * inverseMassA / totalInverseMass)))
	}

	if !b.IsStatic() {
		b.SetPosition(b.Position().Add(info.Normal.Scale(correction * inverseMassB / totalInverseMass)))
	}
}
