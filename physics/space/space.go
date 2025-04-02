// Package space fornisce strutture spaziali per ottimizzare le query spaziali
package space

import (
	"math"
	"sync"

	"github.com/alexanderi96/go-space-engine/core/constants"
	"github.com/alexanderi96/go-space-engine/core/vector"
	"github.com/alexanderi96/go-space-engine/physics/body"
)

// TaskSubmitter rappresenta un'interfaccia per sottomettere task da eseguire in parallelo
type TaskSubmitter interface {
	// Submit sottomette una task da eseguire
	Submit(task func())
	// Wait attende che tutte le task siano completate
	Wait()
}

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
	// UpdateAll aggiorna la posizione di più corpi nella struttura
	UpdateAll(bodies []body.Body, taskSubmitter TaskSubmitter)
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

	// Campi per il calcolo della gravità
	totalMass    float64        // Massa totale di tutti i corpi in questo nodo e nei suoi figli
	centerOfMass vector.Vector3 // Centro di massa di tutti i corpi in questo nodo e nei suoi figli

	// Mutex per proteggere l'accesso concorrente
	mutex sync.RWMutex
}

// NewOctree crea un nuovo octree
func NewOctree(bounds *AABB, maxObjects, maxLevels int) *Octree {
	return &Octree{
		bounds:       bounds,
		maxObjects:   maxObjects,
		maxLevels:    maxLevels,
		level:        0,
		objects:      make([]body.Body, 0),
		divided:      false,
		totalMass:    0,
		centerOfMass: vector.Zero3(),
		mutex:        sync.RWMutex{},
	}
}

// Insert inserisce un corpo nell'octree
func (ot *Octree) Insert(b body.Body) {
	ot.mutex.Lock()
	defer ot.mutex.Unlock()

	ot.insertUnsafe(b)
}

// insertUnsafe inserisce un corpo nell'octree senza bloccare il mutex
func (ot *Octree) insertUnsafe(b body.Body) {
	// Se l'octree è già diviso, inserisci nei figli appropriati
	if ot.divided {
		indices := ot.getIndices(b)
		for _, index := range indices {
			if index != -1 {
				ot.children[index].insertUnsafe(b)
			}
		}

		// Aggiorna il centro di massa e la massa totale
		ot.updateMassAndCenterOfMass(b, true)
		return
	}

	// Aggiungi l'oggetto a questo nodo
	ot.objects = append(ot.objects, b)

	// Aggiorna il centro di massa e la massa totale
	ot.updateMassAndCenterOfMass(b, true)

	// Verifica se è necessario dividere l'octree
	if len(ot.objects) > ot.maxObjects && ot.level < ot.maxLevels {
		// Dividi l'octree
		ot.split()

		// Ridistribuisci gli oggetti nei figli
		for i := 0; i < len(ot.objects); i++ {
			indices := ot.getIndices(ot.objects[i])
			for _, index := range indices {
				if index != -1 {
					ot.children[index].insertUnsafe(ot.objects[i])
				}
			}
		}

		// Svuota gli oggetti di questo nodo
		ot.objects = make([]body.Body, 0)
	}
}

// Remove rimuove un corpo dall'octree
func (ot *Octree) Remove(b body.Body) {
	ot.mutex.Lock()
	defer ot.mutex.Unlock()

	ot.removeUnsafe(b)
}

// removeUnsafe rimuove un corpo dall'octree senza bloccare il mutex
func (ot *Octree) removeUnsafe(b body.Body) {
	// Se l'octree è diviso, rimuovi dai figli appropriati
	if ot.divided {
		indices := ot.getIndices(b)
		for _, index := range indices {
			if index != -1 {
				ot.children[index].removeUnsafe(b)
			}
		}

		// Aggiorna il centro di massa e la massa totale
		ot.updateMassAndCenterOfMass(b, false)
		return
	}

	// Rimuovi l'oggetto da questo nodo
	for i, obj := range ot.objects {
		if obj.ID() == b.ID() {
			// Rimuovi l'oggetto scambiandolo con l'ultimo e troncando la slice
			lastIndex := len(ot.objects) - 1
			ot.objects[i] = ot.objects[lastIndex]
			ot.objects = ot.objects[:lastIndex]

			// Aggiorna il centro di massa e la massa totale
			ot.updateMassAndCenterOfMass(b, false)
			break
		}
	}
}

// UpdateAll aggiorna la posizione di più corpi nell'octree
func (ot *Octree) UpdateAll(bodies []body.Body, taskSubmitter TaskSubmitter) {
	for _, b := range bodies {
		b := b // Cattura la variabile per la goroutine
		taskSubmitter.Submit(func() {
			ot.Update(b)
		})
	}
	taskSubmitter.Wait()
}

// Update aggiorna la posizione di un corpo nell'octree
func (ot *Octree) Update(b body.Body) {
	ot.mutex.Lock()
	defer ot.mutex.Unlock()

	// Rimuovi e reinserisci l'oggetto
	ot.removeUnsafe(b)
	ot.insertUnsafe(b)
}

