// Package space fornisce strutture spaziali per ottimizzare le query spaziali
package space

import (
	"math"

	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/physics/body"
)

// Region rappresenta una regione dello spazio
type Region interface {
	// Contains verifica se un punto è contenuto nella regione
	Contains(point vector.Vector3) bool
	// ContainsSphere verifica se una sfera è contenuta nella regione
	ContainsSphere(center vector.Vector3, radius float64) bool
	// Intersects verifica se la regione interseca un'altra regione
	Intersects(other Region) bool
}

// AABB rappresenta un Axis-Aligned Bounding Box (scatola di delimitazione allineata agli assi)
type AABB struct {
	Min vector.Vector3 // Punto minimo (angolo inferiore sinistro posteriore)
	Max vector.Vector3 // Punto massimo (angolo superiore destro anteriore)
}

// NewAABB crea un nuovo AABB
func NewAABB(min, max vector.Vector3) *AABB {
	return &AABB{
		Min: min,
		Max: max,
	}
}

// Contains verifica se un punto è contenuto nell'AABB
func (aabb *AABB) Contains(point vector.Vector3) bool {
	return point.X() >= aabb.Min.X() && point.X() <= aabb.Max.X() &&
		point.Y() >= aabb.Min.Y() && point.Y() <= aabb.Max.Y() &&
		point.Z() >= aabb.Min.Z() && point.Z() <= aabb.Max.Z()
}

// ContainsSphere verifica se una sfera è contenuta nell'AABB
func (aabb *AABB) ContainsSphere(center vector.Vector3, radius float64) bool {
	// Calcola la distanza al quadrato tra il centro della sfera e il punto più vicino dell'AABB
	closestX := math.Max(aabb.Min.X(), math.Min(center.X(), aabb.Max.X()))
	closestY := math.Max(aabb.Min.Y(), math.Min(center.Y(), aabb.Max.Y()))
	closestZ := math.Max(aabb.Min.Z(), math.Min(center.Z(), aabb.Max.Z()))

	distanceSquared := (closestX-center.X())*(closestX-center.X()) +
		(closestY-center.Y())*(closestY-center.Y()) +
		(closestZ-center.Z())*(closestZ-center.Z())

	// La sfera è contenuta se la distanza al quadrato è minore o uguale al raggio al quadrato
	return distanceSquared <= radius*radius
}

// Intersects verifica se l'AABB interseca un altro AABB
func (aabb *AABB) Intersects(other Region) bool {
	otherAABB, ok := other.(*AABB)
	if !ok {
		// Se l'altra regione non è un AABB, usa un'implementazione generica
		return false
	}

	// Due AABB si intersecano se si sovrappongono in tutte e tre le dimensioni
	return aabb.Min.X() <= otherAABB.Max.X() && aabb.Max.X() >= otherAABB.Min.X() &&
		aabb.Min.Y() <= otherAABB.Max.Y() && aabb.Max.Y() >= otherAABB.Min.Y() &&
		aabb.Min.Z() <= otherAABB.Max.Z() && aabb.Max.Z() >= otherAABB.Min.Z()
}

// Center restituisce il centro dell'AABB
func (aabb *AABB) Center() vector.Vector3 {
	return aabb.Min.Add(aabb.Max).Scale(0.5)
}

// Size restituisce le dimensioni dell'AABB
func (aabb *AABB) Size() vector.Vector3 {
	return aabb.Max.Sub(aabb.Min)
}

// SpatialStructure rappresenta una struttura spaziale per ottimizzare le query spaziali
type SpatialStructure interface {
	// Insert inserisce un corpo nella struttura
	Insert(b body.Body)
	// Remove rimuove un corpo dalla struttura
	Remove(b body.Body)
	// Update aggiorna la posizione di un corpo nella struttura
	Update(b body.Body)
	// Query restituisce tutti i corpi che potrebbero interagire con la regione specificata
	Query(region Region) []body.Body
	// QuerySphere restituisce tutti i corpi che potrebbero interagire con la sfera specificata
	QuerySphere(center vector.Vector3, radius float64) []body.Body
	// Clear rimuove tutti i corpi dalla struttura
	Clear()
}