// Query restituisce tutti i corpi che potrebbero interagire con la regione specificata
func (ot *Octree) Query(region Region) []body.Body {
	ot.mutex.RLock()
	defer ot.mutex.RUnlock()

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
	ot.mutex.RLock()
	defer ot.mutex.RUnlock()

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
	ot.mutex.Lock()
	defer ot.mutex.Unlock()

	ot.objects = make([]body.Body, 0)

	// Resetta il centro di massa e la massa totale
	ot.totalMass = 0
	ot.centerOfMass = vector.Zero3()

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

// updateMassAndCenterOfMass aggiorna il centro di massa e la massa totale
func (ot *Octree) updateMassAndCenterOfMass(b body.Body, adding bool) {
	mass := b.Mass().Value()
	position := b.Position()

	if adding {
		// Aggiunge il corpo
		oldTotalMass := ot.totalMass
		ot.totalMass += mass

		if oldTotalMass > 0 {
			// Aggiorna il centro di massa
			ot.centerOfMass = ot.centerOfMass.Scale(oldTotalMass).Add(position.Scale(mass)).Scale(1.0 / ot.totalMass)
		} else {
			// Se è il primo corpo, il centro di massa è la sua posizione
			ot.centerOfMass = position
		}
	} else {
		// Rimuove il corpo
		if ot.totalMass > mass {
			// Aggiorna il centro di massa
			oldTotalMass := ot.totalMass
			ot.totalMass -= mass
			ot.centerOfMass = ot.centerOfMass.Scale(oldTotalMass).Sub(position.Scale(mass)).Scale(1.0 / ot.totalMass)
		} else {
			// Se era l'ultimo corpo, resetta il centro di massa
			ot.totalMass = 0
			ot.centerOfMass = vector.Zero3()
		}
	}

	// Se l'octree è diviso, propaga l'aggiornamento ai figli
	if ot.divided {
		// Ricalcola il centro di massa dai figli
		ot.totalMass = 0
		weightedPosition := vector.Zero3()

		for i := 0; i < 8; i++ {
			if ot.children[i] != nil {
				childMass := ot.children[i].totalMass
				ot.totalMass += childMass
				if childMass > 0 {
					weightedPosition = weightedPosition.Add(ot.children[i].centerOfMass.Scale(childMass))
				}
			}
		}

		if ot.totalMass > 0 {
			ot.centerOfMass = weightedPosition.Scale(1.0 / ot.totalMass)
		}
	}
}

// CalculateGravity calcola la forza gravitazionale su un corpo utilizzando l'algoritmo Barnes-Hut
func (ot *Octree) CalculateGravity(b body.Body, theta float64) vector.Vector3 {
	ot.mutex.RLock()
	defer ot.mutex.RUnlock()

	force := vector.Zero3()
	ot.calculateGravityRecursive(b, theta, &force)
	return force
}

// calculateGravityRecursive calcola ricorsivamente la forza gravitazionale
func (ot *Octree) calculateGravityRecursive(b body.Body, theta float64, force *vector.Vector3) {
	// Se l'octree non è diviso o non ha corpi, calcola la forza direttamente
	if !ot.divided || ot.totalMass == 0 {
		ot.calculateLeafNodeGravity(b, force)
		return
	}

	// Calcola la larghezza del nodo e la distanza dal corpo al centro di massa
	width := ot.bounds.Max.X() - ot.bounds.Min.X()
	deltaPos := ot.centerOfMass.Sub(b.Position())
	distanceSquared := deltaPos.LengthSquared()

	// Evita divisione per zero
	if distanceSquared < 1e-10 {
		return
	}

	// Se il rapporto larghezza/distanza è inferiore a theta, approssima con il centro di massa
	if (width * width) < (theta * theta * distanceSquared) {
		ot.approximateGravityWithCenterOfMass(b, force)
		return
	}

	// Altrimenti, calcola ricorsivamente per ogni figlio
	for i := 0; i < 8; i++ {
		if ot.children[i] != nil && ot.children[i].totalMass > 0 {
			ot.children[i].calculateGravityRecursive(b, theta, force)
		}
	}
}

// calculateLeafNodeGravity calcola la forza gravitazionale per ogni corpo nel nodo foglia
func (ot *Octree) calculateLeafNodeGravity(b body.Body, force *vector.Vector3) {

	// Massa del corpo
	bodyMass := b.Mass().Value()
	bodyPos := b.Position()

	// Calcola la forza per ogni corpo nel nodo
	for _, obj := range ot.objects {
		// Evita di calcolare la forza su se stesso
		if obj.ID() == b.ID() {
			continue
		}

		// Calcola il vettore direzione
		deltaPos := obj.Position().Sub(bodyPos)
		distanceSquared := deltaPos.LengthSquared()

		// Evita divisione per zero
		if distanceSquared <= 1e-10 {
			continue
		}

		// Calcola la forza gravitazionale
		distance := math.Sqrt(distanceSquared)
		direction := deltaPos.Scale(1.0 / distance)

		// F = G * m1 * m2 / r^2
		forceMagnitude := constants.G * bodyMass * obj.Mass().Value() / distanceSquared

		// Aggiungi la forza al vettore forza totale
		forceVector := *force
		*force = forceVector.Add(direction.Scale(forceMagnitude))
	}
}

// approximateGravityWithCenterOfMass approssima la forza gravitazionale usando il centro di massa
func (ot *Octree) approximateGravityWithCenterOfMass(b body.Body, force *vector.Vector3) {

	// Massa del corpo
	bodyMass := b.Mass().Value()
	bodyPos := b.Position()

	// Calcola il vettore direzione
	deltaPos := ot.centerOfMass.Sub(bodyPos)
	distanceSquared := deltaPos.LengthSquared()

	// Evita divisione per zero
	if distanceSquared <= 1e-10 {
		return
	}

	// Calcola la forza gravitazionale
	distance := math.Sqrt(distanceSquared)
	direction := deltaPos.Scale(1.0 / distance)

	// F = G * m1 * m2 / r^2
	forceMagnitude := constants.G * bodyMass * ot.totalMass / distanceSquared

	// Aggiungi la forza al vettore forza totale
	forceVector := *force
	*force = forceVector.Add(direction.Scale(forceMagnitude))
}