// Octree implementa una struttura spaziale ottimizzata basata su un octree
type Octree struct {
	bounds     *AABB       // Limiti dell'octree
	maxObjects int         // Numero massimo di oggetti per nodo
	maxLevels  int         // Numero massimo di livelli
	level      int         // Livello corrente
	objects    []body.Body // Oggetti in questo nodo
	children   [8]*Octree  // Figli dell'octree
	divided    bool        // Indica se l'octree è stato diviso
}

// NewOctree crea un nuovo octree
func NewOctree(bounds *AABB, maxObjects, maxLevels int) *Octree {
	return &Octree{
		bounds:     bounds,
		maxObjects: maxObjects,
		maxLevels:  maxLevels,
		level:      0,
		objects:    make([]body.Body, 0),
		divided:    false,
	}
}

// Insert inserisce un corpo nell'octree
func (ot *Octree) Insert(b body.Body) {
	// Se l'octree è già diviso, inserisci nei figli appropriati
	if ot.divided {
		indices := ot.getIndices(b)
		for _, index := range indices {
			if index != -1 {
				ot.children[index].Insert(b)
			}
		}
		return
	}

	// Aggiungi l'oggetto a questo nodo
	ot.objects = append(ot.objects, b)

	// Verifica se è necessario dividere l'octree
	if len(ot.objects) > ot.maxObjects && ot.level < ot.maxLevels {
		// Dividi l'octree
		ot.split()

		// Ridistribuisci gli oggetti nei figli
		for i := 0; i < len(ot.objects); i++ {
			indices := ot.getIndices(ot.objects[i])
			for _, index := range indices {
				if index != -1 {
					ot.children[index].Insert(ot.objects[i])
				}
			}
		}

		// Svuota gli oggetti di questo nodo
		ot.objects = make([]body.Body, 0)
	}
}

// Remove rimuove un corpo dall'octree
func (ot *Octree) Remove(b body.Body) {
	// Se l'octree è diviso, rimuovi dai figli appropriati
	if ot.divided {
		indices := ot.getIndices(b)
		for _, index := range indices {
			if index != -1 {
				ot.children[index].Remove(b)
			}
		}
		return
	}

	// Rimuovi l'oggetto da questo nodo
	for i, obj := range ot.objects {
		if obj.ID() == b.ID() {
			// Rimuovi l'oggetto scambiandolo con l'ultimo e troncando la slice
			lastIndex := len(ot.objects) - 1
			ot.objects[i] = ot.objects[lastIndex]
			ot.objects = ot.objects[:lastIndex]
			break
		}
	}
}

// Update aggiorna la posizione di un corpo nell'octree
func (ot *Octree) Update(b body.Body) {
	// Rimuovi e reinserisci l'oggetto
	ot.Remove(b)
	ot.Insert(b)
}

// Query restituisce tutti i corpi che potrebbero interagire con la regione specificata
func (ot *Octree) Query(region Region) []body.Body {
	result := make([]body.Body, 0)

	// Verifica se la regione interseca questo nodo
	if !region.Intersects(ot.bounds) {
		return result
	}

	// Aggiungi gli oggetti di questo nodo
	result = append(result, ot.objects...)

	// Se l'octree è diviso, query i figli
	if ot.divided {
		for i := 0; i < 8; i++ {
			childResult := ot.children[i].Query(region)
			result = append(result, childResult...)
		}
	}

	return result
}

// QuerySphere restituisce tutti i corpi che potrebbero interagire con la sfera specificata
func (ot *Octree) QuerySphere(center vector.Vector3, radius float64) []body.Body {
	result := make([]body.Body, 0)

	// Verifica se la sfera interseca questo nodo
	if !ot.bounds.ContainsSphere(center, radius) {
		return result
	}

	// Aggiungi gli oggetti di questo nodo
	result = append(result, ot.objects...)

	// Se l'octree è diviso, query i figli
	if ot.divided {
		for i := 0; i < 8; i++ {
			childResult := ot.children[i].QuerySphere(center, radius)
			result = append(result, childResult...)
		}
	}

	return result
}

// Clear rimuove tutti i corpi dall'octree
func (ot *Octree) Clear() {
	ot.objects = make([]body.Body, 0)

	if ot.divided {
		for i := 0; i < 8; i++ {
			ot.children[i].Clear()
			ot.children[i] = nil
		}
		ot.divided = false
	}
}

// split divide l'octree in otto figli
func (ot *Octree) split() {
	// Calcola il centro dell'octree
	center := ot.bounds.Center()

	// Crea gli otto figli
	// Ordine: [0] = Bottom Left Back, [1] = Bottom Right Back, [2] = Bottom Right Front, [3] = Bottom Left Front,
	//         [4] = Top Left Back, [5] = Top Right Back, [6] = Top Right Front, [7] = Top Left Front
	childBounds := [8]*AABB{
		// [0] Bottom Left Back
		NewAABB(
			ot.bounds.Min,
			center,
		),
		// [1] Bottom Right Back
		NewAABB(
			vector.NewVector3(center.X(), ot.bounds.Min.Y(), ot.bounds.Min.Z()),
			vector.NewVector3(ot.bounds.Max.X(), center.Y(), center.Z()),
		),
		// [2] Bottom Right Front
		NewAABB(
			vector.NewVector3(center.X(), ot.bounds.Min.Y(), center.Z()),
			vector.NewVector3(ot.bounds.Max.X(), center.Y(), ot.bounds.Max.Z()),
		),
		// [3] Bottom Left Front
		NewAABB(
			vector.NewVector3(ot.bounds.Min.X(), ot.bounds.Min.Y(), center.Z()),
			vector.NewVector3(center.X(), center.Y(), ot.bounds.Max.Z()),
		),
		// [4] Top Left Back
		NewAABB(
			vector.NewVector3(ot.bounds.Min.X(), center.Y(), ot.bounds.Min.Z()),
			vector.NewVector3(center.X(), ot.bounds.Max.Y(), center.Z()),
		),
		// [5] Top Right Back
		NewAABB(
			vector.NewVector3(center.X(), center.Y(), ot.bounds.Min.Z()),
			vector.NewVector3(ot.bounds.Max.X(), ot.bounds.Max.Y(), center.Z()),
		),
		// [6] Top Right Front
		NewAABB(
			center,
			ot.bounds.Max,
		),
		// [7] Top Left Front
		NewAABB(
			vector.NewVector3(ot.bounds.Min.X(), center.Y(), center.Z()),
			vector.NewVector3(center.X(), ot.bounds.Max.Y(), ot.bounds.Max.Z()),
		),
	}

	// Crea gli octree figli
	for i := 0; i < 8; i++ {
		ot.children[i] = NewOctree(childBounds[i], ot.maxObjects, ot.maxLevels)
		ot.children[i].level = ot.level + 1
	}

	ot.divided = true
}

// getIndices determina in quali figli un corpo dovrebbe essere inserito
func (ot *Octree) getIndices(b body.Body) []int {
	result := make([]int, 0, 8)
	center := ot.bounds.Center()
	position := b.Position()
	radius := b.Radius().Value()

	// Determina in quali ottanti il corpo si trova
	top := position.Y()+radius > center.Y()
	bottom := position.Y()-radius < center.Y()
	left := position.X()-radius < center.X()
	right := position.X()+radius > center.X()
	front := position.Z()+radius > center.Z()
	back := position.Z()-radius < center.Z()

	// Bottom Left Back
	if bottom && left && back {
		result = append(result, 0)
	}

	// Bottom Right Back
	if bottom && right && back {
		result = append(result, 1)
	}

	// Bottom Right Front
	if bottom && right && front {
		result = append(result, 2)
	}

	// Bottom Left Front
	if bottom && left && front {
		result = append(result, 3)
	}

	// Top Left Back
	if top && left && back {
		result = append(result, 4)
	}

	// Top Right Back
	if top && right && back {
		result = append(result, 5)
	}

	// Top Right Front
	if top && right && front {
		result = append(result, 6)
	}

	// Top Left Front
	if top && left && front {
		result = append(result, 7)
	}

	if len(result) == 0 {
		result = append(result, -1)
	}

	return result
}
